// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	plan "github.com/SafetyCulture/s12-proto/protobuf/validation/plan"
)

// emitField writes the scaffolding that selects the values to check, then the
// steps that run against each of them.
func emitField(e *emitter, f *plan.Field, holder string) error {
	prop := propertyName(f.Proto)
	value := prop
	closes := 0

	if f.RepeatedLenGte >= 0 {
		e.open("if (!(", prop, ".Count >= ", f.RepeatedLenGte, "))")
		e.p("return ValidationError.Create(", quote(f.Name), ", ",
			quote(fmt.Sprintf("length must be greater than or equal to %d", f.RepeatedLenGte)), ");")
		e.shut()
	}

	if f.RepeatedLenLte >= 0 {
		e.open("if (!(", prop, ".Count <= ", f.RepeatedLenLte, "))")
		e.p("return ValidationError.Create(", quote(f.Name), ", ",
			quote(fmt.Sprintf("length must be lesser than or equal to %d", f.RepeatedLenLte)), ");")
		e.shut()
	}

	if f.Loop {
		item := e.local("item")
		e.open("foreach (var ", item, " in ", prop, ")")
		closes++
		// The Go validators rewrite the loop variable, which does not reach the
		// message. Rewrites here are kept local for the same reason, so both languages
		// check the same value and leave the same one in place.
		value = e.local("value")
		e.p("var ", value, " = ", item, ";")
	}

	if f.Unsupported != "" {
		e.p("// ", f.Unsupported)
		for ; closes > 0; closes-- {
			e.shut()
		}
		return nil
	}

	if f.Oneof {
		oneof := f.Proto.Oneof
		e.open("if (", oneofCaseProperty(oneof), " == ", oneofCaseEnum(oneof), ".", pascalCase(string(f.Proto.Desc.Name())), ")")
		closes++
	}

	if f.Optional {
		e.open("if (", hasProperty(f.Proto), ")")
		closes++
	}

	for _, step := range f.Steps {
		if err := emitStep(e, f, step, value, holder); err != nil {
			return err
		}
	}

	for ; closes > 0; closes-- {
		e.shut()
	}
	return nil
}

func emitStep(e *emitter, f *plan.Field, step plan.Step, value, holder string) error {
	switch s := step.(type) {
	case *plan.LegacyString:
		return emitLegacyString(e, f, s, value)
	case *plan.String:
		return emitString(e, f, s, value, holder)
	case *plan.ID:
		emitID(e, f, s, value)
	case *plan.Email:
		emitEmail(e, f, s, value)
	case *plan.Username:
		emitUsername(e, f, s, value)
	case *plan.URL:
		emitURL(e, f, s, value)
	case *plan.Timezone:
		emitTimezone(e, f, s, value)
	case *plan.SimpleString:
		emitSimpleString(e, f, s, value)
	case *plan.Bytes:
		emitBytes(e, f, s, value)
	case *plan.Int:
		emitInt(e, f, s, value)
	case *plan.Number:
		emitNumber(e, f, s, value)
	case *plan.MessageField:
		emitMessageField(e, f, s, value)
	case *plan.EnumField:
		emitEnum(e, f, value)
	default:
		return fmt.Errorf("no C# emitter for validation step %T", step)
	}
	return nil
}

// fail writes the reaction to an unmet requirement: a log line when the field is
// in log-only mode, and a returned error otherwise.
func fail(e *emitter, f *plan.Field, logOnly bool, value, requirement string) {
	if logOnly {
		e.p("ValidationLog.Report(", quote(f.Name), ", ", quote(requirement), ", ",
			"ValidatorHelpers.TruncateAndEncode(", stringify(f, value), ", ", logOnlyMaxLength, "));")
		return
	}
	e.p("return ValidationError.Create(", quote(f.Name), ", ", quote(requirement), ");")
}

// stringify renders the value as text for a log line.
func stringify(f *plan.Field, value string) string {
	if f.Proto.Desc.Kind() == protoreflect.StringKind {
		return value
	}
	return value + ".ToString()"
}

