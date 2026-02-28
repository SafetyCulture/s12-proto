// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package codegen

import (
	"fmt"
	"strings"
)

// sprintf is an indirection for fmt.Sprintf that prevents the go vet printf
// checker from treating W as a printf wrapper. This is intentional:
// when W is called with zero variadic args, the string is written literally
// (percent signs are safe). WL delegates to W, inheriting this behavior.
// This differs from fmt.Fprintf where %% is always needed.
var sprintf = fmt.Sprintf

// CodeWriter wraps strings.Builder with convenience methods for code generation.
// The zero value is ready to use — no constructor is needed.
//
// W(s, args...) writes text: with zero args the string is written literally
// (percent signs are safe); with args, fmt.Sprintf is used.
//
// WL(args...) is purely variadic: no args writes a bare newline; one string arg
// is written literally with a trailing newline; two+ args uses fmt.Sprintf
// formatting with a trailing newline.
type CodeWriter struct {
	sb strings.Builder
}

// W writes formatted text without a trailing newline.
// With zero args, the format string is written literally (percent signs are safe).
// With args, fmt.Sprintf is used (standard %% escaping applies).
func (w *CodeWriter) W(s string, args ...any) {
	if len(args) > 0 {
		w.sb.WriteString(sprintf(s, args...))
	} else {
		w.sb.WriteString(s)
	}
}

// WL writes formatted text followed by a single newline.
// With no args, writes a bare newline (blank line separator).
// With one arg, the string is written literally (percent signs are safe).
// With two+ args, fmt.Sprintf processes the format string.
func (w *CodeWriter) WL(args ...any) {
	if len(args) == 0 {
		w.sb.WriteByte('\n')
		return
	}
	w.W(args[0].(string), args[1:]...)
	w.sb.WriteByte('\n')
}

// String returns the accumulated output.
func (w *CodeWriter) String() string {
	return w.sb.String()
}

// NewWriter returns a new CodeWriter ready for use.
func NewWriter() *CodeWriter {
	return &CodeWriter{}
}
