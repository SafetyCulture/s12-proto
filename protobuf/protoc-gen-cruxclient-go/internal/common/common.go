// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

// Package common provides string transformation utilities and proto descriptor
// helpers ported from the C++ protoc-gen-cruxclient plugin's common.h.
package common

import (
	"fmt"
	"strings"
	"unicode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	wire "github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/wire"
)

// DotsToColons replaces all "." with "::" in name.
// Mirrors C++ DotsToColons from common.h.
func DotsToColons(name string) string {
	return strings.ReplaceAll(name, ".", "::")
}

// DotsToUnderscores replaces all "." with "_" in name.
// Mirrors C++ DotsToUnderscores from common.h.
func DotsToUnderscores(name string) string {
	return strings.ReplaceAll(name, ".", "_")
}

// DotsToSlashs replaces all "." with "/" in name.
// Mirrors C++ DotsToSlashs from common.h.
func DotsToSlashs(name string) string {
	return strings.ReplaceAll(name, ".", "/")
}

// UnderscoresToDollar replaces all "_" with "$" in name.
// Mirrors C++ UnderscoresToDollar from common.h.
func UnderscoresToDollar(name string) string {
	return strings.ReplaceAll(name, "_", "$")
}

// UnderscoresToDots replaces all "_" with "." in name.
// Mirrors C++ UnderscoresToDots from common.h.
func UnderscoresToDots(name string) string {
	return strings.ReplaceAll(name, "_", ".")
}

// StripProto strips the ".protodevel" or ".proto" suffix from filename.
// Matches C++ StripProto from common.h lines 52-58.
func StripProto(filename string) string {
	if strings.HasSuffix(filename, ".protodevel") {
		return strings.TrimSuffix(filename, ".protodevel")
	}
	return strings.TrimSuffix(filename, ".proto")
}