func emitLegacyString(e *emitter, f *plan.Field, s *plan.LegacyString, value string) error {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	}

	if s.MatchRegex {
		translated, err := translateFreeRegex(s.Regex)
		if err != nil {
			return fmt.Errorf("field %s: (validator.regex) %q: %w", f.Name, s.Regex, err)
		}
		name := e.pattern(translated)
		e.open("if (!", e.holder, ".", name, ".IsMatch(", value, "))")
		fail(e, f, false, value, "value must be a string conforming to predefined pattern")
		e.shut()
	}

	if s.UUID {
		e.open("if (!ValidatorHelpers.IsUuid(", value, "))")
		fail(e, f, false, value, "value must be parsable as a UUID")
		e.shut()
	}

	if s.LegacyID {
		e.open("if (!ValidatorHelpers.IsLegacyId(", value, ", false))")
		fail(e, f, false, value, "value must be parsable as a UUID or a legacy ID")
		e.shut()
	}

	emitLength(e, f, s.Length, value, byteLength(f, value))

	if s.Optional {
		e.shut()
	}
	return nil
}

// byteLength measures a value the way Go's len() does.
func byteLength(f *plan.Field, value string) func(string) string {
	if f.Proto.Desc.Kind() == protoreflect.BytesKind {
		return func(v string) string { return v + ".Length" }
	}
	return func(v string) string { return "ValidatorHelpers.Utf8Length(" + v + ")" }
}

func emitLength(e *emitter, f *plan.Field, l *plan.Length, value string, measure func(string) string) {
	if l == nil {
		return
	}

	if l.TrimFirst {
		trimmed := e.local("trim")
		e.p("var ", trimmed, " = StringMutators.TrimSpace(", value, ");")
		value = trimmed
	}

	if l.Gte >= 0 {
		e.open("if (!(", measure(value), " >= ", l.Gte, "))")
		fail(e, f, false, value, fmt.Sprintf("value must have length greater than or equal to %d", l.Gte))
		e.shut()
	}

	if l.Lte >= 0 {
		e.open("if (!(", measure(value), " <= ", l.Lte, "))")
		fail(e, f, false, value, fmt.Sprintf("value must have length less than or equal to %d", l.Lte))
		e.shut()
	}
}

func emitSimpleString(e *emitter, f *plan.Field, s *plan.SimpleString, value string) {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	}

	if s.Bounded {
		length := e.local("length")
		e.p("var ", length, " = ValidatorHelpers.RuneCount(", value, ");")
		e.open("if (", length, " > ", s.MaxLen, " || ", length, " < ", s.MinLen, ")")
		fail(e, f, s.LogOnly, value, fmt.Sprintf("value must have a length between %d and %d", s.MinLen, s.MaxLen))
		e.shut()
	}

	if s.Optional {
		e.shut()
	}
}

func emitString(e *emitter, f *plan.Field, s *plan.String, value, holder string) error {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	}

	for _, op := range s.Ops {
		if err := emitStringOp(e, f, op, value, holder); err != nil {
			return err
		}
	}

	if s.Optional {
		e.shut()
	}
	return nil
}

