// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"strings"
)

type outputNaming string

const (
	fullPath          outputNaming = "FullPath"
	pathToUnderscores outputNaming = "PathToUnderscores"
	dropPath          outputNaming = "DropPath"
)

type enumGeneration string

const (
	enumGenerationNone              enumGeneration = "None"
	enumGenerationNonexhaustive     enumGeneration = "Nonexhaustive"
	enumGenerationNonexhaustiveWarn enumGeneration = "NonexhaustiveWarn"
)

const maxCasesInSwitch = 500

type hiddenNames uint

const (
	hiddenFields hiddenNames = 1 << iota
	hiddenEnumValues
	hiddenTypes
)

// options holds the protoc-gen-swift options of a request. The keys of the
// other outputs of the plugin pass through.
type options struct {
	outputNaming              outputNaming
	enumGeneration            enumGeneration
	visibility                string
	importDirective           string
	swiftProtobufModule       string
	accessLevelOnImports      bool
	stripNonfunctionalCodegen bool
	hiddenNames               hiddenNames
}

// visibilitySnippet is the access modifier with its trailing space.
func (o *options) visibilitySnippet() string {
	switch o.visibility {
	case "public":
		return "public "
	case "package":
		return "package "
	}
	return ""
}

func (o *options) enumAnnotation() string {
	switch o.enumGeneration {
	case enumGenerationNonexhaustive:
		return "@nonexhaustive "
	case enumGenerationNonexhaustiveWarn:
		return "@nonexhaustive(warn) "
	}
	return ""
}

// splitNonEmpty splits a comma-separated value. An empty value has no part.
func splitNonEmpty(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

func parseBool(value string) (bool, bool) {
	switch value {
	case "true":
		return true, true
	case "false":
		return false, true
	}
	return false, false
}

// parseOptions reads the comma-separated parameter of a request. The options of the
// other outputs of the plugin pass through: isPluginKey names them, and nil names none.
func parseOptions(parameter string, isPluginKey func(key string) bool) (*options, error) {
	opts := &options{
		outputNaming:        fullPath,
		enumGeneration:      enumGenerationNone,
		visibility:          "internal",
		swiftProtobufModule: "SwiftProtobuf",
	}
	implementationOnly := false
	if parameter != "" {
		for _, part := range strings.Split(parameter, ",") {
			key, value, _ := strings.Cut(part, "=")
			key, value = strings.TrimSpace(key), strings.TrimSpace(value)
			if isPluginKey != nil && isPluginKey(key) {
				continue
			}
			invalid := fmt.Errorf("unknown value %q for the option %s", value, key)
			switch key {
			case "FileNaming":
				switch strings.ToLower(value) {
				case "fullpath", "full_path":
					opts.outputNaming = fullPath
				case "pathtounderscores", "path_to_underscores":
					opts.outputNaming = pathToUnderscores
				case "droppath", "drop_path":
					opts.outputNaming = dropPath
				default:
					return nil, invalid
				}
			case "EnumGeneration":
				switch strings.ToLower(value) {
				case "none":
					opts.enumGeneration = enumGenerationNone
				case "nonexhaustive":
					opts.enumGeneration = enumGenerationNonexhaustive
				case "nonexhaustivewarn":
					opts.enumGeneration = enumGenerationNonexhaustiveWarn
				default:
					return nil, invalid
				}
			case "ProtoPathModuleMappings":
				if value != "" {
					return nil, fmt.Errorf("the option ProtoPathModuleMappings is not supported: all files of a run belong to one Swift module")
				}
			case "Visibility":
				switch strings.ToLower(value) {
				case "internal", "public", "package":
					opts.visibility = strings.ToLower(value)
				default:
					return nil, invalid
				}
			case "SwiftProtobufModuleName":
				if !isValidSwiftIdentifier(value) {
					return nil, invalid
				}
				opts.swiftProtobufModule = value
			case "ImplementationOnlyImports":
				b, ok := parseBool(value)
				if !ok {
					return nil, invalid
				}
				implementationOnly = b
			case "UseAccessLevelOnImports":
				b, ok := parseBool(value)
				if !ok {
					return nil, invalid
				}
				opts.accessLevelOnImports = b
			case "experimental_strip_nonfunctional_codegen":
				if value == "" {
					opts.stripNonfunctionalCodegen = true
				} else if b, ok := parseBool(value); ok {
					opts.stripNonfunctionalCodegen = b
				} else {
					return nil, invalid
				}
			case "ExperimentalHiddenNames":
				var names hiddenNames
				for _, name := range splitNonEmpty(value) {
					switch strings.TrimSpace(name) {
					case "Fields":
						names |= hiddenFields
					case "EnumValues":
						names |= hiddenEnumValues
					case "Types":
						names |= hiddenTypes
					case "All":
						names |= hiddenFields | hiddenEnumValues | hiddenTypes
					default:
						return nil, invalid
					}
				}
				opts.hiddenNames = names
			default:
				return nil, fmt.Errorf("unknown generation option %q", key)
			}
		}
	}
	switch {
	case implementationOnly && opts.accessLevelOnImports:
		return nil, fmt.Errorf("when using access levels on imports the @_implementationOnly option is unnecessary; disable @_implementationOnly imports")
	case implementationOnly:
		opts.importDirective = "@_implementationOnly import"
	case opts.accessLevelOnImports:
		opts.importDirective = opts.visibility + " import"
	default:
		opts.importDirective = "import"
	}
	if implementationOnly && opts.visibility != "internal" {
		return nil, fmt.Errorf("cannot use @_implementationOnly imports when the proto visibility is public or package; either change the visibility to internal, or disable @_implementationOnly imports")
	}
	return opts, nil
}

func (o *options) hides(names hiddenNames) bool { return o.hiddenNames&names != 0 }

// outputFilename is the name of the output file of a proto file.
func (o *options) outputFilename(protoName string) string {
	dir, base, _ := splitPath(protoName)
	const ext = ".pb.swift"
	switch o.outputNaming {
	case pathToUnderscores:
		return strings.ReplaceAll(dir, "/", "_") + base + ext
	case dropPath:
		return base + ext
	}
	return dir + base + ext
}
