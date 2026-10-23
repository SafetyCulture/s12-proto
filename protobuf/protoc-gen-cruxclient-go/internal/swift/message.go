// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"sort"
	"strings"
)

// messageGen generates the struct of a message and its runtime support.
type messageGen struct {
	g              *generator
	desc           *message
	visibility     string
	swiftFullName  string
	swiftRelative  string
	fields         []fieldGenerator
	fieldsByNumber []fieldGenerator
	oneofs         []*oneofGen
	storage        *storageGen
	enums          []*enumGen
	messages       []*messageGen
	isExtensible   bool
}

func (g *generator) newMessageGen(m *message) *messageGen {
	mg := &messageGen{
		g: g, desc: m, visibility: g.opts.visibilitySnippet(),
		isExtensible:  len(m.extensionRange) > 0,
		swiftRelative: g.namer.relativeNameMessage(m),
		swiftFullName: g.namer.fullNameMessage(m),
	}
	useHeapStorage := g.storage.shouldUseHeapStorage(m)
	for _, o := range m.realOneofs() {
		mg.oneofs = append(mg.oneofs, g.newOneofGen(o, useHeapStorage))
	}
	for _, f := range m.fields {
		if o := f.realContainingOneof(); o != nil {
			mg.fields = append(mg.fields, mg.oneofs[o.index].fieldGenerator(f.number))
			continue
		}
		mg.fields = append(mg.fields, g.newMessageFieldGen(f, useHeapStorage))
	}
	mg.fieldsByNumber = append([]fieldGenerator(nil), mg.fields...)
	sort.SliceStable(mg.fieldsByNumber, func(i, j int) bool {
		return mg.fieldsByNumber[i].fieldNumber() < mg.fieldsByNumber[j].fieldNumber()
	})

	g.extensions.add(m.extensions)

	for _, e := range m.enums {
		mg.enums = append(mg.enums, g.newEnumGen(e))
	}
	for _, n := range m.nested {
		if n.proto.GetOptions().GetMapEntry() {
			continue
		}
		mg.messages = append(mg.messages, g.newMessageGen(n))
	}
	switch {
	case m.isWellKnownAny():
		if !useHeapStorage {
			panic("google.protobuf.Any must use heap storage")
		}
		mg.storage = &storageGen{fields: mg.fields, isAny: true}
	case useHeapStorage:
		mg.storage = &storageGen{fields: mg.fields}
	}
	return mg
}

func (mg *messageGen) generateMainStruct(p *printer) error {
	g := mg.g
	m := mg.desc
	if m.useMessageSetWireFormat() && len(mg.fields) > 0 {
		return fmt.Errorf("%s has the option message_set_wire_format but it also has fields", m.fullName)
	}
	for _, e := range m.extensions {
		if !e.containing.useMessageSetWireFormat() {
			continue
		}
		if e.typ != typeMessage {
			return fmt.Errorf("%s has the option message_set_wire_format but %s is a non message extension field", e.containing.fullName, e.fullName)
		}
		if e.isRequired() || e.isRepeated() {
			return fmt.Errorf("%s has the option message_set_wire_format but %s cannot be required nor repeated extension field", e.containing.fullName, e.fullName)
		}
	}

	rt := g.namer.runtimePrefix()
	var conformances []string
	if mg.isExtensible {
		conformances = append(conformances, rt+"ExtensibleMessage")
	}
	if mg.storage != nil {
		// Storage classes use copy-on-write semantics, so this Sendable conformance is safe.
		conformances = append(conformances, "@unchecked Sendable")
	} else {
		conformances = append(conformances, "Sendable")
	}
	p.print("", fmt.Sprintf("%s%snonisolated struct %s: %s {", g.messageComments(m), mg.visibility, mg.swiftRelative, strings.Join(conformances, ", ")))
	var err error
	p.withIndentation(func() {
		p.print(fmt.Sprintf("// %s.Message conformance is added in an extension below. See the\n// `Message` and `Message+*Additions` files in the SwiftProtobuf library for\n// methods supported on all messages.", g.namer.runtimeModule))
		for _, f := range mg.fields {
			f.generateInterface(p)
		}
		p.print("", fmt.Sprintf("%svar unknownFields = %sUnknownStorage()", mg.visibility, rt))
		for _, o := range mg.oneofs {
			o.generateMainEnum(p)
		}
		for _, e := range mg.enums {
			e.generateMainEnum(p)
		}
		for _, n := range mg.messages {
			if err == nil {
				err = n.generateMainStruct(p)
			}
		}
		p.print("", mg.visibility+"init() {}")
		if mg.isExtensible {
			p.print("", fmt.Sprintf("%svar _protobuf_extensionFieldValues = %sExtensionFieldValueSet()", mg.visibility, rt))
		}
		if mg.storage != nil {
			if !mg.isExtensible {
				p.print()
			}
			p.print(mg.storage.storageVisibility() + " var _storage = _StorageClass.defaultInstance")
		} else {
			sub := p.sub()
			for _, f := range mg.fields {
				f.generateStorage(sub)
			}
			if !sub.isEmpty() {
				if !mg.isExtensible {
					p.print()
				}
				p.appendPrinter(sub, false)
			}
		}
	})
	p.print("}")
	return err
}

