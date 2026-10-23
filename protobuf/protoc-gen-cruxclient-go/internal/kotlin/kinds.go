// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package kotlin

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Wire types of the protobuf encoding.
const (
	wireVarint  = 0
	wireFixed64 = 1
	wireLength  = 2
	wireFixed32 = 5
)

// wireTypeOf returns the wire type of a value of the kind outside a packed run.
func wireTypeOf(k protoreflect.Kind) int {
	switch k {
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind, protoreflect.DoubleKind:
		return wireFixed64
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind, protoreflect.FloatKind:
		return wireFixed32
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind:
		return wireLength
	}
	return wireVarint
}

// isPackable reports whether a repeated field of the kind can use the packed encoding.
func isPackable(k protoreflect.Kind) bool {
	return wireTypeOf(k) != wireLength
}

// tagOf returns the tag of a field number and a wire type, as the signed 32 bits that the reader returns.
func tagOf(number protoreflect.FieldNumber, wireType int) int32 {
	return int32(uint32(number)<<3 | uint32(wireType))
}

// scalarType returns the Kotlin type of a scalar kind.
func scalarType(k protoreflect.Kind) string {
	switch k {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "kotlin.Int"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "kotlin.Long"
	case protoreflect.BoolKind:
		return "kotlin.Boolean"
	case protoreflect.FloatKind:
		return "kotlin.Float"
	case protoreflect.DoubleKind:
		return "kotlin.Double"
	case protoreflect.StringKind:
		return "kotlin.String"
	}
	return ""
}

// writeMethod returns the WireWriter method that writes a field of the kind with its tag.
func writeMethod(k protoreflect.Kind) string {
	switch k {
	case protoreflect.Int32Kind, protoreflect.EnumKind:
		return "writeInt32"
	case protoreflect.Uint32Kind:
		return "writeUInt32"
	case protoreflect.Sint32Kind:
		return "writeSInt32"
	case protoreflect.Int64Kind, protoreflect.Uint64Kind:
		return "writeInt64"
	case protoreflect.Sint64Kind:
		return "writeSInt64"
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		return "writeFixed32"
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		return "writeFixed64"
	case protoreflect.BoolKind:
		return "writeBool"
	case protoreflect.FloatKind:
		return "writeFloat"
	case protoreflect.DoubleKind:
		return "writeDouble"
	case protoreflect.StringKind:
		return "writeString"
	case protoreflect.BytesKind:
		return "writeBytes"
	}
	return "writeMessage"
}

// readMethod returns the WireReader method that reads a value of the kind.
func readMethod(k protoreflect.Kind) string {
	switch k {
	case protoreflect.Int32Kind, protoreflect.Uint32Kind, protoreflect.EnumKind:
		return "readInt32"
	case protoreflect.Sint32Kind:
		return "readSInt32"
	case protoreflect.Int64Kind, protoreflect.Uint64Kind:
		return "readInt64"
	case protoreflect.Sint64Kind:
		return "readSInt64"
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		return "readFixed32"
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		return "readFixed64"
	case protoreflect.BoolKind:
		return "readBool"
	case protoreflect.FloatKind:
		return "readFloat"
	case protoreflect.DoubleKind:
		return "readDouble"
	case protoreflect.StringKind:
		return "readString"
	case protoreflect.BytesKind:
		return "readBytes"
	}
	return ""
}

// packedSuffix returns the suffix of the WireWriter method that writes a packed run of the kind.
func packedSuffix(k protoreflect.Kind) string {
	switch k {
	case protoreflect.Int32Kind, protoreflect.EnumKind:
		return "Int32"
	case protoreflect.Uint32Kind:
		return "UInt32"
	case protoreflect.Sint32Kind:
		return "SInt32"
	case protoreflect.Int64Kind, protoreflect.Uint64Kind:
		return "Int64"
	case protoreflect.Sint64Kind:
		return "SInt64"
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		return "Fixed32"
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		return "Fixed64"
	case protoreflect.BoolKind:
		return "Bool"
	case protoreflect.FloatKind:
		return "Float"
	}
	return "Double"
}

// packedPayload returns the Kotlin expression of the size of a packed run of the kind.
func (g *fileGen) packedPayload(k protoreflect.Kind, list string) string {
	switch k {
	case protoreflect.Int32Kind, protoreflect.EnumKind:
		return g.rt("WireSize") + ".packedInt32Payload(" + list + ")"
	case protoreflect.Uint32Kind:
		return g.rt("WireSize") + ".packedUInt32Payload(" + list + ")"
	case protoreflect.Sint32Kind:
		return g.rt("WireSize") + ".packedSInt32Payload(" + list + ")"
	case protoreflect.Int64Kind, protoreflect.Uint64Kind:
		return g.rt("WireSize") + ".packedInt64Payload(" + list + ")"
	case protoreflect.Sint64Kind:
		return g.rt("WireSize") + ".packedSInt64Payload(" + list + ")"
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind, protoreflect.FloatKind:
		return list + ".size * 4"
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind, protoreflect.DoubleKind:
		return list + ".size * 8"
	}
	return list + ".size"
}

