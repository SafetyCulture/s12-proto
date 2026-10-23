// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"github.com/bufbuild/protoplugin"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/kotlin"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata/kotlin/golden")

// compileOnDisk compiles proto files from the given include directories, plus the standard imports.
func compileOnDisk(t *testing.T, includes []string, files ...string) linker.Files {
	t.Helper()
	compiler := protocompile.Compiler{
		Resolver:       protocompile.WithStandardImports(&protocompile.SourceResolver{ImportPaths: includes}),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	compiled, err := compiler.Compile(context.Background(), files...)
	require.NoError(t, err)
	return compiled
}

// pluginResponseOnDisk runs the plugin handler in-process on the given files and returns the response, or the error that the plugin reports.
func pluginResponseOnDisk(t *testing.T, includes []string, parameter string, files ...string) (*pluginpb.CodeGeneratorResponse, error) {
	t.Helper()
	compiled := compileOnDisk(t, includes, files...)

	seen := map[string]bool{}
	var protoFiles []*descriptorpb.FileDescriptorProto
	var add func(fd protoreflect.FileDescriptor)
	add = func(fd protoreflect.FileDescriptor) {
		if seen[fd.Path()] {
			return
		}
		seen[fd.Path()] = true
		for i := 0; i < fd.Imports().Len(); i++ {
			add(fd.Imports().Get(i).FileDescriptor)
		}
		protoFiles = append(protoFiles, protodesc.ToFileDescriptorProto(fd))
	}
	for _, fd := range compiled {
		add(fd)
	}

	req := &pluginpb.CodeGeneratorRequest{FileToGenerate: files, ProtoFile: protoFiles, Parameter: proto.String(parameter)}
	reqBytes, err := proto.Marshal(req)
	require.NoError(t, err)

	var stdout bytes.Buffer
	var stderr strings.Builder
	env := protoplugin.Env{Args: []string{}, Stdin: bytes.NewReader(reqBytes), Stdout: &stdout, Stderr: &stderr}
	if err := protoplugin.Run(context.Background(), env, protoplugin.HandlerFunc(handle)); err != nil {
		return nil, err
	}

	resp := &pluginpb.CodeGeneratorResponse{}
	require.NoError(t, proto.Unmarshal(stdout.Bytes(), resp), stderr.String())
	return resp, nil
}

func requirePluginResponseError(t *testing.T, includes []string, parameter, want string, files ...string) {
	t.Helper()
	resp, err := pluginResponseOnDisk(t, includes, parameter, files...)
	require.NoError(t, err)
	require.Contains(t, resp.GetError(), want)
}

// runPluginOnDisk runs the plugin handler in-process on the given files and returns the output files.
func runPluginOnDisk(t *testing.T, includes []string, parameter string, files ...string) map[string]string {
	t.Helper()
	resp, err := pluginResponseOnDisk(t, includes, parameter, files...)
	require.NoError(t, err)
	require.Empty(t, resp.GetError())
	out := map[string]string{}
	for _, f := range resp.GetFile() {
		out[f.GetName()] = f.GetContent()
	}
	return out
}

var kitchenFiles = []string{"s12/kitchen/v1/kitchen.proto", "legacy/legacy.proto", "editions/closed.proto", "editions/utf8.proto", "legacy/strict.proto", "shadow/shadow.proto"}

func kitchenOutput(t *testing.T, parameter string) map[string]string {
	t.Helper()
	return runPluginOnDisk(t, []string{"testdata/kotlin/kitchen"}, parameter, kitchenFiles...)
}

func TestKitchenGolden(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages")
	// TestRuntimeHasNoImports covers the runtime file, which is a copy of its source. The other
	// generated files are not Kotlin messages.
	got := map[string]string{}
	for name, content := range out {
		if strings.HasSuffix(name, ".pb.kt") || strings.HasSuffix(name, ".pro") {
			got[name] = content
		}
	}
	goldenDir := filepath.Join("testdata", "kotlin", "golden")
	if *update {
		require.NoError(t, os.RemoveAll(goldenDir))
		for name, content := range got {
			golden := filepath.Join(goldenDir, name)
			require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
			require.NoError(t, os.WriteFile(golden, []byte(content), 0o644))
		}
		return
	}
	goldens := readTree(t, goldenDir, true, func(string) bool { return true })
	requireSameFiles(t, goldens, got, "golden files", "plugin output; run go test -update to accept")
}

func TestPackageOption(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages=com.example.kt")
	require.Contains(t, out["com/example/kt/kitchen/v1/kitchen.pb.kt"], "\npackage com.example.kt.kitchen.v1\n")
	require.Contains(t, out, "com/example/kt/wire/S12Wire.kt")
	require.Contains(t, out["com/example/kt/wire/S12Wire.kt"], "\npackage com.example.kt.wire\n")
	require.Contains(t, out["com/example/kt/s12kt.pro"], "-keep class com.example.kt.** {\n")

	out = kitchenOutput(t, "kotlin_messages")
	require.Contains(t, out["com/safetyculture/s12kt/kitchen/v1/kitchen.pb.kt"], "\npackage com.safetyculture.s12kt.kitchen.v1\n")
}

func TestRuntimeHasNoImports(t *testing.T) {
	runtime := kitchenOutput(t, "kotlin_messages")["com/safetyculture/s12kt/wire/S12Wire.kt"]
	require.NotEmpty(t, runtime)
	for _, line := range strings.Split(runtime, "\n") {
		require.NotRegexp(t, regexp.MustCompile(`\bimport\b`), line)
	}
}

// paramLine matches one constructor parameter of a generated data class.
var paramLine = regexp.MustCompile("^public val `?(\\w+)`?: (.+?)(?: = .+)?,$")

// constructorTypes returns the declared type of every constructor parameter of the class at the
// given depth, found after offset in the source.
func constructorTypes(t *testing.T, source string, offset int, depth int, md protoreflect.MessageDescriptor) (map[string]string, int) {
	t.Helper()
	name := kotlin.SimpleClassName(md)
	indent := strings.Repeat("    ", depth)
	marker := "\n" + indent + "public data class " + name + "(\n"
	start := strings.Index(source[offset:], marker)
	require.GreaterOrEqual(t, start, 0, "no data class for %s", md.FullName())
	start += offset + len(marker)
	end := strings.Index(source[start:], "\n"+indent+") ")
	require.GreaterOrEqual(t, end, 0, "no end of the constructor of %s", md.FullName())
	types := map[string]string{}
	for _, line := range strings.Split(source[start:start+end], "\n") {
		if strings.HasPrefix(line, indent+"    public val ") {
			m := paramLine.FindStringSubmatch(strings.TrimPrefix(line, indent+"    "))
			require.NotNil(t, m, line)
			types[m[1]] = m[2]
		}
	}
	return types, start
}

// expectedKind describes the rule of the ticket for one field.
func expectedKind(fd protoreflect.FieldDescriptor) string {
	switch {
	case fd.IsMap():
		return "Map"
	case fd.IsList():
		return "List"
	case fd.ContainingOneof() != nil && !fd.ContainingOneof().IsSynthetic():
		return "oneof"
	case fd.Message() != nil && fd.Cardinality() != protoreflect.Required:
		return "nullable"
	case fd.HasPresence() && fd.Cardinality() != protoreflect.Required:
		return "nullable"
	}
	return "plain"
}

func checkNullability(t *testing.T, source string, md protoreflect.MessageDescriptor, offset int, depth int) {
	t.Helper()
	types, start := constructorTypes(t, source, offset, depth, md)
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if fd.ContainingOneof() != nil && !fd.ContainingOneof().IsSynthetic() {
			continue
		}
		got, ok := types[kotlin.PropertyName(fd)]
		require.True(t, ok, "%s.%s is not a constructor parameter", md.FullName(), fd.Name())
		switch expectedKind(fd) {
		case "List":
			require.True(t, strings.HasPrefix(got, "kotlin.collections.List<"), "%s.%s: %s", md.FullName(), fd.Name(), got)
		case "Map":
			require.True(t, strings.HasPrefix(got, "kotlin.collections.Map<"), "%s.%s: %s", md.FullName(), fd.Name(), got)
		case "nullable":
			require.True(t, strings.HasSuffix(got, "?"), "%s.%s must be nullable: %s", md.FullName(), fd.Name(), got)
		case "plain":
			require.False(t, strings.HasSuffix(got, "?"), "%s.%s must not be nullable: %s", md.FullName(), fd.Name(), got)
		}
	}
	for i := 0; i < md.Oneofs().Len(); i++ {
		od := md.Oneofs().Get(i)
		if od.IsSynthetic() {
			continue
		}
		got, ok := types[kotlin.OneofPropertyName(od)]
		require.True(t, ok, "%s.%s is not a constructor parameter", md.FullName(), od.Name())
		require.True(t, strings.HasSuffix(got, "?"), "oneof %s.%s must be nullable: %s", md.FullName(), od.Name(), got)
	}
	nested := md.Messages()
	for i := 0; i < nested.Len(); i++ {
		if nested.Get(i).IsMapEntry() {
			continue
		}
		checkNullability(t, source, nested.Get(i), start, depth+1)
	}
}

