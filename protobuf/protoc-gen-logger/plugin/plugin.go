package plugin

import (
	"fmt"
	"strings"

	logger "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports

	fmtPkg    generator.Single
	logrusPkg generator.Single
	s12proto  generator.Single
}

func New() generator.Plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return "payloadlogger"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.fmtPkg = p.NewImport("fmt")
	p.logrusPkg = p.NewImport("github.com/sirupsen/logrus")
	p.s12proto = p.NewImport("github.com/SafetyCulture/s12-proto/protobuf/s12proto")

	for _, msg := range file.Messages() {
		p.generateParseFunction(file, msg)
	}
}

func (p *plugin) generateParseFunction(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())

	p.P(`func (this *`, ccTypeName, `) Parse(isLevelEnabled func(`, p.s12proto.Use(), `.Level) bool) proto.Message {`)
	p.In()

	oneOfsMap := make(map[string]string)
	p.P(`res:=&`, ccTypeName, `{}`)
	for _, field := range message.Field {
		var (
			fieldNameOneOf = p.GetOneOfFieldName(message, field)
			fieldName      = p.GetFieldName(message, field)
		)

		if fieldName != fieldNameOneOf { // this is oneOf field
			p.generateOneOfFieldRow(ccTypeName, fieldName, fieldNameOneOf, field)
			oneOfsMap[fieldName] = fieldName
			continue
		}

		if !hasPayloadLoggerExtensions(field) {
			p.generateFieldRow(field, fieldName)

			continue
		}

		p.P(`if isLevelEnabled(`, p.s12proto.Use(), `.Level_`, getLevelValue(field).String(), `) {`)
		p.In()
		p.generateFieldRow(field, fieldName)
		p.Out()
		p.P(`}`)
	}
	p.P(`return res`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateOneOfFieldRow(typeName, fieldName, fieldNameOneOf string, field *descriptor.FieldDescriptorProto) {
	fieldTypes := strings.Split(field.GetTypeName(), ".")
	fieldTypeName := fieldTypes[len(fieldTypes)-1]

	p.P(`if reflect.TypeOf(this.`, fieldName, `) == reflect.TypeOf(&`, typeName, `_`, fieldNameOneOf, `{}){`)
	p.In()
	if !hasPayloadLoggerExtensions(field) {
		if field.IsMessage() {
			oneOfVarName := fmt.Sprint(`oneof_`, typeName, `_`, fieldNameOneOf)
			p.P(oneOfVarName, `:=this.`, fieldName, `.(*`, typeName, `_`, fieldNameOneOf, `)`)
			p.P(oneOfVarName, `.`, fieldNameOneOf, `=`, oneOfVarName, `.`, fieldNameOneOf, `.Parse(isLevelEnabled).(*`, fieldTypeName, `)`)
			p.P(`res.`, fieldName, ` = `, oneOfVarName)
		} else {
			p.P(`res.`, fieldName, ` = this.`, fieldName)
		}

	} else {
		p.P(`if isLevelEnabled(`, p.s12proto.Use(), `.Level_`, getLevelValue(field).String(), `) {`)
		p.In()
		if field.IsMessage() {
			oneOfVarName := fmt.Sprintf(`oneof_%s_%s`, typeName, fieldNameOneOf)
			p.P(oneOfVarName, `:=this.`, fieldName, `.(*`, typeName, `_`, fieldNameOneOf, `)`)
			p.P(oneOfVarName, `.`, fieldNameOneOf, `=`, oneOfVarName, `.`, fieldNameOneOf, `.Parse(isLevelEnabled).(*`, fieldTypeName, `)`)
			p.P(`res.`, fieldName, ` = `, oneOfVarName)
		} else {
			p.P(`res.`, fieldName, ` = this.`, fieldName)
		}
		p.Out()
		p.P(`}`)
	}
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateFieldRow(field *descriptor.FieldDescriptorProto, fieldName string) {
	fieldTypes := strings.Split(field.GetTypeName(), ".")
	fieldTypeName := fieldTypes[len(fieldTypes)-1]

	if field.IsMessage() {
		p.P(`res.`, fieldName, `=this.`, fieldName, `.Parse(isLevelEnabled).(*`, fieldTypeName, `)`)
	} else {
		p.P(`res.`, fieldName, `=this.`, fieldName)
	}
}

func getLevelValue(field *descriptor.FieldDescriptorProto) *logger.Level {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, logger.E_Level)
		if err == nil {
			return v.(*logger.Level)
		}

	}

	return nil
}

func hasPayloadLoggerExtensions(field *descriptor.FieldDescriptorProto) bool {
	if field.Options != nil {
		validExts := []*proto.ExtensionDesc{
			logger.E_Level,
		}
		for _, ext := range validExts {
			if proto.HasExtension(field.Options, ext) {
				return true
			}
		}
	}
	return false
}
