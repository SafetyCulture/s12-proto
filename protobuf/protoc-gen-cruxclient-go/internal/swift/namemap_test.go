// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNameInstructionWriterRejectsDescendingNumbers(t *testing.T) {
	w := newNameInstructionWriter()
	w.writeSame(2, "second")
	w.writeSame(1, "first")

	require.EqualError(t, w.err, "field and case numbers must be written in ascending order")
}
