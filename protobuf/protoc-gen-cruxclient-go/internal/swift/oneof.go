// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"sort"
)

// oneofGen generates the enum of a oneof and the code of its fields.
type oneofGen struct {
	g               *generator
	desc            *oneof
	usesHeapStorage bool

	fields               []*oneofMemberGen
	fieldsSortedByNumber []*oneofMemberGen
	groups               [][]*oneofMemberGen
	relativeName         string
	fullName             string
	comments             string
	swiftFieldName       string
	underscoreFieldName  string
	storedProperty       string
}

// oneofMemberGen is the generator of one field of a oneof. It passes the calls to the oneof.
type oneofMemberGen struct {
	oneof             *oneofGen
	desc              *field
	group             int
	swiftName         string
	dottedSwiftName   string
	swiftType         string
	swiftDefaultValue string
	protoGenericType  string
	comments          string
}

func (g *generator) newOneofGen(o *oneof, usesHeapStorage bool) *oneofGen {
	og := &oneofGen{g: g, desc: o, usesHeapStorage: usesHeapStorage}
	og.comments = g.oneofComments(o)
	og.relativeName = g.namer.relativeNameOneof(o)
	og.fullName = g.namer.fullNameOneof(o)
	names := g.namer.messagePropertyNameOneof(o, "_")
	og.swiftFieldName = names.name
	og.underscoreFieldName = names.prefixed
	if usesHeapStorage {
		og.storedProperty = "_storage." + og.underscoreFieldName
	} else {
		og.storedProperty = "self." + og.swiftFieldName
	}
	for _, f := range o.fields {
		names := g.namer.messagePropertyNames(f, ".", false)
		og.fields = append(og.fields, &oneofMemberGen{
			oneof: og, desc: f, group: -1,
			swiftName:         names.name,
			dottedSwiftName:   names.prefixed,
			swiftType:         g.swiftType(f),
			swiftDefaultValue: g.swiftDefaultValue(f),
			protoGenericType:  f.protoGenericType(),
			comments:          g.fieldComments(f),
		})
	}
	og.fieldsSortedByNumber = append([]*oneofMemberGen(nil), og.fields...)
	sort.SliceStable(og.fieldsSortedByNumber, func(i, j int) bool {
		return og.fieldsSortedByNumber[i].desc.number < og.fieldsSortedByNumber[j].desc.number
	})

	// Bucket the fields in continuous chunks of the parent: its own fields and
	// the start of its extension ranges mark where a chunk ends.
	parent := o.containing
	var parentNumbers []int
	for _, f := range parent.fields {
		parentNumbers = append(parentNumbers, int(f.number))
	}
	for _, r := range parent.normalizedExtensionRanges() {
		parentNumbers = append(parentNumbers, int(r[0]))
	}
	sort.Ints(parentNumbers)
	pos := 0
	nextParent := func() (int, bool) {
		if pos >= len(parentNumbers) {
			return 0, false
		}
		n := parentNumbers[pos]
		pos++
		return n, true
	}
	nextParentNumber, hasNext := nextParent()
	var current []*oneofMemberGen
	for _, f := range og.fieldsSortedByNumber {
		number := int(f.desc.number)
		if !hasNext || nextParentNumber != number {
			if len(current) > 0 {
				og.groups = append(og.groups, current)
				current = nil
			}
			for !hasNext || nextParentNumber != number {
				if !hasNext {
					panic(fmt.Sprintf("oneof %s has no parent field %d", o.name, number))
				}
				nextParentNumber, hasNext = nextParent()
			}
		}
		current = append(current, f)
		nextParentNumber, hasNext = nextParent()
	}
	if len(current) > 0 {
		og.groups = append(og.groups, current)
	}
	for i, group := range og.groups {
		for _, f := range group {
			f.group = i
		}
	}
	return og
}

func (og *oneofGen) fieldGenerator(number int32) fieldGenerator {
	for _, f := range og.fields {
		if f.desc.number == number {
			return f
		}
	}
	panic("oneof field not found")
}

