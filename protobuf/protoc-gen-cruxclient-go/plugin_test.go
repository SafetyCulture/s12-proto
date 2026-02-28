// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/protoutil"
	"github.com/bufbuild/protoplugin"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// runPlugin compiles the given proto sources in-process and runs the plugin
// handler on the specified files to generate. Returns a map of output filename
// to content.
//
// protoSources maps logical proto import paths to their source content.
// filesToGenerate is the subset of protoSources that should produce output.
func runPlugin(t *testing.T, protoSources map[string]string, filesToGenerate []string) map[string]string {
	t.Helper()
	ctx := context.Background()

	// Build resolver: custom source map + standard imports (google/protobuf/*)
	sourceResolver := &protocompile.SourceResolver{
		Accessor: protocompile.SourceAccessorFromMap(protoSources),
	}
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(sourceResolver),
	}

	// Compile ALL files in the source map so that transitive dependencies are
	// included in the ProtoFile list of the CodeGeneratorRequest. protoplugin
	// needs the full closure to resolve imports when unmarshaling the request.
	allFiles := make([]string, 0, len(protoSources))
	for path := range protoSources {
		allFiles = append(allFiles, path)
	}
	compiled, err := compiler.Compile(ctx, allFiles...)
	require.NoError(t, err, "protocompile.Compile failed")

	// Collect all compiled file descriptor protos, including transitive imports
	// (standard well-known types like google/protobuf/descriptor.proto) so that
	// protoplugin can resolve them when unmarshaling the CodeGeneratorRequest.
	seen := make(map[string]struct{})
	var protoFiles []*descriptorpb.FileDescriptorProto
	for i := 0; i < len(compiled); i++ {
		collectTransitiveDeps(compiled[i], seen, &protoFiles)
	}

	// Build the CodeGeneratorRequest.
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: filesToGenerate,
		ProtoFile:      protoFiles,
	}
	reqBytes, err := proto.Marshal(req)
	require.NoError(t, err, "proto.Marshal(CodeGeneratorRequest) failed")

	// Run the plugin handler via stdin/stdout.
	var stdout bytes.Buffer
	var stderr strings.Builder
	env := protoplugin.Env{
		Args:   []string{},
		Stdin:  bytes.NewReader(reqBytes),
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err = protoplugin.Run(ctx, env, protoplugin.HandlerFunc(handle))
	require.NoError(t, err, "protoplugin.Run failed; stderr: %s", stderr.String())

	// Unmarshal the response.
	resp := &pluginpb.CodeGeneratorResponse{}
	err = proto.Unmarshal(stdout.Bytes(), resp)
	require.NoError(t, err, "proto.Unmarshal(CodeGeneratorResponse) failed")

	// Collect output files into a map.
	output := make(map[string]string, len(resp.GetFile()))
	for _, f := range resp.GetFile() {
		output[f.GetName()] = f.GetContent()
	}
	return output
}

// mustReadFile reads a file relative to the test's working directory.
func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err, "os.ReadFile(%q) failed", path)
	return data
}

