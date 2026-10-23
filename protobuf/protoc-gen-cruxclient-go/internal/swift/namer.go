// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"strconv"
	"strings"
)

// namer computes the Swift names of the types, cases and properties of a file.
// All the files of a run belong to one Swift module. Only the files that the
// runtime bundles live in the module of the runtime, so a type from such a file
// needs the module as a prefix, unless the file of the namer is one of them.
type namer struct {
	runtimeModule      string
	currentIsBundled   bool
	enumValueNameCache map[string]string
}

// newNamer returns the namer for the Swift code of one file.
func newNamer(runtimeModule string, current *file) *namer {
	return &namer{runtimeModule: runtimeModule, currentIsBundled: current.isBundledProto(), enumValueNameCache: map[string]string{}}
}

// runtimePrefix is the qualifier of the runtime types.
func (n *namer) runtimePrefix() string {
	if n.currentIsBundled {
		return ""
	}
	return n.runtimeModule + "."
}

// modulePrefix is the qualifier of the types of a file.
func (n *namer) modulePrefix(f *file) string {
	if f.isBundledProto() && !n.currentIsBundled {
		return n.runtimeModule + "."
	}
	return ""
}

func (n *namer) forbidden() map[string]bool {
	return map[string]bool{n.runtimeModule: true}
}

func (n *namer) filePrefix(f *file) string {
	return typePrefix(f.pkg, f.swiftPrefix())
}

func (n *namer) relativeNameMessage(m *message) string {
	if m.containing != nil {
		return sanitizeMessageName(m.name, n.forbidden())
	}
	return sanitizeMessageName(n.filePrefix(m.file)+m.name, n.forbidden())
}

func (n *namer) fullNameMessage(m *message) string {
	relative := n.relativeNameMessage(m)
	if m.containing == nil {
		return n.modulePrefix(m.file) + relative
	}
	return n.fullNameMessage(m.containing) + "." + relative
}

func (n *namer) relativeNameEnum(e *enum) string {
	if e.containing != nil {
		return sanitizeEnumName(e.name, n.forbidden())
	}
	return sanitizeEnumName(n.filePrefix(e.file)+e.name, n.forbidden())
}

func (n *namer) fullNameEnum(e *enum) string {
	relative := n.relativeNameEnum(e)
	if e.containing == nil {
		return n.modulePrefix(e.file) + relative
	}
	return n.fullNameMessage(e.containing) + "." + relative
}

func (n *namer) computeRelativeNames(e *enum) {
	stripper := newPrefixStripper(e.name)
	candidateName := func(v *enumValue) string {
		if stripped, ok := stripper.strip(v.name); ok {
			return toLowerCamelCase(stripped)
		}
		return toLowerCamelCase(v.name)
	}
	// Bucket the values by candidate name, in value order.
	var order []string
	candidates := map[string][]*enumValue{}
	for _, v := range e.values {
		name := candidateName(v)
		if _, ok := candidates[name]; !ok {
			order = append(order, name)
		}
		candidates[name] = append(candidates[name], v)
	}
	for _, camelCased := range order {
		values := candidates[camelCased]
		if len(values) == 1 {
			n.enumValueNameCache[values[0].fullName] = sanitizeEnumCaseName(camelCased)
			continue
		}
		first := values[0].number
		hasMultipleValues := false
		for _, v := range values {
			if v.number != first {
				hasMultipleValues = true
				break
			}
		}
		if !hasMultipleValues {
			name := sanitizeEnumCaseName(camelCased)
			for _, v := range values {
				n.enumValueNameCache[v.fullName] = name
			}
			continue
		}
		for _, v := range values {
			suffix := "_" + strconv.Itoa(int(v.number))
			if v.number < 0 {
				suffix = "_n" + strconv.Itoa(int(-int64(v.number)))
			}
			n.enumValueNameCache[v.fullName] = sanitizeEnumCaseName(camelCased + suffix)
		}
	}
}

func (n *namer) relativeNameEnumValue(v *enumValue) string {
	if name, ok := n.enumValueNameCache[v.fullName]; ok {
		return name
	}
	n.computeRelativeNames(v.enum)
	return n.enumValueNameCache[v.fullName]
}

// dottedRelativeName is the name of an enum case with a leading dot.
func (n *namer) dottedRelativeNameEnumValue(v *enumValue) string {
	return "." + trimBackticks(n.relativeNameEnumValue(v))
}

func (n *namer) relativeNameOneof(o *oneof) string {
	return sanitizeOneofName("OneOf_"+toUpperCamelCase(o.name), n.forbidden())
}

func (n *namer) fullNameOneof(o *oneof) string {
	return n.fullNameMessage(o.containing) + "." + n.relativeNameOneof(o)
}

