package valtest

import "testing"

// TestSymbolCategoryIsScopedToTheDeclaredCategory covers symbol_string, which
// declares symbols: [CURRENCY].
//
// A category is written into the character class as a Unicode class token. RE2
// reads an unbraced two-letter name as a one-letter class followed by a literal,
// so the token has to carry braces for the class to mean the category it names.
func TestSymbolCategoryIsScopedToTheDeclaredCategory(t *testing.T) {
	accepted := []struct{ name, value string }{
		{"DollarSign", "$"},
		{"PoundSign", "£"},
		{"EuroSign", "€"},
		{"YenSign", "¥"},
	}
	for _, tc := range accepted {
		t.Run(tc.name, func(t *testing.T) {
			if err := getValMsg(&ValTestMessage{SymbolString: tc.value}).Validate(); err != nil {
				t.Errorf("Validate(%q) = %v, want no error", tc.value, err)
			}
		})
	}

	rejected := []struct{ name, value string }{
		{"MathSymbol", "<"},
		{"MathSymbolPipe", "|"},
		{"ModifierGraveAccent", "`"},
		{"ModifierCircumflex", "^"},
		{"OtherSymbolCopyright", "©"},
	}
	for _, tc := range rejected {
		t.Run(tc.name, func(t *testing.T) {
			if err := getValMsg(&ValTestMessage{SymbolString: tc.value}).Validate(); err == nil {
				t.Errorf("Validate(%q) = nil, want an error", tc.value)
			}
		})
	}
}
