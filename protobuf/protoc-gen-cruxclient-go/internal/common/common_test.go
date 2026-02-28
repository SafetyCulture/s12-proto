// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package common

import (
	"testing"
)

func TestStripProto(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"routeguide/v1/route_guide.proto", "routeguide/v1/route_guide"},
		{"test.protodevel", "test"},
		{"no_extension", "no_extension"},
		{"path/to/file.proto", "path/to/file"},
		{"file.protodevel", "file"},
		// .protodevel takes precedence over .proto
		{"file.protodevel.proto", "file.protodevel"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := StripProto(tt.input)
			if got != tt.expected {
				t.Errorf("StripProto(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDotsToColons(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"routeguide.v1", "routeguide::v1"},
		{"routeguide.v1.RouteGuide", "routeguide::v1::RouteGuide"},
		{"no_dots", "no_dots"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := DotsToColons(tt.input)
			if got != tt.expected {
				t.Errorf("DotsToColons(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDotsToUnderscores(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"routeguide.v1.RouteGuide", "routeguide_v1_RouteGuide"},
		{"no_dots", "no_dots"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := DotsToUnderscores(tt.input)
			if got != tt.expected {
				t.Errorf("DotsToUnderscores(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDotsToSlashs(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"routeguide.v1", "routeguide/v1"},
		{"a.b.c", "a/b/c"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := DotsToSlashs(tt.input)
			if got != tt.expected {
				t.Errorf("DotsToSlashs(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUnderscoresToDollar(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"RouteSummary_Details", "RouteSummary$Details"},
		{"no_underscores", "no$underscores"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := UnderscoresToDollar(tt.input)
			if got != tt.expected {
				t.Errorf("UnderscoresToDollar(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUnderscoresToDots(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"RouteSummary_Details", "RouteSummary.Details"},
		{"no_underscores", "no.underscores"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := UnderscoresToDots(tt.input)
			if got != tt.expected {
				t.Errorf("UnderscoresToDots(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestReplaceCharacters(t *testing.T) {
	tests := []struct {
		input    string
		remove   string
		with     string
		expected string
	}{
		{"a.b.c", ".", "/", "a/b/c"},
		{"a-b_c", "-_", ".", "a.b.c"},
		{"no_match", "xyz", "!", "no_match"},
		{"", ".", "/", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ReplaceCharacters(tt.input, tt.remove, tt.with)
			if got != tt.expected {
				t.Errorf("ReplaceCharacters(%q, %q, %q) = %q, want %q", tt.input, tt.remove, tt.with, got, tt.expected)
			}
		})
	}
}

func TestFilenameIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// C++ isalnum only includes [a-zA-Z0-9]; '_' (0x5f) is encoded as _5f,
		// '.' (0x2e) is encoded as _2e, '/' (0x2f) is encoded as _2f.
		{"route_guide.proto", "route_5fguide_2eproto"},
		{"a/b.c", "a_2fb_2ec"},
		{"abc123", "abc123"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := FilenameIdentifier(tt.input)
			if got != tt.expected {
				t.Errorf("FilenameIdentifier(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestStringReplace(t *testing.T) {
	tests := []struct {
		str        string
		from       string
		to         string
		replaceAll bool
		expected   string
	}{
		{"aababab", "ab", "X", true, "aXXX"},
		{"aababab", "ab", "X", false, "aXabab"},
		{"hello world", "o", "0", true, "hell0 w0rld"},
		{"hello world", "o", "0", false, "hell0 world"},
	}
	for _, tt := range tests {
		got := StringReplace(tt.str, tt.from, tt.to, tt.replaceAll)
		if got != tt.expected {
			t.Errorf("StringReplace(%q, %q, %q, %v) = %q, want %q",
				tt.str, tt.from, tt.to, tt.replaceAll, got, tt.expected)
		}
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input      string
		delimiters string
		expected   []string
	}{
		{"routeguide/v1/route_guide", "/", []string{"routeguide", "v1", "route_guide"}},
		{"a//b", "/", []string{"a", "", "b"}},
		{"", "/", []string{""}},
		{"a", "/", []string{"a"}},
		{"a.b.c", ".", []string{"a", "b", "c"}},
		{"a", ".", []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.input+"|"+tt.delimiters, func(t *testing.T) {
			got := Tokenize(tt.input, tt.delimiters)
			if len(got) != len(tt.expected) {
				t.Errorf("Tokenize(%q, %q) = %v (len %d), want %v (len %d)",
					tt.input, tt.delimiters, got, len(got), tt.expected, len(tt.expected))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Tokenize(%q, %q)[%d] = %q, want %q",
						tt.input, tt.delimiters, i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"route_guide", "RouteGuide"},
		{"get_feature", "GetFeature"},
		{"message", "Message"},
		{"route_guide_v1", "RouteGuideV1"},
		{"", ""},
		{"a", "A"},
		{"_a", "A"},
		{"a_b_c", "ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			if got != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"RouteGuide", "routeguide"},
		{"UPPER", "upper"},
		{"already_lower", "already_lower"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToLower(tt.input)
			if got != tt.expected {
				t.Errorf("ToLower(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		elements  []string
		separator string
		expected  string
	}{
		{[]string{"a", "b", "c"}, "/", "a/b/c"},
		{[]string{"a"}, ",", "a"},
		{[]string{}, ",", ""},
	}
	for _, tt := range tests {
		got := Join(tt.elements, tt.separator)
		if got != tt.expected {
			t.Errorf("Join(%v, %q) = %q, want %q", tt.elements, tt.separator, got, tt.expected)
		}
	}
}
