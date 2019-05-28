// Copyright (c) 2018 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strconv"

	validator "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	regexPkg    generator.Single
	fmtPkg      generator.Single
	s12protoPkg generator.Single
}

// New creates a new generator plugin for go validator
func New() generator.Plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return "govalidator"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.fmtPkg = p.NewImport("fmt")
	p.regexPkg = p.NewImport("regexp")
	p.s12protoPkg = p.NewImport("github.com/SafetyCulture/s12-proto/protobuf/s12proto")

	for _, msg := range file.Messages() {
		if msg.DescriptorProto.GetOptions().GetMapEntry() {
			continue
		}
		p.generateRegexVars(file, msg)
		p.generateValidateFunction(file, msg)
		p.P()
	}
}

func (p *plugin) generateRegexVars(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())
	for _, field := range message.Field {

		if regex := getRegexValue(field); regex != nil {
			fieldName := p.GetFieldName(message, field)
			p.P(`var `, regexName(ccTypeName, fieldName), ` = `, p.regexPkg.Use(), `.MustCompile(`, "`", *regex, "`", `)`)
		}
	}
}

func (p *plugin) generateValidateFunction(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())
	p.P(`func (m *`, ccTypeName, `) Validate() error {`)
	p.In()

	for _, field := range message.Field {

		if !hasValidationExtensions(field) && !field.IsMessage() {
			continue
		}
		var (
			fieldName    = p.GetOneOfFieldName(message, field)
			variableName = "m." + fieldName
			repeated     = field.IsRepeated()
			nullable     = (gogoproto.IsNullable(field) || !gogoproto.ImportsGoGoProto(file.FileDescriptorProto)) && field.IsMessage()
			optional     = proto.GetBoolExtension(field.Options, validator.E_Optional, false)
		)

		if repeated && hasValidationExtensions(field) {
			p.P(`for _, item := range `, variableName, `{`)
			p.In()
			variableName = "item"
		}

		if field.IsString() {
			p.generateStringValidator(variableName, ccTypeName, fieldName, field, optional)
		} else if field.IsBytes() {
			p.generateBytesValidator(variableName, ccTypeName, fieldName, field, optional)
		} else if isSupportedInt(field) {
			p.generateIntegerValidator(variableName, ccTypeName, fieldName, field, optional)
		} else if field.IsMessage() {
			p.generateInnerMessageValidator(variableName, ccTypeName, fieldName, field, nullable)
		}

		if repeated && hasValidationExtensions(field) {
			p.Out()
			p.P(`}`)
		}
	}

	p.P(`return nil`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateStringValidator(variableName string, ccTypeName string, fieldName string, field *descriptor.FieldDescriptorProto, optional bool) {

	if optional {
		p.P(`if `, variableName, ` != "" {`)
		p.In()
	}

	if regex := getRegexValue(field); regex != nil {
		p.P(`if !`, regexName(ccTypeName, fieldName), `.MatchString(`, variableName, `) {`)
		p.In()
		errorStr := "be a string conforming to regex " + strconv.Quote(*regex)
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}
	if getUUIDValue(field) {
		p.P(`if !`, p.s12protoPkg.Use(), `.IsUUID(`, variableName, `) {`)
		p.In()
		errorStr := "be parsable as a UUID"
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}

	p.generateLengthValidator(variableName, ccTypeName, fieldName, field)

	if optional {
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateBytesValidator(variableName string, ccTypeName string, fieldName string, field *descriptor.FieldDescriptorProto, optional bool) {

	if optional {
		p.P(`if `, variableName, ` != nil && len(`, variableName, `) > 0 {`)
		p.In()
	}

	if getUUIDValue(field) {
		p.P(`if len(`, variableName, `) != `, p.s12protoPkg.Use(), `.UUIDSize {`)
		p.In()
		errorStr := "be exactly 16 bytes long to be a valid UUID"
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}

	p.generateLengthValidator(variableName, ccTypeName, fieldName, field)

	if optional {
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateIntegerValidator(variableName string, ccTypeName string, fieldName string, field *descriptor.FieldDescriptorProto, optional bool) {

	if optional {
		p.P(`if `, variableName, ` != 0 {`)
		p.In()
	}

	if v := getIntGtValue(field); v != nil {
		p.P(`if !(`, variableName, ` > `, v, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be greater than '%d'`, *v)
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}
	if v := getIntGteValue(field); v != nil {
		p.P(`if !(`, variableName, ` >= `, v, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be greater than or equal to '%d'`, *v)
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}
	if v := getIntLtValue(field); v != nil {
		p.P(`if !(`, variableName, ` < `, v, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be less than '%d'`, *v)
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}
	if v := getIntLteValue(field); v != nil {
		p.P(`if !(`, variableName, ` <= `, v, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be less than or equal to '%d'`, *v)
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}

	if optional {
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateLengthValidator(variableName string, ccTypeName string, fieldName string, field *descriptor.FieldDescriptorProto) {
	if v := getLengthGteValue(field); v != nil {
		p.P(`if !(len(`, variableName, `) >= `, v, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`have length greater than or equal to '%d'`, *v)
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}
	if v := getLengthLteValue(field); v != nil {
		p.P(`if !(len(`, variableName, `) <= `, v, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`have length less than or equal to '%d'`, *v)
		p.generateErrorString(variableName, fieldName, errorStr)
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateInnerMessageValidator(variableName string, ccTypeName string, fieldName string, field *descriptor.FieldDescriptorProto, nullable bool) {

	if nullable {
		p.P(`if `, variableName, ` != nil {`)
		p.In()
	} else {
		variableName = "&(" + variableName + ")"
	}

	p.P(`if v, ok := interface{}(`, variableName, `).(`, p.s12protoPkg.Use(), `.Validator); ok {`)
	p.In()

	p.P(`if err := v.Validate(); err != nil {`)
	p.In()
	p.P(`return `, p.s12protoPkg.Use(), `.FieldError("`, fieldName, `", err)`)
	p.Out()
	p.P(`}`)

	p.Out()
	p.P(`}`)

	if nullable {
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateErrorString(variableName string, fieldName string, specificError string) {
	p.P(`return `, p.fmtPkg.Use(), ".Errorf(`", fieldName, `: value '%v' must `, specificError, "`, ", variableName, `)`)
}

func regexName(ccTypeName, fieldName string) string {
	return "_regex_" + ccTypeName + "_" + fieldName
}

func hasValidationExtensions(field *descriptor.FieldDescriptorProto) bool {
	if field.Options != nil {
		validExts := []*proto.ExtensionDesc{
			validator.E_Regex,
			validator.E_Uuid,
			validator.E_IntGt,
			validator.E_IntLt,
			validator.E_IntGte,
			validator.E_IntLte,
			validator.E_LengthGte,
			validator.E_LengthLte,
			validator.E_Optional,
		}
		for _, ext := range validExts {
			if proto.HasExtension(field.Options, ext) {
				return true
			}
		}
	}
	return false
}

func getRegexValue(field *descriptor.FieldDescriptorProto) *string {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_Regex)
		if err == nil && v.(*string) != nil {
			return v.(*string)
		}
	}
	return nil
}

func getUUIDValue(field *descriptor.FieldDescriptorProto) bool {
	return proto.GetBoolExtension(field.Options, validator.E_Uuid, false)
}

func getIntGtValue(field *descriptor.FieldDescriptorProto) *int64 {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_IntGt)
		if err == nil && v.(*int64) != nil {
			return v.(*int64)
		}
	}
	return nil
}

func getIntLtValue(field *descriptor.FieldDescriptorProto) *int64 {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_IntLt)
		if err == nil && v.(*int64) != nil {
			return v.(*int64)
		}
	}
	return nil
}

func getIntGteValue(field *descriptor.FieldDescriptorProto) *int64 {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_IntGte)
		if err == nil && v.(*int64) != nil {
			return v.(*int64)
		}
	}
	return nil
}

func getIntLteValue(field *descriptor.FieldDescriptorProto) *int64 {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_IntLte)
		if err == nil && v.(*int64) != nil {
			return v.(*int64)
		}
	}
	return nil
}

func getLengthLteValue(field *descriptor.FieldDescriptorProto) *int64 {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_LengthLte)
		if err == nil && v.(*int64) != nil {
			return v.(*int64)
		}
	}
	return nil
}

func getLengthGteValue(field *descriptor.FieldDescriptorProto) *int64 {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_LengthGte)
		if err == nil && v.(*int64) != nil {
			return v.(*int64)
		}
	}
	return nil
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
