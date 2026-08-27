// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strconv"
	"strings"
)

// translatedPattern is a Go character-class pattern rewritten for .NET.
type translatedPattern struct {
	// Pattern is the .NET equivalent, anchored so it cannot match a prefix.
	Pattern string

	// AstralCategories names the general categories the pattern admits, as
	// UnicodeCategories flag names. .NET matches over UTF-16 code units, so a code
	// point above the Basic Multilingual Plane never matches the category token that
	// admits it; the runtime removes those before matching instead.
	AstralCategories []string

	// NarrowPattern reads each two-letter category token as the category it names,
	// and is empty when the class holds no such token. It admits a subset of
	// Pattern, and reports rather than decides: a value it turns away has
	// characters the field's annotation did not ask for.
	NarrowPattern string

	// NarrowAstralCategories is AstralCategories for NarrowPattern.
	NarrowAstralCategories []string
}

// twoLetterCategories are the Unicode general category names written with two
// letters. A token is only read as one of these when both letters spell a
// category, so a class member that merely follows a one-letter token stays a
// literal.
var twoLetterCategories = map[string]bool{
	"Lu": true, "Ll": true, "Lt": true, "Lm": true, "Lo": true,
	"Mn": true, "Mc": true, "Me": true,
	"Nd": true, "Nl": true, "No": true,
	"Pc": true, "Pd": true, "Ps": true, "Pe": true, "Pi": true, "Pf": true, "Po": true,
	"Sm": true, "Sc": true, "Sk": true, "So": true,
	"Zs": true, "Zl": true, "Zp": true,
	"Cc": true, "Cf": true, "Cs": true, "Co": true, "Cn": true,
}

// translatePattern rewrites a generated Go character-class pattern for .NET.
//
// Five differences make a literal copy wrong, and each one widens the pattern
// rather than narrowing it, so a literal copy would accept input Go rejects:
//
//   - Go writes a one-letter category as \pL, which .NET does not accept at all.
//   - Go writes a code point as \x{0020}, which .NET reads as x repeated 20 times.
//   - ^ and $ match either side of a newline in .NET, so a value can smuggle one in.
//   - RE2 reads \pXy as the one-letter class \pX followed by a literal y, so a
//     declared two-letter category is already the whole one-letter class in Go. That
//     is reproduced rather than corrected: the generated C# has to accept what the
//     Go validators accept, and narrowing it here would reject live traffic. The
//     reading the annotation asks for comes back as NarrowPattern, which reports
//     what it would have turned away without changing the verdict.
//   - .NET rejects an escape it does not define, such as \_, where RE2 allows any
//     punctuation to be escaped.
func translatePattern(goPattern string) (translatedPattern, error) {
	body, ok := strings.CutPrefix(goPattern, "^[")
	if !ok {
		return translatedPattern{}, fmt.Errorf("pattern does not open with ^[: %s", goPattern)
	}
	body, ok = strings.CutSuffix(body, "]+$")
	if !ok {
		return translatedPattern{}, fmt.Errorf("pattern does not close with ]+$: %s", goPattern)
	}

	var out, narrow strings.Builder
	seen := map[string]bool{}
	narrowSeen := map[string]bool{}
	var categories, narrowCategories []string
	widened := false

	collect := func(category string, into *[]string, seen map[string]bool) {
		if seen[category] {
			return
		}
		seen[category] = true
		*into = append(*into, category)
	}

	for i := 0; i < len(body); {
		switch {
		case strings.HasPrefix(body[i:], `\x{`):
			end := strings.IndexByte(body[i:], '}')
			if end < 0 {
				return translatedPattern{}, fmt.Errorf("unterminated \\x{ in %s", goPattern)
			}
			digits := body[i+3 : i+end]
			code, err := strconv.ParseUint(digits, 16, 32)
			if err != nil {
				return translatedPattern{}, fmt.Errorf("unparsable code point \\x{%s}: %w", digits, err)
			}
			if code > 0xFFFF {
				// A character class matches one UTF-16 code unit, so an astral code point
				// would have to be written as its surrogate pair, which would then match
				// either half on its own and admit malformed input.
				return translatedPattern{}, fmt.Errorf("code point U+%X is above the Basic Multilingual Plane and cannot be a character class member", code)
			}
			fmt.Fprintf(&out, `\u%04X`, code)
			fmt.Fprintf(&narrow, `\u%04X`, code)
			i += end + 1

		case strings.HasPrefix(body[i:], `\p`) && i+2 < len(body) && isCategoryLetter(body[i+2]):
			// RE2 takes exactly one letter here; anything after it is a literal.
			category := string(body[i+2])
			consumed := 3

			narrowCategory := category
			if i+3 < len(body) && isCategoryLetter(body[i+3]) && twoLetterCategories[category+string(body[i+3])] {
				narrowCategory = category + string(body[i+3])
				consumed = 4
				widened = true
				fmt.Fprintf(&out, `\p{%s}`, category)
				out.WriteString(escapeClassMember(rune(body[i+3])))
			} else {
				fmt.Fprintf(&out, `\p{%s}`, category)
			}
			fmt.Fprintf(&narrow, `\p{%s}`, narrowCategory)

			collect(category, &categories, seen)
			collect(narrowCategory, &narrowCategories, narrowSeen)
			i += consumed

		default:
			r, size := decodeRune(body[i:])
			out.WriteString(escapeClassMember(r))
			narrow.WriteString(escapeClassMember(r))
			i += size
		}
	}

	translated := translatedPattern{
		Pattern:          `\A[` + out.String() + `]+\z`,
		AstralCategories: categories,
	}
	if widened {
		translated.NarrowPattern = `\A[` + narrow.String() + `]+\z`
		translated.NarrowAstralCategories = narrowCategories
	}
	return translated, nil
}

