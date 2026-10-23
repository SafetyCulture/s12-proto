// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/common"
)

// cppKeywords is kKeywordList of protoc 36.1 (compiler/cpp/helpers.cc). protoc appends "_" to a
// package segment or a message name in this list, so the glue must name the C++ types the same way.
var cppKeywords = map[string]bool{}

func init() {
	for _, keyword := range strings.Fields(`NULL alignas alignof and and_eq asm assert auto bitand bitor bool break case catch char class compl const
		constexpr const_cast continue decltype default delete do double dynamic_cast else enum explicit export extern
		false float for friend goto if inline int long mutable namespace new noexcept not not_eq nullptr operator or
		or_eq private protected public register reinterpret_cast return short signed sizeof static static_assert
		static_cast struct switch template this thread_local throw true try typedef typeid typename union unsigned
		using virtual void volatile wchar_t while xor xor_eq char8_t char16_t char32_t concept consteval constinit
		co_await co_return co_yield requires`) {
		cppKeywords[keyword] = true
	}
}

// cppKnownMethods is MessageKnownMethodsCamelCase of protoc 36.1 (compiler/cpp/helpers.cc). protoc
// appends "_" to a message name in this list, because the name would hide a method of the class.
var cppKnownMethods = map[string]bool{
	"GetDescriptor": true, "GetReflection": true, "default_instance": true, "Swap": true,
	"UnsafeArenaSwap": true, "New": true, "CopyFrom": true, "MergeFrom": true, "IsInitialized": true,
	"GetMetadata": true, "Clear": true,
}

// cppQualifiedName returns the C++ name of a message as protoc writes it: "::" between the parts,
// "_" after a package segment that is a C++ keyword, and "_" after a message name that is a C++
// keyword or a known method name.
func cppQualifiedName(msg protoreflect.MessageDescriptor) string {
	var names []string
	for d := protoreflect.Descriptor(msg); ; d = d.Parent() {
		m, ok := d.(protoreflect.MessageDescriptor)
		if !ok {
			break
		}
		name := string(m.Name())
		if cppKeywords[name] || cppKnownMethods[name] {
			name += "_"
		}
		names = append([]string{name}, names...)
	}
	var parts []string
	if pkg := string(msg.ParentFile().Package()); pkg != "" {
		for _, part := range strings.Split(pkg, ".") {
			if cppKeywords[part] {
				part += "_"
			}
			parts = append(parts, part)
		}
	}
	return strings.Join(append(parts, names...), "::")
}

// ValidateDjinniNames rejects proto names that cannot form a Djinni translator namespace.
func ValidateDjinniNames(files []protoreflect.FileDescriptor, options Options) error {
	if !options.SwiftEnabled() && !options.KotlinEnabled() {
		return nil
	}
	for _, fileDesc := range files {
		segments := strings.Split(string(fileDesc.Package()), ".")
		for index, segment := range segments {
			// The new translators live in djinni::swift:: and djinni::kotlin::, where a legacy translator
			// of a package that starts with that segment would have the same name.
			if options.SwiftEnabled() && index == 0 && segment == "swift" ||
				options.KotlinEnabled() && index == 0 && segment == "kotlin" {
				return fmt.Errorf("cruxclient-go: proto package %q starts with unsupported segment %q", fileDesc.Package(), segment)
			}
			if segment == "std" {
				return fmt.Errorf("cruxclient-go: proto package %q has unsupported segment %q", fileDesc.Package(), segment)
			}
		}
		for _, message := range common.GetMessagesFromFile(fileDesc) {
			// A nested message opens a namespace inside its parent's namespace, which already declares these.
			_, nested := message.Parent().(protoreflect.MessageDescriptor)
			if nested && (message.Name() == "Translator" || options.KotlinEnabled() && message.Name() == "JNIInfo") {
				return fmt.Errorf("cruxclient-go: nested message %q has unsupported name %q", message.FullName(), message.Name())
			}
			if message.Name() == "std" {
				return fmt.Errorf("cruxclient-go: message %q has unsupported segment %q", message.FullName(), message.Name())
			}
		}
	}
	return nil
}
