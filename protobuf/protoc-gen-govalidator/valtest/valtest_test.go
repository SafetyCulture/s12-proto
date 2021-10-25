package valtest

import (
	"bufio"
	fmt "fmt"
	"os"
	"strings"
	"testing"

	"github.com/SafetyCulture/s12-proto/protobuf/s12proto"
)

const (
	id               string = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
	legacyId         string = "56341C6E-35A7-4C97-9C5E-7AC79673EAB2"
	legacyLongIdFail string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-9-00000190327E0675"     // length = 49 (without dashes)
	legacyLongId1    string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-90-00000190327E0675"    // length = 50 (without dashes)
	legacyLongId2    string = "005F2E38-8426-48AF-94DE-5FEA3A396EEA-891-00000153F68896DC"   // length = 51 (without dashes)
	legacyLongId3    string = "007B516E-53F1-4AA0-ABAF-8C78342A2C82-2388-00000221F1C2BD1E"  // length = 52 (without dashes)
	legacyLongId4    string = "00709A17-151F-4CFC-B412-F080343ED84D-11977-000010227B4C60A9" // length = 53 (without dashes)
	email            string = "email@example.com"
	password         string = "1234567!"
	name             string = "Γιώργος Νταλάρας"
	idv4             string = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
	valid            bool   = false
	invalid          bool   = true
	emptyString      string = "<EMPTY>"
)

var invalidEmails = []string{
	"first.last@domain.com.au>", // trailing ">"" is invalid
	"email",
	"example@example.com66",
	"example@domain", // missing domain name
	"example@example.com�",
	"example@example.co uk",
	// "bbbccc.aa@，abcdef.com.uk", // TODO PA: this is currently considered valid by govalidator
	"m _example@example.com", // space is invalid
	"mail@127.0.0.1",
	"test@gmail.com3",
	"1@1.1",
}

var validEmails = []string{
	"first.last@domain.com.au",
	"email@gmail.com",
	"example@safetyculture.io",
	"example1234@non-existing-domain123837.com", // still a valid email format
	"name@組織.香港",
	"example+alias@EXAMPLE.COM",
	"example|2@example.com",
	"mULtiCase@exAMplE.cOm",
}

var validSafeStrings = []string{
	`abc`,
	`Hello`,
	`Γιώργος Νταλάρας`,
	"file:///valid_safe_strings.txt",
}

var validReplaceUnsafeStrings = []string{
	`<script>alert(document.cookie)</script>`,
	`' OR 1=1`,
	`../etc/passwd`,
	`1+1=2/4&a|b`,
}

var testTitles = []string{
	// "file:///test-data/NOCOMMIT_valid_strings.txt",
}

var testPermissive = []string{
	// "file:///test-data/NOCOMMIT_emoticon_strings.txt",
	// "file:///test-data/NOCOMMIT_actions-issues-PUA.txt",
}

var invalidSafeStrings = []string{
	`<script>alert(document.cookie)</script>`,
	`' OR 1=1`,
	`../etc/passwd`,
	`1+1=2/4&a|b`,
	"file:///invalid_safe_strings.txt",
}

var validNames = []string{
	"file:///valid_names.txt",
}

var valMsg = ValTestMessage{
	Id:                    id,
	Ids:                   []string{id, id},
	LegacyId:              legacyId,
	MediaId:               id,
	InnerLegacyId:         &InnerMessageWithLegacyId{Id: id},
	Email:                 email,
	OptEmail:              email,
	Description:           "François Truffaut 久保田 利伸 text",
	Password:              password,
	Title:                 "Short text, ok",
	FixedString:           "abcd",
	RuneString:            "çois",
	ReplaceString:         "With 'unsafe' \"<chars>",
	NotReplaceString:      "Safe chars ç only 123.",
	AllowString:           "Accept ~ and #",
	SymbolString:          "Accept $ £ ¥ €",
	SymbolsString:         "Accept 🌏 💯 and a\u030C", // ǎ as in a + the caron
	NewlineString:         "Accept\nNewlines\n\rYeeha",
	InvalidEncodingString: "Accept invalid \xe9",
	OptString:             "Optional",
	TrimString:            "   Trim me   \t",
	AllString:             " Lot of checks here>",
	Name:                  "Sinéad O'Connor",
	NoValidation:          "<really?>' OR 1=1",
	ContactOneof:          &ValTestMessage_Phone{Phone: "145-456.123"},
	MsgRequired:           &InnerMessage{Id: id},
	NestedMessage: &ValTestMessage_NestedMessage{
		Val: "inner val",
		// InnerNestedMessage: &ValTestMessage_NestedMessage_InnerNestedMessage{
		// 	InnerVal: "abc def",
		// },
		NestedEmail: email,
		MemberEmails: []string{
			email,
			email,
		},
	},
	// NotSupported: ,
}