func TestKitchenNullability(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages")
	compiled := compileOnDisk(t, []string{"testdata/kotlin/kitchen"}, kitchenFiles...)
	for _, fd := range compiled {
		source := out[kotlin.FileName(fd, "com.safetyculture.s12kt")]
		for i := 0; i < fd.Messages().Len(); i++ {
			checkNullability(t, "\n"+source, fd.Messages().Get(i), 0, 0)
		}
	}
}

// constructorField is the expected declaration of one constructor parameter: its name, its type and its default.
type constructorField struct{ name, typ, def string }

// kitchenConstructors lists, for the messages with presence, required fields or a field default, the
// parameters of the constructor as the ticket rules give them. The tables are literal on purpose:
// they do not use the rule that the generator follows.
var kitchenConstructors = map[string][]constructorField{
	"OptionalScalars": {
		{"fDouble", "kotlin.Double?", "null"},
		{"fFloat", "kotlin.Float?", "null"},
		{"fInt32", "kotlin.Int?", "null"},
		{"fInt64", "kotlin.Long?", "null"},
		{"fUint32", "kotlin.Int?", "null"},
		{"fUint64", "kotlin.Long?", "null"},
		{"fSint32", "kotlin.Int?", "null"},
		{"fSint64", "kotlin.Long?", "null"},
		{"fFixed32", "kotlin.Int?", "null"},
		{"fFixed64", "kotlin.Long?", "null"},
		{"fSfixed32", "kotlin.Int?", "null"},
		{"fSfixed64", "kotlin.Long?", "null"},
		{"fBool", "kotlin.Boolean?", "null"},
		{"fString", "kotlin.String?", "null"},
		{"fBytes", "com.safetyculture.s12kt.wire.ByteString?", "null"},
		{"fEnum", "com.safetyculture.s12kt.kitchen.v1.Color?", "null"},
		{"fMessage", "com.safetyculture.s12kt.kitchen.v1.Scalars?", "null"},
	},
	"Legacy": {
		{"id", "kotlin.String", ""},
		{"count", "kotlin.Int?", "null"},
		{"label", "kotlin.String?", "null"},
		{"level", "com.safetyculture.s12kt.kitchen.legacy.Level?", "null"},
		{"numbers", "kotlin.collections.List<kotlin.Int>", "kotlin.collections.emptyList()"},
		{"packedNumbers", "kotlin.collections.List<kotlin.Int>", "kotlin.collections.emptyList()"},
		{"blob", "com.safetyculture.s12kt.wire.ByteString?", "null"},
		{"parent", "com.safetyculture.s12kt.kitchen.legacy.Legacy?", "null"},
		{"score", "kotlin.Double", ""},
		{"pick", "com.safetyculture.s12kt.kitchen.legacy.Legacy.Pick?", "null"},
	},
	"ReqEnum": {
		{"level", "com.safetyculture.s12kt.kitchen.legacy.Level", ""},
	},
	"ReqHolder": {
		{"one", "com.safetyculture.s12kt.kitchen.legacy.Req?", "null"},
		{"many", "kotlin.collections.List<com.safetyculture.s12kt.kitchen.legacy.Req>", "kotlin.collections.emptyList()"},
		{"byName", "kotlin.collections.Map<kotlin.String, com.safetyculture.s12kt.kitchen.legacy.Req>", "kotlin.collections.emptyMap()"},
		{"tags", "kotlin.collections.Map<kotlin.String, kotlin.String>", "kotlin.collections.emptyMap()"},
	},
}

