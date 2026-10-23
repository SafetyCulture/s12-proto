// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestBuildFilesRejectsInvalidIndexes(t *testing.T) {
	_, err := buildFiles([]*descriptorpb.FileDescriptorProto{
		{Name: proto.String("dependency.proto"), Syntax: proto.String("proto3")},
		{
			Name:             proto.String("broken-public.proto"),
			Syntax:           proto.String("proto3"),
			Dependency:       []string{"dependency.proto"},
			PublicDependency: []int32{1},
		},
	})
	require.EqualError(t, err, "add file broken-public.proto: public dependency index 1 is out of range")
	require.NotContains(t, err.Error(), "goroutine")

	_, err = buildFiles([]*descriptorpb.FileDescriptorProto{{
		Name:   proto.String("broken-oneof.proto"),
		Syntax: proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{{
			Name:      proto.String("Message"),
			OneofDecl: []*descriptorpb.OneofDescriptorProto{{Name: proto.String("choice")}},
			Field: []*descriptorpb.FieldDescriptorProto{{
				Name:       proto.String("value"),
				Number:     proto.Int32(1),
				Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:       descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				OneofIndex: proto.Int32(1),
			}},
		}},
	}})
	require.EqualError(t, err, "add file broken-oneof.proto: add message Message: field value has oneof index 1 out of range")
	require.NotContains(t, err.Error(), "goroutine")
}

func TestBuildFilesAddsBindingContext(t *testing.T) {
	_, err := buildFiles([]*descriptorpb.FileDescriptorProto{{
		Name:    proto.String("broken-binding.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("example"),
		MessageType: []*descriptorpb.DescriptorProto{{
			Name: proto.String("Message"),
			Field: []*descriptorpb.FieldDescriptorProto{{
				Name:     proto.String("value"),
				Number:   proto.Int32(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String(".example.Missing"),
			}},
		}},
	}})
	require.ErrorContains(t, err, "add file broken-binding.proto: bind message example.Message: bind field example.Message.value")
	require.ErrorContains(t, err, "unknown message .example.Missing")
}

func TestBuildFilesRejectsUnsupportedEdition(t *testing.T) {
	_, err := buildFiles([]*descriptorpb.FileDescriptorProto{{
		Name:    proto.String("edition-2024.proto"),
		Edition: descriptorpb.Edition_EDITION_2024.Enum(),
	}})
	require.EqualError(t, err, "add file edition-2024.proto: edition EDITION_2024 is not supported")
}

func TestBuildFilesRejectsMissingImport(t *testing.T) {
	_, err := buildFiles([]*descriptorpb.FileDescriptorProto{{
		Name:       proto.String("imp.proto"),
		Syntax:     proto.String("proto3"),
		Dependency: []string{"absent.proto"},
	}})
	require.EqualError(t, err, "add file imp.proto: the import absent.proto has no descriptor")
}
