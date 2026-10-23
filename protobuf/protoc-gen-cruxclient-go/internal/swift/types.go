// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import "fmt"

func (g *generator) swiftType(f *field) string {
	if f.messageType != nil {
		if key, value, ok := f.messageType.mapKeyAndValue(); ok {
			return "Dictionary<" + g.swiftType(key) + "," + g.swiftType(value) + ">"
		}
	}
	var result string
	switch f.typ {
	case typeDouble:
		result = "Double"
	case typeFloat:
		result = "Float"
	case typeInt64, typeSfixed64, typeSint64:
		result = "Int64"
	case typeUint64, typeFixed64:
		result = "UInt64"
	case typeInt32, typeSfixed32, typeSint32:
		result = "Int32"
	case typeFixed32, typeUint32:
		result = "UInt32"
	case typeBool:
		result = "Bool"
	case typeString:
		result = "String"
	case typeGroup, typeMessage:
		result = g.namer.fullNameMessage(f.messageType)
	case typeBytes:
		result = "Data"
	case typeEnum:
		result = g.namer.fullNameEnum(f.enumType)
	}
	if f.isRepeated() {
		return "[" + result + "]"
	}
	return result
}

func (g *generator) swiftStorageType(f *field) string {
	swiftType := g.swiftType(f)
	if f.isRepeated() || f.realContainingOneof() != nil {
		return swiftType
	}
	if f.hasPresence() {
		return swiftType + "?"
	}
	return swiftType
}

func (g *generator) swiftDefaultValue(f *field) string {
	if f.isMap() {
		return "[:]"
	}
	if f.isRepeated() {
		return "[]"
	}
	if f.defaultValue != nil {
		defaultValue := *f.defaultValue
		switch f.typ {
		case typeDouble:
			switch defaultValue {
			case "inf":
				return "Double.infinity"
			case "-inf":
				return "-Double.infinity"
			case "nan", "-nan":
				return "Double.nan"
			}
			return defaultValue
		case typeFloat:
			switch defaultValue {
			case "inf":
				return "Float.infinity"
			case "-inf":
				return "-Float.infinity"
			case "nan", "-nan":
				return "Float.nan"
			}
			return defaultValue
		case typeString:
			return stringToEscapedStringLiteral(defaultValue)
		case typeBytes:
			return escapedToDataLiteral(defaultValue)
		case typeEnum:
			for _, v := range f.enumType.values {
				if v.name == defaultValue {
					return g.namer.dottedRelativeNameEnumValue(v)
				}
			}
			panic(fmt.Sprintf("%s: the default value %s is not a value of its enum", f.fullName, defaultValue))
		}
		return defaultValue
	}
	switch f.typ {
	case typeBool:
		return "false"
	case typeString:
		return "String()"
	case typeBytes:
		return "Data()"
	case typeGroup, typeMessage:
		return g.namer.fullNameMessage(f.messageType) + "()"
	case typeEnum:
		return g.namer.dottedRelativeNameEnumValue(f.enumType.values[0])
	}
	return "0"
}

// traitsType is the runtime type that decodes and visits the field. Maps and extensions use it.
func (g *generator) traitsType(f *field) string {
	rt := g.namer.runtimePrefix()
	if f.messageType != nil {
		if key, value, ok := f.messageType.mapKeyAndValue(); ok {
			keyTraits := g.traitsType(key)
			valueTraits := g.traitsType(value)
			switch value.typ {
			case typeMessage:
				return rt + "_ProtobufMessageMap<" + keyTraits + "," + valueTraits + ">"
			case typeEnum:
				return rt + "_ProtobufEnumMap<" + keyTraits + "," + valueTraits + ">"
			}
			return rt + "_ProtobufMap<" + keyTraits + "," + valueTraits + ">"
		}
	}
	switch f.typ {
	case typeDouble:
		return rt + "ProtobufDouble"
	case typeFloat:
		return rt + "ProtobufFloat"
	case typeInt64:
		return rt + "ProtobufInt64"
	case typeUint64:
		return rt + "ProtobufUInt64"
	case typeInt32:
		return rt + "ProtobufInt32"
	case typeFixed64:
		return rt + "ProtobufFixed64"
	case typeFixed32:
		return rt + "ProtobufFixed32"
	case typeBool:
		return rt + "ProtobufBool"
	case typeString:
		return rt + "ProtobufString"
	case typeGroup, typeMessage:
		return g.namer.fullNameMessage(f.messageType)
	case typeBytes:
		return rt + "ProtobufBytes"
	case typeUint32:
		return rt + "ProtobufUInt32"
	case typeEnum:
		return g.namer.fullNameEnum(f.enumType)
	case typeSfixed32:
		return rt + "ProtobufSFixed32"
	case typeSfixed64:
		return rt + "ProtobufSFixed64"
	case typeSint32:
		return rt + "ProtobufSInt32"
	case typeSint64:
		return rt + "ProtobufSInt64"
	}
	panic("unknown field type")
}
