// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

// TestKotlinOptionOffKeepsOutput checks that a request without the Kotlin option
// produces no Kotlin output.
func TestKotlinOptionOffKeepsOutput(t *testing.T) {
	output := runPlugin(t, messageProtoSources(t), []string{"routeguide/v1/message.proto"})

	require.NotContains(t, output, "routeguide/v1/message.djinni.kotlin.jni.h")
	for _, key := range []string{"jniTypename", "kotlinTranslator", "kotlinHeader", "kotlinTypeSignature"} {
		require.NotContains(t, output["routeguide/v1/message.djinni.yaml"], key)
	}
}

// TestMessageDjinniKotlinGolden is a golden file test for the Kotlin keys of the
// yaml output and the .djinni.kotlin.jni.h output.
func TestMessageDjinniKotlinGolden(t *testing.T) {
	output := runPluginWithParameter(t, messageProtoSources(t),
		[]string{"routeguide/v1/message.proto"}, "kotlin_package")

	for _, name := range []string{"message.djinni.yaml", "message.djinni.kotlin.jni.h"} {
		require.Contains(t, output, "routeguide/v1/"+name,
			"plugin output missing %s; got keys: %v", name, mapKeys(output))
		golden := mustReadFile(t, "testdata/routeguide/kotlin/"+name)
		require.Equal(t, string(golden), output["routeguide/v1/"+name],
			"mismatch — diff output against testdata/routeguide/kotlin/%s", name)
	}

	// The djinni files of the other languages do not change.
	plain := runPlugin(t, messageProtoSources(t), []string{"routeguide/v1/message.proto"})
	for _, name := range []string{"message.djinni.objc.h", "message.djinni.jni.h"} {
		require.Equal(t, plain["routeguide/v1/"+name], output["routeguide/v1/"+name])
	}
}

// TestKotlinPackageOption checks the package root of the Kotlin class names.
func TestKotlinPackageOption(t *testing.T) {
	output := runPluginWithParameter(t, messageProtoSources(t),
		[]string{"routeguide/v1/message.proto"}, "kotlin_package=com.example.kt")

	yaml := output["routeguide/v1/message.djinni.yaml"]
	require.Contains(t, yaml, "  jniTypename: 'com.example.kt.routeguide.v1.Point'\n")
	require.Contains(t, yaml, "  jniTypename: 'com.example.kt.routeguide.v1.RouteSummary.Details.MoreDetails'\n")
	require.Contains(t, yaml, "  kotlinTypeSignature: 'Lcom/example/kt/routeguide/v1/RouteSummary$Details$MoreDetails;'\n")
	require.Contains(t, output["routeguide/v1/message.djinni.kotlin.jni.h"],
		`jniFindClass("com/example/kt/routeguide/v1/RouteSummary$Details")`)
}

// TestKotlinBlockWithoutWirePackage checks that a file without the wire_package option gets a
// kotlin: block with only jniTypename.
func TestKotlinBlockWithoutWirePackage(t *testing.T) {
	sources := map[string]string{
		"swiftnames/loose.proto": string(mustReadFile(t, "testdata/swiftnames/loose.proto")),
	}
	plain := runPlugin(t, sources, []string{"swiftnames/loose.proto"})
	require.NotContains(t, plain["swiftnames/loose.djinni.yaml"], "kotlin:")

	output := runPluginWithParameter(t, sources, []string{"swiftnames/loose.proto"}, "kotlin_package")
	require.Contains(t, output["swiftnames/loose.djinni.yaml"], "kotlin:\n  jniTypename: 'com.safetyculture.s12kt.Loose'\n")
}

// TestKotlinPackageDropsLeadingS12 checks that the package s12.name_test.v1 maps to the
// Kotlin package com.safetyculture.s12kt.name_test.v1, and a nested class has a "$" in its signature.
func TestKotlinPackageDropsLeadingS12(t *testing.T) {
	sources := map[string]string{
		"swiftnames/names.proto": string(mustReadFile(t, "testdata/swiftnames/names.proto")),
	}
	output := runPluginWithParameter(t, sources, []string{"swiftnames/names.proto"}, "kotlin_package")

	yaml := output["swiftnames/names.djinni.yaml"]
	require.Contains(t, yaml, "  jniTypename: 'com.safetyculture.s12kt.name_test.v1.Outer'\n")
	require.Contains(t, yaml, "  jniTypename: 'com.safetyculture.s12kt.name_test.v1.Outer.Foo.Type'\n")
	require.Contains(t, yaml, "  kotlinTypeSignature: 'Lcom/safetyculture/s12kt/name_test/v1/Outer$Foo$Type;'\n")
}

