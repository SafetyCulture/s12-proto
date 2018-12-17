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

	p.generateLevelToLogrusFunction()
}

func (p *plugin) generateParseFunction(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())
	p.P(`func (this *`, ccTypeName, `) Parse(isLevelEnabled func(level `, p.logrusPkg.Use(), `.Level) bool) proto.Message {`)
	p.In()

	p.P(`res:=&`, ccTypeName, `{}`)
	for _, field := range message.Field {
		var (
			fieldName = p.GetOneOfFieldName(message, field)
		)

		if !hasPayloadLoggerExtensions(field) {
			p.P(`res.`, fieldName, `=this.`, fieldName)

			continue
		}

		p.P(`if isLevelEnabled(levelToLogrus(`, p.s12proto.Use(), `.Level_`, getLevelLteValue(field).String(), `)) {`)
		p.In()
		p.P(`res.`, fieldName, `=this.`, fieldName)
		p.Out()
		p.P(`}`)
	}
	p.P(`return res`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateLevelToLogrusFunction() {
	p.P(`func levelToLogrus(level `, p.s12proto.Use(), `.Level) `, p.logrusPkg.Use(), `.Level {`)
	p.In()
	p.P(`switch level {`)
	p.In()
	p.P(`case `, p.s12proto.Use(), `.Level_PANIC:`)
	p.P(`return `, p.logrusPkg.Use(), `.PanicLevel`)
	p.P(`case `, p.s12proto.Use(), `.Level_FATAL:`)
	p.P(`return `, p.logrusPkg.Use(), `.FatalLevel`)
	p.P(`case `, p.s12proto.Use(), `.Level_ERROR:`)
	p.P(`return `, p.logrusPkg.Use(), `.ErrorLevel`)
	p.P(`case `, p.s12proto.Use(), `.Level_WARN:`)
	p.P(`return `, p.logrusPkg.Use(), `.WarnLevel`)
	p.P(`case `, p.s12proto.Use(), `.Level_INFO:`)
	p.P(`return `, p.logrusPkg.Use(), `.InfoLevel`)
	p.P(`case `, p.s12proto.Use(), `.Level_DEBUG:`)
	p.P(`return `, p.logrusPkg.Use(), `.DebugLevel`)
	p.P(`case `, p.s12proto.Use(), `.Level_TRACE:`)
	p.P(`return `, p.logrusPkg.Use(), `.TraceLevel`)
	p.Out()
	p.P(`}`)
	p.P(`return 0`)
	p.Out()
	p.P(`}`)
}

func getLevelLteValue(field *descriptor.FieldDescriptorProto) *logger.Level {
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
