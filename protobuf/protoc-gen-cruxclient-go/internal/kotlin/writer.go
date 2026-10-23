// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package kotlin

import (
	"fmt"
	"strings"
)

// writer builds source text with a four-space indent.
type writer struct {
	b      strings.Builder
	indent int
}

// line writes one line at the current indent. An empty format writes an empty line.
func (w *writer) line(format string, args ...any) {
	if format == "" {
		w.b.WriteByte('\n')
		return
	}
	w.b.WriteString(strings.Repeat("    ", w.indent))
	if len(args) == 0 {
		w.b.WriteString(format)
	} else {
		fmt.Fprintf(&w.b, format, args...)
	}
	w.b.WriteByte('\n')
}

// open writes a line and indents the lines that follow.
func (w *writer) open(format string, args ...any) {
	w.line(format, args...)
	w.indent++
}

// close ends the indented block with a line.
func (w *writer) close(format string, args ...any) {
	w.indent--
	w.line(format, args...)
}

// doc writes a KDoc comment. The comment text of a proto uses "//" markers already removed.
func (w *writer) doc(comment string) {
	comment = strings.ReplaceAll(comment, "\r", "")
	comment = strings.TrimRight(comment, "\n ")
	if strings.TrimSpace(comment) == "" {
		return
	}
	// Kotlin comments nest, so a comment opener inside the text must not stay one.
	comment = strings.ReplaceAll(comment, "/*", "/&#42;")
	comment = strings.ReplaceAll(comment, "*/", "*&#47;")
	w.line("/**")
	for _, l := range strings.Split(comment, "\n") {
		l = strings.TrimPrefix(l, " ")
		if l == "" {
			w.line(" *")
		} else {
			w.line(" * %s", strings.TrimRight(l, " "))
		}
	}
	w.line(" */")
}

func (w *writer) String() string { return w.b.String() }

// text writes a line that has no format verbs.
func (w *writer) text(s string) {
	w.b.WriteString(strings.Repeat("    ", w.indent))
	w.b.WriteString(s)
	w.b.WriteByte('\n')
}
