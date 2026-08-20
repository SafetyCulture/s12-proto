package runtime

import (
	"strings"
	"testing"
)

// The generator writes calls against these types; losing one from the embed turns into a C# compile
// error in the consuming repository, where it is far harder to trace.
var required = []string{
	"Validation/CharacterClass.cs",
	"Validation/IValidatableMessage.cs",
	"Validation/StringMutators.cs",
	"Validation/ValidationError.cs",
	"Validation/ValidationLog.cs",
	"Validation/ValidatorHelpers.cs",
	"Validation/ValidatorPatterns.cs",
}

func TestFiles(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range required {
		content, ok := files[name]
		if !ok {
			t.Errorf("%s missing from the embedded runtime", name)
			continue
		}
		if !strings.Contains(content, "namespace S12.Protobuf.Validation;") {
			t.Errorf("%s is not in the S12.Protobuf.Validation namespace", name)
		}
	}

	if len(files) != len(required) {
		t.Errorf("embedded %d files, expected %d: the required list needs updating", len(files), len(required))
	}
}