// constructorFields returns the parameters of the top-level data class name in source, in order,
// without the unknownFields parameter.
func constructorFields(t *testing.T, source, name string) []constructorField {
	t.Helper()
	marker := "\npublic data class " + name + "(\n"
	start := strings.Index(source, marker)
	require.GreaterOrEqual(t, start, 0, "no data class for %s", name)
	start += len(marker)
	end := strings.Index(source[start:], "\n) ")
	require.GreaterOrEqual(t, end, 0, "no end of the constructor of %s", name)
	var fields []constructorField
	for _, line := range strings.Split(source[start:start+end], "\n") {
		m := declaredParam.FindStringSubmatch(strings.TrimSpace(line))
		if m != nil && m[1] != "unknownFields" {
			fields = append(fields, constructorField{m[1], m[2], m[3]})
		}
	}
	return fields
}

// declaredParam matches one constructor parameter of a generated data class: name, type and default.
var declaredParam = regexp.MustCompile("^public val `?(\\w+)`?: (.+?)(?: = (.+))?,$")

// TestKitchenConstructorTables checks the type and the default of every parameter of the messages
// with explicit presence, required fields and field defaults against literal tables.
func TestKitchenConstructorTables(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages")
	for name, want := range kitchenConstructors {
		file := "com/safetyculture/s12kt/kitchen/v1/kitchen.pb.kt"
		if name != "OptionalScalars" {
			file = "com/safetyculture/s12kt/kitchen/legacy/legacy.pb.kt"
		}
		require.Equal(t, want, constructorFields(t, "\n"+out[file], name), name)
	}
}