func (mg *messageGen) generateRuntimeSupport(p *printer, parent *messageGen) {
	g := mg.g
	rt := g.namer.runtimePrefix()
	m := mg.desc
	p.print("", fmt.Sprintf("nonisolated extension %s: %sMessage, %s_MessageImplementationBase, %s_ProtoNameProviding {", mg.swiftFullName, rt, rt, rt))
	p.withIndentation(func() {
		switch {
		case g.opts.hides(hiddenTypes):
			p.print(mg.visibility + `static let protoMessageName: String = ""`)
		case parent != nil:
			p.print(fmt.Sprintf(`%sstatic let protoMessageName: String = %s.protoMessageName + ".%s"`, mg.visibility, parent.swiftFullName, m.name))
		case m.file.pkg != "":
			p.print(fmt.Sprintf(`%sstatic let protoMessageName: String = _protobuf_package + ".%s"`, mg.visibility, m.name))
		default:
			p.print(fmt.Sprintf(`%sstatic let protoMessageName: String = "%s"`, mg.visibility, m.name))
		}
		mg.generateProtoNameProviding(p)
		if mg.storage != nil {
			p.print()
			mg.storage.generateTypeDeclaration(p)
			p.print()
			mg.storage.generateUniqueStorage(p)
		}
		p.print()
		mg.generateIsInitialized(p)
		mg.generateDecodeMessage(p)
		p.print()
		mg.generateTraverse(p)
		p.print()
		mg.generateMessageEquality(p)
	})
	p.print("}")
	for _, e := range mg.enums {
		e.generateRuntimeSupport(p)
	}
	for _, n := range mg.messages {
		n.generateRuntimeSupport(p, mg)
	}
}

func (mg *messageGen) generateProtoNameProviding(p *printer) {
	rt := mg.g.namer.runtimePrefix()
	if mg.g.opts.hides(hiddenFields) {
		p.print(fmt.Sprintf("%sstatic let _protobuf_nameMap = %s_NameMap()", mg.visibility, rt))
		return
	}
	w := newNameInstructionWriter()
	for _, f := range mg.fieldsByNumber {
		f.writeProtoNameInstruction(w)
	}
	for _, name := range mg.desc.reservedNames {
		w.writeReservedName(name)
	}
	for _, r := range mg.desc.reservedRanges {
		w.writeReservedNumbers(r[0], r[1])
	}
	if w.err != nil {
		mg.g.err = fmt.Errorf("%s: %w", mg.desc.fullName, w.err)
		return
	}
	if w.shouldUseEmptyNameMapInitializer() {
		p.print(fmt.Sprintf("%sstatic let _protobuf_nameMap = %s_NameMap()", mg.visibility, rt))
		return
	}
	p.print(fmt.Sprintf("%sstatic let _protobuf_nameMap = %s_NameMap(bytecode: %s)", mg.visibility, rt, w.bytecode.stringLiteral()))
}

func rangeCaseExpression(r [2]int32) string {
	if r[0] == r[1]-1 {
		return fmt.Sprintf("%d", r[0])
	}
	return fmt.Sprintf("%d..<%d", r[0], r[1])
}

