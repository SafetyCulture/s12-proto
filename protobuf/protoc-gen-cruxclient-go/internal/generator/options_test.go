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
		{"swift_module=A,swift_module=B", "B"},
		{"swift_module=A,swift_module", DefaultSwiftModule},
		{"swift_module=A=B", "A=B"},
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
		{"kotlin_package=a.b,kotlin_package=c.d", "c.d"},
		{"kotlin_package=a.b,kotlin_package", DefaultKotlinPackage},
		{"kotlin_package=a=b", "a=b"},
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

func TestParseKotlinMessages(t *testing.T) {
	tests := []struct {
		parameter string
		root      string
		only      bool
	}{
		{"", "", false},
		{"kotlin_package", "", false},
		{"kotlin_messages", DefaultKotlinPackage, false},
		{"kotlin_messages=", DefaultKotlinPackage, false},
		{"kotlin_messages=com.example.kt", "com.example.kt", false},
		{"kotlin_messages,kotlin_package=com.example.kt", "com.example.kt", false},
		{"kotlin_package=com.example.kt,kotlin_messages", "com.example.kt", false},
		{"kotlin_package=com.example.kt,kotlin_messages=com.example.other", "com.example.other", false},
		{"kotlin_messages=a.b,kotlin_messages=c.d", "c.d", false},
		{"kotlin_messages=a.b,kotlin_messages", DefaultKotlinPackage, false},
		{"kotlin_messages=a=b", "a=b", false},
		{"kotlin_messages, messages_only", DefaultKotlinPackage, true},
		{"messages_only", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			opts := ParseOptions(tt.parameter)
			if opts.KotlinMessages != tt.root || opts.KotlinMessagesEnabled() != (tt.root != "") {
				t.Errorf("ParseOptions(%q).KotlinMessages = %q, want %q", tt.parameter, opts.KotlinMessages, tt.root)
			}
			if opts.MessagesOnly != tt.only {
				t.Errorf("ParseOptions(%q).MessagesOnly = %v, want %v", tt.parameter, opts.MessagesOnly, tt.only)
			}
		})
	}
}
