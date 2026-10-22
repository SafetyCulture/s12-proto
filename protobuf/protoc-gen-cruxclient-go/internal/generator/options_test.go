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

func TestParseKotlinPackage(t *testing.T) {
	tests := []struct {
		parameter string
		pkg       string
	}{
		{"", ""},
		{"swift_module", ""},
		{"kotlin_package", DefaultKotlinPackage},
		{"kotlin_package=", DefaultKotlinPackage},
		{"kotlin_package=com.example.kt", "com.example.kt"},
		{"swift_module, kotlin_package = com.example.kt ", "com.example.kt"},
	}
	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			opts := ParseOptions(tt.parameter)
			if opts.KotlinPackage != tt.pkg {
				t.Errorf("ParseOptions(%q).KotlinPackage = %q, want %q", tt.parameter, opts.KotlinPackage, tt.pkg)
			}
			if opts.KotlinEnabled() != (tt.pkg != "") {
				t.Errorf("ParseOptions(%q).KotlinEnabled() = %v", tt.parameter, opts.KotlinEnabled())
			}
		})
	}
}
