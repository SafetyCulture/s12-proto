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
	require.ErrorContains(t, err, "broken-public.proto: public dependency index 1 is out of range")
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
	require.ErrorContains(t, err, "broken-oneof.proto: field value has oneof index 1 out of range")
	require.NotContains(t, err.Error(), "goroutine")
}