// omit optional fields here
var valMsgOpts = ValTestMessage{
	Id:                    id,
	Ids:                   []string{id, id},
	LegacyId:              legacyId,
	InnerLegacyId:         &InnerMessageWithLegacyId{Id: legacyId},
	Email:                 email,
	Description:           "François Truffaut 久保田 利伸 text",
	Password:              password,
	Title:                 "Short text, ok",
	FixedString:           "abcd",
	RuneString:            "çois",
	ReplaceString:         "With 'unsafe' \"<chars>",
	NotReplaceString:      "Safe chars ç only 123.",
	AllowString:           "Accept ~ and #",
	SymbolString:          "Accept $ £ ¥ €",
	SymbolsString:         "Accept 🌏 💯 and a\u030C", // ǎ as in a + the caron
	NewlineString:         "Accept\nNewlines\n\rYeeha",
	InvalidEncodingString: "Accept invalid \xe9",
	TrimString:            "   Trim me   \t",
	AllString:             " Lot of checks here>",
	NoValidation:          "<really?>' OR 1=1",
	ContactOneof:          &ValTestMessage_Phone{Phone: "145-456.123"},
	MsgRequired:           &InnerMessage{Id: id},
	NestedMessage: &ValTestMessage_NestedMessage{
		Val: "inner val",
		// InnerNestedMessage: &ValTestMessage_NestedMessage_InnerNestedMessage{
		// 	InnerVal: "abc def",
		// },
		NestedEmail: email,
		MemberEmails: []string{
			email,
			email,
		},
	},
	// NotSupported: ,
}

func readFiles(list []string) []string {
	// read string array and add contents from file for file handlers
	outList := []string{}
	for _, filename := range list {
		if !strings.HasPrefix(filename, "file:///") {
			outList = append(outList, filename)
			continue
		}
		// read file
		path, _ := os.Getwd()
		filename = strings.Replace(filename, "file:///", path+"/", 1)
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open: %v\n", err)
			continue
		} else {
			fmt.Println("ADDED file", filename)
		}
		lines := bufio.NewScanner(f)
		// add each line of the file as test input
		for lines.Scan() {
			outList = append(outList, lines.Text())
		}
	}

	return outList
}

func genString(len int) string {
	return strings.Repeat("A", len)
}

func replaceEmpty(s string) string {
	if s == emptyString {
		return ""
	}
	return s
}

