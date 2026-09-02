package valtest

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
)

// TestMachineGeneratedValidation_Validate covers the (validator.machine_generated) validator.
// The headline case (per ADR ARCH-546) is the contrast: machine strings that the display-text
// validator.string REJECTS must PASS machine_generated, which validates size + encoding only
// (no allow-list, no normalisation, no mutation).
func TestMachineGeneratedValidation_Validate(t *testing.T) {
	// A base64/base64url page token: contains + / = which validator.string disallows.
	base64Token := "q1w2e3r4T5+/aB=="
	// A CEL / filter expression: contains == < > | which validator.string disallows.
	celExpr := `score >= 10 && status == "OPEN" || labels.exists(l, l < "z")`

	tests := []struct {
		name        string
		mutate      func(m *ValTestMessage)
		shouldError bool
	}{
		// encoding-only: machine strings pass unchanged (no allow-list, no mutation)
		{"Base64Token_encodingOnly", func(m *ValTestMessage) { m.MachineGeneratedEncodingOnly = base64Token }, valid},
		{"CelExpr_encodingOnly", func(m *ValTestMessage) { m.MachineGeneratedEncodingOnly = celExpr }, valid},
		{"InvalidUTF8_encodingOnly", func(m *ValTestMessage) { m.MachineGeneratedEncodingOnly = "bad\xe9token" }, invalid},
		// len range "1:2048"
		{"LenRange_ok", func(m *ValTestMessage) { m.MachineGeneratedLenRange = base64Token }, valid},
		{"LenRange_tooLong", func(m *ValTestMessage) { m.MachineGeneratedLenRange = strings.Repeat("a", 2049) }, invalid},
		// fixed length "16"
		{"LenFixed_ok16", func(m *ValTestMessage) { m.MachineGeneratedLenFixed = strings.Repeat("a", 16) }, valid},
		{"LenFixed_bad15", func(m *ValTestMessage) { m.MachineGeneratedLenFixed = strings.Repeat("a", 15) }, invalid},
		// min-only "16:"
		{"LenMinOnly_ok", func(m *ValTestMessage) { m.MachineGeneratedLenMinOnly = strings.Repeat("a", 16) }, valid},
		{"LenMinOnly_tooShort", func(m *ValTestMessage) { m.MachineGeneratedLenMinOnly = "short" }, invalid},
		// max-only ":2048"
		{"LenMaxOnly_ok", func(m *ValTestMessage) { m.MachineGeneratedLenMaxOnly = base64Token }, valid},
		{"LenMaxOnly_tooLong", func(m *ValTestMessage) { m.MachineGeneratedLenMaxOnly = strings.Repeat("a", 2049) }, invalid},
		// runes "1:10": 10 multibyte runes = 10 codepoints (20 bytes) must pass; 11 must fail
		{"LenRunes_10multibyte_ok", func(m *ValTestMessage) { m.MachineGeneratedLenRunes = strings.Repeat("é", 10) }, valid},
		{"LenRunes_11_tooLong", func(m *ValTestMessage) { m.MachineGeneratedLenRunes = strings.Repeat("a", 11) }, invalid},
		// validate_encoding:false accepts invalid UTF-8
		{"NoEncoding_invalidUTF8_ok", func(m *ValTestMessage) { m.MachineGeneratedNoEncoding = "bad\xe9token" }, valid},
		// log_only "1:8": over-length is logged, not rejected (soft)
		{"LogOnly_overLength_softPass", func(m *ValTestMessage) { m.MachineGeneratedLogOnly = "waytoolongvalue" }, valid},
		// optional empty is skipped
		{"Optional_empty_ok", func(m *ValTestMessage) { m.MachineGeneratedOptional = "" }, valid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := proto.Clone(&valMsg).(*ValTestMessage)
			tt.mutate(m)
			err := m.Validate()
			if tt.shouldError && err == nil {
				t.Errorf("%s: expected an error, got nil", tt.name)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("%s: expected no error, got %v", tt.name, err)
			}
		})
	}

	// ADR acceptance proof: the SAME values FAIL validator.string (display-text, field `description`),
	// demonstrating why a dedicated machine-string validator is needed.
	for _, c := range []struct {
		name  string
		value string
	}{
		{"Contrast_base64FailsStringValidator", base64Token},
		{"Contrast_celFailsStringValidator", celExpr},
	} {
		t.Run(c.name, func(t *testing.T) {
			m := proto.Clone(&valMsg).(*ValTestMessage)
			m.Description = c.value // description uses (validator.string), which enforces a display-text allow-list
			if err := m.Validate(); err == nil {
				t.Errorf("expected validator.string to reject %q, but it passed", c.value)
			}
		})
	}
}
