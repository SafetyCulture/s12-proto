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
	require.NotRegexp(t, regexp.MustCompile(`(?m)^import `), runtime)
	require.NotContains(t, runtime, "com.google")
}

// paramLine matches one constructor parameter of a generated data class.
var paramLine = regexp.MustCompile("^public val `?(\\w+)`?: (.+?) = .+,$")

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
	case fd.Message() != nil:
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
	out = kitchenOutput(t, "kotlin_package=com.example.kt,kotlin_messages=com.example.other")
	require.Contains(t, out["com/example/other/kitchen/v1/kitchen.pb.kt"], "\npackage com.example.other.kitchen.v1\n")
	require.Contains(t, out["s12/kitchen/v1/kitchen.djinni.yaml"], "com.example.kt.kitchen.v1.")
}

func TestMessagesOnly(t *testing.T) {
	out := kitchenOutput(t, "kotlin_messages,messages_only")
	names := make([]string, 0, len(out))
	for name := range out {
		names = append(names, name)
		require.True(t, strings.HasSuffix(name, ".pb.kt") || strings.HasSuffix(name, "S12Wire.kt") || strings.HasSuffix(name, "s12kt.pro"), name)
	}
	require.Contains(t, names, "com/safetyculture/s12kt/kitchen/v1/kitchen.pb.kt")
	require.Contains(t, names, "com/safetyculture/s12kt/wire/S12Wire.kt")
	require.Contains(t, names, "com/safetyculture/s12kt/s12kt.pro")

	withOutputs := kitchenOutput(t, "kotlin_messages")
	for _, name := range names {
		require.Equal(t, withOutputs[name], out[name], name)
	}
	require.Empty(t, kitchenOutput(t, "messages_only"))
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
	_, err := pluginResponseOnDisk(t, []string{wideMessage(t, "int32", fits)}, "kotlin_messages", "wide.proto")
	require.NoError(t, err)

	_, err = pluginResponseOnDisk(t, []string{wideMessage(t, "int32", fits+1)}, "kotlin_messages", "wide.proto")
	require.ErrorContains(t, err, "wide.v1.Wide: the constructor needs 256 parameter slots and the JVM allows 255")

	// A 64-bit value without presence takes two slots.
	_, err = pluginResponseOnDisk(t, []string{wideMessage(t, "int64", fits-1)}, "kotlin_messages", "wide.proto")
	require.ErrorContains(t, err, "wide.v1.Wide: the constructor needs")
}
