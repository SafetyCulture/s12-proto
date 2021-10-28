package plugin

import (
	"fmt"
	"sort"
	"strings"

	validator "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
)

// Any changes made in this file will require Security Engineering consultation/review

// String consts
const (
	stringLenMinSafe    = uint32(1)
	stringLenMinUnsafe  = uint32(1) // <1 not allowed, use optional instead of empty strings
	stringLenMaxSafe    = uint32(1000)
	stringLenMaxUnsafe  = uint32(10000)
	stringLenMinDefault = uint32(1)
	stringLenMaxDefault = uint32(130)
)

// Default allowed regex tokens for validator.string (safe_string)
var stringReDefaultSafe = []string{
	`\pL`,      // Letter (Unicode Category) - any kind of letter from any language
	`\pN`,      // Number (Unicode Category) - any kind of numeric character in any script
	`\x{0020}`, //    Space                  Unicode Category: Space Separator (Zs)
	`\x{0028}`, // (  Left Parenthesis       Unicode Category: Open Punctuation (Ps)
	`\x{0029}`, // )  Right Parenthesis      Unicode Category: Close Punctuation (Pe)
	`\x{002C}`, // ,  Comma                  Unicode Category: Other Punctuation (Po)
	`\x{002E}`, // .  Full Stop              Unicode Category: Other Punctuation (Po)
	`\x{003A}`, // :  Colon                  Unicode Category: Other Punctuation (Po)
	`\x{003F}`, // ?  Question Mark          Unicode Category: Other Punctuation (Po)
	`\x{0040}`, // @  Commercial At          Unicode Category: Other Punctuation (Po)
	`\x{005B}`, // [  Left Square Bracket    Unicode Category: Open Punctuation (Ps)
	`\x{005D}`, // ]  Right Square Bracket   Unicode Category: Close Punctuation (Pe)
	`\x{005F}`, // _  Low Line               Unicode Category: Connector Punctuation (Pc)
	`\x{00BF}`, // ¿  Inverted Question Mark Unicode Category: Other Punctuation (Po)
	`\x{2013}`, // –  En Dash                Unicode Category: Dash Punctuation (Pd)
}

// Default allowed regex tokens for validator.unsafe_string
// Can be further extended by adding symbol categories in the validation field option
// Notice that by default < and > are not included as they are not often used
var stringReDefaultUnsafe = []string{
	// Anything from stringReDefaultSafe, plus:
	`\x{0021}`, // !  Exclamation Mark  Other Punctuation (Po)
	`\x{0022}`, // "  Quotation Mark    Other Punctuation (Po)
	`\x{0023}`, // #  Number Sign       Other Punctuation (Po)
	`\x{0025}`, // %  Percent Sign      Other Punctuation (Po)
	`\x{0026}`, // &  Ampersand         Other Punctuation (Po)
	`\x{0027}`, // '  Apostrophe        Other Punctuation (Po)
	`\x{002A}`, // *  Asterisk          Other Punctuation (Po)
	`\x{002D}`, // -  Hyphen-Minus      Dash Punctuation (Pd)
	`\x{002F}`, // /  Solidus           Other Punctuation (Po)
	`\x{005C}`, // \  Reverse Solidus   Other Punctuation (Po)
}

var stringReLineBreaks = []string{
	`\x{000A}`, // <End of Line> (EOL, LF, NL, \n)  Control (Cc)
	// \x{000D} <Carriage Return> (CR, \r) is stripped by default for multiline strings
}

// Library that holds all regex patterns in format map[regex_id]map[pattern_token]
var regexLib = make(map[string]map[string]bool)

