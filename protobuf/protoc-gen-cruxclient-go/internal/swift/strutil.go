// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"strings"
	"unicode"
)

// splitPath splits a path into the directory with its trailing slash, the base name and the suffix.
func splitPath(pathname string) (dir, base, suffix string) {
	dirEnd := strings.LastIndexByte(pathname, '/') + 1
	dir = pathname[:dirEnd]
	filename := pathname[dirEnd:]
	suffixStart := strings.LastIndexByte(filename, '.')
	if suffixStart < 0 {
		return dir, filename, ""
	}
	return dir, filename[:suffixStart], filename[suffixStart:]
}

// escapedToDataLiteral converts escaped bytes into a Data literal.
func escapedToDataLiteral(s string) string {
	if s == "" {
		return "Data()"
	}
	var out strings.Builder
	out.WriteString("Data([")
	separator := ""
	emit := func(v int) {
		out.WriteString(separator)
		fmt.Fprintf(&out, "%d", v)
		separator = ", "
	}
	escape := false
	octal := 0
	accumulator := 0
	for i := range len(s) {
		c := s[i]
		switch {
		case octal > 0:
			if c < '0' || c > '7' {
				panic(fmt.Sprintf("non-octal byte %q in an octal escape", c))
			}
			accumulator = accumulator<<3 | int(c-'0')
			octal--
			if octal == 0 {
				emit(accumulator)
			}
		case escape:
			switch {
			case c == 'n':
				emit(10)
			case c == 'r':
				emit(13)
			case c == 't':
				emit(9)
			case c >= '0' && c < '8':
				octal = 2
				accumulator = int(c - '0')
			default:
				emit(int(c))
			}
			escape = false
		case c == '\\':
			escape = true
		default:
			emit(int(c))
		}
	}
	if escape || octal > 0 {
		panic("truncated escape")
	}
	out.WriteString("])")
	return out.String()
}

// stringToEscapedStringLiteral converts a string default value into a Swift string literal.
func stringToEscapedStringLiteral(s string) string {
	if s == "" {
		return "String()"
	}
	var out strings.Builder
	out.WriteByte('"')
	for _, c := range s {
		switch {
		case c == 0:
			out.WriteString("\\0")
		case c >= 1 && c < 32:
			fmt.Fprintf(&out, "\\u{%02x}", c)
		case c == '"':
			out.WriteString("\\\"")
		case c == '\\':
			out.WriteString("\\\\")
		default:
			out.WriteRune(c)
		}
	}
	out.WriteByte('"')
	return out.String()
}

func isSwiftNewline(r rune) bool {
	switch r {
	case '\n', '\v', '\f', '\r', 0x85, 0x2028, 0x2029:
		return true
	}
	return false
}

func isSwiftWhitespace(r rune) bool {
	return r == '\t' || unicode.Is(unicode.Zs, r)
}

// prefixLines adds the prefix to each line of the text and drops the blank lines at the end.
func prefixLines(text, prefix string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	var lines []string
	start := 0
	for i, r := range text {
		if isSwiftNewline(r) {
			lines = append(lines, text[start:i])
			start = i + len(string(r))
		}
	}
	lines = append(lines, text[start:])
	for len(lines) > 0 && strings.TrimFunc(lines[len(lines)-1], isSwiftWhitespace) == "" {
		lines = lines[:len(lines)-1]
	}
	var out strings.Builder
	for _, line := range lines {
		out.WriteString(prefix + line + "\n")
	}
	return out.String()
}
