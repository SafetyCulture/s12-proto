// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestGenerateFileInvalidSwiftPrefixError(t *testing.T) {
	prefix := "not-valid"
	g := &generator{}
	_, err := g.generateFile(&file{
		name: "example.proto",
		proto: &descriptorpb.FileDescriptorProto{
			Options: &descriptorpb.FileOptions{SwiftPrefix: proto.String(prefix)},
		},
	})
	require.EqualError(t, err, "example.proto has an 'swift_prefix' that isn't a valid Swift identifier (not-valid).")
}
