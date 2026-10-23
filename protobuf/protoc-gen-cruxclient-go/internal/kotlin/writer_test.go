// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package kotlin

import "testing"

func TestWriterDocRemovesCarriageReturns(t *testing.T) {
	var w writer
	w.doc("a\r\nb\r\n")
	if got, want := w.String(), "/**\n * a\n * b\n */\n"; got != want {
		t.Errorf("doc = %q, want %q", got, want)
	}
}
