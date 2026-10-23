// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Package main implements protoc-gen-cruxclient-go, a protoc plugin that
// generates CruxClient C++/ObjC/JNI/YAML files from proto service definitions.
// It uses the bufbuild/protoplugin framework to handle stdin/stdout I/O.
package main

import (
	"context"
	"fmt"

	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/generator"
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/kotlin"
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
	responseWriter.SetFeatureProto3Optional()
	responseWriter.SetFeatureSupportsEditions(descriptorpb.Edition_EDITION_PROTO2, descriptorpb.Edition_EDITION_2023)
	fileDescriptors, err := request.FileDescriptorsToGenerate()
	if err != nil {
		return err
	}
	opts, err := generator.ParseOptions(request.Parameter())
	if err != nil {
		return err
	}
	cruxclientFiles := fileDescriptors
	if opts.MessagesOnly {
		cruxclientFiles = nil
	}
	for _, fileDesc := range cruxclientFiles {
		if fileDesc.Services().Len() == 0 {
			// PLUG-05: No services — generate djinni support files only
			generateDjinniSupport(responseWriter, fileDesc, opts)
			continue
		}
		// PLUG-04: Has services — generate all 6 file types
		generateAPIHeader(responseWriter, fileDesc)
		generateAPISource(responseWriter, fileDesc)
		generateMock(responseWriter, fileDesc)
		generateDjinniSupport(responseWriter, fileDesc, opts)
	}
	if opts.KotlinMessagesEnabled() {
		if err := kotlin.Generate(responseWriter, fileDescriptors, opts.KotlinMessages); err != nil {
			return fmt.Errorf("cruxclient-go: kotlin.Generate: %w", err)
		}
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
// .djinni.yaml, .djinni.objc.h, .djinni.jni.h. The Swift option adds .djinni.swiftxx.h. The Kotlin option adds .djinni.kotlin.jni.h.
// Mirrors api_generator.cc GenerateDjinniSupport: guards on
// top-level message count using fileDesc.Messages().Len() to match the C++
// file->message_type_count() check (which counts only top-level messages,
// NOT nested/recursive messages).
func generateDjinniSupport(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor, opts generator.Options) {
	// Match C++ behavior: no djinni files if no top-level messages
	if fileDesc.Messages().Len() == 0 {
		return
	}
	generator.GenerateDjinniYAML(rw, fileDesc, opts)
	generator.GenerateDjinniObjC(rw, fileDesc)
	generator.GenerateDjinniJNI(rw, fileDesc)
	if opts.SwiftEnabled() {
		generator.GenerateDjinniSwiftxx(rw, fileDesc)
	}
	if opts.KotlinEnabled() {
		generator.GenerateDjinniKotlinJNI(rw, fileDesc, opts)
	}
}
