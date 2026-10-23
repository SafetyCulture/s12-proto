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
	// KotlinMessages enables the pure Kotlin message classes and sets their root package. An empty value leaves them off.
	KotlinMessages string
	// MessagesOnly suppresses every cruxclient output, so that a run writes only the message classes.
	MessagesOnly bool
}

// ParseOptions reads a comma-separated protoc parameter string. A bare
// swift_module enables the Swift output with DefaultSwiftModule. A bare
// kotlin_package enables the Kotlin output with DefaultKotlinPackage. A bare
// kotlin_messages enables the Kotlin message classes with the kotlin_package
// value, or DefaultKotlinPackage if kotlin_package is absent. Other keys are ignored.
func ParseOptions(parameter string) Options {
	var opts Options
	var kotlinMessages bool
	for _, item := range strings.Split(parameter, ",") {
		key, value, _ := strings.Cut(item, "=")
		value = strings.TrimSpace(value)
		switch strings.TrimSpace(key) {
		case "swift_module":
			opts.SwiftModule = valueOr(value, DefaultSwiftModule)
		case "kotlin_package":
			opts.KotlinPackage = valueOr(value, DefaultKotlinPackage)
		case "kotlin_messages":
			opts.KotlinMessages = value
			kotlinMessages = true
		case "messages_only":
			opts.MessagesOnly = true
		}
	}
	if kotlinMessages {
		opts.KotlinMessages = valueOr(opts.KotlinMessages, valueOr(opts.KotlinPackage, DefaultKotlinPackage))
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

// KotlinMessagesEnabled reports whether the Kotlin message classes are on.
func (o Options) KotlinMessagesEnabled() bool {
	return o.KotlinMessages != ""
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
