// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitPath(t *testing.T) {
	for _, test := range []struct {
		path, dir, base, suffix string
	}{
		{"file.proto", "", "file", ".proto"},
		{"dir/file", "dir/", "file", ""},
		{"dir/.file", "dir/", "", ".file"},
		{"dir/file.", "dir/", "file", "."},
	} {
		dir, base, suffix := splitPath(test.path)
		require.Equal(t, test.dir, dir)
		require.Equal(t, test.base, base)
		require.Equal(t, test.suffix, suffix)
	}
}

func TestPrinterIndentation(t *testing.T) {
	p := newPrinter()
	p.print("top")
	p.indent()
	p.print("inner\nsecond")
	p.outdent()

	require.Equal(t, "top\n  inner\n  second\n", p.String())
}
