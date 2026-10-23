// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"strings"
)

func setOf(groups ...[]string) map[string]bool {
	set := map[string]bool{}
	for _, group := range groups {
		for _, name := range group {
			set[name] = true
		}
	}
	return set
}

var (
	keywordsInDeclarations = []string{
		"associatedtype", "class", "deinit", "enum", "extension",
		"fileprivate", "func", "import", "init", "inout", "internal",
		"let", "open", "operator", "private", "protocol", "public",
		"static", "struct", "subscript", "typealias", "var",
	}
	keywordsInStatements = []string{
		"break", "case",
		"continue", "default", "defer", "do", "else", "fallthrough",
		"for", "guard", "if", "in", "repeat", "return", "switch", "where",
		"while",
	}
	keywordsInExpressionsAndTypes = []string{
		"as",
		"Any", "catch", "false", "is", "nil", "rethrows", "super", "self",
		"Self", "throw", "throws", "true", "try",
	}
	commonTypes = []string{
		"Bool", "Data", "Double", "Float", "Int",
		"Int32", "Int64", "String", "UInt", "UInt32", "UInt64",
	}
	specialVariables = []string{
		"__COLUMN__",
		"__FILE__", "__FUNCTION__", "__LINE__",
	}
)

var reservedTypeNames = setOf(
	[]string{
		"SwiftProtobuf", "Extensions", "protoMessageName",
		"decodeMessage", "traverse", "isInitialized", "unknownFields",
		"debugDescription", "description", "dynamicType", "hashValue",
		"Type", "Protocol", "Swift", "Equatable", "Hashable", "Sendable",
	},
	keywordsInDeclarations, keywordsInStatements, keywordsInExpressionsAndTypes,
	commonTypes, specialVariables,
)

var quotableFieldNames = setOf(keywordsInDeclarations, keywordsInStatements, keywordsInExpressionsAndTypes)

var reservedFieldNames = setOf(
	[]string{
		"isInitialized", "unknownFields",
		"debugDescription", "description", "dynamicType", "hashValue", "init", "self",
		"Type", "Protocol",
	},
	commonTypes, specialVariables,
)

var quotableEnumCases = setOf(
	[]string{"associativity", "dynamicType", "optional", "required"},
	keywordsInDeclarations, keywordsInStatements, keywordsInExpressionsAndTypes,
)

var reservedEnumCases = setOf([]string{
	"allCases", "debugDescription", "description", "dynamicType", "hashValue", "init", "rawValue", "self",
})

func isAllUnderscore(s string) bool {
	return s != "" && strings.Trim(s, "_") == ""
}

func sanitizeTypeName(s, disambiguator string, forbiddenTypeNames map[string]bool) string {
	switch {
	case reservedTypeNames[s]:
		return s + disambiguator
	case isAllUnderscore(s):
		return s + disambiguator
	case strings.HasSuffix(s, disambiguator):
		return sanitizeTypeName(strings.TrimSuffix(s, disambiguator), disambiguator, forbiddenTypeNames) + disambiguator
	case forbiddenTypeNames[s]:
		return s + disambiguator
	}
	return s
}

func isASCUpper(c rune) bool { return c >= 'A' && c <= 'Z' }
func isASCLower(c rune) bool { return c >= 'a' && c <= 'z' }
func isASCDigit(c rune) bool { return c >= '0' && c <= '9' }

func ascUpper(c rune) rune {
	if isASCLower(c) {
		return c - 32
	}
	return c
}

func ascLower(c rune) rune {
	if isASCUpper(c) {
		return c + 32
	}
	return c
}

func isCharacterUppercase(s string, index int) bool {
	runes := []rune(s)
	if index >= len(runes) {
		return false
	}
	return isASCUpper(runes[index])
}

type charClass int

const (
	classDigit charClass = iota
	classLower
	classUpper
	classUnderscore
	classOther
)

func classOf(c rune) charClass {
	switch {
	case isASCDigit(c):
		return classDigit
	case isASCLower(c):
		return classLower
	case isASCUpper(c):
		return classUpper
	case c == '_':
		return classUnderscore
	}
	return classOther
}

var camelCaseAbbreviations = map[string]bool{"url": true, "http": true, "https": true, "id": true}

