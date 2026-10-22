// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import "testing"

func TestParseOptions(t *testing.T) {
	tests := []struct {
		parameter string
		module    string
	}{
		{"", ""},
		{"unrelated=1", ""},
		{"swift_module", DefaultSwiftModule},
		{"swift_module=", DefaultSwiftModule},
		{"swift_module=Custom", "Custom"},
		{"paths=source_relative, swift_module=Custom", "Custom"},
		{"foo=1, swift_module = Mod ", "Mod"},
		{" swift_module ", DefaultSwiftModule},
	}
	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			opts := ParseOptions(tt.parameter)
			if opts.SwiftModule != tt.module {
				t.Errorf("ParseOptions(%q).SwiftModule = %q, want %q", tt.parameter, opts.SwiftModule, tt.module)
			}
			if opts.SwiftEnabled() != (tt.module != "") {
				t.Errorf("ParseOptions(%q).SwiftEnabled() = %v", tt.parameter, opts.SwiftEnabled())
			}
		})
	}
}
