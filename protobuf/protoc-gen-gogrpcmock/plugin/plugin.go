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

type oneofField struct {
	fieldName string
	subFields []*oneofSubField
}

type oneofSubField struct {
	fieldName       string
	typeName        string
	fieldStructName string
	protoField      *descriptor.FieldDescriptorProto
}

func New() generator.Plugin {
	return &grpcmock{}
}

func (g *grpcmock) Name() string {
	return "grpcmock"
}

func (g *grpcmock) Init(gen *generator.Generator) {
	g.Generator = gen
	fake.Seed(time.Now().UnixNano())
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

	oFields := make(map[int32]*oneofField)

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
		oneof := field.OneofIndex != nil

		if oneof {
			if oFields[*field.OneofIndex] == nil {
				of := oneofField{
					fieldName: fieldName,
				}
				oFields[*field.OneofIndex] = &of
			}

			sfFieldName := g.GetOneOfFieldName(msg, field)
			sfStructName := g.OneOfTypeName(msg, field)

			of := oFields[*field.OneofIndex]
			sf := oneofSubField{
				fieldName:       sfFieldName,
				fieldStructName: sfStructName,
				typeName:        fieldType,
				protoField:      field,
			}
			of.subFields = append(of.subFields, &sf)
			continue
		}

		g.generateMockField(fieldName, fieldType, repeated, nullable, field, depth)
	}

	for _, of := range oFields {
		sf := of.subFields[r.Intn(len(of.subFields))]
		nullable := gogoproto.IsNullable(sf.protoField) && sf.protoField.IsMessage()

		g.P(of.fieldName, `: `, `&`, sf.fieldStructName, `{`)
		g.In()
		g.generateMockField(sf.fieldName, sf.typeName, false, nullable, sf.protoField, 10)
		g.Out()
		g.P(`},`)
	}

	g.Out()
	if inner {
		g.P(`},`)
	} else {
		g.P(`}`)
	}

}

func (g *grpcmock) generateMockField(fieldName, fieldType string, repeated, nullable bool, field *descriptor.FieldDescriptorProto, depth int) {
	if field.IsString() {
		g.generateMockString(fieldName, fieldType, repeated, field)
	} else if isSupportedInt(field) {
		g.generateMockInt(fieldName, fieldType, repeated, field)
	} else if isSupportedFloat(field) {
		g.generateMockFloat(fieldName, fieldType, repeated, field)
	} else if field.IsEnum() {
		g.generateMockEnum(fieldName, fieldType, field)
	} else if field.IsBool() {
		g.generateMockBool(fieldName)
	} else if field.IsMessage() {
		g.generateMockInnerMessage(fieldName, fieldType, repeated, nullable, field, depth)
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
			g.P(generateIntValue(fieldName, field), `,`)
		}

		g.Out()
		g.P(`},`)
		return
	}
	g.P(fieldName, `: `, generateIntValue(fieldName, field), `,`)
}

func (g *grpcmock) generateMockFloat(fieldName, fieldType string, repeated bool, field *descriptor.FieldDescriptorProto) {
	if repeated {
		g.P(fieldName, `: `, fieldType, `{`)
		g.In()

		for i := 0; i < repeatCount; i++ {
			g.P(generateFloatValue(fieldName, field), `,`)
		}

		g.Out()
		g.P(`},`)
		return
	}
	g.P(fieldName, `: `, generateFloatValue(fieldName, field), `,`)
}

func (g *grpcmock) generateMockEnum(fieldName, fieldType string, field *descriptor.FieldDescriptorProto) {
	enum := g.ObjectNamed(field.GetTypeName()).(*generator.EnumDescriptor)
	enumValues := enum.GetValue()
	enumValue := enumValues[r.Intn(len(enumValues))]
	g.P(fieldName, `: `, strconv.Itoa(int(enumValue.GetNumber())), `,`)
}