func (og *oneofGen) generateMainEnum(p *printer) {
	g := og.g
	visibility := g.opts.visibilitySnippet()
	enumDecl := visibility + "nonisolated enum " + og.relativeName + ": Equatable, Sendable {"
	p.print("", og.comments+g.opts.enumAnnotation()+enumDecl)
	p.withIndentation(func() {
		for _, f := range og.fields {
			p.print(fmt.Sprintf("%scase %s(%s)", f.comments, f.swiftName, f.swiftType))
		}
		var fieldsToCheck []*oneofMemberGen
		for _, f := range og.fields {
			if isGroupOrMessage(f.desc) && f.desc.messageType.containsRequiredFields() {
				fieldsToCheck = append(fieldsToCheck, f)
			}
		}
		if len(fieldsToCheck) > 0 {
			p.print("", "fileprivate var isInitialized: Bool {")
			p.withIndentation(func() {
				if len(fieldsToCheck) == 1 {
					f := fieldsToCheck[0]
					p.print(
						fmt.Sprintf("guard case %s(let v) = self else {return true}", f.dottedSwiftName),
						"return v.isInitialized",
					)
				} else {
					p.print("// The use of inline closures is to circumvent an issue where the compiler\n// allocates stack space for every case branch when no optimizations are\n// enabled. https://github.com/apple/swift-protobuf/issues/1034\nswitch self {")
					for _, f := range fieldsToCheck {
						p.print(fmt.Sprintf("case %s: return {", f.dottedSwiftName))
						p.printIndented(
							fmt.Sprintf("guard case %s(let v) = self else { preconditionFailure() }", f.dottedSwiftName),
							"return v.isInitialized",
						)
						p.print("}()")
					}
					if len(fieldsToCheck) != len(og.fields) {
						p.print("default: return true")
					}
					p.print("}")
				}
			})
			p.print("}")
		}
		p.print()
	})
	p.print("}")
}

func (og *oneofGen) generateOneofEnumProperty(p *printer) {
	visibility := og.g.opts.visibilitySnippet()
	p.print()
	if og.usesHeapStorage {
		p.print(fmt.Sprintf("%s%svar %s: %s? {", og.comments, visibility, og.swiftFieldName, og.relativeName))
		p.printIndented(
			fmt.Sprintf("get {return _storage.%s}", og.underscoreFieldName),
			fmt.Sprintf("set {_uniqueStorage().%s = newValue}", og.underscoreFieldName),
		)
		p.print("}")
	} else {
		p.print(fmt.Sprintf("%s%svar %s: %s? = nil", og.comments, visibility, og.swiftFieldName, og.fullName))
	}
}

func (f *oneofMemberGen) fieldNumber() int32 { return f.desc.number }

func (f *oneofMemberGen) writeProtoNameInstruction(w *nameInstructionWriter) {
	writeFieldNameInstruction(f.desc, w)
}

func (f *oneofMemberGen) generateInterface(p *printer) {
	og := f.oneof
	if f == og.fields[0] {
		og.generateOneofEnumProperty(p)
	}
	getter := og.swiftFieldName
	setter := og.swiftFieldName
	if og.usesHeapStorage {
		getter = "_storage." + og.underscoreFieldName
		setter = "_uniqueStorage()." + og.underscoreFieldName
	} else if og.swiftFieldName == "newValue" {
		setter = "self.newValue"
	}
	visibility := og.g.opts.visibilitySnippet()
	p.print("", fmt.Sprintf("%s%svar %s: %s {", f.comments, visibility, f.swiftName, f.swiftType))
	p.withIndentation(func() {
		p.print("get {")
		p.printIndented(
			fmt.Sprintf("if case %s(let v)? = %s {return v}", f.dottedSwiftName, getter),
			fmt.Sprintf("return %s", f.swiftDefaultValue),
		)
		p.print("}", fmt.Sprintf("set {%s = %s(newValue)}", setter, f.dottedSwiftName))
	})
	p.print("}")
}

func (f *oneofMemberGen) generateStorage(p *printer) {
	og := f.oneof
	if f != og.fields[0] || !og.usesHeapStorage {
		return
	}
	p.print(fmt.Sprintf("var %s: %s?", og.underscoreFieldName, og.fullName))
}