func getValMsg(m ValTestMessage) *ValTestMessage {
	newMsg := valMsg
	if m.Id != "" {
		newMsg.Id = replaceEmpty(m.Id)
	}
	if m.Ids != nil {
		newMsg.Ids = m.Ids
	}
	if m.LegacyId != "" {
		newMsg.LegacyId = replaceEmpty(m.LegacyId)
	}
	if m.LegacyId != "" {
		newMsg.LegacyId = replaceEmpty(m.LegacyId)
	}
	if m.MediaId != "" {
		newMsg.MediaId = replaceEmpty(m.MediaId)
	}
	// if m.InnerLegacyId != "" {
	// 	newMsg.InnerLegacyId = m.InnerLegacyId
	// }
	if m.Email != "" {
		newMsg.Email = replaceEmpty(m.Email)
	}
	if m.OptEmail != "" {
		newMsg.OptEmail = replaceEmpty(m.OptEmail)
	}
	if m.Description != "" {
		newMsg.Description = replaceEmpty(m.Description)
	}
	if m.Password != "" {
		newMsg.Password = replaceEmpty(m.Password)
	}
	if m.Title != "" {
		newMsg.Title = replaceEmpty(m.Title)
	}
	if m.FixedString != "" {
		newMsg.FixedString = replaceEmpty(m.FixedString)
	}
	if m.RuneString != "" {
		newMsg.RuneString = replaceEmpty(m.RuneString)
	}
	if m.ReplaceString != "" {
		newMsg.ReplaceString = replaceEmpty(m.ReplaceString)
	}
	if m.NotReplaceString != "" {
		newMsg.NotReplaceString = replaceEmpty(m.NotReplaceString)
	}
	if m.AllowString != "" {
		newMsg.AllowString = replaceEmpty(m.AllowString)
	}
	if m.SymbolString != "" {
		newMsg.SymbolString = replaceEmpty(m.SymbolString)
	}
	if m.SymbolsString != "" {
		newMsg.SymbolsString = replaceEmpty(m.SymbolsString)
	}
	if m.NewlineString != "" {
		newMsg.NewlineString = replaceEmpty(m.NewlineString)
	}
	if m.InvalidEncodingString != "" {
		newMsg.InvalidEncodingString = replaceEmpty(m.InvalidEncodingString)
	}
	if m.OptString != "" {
		newMsg.OptString = replaceEmpty(m.OptString)
	}
	if m.TrimString != "" {
		newMsg.TrimString = replaceEmpty(m.TrimString)
	}
	if m.AllString != "" {
		newMsg.AllString = replaceEmpty(m.AllString)
	}
	if m.Name != "" {
		newMsg.Name = replaceEmpty(m.Name)
	}
	if m.ScTitle != "" {
		newMsg.ScTitle = replaceEmpty(m.ScTitle)
	}
	if m.ScPermissive != "" {
		newMsg.ScPermissive = replaceEmpty(m.ScPermissive)
	}
	if m.NotSanitisePua != "" {
		newMsg.NotSanitisePua = replaceEmpty(m.NotSanitisePua)
	}
	if m.SanitisePua != "" {
		newMsg.SanitisePua = replaceEmpty(m.SanitisePua)
	}
	if m.SanitiseLength != "" {
		newMsg.SanitiseLength = replaceEmpty(m.SanitiseLength)
	}
	if m.NoValidation != "" {
		newMsg.NoValidation = replaceEmpty(m.NoValidation)
	}
	// if m.ContactOneof != "" {
	// 	newMsg.ContactOneof = m.ContactOneof
	// }
	// if m.MsgRequired != "" {
	// 	newMsg.MsgRequired = m.MsgRequired
	// }
	return &newMsg
}

