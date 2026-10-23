// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package kotlin

import (
	"strings"
	"unicode"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// hardKeywords are the Kotlin hard keywords. They cannot name a declaration without backticks.
var hardKeywords = map[string]bool{}

// reservedProperties are property names that clash with a member of the generated class
// or with a method of java.lang.Object.
var reservedProperties = map[string]bool{
	"class": true, "serializedSize": true, "unknownFields": true,
	"missingRequired": true, "findMissingRequired": true,
	// The generated code names library types with their package, so a property called kotlin or com would hide them.
	"kotlin": true, "com": true,
}

// reservedEntries are enum entry names that clash with a member of a Kotlin enum class
// or with a member that the generator adds.
var reservedEntries = map[string]bool{
	"name": true, "ordinal": true, "number": true, "entries": true, "values": true,
	"valueOf": true, "fromNumber": true, "Companion": true, "UNRECOGNIZED": true,
}

func init() {
	for _, k := range strings.Fields(`as break class continue do else false for fun if in interface is null
		object package return super this throw true try typealias typeof val var when while`) {
		hardKeywords[k] = true
	}
}

// escapeIdent puts backticks around a hard keyword.
func escapeIdent(name string) string {
	if hardKeywords[name] {
		return "`" + name + "`"
	}
	return name
}

// className returns the Kotlin class name for a message or enum name. The name
// Companion would clash with the companion object of the enclosing class.
func className(name string) string {
	if name == "Companion" {
		return name + "Message"
	}
	return name
}

// SimpleClassName returns the name that a message or enum has as a class declaration, with the
// backticks that a hard keyword needs.
func SimpleClassName(d protoreflect.Descriptor) string {
	return escapeIdent(className(string(d.Name())))
}

// PropertyName returns the Kotlin property name of a field: the lower camel case field name.
// It is not the JSON name, because a json_name option can hold any text. A name that clashes
// gets a trailing underscore.
func PropertyName(fd protoreflect.FieldDescriptor) string {
	return safeProperty(lowerFirst(upperCamel(string(fd.Name()))))
}

// OneofPropertyName returns the Kotlin property name of a oneof.
func OneofPropertyName(od protoreflect.OneofDescriptor) string {
	return safeProperty(lowerFirst(upperCamel(string(od.Name()))))
}

func safeProperty(name string) string {
	if hardKeywords[name] || reservedProperties[name] {
		return name + "_"
	}
	return name
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// upperCamel converts a snake case name to upper camel case.
func upperCamel(s string) string {
	var b strings.Builder
	up := true
	for _, r := range s {
		if r == '_' {
			up = true
			continue
		}
		if up {
			r = unicode.ToUpper(r)
			up = false
		}
		b.WriteRune(r)
	}
	return b.String()
}

func entryName(name string) string {
	if reservedEntries[name] || hardKeywords[name] {
		return name + "_"
	}
	return name
}

// kotlinPackage returns the Kotlin package of the declarations of a file: the root
// package, then the proto package. The leading "s12" segment of the proto package is
// not repeated, because the root package already names it, as the Java package
// "com.safetyculture.s12.assets.v1" does.
func kotlinPackage(file protoreflect.FileDescriptor, root string) string {
	pkg := string(file.Package())
	if pkg == "s12" {
		return root
	}
	pkg = strings.TrimPrefix(pkg, "s12.")
	if pkg == "" {
		return root
	}
	return root + "." + pkg
}

// escapePackage puts backticks around each segment of a package name that is a hard keyword.
func escapePackage(pkg string) string {
	parts := strings.Split(pkg, ".")
	for i, p := range parts {
		parts[i] = escapeIdent(p)
	}
	return strings.Join(parts, ".")
}

// typeName returns the name of a message or enum inside its package, "Outer.Inner" for a nested one.
func typeName(d protoreflect.Descriptor) string {
	if parent, ok := d.Parent().(protoreflect.MessageDescriptor); ok {
		return typeName(parent) + "." + escapeIdent(className(string(d.Name())))
	}
	return escapeIdent(className(string(d.Name())))
}

// qualifiedName returns the fully qualified Kotlin name of a message or enum. Generated
// code uses it everywhere, so a message named String or List cannot shadow a library type.
func qualifiedName(d protoreflect.Descriptor, root string) string {
	return escapePackage(kotlinPackage(d.ParentFile(), root)) + "." + typeName(d)
}

// BinaryClassName returns the JVM binary name of the Kotlin class of a message or enum,
// with "." between the package parts and "$" before a nested class.
func BinaryClassName(d protoreflect.Descriptor, root string) string {
	return kotlinPackage(d.ParentFile(), root) + "." + strings.ReplaceAll(strings.ReplaceAll(typeName(d), "`", ""), ".", "$")
}
