// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnknownFieldTypePanicsWithFieldName(t *testing.T) {
	g := &generator{namer: &namer{}}
	f := &field{fullName: "example.Message.value", typ: fieldType(-1)}
	message := "example.Message.value: unknown field type -1"

	require.PanicsWithValue(t, message, func() { g.swiftType(f) })
	require.PanicsWithValue(t, message, func() { g.traitsType(f) })
}