// sizeCall returns the Kotlin expression of the encoded size of one value with its tag.
func (g *fileGen) sizeCall(k protoreflect.Kind, number protoreflect.FieldNumber, expr string) string {
	ws := g.rt("WireSize")
	switch k {
	case protoreflect.Int32Kind:
		return fmt.Sprintf("%s.int32(%d, %s)", ws, number, expr)
	case protoreflect.EnumKind:
		return fmt.Sprintf("%s.int32(%d, %s.number)", ws, number, expr)
	case protoreflect.Uint32Kind:
		return fmt.Sprintf("%s.uint32(%d, %s)", ws, number, expr)
	case protoreflect.Sint32Kind:
		return fmt.Sprintf("%s.sint32(%d, %s)", ws, number, expr)
	case protoreflect.Int64Kind, protoreflect.Uint64Kind:
		return fmt.Sprintf("%s.int64(%d, %s)", ws, number, expr)
	case protoreflect.Sint64Kind:
		return fmt.Sprintf("%s.sint64(%d, %s)", ws, number, expr)
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind, protoreflect.FloatKind:
		return fmt.Sprintf("%s.fixed32(%d)", ws, number)
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind, protoreflect.DoubleKind:
		return fmt.Sprintf("%s.fixed64(%d)", ws, number)
	case protoreflect.BoolKind:
		return fmt.Sprintf("%s.bool(%d)", ws, number)
	case protoreflect.StringKind:
		return fmt.Sprintf("%s.string(%d, %s)", ws, number, expr)
	case protoreflect.BytesKind:
		return fmt.Sprintf("%s.bytes(%d, %s)", ws, number, expr)
	}
	return fmt.Sprintf("%s.message(%d, %s)", ws, number, expr)
}

// writeCall returns the Kotlin statement that writes one value with its tag.
func writeCall(k protoreflect.Kind, number protoreflect.FieldNumber, expr string) string {
	if k == protoreflect.EnumKind {
		expr += ".number"
	}
	return fmt.Sprintf("writer.%s(%d, %s)", writeMethod(k), number, expr)
}

// kotlinString returns a Kotlin string literal.
func kotlinString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch {
		case r == '\\':
			b.WriteString(`\\`)
		case r == '"':
			b.WriteString(`\"`)
		case r == '$':
			b.WriteString(`\$`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20 || r == 0x7f:
			fmt.Fprintf(&b, `\u%04x`, r)
		case r > 0xFFFF:
			r -= 0x10000
			fmt.Fprintf(&b, `\u%04x\u%04x`, 0xD800+(r>>10), 0xDC00+(r&0x3FF))
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// defaultLiteral returns the Kotlin expression of the default value of a scalar or enum field.
func (g *fileGen) defaultLiteral(fd protoreflect.FieldDescriptor) string {
	d := fd.Default()
	switch fd.Kind() {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return strconv.FormatInt(d.Int(), 10)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return strconv.FormatInt(int64(int32(uint32(d.Uint()))), 10)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if d.Int() == math.MinInt64 {
			return "kotlin.Long.MIN_VALUE"
		}
		return strconv.FormatInt(d.Int(), 10) + "L"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v := int64(d.Uint())
		if v == math.MinInt64 {
			return "kotlin.Long.MIN_VALUE"
		}
		return strconv.FormatInt(v, 10) + "L"
	case protoreflect.BoolKind:
		return strconv.FormatBool(d.Bool())
	case protoreflect.FloatKind:
		return floatLiteral(d.Float(), "kotlin.Float", 32, "f")
	case protoreflect.DoubleKind:
		return floatLiteral(d.Float(), "kotlin.Double", 64, "")
	case protoreflect.StringKind:
		return kotlinString(d.String())
	case protoreflect.BytesKind:
		b := d.Bytes()
		if len(b) == 0 {
			return g.rt("ByteString") + ".EMPTY"
		}
		parts := make([]string, len(b))
		for i, c := range b {
			parts[i] = strconv.Itoa(int(int8(c)))
		}
		return g.rt("ByteString") + ".copyFrom(kotlin.byteArrayOf(" + strings.Join(parts, ", ") + "))"
	case protoreflect.EnumKind:
		return qualifiedName(fd.Enum(), g.root) + "." + entryName(string(defaultEnumValue(fd).Name()))
	}
	return "null"
}

func floatLiteral(v float64, typ string, bits int, suffix string) string {
	switch {
	case math.IsNaN(v):
		return typ + ".NaN"
	case math.IsInf(v, 1):
		return typ + ".POSITIVE_INFINITY"
	case math.IsInf(v, -1):
		return typ + ".NEGATIVE_INFINITY"
	}
	s := strconv.FormatFloat(v, 'g', -1, bits)
	if !strings.ContainsAny(s, ".e") {
		s += ".0"
	}
	return s + suffix
}

// defaultEnumValue returns the default value of an enum field: the declared default,
// else the entry with the number 0, else the first entry.
func defaultEnumValue(fd protoreflect.FieldDescriptor) protoreflect.EnumValueDescriptor {
	if v := fd.DefaultEnumValue(); v != nil {
		return v
	}
	values := fd.Enum().Values()
	if v := values.ByNumber(0); v != nil {
		return v
	}
	return values.Get(0)
}