// replace_unsafe option
// The following alternative characters are all in the COMMON script but might not be available in all fonts
// None of these characters are allowed in safe_string unless replace_unsafe option is enabled
// Use string literals otherwise the visual unicode symbol will be used in the replacer string which can fail for certain chars
var stringUnsafeReplacerMap = map[string]string{
	`\u0021`: `\uFF01`, //  ! : ！  EXCLAMATION MARK to FULLWIDTH EXCLAMATION MARK
	`\u0022`: `\u201D`, //  " : ”   QUOTATION MARK to RIGHT DOUBLE QUOTATION MARK (alternative: “ \u0093)
	`\u0023`: `\u0020`, // TODO # : ?   NUMBER SIGN to TODO (not found an appropriate alternative, using space for now)
	`\u0025`: `\u2052`, //  % : ⁒   PERCENT SIGN to COMMERCIAL MINUS SIGN
	`\u0026`: `\u0020`, // TODO  & : ?   AMPERSAND to TODO (not found an appropriate alternative, using space for now)
	`\u0027`: `\u2019`, //  ' : ’   APOSTROPHE to RIGHT SINGLE QUOTATION MARK
	`\u002A`: `\u2217`, //  * : ∗   ASTERISK to ASTERISK OPERATOR
	`\u002B`: `\u2795`, //  + : ➕  PLUS SIGN to HEAVY PLUS SIGN
	`\u002D`: `\u2212`, //  - : −   HYPHEN-MINUS to MINUS SIGN (alternative ‐ U+2010)
	`\u002F`: `\u2215`, //  / : ∕   SOLIDUS to DIVISION SLASH
	`\u003B`: `\u037E`, //  ; : ;   SEMICOLON to GREEK QUESTION MARK
	`\u003C`: `\u02C2`, //  < : ˂   LESS-THAN SIGN to MODIFIER LETTER LEFT ARROWHEAD (similar to >, alternative: U+1438)
	`\u003D`: `\u2E40`, //  = : ⹀   EQUALS SIGN to DOUBLE HYPHEN  (alternative ゠U+30A0)
	`\u003E`: `\u02C3`, //  > : ˃   GREATER-THAN SIGN to MODIFIER LETTER RIGHT ARROWHEAD (alternative in other script: U+16F3F)
	`\u005C`: `\uFF3C`, //  \ : ＼  REVERSE SOLIDUS to FULLWIDTH REVERSE SOLIDUS (alternative ⧵ U+29F5)
	`\u0060`: `\u2019`, //  ` : ’   GRAVE ACCENT to RIGHT SINGLE QUOTATION MARK
	`\u007C`: `\uFFE8`, //  | : ￨  VERTICAL LINE to HALFWIDTH FORMS LIGHT VERTICAL
	// NOTE: when updating, also need to update UnsafeCharReplacer in validator_helpers.go (TODO PA: improve this)
}

// Unicode cats not allowed in safe_string without replace_unsafe option
var restrictedSafeStringSymbols = map[validator.SymbolCategory]struct{}{
	validator.SymbolCategory_PUNCTUATION: {},
	validator.SymbolCategory_MODIFIER:    {},
}

// Translates the SymbolCategory to the matching Unicode regex pattern
var stringSymbolMap = map[validator.SymbolCategory][]string{
	validator.SymbolCategory_CURRENCY: {
		`\pSc`, // Currency symbols like $ £ € - https://www.fileformat.info/info/unicode/category/Sc/list.htm
	},
	validator.SymbolCategory_PUNCTUATION: {
		// UNSAFE: contains several unsafe symbols like < | ` !
		// Can be used with string when replace unsafe option is enabled
		`\pSm`, // Math symbols like + < | - https://www.fileformat.info/info/unicode/category/Sm/list.htm
		`\pPo`, // Punctuation, other category symbols like ! # % ? - https://www.fileformat.info/info/unicode/category/Po/list.htm
	},
	validator.SymbolCategory_MODIFIER: {
		// UNSAFE: Contains some unsafe symbols like \x{0060} GRAVE ACCENT
		// Can be used with string when replace unsafe option is enabled
		`\pSk`, // Modifier symbols like ^ ˥ 🏻
	},
	validator.SymbolCategory_OTHER: {
		`\pSo`,     // Symbol Other including symbols like © ♬ 🐳 😗 - https://www.fileformat.info/info/unicode/category/So/list.htm
		`\x{200D}`, // ZERO WIDTH JOINER which is used to combine some emoticons like 👨‍💻 (👨+💻)
	},
	validator.SymbolCategory_MARK: {`\pM`}, // Mark and combining symbols (all sub cats Spacing Combining, Enclosing, Nonspacing)
	validator.SymbolCategory_COMMON: {
		// Do not add unsafe symbols to this category!
		`\x{007E}`, // ~  TILDE
		`\x{104B}`, // ။  MYANMAR SIGN SECTION
		`\x{2018}`, // ‘  LEFT SINGLE QUOTATION MARK
		`\x{2019}`, // ’  RIGHT SINGLE QUOTATION MARK
		`\x{2022}`, // •  BULLET
		`\x{201C}`, // “  LEFT DOUBLE QUOTATION MARK
		`\x{201D}`, // ”  RIGHT DOUBLE QUOTATION MARK
	},
}