// The three encoding messages below repeat "must". The Go validators produce them
// that way, callers already match on the text, and a port that reads differently
// from the implementation it is checked against is not a port. Correcting them is
// a change to every service at once, tracked separately.
func emitStringOp(e *emitter, f *plan.Field, op plan.Op, value, holder string) error {
	switch o := op.(type) {
	case *plan.OpNormalizeNFC:
		normalised := e.local("nfc")
		e.open("if (StringMutators.TryNormalizeToNfc(", value, ", out var ", normalised, "))")
		e.p(value, " = ", normalised, ";")
		e.shut()
		e.open("else")
		fail(e, f, o.LogOnly, value, "value must must be normalisable to NFC")
		e.shut()

	case *plan.OpCheckEncoding:
		e.open("if (StringMutators.ContainsReplacementCharacter(", value, "))")
		fail(e, f, o.LogOnly, value, "value must must have valid encoding")
		e.shut()
		e.open("else if (!ValidatorHelpers.IsWellFormedUtf16(", value, "))")
		fail(e, f, o.LogOnly, value, "value must must be a valid UTF-8-encoded string")
		e.shut()

	case *plan.OpRejectURL:
		e.open("if (ValidatorHelpers.ContainsUrl(", value, "))")
		fail(e, f, false, value, "value must not contain a URL")
		e.shut()

	case *plan.OpBreakPartialURL:
		e.p(value, " = StringMutators.BreakPartialUrls(", value, ");")

	case *plan.OpTrim:
		e.p(value, " = StringMutators.TrimSpace(", value, ");")

	case *plan.OpReplaceLiteral:
		e.p(value, " = ", value, ".Replace(", charStringLiteral(o.From), ", \"", o.To, "\");")

	case *plan.OpReplaceUnsafe:
		e.p(value, " = StringMutators.ReplaceUnsafeCharacters(", value, ");")

	case *plan.OpStripCR:
		e.p(value, " = ", value, `.Replace("\r", "", StringComparison.Ordinal);`)

	case *plan.OpReplaceOther:
		if o.Multiline {
			e.p(value, " = StringMutators.ReplaceSymbolCharactersMultiline(", value, ");")
		} else {
			e.p(value, " = StringMutators.ReplaceSymbolCharacters(", value, ");")
		}

	case *plan.OpStripPUA:
		e.p(value, " = StringMutators.StripPrivateUseArea(", value, ");")

	case *plan.OpLength:
		length := e.local("len")
		if o.Runes {
			e.p("var ", length, " = ValidatorHelpers.RuneCount(", value, ");")
		} else {
			e.p("var ", length, " = ValidatorHelpers.Utf8Length(", value, ");")
		}
		if o.Min == o.Max {
			e.open("if (!(", length, " == ", o.Min, "))")
			fail(e, f, o.LogOnly, value, fmt.Sprintf("value must have length %d", o.Min))
		} else {
			e.open("if (!(", length, " >= ", o.Min, " && ", length, " <= ", o.Max, "))")
			fail(e, f, o.LogOnly, value, fmt.Sprintf("value must have a length between %d and %d", o.Min, o.Max))
		}
		e.shut()

	case *plan.OpHasPrefix:
		e.open("if (!", value, ".StartsWith(", quote(o.Prefix), ", StringComparison.Ordinal))")
		fail(e, f, false, value, "value must start with prefix")
		e.shut()

	case *plan.OpMatchPattern:
		translated, err := translatePattern(o.Pattern)
		if err != nil {
			return fmt.Errorf("field %s: %w", f.Name, err)
		}
		name := e.pattern(translated.Pattern)
		e.open("if (!CharacterClass.IsMatch(", holder, ".", name, ", ", value, ", ",
			astralExpression(translated.AstralCategories), "))")
		fail(e, f, o.LogOnly, value, "value must only have valid characters")
		e.shut()

	default:
		return fmt.Errorf("no C# emitter for string operation %T", op)
	}
	return nil
}

// astralExpression names the categories whose code points above the Basic
// Multilingual Plane the pattern is meant to admit.
func astralExpression(categories []string) string {
	if len(categories) == 0 {
		return "UnicodeCategories.None"
	}
	parts := make([]string, 0, len(categories))
	for _, category := range categories {
		parts = append(parts, "UnicodeCategories."+category)
	}
	return strings.Join(parts, " | ")
}

