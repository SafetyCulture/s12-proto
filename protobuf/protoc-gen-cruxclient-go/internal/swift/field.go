// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import "fmt"

// fieldGenerator generates the code of one field of a message.
type fieldGenerator interface {
	fieldNumber() int32
	writeProtoNameInstruction(w *nameInstructionWriter)
	generateInterface(p *printer)
	generateStorage(p *printer)
	generateStorageClassClone(p *printer)
	generateDecodeFieldCase(p *printer)
	generateTraverseUsesLocals() bool
	generateTraverse(p *printer)
	generateFieldComparison(p *printer)
	generateRequiredFieldCheck(p *printer)
	generateIsInitializedCheck(p *printer)
}

// writeFieldNameInstruction writes the name of a field to the name map program.
func writeFieldNameInstruction(f *field, w *nameInstructionWriter) {
	protoName := f.name
	if f.isGroupLike() {
		protoName = f.messageType.name
	}
	jsonName := f.jsonName
	switch {
	case f.isGroupLike():
		w.writeGroup(f.number, protoName)
	case jsonName == protoName:
		w.writeSame(f.number, protoName)
	default:
		if jsonName == toJSONFieldName(protoName) {
			w.writeStandard(f.number, protoName)
		} else {
			w.writeUnique(f.number, protoName, jsonName)
		}
	}
}

func isGroupOrMessage(f *field) bool {
	return f.typ == typeGroup || f.typ == typeMessage
}

// messageFieldGen generates a field that is not in a oneof.
type messageFieldGen struct {
	g                 *generator
	desc              *field
	usesHeapStorage   bool
	hasFieldPresence  bool
	swiftName         string
	underscoreName    string
	storedProperty    string
	hasName           string
	clearName         string
	swiftType         string
	swiftStorageType  string
	swiftDefaultValue string
	traitsType        string
	comments          string
}

func (g *generator) newMessageFieldGen(f *field, usesHeapStorage bool) *messageFieldGen {
	mf := &messageFieldGen{g: g, desc: f, usesHeapStorage: usesHeapStorage, hasFieldPresence: f.hasPresence()}
	names := g.namer.messagePropertyNames(f, "_", mf.hasFieldPresence)
	mf.swiftName = names.name
	mf.underscoreName = names.prefixed
	mf.hasName = names.has
	mf.clearName = names.clear
	mf.swiftType = g.swiftType(f)
	mf.swiftStorageType = g.swiftStorageType(f)
	mf.swiftDefaultValue = g.swiftDefaultValue(f)
	mf.traitsType = g.traitsType(f)
	mf.comments = g.fieldComments(f)
	switch {
	case usesHeapStorage:
		mf.storedProperty = "_storage." + mf.underscoreName
	case mf.hasFieldPresence:
		mf.storedProperty = "self." + mf.underscoreName
	default:
		mf.storedProperty = "self." + mf.swiftName
	}
	return mf
}

func (mf *messageFieldGen) fieldNumber() int32 { return mf.desc.number }

func (mf *messageFieldGen) writeProtoNameInstruction(w *nameInstructionWriter) {
	writeFieldNameInstruction(mf.desc, w)
}

func (mf *messageFieldGen) generateStorage(p *printer) {
	defaultValue := mf.swiftDefaultValue
	if mf.hasFieldPresence {
		defaultValue = "nil"
	}
	if mf.usesHeapStorage {
		p.print(fmt.Sprintf("var %s: %s = %s", mf.underscoreName, mf.swiftStorageType, defaultValue))
		return
	}
	if mf.hasFieldPresence {
		p.print(fmt.Sprintf("fileprivate var %s: %s = %s", mf.underscoreName, mf.swiftStorageType, defaultValue))
	}
}

func (mf *messageFieldGen) generateInterface(p *printer) {
	visibility := mf.g.opts.visibilitySnippet()
	p.print()
	if mf.usesHeapStorage {
		p.print(fmt.Sprintf("%s%svar %s: %s {", mf.comments, visibility, mf.swiftName, mf.swiftType))
		defaultClause := ""
		if mf.hasFieldPresence {
			defaultClause = " ?? " + mf.swiftDefaultValue
		}
		p.printIndented(
			fmt.Sprintf("get {_storage.%s%s}", mf.underscoreName, defaultClause),
			fmt.Sprintf("set {_uniqueStorage().%s = newValue}", mf.underscoreName),
		)
		p.print("}")
	} else if mf.hasFieldPresence {
		p.print(fmt.Sprintf("%s%svar %s: %s {", mf.comments, visibility, mf.swiftName, mf.swiftType))
		p.printIndented(
			fmt.Sprintf("get {%s ?? %s}", mf.underscoreName, mf.swiftDefaultValue),
			fmt.Sprintf("set {%s = newValue}", mf.underscoreName),
		)
		p.print("}")
	} else {
		p.print(fmt.Sprintf("%s%svar %s: %s = %s", mf.comments, visibility, mf.swiftName, mf.swiftStorageType, mf.swiftDefaultValue))
	}
	if !mf.hasFieldPresence {
		return
	}
	immutablePrefix := "self."
	mutablePrefix := "self."
	if mf.usesHeapStorage {
		immutablePrefix = "_storage."
		mutablePrefix = "_uniqueStorage()."
	}
	p.print(
		fmt.Sprintf("/// Returns true if `%s` has been explicitly set.", mf.swiftName),
		fmt.Sprintf("%svar %s: Bool {%s%s != nil}", visibility, mf.hasName, immutablePrefix, mf.underscoreName),
	)
	p.print(
		fmt.Sprintf("/// Clears the value of `%s`. Subsequent reads from it will return its default value.", mf.swiftName),
		fmt.Sprintf("%smutating func %s() {%s%s = nil}", visibility, mf.clearName, mutablePrefix, mf.underscoreName),
	)
}