func rangeBooleanExpression(r [2]int32, variable string) string {
	if r[0] == r[1]-1 {
		return fmt.Sprintf("%d == %s", r[0], variable)
	}
	return fmt.Sprintf("%d <= %s && %s < %d", r[0], variable, variable, r[1])
}

func rangesCaseExpression(ranges [][2]int32) string {
	parts := make([]string, len(ranges))
	for i, r := range ranges {
		parts[i] = rangeCaseExpression(r)
	}
	return strings.Join(parts, ", ")
}

func rangesBooleanExpression(ranges [][2]int32, variable string) string {
	parts := make([]string, len(ranges))
	for i, r := range ranges {
		parts[i] = "(" + rangeBooleanExpression(r, variable) + ")"
	}
	return strings.Join(parts, " || ")
}

func (mg *messageGen) generateDecodeMessage(p *printer) {
	rt := mg.g.namer.runtimePrefix()
	m := mg.desc
	p.print(fmt.Sprintf("%smutating func decodeMessage<D: %sDecoder>(decoder: inout D) throws {", mg.visibility, rt))
	p.withIndentation(func() {
		if mg.storage != nil {
			p.print("_ = _uniqueStorage()")
		}
		if m.useMessageSetWireFormat() && mg.isExtensible {
			p.print(fmt.Sprintf("try decoder.decodeExtensionFieldsAsMessageSet(values: &_protobuf_extensionFieldValues, messageType: %s.self)", mg.swiftFullName))
			return
		}
		if len(mg.fields) == 0 && !mg.isExtensible {
			p.print("// Load everything into unknown fields", "while try decoder.nextFieldNumber() != nil {}")
			return
		}
		mg.withLifetimeExtension(p, true, false, "", "", func() {
			p.print("while let fieldNumber = try decoder.nextFieldNumber() {")
			p.withIndentation(func() {
				ranges := m.normalizedExtensionRanges()
				switch {
				case len(mg.fields) > 0 || len(ranges) > 3:
					p.print("// The use of inline closures is to circumvent an issue where the compiler\n// allocates stack space for every case branch when no optimizations are\n// enabled. https://github.com/apple/swift-protobuf/issues/1034\nswitch fieldNumber {")
					for _, f := range mg.fieldsByNumber {
						f.generateDecodeFieldCase(p)
					}
					if mg.isExtensible {
						p.print(fmt.Sprintf("case %s:", rangesCaseExpression(ranges)))
						p.printIndented(fmt.Sprintf("try { try decoder.decodeExtensionField(values: &_protobuf_extensionFieldValues, messageType: %s.self, fieldNumber: fieldNumber) }()", mg.swiftFullName))
					}
					p.print("default: break", "}")
				case mg.isExtensible:
					p.print(fmt.Sprintf("if %s {", rangesBooleanExpression(ranges, "fieldNumber")))
					p.printIndented(fmt.Sprintf("try decoder.decodeExtensionField(values: &_protobuf_extensionFieldValues, messageType: %s.self, fieldNumber: fieldNumber)", mg.swiftFullName))
					p.print("}")
				}
			})
			p.print("}")
		})
	})
	p.print("}")
}

func (mg *messageGen) generateTraverse(p *printer) {
	rt := mg.g.namer.runtimePrefix()
	m := mg.desc
	p.print(fmt.Sprintf("%sfunc traverse<V: %sVisitor>(visitor: inout V) throws {", mg.visibility, rt))
	p.withIndentation(func() {
		mg.withLifetimeExtension(p, true, false, "", "", func() {
			if mg.storage != nil {
				mg.storage.generatePreTraverse(p)
			}
			visitExtensionsName := "visitExtensionFields"
			if m.useMessageSetWireFormat() {
				visitExtensionsName = "visitExtensionFieldsAsMessageSet"
			}
			usesLocals := false
			for _, f := range mg.fields {
				usesLocals = usesLocals || f.generateTraverseUsesLocals()
			}
			if usesLocals {
				p.print("// The use of inline closures is to circumvent an issue where the compiler\n// allocates stack space for every if/case branch local when no optimizations\n// are enabled. https://github.com/apple/swift-protobuf/issues/1034 and\n// https://github.com/apple/swift-protobuf/issues/1182")
			}
			ranges := m.ambitiousExtensionRanges()
			next := 0
			for _, f := range mg.fieldsByNumber {
				for next < len(ranges) && int(ranges[next][0]) < int(f.fieldNumber()) {
					p.print(fmt.Sprintf("try visitor.%s(fields: _protobuf_extensionFieldValues, start: %d, end: %d)", visitExtensionsName, ranges[next][0], ranges[next][1]))
					next++
				}
				f.generateTraverse(p)
			}
			for next < len(ranges) {
				p.print(fmt.Sprintf("try visitor.%s(fields: _protobuf_extensionFieldValues, start: %d, end: %d)", visitExtensionsName, ranges[next][0], ranges[next][1]))
				next++
			}
		})
		p.print("try unknownFields.traverse(visitor: &visitor)")
	})
	p.print("}")
}