func (f *oneofMemberGen) generateStorageClassClone(p *printer) {
	og := f.oneof
	if f != og.fields[0] {
		return
	}
	p.print(fmt.Sprintf("%s = source.%s", og.underscoreFieldName, og.underscoreFieldName))
}

func (f *oneofMemberGen) generateDecodeFieldCase(p *printer) {
	og := f.oneof
	p.print(fmt.Sprintf("case %d: try {", f.desc.number))
	p.withIndentation(func() {
		var hadValueTest string
		if isGroupOrMessage(f.desc) {
			p.print(
				fmt.Sprintf("var v: %s?", f.swiftType),
				"var hadOneofValue = false",
				fmt.Sprintf("if let current = %s {", og.storedProperty),
			)
			p.printIndented(
				"hadOneofValue = true",
				fmt.Sprintf("if case %s(let m) = current {v = m}", f.dottedSwiftName),
			)
			p.print("}")
			hadValueTest = "hadOneofValue"
		} else {
			p.print(fmt.Sprintf("var v: %s?", f.swiftType))
			hadValueTest = og.storedProperty + " != nil"
		}
		p.print(
			fmt.Sprintf("try decoder.decodeSingular%sField(value: &v)", f.protoGenericType),
			"if let v = v {",
		)
		p.printIndented(
			fmt.Sprintf("if %s {try decoder.handleConflictingOneOf()}", hadValueTest),
			fmt.Sprintf("%s = %s(v)", og.storedProperty, f.dottedSwiftName),
		)
		p.print("}")
	})
	p.print("}()")
}

func (f *oneofMemberGen) generateTraverseUsesLocals() bool { return true }

func (f *oneofMemberGen) generateTraverse(p *printer) {
	og := f.oneof
	group := og.groups[f.group]
	if f != group[0] {
		return
	}
	if len(group) == 1 {
		p.print(fmt.Sprintf("try { if case %s(let v)? = %s {", f.dottedSwiftName, og.storedProperty))
		p.printIndented(fmt.Sprintf("try visitor.visitSingular%sField(value: v, fieldNumber: %d)", f.protoGenericType, f.desc.number))
		p.print("} }()")
		return
	}
	isOnlyGroup := len(og.groups) == 1
	p.print(fmt.Sprintf("switch %s {", og.storedProperty))
	for idx, m := range group {
		if isOnlyGroup && idx > 0 && idx%maxCasesInSwitch == 0 {
			p.print("default: break", "}")
			p.print(fmt.Sprintf("switch %s {", og.storedProperty))
		}
		p.print(fmt.Sprintf("case %s?: try {", m.dottedSwiftName))
		p.printIndented(
			fmt.Sprintf("guard case %s(let v)? = %s else { preconditionFailure() }", m.dottedSwiftName, og.storedProperty),
			fmt.Sprintf("try visitor.visitSingular%sField(value: v, fieldNumber: %d)", m.protoGenericType, m.desc.number),
		)
		p.print("}()")
	}
	if isOnlyGroup && len(group) <= maxCasesInSwitch {
		p.print("case nil: break")
	} else {
		p.print("default: break")
	}
	p.print("}")
}

func (f *oneofMemberGen) generateFieldComparison(p *printer) {
	og := f.oneof
	if f != og.fields[0] {
		return
	}
	var lhs, rhs string
	if og.usesHeapStorage {
		lhs = "_storage." + og.underscoreFieldName
		rhs = "rhs_storage." + og.underscoreFieldName
	} else {
		lhs = "lhs." + og.swiftFieldName
		rhs = "rhs." + og.swiftFieldName
	}
	p.print(fmt.Sprintf("if %s != %s {return false}", lhs, rhs))
}

func (f *oneofMemberGen) generateRequiredFieldCheck(*printer) {}

func (f *oneofMemberGen) generateIsInitializedCheck(p *printer) {
	og := f.oneof
	if f != og.fields[0] {
		return
	}
	for _, m := range og.fields {
		if isGroupOrMessage(m.desc) && m.desc.messageType.containsRequiredFields() {
			p.print(fmt.Sprintf("if let v = %s, !v.isInitialized {return false}", og.storedProperty))
			return
		}
	}
}
