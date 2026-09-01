// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The names below have to agree with what protoc's C# generator produced for the
// same descriptor, since the validators are emitted as further parts of the
// classes it declared. They follow csharp_helpers.cc.

// pascalCase converts a proto name the way protoc's C# generator does: a letter
// after a separator or a digit is capitalised, and the separator is dropped.
func pascalCase(name string) string {
	var out strings.Builder
	out.Grow(len(name))

	capitaliseNext := true
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			if capitaliseNext {
				out.WriteRune(r - ('a' - 'A'))
			} else {
				out.WriteRune(r)
			}
			capitaliseNext = false
		case r >= 'A' && r <= 'Z':
			out.WriteRune(r)
			capitaliseNext = false
		case r >= '0' && r <= '9':
			out.WriteRune(r)
			capitaliseNext = true
		default:
			capitaliseNext = true
		}
	}

	return out.String()
}

// fileNamespace is the namespace protoc's C# generator puts the file's types in.
func fileNamespace(f *protogen.File) string {
	if ns := f.Proto.GetOptions().GetCsharpNamespace(); ns != "" {
		return ns
	}

	parts := strings.Split(string(f.Desc.Package()), ".")
	for i, part := range parts {
		parts[i] = pascalCase(part)
	}
	return strings.Join(parts, ".")
}

// messageName is the class protoc's C# generator declared for the message. Unlike
// a field, a message keeps the name written in the .proto, underscores included.
func messageName(m protoreflect.MessageDescriptor) string {
	return string(m.Name())
}

// className is the message's name relative to its file's namespace. A message
// declared inside another sits under that message's nested Types container.
func className(m *protogen.Message) string {
	name := messageName(m.Desc)
	for parent, ok := m.Desc.Parent().(protoreflect.MessageDescriptor); ok; parent, ok = parent.Parent().(protoreflect.MessageDescriptor) {
		name = messageName(parent) + ".Types." + name
	}
	return name
}

// propertyName is the property protoc's C# generator declared for the field. A
// name that would collide with its own class gains a trailing underscore.
func propertyName(f *protogen.Field) string {
	name := pascalCase(string(f.Desc.Name()))
	// Types is the container the generator puts a message's nested declarations in,
	// so a field of that name would shadow it.
	if name == "Types" {
		return name + "_"
	}
	if f.Parent != nil && name == messageName(f.Parent.Desc) {
		name += "_"
	}
	return name
}

// oneofCaseProperty reports which member of the oneof is set.
func oneofCaseProperty(o *protogen.Oneof) string {
	return pascalCase(string(o.Desc.Name())) + "Case"
}

// oneofCaseEnum names the enum that oneofCaseProperty returns.
func oneofCaseEnum(o *protogen.Oneof) string {
	return pascalCase(string(o.Desc.Name())) + "OneofCase"
}

// hasProperty reports whether a field with the optional keyword carries a value.
func hasProperty(f *protogen.Field) string {
	return "Has" + propertyName(f)
}

// stringLiteral quotes a value as a C# verbatim string, which needs no escaping
// beyond doubling the quote character.
func stringLiteral(s string) string {
	return `@"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
