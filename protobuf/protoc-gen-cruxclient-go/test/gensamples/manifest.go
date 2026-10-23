// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/kotlin"
)

// manifest describes the accessors of every message of the files, so that the test can read the
// same field from a Kotlin message and from a Java lite message. A message is a line "M", then a
// line "F" for each field with these columns:
//
//	number, Kotlin property, Java name, shape, kind, open enum, Kotlin class of a message value,
//	oneof property, Java name of the oneof, Kotlin class of the oneof member.
//
// The shape is single, optional, list, map or oneof. The kind is the kind of the value: int, long,
// bool, float, double, string, bytes, enum or message.
func manifest(files *protoregistry.Files, paths map[string]bool, root string) string {
	var sorted []string
	for path := range paths {
		sorted = append(sorted, path)
	}
	sort.Strings(sorted)
	var b strings.Builder
	for _, path := range sorted {
		fd, err := files.FindFileByPath(path)
		check(err)
		walkMessages(fd.Messages(), func(md protoreflect.MessageDescriptor) {
			fmt.Fprintf(&b, "M\t%s\n", kotlin.BinaryClassName(md, root))
			names := javaFieldNames(md)
			for i := 0; i < md.Fields().Len(); i++ {
				fd := md.Fields().Get(i)
				b.WriteString(fieldLine(fd, names[fd.Number()], root))
			}
		})
	}
	return b.String()
}

func fieldLine(fd protoreflect.FieldDescriptor, javaName, root string) string {
	value := fd
	shape := "single"
	oneofProp, oneofJava, member := "", "", ""
	switch {
	case fd.IsMap():
		shape = "map"
		value = fd.MapValue()
	case fd.IsList():
		shape = "list"
	case fd.ContainingOneof() != nil && !fd.ContainingOneof().IsSynthetic():
		shape = "oneof"
		oneofProp = kotlin.OneofPropertyName(fd.ContainingOneof())
		oneofJava = javaCamel(string(fd.ContainingOneof().Name()))
		member = kotlin.OneofMember(fd)
	case fd.Message() != nil || (fd.HasPresence() && fd.Cardinality() != protoreflect.Required):
		shape = "optional"
	}
	kind, open, class := valueKind(value, root)
	return fmt.Sprintf("F\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", fd.Number(), kotlin.PropertyName(fd), javaName, shape, kind, open, class, oneofProp, oneofJava, member)
}

// javaFieldNames returns the name that Java lite uses in the accessors of each field of a message.
// When two fields would make the same accessor, protoc appends the field number to the name of both:
// a repeated field and a singular field named like its count or its list, and an open enum field and
// a field named like its value accessor.
func javaFieldNames(md protoreflect.MessageDescriptor) map[protoreflect.FieldNumber]string {
	fields := md.Fields()
	conflict := map[protoreflect.FieldNumber]bool{}
	for i := 0; i < fields.Len(); i++ {
		for j := i + 1; j < fields.Len(); j++ {
			if javaAccessorsClash(fields.Get(i), fields.Get(j)) || javaAccessorsClash(fields.Get(j), fields.Get(i)) {
				conflict[fields.Get(i).Number()] = true
				conflict[fields.Get(j).Number()] = true
			}
		}
	}
	names := map[protoreflect.FieldNumber]string{}
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		names[fd.Number()] = javaCamel(string(fd.Name()))
		if conflict[fd.Number()] {
			names[fd.Number()] += fmt.Sprint(fd.Number())
		}
	}
	return names
}

// javaAccessorsClash reports whether a field makes an accessor with the name of one of the accessors
// of another field.
func javaAccessorsClash(other, field protoreflect.FieldDescriptor) bool {
	a, b := javaCamel(string(other.Name())), javaCamel(string(field.Name()))
	if a == b {
		return true
	}
	if other.IsList() && !field.IsList() && !field.IsMap() {
		return b == a+"Count" || b == a+"List"
	}
	if other.Kind() == protoreflect.EnumKind && !other.Enum().IsClosed() {
		return b == a+"Value"
	}
	return false
}

func valueKind(fd protoreflect.FieldDescriptor, root string) (kind, open, class string) {
	open = "0"
	switch fd.Kind() {
	case protoreflect.BoolKind:
		kind = "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		kind = "int"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		kind = "long"
	case protoreflect.FloatKind:
		kind = "float"
	case protoreflect.DoubleKind:
		kind = "double"
	case protoreflect.StringKind:
		kind = "string"
	case protoreflect.BytesKind:
		kind = "bytes"
	case protoreflect.EnumKind:
		kind = "enum"
		if !fd.Enum().IsClosed() {
			open = "1"
		}
	default:
		kind = "message"
		class = kotlin.BinaryClassName(fd.Message(), root)
	}
	return kind, open, class
}

// javaCamel returns the name that protoc gives a field in a Java accessor: the letter after an
// underscore or after a digit is a capital.
func javaCamel(name string) string {
	var b strings.Builder
	up := true
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			if up {
				r -= 'a' - 'A'
			}
			up = false
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			up = false
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			up = true
			b.WriteRune(r)
		default:
			up = true
		}
	}
	return b.String()
}