// TestRouteGuideGolden is an end-to-end golden file test that:
//  1. Compiles the routeguide proto in-process via protocompile
//  2. Runs the plugin handler
//  3. Compares .crux.api.h and .crux.api.cc output against reference files
func TestRouteGuideGolden(t *testing.T) {
	// Load proto sources from testdata.
	routeGuideProto := mustReadFile(t, "testdata/routeguide/route_guide.proto")
	messageProto := mustReadFile(t, "testdata/routeguide/message.proto")
	wireOptionsProto := mustReadFile(t, "testdata/routeguide/wire_options.proto")

	// Map logical import paths to file content.
	// route_guide.proto uses: import "routeguide/v1/message.proto"; import "wire_options.proto";
	// message.proto uses:     import "wire_options.proto";
	// wire_options.proto uses: import "google/protobuf/descriptor.proto"; (resolved by WithStandardImports)
	protoSources := map[string]string{
		"routeguide/v1/route_guide.proto": string(routeGuideProto),
		"routeguide/v1/message.proto":     string(messageProto),
		"wire_options.proto":              string(wireOptionsProto),
	}

	// Run plugin — only generate output for route_guide.proto.
	output := runPlugin(t, protoSources, []string{"routeguide/v1/route_guide.proto"})

	// Assert both API files are present in the output.
	require.Contains(t, output, "routeguide/v1/route_guide.crux.api.h",
		"plugin output missing .crux.api.h; got keys: %v", mapKeys(output))
	require.Contains(t, output, "routeguide/v1/route_guide.crux.api.cc",
		"plugin output missing .crux.api.cc; got keys: %v", mapKeys(output))

	// Compare .crux.api.h against golden reference.
	goldenH := mustReadFile(t, "testdata/routeguide/route_guide.crux.api.h")
	require.Equal(t, string(goldenH), output["routeguide/v1/route_guide.crux.api.h"],
		"API header mismatch — diff output against testdata/routeguide/route_guide.crux.api.h")

	// Compare .crux.api.cc against golden reference.
	goldenCC := mustReadFile(t, "testdata/routeguide/route_guide.crux.api.cc")
	require.Equal(t, string(goldenCC), output["routeguide/v1/route_guide.crux.api.cc"],
		"API source mismatch — diff output against testdata/routeguide/route_guide.crux.api.cc")

	// Compare .mock.h against golden reference.
	require.Contains(t, output, "routeguide/v1/route_guide.mock.h",
		"plugin output missing .mock.h; got keys: %v", mapKeys(output))

	goldenMock := mustReadFile(t, "testdata/routeguide/route_guide.mock.h")
	require.Equal(t, string(goldenMock), output["routeguide/v1/route_guide.mock.h"],
		"Mock header mismatch — diff output against testdata/routeguide/route_guide.mock.h")
}

// TestMessageDjinniYAMLGolden is a golden file test for the .djinni.yaml output.
// It generates from message.proto (which has messages but no services, triggering
// the PLUG-05 djinni-only path) and compares against the C++ reference output.
//
// message.proto exercises:
//   - 7 messages (Point, Rectangle, Feature, RouteNote, RouteSummary, RouteSummary.Details, RouteSummary.Details.MoreDetails)
//   - DetailsMapEntry synthetic from map<string, Details> that must be skipped
//   - wire_options extension (triggers kotlin block)
//   - Nested messages (tests underscore vs colon name derivation)
func TestMessageDjinniYAMLGolden(t *testing.T) {
	// Load proto sources from testdata.
	messageProto := mustReadFile(t, "testdata/routeguide/message.proto")
	wireOptionsProto := mustReadFile(t, "testdata/routeguide/wire_options.proto")

	// Map logical import paths to file content.
	protoSources := map[string]string{
		"routeguide/v1/message.proto": string(messageProto),
		"wire_options.proto":          string(wireOptionsProto),
	}

	// Run plugin — generate output for message.proto (no services → djinni only path).
	output := runPlugin(t, protoSources, []string{"routeguide/v1/message.proto"})

	// Assert YAML file is present in the output.
	require.Contains(t, output, "routeguide/v1/message.djinni.yaml",
		"plugin output missing .djinni.yaml; got keys: %v", mapKeys(output))

	// Compare .djinni.yaml against golden reference.
	goldenYAML := mustReadFile(t, "testdata/routeguide/message.djinni.yaml")
	require.Equal(t, string(goldenYAML), output["routeguide/v1/message.djinni.yaml"],
		"YAML mismatch — diff output against testdata/routeguide/message.djinni.yaml")
}

