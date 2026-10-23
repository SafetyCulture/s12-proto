// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cruxgen "github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/generator"
	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/protoutil"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// compileProtos compiles proto files of a directory with their comments and returns the descriptors, imports first.
func compileProtos(t *testing.T, dir string, files []string) []*descriptorpb.FileDescriptorProto {
	t.Helper()
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: []string{dir},
		}),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	compiled, err := compiler.Compile(context.Background(), files...)
	require.NoError(t, err)

	seen := map[string]bool{}
	var protos []*descriptorpb.FileDescriptorProto
	var add func(fd protoreflect.FileDescriptor)
	add = func(fd protoreflect.FileDescriptor) {
		if seen[fd.Path()] {
			return
		}
		seen[fd.Path()] = true
		for i := range fd.Imports().Len() {
			add(fd.Imports().Get(i).FileDescriptor)
		}
		protos = append(protos, protoutil.ProtoFromFileDescriptor(fd))
	}
	for _, fd := range compiled {
		add(fd)
	}
	return protos
}

// protoFiles lists the proto files under dir, as the import paths of the files.
func protoFiles(t *testing.T, dir string) []string {
	t.Helper()
	var files []string
	require.NoError(t, filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && filepath.Ext(path) == ".proto" {
			rel, relErr := filepath.Rel(dir, path)
			require.NoError(t, relErr)
			files = append(files, filepath.ToSlash(rel))
		}
		return err
	}))
	return files
}

// requireGoldenDir compares the output with the golden files of a directory. Each output file has a
// golden file, and each golden file has an output file.
func requireGoldenDir(t *testing.T, out []outputFile, goldenDir string) {
	t.Helper()
	var goldenFiles []string
	require.NoError(t, filepath.WalkDir(goldenDir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, relErr := filepath.Rel(goldenDir, path)
			require.NoError(t, relErr)
			goldenFiles = append(goldenFiles, filepath.ToSlash(rel))
		}
		return err
	}))
	var rendered []string
	for _, f := range out {
		rendered = append(rendered, f.Name)
		golden, err := os.ReadFile(filepath.Join(goldenDir, filepath.FromSlash(f.Name)))
		require.NoError(t, err, "missing golden file for %s", f.Name)
		require.Equal(t, string(golden), f.Content, "output differs from the golden file %s", f.Name)
	}
	require.ElementsMatch(t, rendered, goldenFiles, "a golden file has no output")
}

// TestKitchenGolden compares the output for protos of every construct that the whitelist does not
// use with the output of protoc-gen-swift 1.38.1, once for each line of kitchen/variants.txt.
// The golden files come from testdata/regen.sh.
func TestKitchenGolden(t *testing.T) {
	dir := filepath.Join("testdata", "kitchen")
	files := protoFiles(t, filepath.Join(dir, "kitchen"))
	require.Len(t, files, 3)
	for i := range files {
		files[i] = "kitchen/" + filepath.Base(files[i])
	}
	protos := compileProtos(t, dir, files)
	variants := readLines(t, filepath.Join(dir, "variants.txt"))
	tested := 0
	for _, line := range variants {
		if strings.HasPrefix(line, "#") {
			continue
		}
		name, parameter, _ := strings.Cut(line, "|")
		tested++
		t.Run(name, func(t *testing.T) {
			opts, err := parseOptions(parameter, nil)
			require.NoError(t, err)
			out, err := render(protos, files, opts)
			require.NoError(t, err)
			require.Len(t, out, len(files))
			requireGoldenDir(t, out, filepath.Join(dir, "golden", name))
		})
	}
	golden, err := os.ReadDir(filepath.Join(dir, "golden"))
	require.NoError(t, err)
	require.Len(t, golden, tested, "a golden directory has no line in variants.txt")
}