func (n *namer) relativeNameExtension(f *field) string {
	if f.extScope != nil {
		return sanitizeMessageScopedExtensionName(f.namingBase())
	}
	return n.filePrefix(f.file) + "Extensions_" + f.namingBase()
}

func (n *namer) fullNameExtension(f *field) string {
	relative := n.relativeNameExtension(f)
	if f.extScope == nil {
		return n.modulePrefix(f.file) + relative
	}
	return n.fullNameMessage(f.extScope) + ".Extensions." + trimBackticks(relative)
}

type messageFieldNames struct {
	name, prefixed, has, clear string
}

// messagePropertyNames returns the names of the property of a field. The has
// and clear names are empty unless includeHasAndClear is set.
func (n *namer) messagePropertyNames(f *field, prefixed string, includeHasAndClear bool) messageFieldNames {
	lowerName := toLowerCamelCase(f.namingBase())
	fieldName := sanitizeFieldName(lowerName)
	prefixedName := ""
	if prefixed != "" {
		prefixedName = sanitizeFieldNameBasedOn(prefixed+lowerName, lowerName)
	}
	if !includeHasAndClear {
		return messageFieldNames{name: fieldName, prefixed: prefixedName}
	}
	upperName := toUpperCamelCase(f.namingBase())
	return messageFieldNames{
		name:     fieldName,
		prefixed: prefixedName,
		has:      sanitizeFieldNameBasedOn("has"+upperName, lowerName),
		clear:    sanitizeFieldNameBasedOn("clear"+upperName, lowerName),
	}
}

type oneofFieldNames struct{ name, prefixed string }

func (n *namer) messagePropertyNameOneof(o *oneof, prefixed string) oneofFieldNames {
	lowerName := toLowerCamelCase(o.name)
	return oneofFieldNames{
		name:     sanitizeFieldName(lowerName),
		prefixed: sanitizeFieldNameBasedOn(prefixed+lowerName, lowerName),
	}
}

type extensionNames struct{ value, has, clear string }

func (n *namer) messagePropertyNamesExtension(f *field) extensionNames {
	fieldBaseName := toLowerCamelCase(f.namingBase())
	if f.extScope != nil {
		scopeName := n.fullNameMessage(f.extScope)
		fieldName := strings.ReplaceAll(scopeName+"_"+fieldBaseName, ".", "_")
		firstUp := uppercaseFirstCharacter(fieldName)
		return extensionNames{value: fieldName, has: "has" + firstUp, clear: "clear" + firstUp}
	}
	prefix := n.filePrefix(f.file)
	fieldName := sanitizeFieldName(prefix + fieldBaseName)
	if prefix == "" {
		upperCleaned := sanitizeFieldNameBasedOn(toUpperCamelCase(f.namingBase()), fieldBaseName)
		return extensionNames{value: fieldName, has: "has" + upperCleaned, clear: "clear" + upperCleaned}
	}
	firstUp := uppercaseFirstCharacter(fieldName)
	return extensionNames{value: fieldName, has: "has" + firstUp, clear: "clear" + firstUp}
}

// uniquelyNamedValues drops the aliases that get the Swift name of their original value.
func (n *namer) uniquelyNamedValues(info *aliasInfo, e *enum) []*enumValue {
	var out []*enumValue
	for _, v := range e.values {
		original := info.original(v)
		if original == nil {
			out = append(out, v)
			continue
		}
		relative := n.relativeNameEnumValue(v)
		if relative == n.relativeNameEnumValue(original) {
			continue
		}
		aliases := info.aliases(original)
		firstIndex := -1
		for i, a := range aliases {
			if n.relativeNameEnumValue(a) == relative {
				firstIndex = i
				break
			}
		}
		if aliases[firstIndex] == v {
			out = append(out, v)
		}
	}
	return out
}

// aliasInfo records which enum values share a number.
type aliasInfo struct {
	mainValues []*enumValue
	aliasesOf  map[int][]*enumValue
	originals  map[int]*enumValue
}

func newAliasInfo(e *enum) *aliasInfo {
	info := &aliasInfo{aliasesOf: map[int][]*enumValue{}, originals: map[int]*enumValue{}}
	firstValues := map[int32]*enumValue{}
	for _, v := range e.values {
		if first, ok := firstValues[v.number]; ok {
			info.aliasesOf[first.index] = append(info.aliasesOf[first.index], v)
			info.originals[v.index] = first
		} else {
			firstValues[v.number] = v
			info.mainValues = append(info.mainValues, v)
		}
	}
	return info
}

func (a *aliasInfo) aliases(v *enumValue) []*enumValue { return a.aliasesOf[v.index] }

func (a *aliasInfo) original(v *enumValue) *enumValue { return a.originals[v.index] }
