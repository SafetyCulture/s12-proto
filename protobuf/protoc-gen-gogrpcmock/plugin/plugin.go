// Copyright (c) 2019 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/icrowley/fake"

	"github.com/SafetyCulture/s12-proto/protobuf/s12proto"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

var (
	rxID          = regexp.MustCompile("(?i)id")
	rxEmail       = regexp.MustCompile("(?i)email")
	rxPhone       = regexp.MustCompile("(?i)phone")
	rxURL         = regexp.MustCompile("(?i)url")
	rxDescription = regexp.MustCompile("(?i)^(description|desc)$")
	rxLatitude    = regexp.MustCompile("(?i)^(latitude|lat)$")
	rxLongitude   = regexp.MustCompile("(?i)^(longitude|lng)$")
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const repeatCount = 8

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
		// TODO: Depth should be a cmd param
		g.generateMockMessage(m, false, true, 5)
		g.P(`return res, nil`)
	} else {
		// Should this return an error?
		g.P(`return nil, nil`)
	}

	g.Out()
	g.P(`}`)
}

func (g *grpcmock) generateMockMessage(msg *generator.Descriptor, inner, nullable bool, depth int) {
	depth--
	if depth < 0 {
		return
	}

	msgName := g.TypeName(msg)
	if nullable {
		msgName = "&" + msgName
	}

	g.P(msgName, `{`)
	g.In()

	for _, field := range msg.Field {
		fieldName := g.GetFieldName(msg, field)
		fieldType, _ := g.GoType(msg, field)
		repeated := field.IsRepeated()
		nullable := gogoproto.IsNullable(field) && field.IsMessage()

		if field.OneofIndex != nil {
			// TODO: Deal with oneof
			continue
		}

		if field.IsString() {
			g.generateMockString(fieldName, fieldType, repeated, field)
		} else if isSupportedInt(field) {
			g.generateMockInt(fieldName, fieldType, repeated, field)
		} else if field.IsMessage() {
			g.generateMockInnerMessage(fieldName, fieldType, repeated, nullable, field, depth)
		}
	}
	g.Out()
	if inner {
		g.P(`},`)
	} else {
		g.P(`}`)
	}

}

func (g *grpcmock) generateMockString(fieldName, fieldType string, repeated bool, field *descriptor.FieldDescriptorProto) {
	if repeated {
		g.P(fieldName, `: `, fieldType, `{`)
		g.In()

		for i := 0; i < repeatCount; i++ {
			g.P(`"`, generateStringValue(fieldName, field), `",`)
		}

		g.Out()
		g.P(`},`)
		return
	}

	g.P(fieldName, `: "`, generateStringValue(fieldName, field), `",`)

}

func (g *grpcmock) generateMockInt(fieldName, fieldType string, repeated bool, field *descriptor.FieldDescriptorProto) {
	if repeated {
		g.P(fieldName, `: `, fieldType, `{`)
		g.In()

		for i := 0; i < repeatCount; i++ {
			g.P(generateIntValue(fieldName), `,`)
		}

		g.Out()
		g.P(`},`)
		return
	}
	g.P(fieldName, `: `, generateIntValue(fieldName), `,`)
}

func (g *grpcmock) generateMockInnerMessage(fieldName, fieldType string, repeated, nullable bool, field *descriptor.FieldDescriptorProto, depth int) {

	msgObj := g.objectNamed(field.GetTypeName())
	msg, ok := msgObj.(*generator.Descriptor)
	if ok && msg.GetOptions().GetMapEntry() {
		//TODO: return map entries
		return
	}

	length := 1

	if repeated {
		length = repeatCount
		g.P(fieldName, `: `, fieldType, `{`)
		g.In()
	} else {
		g.P(fieldName, `: `)
	}

	//TODO: Create helper methods to deal with *time.Time and *time.Duration

	for i := 0; i < length; i++ {
		switch fieldType {
		case "time.Time":
			g.P(`time.Now(),`)
		case "*time.Time", "*time.Duration":
			// g.P(`&time.Time{},`)
			g.P(`nil,`)
		default:
			g.generateMockMessage(msg, true, nullable, depth)
		}
	}

	if repeated {
		g.Out()
		g.P(`},`)
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

func generateStringValue(fieldName string, field *descriptor.FieldDescriptorProto) string {
	val := ""

	if mocks := getFieldMocksIfAny(field); mocks != nil && len(mocks.String_) > 0 {
		return mocks.String_[r.Intn(len(mocks.String_))]
	}

	if rxID.MatchString(fieldName) {
		val = uuid.Must(uuid.NewV4()).String()
	}

	if rxEmail.MatchString(fieldName) {
		val = fake.EmailAddress()
	}

	if rxPhone.MatchString(fieldName) {
		val = fake.Phone()
	}

	if rxDescription.MatchString(fieldName) {
		val = fake.Paragraph()
	}

	if rxURL.MatchString(fieldName) {
		val = fmt.Sprintf("https://%s/%s", strings.ToLower(fake.DomainName()), strings.Replace(fake.Words(), " ", "/", -1))
	}

	if val == "" {
		val = fake.Word()
	}
	return val
}

func generateIntValue(fieldName string) string {
	val := ""

	if rxLatitude.MatchString(fieldName) {
		val = strconv.Itoa(fake.LatitudeDegrees())
	}

	if rxLongitude.MatchString(fieldName) {
		val = strconv.Itoa(fake.LongitudeDegrees())
	}

	if val == "" {
		val = strconv.Itoa(int(r.Intn(1000)))
	}
	return val
}

func getFieldMocksIfAny(field *descriptor.FieldDescriptorProto) *s12proto.FieldMock {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, s12proto.E_Field)
		if err == nil && v.(*s12proto.FieldMock) != nil {
			return (v.(*s12proto.FieldMock))
		}
	}
	return nil
}