// assetNullability is the expectation for s12.assets.v1.Asset, field by field.
var assetNullability = map[string]string{
	"id":                        "kotlin.String",
	"code":                      "kotlin.String",
	"type":                      "com.safetyculture.s12kt.assets.v1.AssetType?",
	"fields":                    "kotlin.collections.List<com.safetyculture.s12kt.assets.v1.AssetField>",
	"media":                     "kotlin.collections.List<com.safetyculture.s12kt.common.Media>",
	"createdAt":                 "com.safetyculture.s12kt.google.protobuf.Timestamp?",
	"modifiedAt":                "com.safetyculture.s12kt.google.protobuf.Timestamp?",
	"site":                      "com.safetyculture.s12kt.assets.v1.AssetSite?",
	"profileImage":              "com.safetyculture.s12kt.common.Media?",
	"state":                     "com.safetyculture.s12kt.assets.v1.AssetState",
	"inspectedAt":               "com.safetyculture.s12kt.google.protobuf.Timestamp?",
	"sourcedFromExternalSystem": "kotlin.Boolean",
	"allMedia":                  "kotlin.collections.List<com.safetyculture.s12kt.assets.v1.AssetMediaItem>",
	"openIssuesCount":           "kotlin.Int",
	"statuses":                  "kotlin.collections.List<com.safetyculture.s12kt.assets.v1.AssetStatus>",
	"hierarchyMetadata":         "com.safetyculture.s12kt.assets.v1.AssetHierarchyMetadata?",
	"unknownFields":             "com.safetyculture.s12kt.wire.ByteString",
}

// TestAssetNullability checks the pilot file against the table above. It needs an APISchema
// checkout in APISCHEMA, which this repository does not carry.
func TestAssetNullability(t *testing.T) {
	schema := os.Getenv("APISCHEMA")
	if schema == "" {
		t.Skip("APISCHEMA is not set")
	}
	repo, err := filepath.Abs("../..")
	require.NoError(t, err)
	includes := []string{filepath.Join(schema, "protos"), filepath.Join(schema, "third_party"), repo}
	out := runPluginOnDisk(t, includes, "kotlin_messages", "s12/assets/v1/asset.proto")
	source := "\n" + out["com/safetyculture/s12kt/assets/v1/asset.pb.kt"]
	require.NotEqual(t, "\n", source)

	asset := compileOnDisk(t, includes, "s12/assets/v1/asset.proto")[0].Messages().ByName("Asset")
	types, _ := constructorTypes(t, source, 0, 0, asset)
	require.Equal(t, assetNullability, types)

	checkNullability(t, source, asset, 0, 0)
}