func (g *grpcmock) generateMockBool(fieldName string) {
	bVal := "true"
	// air on the side of true vs false
	if n := r.Intn(3); n == 0 {
		bVal = "false"
	}
	g.P(fieldName, `: `, bVal, `,`)
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

func isSupportedFloat(field *descriptor.FieldDescriptorProto) bool {
	switch *(field.Type) {
	case descriptor.FieldDescriptorProto_TYPE_FLOAT, descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return true
	}
	return false
}

func generateStringValue(fieldName string, field *descriptor.FieldDescriptorProto) string {

	if mocks := getFieldMocksIfAny(field); mocks != nil {
		var sb strings.Builder

		if len(mocks.Prefix) > 0 {
			sb.WriteString(mocks.Prefix[r.Intn(len(mocks.Prefix))])
		}

		if len(mocks.String_) > 0 {
			sb.WriteString(mocks.String_[r.Intn(len(mocks.String_))])
		}

		if boolFromPtr(mocks.Word) {
			sb.WriteString(fake.Word())
		}

		if boolFromPtr(mocks.Words) {
			sb.WriteString(fake.Words())
		}

		if mocks.Wordsn != nil {
			sb.WriteString(fake.WordsN(int(*mocks.Wordsn)))
		}

		if boolFromPtr(mocks.Fullname) {
			sb.WriteString(fake.FullName())
		}

		if boolFromPtr(mocks.Firstname) {
			sb.WriteString(fake.FirstName())
		}

		if boolFromPtr(mocks.Lastname) {
			sb.WriteString(fake.LastName())
		}

		if boolFromPtr(mocks.Paragraph) {
			sb.WriteString(fake.Paragraph())
		}

		if boolFromPtr(mocks.Paragraphs) {
			sb.WriteString(fake.Paragraphs())
		}

		if mocks.Paragraphsn != nil {
			sb.WriteString(fake.ParagraphsN(int(*mocks.Paragraphsn)))
		}

		if boolFromPtr(mocks.Uuid) {
			sb.WriteString(uuid.Must(uuid.NewV4()).String())
		}

		if boolFromPtr(mocks.Email) {
			sb.WriteString(fake.EmailAddress())
		}

		if boolFromPtr(mocks.Phone) {
			sb.WriteString(fake.Phone())
		}

		if boolFromPtr(mocks.Company) {
			sb.WriteString(fake.Company())
		}

		if boolFromPtr(mocks.Brand) {
			sb.WriteString(fake.Brand())
		}

		if boolFromPtr(mocks.Product) {
			sb.WriteString(fake.ProductName())
		}

		if boolFromPtr(mocks.Color) {
			sb.WriteString(fake.Color())
		}

		if boolFromPtr(mocks.Hexcolor) {
			sb.WriteString(fake.HexColor())
		}

		return sb.String()
	}

	if rxID.MatchString(fieldName) {
		return uuid.Must(uuid.NewV4()).String()
	}

	if rxEmail.MatchString(fieldName) {
		return fake.EmailAddress()
	}

	if rxPhone.MatchString(fieldName) {
		return fake.Phone()
	}

	if rxDescription.MatchString(fieldName) {
		return fake.Paragraph()
	}

	if rxURL.MatchString(fieldName) {
		return fakeURL()
	}

	return fake.Word()
}

func fakeURL() string {
	return fmt.Sprintf("https://%s/%s", strings.ToLower(fake.DomainName()), strings.Replace(fake.Words(), " ", "/", -1))
}

func generateIntValue(fieldName string, field *descriptor.FieldDescriptorProto) string {
	if rxLatitude.MatchString(fieldName) {
		return strconv.Itoa(fake.LatitudeDegrees())
	}

	if rxLongitude.MatchString(fieldName) {
		return strconv.Itoa(fake.LongitudeDegrees())
	}

	n := 1000

	if mocks := getFieldMocksIfAny(field); mocks != nil {
		if mocks.Intn != nil && *mocks.Intn > 0 {
			n = int(*mocks.Intn)
		}
	}

	return strconv.Itoa(int(r.Intn(n)))
}

func generateFloatValue(fieldName string, field *descriptor.FieldDescriptorProto) string {
	if rxLatitude.MatchString(fieldName) {
		return strconv.FormatFloat(float64(fake.Latitude()), 'f', 6, 32)
	}

	if rxLongitude.MatchString(fieldName) {
		return strconv.FormatFloat(float64(fake.Longitude()), 'f', 6, 32)
	}

	n, m := 1000, 10000
	if mocks := getFieldMocksIfAny(field); mocks != nil {
		if mocks.Floatn != nil && *mocks.Floatn > 0 {
			n = int(*mocks.Floatn)
		}
	}
	return fmt.Sprintf("%d.%d", r.Intn(n), r.Intn(m))
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

func boolFromPtr(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
