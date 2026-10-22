// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Package main implements protoc-gen-crux-kotlin, a protoc plugin that generates
// pure Kotlin message classes with a wire codec from proto files. The classes have
// no dependency except the Kotlin standard library and one generated runtime file.
package main

import (
	"context"

	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-crux-kotlin/internal/kotlin"
)

func main() {
	protoplugin.Main(protoplugin.HandlerFunc(handle))
}

func handle(
	_ context.Context,
	_ protoplugin.PluginEnv,
	responseWriter protoplugin.ResponseWriter,
	request protoplugin.Request,
) error {
	responseWriter.SetFeatureProto3Optional()
	responseWriter.SetFeatureSupportsEditions(descriptorpb.Edition_EDITION_PROTO2, descriptorpb.Edition_EDITION_2023)
	files, err := request.FileDescriptorsToGenerate()
	if err != nil {
		return err
	}
	return kotlin.Generate(responseWriter, files, kotlin.ParseOptions(request.Parameter()))
}