func (mg *messageGen) generateMessageEquality(p *printer) {
	name := mg.swiftFullName
	p.print(fmt.Sprintf("%sstatic func ==(lhs: %s, rhs: %s) -> Bool {", mg.visibility, name, name))
	p.withIndentation(func() {
		compareFields := true
		if mg.storage != nil {
			p.print("if lhs._storage !== rhs._storage {")
			p.indent()
			p.printNoNewline("let storagesAreEqual: Bool = ")
			if mg.storage.providesEqualTo() {
				p.print("lhs._storage.isEqualTo(other: rhs._storage)")
				compareFields = false
			}
		}
		if compareFields {
			mg.withLifetimeExtension(p, false, false, "rhs", "lhs", func() {
				for _, f := range mg.fields {
					f.generateFieldComparison(p)
				}
				if mg.storage != nil {
					p.print("return true")
				}
			})
		}
		if mg.storage != nil {
			p.print("if !storagesAreEqual {return false}")
			p.outdent()
			p.print("}")
		}
		p.print("if lhs.unknownFields != rhs.unknownFields {return false}")
		if mg.isExtensible {
			p.print("if lhs._protobuf_extensionFieldValues != rhs._protobuf_extensionFieldValues {return false}")
		}
		p.print("return true")
	})
	p.print("}")
}

func (mg *messageGen) generateIsInitialized(p *printer) {
	fieldCheck := newPrinter()
	for _, f := range mg.fields {
		f.generateRequiredFieldCheck(fieldCheck)
	}
	for _, f := range mg.fields {
		f.generateIsInitializedCheck(fieldCheck)
	}
	generatedChecks := !fieldCheck.isEmpty()
	if !mg.isExtensible && !generatedChecks {
		return
	}
	p.print("public var isInitialized: Bool {")
	p.withIndentation(func() {
		if mg.isExtensible {
			p.print("if !_protobuf_extensionFieldValues.isInitialized {return false}")
		}
		if generatedChecks {
			mg.withLifetimeExtension(p, false, true, "", "", func() {
				p.appendPrinter(fieldCheck, true)
				p.print("return true")
			})
		} else {
			p.print("return true")
		}
	})
	p.print("}", "")
}

// withLifetimeExtension runs the body inside a withExtendedLifetime call when the message has a storage class.
func (mg *messageGen) withLifetimeExtension(p *printer, canThrow, returns bool, capturedVariable, qualifier string, body func()) {
	if mg.storage != nil {
		prefixKeywords := ""
		if returns {
			prefixKeywords += "return "
		}
		if canThrow {
			prefixKeywords += "try "
		}
		selfQualifier := ""
		if qualifier != "" {
			selfQualifier = qualifier + "."
		}
		if capturedVariable != "" {
			actualArgs := fmt.Sprintf("(%s_storage, %s._storage)", selfQualifier, capturedVariable)
			p.print(fmt.Sprintf("%swithExtendedLifetime(%s) { (_args: (_StorageClass, _StorageClass)) in", prefixKeywords, actualArgs))
			p.indent()
			p.print("let _storage = _args.0")
			p.print(fmt.Sprintf("let %s_storage = _args.1", capturedVariable))
		} else {
			p.print(fmt.Sprintf("%swithExtendedLifetime(%s_storage) { (_storage: _StorageClass) in", prefixKeywords, selfQualifier))
			p.indent()
		}
	}
	body()
	if mg.storage != nil {
		p.outdent()
		p.print("}")
	}
}
