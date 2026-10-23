// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import "testing"

func TestParseOptions(t *testing.T) {
	tests := []struct {
		parameter string
		module    string
	}{
		{"", ""},
		{"swift_module", DefaultSwiftModule},
		{"swift_module=", DefaultSwiftModule},
		{"swift_module=Custom", "Custom"},
		{"paths=source_relative, swift_module=Custom", "Custom"},
		{" swift_module ", DefaultSwiftModule},
		{"swift_module=A,swift_module=B", "B"},
		{"swift_module=A,swift_module", DefaultSwiftModule},
		{"swift_module=A=B", "A=B"},
	}
	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			opts, err := ParseOptions(tt.parameter)
			if err != nil {
				t.Fatalf("ParseOptions(%q) returned %v", tt.parameter, err)
			}
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
	}
	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			opts, err := ParseOptions(tt.parameter)
			if err != nil {
				t.Fatalf("ParseOptions(%q) returned %v", tt.parameter, err)
			}
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
		{"kotlin_messages=a.b,kotlin_messages=c.d", "c.d", false},
		{"kotlin_messages=a.b,kotlin_messages", DefaultKotlinPackage, false},
		{"kotlin_messages, messages_only", DefaultKotlinPackage, true},
	}
	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			opts, err := ParseOptions(tt.parameter)
			if err != nil {
				t.Fatalf("ParseOptions(%q) returned %v", tt.parameter, err)
			}
			if opts.KotlinMessages != tt.root || opts.KotlinMessagesEnabled() != (tt.root != "") {
				t.Errorf("ParseOptions(%q).KotlinMessages = %q, want %q", tt.parameter, opts.KotlinMessages, tt.root)
			}
			if opts.MessagesOnly != tt.only {
				t.Errorf("ParseOptions(%q).MessagesOnly = %v, want %v", tt.parameter, opts.MessagesOnly, tt.only)
			}
		})
	}
}

func TestParseOptionsRejectsInvalidInput(t *testing.T) {
	for _, parameter := range []string{
		"unknown",
		"unknown=value",
		"=x",
		"messages_only=",
		"messages_only=false",
		"messages_only",
		"kotlin_package,messages_only",
		"swift_module,messages_only",
		"kotlin_package=com/example",
		"kotlin_package=com..example",
		"kotlin_package=class",
		"kotlin_messages=com.example.-name",
	} {
		t.Run(parameter, func(t *testing.T) {
			_, err := ParseOptions(parameter)
			if err == nil {
				t.Fatalf("ParseOptions(%q) returned nil error", parameter)
			}
		})
	}
}

func TestParseSwiftMessages(t *testing.T) {
	tests := []struct {
		parameter string
		module    string
	}{
		{"", ""},
		{"swift_module", ""},
		{"swift_messages", DefaultSwiftModule},
		{"swift_messages=", DefaultSwiftModule},
		{"swift_messages=Custom", "Custom"},
		{"Visibility=Public, swift_messages = Mod ", "Mod"},
	}
	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			opts, err := ParseOptions(tt.parameter)
			if err != nil {
				t.Fatalf("ParseOptions(%q) returned %v", tt.parameter, err)
			}
			if opts.SwiftMessages != tt.module {
				t.Errorf("ParseOptions(%q).SwiftMessages = %q, want %q", tt.parameter, opts.SwiftMessages, tt.module)
			}
			if opts.SwiftMessagesEnabled() != (tt.module != "") {
				t.Errorf("ParseOptions(%q).SwiftMessagesEnabled() = %v", tt.parameter, opts.SwiftMessagesEnabled())
			}
		})
	}
}

func TestParseOptionsRejectsDifferentKotlinRoots(t *testing.T) {
	_, err := ParseOptions("kotlin_package=com.example.kt,kotlin_messages=com.example.other")
	if err == nil {
		t.Fatal("ParseOptions returned nil error")
	}
	want := "cruxclient-go: kotlin_package=com.example.kt and kotlin_messages=com.example.other must name the same root package"
	if err.Error() != want {
		t.Errorf("ParseOptions error = %q, want %q", err, want)
	}
}

func TestParseOptionsNamesAnEmptyKey(t *testing.T) {
	_, err := ParseOptions("kotlin_package,=x")
	if err == nil || err.Error() != `cruxclient-go: empty option key in "=x"` {
		t.Fatalf("ParseOptions returned %v", err)
	}
}
