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

// cppQualifiedName returns the C++ name of a proto full name as protoc writes it: "::" between the
// parts and "_" after a C++ keyword.
func cppQualifiedName(fullName protoreflect.FullName) string {
	parts := strings.Split(string(fullName), ".")
	for i, part := range parts {
		if cppKeywords[part] {
			parts[i] = part + "_"
		}
	}
	return strings.Join(parts, "::")
}

// ValidateDjinniNames rejects proto names that cannot form a Djinni translator namespace.
func ValidateDjinniNames(files []protoreflect.FileDescriptor, options Options) error {
	if !options.SwiftEnabled() && !options.KotlinEnabled() {
		return nil
	}
	for _, fileDesc := range files {
		segments := strings.Split(string(fileDesc.Package()), ".")
		for index, segment := range segments {
			if options.SwiftEnabled() && index == 0 && segment == "swift" {
				return fmt.Errorf("proto package %q starts with unsupported segment %q", fileDesc.Package(), segment)
			}
			if segment == "std" {
				return fmt.Errorf("proto package %q has unsupported segment %q", fileDesc.Package(), segment)
			}
		}
		for _, message := range common.GetMessagesFromFile(fileDesc) {
			if _, nested := message.Parent().(protoreflect.MessageDescriptor); nested && message.Name() == "Translator" {
				return fmt.Errorf("nested message %q has unsupported name %q", message.FullName(), message.Name())
			}
			if message.Name() == "std" {
				return fmt.Errorf("message %q has unsupported segment %q", message.FullName(), message.Name())
			}
		}
	}
	return nil
}
