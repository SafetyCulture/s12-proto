// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

// Any changes made in this file will require Security Engineering consultation/review
package plan

import (
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

// Reset clears the state a generation run accumulates. A run that does not call
// it first inherits the pattern pool the previous run left behind.
func Reset() {
	regexLib = make(map[string]map[string]bool)
}

// BuildFile reads the validator options on file and returns the plan for it.
//
// The allow-list pattern pool is shared across every call in a run, so files
// must be built in the order they are emitted, and Reset separates one run from
// the next.
func BuildFile(file *protogen.File) *File {
	prepareStringGenerics()

	plan := &File{Proto: file}
	for _, msg := range file.Messages {
		if opts, ok := msg.Desc.Options().(*descriptorpb.MessageOptions); !ok || opts.GetMapEntry() {
			continue
		}
		plan.Messages = append(plan.Messages, buildMessage(msg))
	}
	return plan
}

// prepareStringGenerics seeds the pattern pool with the characters the replacers
// can produce, and folds the safe defaults into the unsafe ones.
func prepareStringGenerics() {
	prepareStringReplacerRegex(stringUnsafeReplacerMap, "replacer_unsafe_allowed")
	prepareStringReplacerRegex(stringSymbolReplacerMap, "replacer_symbol_allowed")

	stringReDefaultUnsafe = append(append([]string{}, stringReDefaultSafe...), stringReDefaultUnsafeExtra...)
}

func buildMessage(msg *protogen.Message) *Message {
	m := &Message{Proto: msg, IsMapEntry: msg.Desc.IsMapEntry()}

	for _, field := range msg.Fields {
		if regex := getRegexValue(field); regex != "" {
			m.LegacyRegex = append(m.LegacyRegex, &LegacyRegexField{Proto: field, Pattern: regex})
		}
	}

	for _, f := range msg.Fields {
		if field := buildField(f); field != nil {
			m.Fields = append(m.Fields, field)
		}
	}

	for _, inner := range msg.Messages {
		m.Nested = append(m.Nested, buildMessage(inner))
	}
	return m
}

func buildField(f *protogen.Field) *Field {
	hasExt := hasValidationExtensions(f)
	if !hasExt && f.Message == nil {
		return nil
	}

	hasRepeatedExt := hasExt && f.Desc.Cardinality() == protoreflect.Repeated
	isMessageField := f.Desc.Kind() == protoreflect.MessageKind

	field := &Field{
		Proto:          f,
		Name:           string(f.Desc.Name()),
		RepeatedLenGte: Unset,
		RepeatedLenLte: Unset,
		// Only loop when there is per-element work to do. An array of messages is
		// handled by the message step instead.
		Loop:     hasRepeatedExt && !isMessageField && hasNonRepeatedValidationExtensions(f),
		Oneof:    f.Oneof != nil && !f.Oneof.Desc.IsSynthetic(),
		Optional: f.Desc.HasOptionalKeyword(),
	}

	if hasRepeatedExt {
		field.RepeatedLenGte = getIntExtention(f, validator.E_RepeatedLenGte)
		field.RepeatedLenLte = getIntExtention(f, validator.E_RepeatedLenLte)
	}

	if f.Desc.IsMap() {
		field.Unsupported = "Validation of proto3 map<> fields is unsupported."
		return field
	}

	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		// The options that predate (validator.string) can still be combined with it,
		// so both rule sets are planned for the same field.
		if s := buildLegacyString(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildString(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildID(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildEmail(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildUsername(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildURL(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildTimezone(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildSimpleString(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
	case protoreflect.BytesKind:
		if s := buildBytes(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind:
		if s := buildInt(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
		if s := buildNumber(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
	case protoreflect.MessageKind:
		field.Steps = append(field.Steps, buildMessageField(f))
	case protoreflect.EnumKind:
		if s := buildEnum(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		if s := buildNumber(f); s != nil {
			field.Steps = append(field.Steps, s)
		}
	}

	return field
}

func buildLegacyString(f *protogen.Field) *LegacyString {
	s := &LegacyString{
		Optional:   getBoolExtension(f, validator.E_Optional),
		Regex:      getRegexValue(f),
		MatchRegex: getRegexValue(f) != "",
		UUID:       getBoolExtension(f, validator.E_Uuid),
		LegacyID:   getBoolExtension(f, validator.E_LegacyId),
		Length:     buildLength(f),
	}
	if !s.Optional && !s.MatchRegex && !s.UUID && !s.LegacyID && s.Length == nil {
		return nil
	}
	return s
}

func buildLength(f *protogen.Field) *Length {
	l := &Length{
		TrimFirst: getBoolExtension(f, validator.E_TrimLenCheck),
		Gte:       getIntExtention(f, validator.E_LengthGte),
		Lte:       getIntExtention(f, validator.E_LengthLte),
	}
	if !l.TrimFirst && l.Gte < 0 && l.Lte < 0 {
		return nil
	}
	return l
}

func buildSimpleString(f *protogen.Field) *SimpleString {
	rules := getSimpleStringExtension(f, validator.E_SimpleString)
	if rules == nil {
		return nil
	}
	return &SimpleString{
		Optional: rules.GetOptional(),
		Bounded:  rules.GetMinLen() >= 1 || rules.GetMaxLen() >= 1,
		MinLen:   rules.GetMinLen(),
		MaxLen:   rules.GetMaxLen(),
		LogOnly:  rules.GetLogOnly(),
	}
}

func buildString(f *protogen.Field) *String {
	// `string` and `unsafe_string` share most of the underlying validation logic
	stringType := validator.E_String
	rules := getStringExtension(f, stringType)
	if rules == nil {
		stringType = validator.E_UnsafeString
		rules = getStringExtension(f, stringType)
	}
	if rules == nil {
		return nil
	}

	// The allow list is keyed by the generated name, so fields sharing a name also
	// share the characters each of them allows.
	allowListReId := "re_" + f.GoIdent.GoName
	if stringType == validator.E_String {
		prepareRegex(allowListReId, stringReDefaultSafe...)
	} else {
		prepareRegex(allowListReId, stringReDefaultUnsafe...)
	}

	s := &String{Optional: rules.GetOptional()}

	// ### STEP 1: decode
	// TODO PA: Check if we can confirm that data has been decoded, things like %20 might indicate otherwise

	// ### STEP 2: normalise and canonicalise
	// Normalisation for file paths/URLs should not happen here but in URL validator

	// ##### 2A. Normalise Unicode NFD strings to NFC
	//  e.g. NFD string mañana should be normalised to mañana (n + ̃ = ñ)
	// TODO PA: Check if we need support for NFKC/NFKD
	// TODO PA: Do we want to have support for U+202E RIGHT-TO-LEFT OVERRIDE by default for all strings?
	s.Ops = append(s.Ops, &OpNormalizeNFC{LogOnly: rules.GetLogOnly()})

	// ##### 2B. Check for encoding issues
	if rules.GetValidateEncoding() {
		s.Ops = append(s.Ops, &OpCheckEncoding{LogOnly: rules.GetLogOnly()})
	} else {
		// U+FFFD is in Symbol_So category which must be allowed in case validateInvalidEncoding is disabled (to detect issues)
		prepareRegex(allowListReId, `\x{FFFD}`)
	}

	// If the text contains a full URL anywhere in it, reject it. This should help us control spam.
	if rules.GetRejectUrl() {
		s.Ops = append(s.Ops, &OpRejectURL{})
	}

	// ### STEP 3: sanitise
	// #### 3A-1. 'Break' partial URLs by introducing a space after each dot between characters.
	// Note this will PERMANENTLY mutate the message field data. Further iterations will ignore these '. ' patterns,
	// we only care about patterns like 'a.b'
	if rules.GetBreakPartialUrl() {
		s.Ops = append(s.Ops, &OpBreakPartialURL{})
	}

	// ##### 3A-2. sanitise whitespace (trim option)
	// WARNING: Any leading/trailing whitespace will be permanently removed (this is intended here)
	if rules.GetTrim() {
		s.Ops = append(s.Ops, &OpTrim{})
	}

	// ##### 3B. replace restricted characters with a safe alternative
	// WARNING: this could corrupt data, eg. hyperlinks in a description box so disabled by default
	// Ensure to test the alternative characters in all output (database, email, pdf, csv, web UI, mobile UI)
	if rules.GetReplaceUnsafe() {
		// If allow option is defined, only chars from allow list will be allowed
		// Otherwise, all chars in the replace map are accepted
		if rules.GetAllow() != "" {
			for _, runeValue := range rules.GetAllow() {
				unicodeValue := unicodeKey(runeValue)
				replacedUnicodeValue, isUnsafe := stringUnsafeReplacerMap[`\u`+unicodeValue]
				if !isUnsafe {
					continue
				}
				// need to replace this one, replace case by case instead of using the replacer
				s.Ops = append(s.Ops, &OpReplaceLiteral{From: runeValue, To: replacedUnicodeValue})

				// and add it to the regex
				prepareRegex(allowListReId, `\x{`+strings.Replace(replacedUnicodeValue, `\u`, ``, 1)+`}`)
			}
		} else {
			// Add the possible replaced characters to the allow list as they are not all allowed by default
			mergeRegex(allowListReId, "replacer_unsafe_allowed")
			s.Ops = append(s.Ops, &OpReplaceUnsafe{})
		}
	}

	// ##### 3C. strip carriage return characters \r in multiline strings (leave \n)
	if rules.GetMultiline() {
		s.Ops = append(s.Ops, &OpStripCR{})
	}

	// ##### 3D. replace rare symbols with a more common alternative
	if rules.GetReplaceOther() {
		// Add the possible replaced characters to the allow list as they might not be allowed by default
		mergeRegex(allowListReId, "replacer_symbol_allowed")
		s.Ops = append(s.Ops, &OpReplaceOther{Multiline: rules.GetMultiline()})
	}

	// ##### 3E. Sanitise (remove) Private Use Area Codepoints in the Basic Multilingual Plane
	// Note that we do not check for PUA in planes 15 and 16 currently
	if rules.GetSanitisePua() {
		s.Ops = append(s.Ops, &OpStripPUA{})
	}

	// ### STEP 4: validate
	// ##### 4A. validate length
	// Length is checked before the pattern so that oversized input is rejected early
	// and so that the pattern can be shared between fields.
	minLen, maxLen := stringLengthBounds(f, rules, stringType)
	s.Ops = append(s.Ops, &OpLength{
		Runes:   rules.GetRunes(),
		Min:     minLen,
		Max:     maxLen,
		LogOnly: rules.GetLogOnly(),
	})

	// ##### 4B. prepare whitelist/regex
	if rules.GetAllow() != "" {
		// Note that this does inentionally not support a regex pattern or entire categories like \pX
		// All characters in the string are considered as LITERAL value
		for _, runeValue := range rules.GetAllow() {
			unicodeValue := unicodeKey(runeValue)
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
					panic("invalid allow character in field " + f.GoIdent.GoName + ": U+" + unicodeValue + ", this character is potentially dangerous (check charDenyList)")
				}
			}
			prepareRegex(allowListReId, `\x{`+unicodeValue+`}`)
		}
	}

	// Check for multiline option which allows linebreaks
	if rules.GetMultiline() {
		// \r was stripped in step 3C
		prepareRegex(allowListReId, stringReLineBreaks...)
	}

	// Add symbols if defined
	symbols := rules.GetSymbols()
	if len(symbols) > 0 {
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
			fmt.Fprintf(os.Stderr, "WARN: Symbol %v not in stringSymbolMap (not implemented)\n", symbol)
		}
	}

	if rules.GetPrefix() != "" {
		s.Ops = append(s.Ops, &OpHasPrefix{Prefix: rules.GetPrefix()})
	}

	// ##### 4C. validate string against the whitelist/regex
	pattern, err := getPreparedRegex(allowListReId)
	if err != nil {
		// Should not continue as we don't have a regex pattern
		panic("Error generating regex for " + allowListReId + "/" + f.GoIdent.GoName + ": " + fmt.Sprint(err))
	}
	s.Ops = append(s.Ops, &OpMatchPattern{Pattern: pattern, LogOnly: rules.GetLogOnly()})

	return s
}

// stringLengthBounds resolves the length bounds for a field, defaulting anything
// the option leaves out and rejecting a definition outside what the string type
// permits.
func stringLengthBounds(f *protogen.Field, rules *validator.StringRules, stringType protoreflect.ExtensionType) (uint32, uint32) {
	// Always have a min and max length: either default vals or set in the field option
	// This ensures that we no longer accept any length string (safe by default)
	minLen := stringLenMinDefault
	maxLen := stringLenMaxDefault

	if fLen := rules.GetLen(); fLen != "" {
		fLenChunks := strings.SplitN(fLen, ":", 3)
		fMinLen, fMaxLen := -1, -1 // set -1 as default to distinguish between 0 and unset
		switch len(fLenChunks) {
		case 2:
			// min:max notation deffiend, e.g. len: "2:20"
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
			panic("unparsable string validator value for len in field " + f.GoIdent.GoName + ": expected x-y or x, found " + fLen)
		}

		if fMinLen > 0 {
			minLen = uint32(fMinLen)
		}
		if fMaxLen > 0 {
			maxLen = uint32(fMaxLen)
		}
	}

	validMin := stringLenMinSafe
	validMax := stringLenMaxSafe
	if stringType == validator.E_UnsafeString {
		validMin = stringLenMinUnsafe
		validMax = stringLenMaxUnsafe
	}
	if minLen < validMin || maxLen > validMax {
		panic("invalid string validator value for len in field " + f.GoIdent.GoName + ": expected " + fmt.Sprint(validMin) + "<=x<=" + fmt.Sprint(validMax) + ", found " + fmt.Sprint(minLen) + "-" + fmt.Sprint(maxLen))
	}

	return minLen, maxLen
}

func buildID(f *protogen.Field) *ID {
	rules := getIdExtension(f, validator.E_Id)
	if rules == nil {
		return nil
	}

	id := &ID{
		Optional:      rules.GetOptional(),
		Version:       rules.GetVersion(),
		Legacy:        rules.GetLegacy(),
		S12:           rules.GetS12Id(),
		LowercaseOnly: rules.GetLowercaseOnly(),
		LogOnly:       rules.GetLogOnly(),
	}

	switch id.Version {
	case "v4":
		id.Requirement = "be parsable as UUIDv4"
	case "any":
		id.Requirement = "be parsable as UUID"
	default:
		// Unsupported version; do not generate the validators without having an implementation for this version
		panic("unsupported UUID version in field " + f.GoIdent.GoName + ", got: " + id.Version)
	}

	// UUID format is always accepted; the legacy formats are opt-in and the value
	// passes when any one of the enabled checks passes.
	if id.Legacy {
		id.Requirement += " or legacy ID"
	}
	if id.S12 {
		id.Requirement += " or S12 ID"
	}
	// Prefixed, long legacy format Ids require both the s12 and legacy options
	id.LongPrefixed = id.Legacy && id.S12

	return id
}

func buildEmail(f *protogen.Field) *Email {
	rules := getEmailExtension(f, validator.E_Email)
	if rules == nil {
		return nil
	}
	return &Email{Optional: rules.GetOptional()}
}

func buildUsername(f *protogen.Field) *Username {
	rules := getUsernameExtension(f, validator.E_Username)
	if rules == nil {
		return nil
	}
	return &Username{
		Optional:      rules.GetOptional(),
		AllowNonEmail: rules.GetAllowNonEmail(),
		LogOnly:       rules.GetLogOnly(),
	}
}

func buildURL(f *protogen.Field) *URL {
	rules := getURLExtension(f, validator.E_Url)
	if rules == nil {
		return nil
	}

	schemes := rules.GetSchemes()
	if len(schemes) == 0 {
		schemes = append(schemes, "https")
	}
	if rules.GetAllowHttp() {
		schemes = append(schemes, "http")
	}

	return &URL{
		Optional:      rules.GetOptional(),
		Schemes:       schemes,
		AllowFragment: rules.GetAllowFragment(),
	}
}

func buildTimezone(f *protogen.Field) *Timezone {
	rules := getTimezoneExtension(f, validator.E_Timezone)
	if rules == nil {
		return nil
	}
	return &Timezone{Optional: rules.GetOptional()}
}

func buildBytes(f *protogen.Field) *Bytes {
	b := &Bytes{
		Optional: getBoolExtension(f, validator.E_Optional),
		UUIDSize: getBoolExtension(f, validator.E_Uuid),
		Length:   buildLength(f),
	}
	if !b.Optional && !b.UUIDSize && b.Length == nil {
		return nil
	}
	return b
}

func buildInt(f *protogen.Field) *Int {
	i := &Int{
		Optional: getBoolExtension(f, validator.E_Optional),
		Gt:       getIntExtention(f, validator.E_IntGt),
		Gte:      getIntExtention(f, validator.E_IntGte),
		Lt:       getIntExtention(f, validator.E_IntLt),
		Lte:      getIntExtention(f, validator.E_IntLte),
	}
	if !i.Optional && i.Gt < 0 && i.Gte < 0 && i.Lt < 0 && i.Lte < 0 {
		return nil
	}
	return i
}

func buildNumber(f *protogen.Field) *Number {
	rules := getNumberExtension(f, validator.E_Number)

	n := &Number{
		Optional:  rules.GetOptional(),
		RejectNaN: !rules.GetAllowNan(),
		LogOnly:   rules.GetLogOnly(),
	}

	if n.RejectNaN && f.Desc.Kind() != protoreflect.DoubleKind && f.Desc.Kind() != protoreflect.FloatKind {
		panic("cannot use allow_nan option for integers, only supported for float/double")
	}

	if r := rules.GetRange(); r != "" {
		if !strings.Contains(r, ":") {
			// Range must contain : to be used
			panic("unparsable range for number validator")
		}
		rangeVals := strings.Split(r, ":")
		if len(rangeVals) < 1 || len(rangeVals) > 2 {
			panic("unparsable range for number validator")
		}
		n.Min = rangeVals[0]
		n.Max = rangeVals[1]
	}

	if !n.Optional && !n.RejectNaN && n.Min == "" && n.Max == "" {
		return nil
	}
	return n
}

func buildMessageField(f *protogen.Field) *MessageField {
	return &MessageField{
		Required: getBoolExtension(f, validator.E_MsgRequired),
		Repeated: f.Desc.Cardinality() == protoreflect.Repeated,
	}
}

func buildEnum(f *protogen.Field) *EnumField {
	if !getBoolExtension(f, validator.E_EnumRequired) {
		return nil
	}
	return &EnumField{}
}

// unicodeKey renders a rune the way the replacer and deny tables are keyed.
func unicodeKey(r rune) string {
	return strings.Replace(fmt.Sprintf("%U", r), "U+", "", 1)
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
	validator.E_Username,
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

func getUsernameExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.UsernameRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.UsernameRules); ok {
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

func getNumberExtension(f *protogen.Field, xt protoreflect.ExtensionType) *validator.NumberRules {
	if opts := f.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(*validator.NumberRules); ok {
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
			return Unset
		}
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.(int64); ok {
			return v
		}
	}
	return Unset
}
