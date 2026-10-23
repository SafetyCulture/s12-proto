// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import "strings"

// printer writes indented source text. Every print call ends its text with a
// newline, and a blank line gets no indentation.
type printer struct {
	sb          strings.Builder
	indentation string
	atLineStart bool
}

const singleIndent = "  "

func newPrinter() *printer {
	return &printer{atLineStart: true}
}

// sub returns an empty printer with the indentation of p.
func (p *printer) sub() *printer {
	return &printer{indentation: p.indentation, atLineStart: true}
}

func (p *printer) isEmpty() bool { return p.sb.Len() == 0 }

func (p *printer) String() string { return p.sb.String() }

// print writes each text followed by a newline. Without a text it writes a blank line.
func (p *printer) print(texts ...string) {
	if len(texts) == 0 {
		p.sb.WriteByte('\n')
		p.atLineStart = true
		return
	}
	for _, t := range texts {
		p.write(t)
		p.sb.WriteByte('\n')
		p.atLineStart = true
	}
}

// printNoNewline writes the text as it is.
func (p *printer) printNoNewline(text string) {
	p.write(text)
}

// printIndented writes the texts one level deeper.
func (p *printer) printIndented(texts ...string) {
	p.indent()
	p.print(texts...)
	p.outdent()
}

func (p *printer) write(text string) {
	for i := range len(text) {
		c := text[i]
		if p.atLineStart && c != '\n' {
			p.sb.WriteString(p.indentation)
		}
		p.sb.WriteByte(c)
		p.atLineStart = c == '\n'
	}
}

// appendPrinter adds the text of other. With indenting, the text takes the
// indentation of p.
func (p *printer) appendPrinter(other *printer, indenting bool) {
	if indenting {
		p.write(other.String())
		return
	}
	p.sb.WriteString(other.String())
	p.atLineStart = other.atLineStart
}

func (p *printer) indent() { p.indentation += singleIndent }

func (p *printer) outdent() {
	p.indentation = p.indentation[:len(p.indentation)-len(singleIndent)]
}

func (p *printer) withIndentation(body func()) {
	p.indent()
	body()
	p.outdent()
}
