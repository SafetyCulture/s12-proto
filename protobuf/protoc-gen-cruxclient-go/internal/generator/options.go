// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import "strings"

// DefaultSwiftModule is the Swift module that holds the SwiftProtobuf types.
const DefaultSwiftModule = "S12Protos"

// DefaultKotlinPackage is the root package of the pure Kotlin messages.
const DefaultKotlinPackage = "com.safetyculture.s12kt"

// Options holds the plugin options that protoc passes in the request parameter.
type Options struct {
	// SwiftModule enables the Swift output. An empty value leaves it off.
	SwiftModule string
	// KotlinPackage enables the output for the pure Kotlin messages. An empty value leaves it off.
	KotlinPackage string
}

// ParseOptions reads a comma-separated protoc parameter string. A bare
// swift_module enables the Swift output with DefaultSwiftModule. A bare
// kotlin_package enables the Kotlin output with DefaultKotlinPackage. Other keys are ignored.
func ParseOptions(parameter string) Options {
	var opts Options
	for _, item := range strings.Split(parameter, ",") {
		key, value, _ := strings.Cut(item, "=")
		value = strings.TrimSpace(value)
		switch strings.TrimSpace(key) {
		case "swift_module":
			opts.SwiftModule = valueOr(value, DefaultSwiftModule)
		case "kotlin_package":
			opts.KotlinPackage = valueOr(value, DefaultKotlinPackage)
		}
	}
	return opts
}

// SwiftEnabled reports whether the Swift output is on.
func (o Options) SwiftEnabled() bool {
	return o.SwiftModule != ""
}

// KotlinEnabled reports whether the Kotlin output is on.
func (o Options) KotlinEnabled() bool {
	return o.KotlinPackage != ""
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
