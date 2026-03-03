// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package codegen

import (
	"strings"
	"testing"
)

func TestWL_LiteralStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", "hello\n"},
		{"empty string", "", "\n"},
		{"percent literal safe", "100% complete", "100% complete\n"},
		{"double percent literal", "%%s stays", "%%s stays\n"},
		{"unicode content", "namespace \u00e4 {", "namespace \u00e4 {\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w CodeWriter
			w.WL(tt.input)
			if got := w.String(); got != tt.expected {
				t.Errorf("WL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestWL_FormatStrings(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{"one string arg", "hello %s", []any{"world"}, "hello world\n"},
		{"one int arg", "count: %d", []any{42}, "count: 42\n"},
		{"two args", "%s = %d", []any{"x", 42}, "x = 42\n"},
		{"three args", "%s(%s, %s)", []any{"f", "a", "b"}, "f(a, b)\n"},
		{"escaped percent with args", "100%% done %s", []any{"now"}, "100% done now\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w CodeWriter
			args := append([]any{tt.format}, tt.args...)
			w.WL(args...)
			if got := w.String(); got != tt.expected {
				t.Errorf("WL(%q, %v) = %q, want %q", tt.format, tt.args, got, tt.expected)
			}
		})
	}
}

func TestWL_ExtraArgs(t *testing.T) {
	var w CodeWriter
	w.WL("hello %s", "world", "extra")
	got := w.String()
	if !strings.Contains(got, "hello world") {
		t.Errorf("expected output to contain %q, got %q", "hello world", got)
	}
	if !strings.Contains(got, "EXTRA") {
		t.Errorf("expected output to contain EXTRA annotation, got %q", got)
	}
}

func TestW(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{"simple", "hello", nil, "hello"},
		{"with args", "x=%d", []any{5}, "x=5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w CodeWriter
			w.W(tt.format, tt.args...)
			if got := w.String(); got != tt.expected {
				t.Errorf("W(%q, %v) = %q, want %q", tt.format, tt.args, got, tt.expected)
			}
		})
	}
}

func TestWL_NoArgs(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{"single bare WL", 1, "\n"},
		{"double bare WL", 2, "\n\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w CodeWriter
			for i := 0; i < tt.count; i++ {
				w.WL()
			}
			if got := w.String(); got != tt.expected {
				t.Errorf("WL() x%d = %q, want %q", tt.count, got, tt.expected)
			}
		})
	}
}

func TestDoubleNewlinePattern(t *testing.T) {
	tests := []struct {
		name     string
		calls    func(w *CodeWriter)
		expected string
	}{
		{
			"WL then bare WL",
			func(w *CodeWriter) { w.WL("};"); w.WL() },
			"};\n\n",
		},
		{
			"bare WL then WL",
			func(w *CodeWriter) { w.WL(); w.WL("namespace %s {", "foo") },
			"\nnamespace foo {\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w CodeWriter
			tt.calls(&w)
			if got := w.String(); got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestPartialLineBuilding(t *testing.T) {
	tests := []struct {
		name     string
		calls    func(w *CodeWriter)
		expected string
	}{
		{
			"W then WL",
			func(w *CodeWriter) { w.W("hello "); w.WL("world") },
			"hello world\n",
		},
		{
			"multiple W then WL",
			func(w *CodeWriter) { w.W("a"); w.W("b"); w.WL("c") },
			"abc\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w CodeWriter
			tt.calls(&w)
			if got := w.String(); got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestEmptyWriter(t *testing.T) {
	var w CodeWriter
	if got := w.String(); got != "" {
		t.Errorf("empty CodeWriter.String() = %q, want %q", got, "")
	}
}

func TestMultipleLines(t *testing.T) {
	var w CodeWriter
	w.WL("namespace %s {", "foo")
	w.WL()
	w.WL("class Bar {")
	w.WL("};")
	w.WL()

	expected := "namespace foo {\n\nclass Bar {\n};\n\n"
	if got := w.String(); got != expected {
		t.Errorf("multi-line output = %q, want %q", got, expected)
	}
}