func (mf *messageFieldGen) generateStorageClassClone(p *printer) {
	p.print(fmt.Sprintf("%s = source.%s", mf.underscoreName, mf.underscoreName))
}

func (mf *messageFieldGen) generateFieldComparison(p *printer) {
	var lhs, rhs string
	if mf.usesHeapStorage {
		lhs = "_storage." + mf.underscoreName
		rhs = "rhs_storage." + mf.underscoreName
	} else {
		name := mf.swiftName
		if mf.hasFieldPresence {
			name = mf.underscoreName
		}
		lhs = "lhs." + name
		rhs = "rhs." + name
	}
	p.print(fmt.Sprintf("if %s != %s {return false}", lhs, rhs))
}

func (mf *messageFieldGen) generateRequiredFieldCheck(p *printer) {
	if !mf.desc.isRequired() {
		return
	}
	p.print(fmt.Sprintf("if %s == nil {return false}", mf.storedProperty))
}

func (mf *messageFieldGen) generateIsInitializedCheck(p *printer) {
	if !isGroupOrMessage(mf.desc) || !mf.desc.messageType.containsRequiredFields() {
		return
	}
	if mf.desc.isRepeated() {
		p.print(fmt.Sprintf("if !%sInternal.areAllInitialized(%s) {return false}", mf.g.namer.runtimePrefix(), mf.storedProperty))
	} else {
		p.print(fmt.Sprintf("if let v = %s, !v.isInitialized {return false}", mf.storedProperty))
	}
}

func (mf *messageFieldGen) generateDecodeFieldCase(p *printer) {
	var decoderMethod, traitsArg string
	if mf.desc.isMap() {
		decoderMethod = "decodeMapField"
		traitsArg = "fieldType: " + mf.traitsType + ".self, "
	} else {
		modifier := "Singular"
		if mf.desc.isRepeated() {
			modifier = "Repeated"
		}
		decoderMethod = "decode" + modifier + mf.desc.protoGenericType() + "Field"
	}
	p.print(fmt.Sprintf("case %d: try { try decoder.%s(%svalue: &%s) }()", mf.desc.number, decoderMethod, traitsArg, mf.storedProperty))
}

func (mf *messageFieldGen) generateTraverseUsesLocals() bool {
	return !mf.desc.isRepeated() && mf.hasFieldPresence
}

func (mf *messageFieldGen) generateTraverse(p *printer) {
	f := mf.desc
	var visitMethod, traitsArg string
	if f.isMap() {
		visitMethod = "visitMapField"
		traitsArg = "fieldType: " + mf.traitsType + ".self, "
	} else {
		modifier := "Singular"
		switch {
		case f.isPacked():
			modifier = "Packed"
		case f.isRepeated():
			modifier = "Repeated"
		}
		visitMethod = "visit" + modifier + f.protoGenericType() + "Field"
	}
	varName := mf.storedProperty
	if mf.hasFieldPresence {
		varName = "v"
	}
	usesLocals := false
	var conditional string
	switch {
	case f.isRepeated():
		conditional = "!" + varName + ".isEmpty"
	case mf.hasFieldPresence:
		conditional = "let v = " + mf.storedProperty
		usesLocals = true
	default:
		switch f.typ {
		case typeString, typeBytes:
			conditional = "!" + varName + ".isEmpty"
		case typeFloat, typeDouble:
			conditional = varName + ".bitPattern != 0"
		default:
			conditional = varName + " != " + mf.swiftDefaultValue
		}
	}
	prefix, suffix := "", ""
	if usesLocals {
		prefix, suffix = "try { ", " }()"
	}
	p.print(fmt.Sprintf("%sif %s {", prefix, conditional))
	p.printIndented(fmt.Sprintf("try visitor.%s(%svalue: %s, fieldNumber: %d)", visitMethod, traitsArg, varName, f.number))
	p.print("}" + suffix)
}