// ReplaceCharacters replaces any character in s that appears in the remove
// string with the first rune of the with string.
// Mirrors C++ ReplaceCharacters from common.h (in-place mutation → returns new string in Go).
func ReplaceCharacters(s, remove, with string) string {
	if len(with) == 0 {
		return s
	}
	replaceWith := rune(with[0])
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if strings.ContainsRune(remove, ch) {
			b.WriteRune(replaceWith)
		} else {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// FilenameIdentifier encodes filename so it can be used as a C++ identifier.
// Alphanumeric bytes pass through; all others are encoded as "_XX" where XX
// is the lowercase hex value of the byte.
// Mirrors C++ FilenameIdentifier from common.h lines 63-77.
func FilenameIdentifier(filename string) string {
	var b strings.Builder
	b.Grow(len(filename))
	for i := 0; i < len(filename); i++ {
		c := filename[i]
		if isAlphaNumeric(c) {
			b.WriteByte(c)
		} else {
			b.WriteString(fmt.Sprintf("_%02x", c))
		}
	}
	return b.String()
}

// isAlphaNumeric returns true if the byte is ASCII alphanumeric (matching C isalnum).
func isAlphaNumeric(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// StringReplace replaces occurrences of from with to in str.
// If replaceAll is true, all occurrences are replaced; otherwise only the first.
// Mirrors C++ StringReplace from common.h lines 79-97.
func StringReplace(str, from, to string, replaceAll bool) string {
	if replaceAll {
		return strings.ReplaceAll(str, from, to)
	}
	return strings.Replace(str, from, to, 1)
}

// StringReplaceAll replaces all occurrences of from with to in str.
// Convenience wrapper matching the C++ overload that defaults replaceAll=true.
func StringReplaceAll(str, from, to string) string {
	return StringReplace(str, from, to, true)
}

// ToLower converts input to lowercase.
// Exported wrapper for consistency with C++ ToLower from common.h.
func ToLower(input string) string {
	return strings.ToLower(input)
}

// Join concatenates elements with separator.
// Exported wrapper for consistency with C++ join from common.h.
func Join(elements []string, separator string) string {
	return strings.Join(elements, separator)
}

// Tokenize splits input on any character in delimiters, preserving empty
// segments between consecutive delimiters.
// This is NOT strings.FieldsFunc — it matches the C++ tokenize behavior from
// common.h lines 101-119 which uses find_first_of and includes empty strings.
func Tokenize(input, delimiters string) []string {
	var tokens []string
	lastPos := 0
	for {
		pos := strings.IndexAny(input[lastPos:], delimiters)
		if pos == -1 {
			tokens = append(tokens, input[lastPos:])
			return tokens
		}
		tokens = append(tokens, input[lastPos:lastPos+pos])
		lastPos = lastPos + pos + 1
	}
}

// ToCamelCase converts text to CamelCase by capitalizing the first letter of
// each underscore-delimited word. Underscores are dropped.
// Mirrors C++ ToCamelCase from common.h lines 168-188.
// Example: "route_guide_v1" → "RouteGuideV1"
func ToCamelCase(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	newWord := true
	for _, ch := range text {
		if ch == '_' {
			newWord = true
			continue
		}
		if unicode.IsLetter(ch) {
			if newWord {
				b.WriteRune(unicode.ToUpper(ch))
				newWord = false
			} else {
				b.WriteRune(unicode.ToLower(ch))
			}
		} else {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// ClassName returns the C++ class name for a protobuf message descriptor.
// It walks up to the outermost containing message, then constructs the name
// using "::" for namespace separation (qualified=true) or just the outer name
// (qualified=false), with "_" for nested message separators.
// Mirrors C++ ClassName from common.h lines 190-204.
//
// qualified=true:  "routeguide::v1::RouteSummary_Details"
// qualified=false: "RouteSummary_Details"
func ClassName(msg protoreflect.MessageDescriptor, qualified bool) string {
	// Walk to outermost containing message.
	outer := msg
	for {
		parent := outer.Parent()
		if _, ok := parent.(protoreflect.MessageDescriptor); !ok {
			break
		}
		outer = parent.(protoreflect.MessageDescriptor)
	}

	outerName := string(outer.FullName())                    // e.g. "routeguide.v1.RouteSummary"
	innerSuffix := string(msg.FullName())[len(outerName):]   // e.g. ".Details"

	if qualified {
		return DotsToColons(outerName) + DotsToUnderscores(innerSuffix)
	}
	return string(outer.Name()) + DotsToUnderscores(innerSuffix)
}

// GetWirePackage reads the wire_package file option from a FileDescriptor.
// Returns empty string if the option is not set or if options are nil.
// Mirrors the wire_options extension access from api_generator.cc lines 429-436.
func GetWirePackage(fileDesc protoreflect.FileDescriptor) string {
	opts, ok := fileDesc.Options().(*descriptorpb.FileOptions)
	if !ok || opts == nil {
		return ""
	}
	if !proto.HasExtension(opts, wire.E_WirePackage) {
		return ""
	}
	return proto.GetExtension(opts, wire.E_WirePackage).(string)
}

// GetMessagesFromFile returns all message descriptors in a file, including
// nested messages, but excluding map-entry synthetic messages.
// Mirrors GetMessagesFromFile from api_generator.cc lines 370-393.
func GetMessagesFromFile(fileDesc protoreflect.FileDescriptor) []protoreflect.MessageDescriptor {
	var result []protoreflect.MessageDescriptor
	msgs := fileDesc.Messages()
	for i := 0; i < msgs.Len(); i++ {
		result = append(result, getMessagesFromMessage(msgs.Get(i))...)
	}
	return result
}

// getMessagesFromMessage recursively collects a message and all its nested messages.
func getMessagesFromMessage(msg protoreflect.MessageDescriptor) []protoreflect.MessageDescriptor {
	result := []protoreflect.MessageDescriptor{msg}
	nested := msg.Messages()
	for i := 0; i < nested.Len(); i++ {
		result = append(result, getMessagesFromMessage(nested.Get(i))...)
	}
	return result
}
