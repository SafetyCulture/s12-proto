// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-tsvalidator/plugin"
)

func main() {
	var flags flag.FlagSet
	opts := &protogen.Options{
		ParamFunc: flags.Set,
	}

	opts.Run(func(p *protogen.Plugin) error {
		for _, f := range p.Files {
			if f.Generate {
				plugin.GenerateFile(p, f)
			}
		}
		return nil
	})
}
