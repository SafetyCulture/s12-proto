// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	validator "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
)

// Standard library dependencies.
const (
	fmtPackage     = protogen.GoImportPath("fmt")
	mathPackage    = protogen.GoImportPath("math")
	regexpPackage  = protogen.GoImportPath("regexp")
	stringsPackage = protogen.GoImportPath("strings")
)

// Other library dependencies.
const (
	s12protoPackage = protogen.GoImportPath("github.com/SafetyCulture/s12-proto/s12/protobuf/proto")
)

// GenerateFile generates the validator.pb.go file
func GenerateFile(p *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".validator.pb.go"
	g := p.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-govalidator. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	// Generate a blank import statement if not used by any validation functions,
	// as we still want to register the protobuf validation types for the Reflection API.
	g.Import(s12protoPackage)

	for _, msg := range file.Messages {
		if opts, ok := msg.Desc.Options().(*descriptorpb.MessageOptions); !ok || opts.GetMapEntry() {
			continue
		}
		genRegexVars(g, msg)
		genValidateFunc(g, msg)
		g.P()
	}

	return g
}

func genRegexVars(g *protogen.GeneratedFile, msg *protogen.Message) {
	for _, field := range msg.Fields {
		if regex := getRegexValue(field); regex != "" {
			g.P("const ", "_regex_val_", field.GoIdent, " = `", regex, "`")
			g.P("var ", "_regex_", field.GoIdent, " = ", regexpPackage.Ident("MustCompile"), "(_regex_val_", field.GoIdent, ")")
		}
	}
	for _, innerMsg := range msg.Messages {
		genRegexVars(g, innerMsg)
	}
}

func genValidateFunc(g *protogen.GeneratedFile, msg *protogen.Message) {
	g.P("func (m *", msg.GoIdent.GoName, `) Validate() error {`)
	for _, f := range msg.Fields {
		hasExt := hasValidationExtensions(f)
		if !hasExt && f.Message == nil {
			continue
		}

		varName := "m." + f.GoName
		repeated := f.Desc.Cardinality() == protoreflect.Repeated

		if hasExt && repeated {
			g.P("for _, item := range ", varName, "{")
			varName = "item"
		}

		if f.Desc.IsMap() {
			g.P(`// Validation of proto3 map<> fields is unsupported.`)
			continue
		}

		if f.Oneof != nil {
			g.P(`// Validation of oneof fields is unsupported.`)
			continue
		}

		switch f.Desc.Kind() {
		case protoreflect.StringKind:
			genStringValidator(g, f, varName)
		case protoreflect.BytesKind:
			genBytesValidator(g, f, varName)
		case protoreflect.Int32Kind, protoreflect.Int64Kind,
			protoreflect.Uint32Kind, protoreflect.Uint64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind:
			genIntValidator(g, f, varName)
		case protoreflect.MessageKind:
			genMsgValidator(g, f, varName)
		}

		if hasExt && repeated {
			g.P(`}`)
		}
	}

	g.P(`return nil`)
	g.P(`}`)

	for _, innerMsg := range msg.Messages {
		if innerMsg.Desc.IsMapEntry() {
			// Do not generate validation func for map entry type
			continue
		}
		g.P()
		genValidateFunc(g, innerMsg)
	}
}

func genStringValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	optional := getBoolExtension(f, validator.E_Optional)

	if optional {
		g.P("if ", varName, " != \"\" {")
	}

	if regex := getRegexValue(f); regex != "" {
		g.P("if !_regex_", f.GoIdent, ".MatchString(", varName, ") {")
		errStr := "be a string conforming to regex '%s'"
		genErrorStringWithParams(g, varName, string(f.Desc.Name()), errStr, "_regex_val_"+f.GoIdent.GoName)
		g.P("}")
	}

	if getBoolExtension(f, validator.E_Uuid) {
		g.P("if !", s12protoPackage.Ident("IsUUID"), "(", varName, ") {")
		errStr := "be parsable as a UUID"
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	if getBoolExtension(f, validator.E_LegacyId) {
		g.P("if !", s12protoPackage.Ident("IsLegacyID"), "(", varName, ") {")
		errStr := "be parsable as a UUID or a legacy ID"
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	genLenValidator(g, f, varName)

	if optional {
		g.P("}")
	}
}

func genLenValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	if getBoolExtension(f, validator.E_TrimLenCheck) {
		g.P("_trim_", f.GoIdent, " := ", stringsPackage.Ident("TrimSpace"), "(", varName, ")")
		varName = "_trim_" + f.GoIdent.GoName
		g.P("_ = ", varName)
	}

	if v := getIntExtention(f, validator.E_LengthGte); v >= 0 {
		g.P("if !(len(", varName, ") >= ", v, ") {")
		errStr := fmt.Sprintf(`have length greater than or equal to '%d'`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	if v := getIntExtention(f, validator.E_LengthLte); v >= 0 {
		g.P("if !(len(", varName, ") <= ", v, ") {")
		errStr := fmt.Sprintf(`have length less than or equal to '%d'`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}
}

func genBytesValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	optional := getBoolExtension(f, validator.E_Optional)

	if optional {
		g.P("if ", varName, " != nil && len(", varName, ") > 0 {")
	}

	if getBoolExtension(f, validator.E_Uuid) {
		g.P("if len(", varName, ") != ", s12protoPackage.Ident("UUIDSize"), " {")
		errStr := "be exactly 16 bytes long to be a valid UUID"
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	genLenValidator(g, f, varName)

	if optional {
		g.P("}")
	}
}

func genIntValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	optional := getBoolExtension(f, validator.E_Optional)

	if optional {
		g.P("if ", varName, " != 0 {")
	}

	if v := getIntExtention(f, validator.E_IntGt); v >= 0 {
		g.P("if !(", varName, " > ", v, ") {")
		errStr := fmt.Sprintf(`be greater than '%d'`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}
	if v := getIntExtention(f, validator.E_IntGte); v >= 0 {
		g.P("if !(", varName, " >= ", v, ") {")
		errStr := fmt.Sprintf(`be greater than or equal to '%d'`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}
	if v := getIntExtention(f, validator.E_IntLt); v >= 0 {
		g.P("if !(", varName, " < ", v, ") {")
		errStr := fmt.Sprintf(`be less than '%d'`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}
	if v := getIntExtention(f, validator.E_IntLte); v >= 0 {
		g.P("if !(", varName, " <= ", v, ") {")
		errStr := fmt.Sprintf(`be less than or equal to '%d'`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	if optional {
		g.P("}")
	}
}

func genMsgValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	if getBoolExtension(f, validator.E_MsgRequired) {
		g.P("if ", varName, " == nil {")
		g.P("return ", fmtPackage.Ident("Errorf"), "(\"field ", f.Desc.Name(), " is required\")")
		g.P("}")
	}

	g.P("if ", varName, " != nil {")
	g.P("if v, ok := interface{}(", varName, ").(", s12protoPackage.Ident("Validator"), "); ok {")
	g.P("if err := v.Validate(); err != nil {")
	g.P("return ", s12protoPackage.Ident("FieldError"), "(\"", f.Desc.Name(), "\", err)")
	g.P("}")
	g.P("}")
	g.P("}")
}

func genErrorString(g *protogen.GeneratedFile, varName, fieldName, specificErr string) {
	g.P(`return `, fmtPackage.Ident("Errorf"), "(`", fieldName, `: value '%v' must `, specificErr, "`, ", varName, `)`)
}

func genErrorStringWithParams(g *protogen.GeneratedFile, varName, fieldName, specificErr string, errParams ...string) {
	g.P(`return `, fmtPackage.Ident("Errorf"), "(`", fieldName, `: value '%v' must `, specificErr, "`, ", varName, ", ", strings.Join(errParams, ", "), `)`)
}

var validExts = []protoreflect.ExtensionType{
	validator.E_Regex,
	validator.E_Uuid,
	validator.E_IntGt,
	validator.E_IntLt,
	validator.E_IntGte,
	validator.E_IntLte,
	validator.E_LengthGte,
	validator.E_LengthLte,
	validator.E_Optional,
	validator.E_MsgRequired,
	validator.E_LegacyId,
	validator.E_TrimLenCheck,
}

func hasValidationExtensions(f *protogen.Field) bool {
	if opts := f.Desc.Options(); opts != nil {
		for _, ext := range validExts {
			if proto.HasExtension(opts, ext) {
				return true
			}
		}
	}
	return false
}

func getRegexValue(f *protogen.Field) string {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, validator.E_Regex)
		if v, ok := ext.(string); ok {
			return v
		}
	}
	return ""
}

func getBoolExtension(f *protogen.Field, xt protoreflect.ExtensionType) bool {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(bool); ok {
			return v
		}
	}
	return false
}

func getIntExtention(f *protogen.Field, xt protoreflect.ExtensionType) int64 {
	if opts := f.Desc.Options(); opts != nil {
		if !proto.HasExtension(opts, xt) {
			return -1
		}
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(int64); ok {
			return v
		}
	}
	return -1
}
