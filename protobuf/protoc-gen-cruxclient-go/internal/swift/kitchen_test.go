// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/protoutil"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var kitchenFiles = []string{"kitchen/editions.proto", "kitchen/proto2.proto", "kitchen/proto3.proto"}

// compileKitchen compiles the kitchen protos with their comments and returns the descriptors, imports first.
func compileKitchen(t *testing.T) []*descriptorpb.FileDescriptorProto {
	t.Helper()
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: []string{filepath.Join("testdata", "kitchen")},
		}),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	compiled, err := compiler.Compile(context.Background(), kitchenFiles...)
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

// TestKitchenGolden compares the output for protos of every construct that the
// whitelist does not use with the output of protoc-gen-swift 1.38.1. The
// golden files come from the command in testdata/regen.sh.
func TestKitchenGolden(t *testing.T) {
	protos := compileKitchen(t)
	for _, tt := range []struct{ dir, parameter string }{
		{"default", ""},
		{"public", "Visibility=Public,FileNaming=PathToUnderscores"},
	} {
		t.Run(tt.dir, func(t *testing.T) {
			opts, err := parseOptions(tt.parameter)
			require.NoError(t, err)
			out, err := render(protos, kitchenFiles, opts)
			require.NoError(t, err)
			require.Len(t, out, len(kitchenFiles))
			for _, f := range out {
				golden, err := os.ReadFile(filepath.Join("testdata", "kitchen", "golden", tt.dir, filepath.FromSlash(f.Name)))
				require.NoError(t, err, "missing golden file for %s", f.Name)
				require.Equal(t, string(golden), f.Content, "output differs from the golden file %s", f.Name)
			}
		})
	}
}

func TestParseOptions(t *testing.T) {
	opts, err := parseOptions("swift_messages=S12Protos,messages_only,Visibility=Public,FileNaming=pathToUnderscores,kotlin_messages")
	require.NoError(t, err)
	require.Equal(t, "public ", opts.visibilitySnippet())
	require.Equal(t, "s12_assets_v1_asset.pb.swift", opts.outputFilename("s12/assets/v1/asset.proto"))

	opts, err = parseOptions("")
	require.NoError(t, err)
	require.Equal(t, "", opts.visibilitySnippet())
	require.Equal(t, "s12/assets/v1/asset.pb.swift", opts.outputFilename("s12/assets/v1/asset.proto"))

	opts, err = parseOptions("FileNaming=DropPath")
	require.NoError(t, err)
	require.Equal(t, "asset.pb.swift", opts.outputFilename("s12/assets/v1/asset.proto"))

	for _, bad := range []string{
		"Visibility=open", "FileNaming=", "Unknown=1", "ProtoPathModuleMappings=map.txt",
		"ImplementationOnlyImports=true,UseAccessLevelOnImports=true", "ImplementationOnlyImports=true,Visibility=Public",
		"SwiftProtobufModuleName=1x",
	} {
		_, err := parseOptions(bad)
		require.Error(t, err, bad)
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
