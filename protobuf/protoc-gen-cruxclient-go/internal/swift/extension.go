// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"sort"
)

// extensionSet generates the code of the extension fields of a file.
type extensionSet struct {
	g          *generator
	f          *file
	extensions []*extensionGen
}

type extensionGen struct {
	g                           *generator
	desc                        *field
	comments                    string
	containingTypeSwiftFullName string
	swiftFullExtensionName      string
}

func (g *generator) newExtensionSet(f *file) *extensionSet {
	return &extensionSet{g: g, f: f}
}

func (s *extensionSet) isEmpty() bool { return len(s.extensions) == 0 }

func (s *extensionSet) add(fields []*field) {
	for _, f := range fields {
		s.extensions = append(s.extensions, &extensionGen{
			g: s.g, desc: f,
			comments:                    s.g.fieldComments(f),
			containingTypeSwiftFullName: s.g.namer.fullNameMessage(f.containing),
			swiftFullExtensionName:      s.g.namer.fullNameExtension(f),
		})
	}
}

func (e *extensionGen) extensionFieldType() string {
	f := e.desc
	var label string
	switch {
	case f.isRequired():
		label = "Required"
	case f.isRepeated():
		if f.isPacked() {
			label = "Packed"
		} else {
			label = "Repeated"
		}
	default:
		label = "Optional"
	}
	modifier := ""
	switch f.typ {
	case typeGroup:
		modifier = "Group"
	case typeMessage:
		modifier = "Message"
	case typeEnum:
		modifier = "Enum"
	}
	return e.g.namer.runtimePrefix() + label + modifier + "ExtensionField"
}

func (e *extensionGen) generateProtobufExtensionDeclarations(p *printer) {
	g := e.g
	f := e.desc
	visibility := g.opts.visibilitySnippet()
	isFileScope := f.extScope == nil
	scope, nonisolated := "static ", ""
	if isFileScope {
		scope, nonisolated = "", "nonisolated "
	}
	traitsType := g.traitsType(f)
	relativeName := g.namer.relativeNameExtension(f)
	fieldNamePath := f.fullName
	if f.containing.useMessageSetWireFormat() && f.typ == typeMessage && !f.isRepeated() && !f.isRequired() && f.messageType == f.extScope {
		fieldNamePath = f.messageType.fullName
	}
	p.print(fmt.Sprintf("%s%s%s%slet %s = %sMessageExtension<%s<%s>, %s>(",
		e.comments, visibility, nonisolated, scope, relativeName, g.namer.runtimePrefix(), e.extensionFieldType(), traitsType, e.containingTypeSwiftFullName))
	p.printIndented(
		fmt.Sprintf("_protobuf_fieldNumber: %d,", f.number),
		fmt.Sprintf("fieldName: \"%s\"", fieldNamePath),
	)
	p.print(")")
}

func (e *extensionGen) generateMessageSwiftExtension(p *printer) {
	g := e.g
	f := e.desc
	visibility := g.opts.visibilitySnippet()
	apiType := g.swiftType(f)
	names := g.namer.messagePropertyNamesExtension(f)
	defaultValue := g.swiftDefaultValue(f)
	p.print("", fmt.Sprintf("%s%svar %s: %s {", e.comments, visibility, names.value, apiType))
	p.printIndented(
		fmt.Sprintf("get {return getExtensionValue(ext: %s) ?? %s}", e.swiftFullExtensionName, defaultValue),
		fmt.Sprintf("set {setExtensionValue(ext: %s, value: newValue)}", e.swiftFullExtensionName),
	)
	p.print("}")
	if !f.isRepeated() {
		p.print(
			fmt.Sprintf("/// Returns true if extension `%s`\n/// has been explicitly set.", e.swiftFullExtensionName),
			fmt.Sprintf("%svar %s: Bool {", visibility, names.has),
		)
		p.printIndented(fmt.Sprintf("return hasExtensionValue(ext: %s)", e.swiftFullExtensionName))
		p.print("}")
		p.print(
			fmt.Sprintf("/// Clears the value of extension `%s`.", e.swiftFullExtensionName),
			"/// Subsequent reads from it will return its default value.",
			fmt.Sprintf("%smutating func %s() {", visibility, names.clear),
		)
		p.printIndented(fmt.Sprintf("clearExtensionValue(ext: %s)", e.swiftFullExtensionName))
		p.print("}")
	}
}