// swiftTypeNames returns the swift: typename values of a yaml output, in file order.
func swiftTypeNames(yaml string) []string {
	var names []string
	for _, match := range regexp.MustCompile(`(?m)^swift:\n  typename: '([^']+)'`).FindAllStringSubmatch(yaml, -1) {
		names = append(names, match[1])
	}
	return names
}

// TestSwiftTypeNames checks the type names against the output of protoc-gen-swift for the same
// protos: reserved nested names, no package, swift_prefix. The names of reserved.proto are the
// output of protoc-gen-swift 1.38.1 and 1.39.0, and the other names are the output of 1.26.0.
func TestSwiftTypeNames(t *testing.T) {
	sources := map[string]string{}
	var files []string
	for _, name := range []string{"names", "loose", "prefixed", "wellknown", "reserved"} {
		sources["swiftnames/"+name+".proto"] = string(mustReadFile(t, "testdata/swiftnames/"+name+".proto"))
		files = append(files, "swiftnames/"+name+".proto")
	}
	output := runPluginWithParameter(t, sources, files, "swift_module")

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
		"S12_NameTest_V1_Outer.SendableMessage",
		"S12_NameTest_V1_Outer.EquatableMessage",
		"S12_NameTest_V1_Outer.HashableMessage",
		"S12_NameTest_V1_Outer.SwiftMessage",
		"S12_NameTest_V1_Outer.Date",
		"S12_NameTest_V1_Outer.Foo",
		"S12_NameTest_V1_Outer.Foo.TypeMessage",
		"S12_NameTest_V1_Data",
	}, swiftTypeNames(output["swiftnames/names.djinni.yaml"]))

	require.Equal(t, []string{
		"SendableMessage", "EquatableMessage", "HashableMessage", "SwiftMessage", "SwiftProtobufMessage",
		"ExtensionsMessage", "letMessage", "selfMessage", "funcMessage", "_Message", "__Message",
		"SendableMessageMessage", "Outer", "Outer.SendableMessage", "Outer.SwiftMessage", "Outer.inMessage",
		"Outer._Message",
	}, swiftTypeNames(output["swiftnames/reserved.djinni.yaml"]))

	for _, name := range []string{"names", "loose", "prefixed", "wellknown", "reserved"} {
		golden := mustReadFile(t, "testdata/swiftnames/"+name+".djinni.yaml")
		require.Equal(t, string(golden), output["swiftnames/"+name+".djinni.yaml"])
	}

	// A file outside the bundled google/protobuf files keeps the module of the option, even in package google.protobuf.
	require.Equal(t, []string{"Google_Protobuf_Timestamp"}, swiftTypeNames(output["swiftnames/wellknown.djinni.yaml"]))
	require.Contains(t, output["swiftnames/wellknown.djinni.yaml"], "  module: 'S12Protos'\n")
	require.Contains(t, output["swiftnames/names.djinni.yaml"], "  module: 'S12Protos'\n")

	require.Equal(t, []string{"Loose", "Loose.Inner"}, swiftTypeNames(output["swiftnames/loose.djinni.yaml"]))
	require.NotContains(t, output["swiftnames/loose.djinni.yaml"], "  prefix:")

	require.Equal(t, []string{"SPWithSwiftPrefix", "SPWithSwiftPrefix.DataMessage"}, swiftTypeNames(output["swiftnames/prefixed.djinni.yaml"]))
	require.NotContains(t, output["swiftnames/prefixed.djinni.yaml"], "  prefix: 'SP'")
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

// readTree returns the files in dir that keep accepts, keyed by their path relative to dir. It
// reads the directories below dir only if recursive is true.
func readTree(t *testing.T, dir string, recursive bool, keep func(name string) bool) map[string]string {
	t.Helper()
	files := map[string]string{}
	require.NoError(t, filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if !recursive && path != dir {
				return fs.SkipDir
			}
			return nil
		}
		if !keep(entry.Name()) {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		files[rel] = string(content)
		return err
	}))
	return files
}

// requireSameFiles fails with the path of every file that only one side has and of every file whose bytes differ.
func requireSameFiles(t *testing.T, a, b map[string]string, aName, bName string) {
	t.Helper()
	var problems []string
	for name, content := range a {
		other, ok := b[name]
		switch {
		case !ok:
			problems = append(problems, "missing in the "+bName+": "+name)
		case other != content:
			problems = append(problems, "differs: "+name)
		}
	}
	for name := range b {
		if _, ok := a[name]; !ok {
			problems = append(problems, "missing in the "+aName+": "+name)
		}
	}
	sort.Strings(problems)
	require.Empty(t, problems, "%s and %s are not the same files", aName, bName)
}

