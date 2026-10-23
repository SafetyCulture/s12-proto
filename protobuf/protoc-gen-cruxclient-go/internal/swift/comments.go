// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// Field numbers of the descriptor messages, used in source info paths.
const (
	filePathSyntax      = 12
	filePathEdition     = 14
	filePathMessageType = 4
	filePathEnumType    = 5
	filePathExtension   = 7

	messagePathField      = 2
	messagePathNestedType = 3
	messagePathEnumType   = 4
	messagePathExtension  = 6
	messagePathOneofDecl  = 8

	enumPathValue = 2
)

func (m *message) locationPath() []int32 {
	if m.containing != nil {
		return append(m.containing.locationPath(), messagePathNestedType, int32(m.index))
	}
	return []int32{filePathMessageType, int32(m.index)}
}

func (e *enum) locationPath() []int32 {
	if e.containing != nil {
		return append(e.containing.locationPath(), messagePathEnumType, int32(e.index))
	}
	return []int32{filePathEnumType, int32(e.index)}
}

func (v *enumValue) locationPath() []int32 {
	return append(v.enum.locationPath(), enumPathValue, int32(v.index))
}

func (o *oneof) locationPath() []int32 {
	return append(o.containing.locationPath(), messagePathOneofDecl, int32(o.index))
}

func (f *field) locationPath() []int32 {
	if f.isExtension {
		if f.extScope != nil {
			return append(f.extScope.locationPath(), messagePathExtension, int32(f.index))
		}
		return []int32{filePathExtension, int32(f.index)}
	}
	return append(f.containing.locationPath(), messagePathField, int32(f.index))
}

// leadingOrTrailing returns the leading comment of a location, or the trailing comment when it has none.
func leadingOrTrailing(loc *descriptorpb.SourceCodeInfo_Location) string {
	if loc.LeadingComments != nil {
		return loc.GetLeadingComments()
	}
	return loc.GetTrailingComments()
}

// asSourceComment builds the comment lines of a location.
func asSourceComment(loc *descriptorpb.SourceCodeInfo_Location, commentPrefix string) string {
	return prefixLines(leadingOrTrailing(loc), commentPrefix)
}

// asSourceCommentWithDetached builds the comment lines of a location, with its detached comments before them.
func asSourceCommentWithDetached(loc *descriptorpb.SourceCodeInfo_Location, commentPrefix, detachedPrefix string) string {
	var result strings.Builder
	for _, detached := range loc.GetLeadingDetachedComments() {
		comment := prefixLines(detached, detachedPrefix)
		if comment != "" {
			result.WriteString(comment)
			// A blank line follows a detached comment.
			result.WriteString("\n")
		}
	}
	result.WriteString(asSourceComment(loc, commentPrefix))
	return result.String()
}

// comments returns the doc comment of a descriptor, which is empty when the options strip comments.
func (g *generator) comments(f *file, path []int32) string {
	if g.opts.stripNonfunctionalCodegen {
		return ""
	}
	loc := f.location(path...)
	if loc == nil {
		return ""
	}
	return asSourceComment(loc, "///")
}

func withDeprecation(comments, deprecation string) string {
	switch {
	case deprecation == "":
		return comments
	case comments == "":
		return deprecation
	}
	return comments + "///\n" + deprecation
}

func (g *generator) commentsWithDeprecation(f *file, path []int32, deprecation string) string {
	return withDeprecation(g.comments(f, path), deprecation)
}

// typeDeprecation is the deprecation note of a message or an enum, which a deprecated file also sets.
func typeDeprecation(typeName string, deprecated bool, f *file) string {
	if deprecated {
		return "/// NOTE: This " + typeName + " was marked as deprecated in the .proto file.\n"
	}
	if f.proto.GetOptions().GetDeprecated() {
		return "/// NOTE: The whole .proto file that defined this " + typeName + " was marked as deprecated.\n"
	}
	return ""
}

func (g *generator) messageComments(m *message) string {
	return g.commentsWithDeprecation(m.file, m.locationPath(), typeDeprecation("message", m.proto.GetOptions().GetDeprecated(), m.file))
}

func (g *generator) enumComments(e *enum) string {
	return g.commentsWithDeprecation(e.file, e.locationPath(), typeDeprecation("enum", e.proto.GetOptions().GetDeprecated(), e.file))
}

func (g *generator) enumValueComments(v *enumValue) string {
	deprecation := ""
	if v.proto.GetOptions().GetDeprecated() {
		deprecation = "/// NOTE: This enum value was marked as deprecated in the .proto file\n"
	}
	return g.commentsWithDeprecation(v.enum.file, v.locationPath(), deprecation)
}

func (g *generator) fieldComments(f *field) string {
	deprecation := ""
	switch {
	case f.proto.GetOptions().GetDeprecated():
		kind := "field"
		if f.isExtension {
			kind = "extension field"
		}
		deprecation = "/// NOTE: This " + kind + " was marked as deprecated in the .proto file.\n"
	case f.isExtension && f.file.proto.GetOptions().GetDeprecated():
		deprecation = "/// NOTE: The whole .proto file that defined this extension field was marked as deprecated.\n"
	}
	return g.commentsWithDeprecation(f.file, f.locationPath(), deprecation)
}

func (g *generator) oneofComments(o *oneof) string {
	return g.comments(o.containing.file, o.locationPath())
}