// replace_other option
// Translate rare symbols to a more common alternative or space to sanitise the data
// Most of these occur in less than 0.02% of our data
var stringSymbolReplacerMap = map[string]string{
	// Separator, Space (Zs) to normal space - https://www.fileformat.info/info/unicode/category/Zs/list.htm
	`\u00A0`: `\u0020`, // NO-BREAK SPACE
	`\u1680`: `\u0020`, // OGHAM SPACE MARK
	`\u2000`: `\u0020`, // EN QUAD
	`\u2001`: `\u0020`, // EM QUAD
	`\u2002`: `\u0020`, // EN SPACE
	`\u2003`: `\u0020`, // EM SPACE
	`\u2004`: `\u0020`, // THREE-PER-EM SPACE
	`\u2005`: `\u0020`, // FOUR-PER-EM SPACE
	`\u2006`: `\u0020`, // SIX-PER-EM SPACE
	`\u2007`: `\u0020`, // FIGURE SPACE
	`\u2008`: `\u0020`, // PUNCTUATION SPACE
	`\u2009`: `\u0020`, // THIN SPACE
	`\u200A`: `\u0020`, // HAIR SPACE
	`\u202F`: `\u0020`, // NARROW NO-BREAK SPACE
	`\u205F`: `\u0020`, // MEDIUM MATHEMATICAL SPACE
	`\u3000`: `\u0020`, // IDEOGRAPHIC SPACE
	// Other, format (Control)
	`\u200C`: ``, // ZERO WIDTH NON-JOINER to <nil> (strip)
	`\u200D`: ``, // ZERO WIDTH JOINER to <nil> (strip)
	// TODO: `\u200E`: ``, // LEFT-TO-RIGHT MARK - we do not actively support LTR or RTR atm, likely resulting in UI issues
	// TODO: `\u200F`: ``, // RIGHT-TO-LEFT MARK - needs further work, currently we will reject this input, uncomment this to strip
	`\uFEFF`: ``, // ZERO WIDTH NO-BREAK SPACE
	// Punctuation
	`\u2014`: `\u2013`, // EM DASH to EN dash
	`\u2018`: `\u2019`, // LEFT SINGLE QUOTATION MARK to right (allowed by default)
	`\u3002`: `\u002E`, // IDEOGRAPHIC FULL STOP to normal full stop
	`\uFF0C`: `\u002C`, // FULLWIDTH COMMA to normal comma (from [Po])
	`\uFF1A`: `\u003A`, // FULLWIDTH COLON to normal colon
	// Control chars to normal space (other control chars, especially the ones listed in charDenyList should still be denied)
	`\u0009`: `\u0020`, // <TAB> \t HORIZONTAL TAB
	`\u000A`: `\u0020`, // \n LINE FEED (will not be replaced with multiline option enabled)
	`\u000D`: `\u0020`, // \r CARRIAGE RETURN
	// NOTE: when updating, also need to update SymbolCharReplacer and SymbolCharReplacerMultiline in validator_helpers.go (TODO PA: improve this)
}

