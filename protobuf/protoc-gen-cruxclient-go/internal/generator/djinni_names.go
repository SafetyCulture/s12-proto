// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/common"
)

var cppKeywords = map[string]struct{}{
	"alignas": {}, "alignof": {}, "and": {}, "and_eq": {}, "asm": {}, "atomic_cancel": {},
	"atomic_commit": {}, "atomic_noexcept": {}, "auto": {}, "bitand": {}, "bitor": {}, "bool": {},
	"break": {}, "case": {}, "catch": {}, "char": {}, "char8_t": {}, "char16_t": {}, "char32_t": {},
	"class": {}, "compl": {}, "concept": {}, "const": {}, "consteval": {}, "constexpr": {},
	"constinit": {}, "const_cast": {}, "continue": {}, "co_await": {}, "co_return": {}, "co_yield": {},
	"decltype": {}, "default": {}, "delete": {}, "do": {}, "double": {}, "dynamic_cast": {}, "else": {},
	"enum": {}, "explicit": {}, "export": {}, "extern": {}, "false": {}, "final": {}, "float": {}, "for": {},
	"friend": {}, "goto": {}, "if": {}, "import": {}, "inline": {}, "int": {}, "long": {}, "module": {},
	"mutable": {}, "namespace": {}, "new": {}, "noexcept": {}, "not": {}, "not_eq": {}, "nullptr": {},
	"operator": {}, "or": {}, "or_eq": {}, "override": {}, "private": {}, "protected": {}, "public": {},
	"reflexpr": {}, "register": {}, "reinterpret_cast": {}, "requires": {}, "return": {}, "short": {},
	"signed": {}, "sizeof": {}, "static": {}, "static_assert": {}, "static_cast": {}, "struct": {}, "switch": {},
	"synchronized": {}, "template": {}, "this": {}, "thread_local": {}, "throw": {}, "true": {}, "try": {},
	"typedef": {}, "typeid": {}, "typename": {}, "union": {}, "unsigned": {}, "using": {}, "virtual": {},
	"void": {}, "volatile": {}, "wchar_t": {}, "while": {}, "xor": {}, "xor_eq": {},
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
			if err := validateCPPName("proto package", string(fileDesc.Package()), segment); err != nil {
				return err
			}
		}
		for _, message := range common.GetMessagesFromFile(fileDesc) {
			if _, nested := message.Parent().(protoreflect.MessageDescriptor); nested && message.Name() == "Translator" {
				return fmt.Errorf("nested message %q has unsupported name %q", message.FullName(), message.Name())
			}
			if err := validateCPPName("message", string(message.FullName()), string(message.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateCPPName(kind, fullName, segment string) error {
	if segment == "std" {
		return fmt.Errorf("%s %q has unsupported segment %q", kind, fullName, segment)
	}
	if _, keyword := cppKeywords[segment]; keyword {
		return fmt.Errorf("%s %q has unsupported C++ keyword %q", kind, fullName, segment)
	}
	return nil
}
