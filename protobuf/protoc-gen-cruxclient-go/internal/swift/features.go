// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// defaultFeatures returns the feature values that an edition gives to a file.
func defaultFeatures(edition descriptorpb.Edition) *descriptorpb.FeatureSet {
	fs := &descriptorpb.FeatureSet{
		FieldPresence:         descriptorpb.FeatureSet_EXPLICIT.Enum(),
		EnumType:              descriptorpb.FeatureSet_OPEN.Enum(),
		RepeatedFieldEncoding: descriptorpb.FeatureSet_PACKED.Enum(),
		Utf8Validation:        descriptorpb.FeatureSet_VERIFY.Enum(),
		MessageEncoding:       descriptorpb.FeatureSet_LENGTH_PREFIXED.Enum(),
		JsonFormat:            descriptorpb.FeatureSet_ALLOW.Enum(),
	}
	switch edition {
	case descriptorpb.Edition_EDITION_PROTO2:
		fs.EnumType = descriptorpb.FeatureSet_CLOSED.Enum()
		fs.RepeatedFieldEncoding = descriptorpb.FeatureSet_EXPANDED.Enum()
		fs.Utf8Validation = descriptorpb.FeatureSet_NONE.Enum()
		fs.JsonFormat = descriptorpb.FeatureSet_LEGACY_BEST_EFFORT.Enum()
	case descriptorpb.Edition_EDITION_PROTO3:
		fs.FieldPresence = descriptorpb.FeatureSet_IMPLICIT.Enum()
	}
	return fs
}

// resolveFeatures applies the features of a descriptor on top of the features it inherits.
func resolveFeatures(parent, own *descriptorpb.FeatureSet) *descriptorpb.FeatureSet {
	if own == nil {
		return parent
	}
	result := proto.Clone(parent).(*descriptorpb.FeatureSet)
	proto.Merge(result, own)
	return result
}

// resolveFieldFeatures also infers the features that proto2 and proto3 express with labels and options.
func resolveFieldFeatures(edition descriptorpb.Edition, fp *descriptorpb.FieldDescriptorProto, parent *descriptorpb.FeatureSet) *descriptorpb.FeatureSet {
	if edition >= descriptorpb.Edition_EDITION_2023 {
		return resolveFeatures(parent, fp.GetOptions().GetFeatures())
	}
	inferred := &descriptorpb.FeatureSet{}
	if fp.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REQUIRED {
		inferred.FieldPresence = descriptorpb.FeatureSet_LEGACY_REQUIRED.Enum()
	}
	if fp.GetType() == descriptorpb.FieldDescriptorProto_TYPE_GROUP {
		inferred.MessageEncoding = descriptorpb.FeatureSet_DELIMITED.Enum()
	}
	options := fp.GetOptions()
	if options.GetPacked() {
		inferred.RepeatedFieldEncoding = descriptorpb.FeatureSet_PACKED.Enum()
	}
	if edition == descriptorpb.Edition_EDITION_PROTO3 && options != nil && options.Packed != nil && !options.GetPacked() {
		inferred.RepeatedFieldEncoding = descriptorpb.FeatureSet_EXPANDED.Enum()
	}
	return resolveFeatures(defaultFeatures(edition), inferred)
}
