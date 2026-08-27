// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"slices"
	"testing"
)

// TestTranslatePatternNarrowsTwoLetterCategories covers the second reading of a
// class holding a two-letter category token: the token names the category, and
// the letter that the widened reading leaves behind as a literal is part of it.
func TestTranslatePatternNarrowsTwoLetterCategories(t *testing.T) {
	tests := []struct {
		name             string
		goPattern        string
		wantPattern      string
		wantCategories   []string
		wantNarrow       string
		wantNarrowCatgys []string
	}{
		{
			name:             "CurrencyWidensToEverySymbol",
			goPattern:        `^[\pL\pN\pSc\x{0020}]+$`,
			wantPattern:      `\A[\p{L}\p{N}\p{S}c\u0020]+\z`,
			wantCategories:   []string{"L", "N", "S"},
			wantNarrow:       `\A[\p{L}\p{N}\p{Sc}\u0020]+\z`,
			wantNarrowCatgys: []string{"L", "N", "Sc"},
		},
		{
			name:             "PunctuationAndSymbolTokensTogether",
			goPattern:        `^[\pM\pPo\pSk\x{2013}]+$`,
			wantPattern:      `\A[\p{M}\p{P}o\p{S}k\u2013]+\z`,
			wantCategories:   []string{"M", "P", "S"},
			wantNarrow:       `\A[\p{M}\p{Po}\p{Sk}\u2013]+\z`,
			wantNarrowCatgys: []string{"M", "Po", "Sk"},
		},
		{
			name:           "OneLetterTokensAreAlreadyWhatTheyName",
			goPattern:      `^[\pL\pN\x{0020}]+$`,
			wantPattern:    `\A[\p{L}\p{N}\u0020]+\z`,
			wantCategories: []string{"L", "N"},
			wantNarrow:     "",
		},
		{
			name:           "LetterAfterATokenIsNotAlwaysACategory",
			goPattern:      `^[\pLz]+$`,
			wantPattern:    `\A[\p{L}z]+\z`,
			wantCategories: []string{"L"},
			wantNarrow:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := translatePattern(tc.goPattern)
			if err != nil {
				t.Fatalf("translatePattern(%q) = %v", tc.goPattern, err)
			}
			if got.Pattern != tc.wantPattern {
				t.Errorf("Pattern = %q, want %q", got.Pattern, tc.wantPattern)
			}
			if !slices.Equal(got.AstralCategories, tc.wantCategories) {
				t.Errorf("AstralCategories = %v, want %v", got.AstralCategories, tc.wantCategories)
			}
			if got.NarrowPattern != tc.wantNarrow {
				t.Errorf("NarrowPattern = %q, want %q", got.NarrowPattern, tc.wantNarrow)
			}
			if !slices.Equal(got.NarrowAstralCategories, tc.wantNarrowCatgys) {
				t.Errorf("NarrowAstralCategories = %v, want %v", got.NarrowAstralCategories, tc.wantNarrowCatgys)
			}
		})
	}
}