// TestClosedEnum checks that only an enum that the schema marks as closed keeps an unknown number
// in the unknown fields.
func TestClosedEnum(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages")
	source := out["com/safetyculture/s12kt/kitchen/editions/closed.pb.kt"]
	require.Contains(t, source, "Closed.UNRECOGNIZED) _unknown.addVarint(16,")
	require.Contains(t, source, "v_open = com.safetyculture.s12kt.kitchen.editions.Open.fromNumber(reader.readInt32())")
	require.NotContains(t, source, "Open.UNRECOGNIZED) _unknown")
	require.NotContains(t, out["com/safetyculture/s12kt/kitchen/v1/kitchen.pb.kt"], "_unknown.addVarint")
	require.Contains(t, out["com/safetyculture/s12kt/kitchen/legacy/legacy.pb.kt"], "_unknown.addVarint")
}

func TestKotlinMessagesOptionOff(t *testing.T) {
	for _, parameter := range []string{"", "kotlin_package", "kotlin_package=com.example.kt", "swift_module"} {
		for name := range kitchenOutput(t, parameter) {
			require.False(t, strings.HasSuffix(name, ".kt"), "parameter %q wrote %s", parameter, name)
		}
	}
}

func TestKotlinMessagesPackageFallback(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages,kotlin_package=com.example.kt")
	require.Contains(t, out["com/example/kt/kitchen/v1/kitchen.pb.kt"], "\npackage com.example.kt.kitchen.v1\n")
}

func TestMessagesOnly(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages,messages_only")
	expected := map[string]bool{
		"com/safetyculture/s12kt/kitchen/editions/closed.pb.kt": true,
		"com/safetyculture/s12kt/kitchen/editions/utf8.pb.kt":   true,
		"com/safetyculture/s12kt/kitchen/legacy/legacy.pb.kt":   true,
		"com/safetyculture/s12kt/kitchen/legacy/strict.pb.kt":   true,
		"com/safetyculture/s12kt/kitchen/shadow/shadow.pb.kt":   true,
		"com/safetyculture/s12kt/kitchen/v1/kitchen.pb.kt":      true,
		"com/safetyculture/s12kt/s12kt.pro":                     true,
		"com/safetyculture/s12kt/wire/S12Wire.kt":               true,
	}
	got := map[string]bool{}
	for name := range out {
		got[name] = true
	}
	require.Equal(t, expected, got)

	withOutputs := kitchenOutput(t, "kotlin_messages")
	for name := range expected {
		require.Equal(t, withOutputs[name], out[name], name)
	}
	requirePluginResponseError(t, []string{"testdata/kotlin/kitchen"}, "messages_only", "messages_only requires kotlin_messages", kitchenFiles...)
}

func TestKotlinNameCollisions(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{"fields", "string foo = 1; string _foo = 2;", `foo and _foo map to the same Kotlin name "foo"`},
		{"oneof property", "oneof my_choice { string value = 1; } string myChoice = 2;", `my_choice and myChoice map to the same Kotlin name "myChoice"`},
		{"oneof members", "oneof choice { string foo = 1; string _foo = 2; }", `foo and _foo map to the same Kotlin name "Foo"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			source := "syntax = \"proto3\";\npackage collision.v1;\nmessage M { " + tc.body + " }\n"
			require.NoError(t, os.WriteFile(filepath.Join(dir, "collision.proto"), []byte(source), 0o644))
			requirePluginResponseError(t, []string{dir}, "kotlin_messages", tc.want, "collision.proto")
		})
	}
}

func TestKotlinClassNameCollisions(t *testing.T) {
	cases := []struct {
		name string
		root string
		body string
		want string
	}{
		{"root segment", "io.acme", "message io {} message ioMessage {}", `io and ioMessage map to the same Kotlin name "ioMessage"`},
		{"reserved root", "com.safetyculture.s12kt", "message kotlin {} message kotlinMessage {}", `kotlin and kotlinMessage map to the same Kotlin name "kotlinMessage"`},
		{"nested", "io.acme", "message Outer { message io {} message ioMessage {} }", `io and ioMessage map to the same Kotlin name "ioMessage"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			source := "syntax = \"proto3\";\npackage collision.v1;\n" + tc.body + "\n"
			require.NoError(t, os.WriteFile(filepath.Join(dir, "collision.proto"), []byte(source), 0o644))
			requirePluginResponseError(t, []string{dir}, "kotlin_messages="+tc.root, tc.want, "collision.proto")
		})
	}
}

