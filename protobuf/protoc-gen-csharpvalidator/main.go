// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-csharpvalidator/plugin"
)

// includeFlag collects the --csharpvalidator_out=include=<glob> options. Nothing
// is generated until at least one is given, so adding the plugin to a build is
// not on its own a change to that build's output.
type includeFlag []string

func (f *includeFlag) String() string { return "" }

func (f *includeFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var flags flag.FlagSet
	var include includeFlag
	flags.Var(&include, "include", "path glob selecting the proto files to generate validators for; repeatable, and /** matches a subtree")
	runtime := flags.Bool("runtime", true, "write the support library the generated code calls into")

	opts := &protogen.Options{ParamFunc: flags.Set}
	opts.Run(func(p *protogen.Plugin) error {
		return plugin.Generate(p, plugin.Options{Include: include, SkipRuntime: !*runtime})
	})
}
