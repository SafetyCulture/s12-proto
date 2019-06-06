// Copyright (c) 2019 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-gogrpcmock/plugin"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func main() {
	gen := generator.New()
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err, "gogrpcmock: reading input")
	}
	if err := proto.Unmarshal(data, gen.Request); err != nil {
		gen.Error(err, "gogrpcmock: parsing input proto")
	}

	filesToGenerate := make([]string, 0, len(gen.Request.ProtoFile))
	for _, file := range gen.Request.ProtoFile {
		if len(file.Service) > 0 {
			filesToGenerate = append(filesToGenerate, *file.Name)
		}
	}
	gen.Request.FileToGenerate = filesToGenerate

	if len(gen.Request.FileToGenerate) == 0 {
		log.Println("gogrpcmock: no files to generate")
		return
	}

	gen.CommandLineParameters(gen.Request.GetParameter())
	gen.WrapTypes()
	gen.SetPackageNames()
	gen.BuildTypeNameMap()

	gen.GeneratePlugin(plugin.New())

	for i := 0; i < len(gen.Response.File); i++ {
		gen.Response.File[i].Name = proto.String(strings.Replace(*gen.Response.File[i].Name, ".pb.go", ".mock.go", -1))
	}

	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "gogrpcmock: failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		gen.Error(err, "gogrpcmock: failed to write output proto")
	}
}