// TestMessageDjinniObjCGolden is a golden file test for the .djinni.objc.h output.
// It generates from message.proto (which has messages but no services, triggering
// the PLUG-05 djinni-only path) and compares against the C++ reference output.
//
// Validates: 7 Translator blocks (one per non-map-entry message), per-message
// #include/#import, namespace close comment format (//namespace, no space),
// NSData serialization bridge in toCpp/fromCpp.
func TestMessageDjinniObjCGolden(t *testing.T) {
	// Load proto sources from testdata.
	messageProto := mustReadFile(t, "testdata/routeguide/message.proto")
	wireOptionsProto := mustReadFile(t, "testdata/routeguide/wire_options.proto")

	// Map logical import paths to file content.
	protoSources := map[string]string{
		"routeguide/v1/message.proto": string(messageProto),
		"wire_options.proto":          string(wireOptionsProto),
	}

	// Run plugin — generate output for message.proto (no services -> djinni only path).
	output := runPlugin(t, protoSources, []string{"routeguide/v1/message.proto"})

	// Assert ObjC file is present in the output.
	require.Contains(t, output, "routeguide/v1/message.djinni.objc.h",
		"plugin output missing .djinni.objc.h; got keys: %v", mapKeys(output))

	// Compare .djinni.objc.h against golden reference.
	goldenObjC := mustReadFile(t, "testdata/routeguide/message.djinni.objc.h")
	require.Equal(t, string(goldenObjC), output["routeguide/v1/message.djinni.objc.h"],
		"ObjC mismatch — diff output against testdata/routeguide/message.djinni.objc.h")
}

// TestMessageDjinniJNIGolden is a golden file test for the .djinni.jni.h output.
// It generates from message.proto (which has messages but no services, triggering
// the PLUG-05 djinni-only path) and compares against the C++ reference output.
//
// Validates: 7 JNIInfo+Translator block pairs (one per non-map-entry message),
// djinni_support.hpp include at top, per-message #include (no #import),
// jbyteArray marshaling in toCpp/fromCpp, UnderscoresToDollar for nested names.
func TestMessageDjinniJNIGolden(t *testing.T) {
	// Load proto sources from testdata.
	messageProto := mustReadFile(t, "testdata/routeguide/message.proto")
	wireOptionsProto := mustReadFile(t, "testdata/routeguide/wire_options.proto")

	// Map logical import paths to file content.
	protoSources := map[string]string{
		"routeguide/v1/message.proto": string(messageProto),
		"wire_options.proto":          string(wireOptionsProto),
	}

	// Run plugin — generate output for message.proto (no services -> djinni only path).
	output := runPlugin(t, protoSources, []string{"routeguide/v1/message.proto"})

	// Assert JNI file is present in the output.
	require.Contains(t, output, "routeguide/v1/message.djinni.jni.h",
		"plugin output missing .djinni.jni.h; got keys: %v", mapKeys(output))

	// Compare .djinni.jni.h against golden reference.
	goldenJNI := mustReadFile(t, "testdata/routeguide/message.djinni.jni.h")
	require.Equal(t, string(goldenJNI), output["routeguide/v1/message.djinni.jni.h"],
		"JNI mismatch — diff output against testdata/routeguide/message.djinni.jni.h")
}

// collectTransitiveDeps recursively walks a FileDescriptor's transitive import
// closure and appends each unique FileDescriptorProto to out. seen tracks
// already-visited paths to avoid duplicates.
func collectTransitiveDeps(fd protoreflect.FileDescriptor, seen map[string]struct{}, out *[]*descriptorpb.FileDescriptorProto) {
	path := string(fd.Path())
	if _, already := seen[path]; already {
		return
	}
	seen[path] = struct{}{}

	// First recurse into imports so dependencies come before dependants.
	imports := fd.Imports()
	for i := 0; i < imports.Len(); i++ {
		imp := imports.Get(i)
		collectTransitiveDeps(imp.FileDescriptor, seen, out)
	}
	*out = append(*out, protoutil.ProtoFromFileDescriptor(fd))
}

// mapKeys returns the keys of a map for use in error messages.
func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
