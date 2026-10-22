// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package kotlin

import "strings"

// DefaultPackage is the root package of the generated classes.
const DefaultPackage = "com.safetyculture.s12kt"

// Options holds the plugin options that protoc passes in the request parameter.
type Options struct {
	// Package is the root package. The proto package follows it.
	Package string
}

// ParseOptions reads a comma-separated protoc parameter string. The option
// kotlin_package sets the root package, and a bare kotlin_package keeps
// DefaultPackage. Other keys are ignored.
func ParseOptions(parameter string) Options {
	opts := Options{Package: DefaultPackage}
	for _, item := range strings.Split(parameter, ",") {
		key, value, _ := strings.Cut(item, "=")
		if strings.TrimSpace(key) != "kotlin_package" {
			continue
		}
		if value = strings.TrimSpace(value); value != "" {
			opts.Package = value
		}
	}
	return opts
}
