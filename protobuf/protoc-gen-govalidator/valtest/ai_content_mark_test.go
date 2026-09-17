package valtest

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
)

// mark is the Mitti AI content mark, the sequence a generated string carries so
// it is detectable as artificially generated. U+2062 INVISIBLE TIMES, U+2064
// INVISIBLE PLUS, U+2062 INVISIBLE TIMES.
const mark = "\u2062\u2064\u2062"

// TestAIContentMark_IsAllowed covers the two codepoints added to the default
// allow list so that generated text can carry its mark through a write.
//
// Before this, every field using validator.string or validator.unsafe_string
// rejected a marked value with "value must only have valid characters". The
// allow list has no token for Unicode category Cf and no symbol category maps
// to it, so no field option could admit them.
func TestAIContentMark_IsAllowed(t *testing.T) {
	// No hyphen. validator.string allows it only through unsafe_string or the
	// allow option, and a hyphen in the fixture would hide what is under test.
	marked := "Cooler temperature" + mark + " excursion"

	tests := []struct {
		name        string
		mutate      func(m *ValTestMessage)
		shouldError bool
	}{
		// The two validators that carry display text.
		{"SafeString_marked", func(m *ValTestMessage) { m.Description = marked }, valid},
		{"UnsafeString_marked", func(m *ValTestMessage) { m.LongString = marked }, valid},

		// A permissive field that mutates on the way in. replace_other strips the
		// other invisible format characters, so this proves the mark is not one of
		// them. A silent strip is worse than a rejection, because the write appears
		// to succeed and the content ends up unmarked.
		{"Permissive_marked", func(m *ValTestMessage) { m.ScPermissive = marked }, valid},

		// trim removes leading and trailing whitespace. The mark is not whitespace,
		// so a mark at either end survives.
		{"Trim_marked", func(m *ValTestMessage) { m.TrimString = mark + "a blocked fire exit" + mark }, valid},

		// Every other format character stays rejected. The mark is two named
		// codepoints, not the category, so the bidi overrides and the tag block
		// used for spoofing and for prompt injection are untouched.
		{"WordJoiner_rejected", func(m *ValTestMessage) { m.Description = "sc\u2060am" }, invalid},
		{"FunctionApplication_rejected", func(m *ValTestMessage) { m.Description = "sc\u2061am" }, invalid},
		{"InvisibleSeparator_rejected", func(m *ValTestMessage) { m.Description = "sc\u2063am" }, invalid},
		{"ZeroWidthSpace_rejected", func(m *ValTestMessage) { m.Description = "sc\u200bam" }, invalid},
		{"ZeroWidthNonJoiner_rejected", func(m *ValTestMessage) { m.Description = "sc\u200cam" }, invalid},
		{"ZeroWidthJoiner_rejected", func(m *ValTestMessage) { m.Description = "sc\u200dam" }, invalid},
		{"LeftToRightMark_rejected", func(m *ValTestMessage) { m.Description = "sc\u200eam" }, invalid},
		{"RightToLeftOverride_rejected", func(m *ValTestMessage) { m.Description = "sc\u202eam" }, invalid},
		{"ByteOrderMark_rejected", func(m *ValTestMessage) { m.Description = "sc\ufeffam" }, invalid},
		{"LanguageTag_rejected", func(m *ValTestMessage) { m.Description = "sc\U000e0001am" }, invalid},
		{"TagSpace_rejected", func(m *ValTestMessage) { m.Description = "sc\U000e0020am" }, invalid},
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
}

// TestAIContentMark_IsNotMutated proves the mark reaches the store intact.
//
// A field that accepts the mark and then removes it is the worst outcome, so
// this checks the value after Validate() on the field with every mutating
// option turned on.
func TestAIContentMark_IsNotMutated(t *testing.T) {
	marked := "A blocked fire exit was found" + mark + " on level two"

	m := proto.Clone(&valMsg).(*ValTestMessage)
	m.ScPermissive = marked
	if err := m.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got := m.GetScPermissive(); got != marked {
		t.Errorf("the mark did not survive validation\n got %q\nwant %q", got, marked)
	}
}

// TestAIContentMark_CostsNoLength covers the second blocker the allow list alone
// does not fix.
//
// A value already at a field maximum fails the moment it is marked. The mark is
// up to 8 copies of 3 codepoints, so between 3 and 72 units over depending on
// the cap and on whether the field counts runes or bytes. The length is
// therefore measured on a copy with the mark removed.
func TestAIContentMark_CostsNoLength(t *testing.T) {
	tests := []struct {
		name        string
		mutate      func(m *ValTestMessage)
		shouldError bool
	}{
		// description is len ":750" counted in bytes. 750 bytes of content plus
		// the mark is 759 bytes on the wire.
		{"MaxLengthBytes_marked", func(m *ValTestMessage) {
			m.Description = strings.Repeat("a", 750) + mark
		}, valid},
		{"OverMaxLength_stillRejected", func(m *ValTestMessage) {
			m.Description = strings.Repeat("a", 751) + mark
		}, invalid},

		// fixed_string is len "4". A fixed length is the strictest case, because
		// any carrier at all breaks an equality check.
		{"FixedLength_marked", func(m *ValTestMessage) { m.FixedString = "abcd" + mark }, valid},
		{"FixedLengthShort_stillRejected", func(m *ValTestMessage) { m.FixedString = "abc" + mark }, invalid},

		// rune_string is len "4" with runes: true.
		{"FixedRunes_marked", func(m *ValTestMessage) { m.RuneString = "abcd" + mark }, valid},

		// title is len "3:50".
		{"MinLength_marked", func(m *ValTestMessage) { m.Title = "abc" + mark }, valid},
		{"UnderMinLength_stillRejected", func(m *ValTestMessage) { m.Title = "ab" + mark }, invalid},

		// A value that is nothing but the mark measures zero and fails the
		// minimum, which is what should happen. An invisible mark is not content.
		{"MarkOnly_rejected", func(m *ValTestMessage) { m.Description = mark }, invalid},
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
}

// TestAIContentMark_DoesNotHideAURL covers the hardening the allow list needs.
//
// reject_url and break_partial_url match the literal characters of a URL. The
// mark is invisible and now allowed, so a mark between the host and the dot
// reads to a person as a URL while failing the pattern. Both controls run
// against a copy with the mark removed for that reason.
func TestAIContentMark_DoesNotHideAURL(t *testing.T) {
	t.Run("RejectUrl_seesThroughTheMark", func(t *testing.T) {
		m := &NonUrlMessage{RejectUrlTest: "Look at https://evil" + mark + ".com now"}
		if err := m.Validate(); err == nil {
			t.Error("expected the URL to be rejected, got nil")
		}
	})

	t.Run("RejectUrl_allowsAMarkedSentence", func(t *testing.T) {
		m := &NonUrlMessage{RejectUrlTest: "A blocked fire exit" + mark + " on level two"}
		if err := m.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("BreakPartialUrl_seesThroughTheMark", func(t *testing.T) {
		m := &NonUrlMessage{BreakPartialUrlTest: "Check out example" + mark + ".com for info"}
		if err := m.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		const want = "Check out example. com for info"
		if got := m.GetBreakPartialUrlTest(); got != want {
			t.Errorf("the partial URL was not broken\n got %q\nwant %q", got, want)
		}
	})

	// Only the run in front of the dot goes. Everything else about the value is
	// left alone, so the field keeps its mark and stays compliant.
	t.Run("BreakPartialUrl_keepsEveryOtherCopy", func(t *testing.T) {
		m := &NonUrlMessage{
			BreakPartialUrlTest: "A blocked fire exit" + mark +
				" was found at example" + mark + ".com on level" + mark + " two",
		}
		if err := m.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := "A blocked fire exit" + mark +
			" was found at example. com on level" + mark + " two"
		if got := m.GetBreakPartialUrlTest(); got != want {
			t.Errorf("the surviving marks are wrong\n got %q\nwant %q", got, want)
		}
		if !strings.Contains(m.GetBreakPartialUrlTest(), mark) {
			t.Error("the field lost its mark entirely")
		}
	})

	// A mark straight after the dot never hid anything, because the boundary is
	// in front of the dot. It is captured as the character to push along.
	t.Run("BreakPartialUrl_markAfterTheDotIsKept", func(t *testing.T) {
		m := &NonUrlMessage{BreakPartialUrlTest: "Check out example." + mark + "com for info"}
		if err := m.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(m.GetBreakPartialUrlTest(), mark) {
			t.Error("a mark after the dot should have survived")
		}
		if !strings.Contains(m.GetBreakPartialUrlTest(), ". ") {
			t.Error("the partial URL was not broken")
		}
	})
}
