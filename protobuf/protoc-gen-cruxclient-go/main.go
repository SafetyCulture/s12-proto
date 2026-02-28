// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

// Package main implements protoc-gen-cruxclient-go, a protoc plugin that
// generates CruxClient C++/ObjC/JNI/YAML files from proto service definitions.
// It uses the bufbuild/protoplugin framework to handle stdin/stdout I/O.
package main

import (
	"context"

	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/generator"
)

func main() {
	protoplugin.Main(protoplugin.HandlerFunc(handle))
}

// handle is the plugin handler that mirrors cruxclient_generator.cc Generate()
// (lines 31-49). It dispatches on service count:
//   - services == 0: generate djinni support files only (if messages exist)
//   - services > 0:  generate all 6 file types
func handle(
	_ context.Context,
	_ protoplugin.PluginEnv,
	responseWriter protoplugin.ResponseWriter,
	request protoplugin.Request,
) error {
	fileDescriptors, err := request.FileDescriptorsToGenerate()
	if err != nil {
		return err
	}
	for _, fileDesc := range fileDescriptors {
		if fileDesc.Services().Len() == 0 {
			// PLUG-05: No services — generate djinni support files only
			generateDjinniSupport(responseWriter, fileDesc)
			continue
		}
		// PLUG-04: Has services — generate all 6 file types
		generateAPIHeader(responseWriter, fileDesc)
		generateAPISource(responseWriter, fileDesc)
		generateMock(responseWriter, fileDesc)
		generateDjinniSupport(responseWriter, fileDesc)
	}
	return nil
}

// generateAPIHeader generates the .crux.api.h file.
// Delegates to generator.GenerateAPIHeader for real C++ code generation.
func generateAPIHeader(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor) {
	generator.GenerateAPIHeader(rw, fileDesc)
}

// generateAPISource generates the .crux.api.cc file.
// Delegates to generator.GenerateAPISource for real C++ code generation.
func generateAPISource(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor) {
	generator.GenerateAPISource(rw, fileDesc)
}

// generateMock generates the .mock.h file.
// Delegates to generator.GenerateMock for real C++ code generation.
func generateMock(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor) {
	generator.GenerateMock(rw, fileDesc)
}

// generateDjinniSupport generates the three djinni support files:
// .djinni.yaml, .djinni.objc.h, .djinni.jni.h.
// Mirrors api_generator.cc GenerateDjinniSupport (line 672): guards on
// top-level message count using fileDesc.Messages().Len() to match the C++
// file->message_type_count() check (which counts only top-level messages,
// NOT nested/recursive messages).
func generateDjinniSupport(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor) {
	// Match C++ behavior: no djinni files if no top-level messages
	if fileDesc.Messages().Len() == 0 {
		return
	}
	generator.GenerateDjinniYAML(rw, fileDesc)
	generator.GenerateDjinniObjC(rw, fileDesc)
	generator.GenerateDjinniJNI(rw, fileDesc)
}

// nonClientStreamingMethods returns methods that are not client-streaming.
// Both client-only and bidirectional streaming methods are excluded because
// both have IsStreamingClient() == true.
// This implements PLUG-06, mirroring api_generator.cc line 110:
//
//	if (method->client_streaming()) { continue; }
//
// Examples from routeguide RouteGuide:
//   - GetFeature (unary): INCLUDED
//   - UpdateFeature (unary): INCLUDED
//   - ListFeatures (server-streaming): INCLUDED
//   - RecordRoute (client-streaming): SKIPPED
//   - RouteChat (bidi): SKIPPED (client_streaming=true)
func nonClientStreamingMethods(service protoreflect.ServiceDescriptor) []protoreflect.MethodDescriptor {
	var methods []protoreflect.MethodDescriptor
	allMethods := service.Methods()
	for i := 0; i < allMethods.Len(); i++ {
		method := allMethods.Get(i)
		if method.IsStreamingClient() {
			continue
		}
		methods = append(methods, method)
	}
	return methods
}