// cppOutputCase is one run of the C++ plugin that has its output in a directory of the repository.
type cppOutputCase struct {
	cppDir    string   // the directory of the C++ output
	sources   []string // proto paths below testdata or, for routeguide, the cruxclient proto directory
	parameter string
	prefix    string // the output path of the files of the run
}

// requireGoMatchesCpp runs the Go plugin for each case and compares every file that it writes with
// the files of the C++ plugin in the directory of the case: the yaml, the swiftxx, jni, objc and
// Kotlin jni headers, in both directions.
func requireGoMatchesCpp(t *testing.T, cases []cppOutputCase) {
	t.Helper()
	goFiles, cppFiles := map[string]string{}, map[string]string{}
	for _, c := range cases {
		sources := map[string]string{}
		for _, path := range c.sources {
			sources[path] = string(mustReadFile(t, sourcePath(path)))
		}
		if _, ok := sources["routeguide/v1/message.proto"]; ok {
			sources["wire_options.proto"] = string(mustReadFile(t, "testdata/routeguide/wire_options.proto"))
		}
		for name, content := range runPluginWithParameter(t, sources, c.sources, c.parameter) {
			if strings.HasPrefix(name, c.prefix) && strings.Contains(name, ".djinni.") {
				goFiles[filepath.Join(c.cppDir, strings.TrimPrefix(name, c.prefix))] = content
			}
		}
		for name, content := range readTree(t, c.cppDir, false, func(string) bool { return true }) {
			cppFiles[filepath.Join(c.cppDir, name)] = content
		}
	}
	requireSameFiles(t, goFiles, cppFiles, "Go plugin output", "C++ plugin output")
}

// sourcePath returns the file that holds the proto source of a path in a test request.
func sourcePath(path string) string {
	if path == "routeguide/v1/message.proto" {
		return "testdata/routeguide/message.proto"
	}
	return filepath.Join("testdata", path)
}

// TestCppSwiftOutputMatchesGo checks that the Go plugin and the C++ plugin write the same bytes with
// the Swift option, for every file that the C++ plugin wrote into the repository.
func TestCppSwiftOutputMatchesGo(t *testing.T) {
	const swift = "../protoc-gen-cruxclient/generated-swift/"
	// A request cannot hold wellknown.proto and timestamp.proto: both declare google.protobuf.Timestamp.
	requireGoMatchesCpp(t, []cppOutputCase{
		{swift + "routeguide/v1", []string{"routeguide/v1/message.proto"}, "swift_module=RouteGuideProtos", "routeguide/v1/"},
		{swift + "swiftnames", []string{"swiftnames/names.proto", "swiftnames/loose.proto", "swiftnames/prefixed.proto", "swiftnames/reserved.proto"}, "swift_module", "swiftnames/"},
		{swift + "swiftnames", []string{"swiftnames/wellknown.proto"}, "swift_module", "swiftnames/"},
		{swift + "google/protobuf", []string{"google/protobuf/timestamp.proto"}, "swift_module", "google/protobuf/"},
	})
}

// TestKotlinSwiftnamesGolden checks the Kotlin output for a package with the leading s12
// segment, nested classes, and a file without a package.
func TestKotlinSwiftnamesGolden(t *testing.T) {
	sources := map[string]string{}
	var files []string
	for _, name := range []string{"names", "loose"} {
		sources["swiftnames/"+name+".proto"] = string(mustReadFile(t, "testdata/swiftnames/"+name+".proto"))
		files = append(files, "swiftnames/"+name+".proto")
	}
	output := runPluginWithParameter(t, sources, files, "kotlin_package")

	got := map[string]string{}
	for name, content := range output {
		if strings.HasSuffix(name, ".djinni.yaml") || strings.HasSuffix(name, ".djinni.kotlin.jni.h") {
			got[strings.TrimPrefix(name, "swiftnames/")] = content
		}
	}
	golden := readTree(t, "testdata/swiftnames/kotlin", true, func(string) bool { return true })
	requireSameFiles(t, golden, got, "golden files", "plugin output")
}

// TestCppKotlinOutputMatchesGo checks that the Go plugin and the C++ plugin write the same bytes with
// the Kotlin option, for every file that the C++ plugin wrote into the repository.
func TestCppKotlinOutputMatchesGo(t *testing.T) {
	const kotlin = "../protoc-gen-cruxclient/generated-kotlin/"
	requireGoMatchesCpp(t, []cppOutputCase{
		{kotlin + "routeguide/v1", []string{"routeguide/v1/message.proto"}, "kotlin_package", "routeguide/v1/"},
	})
	requireGoMatchesCpp(t, []cppOutputCase{
		{kotlin + "swiftnames", []string{"swiftnames/names.proto", "swiftnames/loose.proto"}, "kotlin_package", "swiftnames/"},
	})
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