func TestValidationRules(t *testing.T) {

	sb := strings.Builder{}
	for i := 0; i < 2001; i++ {
		sb.WriteRune('a')
	}

	type TestSet struct {
		name        string
		input       s12proto.Validator
		shouldError bool
	}

	tests := []TestSet{
		{
			"ValidMessage",
			&valMsg,
			valid,
		},
		{
			"ValidMessageNoOptFields",
			&valMsgOpts,
			valid,
		},
		{
			"InvalidId",
			getValMsg(ValTestMessage{Id: "abc"}),
			invalid,
		},
		{
			"InvalidLegacyId",
			getValMsg(ValTestMessage{LegacyId: legacyLongIdFail}),
			invalid,
		},
		{
			"ValidLegacyIdLong1",
			getValMsg(ValTestMessage{LegacyId: legacyLongId1}),
			valid,
		},
		{
			"ValidLegacyIdLong2",
			getValMsg(ValTestMessage{LegacyId: legacyLongId2}),
			valid,
		},
		{
			"ValidLegacyIdLong3",
			getValMsg(ValTestMessage{LegacyId: legacyLongId3}),
			valid,
		},
		{
			"ValidLegacyIdLong4",
			getValMsg(ValTestMessage{LegacyId: legacyLongId4}),
			valid,
		},
		{
			"InvalidIds",
			getValMsg(ValTestMessage{Ids: []string{id, "audit_invalid", id}}),
			invalid,
		},
		{
			"ValidEmail",
			getValMsg(ValTestMessage{Email: email}),
			valid,
		}, {
			"ValidDescriptionLen1",
			getValMsg(ValTestMessage{Description: genString(1)}),
			valid,
		}, {
			"ValidDescriptionLen2",
			getValMsg(ValTestMessage{Description: genString(2)}),
			valid,
		}, {
			"ValidDescriptionLen50",
			getValMsg(ValTestMessage{Description: genString(50)}),
			valid,
		}, {
			"ValidDescriptionLen100",
			getValMsg(ValTestMessage{Description: genString(100)}),
			valid,
		}, {
			"ValidDescriptionLen750",
			getValMsg(ValTestMessage{Description: genString(750)}),
			valid,
		}, {
			"InvalidDescriptionLen751",
			getValMsg(ValTestMessage{Description: genString(751)}),
			invalid,
		}, {
			"InvalidDescriptionLen0",
			getValMsg(ValTestMessage{Description: emptyString}),
			invalid,
		}, {
			"InvalidDescriptionLen10k",
			getValMsg(ValTestMessage{Description: genString(10000)}),
			invalid,
		}, {
			"InvalidEncoding",
			getValMsg(ValTestMessage{Description: "X\xe9X invalid X\xa3X encoding"}),
			invalid,
		}, {
			"InvalidPUA",
			getValMsg(ValTestMessage{NotSanitisePua: "X\uf0a7X InvalidPUA"}),
			invalid,
		}, {
			"ValidSanitisedPUA",
			getValMsg(ValTestMessage{SanitisePua: "X\uf0a7X ValidSanitisedPUA"}),
			valid,
		}, {
			"ValidSanitisedLength",
			getValMsg(ValTestMessage{SanitiseLength: "A\uf0a7B"}), // sanitised to length 2 (AB)
			valid,
		}, {
			"InvalidSanitisedLength",
			getValMsg(ValTestMessage{SanitiseLength: "A\uf0a7"}), // sanitised to length 1 (A)
			invalid,
		}, {
			"InvalidBidiRtL",
			getValMsg(ValTestMessage{ScPermissive: "InvalidBidi RtL\u200f"}),
			invalid,
		}, {
			"InvalidBidiLtR",
			getValMsg(ValTestMessage{ScPermissive: "\u200eInvalidBidi LtR"}),
			invalid,
		}, {
			"ValidPasswordLong",
			getValMsg(ValTestMessage{Password: "XCzUDDdpwvFR@MoGzsVP@hvjmNqjPG2bNb9G6uz7"}), // not a real password
			valid,
		}, {
			"ValidPasswordShort",
			getValMsg(ValTestMessage{Password: "12345678"}), // not a real password
			valid,
		}, {
			"ValidPasswordUnsafeChars",
			getValMsg(ValTestMessage{Password: "<$pass>'|\"../etc/password"}), // not a real password
			valid,
		}, {
			"InvalidPasswordLength",
			getValMsg(ValTestMessage{Password: "1234567"}), // not a real password
			invalid,
		}, {
			"ValidFixedLenString",
			getValMsg(ValTestMessage{FixedString: "1234"}),
			valid,
		}, {
			"InvalidFixedLenStringMB",
			getValMsg(ValTestMessage{FixedString: "123é"}), // é takes up 2 bytes, making this len 5 instead of 4
			invalid,
		}, {
			"InvalidFixedLenStringTooShort",
			getValMsg(ValTestMessage{FixedString: "123"}),
			invalid,
		}, {
			"InvalidFixedLenStringTooLong",
			getValMsg(ValTestMessage{FixedString: "12345"}),
			invalid,
		}, {
			"ValidRuneLenStringSB",
			getValMsg(ValTestMessage{RuneString: "123é"}), // len 4 in runes, 5 in bytes
			valid,
		}, {
			"ValidRuneLenStringNFD",
			getValMsg(ValTestMessage{RuneString: "123e\u0301"}), // this NFD string (e + ´) will be normalised to NFC string before len check so still 4 bytes
			valid,
		}, {
			"ValidRuneLenStringMB",
			getValMsg(ValTestMessage{RuneString: "AAA🌏"}), // will result in 7 bytes and 4 runes
			valid,
		}, {
			"ValidMinMaxLenStringL3",
			getValMsg(ValTestMessage{Title: "123"}), // len between 3 and 50
			valid,
		}, {
			"ValidMinMaxLenStringL4",
			getValMsg(ValTestMessage{Title: "1234"}), // len between 3 and 50
			valid,
		}, {
			"ValidMinMaxLenStringL49",
			getValMsg(ValTestMessage{Title: genString(49)}), // len between 3 and 50
			valid,
		}, {
			"ValidMinMaxLenStringL50",
			getValMsg(ValTestMessage{Title: genString(50)}), // len between 3 and 50
			valid,
		}, {
			"InvalidMinMaxLenStringL51",
			getValMsg(ValTestMessage{Title: genString(51)}), // len between 3 and 50
			invalid,
		}, {
			"InvalidMinMaxLenStringL2",
			getValMsg(ValTestMessage{Title: "12"}), // len between 3 and 50
			invalid,
		}, {
			"InvalidMinMaxLenStringL0",
			getValMsg(ValTestMessage{Title: emptyString}), // len between 3 and 50
			invalid,
		}, {
			"ValidReplaceUnsafe",
			getValMsg(ValTestMessage{ReplaceString: "<script>"}),
			valid,
		}, {
			"ValidTrimStringCustomTest",
			getValMsg(ValTestMessage{TrimString: " \t 1 2 3 \n\t"}), // expected result is 1 2 3
			valid,
		},
	}

	// Tests for emails
	for _, invalidEmail := range invalidEmails {
		tests = append(tests, TestSet{
			"InvalidEmail_" + invalidEmail,
			getValMsg(ValTestMessage{Email: invalidEmail}),
			invalid,
		})
	}
	for _, validEmail := range validEmails {
		tests = append(tests, TestSet{
			"ValidEmail_" + validEmail,
			getValMsg(ValTestMessage{Email: validEmail}),
			valid,
		})
	}

	// Safe string test payloads
	for _, input := range readFiles(invalidSafeStrings) {
		tests = append(tests, TestSet{
			"InvalidSafeString_" + input,
			getValMsg(ValTestMessage{Description: input}),
			invalid,
		})
	}

	// Replace unsafe option enabled:
	for _, input := range validReplaceUnsafeStrings {
		tests = append(tests, TestSet{
			"ValidReplaceUnsafeStrings_" + input,
			getValMsg(ValTestMessage{ReplaceString: input}),
			valid,
		})
	}
	// Replace unsafe option disabled:
	for _, input := range validReplaceUnsafeStrings {
		tests = append(tests, TestSet{
			"InvalidNotReplaceUnsafeStrings_" + input,
			getValMsg(ValTestMessage{NotReplaceString: input}),
			invalid,
		})
	}

	for _, input := range readFiles(validSafeStrings) {
		tests = append(tests, TestSet{
			"ValidSafeString_" + input,
			getValMsg(ValTestMessage{Description: input}),
			valid,
		})
	}

	// Names
	// k, fmt.Sprint(k)
	for _, input := range readFiles(validNames) {
		tests = append(tests, TestSet{
			"ValidName_" + input,
			getValMsg(ValTestMessage{Name: input}),
			valid,
		})
	}

	for _, input := range readFiles(testTitles) {
		tests = append(tests, TestSet{
			"ValidTitle_" + input,
			getValMsg(ValTestMessage{ScTitle: input}),
			valid,
		})
	}

	for _, input := range readFiles(testPermissive) {
		tests = append(tests, TestSet{
			"ValidTitlePermissive_" + input,
			getValMsg(ValTestMessage{ScPermissive: input}),
			valid,
		})
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.input.Validate()
			if test.shouldError && err == nil {
				t.Errorf("expected error for %v, but nil", test.name)
			}
			if !test.shouldError && err != nil {
				t.Errorf("unexpected error for %v: %v", test.name, err)
			}
			// Custom test for trim string
			// [whitespace]string string[whitespace] should result in "string string" (leading/trailing whitespace stripped)
			if test.name == "ValidTrimStringCustomTest" {
				if test.input.(*ValTestMessage).TrimString != "1 2 3" {
					t.Errorf("unexpected error for %v: %v", test.name, "unexpected value after trim")
				}
			}
			// if test.name == "InvalidBidi" {
			// 	// if strings.HasPrefix(test.name, "ValidPassword") {
			// 	// dump the message for debugging
			// 	fmt.Println("Validation error for", test.name, ":", err)
			// 	t.Errorf("Validation error for %v %v", test.name, err)
			// 	fmt.Println("Trim string: ")
			// 	fmt.Println(test.input)
			// }
		})
	}

}

func (m *ValTestMessage) getMsgField() *ValTestMessage {
	return m
}