func (s *extensionSet) generateMessageSwiftExtensions(p *printer) {
	if s.isEmpty() {
		return
	}
	p.print("\n// MARK: - Extension Properties\n\n// Swift Extensions on the extended Messages to add easy access to the declared\n// extension fields. The names are based on the extension field name from the proto\n// declaration. To avoid naming collisions, the names are prefixed with the name of\n// the scope where the extend directive occurs.")

	grouped := append([]*extensionGen(nil), s.extensions...)
	index := map[*extensionGen]int{}
	for i, e := range s.extensions {
		index[e] = i
	}
	sort.SliceStable(grouped, func(i, j int) bool {
		a, b := grouped[i], grouped[j]
		if a.containingTypeSwiftFullName == b.containingTypeSwiftFullName {
			return index[a] < index[b]
		}
		return a.containingTypeSwiftFullName < b.containingTypeSwiftFullName
	})
	currentType := ""
	for _, e := range grouped {
		if currentType != e.containingTypeSwiftFullName {
			if currentType != "" {
				p.outdent()
				p.print("}")
			}
			currentType = e.containingTypeSwiftFullName
			p.print("", fmt.Sprintf("nonisolated extension %s {", currentType))
			p.indent()
		}
		e.generateMessageSwiftExtension(p)
	}
	p.outdent()
	p.print("", "}")
}

func (s *extensionSet) generateFileProtobufExtensionRegistry(p *printer) {
	if s.isEmpty() {
		return
	}
	g := s.g
	_, base, _ := splitPath(s.f.name)
	filenameAsIdentifier := toUpperCamelCase(base)
	filePrefix := g.namer.filePrefix(s.f)
	rt := g.namer.runtimePrefix()
	p.print(fmt.Sprintf("\n// MARK: - File's ExtensionMap: %s%s_Extensions\n\n"+
		"/// A `SwiftProtobuf.SimpleExtensionMap` that includes all of the extensions defined by\n"+
		"/// this .proto file. It can be used any place an `SwiftProtobuf.ExtensionMap` is needed\n"+
		"/// in parsing, or it can be combined with other `SwiftProtobuf.SimpleExtensionMap`s to create\n"+
		"/// a larger `SwiftProtobuf.SimpleExtensionMap`.\n"+
		"%snonisolated let %s%s_Extensions: %sSimpleExtensionMap = [",
		filePrefix, filenameAsIdentifier,
		g.opts.visibilitySnippet(), filePrefix, filenameAsIdentifier, rt))
	p.withIndentation(func() {
		last := len(s.extensions) - 1
		for i, e := range s.extensions {
			comma := ","
			if i == last {
				comma = ""
			}
			p.print(e.swiftFullExtensionName + comma)
		}
	})
	p.print("]")
}

func (s *extensionSet) generateProtobufExtensionDeclarations(p *printer) {
	if s.isEmpty() {
		return
	}
	p.print("\n// Extension Objects - The only reason these might be needed is when manually\n// constructing a `SimpleExtensionMap`, otherwise, use the above _Extension Properties_\n// accessors for the extension fields on the messages directly.")
	endScope := func() {
		p.outdent()
		p.print("}")
		p.outdent()
		p.print("}")
	}
	visibility := s.g.opts.visibilitySnippet()
	var currentScope *message
	addNewline := true
	for _, e := range s.extensions {
		if currentScope != e.desc.extScope {
			if currentScope != nil {
				endScope()
			}
			currentScope = e.desc.extScope
			// The extensions of the file come first, and they have no scope message.
			if currentScope != nil {
				p.print("", fmt.Sprintf("nonisolated extension %s {", s.g.namer.fullNameMessage(currentScope)))
				p.indent()
				p.print(visibility + "enum Extensions {")
				p.indent()
				addNewline = false
			}
		}
		if addNewline {
			p.print()
		} else {
			addNewline = true
		}
		e.generateProtobufExtensionDeclarations(p)
	}
	if currentScope != nil {
		endScope()
	}
}
