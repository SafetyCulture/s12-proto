// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-csharpvalidator/plugin"
)

func main() {
	var flags flag.FlagSet
	opts := &protogen.Options{ParamFunc: flags.Set}

	opts.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return plugin.Generate(p)
	})
}
