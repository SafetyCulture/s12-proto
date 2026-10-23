// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import (
	"fmt"
	"strings"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/kotlin"
)

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

// ParseOptions reads a comma-separated protoc parameter string.
func ParseOptions(parameter string) (Options, error) {
	var opts Options
	var kotlinMessages bool
	for _, item := range strings.Split(parameter, ",") {
		key, value, hasValue := strings.Cut(item, "=")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			if item == "" {
				continue
			}
			return Options{}, fmt.Errorf("cruxclient-go: unknown option %q", key)
		}
		switch key {
		case "swift_module":
			opts.SwiftModule = valueOr(value, DefaultSwiftModule)
		case "kotlin_package":
			opts.KotlinPackage = valueOr(value, DefaultKotlinPackage)
		case "kotlin_messages":
			opts.KotlinMessages = value
			kotlinMessages = true
		case "messages_only":
			if hasValue {
				return Options{}, fmt.Errorf("cruxclient-go: option %q does not accept a value", key)
			}
			opts.MessagesOnly = true
		case "paths":
		default:
			return Options{}, fmt.Errorf("cruxclient-go: unknown option %q", key)
		}
	}
	if opts.MessagesOnly && !kotlinMessages {
		return Options{}, fmt.Errorf("cruxclient-go: messages_only requires kotlin_messages")
	}
	if opts.KotlinPackage != "" {
		if err := kotlin.ValidateRoot(opts.KotlinPackage); err != nil {
			return Options{}, err
		}
	}
	if kotlinMessages {
		opts.KotlinMessages = valueOr(opts.KotlinMessages, valueOr(opts.KotlinPackage, DefaultKotlinPackage))
		if err := kotlin.ValidateRoot(opts.KotlinMessages); err != nil {
			return Options{}, err
		}
		if opts.KotlinPackage != "" && opts.KotlinPackage != opts.KotlinMessages {
			return Options{}, fmt.Errorf(
				"kotlin_package=%s and kotlin_messages=%s must name the same root package",
				opts.KotlinPackage,
				opts.KotlinMessages,
			)
		}
	}
	return opts, nil
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
