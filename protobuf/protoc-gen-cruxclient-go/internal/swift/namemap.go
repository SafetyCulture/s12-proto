// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"strings"
)

// The opcodes of the name map program that the SwiftProtobuf runtime reads.
const (
	opSameNext        = 1
	opSameDelta       = 2
	opStandardNext    = 3
	opStandardDelta   = 4
	opUniqueNext      = 5
	opUniqueDelta     = 6
	opGroupNext       = 7
	opGroupDelta      = 8
	opAliasNext       = 9
	opAliasDelta      = 10
	opReservedName    = 11
	opReservedNumbers = 12

	latestBytecodeProgramFormat = 0
)

// bytecodeWriter writes a program as the content of a Swift string literal.
type bytecodeWriter struct {
	code    strings.Builder
	hasData bool
}

func newBytecodeWriter() *bytecodeWriter {
	w := &bytecodeWriter{}
	w.writeUInt64(latestBytecodeProgramFormat)
	w.hasData = false
	return w
}

// stringLiteral is the literal with its quotes.
func (w *bytecodeWriter) stringLiteral() string {
	return `"` + w.code.String() + `"`
}

func (w *bytecodeWriter) writeOpcode(op uint64) { w.writeUInt64(op) }

func (w *bytecodeWriter) writeInt32(v int32) { w.writeUInt64(uint64(uint32(v))) }

func (w *bytecodeWriter) writeUInt64(value uint64) {
	w.hasData = true
	appendValue := func(v uint64) {
		switch {
		case v == 0:
			w.code.WriteString("\\0")
		case v >= 0x20 && v <= 0x7e:
			w.appendEscaping(rune(v))
		default:
			fmt.Fprintf(&w.code, "\\u{%x}", v)
		}
	}
	v := value
	for v > 0x3f {
		appendValue(v&0x3f | 0x40)
		v >>= 6
	}
	appendValue(v)
}

func (w *bytecodeWriter) writeNullTerminatedString(s string) {
	for _, r := range s {
		w.appendEscaping(r)
	}
	w.writeUInt64(0)
}

func (w *bytecodeWriter) writeNullTerminatedStringArray(items []string) {
	w.writeUInt64(uint64(len(items)))
	for _, item := range items {
		w.writeNullTerminatedString(item)
	}
}

func (w *bytecodeWriter) appendEscaping(r rune) {
	w.hasData = true
	switch {
	case r == '\\' || r == '"':
		w.code.WriteByte('\\')
		w.code.WriteRune(r)
	case r < ' ' || r == 0x7f:
		fmt.Fprintf(&w.code, "\\u{%x}", r)
	default:
		w.code.WriteRune(r)
	}
}

// nameInstructionWriter writes the name instructions of the fields of a message or the values of an enum.
type nameInstructionWriter struct {
	bytecode       *bytecodeWriter
	previousNumber int32
	hasPrevious    bool
	err            error
}

func newNameInstructionWriter() *nameInstructionWriter {
	return &nameInstructionWriter{bytecode: newBytecodeWriter()}
}

func (w *nameInstructionWriter) shouldUseEmptyNameMapInitializer() bool {
	return !w.bytecode.hasData
}

func (w *nameInstructionWriter) delta(number int32) int32 {
	if w.err != nil {
		return 0
	}
	if !w.hasPrevious {
		w.hasPrevious = true
		w.previousNumber = number
		return number
	}
	if number < w.previousNumber {
		w.err = fmt.Errorf("field and case numbers must be written in ascending order")
		return 0
	}
	d := number - w.previousNumber
	w.previousNumber = number
	return d
}

func (w *nameInstructionWriter) writeOp(number int32, next, deltaOp uint64) {
	d := w.delta(number)
	if w.err != nil {
		return
	}
	if d == 1 {
		w.bytecode.writeOpcode(next)
		return
	}
	w.bytecode.writeOpcode(deltaOp)
	w.bytecode.writeInt32(d)
}

func (w *nameInstructionWriter) writeSame(number int32, name string) {
	w.writeOp(number, opSameNext, opSameDelta)
	if w.err == nil {
		w.bytecode.writeNullTerminatedString(name)
	}
}

func (w *nameInstructionWriter) writeStandard(number int32, name string) {
	w.writeOp(number, opStandardNext, opStandardDelta)
	if w.err == nil {
		w.bytecode.writeNullTerminatedString(name)
	}
}

func (w *nameInstructionWriter) writeUnique(number int32, protoName, jsonName string) {
	w.writeOp(number, opUniqueNext, opUniqueDelta)
	if w.err == nil {
		w.bytecode.writeNullTerminatedString(protoName)
		w.bytecode.writeNullTerminatedString(jsonName)
	}
}

func (w *nameInstructionWriter) writeGroup(number int32, name string) {
	w.writeOp(number, opGroupNext, opGroupDelta)
	if w.err == nil {
		w.bytecode.writeNullTerminatedString(name)
	}
}

func (w *nameInstructionWriter) writeAliased(v *enumValue, aliases []*enumValue) {
	w.writeOp(v.number, opAliasNext, opAliasDelta)
	if w.err != nil {
		return
	}
	w.bytecode.writeNullTerminatedString(v.name)
	names := make([]string, len(aliases))
	for i, a := range aliases {
		names[i] = a.name
	}
	w.bytecode.writeNullTerminatedStringArray(names)
}

func (w *nameInstructionWriter) writeReservedName(name string) {
	w.bytecode.writeOpcode(opReservedName)
	w.bytecode.writeNullTerminatedString(name)
}

func (w *nameInstructionWriter) writeReservedNumbers(start, end int32) {
	w.bytecode.writeOpcode(opReservedNumbers)
	w.bytecode.writeInt32(start)
	w.bytecode.writeInt32(end - start)
}
