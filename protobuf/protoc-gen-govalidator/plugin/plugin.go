// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

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
var regexGeneratedFile *protogen.GeneratedFile
var regexHashLib = make(map[string]struct{})

// Validator plugin version
var validatorVersion = "v2.4.0"

// Write a preamble in the auto generated files
func genGeneratedHeader(gen *protogen.Plugin, g *protogen.GeneratedFile, f *protogen.File) {
	g.P("// Code generated by protoc-gen-govalidator. DO NOT EDIT.")
	g.P("// versions:")
	protocVersion := "(unknown)"
	if v := gen.Request.GetCompilerVersion(); v != nil {
		protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
		if s := v.GetSuffix(); s != "" {
			protocVersion += "-" + s
		}
	}
	g.P("// \tprotoc-gen-govalidator ", validatorVersion)
	g.P("// \tprotoc                 ", protocVersion)

	if f.Proto.GetOptions().GetDeprecated() {
		g.P("// ", f.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", f.Desc.Path())
	}
	g.P()
}

// GenerateFile generates the validator.pb.go file
func GenerateFile(p *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".validator.pb.go"
	g := p.NewGeneratedFile(filename, file.GoImportPath)

	// Add preamble and version for validator.pb.go
	genGeneratedHeader(p, g, file)

	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	// Generate a blank import statement if not used by any validation functions,
	// as we still want to register the protobuf validation types for the Reflection API.
	g.Import(s12protoPackage)

	// Also prepare a separate file for regex patterns
	// Writing these to a separate file otherwise the generation logic would be confusing
	// and require two iterations for each message field (one to generate the regex, one for the rest)
	regexFilename := file.GeneratedFilenamePrefix + ".validator_regex.pb.go"
	regexGeneratedFile = p.NewGeneratedFile(regexFilename, file.GoImportPath)

	// Add preamble and version for validator_regex.pb.go
	genGeneratedHeader(p, regexGeneratedFile, file)
	regexGeneratedFile.P()
	regexGeneratedFile.P("package ", file.GoPackageName)
	regexGeneratedFile.P()
	regexGeneratedFile.Import(s12protoPackage)

	genStringGenerics()
	for _, msg := range file.Messages {
		if opts, ok := msg.Desc.Options().(*descriptorpb.MessageOptions); !ok || opts.GetMapEntry() {
			continue
		}
		genLegacyRegexVars(g, msg)
		genValidateFunc(g, msg)
		g.P()
	}

	return g
}

func addRegexVar(fieldName, regexId string) string {

	r, e := getPreparedRegex(regexId)
	if e != nil {
		// Should not continue as we don't have a regex pattern
		panic("Error generating regex for " + regexId + "/" + fieldName + ": " + fmt.Sprint(e))
	}

	// Generate a unique hash based on regex pattern so we can reuse patterns
	hash := md5.New()
	hash.Write([]byte(r))
	name := hex.EncodeToString(hash.Sum(nil))

	// Check if this is a duplicate regex pattern in which case we can reuse the regex
	if _, found := regexHashLib[name]; found {
		regexGeneratedFile.P("// " + fieldName + " is using regex " + name)
		return name
	}

	regexGeneratedFile.P("// Pattern for " + fieldName)
	regexGeneratedFile.P("const ", "_regex_val_", name, " = `", r, "`")
	regexGeneratedFile.P("var ", "_regex_", name, " = ", regexpPackage.Ident("MustCompile"), "(_regex_val_", name, ")")

	// Keep track of this regex to avoid duplicates
	if _, found := regexHashLib[name]; !found {
		regexHashLib[name] = struct{}{}
	}
	return name
}

// Generic preparations/initialisations for the string validator
func genStringGenerics() {

	// Prepare the regex pattern to allow replaced chars
	prepareStringReplacerRegex(stringUnsafeReplacerMap, "replacer_unsafe_allowed")
	prepareStringReplacerRegex(stringSymbolReplacerMap, "replacer_symbol_allowed")

	// Prepare the default `unsafe_string` pattern by adding tokens from `string`
	stringReDefaultUnsafe = append(stringReDefaultSafe, stringReDefaultUnsafe...)
}

// Generator for legacy regex validator
func genLegacyRegexVars(g *protogen.GeneratedFile, msg *protogen.Message) {
	for _, field := range msg.Fields {
		if regex := getRegexValue(field); regex != "" {
			g.P("const ", "_regex_val_", field.GoIdent, " = `", regex, "`")
			g.P("var ", "_regex_", field.GoIdent, " = ", regexpPackage.Ident("MustCompile"), "(_regex_val_", field.GoIdent, ")")
		}
	}
	for _, innerMsg := range msg.Messages {
		genLegacyRegexVars(g, innerMsg)
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
		hasRepeatedExt := hasExt && f.Desc.Cardinality() == protoreflect.Repeated
		isMessageField := f.Desc.Kind() == protoreflect.MessageKind

		// shouldLoopOverField decides whether to generate the for loop or not:
		// - don't generate the for loop for an array of msg because it wil be handled differently
		// - only generate the for loop if there are validations for each elment
		shouldLoopOverField := hasRepeatedExt && !isMessageField && hasNonRepeatedValidationExtensions(f)

		if hasRepeatedExt {
			if v := getIntExtention(f, validator.E_RepeatedLenGte); v >= 0 {
				g.P("if !(len(", varName, ") >= ", v, ") {")
				errStr := fmt.Sprintf(`be greater than or equal to %d`, v)
				genLenErrorString(g, varName, string(f.Desc.Name()), errStr)
				g.P("}")
			}

			if v := getIntExtention(f, validator.E_RepeatedLenLte); v >= 0 {
				g.P("if !(len(", varName, ") <= ", v, ") {")
				errStr := fmt.Sprintf(`be lesser than or equal to %d`, v)
				genLenErrorString(g, varName, string(f.Desc.Name()), errStr)
				g.P("}")
			}
		}

		if shouldLoopOverField {
			g.P("for _, item := range ", varName, "{")
			varName = "item"
		}

		if f.Desc.IsMap() {
			g.P(`// Validation of proto3 map<> fields is unsupported.`)
			continue
		}

		if f.Oneof != nil {
			g.P("if x, ok := m.", f.Oneof.GoName, ".(*", g.QualifiedGoIdent(f.GoIdent), "); ok {")
			varName = "x." + f.GoName
		}

		switch f.Desc.Kind() {
		case protoreflect.StringKind:
			// Run both legacy and newer validation logic as both options can be combined at this time
			// Should exclusively use newer validation logic (validator.string/unsafe_string) for newly added fields
			// For existing fields, should be fairly easy to update the rules to the new validator
			genLegacyStringValidator(g, f, varName)
			genStringValidator(g, f, varName)
			genIdValidator(g, f, varName)
			genEmailValidator(g, f, varName)
			genURLValidator(g, f, varName)
			genTimezoneValidator(g, f, varName)
		case protoreflect.BytesKind:
			// IdValidator not supported for bytes at this point
			genBytesValidator(g, f, varName)
		case protoreflect.Int32Kind, protoreflect.Int64Kind,
			protoreflect.Uint32Kind, protoreflect.Uint64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind:
			genIntValidator(g, f, varName)
		case protoreflect.MessageKind:
			genMsgValidator(g, f, varName)
		case protoreflect.EnumKind:
			genEnumValidator(g, f, varName)
		case protoreflect.DoubleKind, protoreflect.FloatKind:
			genFloatValidator(g, f, varName)

		}

		if f.Oneof != nil {
			g.P("}")
		}

		if shouldLoopOverField {
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

func genLegacyStringValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	optional := getBoolExtension(f, validator.E_Optional)

	if optional {
		g.P("if ", varName, " != \"\" {")
	}

	if regex := getRegexValue(f); regex != "" {
		g.P("if !_regex_", f.GoIdent, ".MatchString(", varName, ") {")
		// Let's not expose the regex pattern and find a better way to reflect what was wrong with the data
		// Might want to replace with a more helpful message
		errStr := "be a string conforming to predefined pattern"
		genErrorStringWithParams(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	// Legacy UUIDv4 validator
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

func genStringValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
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

	// Prepare the allow list regex
	allowListReId := "re_" + f.GoIdent.GoName
	if stringType == validator.E_String {
		prepareRegex(allowListReId, stringReDefaultSafe...)
	} else {
		prepareRegex(allowListReId, stringReDefaultUnsafe...)
	}

	// "optional" can be provided via the legacy directive or set as option on the `string`/`unsafe_string` validator
	// To allow deprecation of the 'old' directive later, we duplicate the logic here
	if rules.GetOptional() {
		g.P("if ", varName, " != \"\" {")
	}

	// ### STEP 1: decode
	// TODO PA: Check if we can confirm that data has been decoded, things like %20 might indicate otherwise

	// ### STEP 2: normalise and canonicalise
	// Normalisation for file paths/URLs should not happen here but in URL validator

	// ##### 2A. Normalise Unicode NFD strings to NFC
	// Convert to unicode NFC if not already in that format
	//  e.g. NFD string man\u0303ana should be normalised to ma\u00F1ana (n + \u0303 = \u00F1)
	// TODO PA: Check if we need support for NFKC/NFKD
	// TODO PA: Do we want to have support for U+202E RIGHT-TO-LEFT OVERRIDE by default for all strings?
	g.P("if !", normPackage.Ident("NFC.IsNormalString"), "(", varName, ") && ", normPackage.Ident("NFD.IsNormalString"), "(", varName, ") {")
	g.P("// normalise NFD to NFC string")
	g.P("var normErr error")
	g.P(varName, ", _, normErr = ", transformPackage.Ident("String"), "(", transformPackage.Ident("Chain"), "(", normPackage.Ident("NFD"), ",", normPackage.Ident("NFC"), ")", ",", varName, ")")
	g.P("if normErr != nil {")
	genErrorString(g, varName, string(f.Desc.Name()), `must be normalisable to NFC`)
	g.P("}")
	g.P("}")

	// ##### 2B. Check for encoding issues
	if rules.GetValidateEncoding() {
		// U+FFFD in character indicates that the encoding was likely incorrect, resulting in � U+FFFD (RuneError)
		// Test string: $\xa35 for Pepp\xe9 which is encoded as latin1
		g.P("if ", stringsPackage.Ident("ContainsRune"), "(", varName, ", ", utfPackage.Ident("RuneError"), ") {")
		errStr := `must have valid encoding`
		genErrorString(g, varName, string(f.Desc.Name()), errStr)

		// Also validate using utf8.ValidString: reports whether string consists entirely of valid UTF-8-encoded runes
		g.P("} else if !", utfPackage.Ident("ValidString"), "(", varName, ") {")
		errStr = `must be a valid UTF-8-encoded string`
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	} else {
		// Allow RuneError as valid character if we do not check for invalid encoding otherwise the string might be denied
		// U+FFFD is in Symbol_So category which must be allowed in case validateInvalidEncoding is disabled (to detect issues)
		prepareRegex(allowListReId, `\x{FFFD}`)
	}

	// ### STEP 3: sanitise
	// ##### 3A. sanitise whitespace (trim option)
	// Note this will PERMANENTLY mutate the message field data, this trim is not for length check only but for sanitisation
	// WARNING: Any leading/trailing whitespace will be permanently removed (this is intended here)
	if rules.GetTrim() {
		g.P(varName, " = ", stringsPackage.Ident("TrimSpace"), "(", varName, ")")
	}

	// ##### 3B. replace restricted characters with a safe alternative
	// WARNING: this could corrupt data, eg. hyperlinks in a description box so disabled by default
	// Ensure to test the alternative characters in all output (database, email, pdf, csv, web UI, mobile UI)
	if rules.GetReplaceUnsafe() {
		// If allow option is defined, only chars from allow list will be allowed
		// Otherwise, all chars in the replace map are accepted
		// If none of the allow chars are in the unsafe map, there is no work to be done here
		// throw a warning that the option is not required in that case
		if rules.GetAllow() != "" {

			for _, runeValue := range rules.GetAllow() {
				unicodeValue := fmt.Sprintf("%U", runeValue)
				unicodeValue = strings.Replace(unicodeValue, "U+", "", 1)
				if _, isUnsafe := stringUnsafeReplacerMap[`\u`+unicodeValue]; isUnsafe {
					// need to replace this one, replace case by case instead of using the replacer
					replacedUnicodeValue := stringUnsafeReplacerMap[`\u`+unicodeValue]
					g.P(varName, " = ", stringsPackage.Ident("ReplaceAll"), "(", varName, `, "\u`, unicodeValue, `", "`, replacedUnicodeValue, `")`)

					// and add it to the regex
					prepareRegex(allowListReId, `\x{`+strings.Replace(replacedUnicodeValue, `\u`, ``, 1)+`}`)
				}
			}
		} else {
			// Add the possible replaced characters to the allow list as they are not all allowed by default
			mergeRegex(allowListReId, "replacer_unsafe_allowed")
			// Replace the unsafe values
			g.P(varName, " = ", s12protoPackage.Ident("UnsafeCharReplacer"), ".Replace(", varName, ")")
		}
	}

	// ##### 3C. strip carriage return characters \r in multiline strings (leave \n)
	if rules.GetMultiline() {
		// Generates: varName = strings.ReplaceAll(varName, "\r", "")
		g.P(varName, " = ", stringsPackage.Ident("ReplaceAll"), "(", varName, ", \"\\r\", \"\")")
	}

	// ##### 3D. replace rare symbols with a more common alternative
	if rules.GetReplaceOther() {
		// Add the possible replaced characters to the allow list as they might not be allowed by default
		mergeRegex(allowListReId, "replacer_symbol_allowed")
		// Replace the values
		if rules.GetMultiline() {
			// Should not replace \n as the normal symbol replacer does otherwise multiline breaks
			g.P(varName, " = ", s12protoPackage.Ident("SymbolCharReplacerMultiline"), ".Replace(", varName, ")")
		} else {
			g.P(varName, " = ", s12protoPackage.Ident("SymbolCharReplacer"), ".Replace(", varName, ")")

		}
	}

	// ##### 3E. Sanitise (remove) Private Use Area Codepoints in the Basic Multilingual Plane
	// Note that we do not check for PUA in planes 15 and 16 currently
	if rules.GetSanitisePua() {
		g.P(varName, "= ", s12protoPackage.Ident("RegexPua"), ".ReplaceAllString(", varName, ", \"\")")
	}

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

	// Determine which length check method to use: bytes or runes
	// Byte length check is default [using len()], unless rune length option is enabled
	lenVar := "_len_" + f.GoIdent.GoName
	if rules.GetRunes() {
		// Use utf8.RuneCountInString method, requires import of utf8 package
		g.P("var "+lenVar+" = ", utfPackage.Ident("RuneCountInString"), "(", varName, ")")
	} else {
		g.P("var "+lenVar+" = len(", varName, ")")
	}

	// Write the len check logic to the validator
	errStr := ""
	if minLen == maxLen {
		// Could use the same if statement as for min+max checks but this is a bit cleaner
		g.P("if !("+lenVar+" == ", minLen, ") {")
		errStr = fmt.Sprintf(`have length %d`, minLen)
	} else {
		g.P("if !("+lenVar+" >= ", minLen, " && "+lenVar+" <= ", maxLen, ") {")
		errStr = fmt.Sprintf(`have length between %d and %d`, minLen, maxLen)
	}
	genErrorString(g, varName, string(f.Desc.Name()), errStr)
	g.P("}")

	// ##### 4B. prepare whitelist/regex
	// Check if allow option is set: these will be allowed in addition to the default whitelist
	if rules.GetAllow() != "" {
		// Note that this does inentionally not support a regex pattern or entire categories like \pX
		// All characters in the string are considered as LITERAL value

		// Convert chars to \x{dddd} notation
		for _, runeValue := range rules.GetAllow() {
			unicodeValue := fmt.Sprintf("%U", runeValue)
			unicodeValue = strings.Replace(unicodeValue, "U+", "", 1)
			if stringType == validator.E_String {
				if _, isUnsafe := stringUnsafeReplacerMap[`\u`+unicodeValue]; isUnsafe {
					// Chars in stringUnsafeReplacerMap can not be used in `string` allow option when `replace_unsafe` is not enabled
					if !rules.GetReplaceUnsafe() {
						panic("invalid allow character in field " + f.GoIdent.GoName + ": " + string(runeValue) + ", enable replace_unsafe option or use validator.unsafe_string instead")
					}
					// else: this char was replaced with a safe equivalent so we do not need to add the unsafe one to the allow regex
					continue
				}
				// Chars in charDenyList can not be used in allow option as they are potentially dangerous
				if _, isRestricted := charDenyList[`\u`+unicodeValue]; isRestricted {
					panic("invalid allow character in field " + f.GoIdent.GoName + ": U+" + string(unicodeValue) + ", this character is potentially dangerous (check charDenyList)")
				}
			}
			prepareRegex(allowListReId, `\x{`+unicodeValue+`}`)
		}
	}

	// Check for multiline option which allows linebreaks
	if rules.GetMultiline() {
		// Add predefined linebreak chars (currently just \n)
		// \r was stripped in step 2B
		prepareRegex(allowListReId, stringReLineBreaks...)
	}

	// Add symbols if defined
	symbols := rules.GetSymbols()
	if len(symbols) > 0 {
		// Add the corresponding regex pattern to the allow list
		for i := range symbols {
			symbol := symbols[i]
			if stringType == validator.E_String && !rules.GetReplaceUnsafe() {
				// Check and error on restricted symbol classes
				if _, ok := restrictedSafeStringSymbols[symbol]; ok {
					panic("invalid symbol class in field " + f.GoIdent.GoName + ": " + symbol.String() + ", enable replace_unsafe option or use unsafe_string instead.")
				}
			}
			if symbolRe, ok := stringSymbolMap[symbol]; ok {
				prepareRegex(allowListReId, symbolRe...)
				continue
			}
			// Ignore this as it would not weaken the validation.
			// Warn as it likely indicates something that requires follow up
			fmt.Fprintf(os.Stderr, "WARN: Symbol "+fmt.Sprint(symbol)+" not in stringSymbolMap (not implemented)\n")
		}
	}

	// ##### 4C. validate string against the whitelist/regex
	regexName := addRegexVar(f.GoIdent.GoName, allowListReId)
	g.P("if !_regex_", regexName, ".MatchString(", varName, ") {")
	// Might want to improve the error message: make it more helpful so the end-user can fix it if applicable
	errStr = "only have valid characters"
	genErrorStringWithParams(g, varName, string(f.Desc.Name()), errStr)
	g.P("}")

	// Optional value: close the if statement
	if rules.GetOptional() {
		g.P("}")
	}
}

func genIdValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {

	// ID validator with support for legacy ids and UUIDv4
	// Might want to move the shared validation logic such as optional to a separate method
	rules := getIdExtension(f, validator.E_Id)
	if rules == nil {
		return
	}

	if rules.GetOptional() {
		g.P("if ", varName, " != \"\" {")
	}

	if rules.GetVersion() != "v4" {
		// Unsupported version; do not generate the validators without having an implementation for this version
		panic("unsupported UUID version in field " + f.GoIdent.GoName + ": expected v4, got: " + rules.GetVersion())
	}

	// Check what ID formats are accepted in addition to UUID
	// Validation passes if the provided value passes at least one of the validators
	// UUID format is always accepted where LegacyId and S12 Id can be enabled as an option
	errMsg := "be parsable as UUIDv4"

	// Always try to validate the value for valid UUIDv4
	g.P("if !", s12protoPackage.Ident("IsUUIDv4"), "(", varName, ") {")
	// Invalid UUID: check other formats or error if no other options enabled
	g.P("isValidId := false")

	// Check for valid legacy id if option enabled
	if rules.GetLegacy() {
		errMsg += " or legacy ID"
		g.P("if ", s12protoPackage.Ident("IsLegacyID"), "(", varName, ") {")
		g.P("isValidId = true")
		g.P("}")
	}

	// Finally, if still no valid id found, also check for S12 Id if option enabled
	if rules.GetS12Id() {
		errMsg += " or S12 ID"
		g.P("if !isValidId && ", s12protoPackage.Ident("IsS12ID"), "(", varName, ") {")
		g.P("isValidId = true")
		g.P("}")
	}

	// Check if any of the validations passed
	g.P("if !isValidId {")
	if rules.GetLogOnly() {
		printErrorString(g, varName, string(f.Desc.Name()), errMsg, 50)
	} else {
		genErrorString(g, varName, string(f.Desc.Name()), errMsg)
	}

	g.P("}")
	g.P("}")

	if rules.GetOptional() {
		g.P("}")
	}
}

func genEmailValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {

	rules := getEmailExtension(f, validator.E_Email)
	if rules == nil {
		return
	}

	if rules.GetOptional() {
		g.P("if ", varName, " != \"\" {")
	}

	checkDomain := false // rules.GetCheckDomain() == true

	// Validate the value for valid email using govalidator validation package
	g.P("if !", s12protoPackage.Ident("IsValidEmail"), "(", varName, ", ", checkDomain, ") {")
	errStr := "be parsable as an email address"
	if checkDomain {
		errStr = "be a valid email address"
	}
	// Should not return/log email address as this is PII (genErrorString will no longer reflect values)
	genErrorString(g, varName, string(f.Desc.Name()), errStr)
	g.P("}")

	if rules.GetOptional() {
		g.P("}")
	}
}

func genURLValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {

	rules := getURLExtension(f, validator.E_Url)
	if rules == nil {
		return
	}

	if rules.GetOptional() {
		g.P("if ", varName, " != \"\" {")
	}

	// Check defined schemes
	schemes := rules.GetSchemes()
	if len(schemes) == 0 {
		schemes = append(schemes, "https")
	}

	if rules.GetAllowHttp() {
		schemes = append(schemes, "http")
	}

	// Validate the url using the helper method
	g.P(`_schemes_`+f.GoIdent.GoName+` := []string{"`, strings.Join(schemes, `", "`), `"}`)
	g.P("if _, err := ", s12protoPackage.Ident("IsValidURL"), "(", varName, ", _schemes_"+f.GoIdent.GoName+", ", rules.GetAllowFragment(), "); err != nil {")
	genErrorStringWithParams(g, varName, string(f.Desc.Name()), "be parsable as a URL: %v", "err")
	g.P("}")

	if rules.GetOptional() {
		g.P("}")
	}
}

func genTimezoneValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	rules := getTimezoneExtension(f, validator.E_Timezone)
	if rules == nil {
		return
	}

	if rules.GetOptional() {
		g.P("if ", varName, " != \"\" {")
	} else {
		g.P("if ", varName, " == \"\" {")
		g.P("return ", fmtPackage.Ident("Errorf"), "(\"field ", f.Desc.Name(), " is required\")")
		g.P("}")
	}

	g.P("if tz, err := ", timePackage.Ident("LoadLocation"), "(", varName, "); err != nil || tz == nil {")
	errStr := "be a valid IANA TZ database value"
	genErrorStringWithParams(g, varName, string(f.Desc.Name()), errStr)
	g.P("}")

	if rules.GetOptional() {
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
		errStr := fmt.Sprintf(`have length greater than or equal to %d`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	if v := getIntExtention(f, validator.E_LengthLte); v >= 0 {
		g.P("if !(len(", varName, ") <= ", v, ") {")
		errStr := fmt.Sprintf(`have length less than or equal to %d`, v)
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
		errStr := fmt.Sprintf(`be greater than %d`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}
	if v := getIntExtention(f, validator.E_IntGte); v >= 0 {
		g.P("if !(", varName, " >= ", v, ") {")
		errStr := fmt.Sprintf(`be greater than or equal to %d`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}
	if v := getIntExtention(f, validator.E_IntLt); v >= 0 {
		g.P("if !(", varName, " < ", v, ") {")
		errStr := fmt.Sprintf(`be less than %d`, v)
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}
	if v := getIntExtention(f, validator.E_IntLte); v >= 0 {
		g.P("if !(", varName, " <= ", v, ") {")
		errStr := fmt.Sprintf(`be less than or equal to %d`, v)
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

	// For repeated messages, we need to run the validator on each message instead of this field
	repeated := f.Desc.Cardinality() == protoreflect.Repeated
	if repeated {
		// Please leave the len check in place here as I have identified edge cases where it is required (will add test case for it later)
		g.P("if len(", varName, ") > 0 {")
		g.P("for _, item := range ", varName, "{")
		varName = "item"
	}

	g.P("if ", varName, " != nil {")
	g.P("if v, ok := interface{}(", varName, ").(", s12protoPackage.Ident("Validator"), "); ok {")
	g.P("if err := v.Validate(); err != nil {")
	g.P("return ", s12protoPackage.Ident("FieldError"), "(\"", f.Desc.Name(), "\", err)")
	g.P("}")
	g.P("}")
	g.P("}")

	if repeated {
		g.P("}") // for
		g.P("}") // if len
	}
}

func genEnumValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	if !getBoolExtension(f, validator.E_EnumRequired) {
		return
	}

	g.P("if int(", varName, ") == 0 {")
	g.P("return ", fmtPackage.Ident("Errorf"), "(\"field ", f.Desc.Name(), " must be specified and a non-zero value\")")
	g.P("}")
}

func printErrorString(g *protogen.GeneratedFile, varName, fieldName, specificErr string, maxLen int) {
	// Do not reflect untrusted value in error, certainly not for sensitive fields like password or PII like email
	g.P(fmtPackage.Ident("Printf"), `("[log-only] %s: value must %s: %s\n", `, `"`, fieldName, `", "`, specificErr, `" `, ", proto.FirstCharactersFromString(", varName, ", ", maxLen, "))")
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
	validator.E_Float,
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

func genFloatValidator(g *protogen.GeneratedFile, f *protogen.Field, varName string) {
	rules := getFloatExtension(f, validator.E_Float)

	if rules.GetOptional() {
		g.P("if ", varName, " != 0 {")
	}

	if !rules.GetAllowNan() {
		g.P("// This statement checks for NaN value without using Math package")
		g.P("if ", varName, " != ", varName, " {")
		errStr := "not be NaN"
		genErrorString(g, varName, string(f.Desc.Name()), errStr)
		g.P("}")
	}

	if rules.GetRange() != "" {
		if !strings.Contains(rules.GetRange(), ":") {
			// Range must contain : to be used
			panic("unparsable range for float validator")
		}

		var rangeVals = strings.Split(rules.GetRange(), ":")
		if len(rangeVals) < 1 || len(rangeVals) > 2 {
			panic("unparseable range for float validator")
		}

		if rangeVals[0] != "" {
			g.P("// Range check lower bounds")
			g.P("if ", varName, " < ", rangeVals[0], " {")
			errStr := fmt.Sprintf("be greater than or equal to %v", rangeVals[0])
			genErrorString(g, varName, string(f.Desc.Name()), errStr)
			g.P("}")
		}

		if rangeVals[1] != "" {
			g.P("// Range check upper bounds")
			g.P("if ", varName, " > ", rangeVals[1], " {")
			errStr := fmt.Sprintf("be less than or equal to %v", rangeVals[1])
			genErrorString(g, varName, string(f.Desc.Name()), errStr)
			g.P("}")
		}
	}

	if rules.GetOptional() {
		g.P("}")
	}
}

func getFloatExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.FloatRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.FloatRules); ok {
			return v
		}
	}
	return nil
}
