// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestResolveFeaturesCopiesParent(t *testing.T) {
	parent := &descriptorpb.FeatureSet{EnumType: descriptorpb.FeatureSet_OPEN.Enum()}
	child := resolveFeatures(parent, nil)

	require.NotSame(t, parent, child)
	child.EnumType = descriptorpb.FeatureSet_CLOSED.Enum()
	require.Equal(t, descriptorpb.FeatureSet_OPEN, parent.GetEnumType())
}