func emitID(e *emitter, f *plan.Field, s *plan.ID, value string) {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	}

	switch s.Version {
	case "v4":
		e.open("if (!ValidatorHelpers.IsUuidV4(", value, "))")
	default:
		e.open("if (!ValidatorHelpers.IsUuid(", value, "))")
	}

	valid := e.local("validId")
	e.p("var ", valid, " = false;")

	if s.Legacy {
		e.open("if (ValidatorHelpers.IsLegacyId(", value, ", ", s.LowercaseOnly, "))")
		e.p(valid, " = true;")
		e.shut()
	}

	if s.S12 {
		e.open("if (!", valid, " && ValidatorHelpers.IsS12Id(", value, ", ", s.LowercaseOnly, "))")
		e.p(valid, " = true;")
		e.shut()
	}

	if s.LongPrefixed {
		e.open("if (!", valid, " && ValidatorHelpers.IsLongPrefixedLegacyId(", value, "))")
		e.p(valid, " = true;")
		e.shut()
	}

	e.open("if (!", valid, ")")
	fail(e, f, s.LogOnly, value, "value must "+s.Requirement)
	e.shut()

	e.shut()

	if s.Optional {
		e.shut()
	}
}

func emitEmail(e *emitter, f *plan.Field, s *plan.Email, value string) {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	}

	e.open("if (!ValidatorHelpers.IsValidEmail(", value, "))")
	// Should not return the address, which is PII.
	fail(e, f, false, value, "value must be parsable as an email address")
	e.shut()

	if s.Optional {
		e.shut()
	}
}

func emitUsername(e *emitter, f *plan.Field, s *plan.Username, value string) {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	}

	e.open("if (!ValidatorHelpers.IsValidEmail(", value, "))")
	valid := e.local("validUsername")
	e.p("var ", valid, " = false;")
	if s.AllowNonEmail {
		e.open("if (ValidatorHelpers.IsValidNonEmail(", value, "))")
		e.p(valid, " = true;")
		e.shut()
	}
	e.open("if (!", valid, ")")
	fail(e, f, s.LogOnly, value, "value must be a valid email address or username")
	e.shut()
	e.shut()

	if s.Optional {
		e.shut()
	}
}

func emitURL(e *emitter, f *plan.Field, s *plan.URL, value string) {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	}

	schemes := make([]string, 0, len(s.Schemes))
	for _, scheme := range s.Schemes {
		schemes = append(schemes, quote(scheme))
	}

	failure := e.local("urlFailure")
	e.open("if (!ValidatorHelpers.TryValidateUrl(", value, ", new[] { ", strings.Join(schemes, ", "), " }, ",
		s.AllowFragment, ", out var ", failure, "))")
	e.p("return ValidationError.Create(", quote(f.Name), ", $\"value must be parsable as a URL: {", failure, "}\");")
	e.shut()

	if s.Optional {
		e.shut()
	}
}

func emitTimezone(e *emitter, f *plan.Field, s *plan.Timezone, value string) {
	if s.Optional {
		e.open(`if (`, value, ` != "")`)
	} else {
		e.open(`if (`, value, ` == "")`)
		e.p("return ValidationError.Required(", quote(f.Name), ");")
		e.shut()
	}

	e.open("if (!ValidatorHelpers.IsValidTimeZone(", value, "))")
	fail(e, f, false, value, "value must be a valid IANA TZ database value")
	e.shut()

	if s.Optional {
		e.shut()
	}
}

func emitBytes(e *emitter, f *plan.Field, s *plan.Bytes, value string) {
	if s.Optional {
		e.open("if (", value, ".Length > 0)")
	}

	if s.UUIDSize {
		e.open("if (", value, ".Length != 16)")
		fail(e, f, false, value, "value must be exactly 16 bytes long to be a valid UUID")
		e.shut()
	}

	emitLength(e, f, s.Length, value, byteLength(f, value))

	if s.Optional {
		e.shut()
	}
}

func emitInt(e *emitter, f *plan.Field, s *plan.Int, value string) {
	if s.Optional {
		e.open("if (", value, " != 0)")
	}

	bounds := []struct {
		limit int64
		op    string
		text  string
	}{
		{s.Gt, ">", "be greater than %d"},
		{s.Gte, ">=", "be greater than or equal to %d"},
		{s.Lt, "<", "be less than %d"},
		{s.Lte, "<=", "be less than or equal to %d"},
	}
	for _, bound := range bounds {
		if bound.limit < 0 {
			continue
		}
		e.open("if (!(", value, " ", bound.op, " ", bound.limit, "))")
		fail(e, f, false, value, "value must "+fmt.Sprintf(bound.text, bound.limit))
		e.shut()
	}

	if s.Optional {
		e.shut()
	}
}

