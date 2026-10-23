// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"bytes"
	"context"
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

// runPluginWithComments is runPluginWithParameter with the comments of the proto
// files in the request, and it returns what the plugin writes to stderr.
func runPluginWithComments(t *testing.T, protoSources map[string]string, filesToGenerate []string, parameter string) (map[string]string, string) {
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
	require.NoError(t, protoplugin.Run(context.Background(), env, protoplugin.HandlerFunc(handle)), stderr.String())

	resp := &pluginpb.CodeGeneratorResponse{}
	require.NoError(t, proto.Unmarshal(stdout.Bytes(), resp))
	require.Empty(t, resp.GetError())
	output := map[string]string{}
	for _, f := range resp.GetFile() {
		output[f.GetName()] = f.GetContent()
	}
	return output, stderr.String()
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

// TestMessagesOnlyWithoutMessagesWritesNothing checks that messages_only with no
// messages option writes no file and says so on stderr.
func TestMessagesOnlyWithoutMessagesWritesNothing(t *testing.T) {
	output, stderr := runPluginWithComments(t, messageProtoSources(t), []string{"routeguide/v1/message.proto"}, "messages_only")

	require.Empty(t, output)
	require.Contains(t, stderr, "messages_only")
}

// TestSwiftMessagesRejectsUnknownOption checks that a mistyped option of the Swift output fails the run.
func TestSwiftMessagesRejectsUnknownOption(t *testing.T) {
	compiler := protocompile.Compiler{Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
		Accessor: protocompile.SourceAccessorFromMap(messageProtoSources(t)),
	})}
	compiled, err := compiler.Compile(context.Background(), "routeguide/v1/message.proto", "wire_options.proto")
	require.NoError(t, err)
	seen := map[string]struct{}{}
	var protoFiles []*descriptorpb.FileDescriptorProto
	for i := range compiled {
		collectTransitiveDeps(compiled[i], seen, &protoFiles)
	}
	reqBytes, err := proto.Marshal(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"routeguide/v1/message.proto"},
		ProtoFile:      protoFiles,
		Parameter:      proto.String("swift_messages,Visibilty=Public"),
	})
	require.NoError(t, err)

	var stdout bytes.Buffer
	var stderr strings.Builder
	env := protoplugin.Env{Args: []string{}, Stdin: bytes.NewReader(reqBytes), Stdout: &stdout, Stderr: &stderr}
	err = protoplugin.Run(context.Background(), env, protoplugin.HandlerFunc(handle))
	require.ErrorContains(t, err, "Visibilty")
}