// TestLargeGolden checks the code for an enum and a oneof with more cases than the limit of one switch statement.
func TestLargeGolden(t *testing.T) {
	dir := filepath.Join("testdata", "large")
	protos := compileProtos(t, dir, []string{"large.proto"})
	opts, err := parseOptions(goldenParameter, nil)
	require.NoError(t, err)
	out, err := render(protos, []string{"large.proto"}, opts)
	require.NoError(t, err)
	require.Contains(t, out[0].Content, "// Can't get here, all the cases are listed in the above switches.")
	require.Contains(t, out[0].Content, "    default: break\n    }\n    switch self.pick {\n    case .f501?: try {")
	require.Contains(t, out[0].Content, "    default: break\n    }\n    switch self {\n    case .v500: return 500")
	requireGoldenDir(t, out, filepath.Join(dir, "golden"))
}

// TestUpstreamGolden compares the output for protos of swift-protobuf, which test the naming of
// types, fields and cases, reserved words, deprecation and groups, with the output of
// protoc-gen-swift 1.38.1 with the default options.
func TestUpstreamGolden(t *testing.T) {
	dir := filepath.Join("testdata", "upstream")
	files := protoFiles(t, dir)
	require.NotEmpty(t, files)
	protos := compileProtos(t, dir, files)
	opts, err := parseOptions("", nil)
	require.NoError(t, err)
	out, err := render(protos, files, opts)
	require.NoError(t, err)
	requireGoldenDir(t, out, filepath.Join(dir, "golden"))
}

// TestGoldenVersion checks that the golden files come from the version of protoc-gen-swift that the tests name.
func TestGoldenVersion(t *testing.T) {
	version, err := os.ReadFile(filepath.Join("testdata", "GOLDEN_VERSION"))
	require.NoError(t, err)
	require.Equal(t, "protoc-gen-swift 1.38.1", strings.TrimSpace(string(version)))
}

func TestParseOptions(t *testing.T) {
	opts, err := parseOptions("swift_messages=S12Protos,messages_only,Visibility=Public,FileNaming=pathToUnderscores,kotlin_messages", cruxgen.IsOptionKey)
	require.NoError(t, err)
	require.Equal(t, "public ", opts.visibilitySnippet())
	require.Equal(t, "s12_assets_v1_asset.pb.swift", opts.outputFilename("s12/assets/v1/asset.proto"))

	opts, err = parseOptions("", nil)
	require.NoError(t, err)
	require.Equal(t, "", opts.visibilitySnippet())
	require.Equal(t, "s12/assets/v1/asset.pb.swift", opts.outputFilename("s12/assets/v1/asset.proto"))

	opts, err = parseOptions("ExperimentalHiddenNames=", nil)
	require.NoError(t, err)
	require.Zero(t, opts.hiddenNames)

	opts, err = parseOptions("FileNaming=DropPath", nil)
	require.NoError(t, err)
	require.Equal(t, "asset.pb.swift", opts.outputFilename("s12/assets/v1/asset.proto"))

	for _, bad := range []string{
		"Visibility=open", "FileNaming=", "Unknown=1", "ProtoPathModuleMappings=map.txt",
		"ImplementationOnlyImports=true,UseAccessLevelOnImports=true", "ImplementationOnlyImports=true,Visibility=Public",
		"SwiftProtobufModuleName=1x",
	} {
		_, err := parseOptions(bad, nil)
		require.Error(t, err, bad)
	}
}

func TestGeneratorSwiftOptionKeysParse(t *testing.T) {
	values := map[string]string{
		"EnumGeneration":                           "None",
		"ExperimentalHiddenNames":                  "Fields",
		"FileNaming":                               "FullPath",
		"ImplementationOnlyImports":                "false",
		"ProtoPathModuleMappings":                  "",
		"SwiftProtobufModuleName":                  "SPB",
		"UseAccessLevelOnImports":                  "false",
		"Visibility":                               "internal",
		"experimental_strip_nonfunctional_codegen": "",
	}
	for key := range cruxgen.SwiftOptionKeys {
		t.Run(key, func(t *testing.T) {
			_, err := parseOptions("swift_messages,"+key+"="+values[key], cruxgen.IsOptionKey)
			require.NoError(t, err)
		})
	}
}

