// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package kotlin

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Numbers in google/protobuf/java_features.proto: the extension (pb.java) of FeatureSet, its field
// utf8_validation, and the value VERIFY of that field.
const (
	javaFeaturesExtension    = 1001
	javaUTF8ValidationTag    = 2
	javaUTF8ValidationVerify = 2
)

// javaUTF8Verify reports whether protoc Java checks the UTF-8 of a string field for the reason that
// is specific to Java: the file option java_string_check_utf8, or the Java feature utf8_validation
// set to VERIFY on the field, an enclosing message or the file, the nearest one winning.
func javaUTF8Verify(fd protoreflect.FieldDescriptor) bool {
	if opts, ok := fd.ParentFile().Options().(*descriptorpb.FileOptions); ok && opts.GetJavaStringCheckUtf8() {
		return true
	}
	if features, ok := fd.Options().(*descriptorpb.FieldOptions); ok {
		if verify, set := javaUTF8Feature(features.GetFeatures()); set {
			return verify
		}
	}
	for d := fd.Parent(); d != nil; d = d.Parent() {
		var features *descriptorpb.FeatureSet
		switch opts := d.Options().(type) {
		case *descriptorpb.MessageOptions:
			features = opts.GetFeatures()
		case *descriptorpb.FileOptions:
			features = opts.GetFeatures()
		}
		if verify, set := javaUTF8Feature(features); set {
			return verify
		}
	}
	return false
}

// javaUTF8Feature reads the Java feature utf8_validation of a set of features. The extension can be
// a known field or an unknown field of the set, so the function reads the wire form.
func javaUTF8Feature(features *descriptorpb.FeatureSet) (verify, set bool) {
	if features == nil {
		return false, false
	}
	raw, err := proto.Marshal(features)
	if err != nil {
		return false, false
	}
	for len(raw) > 0 {
		number, wireType, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return false, false
		}
		raw = raw[n:]
		if number == javaFeaturesExtension && wireType == protowire.BytesType {
			inner, n := protowire.ConsumeBytes(raw)
			if n < 0 {
				return false, false
			}
			return readJavaUTF8(inner)
		}
		n = protowire.ConsumeFieldValue(number, wireType, raw)
		if n < 0 {
			return false, false
		}
		raw = raw[n:]
	}
	return false, false
}

func readJavaUTF8(raw []byte) (verify, set bool) {
	for len(raw) > 0 {
		number, wireType, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return false, false
		}
		raw = raw[n:]
		if number == javaUTF8ValidationTag && wireType == protowire.VarintType {
			value, n := protowire.ConsumeVarint(raw)
			if n < 0 {
				return false, false
			}
			return value == javaUTF8ValidationVerify, true
		}
		n = protowire.ConsumeFieldValue(number, wireType, raw)
		if n < 0 {
			return false, false
		}
		raw = raw[n:]
	}
	return false, false
}