// camelCase turns a proto identifier into a Swift identifier, with an upper
// case first letter when initialUpperCase is set.
func camelCase(s string, initialUpperCase bool) string {
	var result strings.Builder
	var current []rune
	lastClass := classOf(0)

	addCurrent := func() {
		if len(current) == 0 {
			return
		}
		text := string(current)
		switch {
		case result.Len() == 0 && !initialUpperCase:
		case camelCaseAbbreviations[text]:
			text = strings.ToUpper(text)
		default:
			text = uppercaseFirstCharacter(text)
		}
		result.WriteString(text)
		current = nil
	}

	for _, scalar := range s {
		class := classOf(scalar)
		switch class {
		case classDigit:
			if lastClass != classDigit {
				addCurrent()
			}
			if result.Len() == 0 {
				result.WriteString("_")
			}
			current = append(current, scalar)
		case classUpper:
			if lastClass != classUpper {
				addCurrent()
			}
			current = append(current, ascLower(scalar))
		case classLower:
			if lastClass != classLower && lastClass != classUpper {
				addCurrent()
			}
			current = append(current, scalar)
		case classUnderscore:
			addCurrent()
			if lastClass == classUnderscore {
				result.WriteString("_")
			}
		case classOther:
			addCurrent()
			var escape bool
			if result.Len() == 0 {
				escape = !isSwiftIdentifierHead(scalar)
			} else {
				escape = !isSwiftIdentifierCharacter(scalar)
			}
			if escape {
				fmt.Fprintf(&result, "_u%d", scalar)
			} else {
				current = append(current, scalar)
			}
		}
		lastClass = class
	}
	addCurrent()
	if lastClass == classUnderscore {
		result.WriteString("_")
	}
	return result.String()
}

func toUpperCamelCase(s string) string { return camelCase(s, true) }
func toLowerCamelCase(s string) string { return camelCase(s, false) }

// uppercaseFirstCharacter forces the first character to upper case and keeps the rest.
func uppercaseFirstCharacter(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = ascUpper(runes[0])
	return string(runes)
}

func trimBackticks(s string) string {
	if strings.HasPrefix(s, "`") {
		return s[1 : len(s)-1]
	}
	return s
}

// typePrefix returns the prefix of the type names of a file: the swift_prefix
// option when it is set, otherwise a prefix from the proto package.
func typePrefix(protoPackage string, swiftPrefix *string) string {
	if swiftPrefix != nil {
		return *swiftPrefix
	}
	if protoPackage == "" {
		return ""
	}
	var prefix []rune
	makeUpper := true
	for _, c := range protoPackage {
		switch c {
		case '_':
			makeUpper = true
		case '.':
			makeUpper = true
			prefix = append(prefix, '_')
		default:
			if len(prefix) == 0 && isASCDigit(c) {
				prefix = append(prefix, '_')
			}
			if makeUpper {
				prefix = append(prefix, ascUpper(c))
				makeUpper = false
			} else {
				prefix = append(prefix, c)
			}
		}
	}
	return string(prefix) + "_"
}

// prefixStripper removes an enum name prefix from the names of its values.
type prefixStripper struct {
	prefixChars []rune
}

func newPrefixStripper(prefix string) prefixStripper {
	prefix = strings.ReplaceAll(strings.ToLower(prefix), "_", "")
	return prefixStripper{prefixChars: []rune(prefix)}
}

// strip returns the name without the prefix, or false when the prefix cannot be stripped.
func (s prefixStripper) strip(from string) (string, bool) {
	fromChars := []rune(strings.ToLower(from))
	original := []rune(from)
	pi, fi := 0, 0
	for pi != len(s.prefixChars) {
		if fi == len(fromChars) {
			return "", false
		}
		if fromChars[fi] == '_' {
			fi++
			continue
		}
		if s.prefixChars[pi] != fromChars[fi] {
			return "", false
		}
		pi++
		fi++
	}
	for fi != len(fromChars) && fromChars[fi] == '_' {
		fi++
	}
	if fi == len(fromChars) {
		return "", false
	}
	if !isASCLower(fromChars[fi]) {
		return "", false
	}
	return string(original[fi:]), true
}

func sanitizeMessageName(s string, forbidden map[string]bool) string {
	return sanitizeTypeName(s, "Message", forbidden)
}

func sanitizeEnumName(s string, forbidden map[string]bool) string {
	return sanitizeTypeName(s, "Enum", forbidden)
}

func sanitizeOneofName(s string, forbidden map[string]bool) string {
	return sanitizeTypeName(s, "Oneof", forbidden)
}

