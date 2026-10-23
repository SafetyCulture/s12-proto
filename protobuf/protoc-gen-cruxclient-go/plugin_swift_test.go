// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"bytes"
	"context"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protoplugin"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// runPluginRequest compiles the proto files with their comments and runs the plugin handler on them.
// It returns the files that the plugin writes, what it writes to stderr, and its error.
func runPluginRequest(t *testing.T, protoSources map[string]string, filesToGenerate []string, parameter string) (map[string]string, string, error) {
	t.Helper()
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			Accessor: protocompile.SourceAccessorFromMap(protoSources),
		}),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	allFiles := make([]string, 0, len(protoSources))
	for path := range protoSources {
		allFiles = append(allFiles, path)
	}
	compiled, err := compiler.Compile(context.Background(), allFiles...)
	require.NoError(t, err)

	seen := map[string]struct{}{}
	var protoFiles []*descriptorpb.FileDescriptorProto
	for i := range compiled {
		collectTransitiveDeps(compiled[i], seen, &protoFiles)
	}
	reqBytes, err := proto.Marshal(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: filesToGenerate,
		ProtoFile:      protoFiles,
		Parameter:      proto.String(parameter),
	})
	require.NoError(t, err)

	var stdout bytes.Buffer
	var stderr strings.Builder
	env := protoplugin.Env{Args: []string{}, Stdin: bytes.NewReader(reqBytes), Stdout: &stdout, Stderr: &stderr}
	if err := protoplugin.Run(context.Background(), env, protoplugin.HandlerFunc(handle)); err != nil {
		return nil, stderr.String(), err
	}
	resp := &pluginpb.CodeGeneratorResponse{}
	require.NoError(t, proto.Unmarshal(stdout.Bytes(), resp))
	require.Empty(t, resp.GetError())
	output := map[string]string{}
	for _, f := range resp.GetFile() {
		output[f.GetName()] = f.GetContent()
	}
	return output, stderr.String(), nil
}

// runPluginWithComments is runPluginRequest for a run that must succeed.
func runPluginWithComments(t *testing.T, protoSources map[string]string, filesToGenerate []string, parameter string) (map[string]string, string) {
	t.Helper()
	output, stderr, err := runPluginRequest(t, protoSources, filesToGenerate, parameter)
	require.NoError(t, err, stderr)
	return output, stderr
}

func sortedKeys(m map[string]string) []string {
	keys := mapKeys(m)
	sort.Strings(keys)
	return keys
}

// TestSwiftMessagesOffWritesNoSwiftFile checks that a request without
// swift_messages gives the same files as before, and that none is a Swift file.
func TestSwiftMessagesOffWritesNoSwiftFile(t *testing.T) {
	files := []string{"routeguide/v1/message.proto"}
	plain := runPlugin(t, messageProtoSources(t), files)
	for name := range plain {
		require.False(t, strings.HasSuffix(name, ".pb.swift"), name)
	}
	with, _ := runPluginWithComments(t, messageProtoSources(t), files, "swift_messages")
	for name, content := range plain {
		require.Equal(t, content, with[name], "swift_messages changed %s", name)
	}
	require.Contains(t, with, "routeguide/v1/message.pb.swift")
	require.Len(t, with, len(plain)+1)
}

// TestSwiftMessagesOnly checks that messages_only leaves only the Swift file,
// and that its content is the output of protoc-gen-swift 1.38.1.
func TestSwiftMessagesOnly(t *testing.T) {
	output, stderr := runPluginWithComments(t, messageProtoSources(t), []string{"routeguide/v1/message.proto"},
		"swift_messages,messages_only,Visibility=Public,FileNaming=PathToUnderscores")

	require.Equal(t, []string{"routeguide_v1_message.pb.swift"}, sortedKeys(output))
	require.Empty(t, stderr)
	golden := mustReadFile(t, "testdata/swiftmessages/routeguide_v1_message.pb.swift")
	require.Equal(t, string(golden), output["routeguide_v1_message.pb.swift"])
}

// TestMessagesOnlyWithoutMessagesFails checks that messages_only with no messages option is an error.
func TestMessagesOnlyWithoutMessagesFails(t *testing.T) {
	_, _, err := runPluginRequest(t, messageProtoSources(t), []string{"routeguide/v1/message.proto"}, "messages_only")

	require.ErrorContains(t, err, "messages_only requires swift_messages or kotlin_messages")
}

// TestSwiftMessagesRejectsUnknownOption checks that a mistyped option of the Swift output fails the run.
func TestSwiftMessagesRejectsUnknownOption(t *testing.T) {
	_, _, err := runPluginRequest(t, messageProtoSources(t), []string{"routeguide/v1/message.proto"}, "swift_messages,Visibilty=Public")

	require.ErrorContains(t, err, "cruxclient-go: swift.Generate: ")
	require.ErrorContains(t, err, "Visibilty")
}

// TestSwiftMessagesKitchenThroughPlugin runs protos with a message set, groups and
// editions through the plugin handler, as the binary runs. It compares the output with the golden files of
// protoc-gen-swift. The runtime descriptors of Go refuse a message set, so the plugin
// must not build them when it writes only the Swift messages.
func TestSwiftMessagesKitchenThroughPlugin(t *testing.T) {
	root := filepath.Join("internal", "swift", "testdata", "kitchen")
	files := []string{"kitchen/editions.proto", "kitchen/proto2.proto", "kitchen/proto3.proto"}
	sources := map[string]string{}
	for _, name := range files {
		sources[name] = string(mustReadFile(t, filepath.Join(root, filepath.FromSlash(name))))
	}
	output, stderr := runPluginWithComments(t, sources, files, "swift_messages,messages_only,Visibility=Public,FileNaming=PathToUnderscores")

	require.Empty(t, stderr)
	golden, err := filepath.Glob(filepath.Join(root, "golden", "public", "*.pb.swift"))
	require.NoError(t, err)
	require.Len(t, output, len(golden))
	for _, path := range golden {
		require.Equal(t, string(mustReadFile(t, path)), output[filepath.Base(path)], "output differs from %s", path)
	}
}
