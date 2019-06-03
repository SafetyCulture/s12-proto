// Copyright (c) 2019 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type grpcmock struct {
	*generator.Generator
	contextPkg string
}

func New() generator.Plugin {
	return &grpcmock{}
}

func (g *grpcmock) Name() string {
	return "grpcmock"
}

func (g *grpcmock) Init(gen *generator.Generator) {
	g.Generator = gen
}

func (g *grpcmock) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	for _, service := range file.FileDescriptorProto.Service {
		g.mockService(file, service)
	}
}

func (g *grpcmock) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	imports := generator.NewPluginImports(g.Generator)
	g.contextPkg = imports.NewImport("context").Use()
	imports.GenerateImports(file)
}

func (g *grpcmock) mockService(file *generator.FileDescriptor, service *descriptor.ServiceDescriptorProto) {
	origServName := service.GetName()
	// fullServName := origServName
	// if pkg := file.GetPackage(); pkg != "" {
	// 	fullServName = pkg + "." + fullServName
	// }
	servName := generator.CamelCase(origServName)
	servTypeName := fmt.Sprintf("%sMock", servName)

	g.P(`type `, servTypeName, ` struct {}`)
	g.P()
	for _, method := range service.Method {
		g.mockMethod(servTypeName, method)
	}
}

func (g *grpcmock) mockMethod(servTypeName string, method *descriptor.MethodDescriptorProto) {
	if method.GetServerStreaming() || method.GetClientStreaming() {
		// TODO: [RC] Mock streaming methods.
		return
	}

	methName := generator.CamelCase(method.GetName())
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	g.P(`func (m *`, servTypeName, `) `, methName, `(ctx context.Context, req *`, inType, `) (*`, outType, `, error){`)
	g.In()

	msg := g.objectNamed(method.GetOutputType())
	if m, ok := msg.(*generator.Descriptor); ok && !m.GetOptions().GetMapEntry() {
		g.P(`res := `)
		g.generateMockMessage(m, false)
		g.P(`return res, nil`)
	} else {
		// Should this return an error?
		g.P(`return nil, nil`)
	}

	g.Out()
	g.P(`}`)
}

func (g *grpcmock) generateMockMessage(msg *generator.Descriptor, inner bool) {
	msgName := generator.CamelCaseSlice(msg.TypeName())
	g.P(`&`, msgName, `{`)
	g.In()
	for _, field := range msg.Field {
		fieldName := g.GetOneOfFieldName(msg, field)
		// repeated := field.IsRepeated()
		if field.IsString() {
			g.generateMockString(fieldName, field)
		} else if isSupportedInt(field) {
			g.generateMockInt(fieldName, field)
		} else if field.IsMessage() {
			g.generateMockInnerMessage(fieldName, field)
		}
	}
	g.Out()
	if inner {
		g.P(`},`)
	} else {
		g.P(`}`)
	}

}

func (g *grpcmock) generateMockString(fieldName string, field *descriptor.FieldDescriptorProto) {
	g.P(fieldName, `: "sddf",`)
}

func (g *grpcmock) generateMockInt(fieldName string, field *descriptor.FieldDescriptorProto) {
	g.P(fieldName, `: 100,`)
}

func (g *grpcmock) generateMockInnerMessage(fieldName string, field *descriptor.FieldDescriptorProto) {
	g.P(fieldName, `: `)
	msg := g.objectNamed(field.GetTypeName())
	if m, ok := msg.(*generator.Descriptor); ok && !m.GetOptions().GetMapEntry() {
		g.generateMockMessage(m, true)
	}
}

func (g *grpcmock) objectNamed(name string) generator.Object {
	g.RecordTypeUse(name)
	return g.ObjectNamed(name)
}

func (g *grpcmock) typeName(str string) string {
	return g.TypeName(g.objectNamed(str))
}

func isSupportedInt(field *descriptor.FieldDescriptorProto) bool {
	switch *(field.Type) {
	case descriptor.FieldDescriptorProto_TYPE_INT32, descriptor.FieldDescriptorProto_TYPE_INT64:
		return true
	case descriptor.FieldDescriptorProto_TYPE_UINT32, descriptor.FieldDescriptorProto_TYPE_UINT64:
		return true
	case descriptor.FieldDescriptorProto_TYPE_SINT32, descriptor.FieldDescriptorProto_TYPE_SINT64:
		return true
	}
	return false
}