func sanitizeFieldNameBasedOn(s, basedOn string) string {
	switch {
	case strings.HasPrefix(basedOn, "clear") && isCharacterUppercase(basedOn, 5):
		return s + "_p"
	case strings.HasPrefix(basedOn, "has") && isCharacterUppercase(basedOn, 3):
		return s + "_p"
	case reservedFieldNames[basedOn]:
		return s + "_p"
	case basedOn == s && quotableFieldNames[basedOn]:
		return "`" + s + "`"
	case isAllUnderscore(basedOn):
		return s + "__"
	}
	return s
}

func sanitizeFieldName(s string) string { return sanitizeFieldNameBasedOn(s, s) }

func sanitizeEnumCaseName(s string) string {
	switch {
	case reservedEnumCases[s]:
		return s + "_"
	case quotableEnumCases[s]:
		return "`" + s + "`"
	case isAllUnderscore(s):
		return s + "__"
	}
	return s
}

// sanitizeMessageScopedExtensionName uses the same sets as the enum cases.
func sanitizeMessageScopedExtensionName(s string) string {
	return sanitizeEnumCaseName(s)
}

// toJSONFieldName must match the function of the SwiftProtobuf runtime.
func toJSONFieldName(s string) string {
	var result []rune
	capitalizeNext := false
	for _, c := range s {
		switch {
		case c == '_':
			capitalizeNext = true
		case capitalizeNext:
			result = append(result, ascUpper(c))
			capitalizeNext = false
		default:
			result = append(result, c)
		}
	}
	return string(result)
}

func isSwiftIdentifierHead(c rune) bool {
	switch {
	case c >= 0x61 && c <= 0x7a, c >= 0x41 && c <= 0x5a:
		return true
	case c == 0x5f:
		return true
	case c == 0xa8, c == 0xaa, c == 0xad, c == 0xaf, c >= 0xb2 && c <= 0xb5, c >= 0xb7 && c <= 0xba:
		return true
	case c >= 0xbc && c <= 0xbe, c >= 0xc0 && c <= 0xd6, c >= 0xd8 && c <= 0xf6, c >= 0xf8 && c <= 0xff:
		return true
	case c >= 0x100 && c <= 0x2ff, c >= 0x370 && c <= 0x167f, c >= 0x1681 && c <= 0x180d, c >= 0x180f && c <= 0x1dbf:
		return true
	case c >= 0x1e00 && c <= 0x1fff:
		return true
	case c >= 0x200b && c <= 0x200d, c >= 0x202a && c <= 0x202e, c == 0x203f, c == 0x2040, c == 0x2054, c >= 0x2060 && c <= 0x206f:
		return true
	case c >= 0x2070 && c <= 0x20cf, c >= 0x2100 && c <= 0x218f, c >= 0x2460 && c <= 0x24ff, c >= 0x2776 && c <= 0x2793:
		return true
	case c >= 0x2c00 && c <= 0x2dff, c >= 0x2e80 && c <= 0x2fff:
		return true
	case c >= 0x3004 && c <= 0x3007, c >= 0x3021 && c <= 0x302f, c >= 0x3031 && c <= 0x303f, c >= 0x3040 && c <= 0xd7ff:
		return true
	case c >= 0xf900 && c <= 0xfd3d, c >= 0xfd40 && c <= 0xfdcf, c >= 0xfdf0 && c <= 0xfe1f, c >= 0xfe30 && c <= 0xfe44:
		return true
	case c >= 0xfe47 && c <= 0xfffd:
		return true
	case c >= 0x10000 && c <= 0xdfffd && c&0xffff <= 0xfffd:
		return true
	case c >= 0xe0000 && c <= 0xefffd:
		return true
	}
	return false
}

func isSwiftIdentifierCharacter(c rune) bool {
	switch {
	case c >= 0x30 && c <= 0x39:
		return true
	case c >= 0x300 && c <= 0x36f, c >= 0x1dc0 && c <= 0x1dff, c >= 0x20d0 && c <= 0x20ff, c >= 0xfe20 && c <= 0xfe2f:
		return true
	}
	return isSwiftIdentifierHead(c)
}

// isValidSwiftIdentifier checks the syntax of an identifier. It does not check for keywords.
func isValidSwiftIdentifier(s string) bool {
	if s == "_" {
		return false
	}
	for i, c := range s {
		if i == 0 {
			if !isSwiftIdentifierHead(c) {
				return false
			}
		} else if !isSwiftIdentifierCharacter(c) {
			return false
		}
	}
	return s != ""
}