func isCategoryLetter(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func decodeRune(s string) (rune, int) {
	for i, r := range s {
		if i == 0 {
			return r, len(string(r))
		}
	}
	return 0, 0
}

// escapeClassMember writes a literal so .NET reads it as itself inside a class.
func escapeClassMember(r rune) string {
	switch r {
	case '\\', ']', '^', '-', '[':
		return `\` + string(r)
	}
	if r < 0x20 || r > 0x7E {
		return fmt.Sprintf(`\u%04X`, r)
	}
	return string(r)
}

// dotNetDefinedEscapes are the escapes .NET gives the same meaning as RE2. RE2
// additionally permits a backslash before any punctuation, where .NET rejects an
// escape it does not define, so those lose the backslash instead.
var dotNetDefinedEscapes = map[byte]bool{
	'.': true, '+': true, '*': true, '?': true, '(': true, ')': true,
	'[': true, ']': true, '{': true, '}': true, '|': true, '^': true,
	'$': true, '\\': true, 't': true, 'n': true, 'r': true, 'f': true,
	'v': true, '0': true, 'a': true, 'e': true,
}

// asciiClasses rewrite the shorthand classes RE2 defines over ASCII. .NET reads
// the same shorthand over the whole of Unicode, so copying it through would admit
// input Go rejects.
var asciiClasses = map[byte]struct{ inside, outside string }{
	'd': {`0-9`, `[0-9]`},
	'D': {"", `[^0-9]`},
	'w': {`0-9A-Za-z_`, `[0-9A-Za-z_]`},
	'W': {"", `[^0-9A-Za-z_]`},
	's': {`\t\n\f\r `, `[\t\n\f\r ]`},
	'S': {"", `[^\t\n\f\r ]`},
}

// translateFreeRegex rewrites a hand-written (validator.regex) pattern for .NET.
//
// Anything whose .NET meaning cannot be shown to match RE2's is refused rather
// than approximated, because every difference between the two engines widens the
// pattern, and a widened pattern accepts input the Go validators reject.
func translateFreeRegex(goPattern string) (string, error) {
	var out strings.Builder
	inClass := false

	for i := 0; i < len(goPattern); i++ {
		c := goPattern[i]

		if c == '\\' {
			if i+1 >= len(goPattern) {
				return "", fmt.Errorf("pattern ends in a backslash")
			}
			i++
			next := goPattern[i]

			if class, ok := asciiClasses[next]; ok {
				if inClass {
					if class.inside == "" {
						return "", fmt.Errorf(`\%c inside a character class has no .NET equivalent that stays ASCII`, next)
					}
					out.WriteString(class.inside)
				} else {
					out.WriteString(class.outside)
				}
				continue
			}

			switch next {
			case 'b', 'B':
				// RE2 decides a word boundary over ASCII and .NET over Unicode.
				return "", fmt.Errorf(`\%c is ASCII-only in RE2 and Unicode-aware in .NET`, next)
			case 'A', 'z', 'Z':
				out.WriteByte('\\')
				out.WriteByte(next)
				continue
			case 'p', 'P':
				return "", fmt.Errorf(`\%c is read as a one-letter category by RE2 and needs braces in .NET; rewrite the pattern`, next)
			case 'x', 'u':
				return "", fmt.Errorf(`\%c escapes differ between the engines; write the character itself`, next)
			}

			if next >= '1' && next <= '9' {
				return "", fmt.Errorf(`\%c is a backreference, which the linear-time engine does not support`, next)
			}

			if dotNetDefinedEscapes[next] {
				out.WriteByte('\\')
				out.WriteByte(next)
				continue
			}

			if isCategoryLetter(next) {
				return "", fmt.Errorf(`\%c is not an escape .NET defines`, next)
			}

			// RE2 allows a backslash before any punctuation; .NET rejects the ones it
			// does not define, so the character stands for itself.
			out.WriteByte(next)
			continue
		}

		switch c {
		case '[':
			if !inClass {
				if end, total := totalClass(goPattern[i:]); total {
					// A class holding a shorthand and its complement matches every character
					// in either engine, wherever each one puts the boundary between them.
					out.WriteString(`[\s\S]`)
					i += end
					continue
				}
				inClass = true
			}
			out.WriteByte(c)
		case ']':
			inClass = false
			out.WriteByte(c)
		case '^':
			if inClass {
				out.WriteByte(c)
				continue
			}
			// Equivalent to .NET's ^ while Multiline is off, but stated outright so the
			// option cannot change what the pattern means.
			out.WriteString(`\A`)
		case '$':
			if inClass {
				out.WriteByte(c)
				continue
			}
			// .NET's $ also matches before a trailing newline, so a value could end in
			// one and still pass.
			out.WriteString(`\z`)
		case '(':
			if strings.HasPrefix(goPattern[i:], "(?=") || strings.HasPrefix(goPattern[i:], "(?!") ||
				strings.HasPrefix(goPattern[i:], "(?<") {
				return "", fmt.Errorf("lookaround is not supported by the linear-time engine")
			}
			out.WriteByte(c)
		default:
			out.WriteByte(c)
		}
	}

	return out.String(), nil
}

// totalClass reports whether the character class opening at the start of s holds
// a shorthand together with its complement, and where the class ends. Such a
// class admits every character, so the difference between RE2's ASCII shorthands
// and .NET's Unicode ones cannot change what it matches.
func totalClass(s string) (int, bool) {
	if !strings.HasPrefix(s, "[") {
		return 0, false
	}

	var body strings.Builder
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			if i+1 >= len(s) {
				return 0, false
			}
			body.WriteByte(s[i+1])
			i++
		case ']':
			if i == 1 {
				// A closing bracket in first position is a member, not the terminator.
				body.WriteByte(']')
				continue
			}
			members := body.String()
			for _, pair := range []struct{ lower, upper byte }{{'s', 'S'}, {'d', 'D'}, {'w', 'W'}} {
				if strings.IndexByte(members, pair.lower) >= 0 && strings.IndexByte(members, pair.upper) >= 0 {
					return i, true
				}
			}
			return 0, false
		}
	}
	return 0, false
}
