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
	// SwiftMessages enables the Swift messages. The option takes no value.
	SwiftMessages bool
	// KotlinPackage enables the output for the pure Kotlin messages. An empty value leaves it off.
	KotlinPackage string
	// KotlinMessages enables the pure Kotlin message classes and sets their root package. An empty value leaves them off.
	KotlinMessages string
	// MessagesOnly suppresses every cruxclient output, so that a run writes only the message classes.
	MessagesOnly bool
}

type parsedOptions struct {
	Options
	kotlinMessages bool
}

type optionSetter struct {
	set          func(p *parsedOptions, value string)
	acceptsValue bool
}

var optionSetters = map[string]optionSetter{
	"swift_module":   {set: func(p *parsedOptions, value string) { p.SwiftModule = valueOr(value, DefaultSwiftModule) }, acceptsValue: true},
	"swift_messages": {set: func(p *parsedOptions, _ string) { p.SwiftMessages = true }},
	"kotlin_package": {set: func(p *parsedOptions, value string) { p.KotlinPackage = valueOr(value, DefaultKotlinPackage) }, acceptsValue: true},
	"kotlin_messages": {set: func(p *parsedOptions, value string) {
		p.KotlinMessages = value
		p.kotlinMessages = true
	}, acceptsValue: true},
	"messages_only": {set: func(p *parsedOptions, _ string) { p.MessagesOnly = true }},
}

var swiftOptionKeys = map[string]bool{
	"EnumGeneration":                           true,
	"ExperimentalHiddenNames":                  true,
	"FileNaming":                               true,
	"ImplementationOnlyImports":                true,
	"ProtoPathModuleMappings":                  true,
	"SwiftProtobufModuleName":                  true,
	"UseAccessLevelOnImports":                  true,
	"Visibility":                               true,
	"experimental_strip_nonfunctional_codegen": true,
}

// SwiftOptionKeys returns the supported Swift generator option keys.
func SwiftOptionKeys() []string {
	keys := make([]string, 0, len(swiftOptionKeys))
	for key := range swiftOptionKeys {
		keys = append(keys, key)
	}
	return keys
}

// IsOptionKey reports whether key is a plugin option.
func IsOptionKey(key string) bool {
	_, ok := optionSetters[key]
	return ok || key == "paths"
}

// ParseOptions reads a comma-separated protoc parameter string.
func ParseOptions(parameter string) (Options, error) {
	var parsed parsedOptions
	var firstSwiftOption string
	for _, item := range strings.Split(parameter, ",") {
		key, value, hasValue := strings.Cut(item, "=")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			if item == "" {
				continue
			}
			return Options{}, fmt.Errorf("cruxclient-go: empty option key in %q", item)
		}
		if setter, ok := optionSetters[key]; ok {
			if hasValue && !setter.acceptsValue {
				return Options{}, fmt.Errorf("cruxclient-go: option %q does not accept a value", key)
			}
			setter.set(&parsed, value)
			continue
		}
		if key == "paths" {
			continue
		}
		if swiftOptionKeys[key] {
			if firstSwiftOption == "" {
				firstSwiftOption = key
			}
			continue
		}
		return Options{}, fmt.Errorf("cruxclient-go: unknown option %q", key)
	}
	if firstSwiftOption != "" && !parsed.SwiftMessages {
		return Options{}, fmt.Errorf("cruxclient-go: option %q requires swift_messages", firstSwiftOption)
	}
	if parsed.MessagesOnly && !parsed.kotlinMessages && !parsed.SwiftMessages {
		return Options{}, fmt.Errorf("cruxclient-go: messages_only requires swift_messages or kotlin_messages")
	}
	if parsed.KotlinPackage != "" {
		if err := kotlin.ValidateRoot(parsed.KotlinPackage); err != nil {
			return Options{}, fmt.Errorf("cruxclient-go: %w", err)
		}
	}
	opts := parsed.Options
	if parsed.kotlinMessages {
		opts.KotlinMessages = valueOr(opts.KotlinMessages, valueOr(opts.KotlinPackage, DefaultKotlinPackage))
		if err := kotlin.ValidateRoot(opts.KotlinMessages); err != nil {
			return Options{}, fmt.Errorf("cruxclient-go: %w", err)
		}
		if opts.KotlinPackage != "" && opts.KotlinPackage != opts.KotlinMessages {
			return Options{}, fmt.Errorf(
				"cruxclient-go: kotlin_package=%s and kotlin_messages=%s must name the same root package",
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

// SwiftMessagesEnabled reports whether the Swift messages are on.
func (o Options) SwiftMessagesEnabled() bool {
	return o.SwiftMessages
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
