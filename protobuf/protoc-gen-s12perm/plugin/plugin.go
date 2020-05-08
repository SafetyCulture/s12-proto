// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	perms "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type s12perm struct {
	*generator.Generator
}

func New() generator.Plugin {
	return &s12perm{}
}

func (g *s12perm) Name() string {
	return "s12perm"
}

func (g *s12perm) Init(gen *generator.Generator) {
	g.Generator = gen
}

func (g *s12perm) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	for _, service := range file.FileDescriptorProto.Service {
		for _, method := range service.Method {
			flags := getPermissions(method)
			if len(flags) == 0 {
				continue
			}
			g.P("// Permissions found!")

		}
	}

}

func (g *s12perm) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	imports := generator.NewPluginImports(g.Generator)
	imports.GenerateImports(file)
}

func getPermissions(method *descriptor.MethodDescriptorProto) []string {
	if method.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(method.Options, perms.E_RequiredFlags)
	if err != nil {
		// option is optional so continue
		return nil
	}
	return v.([]string)
}
