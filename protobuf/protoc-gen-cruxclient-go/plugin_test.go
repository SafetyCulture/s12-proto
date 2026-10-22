// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"bytes"
	"context"
	"os"
	"regexp"
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
	return runPluginWithParameter(t, protoSources, filesToGenerate, "")
}

// runPluginWithParameter is runPlugin with a protoc plugin parameter string.
func runPluginWithParameter(t *testing.T, protoSources map[string]string, filesToGenerate []string, parameter string) map[string]string {
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
		Parameter:      proto.String(parameter),
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

// messageProtoSources returns the sources that the message.proto golden tests compile.
func messageProtoSources(t *testing.T) map[string]string {
	t.Helper()
	return map[string]string{
		"routeguide/v1/message.proto": string(mustReadFile(t, "testdata/routeguide/message.proto")),
		"wire_options.proto":          string(mustReadFile(t, "testdata/routeguide/wire_options.proto")),
	}
}

// TestSwiftOptionOffKeepsOutput checks that a request without the Swift option
// produces no Swift output and the same djinni files as the request with an unrelated option.
func TestSwiftOptionOffKeepsOutput(t *testing.T) {
	files := []string{"routeguide/v1/message.proto"}
	plain := runPlugin(t, messageProtoSources(t), files)
	other := runPluginWithParameter(t, messageProtoSources(t), files, "unrelated=1")

	require.Equal(t, plain, other)
	require.NotContains(t, plain, "routeguide/v1/message.djinni.swiftxx.h")
	require.NotContains(t, plain["routeguide/v1/message.djinni.yaml"], "swift")
}

// TestMessageDjinniSwiftGolden is a golden file test for the swift: and swiftxx:
// yaml blocks and the .djinni.swiftxx.h output, with the module name option.
func TestMessageDjinniSwiftGolden(t *testing.T) {
	output := runPluginWithParameter(t, messageProtoSources(t),
		[]string{"routeguide/v1/message.proto"}, "swift_module=RouteGuideProtos")

	for _, name := range []string{"message.djinni.yaml", "message.djinni.swiftxx.h"} {
		require.Contains(t, output, "routeguide/v1/"+name,
			"plugin output missing %s; got keys: %v", name, mapKeys(output))
		golden := mustReadFile(t, "testdata/routeguide/swift/"+name)
		require.Equal(t, string(golden), output["routeguide/v1/"+name],
			"mismatch — diff output against testdata/routeguide/swift/%s", name)
	}

	// The djinni files of the other languages do not change.
	plain := runPlugin(t, messageProtoSources(t), []string{"routeguide/v1/message.proto"})
	for _, name := range []string{"message.djinni.objc.h", "message.djinni.jni.h"} {
		require.Equal(t, plain["routeguide/v1/"+name], output["routeguide/v1/"+name])
	}
}

// TestSwiftOptionDefaultModule checks that a bare swift_module option selects S12Protos.
func TestSwiftOptionDefaultModule(t *testing.T) {
	output := runPluginWithParameter(t, messageProtoSources(t),
		[]string{"routeguide/v1/message.proto"}, "swift_module")

	require.Contains(t, output["routeguide/v1/message.djinni.yaml"], "  module: 'S12Protos'\n")
	require.Contains(t, output, "routeguide/v1/message.djinni.swiftxx.h")
}

// swiftTypeNames returns the swift: typename values of a yaml output, in file order.
func swiftTypeNames(yaml string) []string {
	var names []string
	for _, match := range regexp.MustCompile(`(?m)^swift:\n  typename: '([^']+)'`).FindAllStringSubmatch(yaml, -1) {
		names = append(names, match[1])
	}
	return names
}

// TestSwiftTypeNames checks the type names against the output of protoc-gen-swift 1.26.0
// (protoc-swift:1.3.2) for the same protos: reserved nested names, no package, swift_prefix.
func TestSwiftTypeNames(t *testing.T) {
	sources := map[string]string{}
	for _, name := range []string{"names", "loose", "prefixed", "wellknown"} {
		sources["swiftnames/"+name+".proto"] = string(mustReadFile(t, "testdata/swiftnames/"+name+".proto"))
	}
	output := runPluginWithParameter(t, sources, []string{"swiftnames/names.proto", "swiftnames/loose.proto", "swiftnames/prefixed.proto", "swiftnames/wellknown.proto"}, "swift_module")

	require.Equal(t, []string{
		"S12_NameTest_V1_Outer",
		"S12_NameTest_V1_Outer.DataMessage",
		"S12_NameTest_V1_Outer.IntMessage",
		"S12_NameTest_V1_Outer.StringMessage",
		"S12_NameTest_V1_Outer.TypeMessage",
		"S12_NameTest_V1_Outer.ProtocolMessage",
		"S12_NameTest_V1_Outer.AnyMessage",
		"S12_NameTest_V1_Outer.SelfMessage",
		"S12_NameTest_V1_Outer.Message",
		"S12_NameTest_V1_Outer.DataMessageMessage",
		"S12_NameTest_V1_Outer.Error",
		"S12_NameTest_V1_Outer.Date",
		"S12_NameTest_V1_Outer.Foo",
		"S12_NameTest_V1_Outer.Foo.TypeMessage",
		"S12_NameTest_V1_Data",
	}, swiftTypeNames(output["swiftnames/names.djinni.yaml"]))

	for _, name := range []string{"names", "loose", "prefixed", "wellknown"} {
		golden := mustReadFile(t, "testdata/swiftnames/"+name+".djinni.yaml")
		require.Equal(t, string(golden), output["swiftnames/"+name+".djinni.yaml"])
	}

	// A file outside the bundled google/protobuf files keeps the module of the option, even in package google.protobuf.
	require.Equal(t, []string{"Google_Protobuf_Timestamp"}, swiftTypeNames(output["swiftnames/wellknown.djinni.yaml"]))
	require.Contains(t, output["swiftnames/wellknown.djinni.yaml"], "  module: 'S12Protos'\n")
	require.Contains(t, output["swiftnames/names.djinni.yaml"], "  module: 'S12Protos'\n")

	require.Equal(t, []string{"Loose", "Loose.Inner"}, swiftTypeNames(output["swiftnames/loose.djinni.yaml"]))
	require.Contains(t, output["swiftnames/loose.djinni.yaml"], "  prefix: ''\n")

	require.Equal(t, []string{"SPWithSwiftPrefix", "SPWithSwiftPrefix.DataMessage"}, swiftTypeNames(output["swiftnames/prefixed.djinni.yaml"]))
	require.Contains(t, output["swiftnames/prefixed.djinni.yaml"], "  prefix: 'SP'\n")
}

// TestSwiftBundledProtoModule checks that the types of a bundled google/protobuf file use the
// SwiftProtobuf module, as protoc-gen-swift does.
func TestSwiftBundledProtoModule(t *testing.T) {
	sources := map[string]string{
		"google/protobuf/timestamp.proto": string(mustReadFile(t, "testdata/google/protobuf/timestamp.proto")),
	}
	output := runPluginWithParameter(t, sources, []string{"google/protobuf/timestamp.proto"}, "swift_module")

	yaml := output["google/protobuf/timestamp.djinni.yaml"]
	require.Equal(t, []string{"Google_Protobuf_Timestamp"}, swiftTypeNames(yaml))
	require.Contains(t, yaml, "  module: 'SwiftProtobuf'\n")
	require.NotContains(t, yaml, "S12Protos")
	require.Equal(t, string(mustReadFile(t, "testdata/google/protobuf/timestamp.djinni.yaml")), yaml)
}

// TestCppSwiftGoldenMatchesGo checks that the golden files of the C++ plugin and of
// the Go plugin are the same bytes.
func TestCppSwiftGoldenMatchesGo(t *testing.T) {
	for _, name := range []string{"message.djinni.yaml", "message.djinni.swiftxx.h"} {
		goGolden := mustReadFile(t, "testdata/routeguide/swift/"+name)
		cppGolden := mustReadFile(t, "../protoc-gen-cruxclient/generated-swift/routeguide/v1/"+name)
		require.Equal(t, string(goGolden), string(cppGolden), "%s differs between the C++ and the Go golden files", name)
	}
	for _, name := range []string{"names", "loose", "prefixed", "wellknown"} {
		goGolden := mustReadFile(t, "testdata/swiftnames/"+name+".djinni.yaml")
		cppGolden := mustReadFile(t, "../protoc-gen-cruxclient/generated-swift/swiftnames/"+name+".djinni.yaml")
		require.Equal(t, string(goGolden), string(cppGolden), "%s differs between the C++ and the Go golden files", name)
	}
	goGolden := mustReadFile(t, "testdata/google/protobuf/timestamp.djinni.yaml")
	cppGolden := mustReadFile(t, "../protoc-gen-cruxclient/generated-swift/google/protobuf/timestamp.djinni.yaml")
	require.Equal(t, string(goGolden), string(cppGolden), "timestamp.djinni.yaml differs between the C++ and the Go golden files")
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
