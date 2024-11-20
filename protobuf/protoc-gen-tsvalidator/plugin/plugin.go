// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	validator "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
)

// Standard library dependencies.
const (
	fmtPackage       = protogen.GoImportPath("fmt")
	timePackage      = protogen.GoImportPath("time")
	mathPackage      = protogen.GoImportPath("math")
	regexpPackage    = protogen.GoImportPath("regexp")
	stringsPackage   = protogen.GoImportPath("strings")
	utfPackage       = protogen.GoImportPath("unicode/utf8")
	normPackage      = protogen.GoImportPath("golang.org/x/text/unicode/norm")
	transformPackage = protogen.GoImportPath("golang.org/x/text/transform")
)

// Other library dependencies.
const (
	s12protoPackage = protogen.GoImportPath("github.com/SafetyCulture/s12-proto/s12/protobuf/proto")
)

// Some global vars for regex generations
var regexHashLib = make(map[string]struct{})

var dependencies = make(map[*protogen.Message][]*protogen.Message)
var tsVars = make(map[string]bool)

// Validator plugin version
var validatorVersion = "v2.6.0"

// Write a preamble in the auto generated files
func genGeneratedHeader(gen *protogen.Plugin, g *protogen.GeneratedFile, f *protogen.File) {
	g.P("// Code generated by protoc-gen-tsvalidator. DO NOT EDIT.")
	g.P("// versions:")
	protocVersion := "(unknown)"
	if v := gen.Request.GetCompilerVersion(); v != nil {
		protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
		if s := v.GetSuffix(); s != "" {
			protocVersion += "-" + s
		}
	}
	g.P("// \tprotoc-gen-tsvalidator ", validatorVersion)
	g.P("// \tprotoc                 ", protocVersion)

	if f.Proto.GetOptions().GetDeprecated() {
		g.P("// ", f.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", f.Desc.Path())
	}
	g.P()
}

func TopologicalSort(dependencies map[*protogen.Message][]*protogen.Message) ([]*protogen.Message, error) {
	sorted := make([]*protogen.Message, 0, len(dependencies))
	visited := make(map[*protogen.Message]bool)
	temp := make(map[*protogen.Message]bool)

	var visit func(m *protogen.Message) error
	visit = func(m *protogen.Message) error {
		if temp[m] {
			// log.Printf("circular dependency detected: %s ", m)
			return nil
		}
		if !visited[m] {
			temp[m] = true
			for _, dep := range dependencies[m] {
				if err := visit(dep); err != nil {
					return err
				}
			}
			visited[m] = true
			temp[m] = false
			sorted = append(sorted, m)
		}
		return nil
	}

	for m := range dependencies {
		if !visited[m] {
			if err := visit(m); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}

// GenerateFile generates the validator.pb.go file
func GenerateFile(p *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	p.SupportedFeatures |= uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	filename := file.GeneratedFilenamePrefix + ".validator.ts"
	g := p.NewGeneratedFile(filename, file.GoImportPath)

	// Add preamble and version for validator.pb.go
	genGeneratedHeader(p, g, file)

	g.P()
	g.P(`import { z } from "zod";`)
	g.P()

	// Generate a blank import statement if not used by any validation functions,
	// as we still want to register the protobuf validation types for the Reflection API.
	g.Import(s12protoPackage)

	genStringGenerics()
	for _, msg := range file.Messages {
		if opts, ok := msg.Desc.Options().(*descriptorpb.MessageOptions); !ok || opts.GetMapEntry() {
			continue
		}
		updateDependencies(msg)
		// genValidateFunc(g, msg)
		g.P()
	}

	orderedMessages, _ := TopologicalSort(dependencies)
	// log.Printf("orderedMessages: %s", orderedMessages)

	for _, msg := range orderedMessages {
		if opts, ok := msg.Desc.Options().(*descriptorpb.MessageOptions); !ok || opts.GetMapEntry() {
			continue
		}
		// log.Printf("msg: %s", msg.GoIdent.GoName)
		genValidateFunc(g, msg)
		g.P()
	}

	return g
}

func updateDependencies(msg *protogen.Message) {

	dependant := false
	for _, f := range msg.Fields {
		if f.Desc.Kind() == protoreflect.MessageKind {
			if dependencies[msg] == nil {
				dependencies[msg] = []*protogen.Message{}
			}
			dependant = true
			dependencies[msg] = append(dependencies[msg], f.Message)
			// log.Printf("dependencies: %s", dependencies)

		}
	}
	if !dependant {
		dependencies[msg] = []*protogen.Message{}
	}
	// log.Printf("dependencies: %s", dependencies)

}

// Generic preparations/initialisations for the string validator
func genStringGenerics() {

	// Prepare the regex pattern to allow replaced chars
	prepareStringReplacerRegex(stringUnsafeReplacerMap, "replacer_unsafe_allowed")
	prepareStringReplacerRegex(stringSymbolReplacerMap, "replacer_symbol_allowed")

	// Prepare the default `unsafe_string` pattern by adding tokens from `string`
	stringReDefaultUnsafe = append(stringReDefaultSafe, stringReDefaultUnsafe...)
}

func genValidateFunc(g *protogen.GeneratedFile, msg *protogen.Message) {
	validatorConst := `export const ` + msg.GoIdent.GoName + `Validator = z.object({`
	// g.P(validator + "z.string();")
	hasAnyExt := false
	var messagePropsMap = make(map[string]string)

	for _, f := range msg.Fields {
		hasExt := hasValidationExtensions(f)
		if !hasExt && f.Message == nil {
			continue
		}
		hasAnyExt = true

		varName := f.GoIdent.GoName
		hasRepeatedExt := hasExt && f.Desc.Cardinality() == protoreflect.Repeated

		validatorAttrConst := "export const " + varName + " = "

		if f.Desc.IsMap() {
			g.P(validatorAttrConst + "z.object({})")
			messagePropsMap[f.Desc.JSONName()] = varName
			continue
		}

		// // exclude synthetic oneof fields, e.g. fields with the optional keyword
		// if f.Oneof != nil && !f.Oneof.Desc.IsSynthetic() {
		// 	g.P("if x, ok := m.", f.Oneof.GoName, ".(*", g.QualifiedGoIdent(f.GoIdent), "); ok {")
		// 	varName = "x." + f.GoName
		// }

		switch f.Desc.Kind() {
		case protoreflect.StringKind:
			// Run both legacy and newer validation logic as both options can be combined at this time
			// Should exclusively use newer validation logic (validator.string/unsafe_string) for newly added fields
			// For existing fields, should be fairly easy to update the rules to the new validator
			stringValidator := ".string()"
			genLegacyStringValidator(g, f, &stringValidator)
			genStringValidator(g, f, &stringValidator)
			genIdValidator(g, f, &stringValidator)
			genEmailValidator(g, f, &stringValidator)
			genURLValidator(g, f, &stringValidator)
			genTimezoneValidator(g, f, &stringValidator)
			genSimpleStringValidator(g, f, &stringValidator)
			g.P(validatorAttrConst + "z" + stringValidator)
			messagePropsMap[f.Desc.JSONName()] = varName
		case protoreflect.BytesKind:
			// IdValidator not supported for bytes at this point
			g.P(validatorAttrConst + "z" + genBytesValidator(g, f))
			messagePropsMap[f.Desc.JSONName()] = varName
		case protoreflect.Int32Kind, protoreflect.Int64Kind,
			protoreflect.Uint32Kind, protoreflect.Uint64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind:
			g.P(validatorAttrConst + "z" + genNumberValidator(g, f, varName, true) + ";")
			messagePropsMap[f.Desc.JSONName()] = varName
		case protoreflect.MessageKind:
			g.P(validatorAttrConst + genMsgValidator(g, f, varName))
		case protoreflect.EnumKind:
			g.P(validatorAttrConst + "z" + genEnumValidator(g, f, varName) + ";")
			messagePropsMap[f.Desc.JSONName()] = varName
		case protoreflect.DoubleKind, protoreflect.FloatKind:
			g.P(validatorAttrConst + "z" + genNumberValidator(g, f, varName, false) + ";")
			messagePropsMap[f.Desc.JSONName()] = varName
		default:
			g.P(validatorAttrConst + "z.any();")
			messagePropsMap[f.Desc.JSONName()] = varName
		}

		if hasRepeatedExt && messagePropsMap[f.Desc.JSONName()] != "" {
			arrayName := f.GoIdent.GoName + "Array"
			arrConst := "export const " + arrayName + " = " + varName + ".array()"
			if v := getIntExtention(f, validator.E_RepeatedLenGte); v >= 0 {
				arrConst += ".min(" + strconv.FormatInt(v, 10) + ")"
			}
			if v := getIntExtention(f, validator.E_RepeatedLenLte); v >= 0 {
				arrConst += ".max(" + strconv.FormatInt(v, 10) + ")"
			}
			if f.Desc.HasOptionalKeyword() {
				arrConst += ".optional()"
			}
			g.P(arrConst + ";")
			messagePropsMap[f.Desc.JSONName()] = arrayName
		}

		// if f.Oneof != nil && !f.Oneof.Desc.IsSynthetic() {
		// 	g.P("}")
		// }

	}
	if hasAnyExt {
		for key, value := range messagePropsMap {
			validatorConst += fmt.Sprintf(`%s: %s,`, key, value)
		}
		g.P(validatorConst + "})")
		g.P()
	}

}

func genLegacyStringValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {
	optional := getBoolExtension(f, validator.E_Optional)

	if regex := getRegexValue(f); regex != "" {
		*stringValidator += ".regex(/" + regex + "/)"
	}

	// Legacy UUIDv4 validator
	if getBoolExtension(f, validator.E_Uuid) {
		*stringValidator += ".uuid()"
	}

	// if getBoolExtension(f, validator.E_LegacyId) {
	// 	g.P("if !", s12protoPackage.Ident("IsLegacyID"), "(", varName, ", false) {")
	// 	errStr := "be parsable as a UUID or a legacy ID"
	// 	genErrorString(g, varName, string(f.Desc.Name()), errStr)
	// 	g.P("}")
	// }

	genLenValidator(g, f, stringValidator)

	if optional {
		*stringValidator += ".optional()"
	}

}

func genSimpleStringValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {
	rules := getSimpleStringExtension(f, validator.E_SimpleString)
	if rules == nil {
		return
	}

	// "optional" can be provided via the legacy directive or set as option on the `string`/`unsafe_string` validator
	// To allow deprecation of the 'old' directive later, we duplicate the logic here

	if rules.GetMinLen() >= 1 {
		*stringValidator += ".min(" + strconv.Itoa(int(rules.GetMinLen())) + ")"
	}

	if rules.GetMaxLen() >= 1 {
		*stringValidator += ".max(" + strconv.Itoa(int(rules.GetMaxLen())) + ")"
	}
	if rules.GetOptional() {
		*stringValidator += ".optional()"
	}
}

func genStringValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {
	// `string` and `unsafe_string` share most of the underlying validation logic

	stringType := validator.E_String
	rules := getStringExtension(f, stringType)
	if rules == nil {
		stringType = validator.E_UnsafeString
		rules = getStringExtension(f, stringType)
	}
	if rules == nil {
		// not E_String or E_UnsafeString
		return
	}

	// // Prepare to allow list regex
	// allowListReId := "re_" + f.GoIdent.GoName
	// if stringType == validator.E_String {
	// 	prepareRegex(allowListReId, stringReDefaultSafe...)
	// } else {
	// 	prepareRegex(allowListReId, stringReDefaultUnsafe...)
	// }

	// "optional" can be provided via the legacy directive or set as option on the `string`/`unsafe_string` validator
	// To allow deprecation of the 'old' directive later, we duplicate the logic here

	// // ### STEP 1: decode
	// // TODO PA: Check if we can confirm that data has been decoded, things like %20 might indicate otherwise

	// // ### STEP 2: normalise and canonicalise
	// // Normalisation for file paths/URLs should not happen here but in URL validator

	// // ##### 2A. Normalise Unicode NFD strings to NFC
	// // Convert to unicode NFC if not already in that format
	// //  e.g. NFD string man\u0303ana should be normalised to ma\u00F1ana (n + \u0303 = \u00F1)
	// // TODO PA: Check if we need support for NFKC/NFKD
	// // TODO PA: Do we want to have support for U+202E RIGHT-TO-LEFT OVERRIDE by default for all strings?
	// g.P("if !", normPackage.Ident("NFC.IsNormalString"), "(", varName, ") && ", normPackage.Ident("NFD.IsNormalString"), "(", varName, ") {")
	// g.P("// normalise NFD to NFC string")
	// g.P("var normErr error")
	// g.P(varName, ", _, normErr = ", transformPackage.Ident("String"), "(", transformPackage.Ident("Chain"), "(", normPackage.Ident("NFD"), ",", normPackage.Ident("NFC"), ")", ",", varName, ")")
	// g.P("if normErr != nil {")
	// if rules.GetLogOnly() {
	// 	printErrorString(g, varName, string(f.Desc.Name()), `must be normalisable to NFC`, 50)
	// } else {
	// 	genErrorString(g, varName, string(f.Desc.Name()), `must be normalisable to NFC`)
	// }
	// g.P("}")
	// g.P("}")

	// // ##### 2B. Check for encoding issues
	// if rules.GetValidateEncoding() {
	// 	// U+FFFD in character indicates that the encoding was likely incorrect, resulting in � U+FFFD (RuneError)
	// 	// Test string: $\xa35 for Pepp\xe9 which is encoded as latin1
	// 	g.P("if ", stringsPackage.Ident("ContainsRune"), "(", varName, ", ", utfPackage.Ident("RuneError"), ") {")
	// 	errStr := `must have valid encoding`
	// 	if rules.GetLogOnly() {
	// 		printErrorString(g, varName, string(f.Desc.Name()), errStr, 50)
	// 	} else {
	// 		genErrorString(g, varName, string(f.Desc.Name()), errStr)
	// 	}

	// 	// Also validate using utf8.ValidString: reports whether string consists entirely of valid UTF-8-encoded runes
	// 	g.P("} else if !", utfPackage.Ident("ValidString"), "(", varName, ") {")
	// 	errStr = `must be a valid UTF-8-encoded string`
	// 	if rules.GetLogOnly() {
	// 		printErrorString(g, varName, string(f.Desc.Name()), errStr, 50)
	// 	} else {
	// 		genErrorString(g, varName, string(f.Desc.Name()), errStr)
	// 	}
	// 	g.P("}")
	// } else {
	// 	// Allow RuneError as valid character if we do not check for invalid encoding otherwise the string might be denied
	// 	// U+FFFD is in Symbol_So category which must be allowed in case validateInvalidEncoding is disabled (to detect issues)
	// 	prepareRegex(allowListReId, `\x{FFFD}`)
	// }

	// // ### STEP 3: sanitise
	// // ##### 3A. sanitise whitespace (trim option)
	// // Note this will PERMANENTLY mutate the message field data, this trim is not for length check only but for sanitisation
	// // WARNING: Any leading/trailing whitespace will be permanently removed (this is intended here)
	// if rules.GetTrim() {
	// 	g.P(varName, " = ", stringsPackage.Ident("TrimSpace"), "(", varName, ")")
	// }

	// // ##### 3B. replace restricted characters with a safe alternative
	// // WARNING: this could corrupt data, eg. hyperlinks in a description box so disabled by default
	// // Ensure to test the alternative characters in all output (database, email, pdf, csv, web UI, mobile UI)
	// if rules.GetReplaceUnsafe() {
	// 	// If allow option is defined, only chars from allow list will be allowed
	// 	// Otherwise, all chars in the replace map are accepted
	// 	// If none of the allow chars are in the unsafe map, there is no work to be done here
	// 	// throw a warning that the option is not required in that case
	// 	if rules.GetAllow() != "" {

	// 		for _, runeValue := range rules.GetAllow() {
	// 			unicodeValue := fmt.Sprintf("%U", runeValue)
	// 			unicodeValue = strings.Replace(unicodeValue, "U+", "", 1)
	// 			if _, isUnsafe := stringUnsafeReplacerMap[`\u`+unicodeValue]; isUnsafe {
	// 				// need to replace this one, replace case by case instead of using the replacer
	// 				replacedUnicodeValue := stringUnsafeReplacerMap[`\u`+unicodeValue]
	// 				g.P(varName, " = ", stringsPackage.Ident("ReplaceAll"), "(", varName, `, "\u`, unicodeValue, `", "`, replacedUnicodeValue, `")`)

	// 				// and add it to the regex
	// 				prepareRegex(allowListReId, `\x{`+strings.Replace(replacedUnicodeValue, `\u`, ``, 1)+`}`)
	// 			}
	// 		}
	// 	} else {
	// 		// Add the possible replaced characters to the allow list as they are not all allowed by default
	// 		mergeRegex(allowListReId, "replacer_unsafe_allowed")
	// 		// Replace the unsafe values
	// 		g.P(varName, " = ", s12protoPackage.Ident("UnsafeCharReplacer"), ".Replace(", varName, ")")
	// 	}
	// }

	// // ##### 3C. strip carriage return characters \r in multiline strings (leave \n)
	// if rules.GetMultiline() {
	// 	// Generates: varName = strings.ReplaceAll(varName, "\r", "")
	// 	g.P(varName, " = ", stringsPackage.Ident("ReplaceAll"), "(", varName, ", \"\\r\", \"\")")
	// }

	// // ##### 3D. replace rare symbols with a more common alternative
	// if rules.GetReplaceOther() {
	// 	// Add the possible replaced characters to the allow list as they might not be allowed by default
	// 	mergeRegex(allowListReId, "replacer_symbol_allowed")
	// 	// Replace the values
	// 	if rules.GetMultiline() {
	// 		// Should not replace \n as the normal symbol replacer does otherwise multiline breaks
	// 		g.P(varName, " = ", s12protoPackage.Ident("SymbolCharReplacerMultiline"), ".Replace(", varName, ")")
	// 	} else {
	// 		g.P(varName, " = ", s12protoPackage.Ident("SymbolCharReplacer"), ".Replace(", varName, ")")

	// 	}
	// }

	// // ##### 3E. Sanitise (remove) Private Use Area Codepoints in the Basic Multilingual Plane
	// // Note that we do not check for PUA in planes 15 and 16 currently
	// if rules.GetSanitisePua() {
	// 	g.P(varName, "= ", s12protoPackage.Ident("RegexPua"), ".ReplaceAllString(", varName, ", \"\")")
	// }

	// ### STEP 4: validate
	// ##### 4A. validate length
	// Length could be checked as part of the regex in next step but currently implemented
	//  as a separate step before testing against the regex. This might be more efficient
	//  for large strings/invalid data as we reject it early. It also allows for regex reuse.

	// Always have a min and max length: either default vals or set in the field option
	// This ensures that we no longer accept any length string (safe by default)
	minLen := stringLenMinDefault
	maxLen := stringLenMaxDefault

	// Check if friendly format length defined ("x:y"), overrides the above
	if fLen := rules.GetLen(); fLen != "" {
		// Check for a single value or min:max notation
		fLenChunks := strings.SplitN(fLen, ":", 3)
		fMinLen, fMaxLen := -1, -1 // set -1 as default to distinguish between 0 and unset
		switch len(fLenChunks) {
		case 2:
			// min:max notation deffiend, e.g. len: "2:20"
			// Try casting from string to int
			fMinLen, _ = strconv.Atoi(fLenChunks[0])
			fMaxLen, _ = strconv.Atoi(fLenChunks[1])
			if fMinLen == 0 {
				// Use default min length (not unlimited) for missing min value, e.g. :X
				fMinLen = int(minLen)
			}
			if fMaxLen == 0 {
				// Use default min length (not unlimited) for missing max value, e.g. X:
				fMaxLen = int(maxLen)
			}
			if fMaxLen <= fMinLen || fMinLen < 0 || fMaxLen < 0 {
				// Invalid syntax; don't just skip validation but break the compilation so this can be fixed in the definition
				panic("unparsable string validator value for len in field " + f.GoIdent.GoName + ": expected 0<x<y, found 0>=x>=y " + fLen)
			}
		case 1:
			// Fixed length defined, e.g. len: "8" - min and max length are equal in this case
			fMinLen, _ = strconv.Atoi(fLenChunks[0])
			fMaxLen = fMinLen
		default:
			// Invalid definition, more than two -, e.g. -5 or 5-10-12 or -5-5 or something we don't understand
			// This should be fixed so break the compilation
			panic("unparsable string validator value for len in field " + f.GoIdent.GoName + ": expected x-y or x, found " + fLen)
		}

		if fMinLen > 0 {
			minLen = uint32(fMinLen)
		}
		if fMaxLen > 0 {
			maxLen = uint32(fMaxLen)
		}
	}

	// Now both minLen and maxLen are set
	// First make sure this satisfies our predefined absolute min/max values for `string` and `unsafe_string`
	validMin := stringLenMinSafe
	validMax := stringLenMaxSafe
	if stringType == validator.E_UnsafeString {
		validMin = stringLenMinUnsafe
		validMax = stringLenMaxUnsafe
	}
	if minLen < validMin || maxLen > validMax {
		// Could just warn and ignore this instead of breaking initially until we are confident that we have appropriate safe guard values
		panic("invalid string validator value for len in field " + f.GoIdent.GoName + ": expected " + fmt.Sprint(validMin) + "<=x<=" + fmt.Sprint(validMax) + ", found " + fmt.Sprint(minLen) + "-" + fmt.Sprint(maxLen))
	}

	if minLen >= 1 {
		*stringValidator += ".min(" + strconv.Itoa(int(minLen)) + ")"
	}

	if maxLen >= 1 {
		*stringValidator += ".max(" + strconv.Itoa(int(maxLen)) + ")"
	}

	if rules.GetOptional() || f.Desc.HasOptionalKeyword() {
		*stringValidator += ".optional()"
	}

	// // ##### 4B. prepare whitelist/regex
	// // Check if allow option is set: these will be allowed in addition to the default whitelist
	// if rules.GetAllow() != "" {
	// 	// Note that this does inentionally not support a regex pattern or entire categories like \pX
	// 	// All characters in the string are considered as LITERAL value

	// 	// Convert chars to \x{dddd} notation
	// 	for _, runeValue := range rules.GetAllow() {
	// 		unicodeValue := fmt.Sprintf("%U", runeValue)
	// 		unicodeValue = strings.Replace(unicodeValue, "U+", "", 1)
	// 		if stringType == validator.E_String {
	// 			if _, isUnsafe := stringUnsafeReplacerMap[`\u`+unicodeValue]; isUnsafe {
	// 				// Chars in stringUnsafeReplacerMap can not be used in `string` allow option when `replace_unsafe` is not enabled
	// 				if !rules.GetReplaceUnsafe() {
	// 					panic("invalid allow character in field " + f.GoIdent.GoName + ": " + string(runeValue) + ", enable replace_unsafe option or use validator.unsafe_string instead")
	// 				}
	// 				// else: this char was replaced with a safe equivalent so we do not need to add the unsafe one to the allow regex
	// 				continue
	// 			}
	// 			// Chars in charDenyList can not be used in allow option as they are potentially dangerous
	// 			if _, isRestricted := charDenyList[`\u`+unicodeValue]; isRestricted {
	// 				panic("invalid allow character in field " + f.GoIdent.GoName + ": U+" + string(unicodeValue) + ", this character is potentially dangerous (check charDenyList)")
	// 			}
	// 		}
	// 		prepareRegex(allowListReId, `\x{`+unicodeValue+`}`)
	// 	}
	// }

	// // Check for multiline option which allows linebreaks
	// if rules.GetMultiline() {
	// 	// Add predefined linebreak chars (currently just \n)
	// 	// \r was stripped in step 2B
	// 	prepareRegex(allowListReId, stringReLineBreaks...)
	// }

	// // Add symbols if defined
	// symbols := rules.GetSymbols()
	// if len(symbols) > 0 {
	// 	// Add the corresponding regex pattern to the allow list
	// 	for i := range symbols {
	// 		symbol := symbols[i]
	// 		if stringType == validator.E_String && !rules.GetReplaceUnsafe() {
	// 			// Check and error on restricted symbol classes
	// 			if _, ok := restrictedSafeStringSymbols[symbol]; ok {
	// 				panic("invalid symbol class in field " + f.GoIdent.GoName + ": " + symbol.String() + ", enable replace_unsafe option or use unsafe_string instead.")
	// 			}
	// 		}
	// 		if symbolRe, ok := stringSymbolMap[symbol]; ok {
	// 			prepareRegex(allowListReId, symbolRe...)
	// 			continue
	// 		}
	// 		// Ignore this as it would not weaken the validation.
	// 		// Warn as it likely indicates something that requires follow up
	// 		fmt.Fprintf(os.Stderr, "WARN: Symbol "+fmt.Sprint(symbol)+" not in stringSymbolMap (not implemented)\n")
	// 	}
	// }

	// // ##### 4C. validate string against the whitelist/regex
	// g.P("if !_regex_", "regexName", ".MatchString(", varName, ") {")
	// // Might want to improve the error message: make it more helpful so the end-user can fix it if applicable
	// errStr = "only have valid characters"
	// if rules.GetLogOnly() {
	// 	printErrorString(g, varName, string(f.Desc.Name()), errStr, 50)
	// } else {
	// 	genErrorStringWithParams(g, varName, string(f.Desc.Name()), errStr)
	// }
	// g.P("}")

	// // Optional value: close the if statement
	// if rules.GetOptional() {
	// 	g.P("}")
	// }
}

func genIdValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {

	// ID validator with support for legacy ids and UUIDv4
	// Might want to move the shared validation logic such as optional to a separate method
	rules := getIdExtension(f, validator.E_Id)
	if rules == nil {
		return
	}

	// var errMsg string

	switch rules.GetVersion() {
	case "v4":
		// Try to validate the value for valid UUIDv4
		*stringValidator += ".uuid()"
		// errMsg = "be parsable as UUIDv4"
	case "any":
		// Try to validate the value for valid UUID (any version)
		// g.P("if !", s12protoPackage.Ident("IsUUID"), "(", varName, ") {")
		// errMsg = "be parsable as UUID"
	default:
		// Unsupported version; do not generate the validators without having an implementation for this version
		panic("unsupported UUID version in field " + f.GoIdent.GoName + ", got: " + rules.GetVersion())
	}

	if rules.GetOptional() {
		*stringValidator += ".optional()"
	}

	// Check what ID formats are accepted in addition to UUID
	// Validation passes if the provided value passes at least one of the validators
	// UUID format is always accepted where LegacyId and S12 Id can be enabled as an option

	// Invalid UUID: check other formats or error if no other options enabled
	// g.P("isValidId := false")

	// // Check for valid legacy id if option enabled
	// if rules.GetLegacy() {
	// 	errMsg += " or legacy ID"
	// 	g.P("if ", s12protoPackage.Ident("IsLegacyID"), "(", varName, ", ", rules.GetLowercaseOnly(), ") {")
	// 	g.P("isValidId = true")
	// 	g.P("}")
	// }

	// // If still no valid id found, also check for S12 Id if option enabled
	// if rules.GetS12Id() {
	// 	errMsg += " or S12 ID"
	// 	g.P("if !isValidId && ", s12protoPackage.Ident("IsS12ID"), "(", varName, ", ", rules.GetLowercaseOnly(), ") {")
	// 	g.P("isValidId = true")
	// 	g.P("}")
	// }

	// // Finally, we have a special case for prefixed, long legacy format Ids which requires both s12 and legacy option
	// if rules.GetLegacy() && rules.GetS12Id() {
	// 	g.P("if !isValidId && ", s12protoPackage.Ident("IsLongPrefixedLegacyID"), "(", varName, ") {")
	// 	g.P("isValidId = true")
	// 	g.P("}")
	// }

	// // Check if any of the validations passed
	// g.P("if !isValidId {")
	// if rules.GetLogOnly() {
	// 	printErrorString(g, varName, string(f.Desc.Name()), errMsg, 50)
	// } else {
	// 	genErrorString(g, varName, string(f.Desc.Name()), errMsg)
	// }

	// g.P("}")
	// g.P("}")

	// if rules.GetOptional() {
	// 	g.P("}")
	// }
}

func genEmailValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {

	rules := getEmailExtension(f, validator.E_Email)
	if rules == nil {
		return
	}

	*stringValidator += ".email().max(231)"

	if rules.GetOptional() {
		*stringValidator += ".optional()"
	}

}

func genURLValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {

	rules := getURLExtension(f, validator.E_Url)
	if rules == nil {
		return
	}

	*stringValidator += ".url()"

	// Check defined schemes
	schemes := rules.GetSchemes()
	if len(schemes) == 0 {
		schemes = append(schemes, "https")
	}

	if rules.GetAllowHttp() {
		schemes = append(schemes, "http")
	}

	*stringValidator += ".regex(/^(" + strings.Join(schemes, "|") + ")\\:\\/\\//)"

	if rules.GetOptional() {
		*stringValidator += ".optional()"
	}

}

func genTimezoneValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {
	rules := getTimezoneExtension(f, validator.E_Timezone)
	if rules == nil {
		return
	}

	*stringValidator = ".enum([...Intl.supportedValuesOf('timeZone')] as [string, ...string[]])"

	if rules.GetOptional() {
		*stringValidator += ".optional()"
	}

}

func genLenValidator(g *protogen.GeneratedFile, f *protogen.Field, stringValidator *string) {
	if getBoolExtension(f, validator.E_TrimLenCheck) {
		*stringValidator += ".trim()"
	}

	if v := getIntExtention(f, validator.E_LengthGte); v >= 0 {
		*stringValidator += ".min(" + strconv.FormatInt(v, 10) + ")"

	}

	if v := getIntExtention(f, validator.E_LengthLte); v >= 0 {
		*stringValidator += ".max(" + strconv.FormatInt(v, 10) + ")"
	}
}

func genBytesValidator(g *protogen.GeneratedFile, f *protogen.Field) string {
	optional := getBoolExtension(f, validator.E_Optional)

	bytesValidator := ".string()"

	if getBoolExtension(f, validator.E_Uuid) {
		bytesValidator += ".refine((val) => new Blob([val]).size === 16)"
	}

	genLenValidator(g, f, &bytesValidator)

	if optional || f.Desc.HasOptionalKeyword() {
		bytesValidator += ".optional()"
	}

	return bytesValidator
}

func genMsgValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) string {

	// g.P(f.Desc.Name())
	objectValidator := f.Message.GoIdent.GoName + "Validator"
	if f.Message.GoIdent.GoName == "Timestamp" {
		objectValidator = "z.string().datetime()"
	}
	if !getBoolExtension(f, validator.E_MsgRequired) {
		objectValidator += ".optional()"
	}

	return objectValidator

	// // For repeated messages, we need to run the validator on each message instead of this field
	// repeated := f.Desc.Cardinality() == protoreflect.Repeated
	// if repeated {
	// 	// Please leave the len check in place here as I have identified edge cases where it is required (will add test case for it later)
	// 	g.P("if len(", varName, ") > 0 {")
	// 	g.P("for _, item := range ", varName, "{")
	// 	varName = "item"
	// }

	// g.P("if ", varName, " != nil {")
	// g.P("if v, ok := interface{}(", varName, ").(", s12protoPackage.Ident("Validator"), "); ok {")
	// g.P("if err := v.Validate(); err != nil {")
	// g.P("return ", s12protoPackage.Ident("FieldError"), "(\"", f.Desc.Name(), "\", err)")
	// g.P("}")
	// g.P("}")
	// g.P("}")

	// if repeated {
	// 	g.P("}") // for
	// 	g.P("}") // if len
	// }
}

func genEnumValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) string {
	enumName := f.Desc.JSONName() + "Enum"
	if !tsVars[enumName] {
		g.P("const " + enumName + " = {")
		for _, v := range f.Enum.Values {
			g.P("  ", v.Desc.Name(), ": ", v.Desc.Number(), ",")
		}
		g.P("} as const")
		tsVars[enumName] = true
	}

	enumValidator := ".nativeEnum(" + enumName + ")"

	if (!getBoolExtension(f, validator.E_EnumRequired)) || f.Desc.HasOptionalKeyword() {
		enumValidator += ".optional()"
	}

	return enumValidator
}

func printErrorString(g *protogen.GeneratedFile, varName, fieldName, specificErr string, maxLen int) {
	// Do not reflect untrusted value in error, certainly not for sensitive fields like password or PII like email
	g.P(fmtPackage.Ident("Printf"), `("[log-only] %s: value must be %s: Base64Encoded input: %s\n", `, `"`, fieldName, `", "`, specificErr, `" `, ", ",
		s12protoPackage.Ident("Base64Encode("),
		s12protoPackage.Ident("FirstCharactersFromString("),
		fmtPackage.Ident("Sprintf"), "(\"%v\", ", varName, "), ", maxLen, ")))")
}

func genErrorString(g *protogen.GeneratedFile, varName, fieldName, specificErr string) {
	// Do not reflect untrusted value in error, certainly not for sensitive fields like password or PII like email
	g.P(`return `, fmtPackage.Ident("Errorf"), "(`", fieldName, `: value must `, specificErr, "`)")
}

func genErrorStringWithParams(g *protogen.GeneratedFile, varName, fieldName, specificErr string, errParams ...string) {
	// Might want to improve the error message: make it more helpful so the end-user can fix it if applicable
	g.P(`return `, fmtPackage.Ident("Errorf"), "(`", fieldName, `: value must `, specificErr, "`, ", strings.Join(errParams, ", "), `)`)
}

func genLenErrorString(g *protogen.GeneratedFile, varName, fieldName, specificErr string) {
	g.P(`return `, fmtPackage.Ident("Errorf"), "(`", fieldName, `: length must `, specificErr, "`)")
}

var validNonRepeatedExts = []protoreflect.ExtensionType{
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
	validator.E_Email,
	validator.E_Id,
	validator.E_String,
	validator.E_UnsafeString,
	validator.E_EnumRequired,
	validator.E_Url,
	validator.E_Timezone,
	validator.E_Number,
	validator.E_SimpleString,
}

var validRepeatedExts = []protoreflect.ExtensionType{
	validator.E_RepeatedLenGte,
	validator.E_RepeatedLenLte,
}

var validExts = append(validNonRepeatedExts, validRepeatedExts...)

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

func hasNonRepeatedValidationExtensions(f *protogen.Field) bool {
	if opts := f.Desc.Options(); opts != nil {
		for _, ext := range validNonRepeatedExts {
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

func getSimpleStringExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.SimpleStringRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.SimpleStringRules); ok {
			return v
		}
	}
	return nil
}

func getStringExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.StringRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.StringRules); ok {
			return v
		}
	}
	return nil
}

func getIdExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.IdRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.IdRules); ok {
			return v
		}
	}
	return nil
}

func getEmailExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.EmailRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.EmailRules); ok {
			return v
		}
	}
	return nil
}

func getURLExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.URLRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.URLRules); ok {
			return v
		}
	}
	return nil
}

func getTimezoneExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.TimezoneRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.TimezoneRules); ok {
			return v
		}
	}
	return nil
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

func genNumberValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string, isInt bool) string {
	rules := getNumberExtension(f, validator.E_Number)
	numberValidation := ".number()"
	if isInt {
		numberValidation += ".int()"
	}

	if !rules.GetAllowNan() {
		if f.Desc.Kind() != protoreflect.DoubleKind && f.Desc.Kind() != protoreflect.FloatKind {
			panic("cannot use allow_nan option for integers, only supported for float/double")
		}
		// g.P("// This statement checks for NaN value without using Math package")
		// g.P("if ", varName, " != ", varName, " {")
		// errStr := "not be NaN"
		// genErrorString(g, varName, string(f.Desc.Name()), errStr)
		// g.P("}")
	}

	if rules.GetRange() != "" {
		if !strings.Contains(rules.GetRange(), ":") {
			// Range must contain : to be used
			panic("unparsable range for number validator")
		}

		var rangeVals = strings.Split(rules.GetRange(), ":")
		if len(rangeVals) < 1 || len(rangeVals) > 2 {
			panic("unparsable range for number validator")
		}

		if rangeVals[0] != "" {
			numberValidation += ".min(" + rangeVals[0] + ")"
		}

		if rangeVals[1] != "" {
			numberValidation += ".max(" + rangeVals[1] + ")"
		}
	}
	if rules.GetOptional() || f.Desc.HasOptionalKeyword() {
		numberValidation += ".optional()"
	}
	return numberValidation
}

func getNumberExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.NumberRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.NumberRules); ok {
			return v
		}
	}
	return nil
}
