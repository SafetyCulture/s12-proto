// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Package runtime holds the C# that generated validators call. The generator emits these files
// alongside the validators it produces, so the two always come from the same version of this
// repository and the calls the generator writes cannot outrun the code that serves them.
package runtime

import (
	"embed"
	"io/fs"
)

//go:embed Validation/*.cs
var sources embed.FS

// Dir is the directory the runtime is emitted into, relative to the generation root.
const Dir = "Validation"

// Files returns each runtime source file keyed by its path relative to the generation root.
func Files() (map[string]string, error) {
	entries, err := fs.ReadDir(sources, Dir)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(entries))
	for _, e := range entries {
		name := Dir + "/" + e.Name()
		b, err := fs.ReadFile(sources, name)
		if err != nil {
			return nil, err
		}
		out[name] = string(b)
	}
	return out, nil
}