func TestRequiredMessageCyclesRejectMissingFields(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"self", "message A { required A a = 1; optional int32 x = 2; }"},
		{"indirect", "message B { required C c = 1; } message C { required B b = 1; }"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			source := "syntax = \"proto2\";\npackage cycle.v1;\n" + tc.body + "\n"
			require.NoError(t, os.WriteFile(filepath.Join(dir, "cycle.proto"), []byte(source), 0o644))
			out := runPluginOnDisk(t, []string{dir}, "kotlin_messages", "cycle.proto")
			code := out["com/safetyculture/s12kt/cycle/v1/cycle.pb.kt"]
			require.NotContains(t, code, ".decode(com.safetyculture.s12kt.wire.WireReader(kotlin.ByteArray(0)))")
			require.Contains(t, code, "throw com.safetyculture.s12kt.wire.InvalidWireException(\"missing required field ")
		})
	}
}

func TestKotlinKeywordNameCollisions(t *testing.T) {
	cases := []struct {
		name string
		file *descriptorpb.FileDescriptorProto
		want string
	}{
		{
			name: "field",
			file: &descriptorpb.FileDescriptorProto{
				Name:    proto.String("collision.proto"),
				Syntax:  proto.String("proto3"),
				Package: proto.String("collision.v1"),
				MessageType: []*descriptorpb.DescriptorProto{{
					Name: proto.String("M"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{Name: proto.String("object"), Number: proto.Int32(1), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()},
						{Name: proto.String("object_"), Number: proto.Int32(2), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()},
					},
				}},
			},
			want: `object and object_ map to the same Kotlin name "object_"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := protodesc.NewFile(tc.file, nil)
			require.NoError(t, err)
			_, err = kotlin.GenerateFile(file, "com.safetyculture.s12kt")
			require.ErrorContains(t, err, tc.want)
		})
	}
}

func TestKotlinRootShadowsDeclarations(t *testing.T) {
	dir := t.TempDir()
	source := "syntax = \"proto3\";\npackage shadow.v1;\nmessage io { string io = 1; }\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "shadow.proto"), []byte(source), 0o644))
	out := runPluginOnDisk(t, []string{dir}, "kotlin_messages=io.acme", "shadow.proto")
	require.Contains(t, out["io/acme/shadow/v1/shadow.pb.kt"], "public data class ioMessage(")
	require.Contains(t, out["io/acme/shadow/v1/shadow.pb.kt"], "public val io_: kotlin.String")
}

// wideMessage writes a proto file whose only message has count fields of the given type.
func wideMessage(t *testing.T, fieldType string, count int) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("syntax = \"proto3\";\npackage wide.v1;\nmessage Wide {\n")
	for i := 1; i <= count; i++ {
		fmt.Fprintf(&b, "  %s f%d = %d;\n", fieldType, i, i)
	}
	b.WriteString("}\n")
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "wide.proto"), []byte(b.String()), 0o644))
	return dir
}

// TestConstructorSlotLimit pins the largest message whose constructor fits in a JVM method.
func TestConstructorSlotLimit(t *testing.T) {
	const fits = 244 // 245 parameters, 8 masks and 2 more slots make 255
	resp, err := pluginResponseOnDisk(t, []string{wideMessage(t, "int32", fits)}, "kotlin_messages", "wide.proto")
	require.NoError(t, err)
	require.Empty(t, resp.GetError())

	requirePluginResponseError(t, []string{wideMessage(t, "int32", fits+1)}, "kotlin_messages", "wide.v1.Wide: the constructor needs 256 parameter slots and the JVM allows 255", "wide.proto")

	// A 64-bit value without presence takes two slots.
	requirePluginResponseError(t, []string{wideMessage(t, "int64", fits-1)}, "kotlin_messages", "wide.v1.Wide: the constructor needs", "wide.proto")
}