// Prepare the string replacer arguments based on provided map contents
func prepareStringReplacerRegex(replacerMap map[string]string, regexId string) []string {
	stringReplacerArguments := []string{}
	for _, safe := range replacerMap {
		// Add to whitelist regex so that the replaced characters are allowed in the validation pattern
		r := strings.Replace(safe, `\u`, `\x{`, 1) + "}"
		prepareRegex(regexId, r)
	}
	return stringReplacerArguments
}

// Prepare a regex by adding one or multiple tokens to the pattern
func prepareRegex(regexId string, patternTokens ...string) {
	for _, patternToken := range patternTokens {
		if _, ok := regexLib[regexId]; !ok {
			regexLib[regexId] = make(map[string]bool)
		}
		regexLib[regexId][patternToken] = true
	}
}

// Add tokens from a prepared regex to another one
// This is used to add safe_string pattern to unsafe_string pattern
func mergeRegex(toRegexId, fromRegexId string) {
	for k, v := range regexLib[fromRegexId] {
		regexLib[toRegexId][k] = v
	}
}

// Return the regex pattern based on tokens added via prepareRegex
func getPreparedRegex(regexId string) (string, error) {

	if _, ok := regexLib[regexId]; !ok {
		// regexId not initialised via prepareRegex
		// might indicate typo, return a pattern that will always fail validation
		// should not continue in the validator though, so just in case errors are not handled
		return "^\b$", fmt.Errorf("invalid regexId %s - initialise the regex with prepareRegex first", regexId)
	}
	// Sort the tokens so that we can consistent results
	tokens := make([]string, 0, len(regexLib[regexId]))
	for k := range regexLib[regexId] {
		tokens = append(tokens, k)
	}
	sort.Strings(tokens)
	return `^[` + strings.Join(tokens, "") + `]+$`, nil
}

// Deny list of characters not allowed in both safe_string and/or unsafe_string
// Normally, none of these characters are allowed by our default lists as they are potentially dangerous
// This deny list is a safeguard to avoid that these characters are added to the 'allow' validator option
// Use string literals as we do not need the actual character
var charDenyList = map[string]struct{}{
	// For safe_string, this is in addition to stringUnsafeReplacerMap and restrictedSafeStringSymbols
	// ASCII Control Characters (0-31 and 127):
	`\u0000`: {}, // Null byte
	`\u0001`: {}, // Start of Heading
	`\u0002`: {}, // Start of Text
	`\u0003`: {}, // End of Text
	`\u0004`: {}, // End of Transmission
	`\u0005`: {}, // Enquiry
	`\u0006`: {}, // Acknowledge
	`\u0007`: {}, // Bell
	`\u0008`: {}, // \b Backspace
	`\u000B`: {}, // \v Vertical Tab
	`\u000C`: {}, // \f Form Feed
	`\u000E`: {}, // Shift Out
	`\u000F`: {}, // Shift In
	`\u0010`: {}, // Data Link Escape
	`\u0011`: {}, // Device Control 1 (XON)
	`\u0012`: {}, // Device Control 2
	`\u0013`: {}, // Device Control 3 (XOFF)
	`\u0014`: {}, // Device Control 4
	`\u0015`: {}, // Negative Acknowledge
	`\u0016`: {}, // Synchronous Idle
	`\u0017`: {}, // End of Transmission Block
	`\u0018`: {}, // Cancel
	`\u0019`: {}, // End of Medium
	`\u001A`: {}, // Substitute
	`\u001B`: {}, // \e Escape
	`\u001C`: {}, // File Separator
	`\u001D`: {}, // Group Separator
	`\u001E`: {}, // Record Separator
	`\u001F`: {}, // Unit Separator
	`\u007F`: {}, // Delete
	// ANSI Characters (128-159), ISO-8859-1 treats these as control characters:
	`\u0081`: {}, // Not used
	`\u008D`: {}, // Not used
	`\u008F`: {}, // Not used
	`\u0090`: {}, // Not used
	`\u009D`: {}, // Not used
}
