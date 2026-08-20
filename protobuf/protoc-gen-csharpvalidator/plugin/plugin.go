// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-csharpvalidator/runtime"
)

// Generate writes the C# validation runtime and the validators for each selected file.
func Generate(p *protogen.Plugin) error {
	return emitRuntime(p)
}

// emitRuntime writes the runtime once per invocation, whatever the input set.
func emitRuntime(p *protogen.Plugin) error {
	files, err := runtime.Files()
	if err != nil {
		return err
	}

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if _, err := p.NewGeneratedFile(name, "").Write([]byte(files[name])); err != nil {
			return err
		}
	}
	return nil
}
