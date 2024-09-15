package main

import (
	"fmt"
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/plugin"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"log"
	"os"
)

func main() {
	// the file name is hard coded to match the output from the makefile target `playground-req-object`
	req, err := os.ReadFile("./code_generator_request.pb.bin")
	if err != nil {
		log.Fatal(err)
	}

	cgReq := pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(req, &cgReq); err != nil {
		panic(fmt.Errorf("proto.Unmarshal(): %w", err))
	}

	opts := protogen.Options{}
	protoPlugins, err := opts.New(&cgReq)
	if err != nil {
		panic(fmt.Errorf("protogen.Options.New(): %w", err))
	}

	for _, f := range protoPlugins.Files {
		if !f.Generate {
			continue
		}

		plugin.GenerateFile(protoPlugins, f)
		resp := protoPlugins.Response()

		for _, f := range resp.File {
			os.Stdout.WriteString(">>>>>>>>>>>>>>> " + *f.Name + "\n")
			_, err := os.Stdout.WriteString(*f.Content)
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.WriteString("<<<<<<<<<<<<<<< " + *f.Name + "\n")
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString("\n")
		}
	}
}
