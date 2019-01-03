package plugin

import (
	logger "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports

	reflectPkg generator.Single
	fmtPkg     generator.Single
	logrusPkg  generator.Single
	s12proto   generator.Single
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
	p.reflectPkg = p.NewImport("reflect")
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

	p.P(`res:=&`, ccTypeName, `{}`)
	for _, field := range message.Field {
		var (
			fieldName      = p.GetFieldName(message, field)
			fieldNameOneOf = p.GetOneOfFieldName(message, field)
		)

		if fieldName != fieldNameOneOf {
			if !hasPayloadLoggerExtensions(field) {
				p.P(`if `, p.reflectPkg.Use(), `.TypeOf(this.`, fieldName, `) == `, p.reflectPkg.Use(), `.TypeOf(&`, ccTypeName, `_`, fieldNameOneOf, `{}){`)
				p.In()
				p.P(`res.`, fieldName, `=this.`, fieldName)
				p.Out()
				p.P(`}`)
			} else {
				p.P(`if isLevelEnabled(`, p.s12proto.Use(), `.Level_`, getLevelValue(field).String(), `) {`)
				p.In()
				p.P(`if `, p.reflectPkg.Use(), `.TypeOf(this.`, fieldName, `) == `, p.reflectPkg.Use(), `.TypeOf(&`, ccTypeName, `_`, fieldNameOneOf, `{}){`)
				p.In()
				p.P(`res.`, fieldName, `=this.`, fieldName)
				p.Out()
				p.P(`}`)
				p.Out()
				p.P(`}`)
			}
			continue
		}

		if !hasPayloadLoggerExtensions(field) {
			p.P(`res.`, fieldName, `=this.`, fieldName)
			continue
		} else {
			p.P(`if isLevelEnabled(`, p.s12proto.Use(), `.Level_`, getLevelValue(field).String(), `) {`)
			p.In()
			p.P(`res.`, fieldName, `=this.`, fieldName)
			p.Out()
			p.P(`}`)
		}
	}
	p.P(`return res`)
	p.Out()
	p.P(`}`)
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