func emitNumber(e *emitter, f *plan.Field, s *plan.Number, value string) {
	if s.Optional {
		e.open("if (", value, " != 0)")
	}

	if s.RejectNaN {
		e.p("// This statement checks for NaN without a call into Math.")
		e.open("if (", value, " != ", value, ")")
		fail(e, f, false, value, "value must not be NaN")
		e.shut()
	}

	if s.Min != "" {
		e.p("// Range check lower bounds")
		e.open("if (", value, " < ", numericLiteral(f, s.Min), ")")
		fail(e, f, s.LogOnly, value, "value must be greater than or equal to "+s.Min)
		e.shut()
	}

	if s.Max != "" {
		e.p("// Range check upper bounds")
		e.open("if (", value, " > ", numericLiteral(f, s.Max), ")")
		fail(e, f, s.LogOnly, value, "value must be less than or equal to "+s.Max)
		e.shut()
	}

	if s.Optional {
		e.shut()
	}
}

// numericLiteral suffixes a bound so it compares against the field's own type
// rather than promoting it.
func numericLiteral(f *plan.Field, bound string) string {
	switch f.Proto.Desc.Kind() {
	case protoreflect.DoubleKind:
		return bound + "d"
	case protoreflect.FloatKind:
		return bound + "f"
	default:
		return bound
	}
}

func emitMessageField(e *emitter, f *plan.Field, s *plan.MessageField, value string) {
	if s.Required {
		e.open("if (", value, " == null)")
		e.p("return ValidationError.Required(", quote(f.Name), ");")
		e.shut()
	}

	if s.Repeated {
		e.open("if (", value, ".Count > 0)")
		item := e.local("item")
		e.open("foreach (var ", item, " in ", value, ")")
		value = item
	}

	validatable := e.local("validatable")
	// Cast through object first. A well-known wrapper type surfaces in C# as the
	// primitive it wraps, and a message from another assembly does not implement the
	// interface, so a direct pattern match on those is a compile error rather than a
	// test that fails at runtime.
	e.open("if ((object?)", value, " is IValidatableMessage ", validatable, ")")
	result := e.local("error")
	e.p("var ", result, " = ", validatable, ".Validate();")
	e.open("if (", result, " != null)")
	e.p("return ", result, ".Nest(", quote(f.Name), ");")
	e.shut()
	e.shut()

	if s.Repeated {
		e.shut()
		e.shut()
	}
}

func emitEnum(e *emitter, f *plan.Field, value string) {
	e.open("if ((int)", value, " == 0)")
	e.p("return ValidationError.Create(", quote(f.Name), ", ",
		quote("must be specified and a non-zero value"), ");")
	e.shut()
}

// quote renders a C# string literal.
func quote(s string) string {
	var out strings.Builder
	out.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			out.WriteString(`\"`)
		case '\\':
			out.WriteString(`\\`)
		case '\n':
			out.WriteString(`\n`)
		case '\r':
			out.WriteString(`\r`)
		case '\t':
			out.WriteString(`\t`)
		default:
			if r < 0x20 || r > 0x7E {
				fmt.Fprintf(&out, `\u%04X`, r)
				continue
			}
			out.WriteRune(r)
		}
	}
	out.WriteByte('"')
	return out.String()
}

// charLiteral renders a C# char literal.
func charLiteral(r rune) string {
	return fmt.Sprintf(`'\u%04X'`, r)
}

// charStringLiteral writes a character as a one-character C# string, for a
// replacement whose result may be any number of characters.
func charStringLiteral(r rune) string {
	return fmt.Sprintf(`"\u%04X"`, r)
}
