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
	refectPkg generator.Single
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
	p.refectPkg = p.NewImport("reflect")
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

	p.P(`type foo interface{`)
	p.In()
	p.P()
	p.P(`Parse(func(`, p.s12proto.Use(), `.Level) bool) proto.Message`)
	p.Out()
	p.P(`}`)

	oneOfsMap := make(map[string]string)
	p.P(`res:=&`, ccTypeName, `{}`)
	for _, field := range message.Field {
		var (
			fieldNameOneOf = p.GetOneOfFieldName(message, field)
			fieldName      = p.GetFieldName(message, field)
		)

		if fieldName != fieldNameOneOf { // this is oneOf field

			p.generateOneOfFieldRow(ccTypeName, fieldName, fieldNameOneOf, file.GoPackageName(), field)
			oneOfsMap[fieldName] = fieldName
			continue
		}

		if !hasPayloadLoggerExtensions(field) {
			p.generateParseCallRow(field, fieldName, file.GoPackageName())

			continue
		}

		p.P(`if isLevelEnabled(`, p.s12proto.Use(), `.Level_`, getLevelValue(field).String(), `) {`)
		p.In()
		p.generateParseCallRow(field, fieldName, file.GoPackageName())
		p.Out()
		p.P(`}`)
	}
	p.P(`return res`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateOneOfFieldRow(typeName, fieldName, fieldNameOneOf, namespace string, field *descriptor.FieldDescriptorProto) {
	fieldTypename := strings.TrimPrefix(field.GetTypeName(), ".")
	fieldTypename = strings.TrimPrefix(fieldTypename, namespace+".")

	p.P(`if reflect.TypeOf(this.`, fieldName, `) == reflect.TypeOf(&`, typeName, `_`, fieldNameOneOf, `{}){`)
	p.In()
	if !hasPayloadLoggerExtensions(field) {
		if field.IsMessage() {
			oneOfVarName := fmt.Sprint(`oneof_`, typeName, `_`, fieldNameOneOf)
			p.P(oneOfVarName, `:=this.`, fieldName, `.(*`, typeName, `_`, fieldNameOneOf, `)`)

			p.P(`mtype := `, p.refectPkg.Use(), `.TypeOf((*foo)(nil)).Elem()`)
			p.P(`fieldtype := `, p.refectPkg.Use(), `.TypeOf(`, oneOfVarName, `.`, fieldNameOneOf, `)`)

			p.P(`if `, oneOfVarName, `!=nil && fieldtype.Implements(mtype) {`)
			p.In()
			p.P(oneOfVarName, `.`, fieldNameOneOf, `=`, oneOfVarName, `.`, fieldNameOneOf, `.Parse(isLevelEnabled).(*`, fieldTypename, `)`)
			p.P(`res.`, fieldName, ` = `, oneOfVarName)
			p.Out()
			p.P(`}`)
		} else {
			p.P(`res.`, fieldName, ` = this.`, fieldName)
		}
	} else {
		p.P(`if isLevelEnabled(`, p.s12proto.Use(), `.Level_`, getLevelValue(field).String(), `) {`)
		p.In()
		if field.IsMessage() {
			oneOfVarName := fmt.Sprintf(`oneof_%s_%s`, typeName, fieldNameOneOf)
			p.P(oneOfVarName, `:=this.`, fieldName, `.(*`, typeName, `_`, fieldNameOneOf, `)`)
			p.P(`mtype := `, p.refectPkg.Use(), `.TypeOf((*foo)(nil)).Elem()`)
			p.P(`fieldtype := `, p.refectPkg.Use(), `.TypeOf(`, oneOfVarName, `.`, fieldNameOneOf, `)`)

			p.P(`if `, oneOfVarName, `!=nil && fieldtype.Implements(mtype) {`)
			p.In()
			p.P(oneOfVarName, `.`, fieldNameOneOf, `=`, oneOfVarName, `.`, fieldNameOneOf, `.Parse(isLevelEnabled).(*`, fieldTypename, `)`)
			p.P(`res.`, fieldName, ` = `, oneOfVarName)
			p.Out()
			p.P(`}`)
		} else {
			p.P(`res.`, fieldName, ` = this.`, fieldName)
		}
		p.Out()
		p.P(`}`)
	}
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateParseCallRow(field *descriptor.FieldDescriptorProto, fieldName, namespace string) {
	fieldTypename := strings.TrimPrefix(field.GetTypeName(), ".")
	fieldTypename = strings.TrimPrefix(fieldTypename, namespace+".")

	if field.IsMessage() {
		p.P(`mtype:=`, p.refectPkg.Use(), `.TypeOf((*foo)(nil)).Elem()`)
		p.P(`fieldtype:=`, p.refectPkg.Use(), `.TypeOf(this.`, fieldName, `)`)
		p.P(`if this.`, fieldName, `!=nil && fieldtype.Implements(mtype){`)
		p.In()
		p.P(`res.`, fieldName, `=this.`, fieldName, `.Parse(isLevelEnabled).(*`, fieldTypename, `)`)
		p.Out()
		p.P(`}`)
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
