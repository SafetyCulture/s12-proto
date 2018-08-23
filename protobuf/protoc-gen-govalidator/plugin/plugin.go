// Copyright (c) 2018 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strconv"

	validator "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	regexPkg  generator.Single
	fmtPkg    generator.Single
	errrosPkg generator.Single
	uuidPkg   generator.Single
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
	p.errrosPkg = p.NewImport("github.com/pkg/errors")
	p.uuidPkg = p.NewImport("github.com/satori/go.uuid")

	for _, msg := range file.Messages() {
		if msg.DescriptorProto.GetOptions().GetMapEntry() {
			continue
		}
		p.generateRegexVars(file, msg)
		p.generateValidateFunction(file, msg)

	}
}

func (p *plugin) generateRegexVars(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())
	for _, field := range message.Field {
		valid := getFieldValidatorIfAny(field)
		if valid != nil && valid.Regex != nil {
			fieldName := p.GetFieldName(message, field)
			p.P(`var `, regexName(ccTypeName, fieldName), ` = `, p.regexPkg.Use(), `.MustCompile(`, "`", valid.Regex, "`", `)`)
		}
	}
}

func (p *plugin) generateValidateFunction(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())
	p.P(`func (this *`, ccTypeName, `) Validate() error {`)
	p.In()

	for _, field := range message.Field {
		fieldValidator := getFieldValidatorIfAny(field)
		if fieldValidator == nil && !field.IsMessage() {
			continue
		}
		var (
			fieldName    = p.GetOneOfFieldName(message, field)
			variableName = "this." + fieldName
		)

		if field.IsString() {
			p.generateStringValidator(variableName, ccTypeName, fieldName, fieldValidator)
		} else if field.IsBytes() {
			p.generateBytesValidator(variableName, ccTypeName, fieldName, fieldValidator)
		} else if isSupportedInt(field) {
			p.generateIntegerValidator(variableName, ccTypeName, fieldName, fieldValidator)
		}
	}

	p.P(`return nil`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateStringValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator) {
	if fv.Regex != nil {
		p.P(`if !`, regexName(ccTypeName, fieldName), `.MatchString(`, variableName, `) {`)
		p.In()
		errorStr := "be a string conforming to regex " + strconv.Quote(fv.GetRegex())
		p.generateErrorString(variableName, fieldName, errorStr, fv)
		p.Out()
		p.P(`}`)
	}
	if fv.Uuid != nil && *fv.Uuid {
		p.P(`if _, err := `, p.uuidPkg.Use(), `.FromString(`, variableName, `); err != nil {`)
		p.In()
		errorStr := "be a parsable as a UUID"
		p.generateErrorString(variableName, fieldName, errorStr, fv)
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateBytesValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator) {
	if fv.Uuid != nil && *fv.Uuid {
		p.P(`if _, err := `, p.uuidPkg.Use(), `.FromBytes(`, variableName, `); err != nil {`)
		p.In()
		errorStr := "be a parsable as a UUID"
		p.generateErrorString(variableName, fieldName, errorStr, fv)
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateIntegerValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator) {
	if fv.IntGt != nil {
		p.P(`if !(`, variableName, ` > `, fv.IntGt, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be greater than '%d'`, fv.GetIntGt())
		p.generateErrorString(variableName, fieldName, errorStr, fv)
		p.Out()
		p.P(`}`)
	}
	if fv.IntGte != nil {
		p.P(`if !(`, variableName, ` >= `, fv.IntGte, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be greater than or equal to '%d'`, fv.GetIntGte())
		p.generateErrorString(variableName, fieldName, errorStr, fv)
		p.Out()
		p.P(`}`)
	}
	if fv.IntLt != nil {
		p.P(`if !(`, variableName, ` < `, fv.IntLt, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be less than '%d'`, fv.GetIntLt())
		p.generateErrorString(variableName, fieldName, errorStr, fv)
		p.Out()
		p.P(`}`)
	}
	if fv.IntLte != nil {
		p.P(`if !(`, variableName, ` <= `, fv.IntLte, `) {`)
		p.In()
		errorStr := fmt.Sprintf(`be less than or equal to '%d'`, fv.GetIntLte())
		p.generateErrorString(variableName, fieldName, errorStr, fv)
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateErrorString(variableName string, fieldName string, specificError string, fv *validator.FieldValidator) {
	p.P(`return `, p.errrosPkg.Use(), ".Errorf(`", fieldName, `: value '%s' must `, specificError, "`, ", variableName, `)`)
}

func regexName(ccTypeName, fieldName string) string {
	return "_regex_" + ccTypeName + "_" + fieldName
}

func getFieldValidatorIfAny(field *descriptor.FieldDescriptorProto) *validator.FieldValidator {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, validator.E_Field)
		if err == nil && v.(*validator.FieldValidator) != nil {
			return (v.(*validator.FieldValidator))
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
