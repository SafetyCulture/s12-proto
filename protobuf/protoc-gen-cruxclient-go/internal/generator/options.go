// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import "strings"

// DefaultSwiftModule is the Swift module that holds the SwiftProtobuf types.
const DefaultSwiftModule = "S12Protos"

// Options holds the plugin options that protoc passes in the request parameter.
type Options struct {
	// SwiftModule enables the Swift output. An empty value leaves it off.
	SwiftModule string
}

// ParseOptions reads a comma-separated protoc parameter string. A bare
// swift_module enables the Swift output with DefaultSwiftModule. Other keys are ignored.
func ParseOptions(parameter string) Options {
	var opts Options
	for _, item := range strings.Split(parameter, ",") {
		key, value, _ := strings.Cut(item, "=")
		if strings.TrimSpace(key) != "swift_module" {
			continue
		}
		opts.SwiftModule = strings.TrimSpace(value)
		if opts.SwiftModule == "" {
			opts.SwiftModule = DefaultSwiftModule
		}
	}
	return opts
}

// SwiftEnabled reports whether the Swift output is on.
func (o Options) SwiftEnabled() bool {
	return o.SwiftModule != ""
}
