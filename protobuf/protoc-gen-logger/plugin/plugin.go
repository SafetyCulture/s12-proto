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
	p.reflectPkg = p.NewImport("reflect")

	for _, msg := range file.Messages() {
		p.generateParseFunction(file, msg)
	}
}

func (p *plugin) generateParseFunction(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())

	p.P(`func (this *`, ccTypeName, `) LogPayload(logger interface {`)
	p.In()
	p.P(`Debug(args ...interface{})`)
	p.P(`Info(args ...interface{})`)
	p.P(`Warn(args ...interface{})`)
	p.P(`Error(args ...interface{})`)
	p.P(`Fatal(args ...interface{})`)
	p.P(`Panic(args ...interface{})`)
	p.Out()
	p.P(`}){`)
	p.In()

	for _, field := range message.Field {
		var (
			fieldName      = p.GetFieldName(message, field)
			fieldNameOneOf = p.GetOneOfFieldName(message, field)
		)
		if fieldName != fieldNameOneOf {
			if hasPayloadLoggerExtensions(field) {
				p.P(`if `, p.reflectPkg.Use(), `.TypeOf(this.`, fieldName, `) == `, p.reflectPkg.Use(), `.TypeOf(&`, ccTypeName, `_`, fieldNameOneOf, `{}){`)
				p.In()
				switch *getLevelValue(field) {
				case logger.Level_PANIC:
					p.P(`logger.Panic(this.`, fieldName, `)`)
					break
				case logger.Level_DEBUG:
					p.P(`logger.Debug(this.`, fieldName, `)`)
					break
				case logger.Level_ERROR:
					p.P(`logger.Error(this.`, fieldName, `)`)
					break
				case logger.Level_FATAL:
					p.P(`logger.Fatal(this.`, fieldName, `)`)
					break
				case logger.Level_INFO:
					p.P(`logger.Info(this.`, fieldName, `)`)
					break
				case logger.Level_WARN:
					p.P(`logger.Warn(this.`, fieldName, `)`)
					break
				}
				p.Out()
				p.P(`}`)
			}
			continue
		}

		if hasPayloadLoggerExtensions(field) {
			switch *getLevelValue(field) {
			case logger.Level_PANIC:
				p.P(`logger.Panic(this.`, fieldName, `)`)
				break
			case logger.Level_DEBUG:
				p.P(`logger.Debug(this.`, fieldName, `)`)
				break
			case logger.Level_ERROR:
				p.P(`logger.Error(this.`, fieldName, `)`)
				break
			case logger.Level_FATAL:
				p.P(`logger.Fatal(this.`, fieldName, `)`)
				break
			case logger.Level_INFO:
				p.P(`logger.Info(this.`, fieldName, `)`)
				break
			case logger.Level_WARN:
				p.P(`logger.Warn(this.`, fieldName, `)`)
				break
			}
		}
	}
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