func TestNames(t *testing.T) {
	require.Equal(t, "S12ASSETSAsset", sanitizeMessageName("S12ASSETS"+"Asset", map[string]bool{"SwiftProtobuf": true}))
	require.Equal(t, "assetID", toLowerCamelCase("asset_id"))
	require.Equal(t, "someURL", toLowerCamelCase("some_url"))
	require.Equal(t, "HTTPServer", toUpperCamelCase("http_server"))
	require.Equal(t, "_9Lives", toLowerCamelCase("9lives"))
	require.Equal(t, "S12_Assets_V1_", typePrefix("s12.assets.v1", nil))
	prefix := "S12"
	require.Equal(t, "S12", typePrefix("s12.assets.v1", &prefix))
	require.Equal(t, "self_p", sanitizeFieldName("self"))
	require.Equal(t, "`class`", sanitizeFieldName("class"))
	require.Equal(t, "hasFoo_p", sanitizeFieldNameBasedOn("hasFoo", "hasFoo"))
	require.Equal(t, "MessageMessage", sanitizeMessageName("MessageMessage", nil))
	require.Equal(t, "TypeMessage", sanitizeMessageName("Type", nil))
}

func TestTrimBackticks(t *testing.T) {
	tests := map[string]string{
		"`":   "`",
		"`a":  "`a",
		"a`":  "a`",
		"`a`": "a",
		"a":   "a",
	}
	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			require.Equal(t, want, trimBackticks(input))
		})
	}
}

func TestBytecode(t *testing.T) {
	w := newNameInstructionWriter()
	require.True(t, w.shouldUseEmptyNameMapInitializer())
	w.writeSame(1, "id")
	w.writeStandard(2, "created_at")
	w.writeUnique(5, "a", "b")
	w.writeReservedName("old")
	w.writeReservedNumbers(100, 105)
	require.False(t, w.shouldUseEmptyNameMapInitializer())
	require.Equal(t, `"\0\u{1}id\0\u{3}created_at\0\u{6}\u{3}a\0b\0\u{b}old\0\u{c}d\u{1}\u{5}"`, w.bytecode.stringLiteral())
}

func TestPrintSwiftLiterals(t *testing.T) {
	require.Equal(t, `"a\"b\\c\u{0a}"`, stringToEscapedStringLiteral("a\"b\\c\n"))
	require.Equal(t, "String()", stringToEscapedStringLiteral(""))
	require.Equal(t, "Data([0, 1, 97, 255, 10])", escapedToDataLiteral(`\000\001a\377\n`))
	require.Equal(t, "Data()", escapedToDataLiteral(""))
}

func TestEscapedDataLiteralNamesNonOctalByte(t *testing.T) {
	require.PanicsWithValue(t, "non-octal byte 'x' in an octal escape", func() {
		escapedToDataLiteral(`\00x`)
	})

	broken := &descriptorpb.FileDescriptorProto{
		Name:   proto.String("broken.proto"),
		Syntax: proto.String("proto2"),
		MessageType: []*descriptorpb.DescriptorProto{{
			Name: proto.String("Message"),
			Field: []*descriptorpb.FieldDescriptorProto{{
				Name:         proto.String("data"),
				Number:       proto.Int32(1),
				Label:        descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:         descriptorpb.FieldDescriptorProto_TYPE_BYTES.Enum(),
				DefaultValue: proto.String(`\00x`),
			}},
		}},
	}
	opts, err := parseOptions("", nil)
	require.NoError(t, err)
	_, err = render([]*descriptorpb.FileDescriptorProto{broken}, []string{"broken.proto"}, opts)
	require.ErrorContains(t, err, "internal error of the Swift generator on broken.proto")
	require.ErrorContains(t, err, "non-octal byte 'x' in an octal escape")
}

// TestInternalErrorNamesFile checks that a descriptor that breaks the generator gives an error with the file name and the stack.
func TestInternalErrorNamesFile(t *testing.T) {
	broken := &descriptorpb.FileDescriptorProto{
		Name:     proto.String("broken.proto"),
		Syntax:   proto.String("proto3"),
		EnumType: []*descriptorpb.EnumDescriptorProto{{Name: proto.String("Empty")}},
	}
	opts, err := parseOptions("", nil)
	require.NoError(t, err)
	_, err = render([]*descriptorpb.FileDescriptorProto{broken}, []string{"broken.proto"}, opts)
	require.ErrorContains(t, err, "internal error of the Swift generator on broken.proto")
	require.ErrorContains(t, err, "goroutine")
}
