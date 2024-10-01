package valtest

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
)

const (
	id               string = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
	legacyId         string = "56341C6E-35A7-4C97-9C5E-7AC79673EAB2"
	s12Id            string = "audit_f6dad1c9334040739b1e67ca70f4cf4d"
	legacyLongIdFail string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-9-00000190327E0675"     // length = 49 (without dashes)
	legacyLongId1    string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-90-00000190327E0675"    // length = 50 (without dashes)
	legacyLongId2    string = "005F2E38-8426-48AF-94DE-5FEA3A396EEA-891-00000153F68896DC"   // length = 51 (without dashes)
	legacyLongId3    string = "007B516E-53F1-4AA0-ABAF-8C78342A2C82-2388-00000221F1C2BD1E"  // length = 52 (without dashes)
	legacyLongId4    string = "00709A17-151F-4CFC-B412-F080343ED84D-11977-000010227B4C60A9" // length = 53 (without dashes)
	email            string = "email@example.com"
	password         string = "1234567!"
	name             string = "Γιώργος Νταλάρας"
	idv4             string = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
	uuidNotV4        string = "e07b6ac0-8a05-11e2-9951-ddd1182f65d9"
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
	"email.ending@dot.com.",
	"invalid.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa64@aa257aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com",
	"invalid.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa65@aa256aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com",
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
	// user part max 64 chars, domain part max 256 chars:
	"valid.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa64@aa256aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com",
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
	// "file:///testdata/NOCOMMIT_valid_strings.txt",
}

var testPermissive = []string{
	// "file:///testdata/NOCOMMIT_emoticon_strings.txt",
	// "file:///testdata/NOCOMMIT_actions-issues-PUA.txt",
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

var validURLs = []string{
	"https://example/", // technically valid - might not want to accept such format
	"https://example.com",
	"https://example.com/a",
	"https://example.com/index.html?a=b&c=d",
	"https://user@example.com",
	"https://user:cred@example.com", // not an actual credential
	"https://example.com/" + strings.Repeat("a", 900),
	"https://example.com/a+",   // trailing whitespace encoded
	"https://example.com/a%20", // trailing whitespace encoded
	"https://example.com/a?b=ftp%3A//user%3Acred%40example.com/file.html%3Fq%3D%27b%27%26test%3D%22something%2526%22", // several encoded and double-encoded URL-reserved chars
	"https://app.safetyculture.com/path/action?lang=en-US",
	"https://host.tld:88",
	"https://example.com:88/test",
	"https://xn--tlphone-byab/abc", // punnycode, technically valid - might not want to accept such format
}

var invalidURLs = []string{
	"https",
	"https:/",
	"https://",
	"https:///",
	"https:example.com",
	"https:/example.com",
	"https//:example.com",
	"https://example.com/" + strings.Repeat("a", 1000), // too long
	"ftp://example.com/",                               // default scheme is https
	"https://example.com/a#fragment",                   // fragment not allowed unless option enabled
	"https://example.com/\na",
	" https://example.com/a",  // leading whitespace
	"\thttps://example.com/a", // leading whitespace
	"https://example.com/a ",  // trailing whitespace not encoded, most browsers will just strip it but technically invalid
	"file:///etc/passwd",
	"javascript:alert(1)",
	"/absolute_url", // require URL (with scheme, domain), not URI
	"./relative_url",
	"../relative/url",
	"https://example.com/λά.html",              // invalid runes in path
	"https://invalid-unicodé-domainλά.com/abc", // invalid runes in domain
}

var valMsg = ValTestMessage{
	Id:                    id,
	Ids:                   []string{id, id},
	LegacyId:              legacyId,
	S12Id:                 s12Id,
	MediaId:               id,
	InnerLegacyId:         &InnerMessageWithLegacyId{Id: id},
	InnerS12Id:            &InnerMessageWithS12Id{Id: s12Id},
	Uuid:                  id,
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
	ContactOneof:          &ValTestMessage_Phone{Phone: "14574560123"},
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
	ContactsWithLengthConstraint: []*ValTestMessage_Contact{
		{Phone: "abc", Email: "test@example.com"},
		{Phone: "", Email: "test2@example.com"},
	},
	Url:        "https://example.com/test",
	UrlAllOpts: "http://app.safetyculture.com/report/media?param=test#fragment",
	// NotSupported: ,
	Timezone:   "Australia/Sydney",
	LongString: strings.Repeat("x", 30000),
}

// omit optional fields here
var valMsgOpts = ValTestMessage{
	Id:                    id,
	Ids:                   []string{id, id},
	LegacyId:              legacyId,
	S12Id:                 s12Id,
	InnerLegacyId:         &InnerMessageWithLegacyId{Id: legacyId},
	InnerS12Id:            &InnerMessageWithS12Id{Id: s12Id},
	Uuid:                  id,
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
	ContactOneof:          &ValTestMessage_Phone{Phone: "14574560123"},
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
	ContactsWithLengthConstraint: []*ValTestMessage_Contact{
		{Phone: "abc", Email: "test@example.com"},
	},
	Url: "https://example.com/test",
	// NotSupported: ,
	Timezone: "Australia/Sydney",
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
	if m.S12Id != "" {
		newMsg.S12Id = replaceEmpty(m.S12Id)
	}
	if m.AllId != "" {
		newMsg.AllId = replaceEmpty(m.AllId)
	}
	if m.MediaId != "" {
		newMsg.MediaId = replaceEmpty(m.MediaId)
	}
	if m.Uuid != "" {
		newMsg.Uuid = replaceEmpty(m.Uuid)
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
	if m.OptionalString != nil {
		newMsg.OptionalString = m.OptionalString
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
	if m.Url != "" {
		newMsg.Url = replaceEmpty(m.Url)
	}
	if m.UrlAllOpts != "" {
		newMsg.UrlAllOpts = replaceEmpty(m.UrlAllOpts)
	}
	if m.Timezone != "" {
		newMsg.Timezone = replaceEmpty(m.Timezone)
	}
	if m.TimezoneOptional != "" {
		newMsg.TimezoneOptional = replaceEmpty(m.TimezoneOptional)
	}
	if m.LongString != "" {
		newMsg.LongString = replaceEmpty(m.LongString)
	}
	return &newMsg
}

func TestValidationRules(t *testing.T) {

	sb := strings.Builder{}
	for i := 0; i < 2001; i++ {
		sb.WriteRune('a')
	}

	emptyStr := ""
	optionalStr := "optional"

	type TestSet struct {
		name        string
		input       proto.Validator
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
			"InvalidIdWithUUIDNotV4",
			getValMsg(ValTestMessage{Id: uuidNotV4}),
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
			"ValidS12ID",
			getValMsg(ValTestMessage{S12Id: s12Id}),
			valid,
		},
		{
			"InvalidS12ID",
			getValMsg(ValTestMessage{S12Id: "fake_id"}),
			invalid,
		},
		{
			"InvalidUUIDwithS12Id",
			getValMsg(ValTestMessage{Id: s12Id}),
			invalid,
		},
		{
			"InvalidUUIDwithLegacyId",
			getValMsg(ValTestMessage{Id: legacyId}),
			invalid,
		},
		{
			"InvalidLegacyIdwithS12Id",
			getValMsg(ValTestMessage{LegacyId: s12Id}),
			invalid,
		},
		{
			"InvalidS12IdWithLegacyId",
			getValMsg(ValTestMessage{S12Id: legacyId}),
			invalid,
		},
		{
			"ValidIdAllOptsUUID",
			getValMsg(ValTestMessage{AllId: id}),
			valid,
		},
		{
			"ValidIdAllOptsLegacyId",
			getValMsg(ValTestMessage{AllId: legacyId}),
			valid,
		},
		{
			"ValidIdAllOptsS12Id",
			getValMsg(ValTestMessage{AllId: s12Id}),
			valid,
		},
		{
			"ValidS12IDUppercase",
			getValMsg(ValTestMessage{S12Id: "audit_C6C011ED4ADE460DA04BDA730834B667"}),
			valid,
		},
		{
			"ValidIDAnyVersionWithUUIDv4",
			getValMsg(ValTestMessage{Uuid: id}),
			valid,
		},
		{
			"ValidIDAnyVersionWithUUIDNotV4",
			getValMsg(ValTestMessage{Uuid: uuidNotV4}),
			valid,
		},
		{
			"InvalidIDAnyVersionWithInvalidId",
			getValMsg(ValTestMessage{Uuid: "invalid"}),
			invalid,
		},
		{
			"ValidIdAllOptsLongPrefixedLegacyId50",
			getValMsg(ValTestMessage{AllId: "audit_1CC51EFA600C4F6285B9652A32D714D69500000013D6801111"}),
			valid,
		},
		{
			"ValidIdAllOptsLongPrefixedLegacyId51",
			getValMsg(ValTestMessage{AllId: "audit_0CEE3E70B3584D40B1E7B4A485C19D8A23500000154D3322222"}),
			valid,
		},
		{
			"ValidIdAllOptsLongPrefixedLegacyId52",
			getValMsg(ValTestMessage{AllId: "audit_072F8B1B16594188A8B906ACED430CA889670000153EA8333333"}),
			valid,
		},
		{
			"ValidIdAllOptsLongPrefixedLegacyId53",
			getValMsg(ValTestMessage{AllId: "template_5033F701076E47FCA698F84E05557A6D2024300002267A4444444"}),
			valid,
		},
		{
			"InValidIdAllOptsLongPrefixedLegacyId54",
			getValMsg(ValTestMessage{AllId: "template_5033F701076E47FCA698F84E05557A6D2024300002267A44444445"}),
			invalid,
		},
		{
			"InValidIdAllOptsLongPrefixedLegacyId44",
			getValMsg(ValTestMessage{AllId: "template_5033F7A698F84E05557A6D2024300002267A44444445"}),
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
		}, {
			"ValidNilOptionalString",
			getValMsg(ValTestMessage{OptionalString: nil}),
			valid,
		}, {
			"InvalidEmptyOptionalString",
			getValMsg(ValTestMessage{OptionalString: &emptyStr}),
			invalid,
		}, {
			"ValidNonEmptyOptionalString",
			getValMsg(ValTestMessage{OptionalString: &optionalStr}),
			valid,
		}, {
			// Custom test for deeply nested messages; these should still be validated, expecting an error at level 3
			"CustomNestedMessage",
			&MyReqMessage{
				UserId: "12",
				OrgNested: &NestedLevel1Message{
					OrgId3: "123",
					OrgNested: &NestedLevel2Message{
						OrgId4: "1234",
						OrgNested: &NestedLevel3Message{
							OrgId5: "123456", // length should be 5, 6 given
						},
					},
				},
			},
			invalid,
		}, {
			"ScimUser",
			&ScimUser{
				Emails: []*ScimEmail{
					{Value: "valid@example.com"},
					{Value: "invalid_email"},
				},
			},
			invalid,
		}, {
			"ScimUserEmpty",
			&ScimUser{},
			valid,
		},
		{
			"enum_required is failed",
			&MyMessageWithEnum{},
			invalid,
		},
		{
			"enum_required is passed",
			&MyMessageWithEnum{
				Enum: MyMessageWithEnum_MY_ENUM_FIRST,
			},
			valid,
		},
		{
			"repeated enum_required is failed",
			&MyMessageWithRepeatedEnum{
				Enums: []MyMessageWithRepeatedEnum_MyEnum{
					MyMessageWithRepeatedEnum_MY_ENUM_UNSPECIFIED,
				},
			},
			invalid,
		},
		{
			"repeated enum_required is passed",
			&MyMessageWithRepeatedEnum{
				Enums: []MyMessageWithRepeatedEnum_MyEnum{
					MyMessageWithRepeatedEnum_MY_ENUM_FIRST,
				},
			},
			valid,
		},
		{
			"my_second_field validation is failed",
			&MyOneOfMsg{
				MyField: &MyOneOfMsg_MySecondField{
					MySecondField: &MyOneOfMsg_SecondType{
						Value: 2,
					},
				},
			},
			invalid,
		},
		{
			"my_second_field validation is passed",
			&MyOneOfMsg{
				MyField: &MyOneOfMsg_MySecondField{
					MySecondField: &MyOneOfMsg_SecondType{
						Value: 3,
					},
				},
			},
			valid,
		},
		{
			"my_first_field validation is passed",
			&MyOneOfMsg{
				MyField: &MyOneOfMsg_MyFirstField{
					MyFirstField: &MyOneOfMsg_FirstType{
						Value: 2,
					},
				},
			},
			valid,
		},
		{
			"my_first_field validation is failed",
			&MyOneOfMsg{
				MyField: &MyOneOfMsg_MyFirstField{
					MyFirstField: &MyOneOfMsg_FirstType{
						Value: 0,
					},
				},
			},
			invalid,
		},
		{
			"my_third_field validation is failed",
			&MyOneOfMsg{
				MyField: &MyOneOfMsg_MyThirdField{
					MyThirdField: "",
				},
			},
			invalid,
		},
		{
			"my_third_field validation is passed",
			&MyOneOfMsg{
				MyField: &MyOneOfMsg_MyThirdField{
					MyThirdField: "abc",
				},
			},
			valid,
		}, {
			"ValidUrlCustomSchemeFtp",
			getValMsg(ValTestMessage{UrlAllOpts: "ftp://example.com/abc"}),
			valid,
		}, {
			"ValidUrlCustomSchemeFtps",
			getValMsg(ValTestMessage{UrlAllOpts: "ftps://example.com/abc#fragment"}),
			valid,
		}, {
			"InvalidUrlCustomSchemeHttps",
			getValMsg(ValTestMessage{UrlAllOpts: "https://example.com/abc"}),
			invalid,
		}, {
			"InvalidTimezone",
			getValMsg(ValTestMessage{Timezone: "Australia/BonnieDoon"}),
			invalid,
		}, {
			"InvalidTimezoneOptional",
			getValMsg(ValTestMessage{TimezoneOptional: "Australia/BonnieDoon"}),
			invalid,
		},
	}

	// Custom test for repeated contact
	contactsMsgValid1 := valMsg
	contactsMsgValid1.ContactsWithLengthConstraint = []*ValTestMessage_Contact{
		{Phone: "1 message", Email: "test1@example.com"},
	}
	contactsMsgValid10 := valMsg
	contactsMsgValid10.ContactsWithLengthConstraint = []*ValTestMessage_Contact{
		{Phone: "10 messages", Email: "test1@example.com"},
		{Phone: "", Email: "test2@example.com"},
		{Phone: "", Email: "test3@example.com"},
		{Phone: "", Email: "test4@example.com"},
		{Phone: "", Email: "test5@example.com"},
		{Phone: "", Email: "test6@example.com"},
		{Phone: "", Email: "test7@example.com"},
		{Phone: "", Email: "test8@example.com"},
		{Phone: "", Email: "test9@example.com"},
		{Phone: "", Email: "test10@example.com"},
	}
	contactsMsgInvalid0 := valMsg
	contactsMsgInvalid0.ContactsWithLengthConstraint = []*ValTestMessage_Contact{}
	contactsMsgInvalid11 := valMsg
	contactsMsgInvalid11.ContactsWithLengthConstraint = []*ValTestMessage_Contact{
		{Phone: "11 messages", Email: "test1@example.com"},
		{Phone: "", Email: "test2@example.com"},
		{Phone: "", Email: "test3@example.com"},
		{Phone: "", Email: "test4@example.com"},
		{Phone: "", Email: "test5@example.com"},
		{Phone: "", Email: "test6@example.com"},
		{Phone: "", Email: "test7@example.com"},
		{Phone: "", Email: "test8@example.com"},
		{Phone: "", Email: "test9@example.com"},
		{Phone: "", Email: "test10@example.com"},
		{Phone: "", Email: "test11@example.com"},
	}
	contactsMsgInvalidEmail := valMsg
	contactsMsgInvalidEmail.ContactsWithLengthConstraint = []*ValTestMessage_Contact{
		{Phone: "", Email: "test1@invalid"},
	}
	contactsMsgNoConstraintsValid0 := valMsg
	contactsMsgNoConstraintsValid0.ContactsWithoutLengthConstraint = []*ValTestMessage_Contact{}
	contactsMsgNoConstraintsValid2 := valMsg
	contactsMsgNoConstraintsValid2.ContactsWithoutLengthConstraint = []*ValTestMessage_Contact{
		{Phone: "2 messages", Email: "test1@example.com"},
		{Phone: "", Email: "test2@example.com"},
	}
	tests = append(tests, TestSet{
		"ContactWithLengthConstraintValid1",
		&contactsMsgValid1,
		valid,
	}, TestSet{
		"ContactWithLengthConstraintValid10",
		&contactsMsgValid10,
		valid,
	}, TestSet{
		"ContactWithLengthConstraintInvalid0",
		&contactsMsgInvalid0,
		invalid,
	}, TestSet{
		"ContactWithLengthConstraintInvalid11",
		&contactsMsgInvalid11,
		invalid,
	}, TestSet{
		"ContactWithLengthConstraintInvalidEmail",
		&contactsMsgInvalidEmail,
		invalid,
	}, TestSet{
		"ContactWithoutLengthConstraintValid0",
		&contactsMsgNoConstraintsValid0,
		valid,
	}, TestSet{
		"ContactWithoutLengthConstraintValid2",
		&contactsMsgNoConstraintsValid0,
		valid,
	})

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

	// URLs
	for _, invalidURL := range invalidURLs {
		tests = append(tests, TestSet{
			"InvalidURL_" + invalidURL,
			getValMsg(ValTestMessage{Url: invalidURL}),
			invalid,
		})
	}
	for _, validURL := range validURLs {
		tests = append(tests, TestSet{
			"ValidURL_" + validURL,
			getValMsg(ValTestMessage{Url: validURL}),
			valid,
		})
	}

	emptyTimezoneMsg := getValMsg(ValTestMessage{})
	emptyTimezoneMsg.Timezone = ""
	// Empty Timezone
	tests = append(tests, TestSet{
		"InvalidTimezoneRequired",
		emptyTimezoneMsg,
		invalid,
	})
	tests = append(tests, TestSet{
		"InvalidLongString",
		getValMsg(ValTestMessage{LongString: strings.Repeat("y", 30002)}),
		invalid,
	})
	fmt.Println("###### LEN = ", len(strings.Repeat("y", 30002)))

	// RejectUrl
	rejectUrlTestUrlsToReject := map[string]string{
		"[RejectUrl] Text is only a full URL":                   "http://example.com",
		"[RejectUrl] Text with URL and query params":            "https://example.com/path?query=test#hello",
		"[RejectUrl] Text with URL and non-standard characters": "Claim your bonus http://example.com/xyz:D",
		"[RejectUrl] Text with URL with emoji":                  "Claim your bonus http://example.com/xyz😊",
		"[RejectUrl] Text with multiple URLs":                   "Click https://example.com/here and http://example.com.",
	}

	for testName, input := range rejectUrlTestUrlsToReject {
		tests = append(tests, TestSet{
			name: testName,
			input: &NonUrlMessage{
				RejectUrlTest: input,
			},
			shouldError: invalid,
		})
	}
	RejectUrlTestUrlToAllow := map[string]string{
		"[RejectUrl] Text with only partial URL":           "example.com",
		"[RejectUrl] Text with partial URL and other text": "Hello example.com!",
		"[RejectUrl] Text without any URL parts":           "This is a sentence. And this is another sentence!",
	}
	for testName, input := range RejectUrlTestUrlToAllow {
		tests = append(tests, TestSet{
			name: testName,
			input: &NonUrlMessage{
				RejectUrlTest: input,
			},
			shouldError: valid,
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

// TestBreakPartialUrl is a special test case that verifies that a message contents is modified, not that an error is
// returned
func TestBreakPartialUrl(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"Does not modify text without URLs": {
			input:    "This is a normal sentence. It has no URLs",
			expected: "This is a normal sentence. It has no URLs",
		},
		"Break partial URLs in text": {
			input:    "Check out example.com and another.example.com for more info.",
			expected: "Check out example. com and another. example. com for more info.",
		},
		"Handles multiple URLs in a single line": {
			input:    "Visit site1.com, site2.org, and site3.net for examples.",
			expected: "Visit site1. com, site2. org, and site3. net for examples.",
		},
		"Also affects email addresses": {
			input:    "Contact me at user@example.com or support@company.org.",
			expected: "Contact me at user@example. com or support@company. org.",
		},
		"Handles URLs at the end of sentences": {
			input:    "Check this website: example.com. Then visit another.site.",
			expected: "Check this website: example. com. Then visit another. site.",
		},
		"Handles multiple subdomains, spacing after each dot": {
			input:    "Look at https://multiple-subdomains.gov.qld.au",
			expected: "Look at https://multiple-subdomains. gov. qld. au",
		},
		"Handles empty inputs": {
			input:    "",
			expected: "",
		},
		"Handles input with only URLs": {
			input:    "example.com another.example.com",
			expected: "example. com another. example. com",
		},
		"Does not affect numbers or other punctuation": {
			input:    "Pi is about 3.14159. Visit example.com for more math!",
			expected: "Pi is about 3.14159. Visit example. com for more math!",
		},
		"Does not affect strings that already has spaces afterwards": {
			input:    "Already. Split",
			expected: "Already. Split",
		},
		"Does not affect a split up URL": {
			input:    "https:// split-url .com",
			expected: "https:// split-url .com",
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			msg := &NonUrlMessage{
				BreakPartialUrlTest: test.input,
			}
			require.NoError(t, msg.Validate())
			assert.Equal(t, test.expected, msg.GetBreakPartialUrlTest())
		})
	}
}

func (m *ValTestMessage) getMsgField() *ValTestMessage {
	return m
}

func genLogOnlyValidationMessage() *LogOnlyValidationMessage {
	return &LogOnlyValidationMessage{
		ImageId:      id,
		InspectionId: id,
		OwnerId:      id,
		Title:        "123",
		Name:         "123",
		Latitude:     0,
	}
}

func genLowercaseValidationMessage() *LowercaseValidationMessage {
	return &LowercaseValidationMessage{
		Uuidv4:        "92b6c2f9-abd8-48bc-a2c9-bf70e969751a",
		Uuidv4LogOnly: "92b6c2f9-abd8-48bc-a2c9-bf70e969751a",
		S12Id:         "audit_f6dad1c9334040739b1e67ca70f4cf4d",
		S12IdLogOnly:  "audit_f6dad1c9334040739b1e67ca70f4cf4d",
		Legacy:        "92b6c2f9abd848bca2c9bf70e969751a",
		LegacyLogOnly: "92b6c2f9abd848bca2c9bf70e969751a",
		RevId:         "",
		QuestionId:    "92b6c2f9-abd8-48bc-a2c9-bf70e969751a",
	}
}

func TestSoftValidation_Validate(t *testing.T) {
	tests := []struct {
		name      string
		msg       *LogOnlyValidationMessage
		shouldErr bool
	}{
		{
			name: "should fail imageId when passing bad format UUID",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.ImageId = "c1c01a8f-f724-42bf-ac6f-5478a0f1292x"
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "should fail inspectionId when passing bad format UUID",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "c1c01a8f-f724-42bf-ac6f-5478a0f1292x"
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "should not fail ownerId when passing bad format UUID and softValidation is enabled",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.OwnerId = "c1c01a8f-f724-42bf-ac6f-5478a0f1292x"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should not fail when the input exceeds 50 characters",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.OwnerId = "c1c01a8f-f724-42bf-ac6f-5478a0f1292xc1c01a8f-f724-42bf-ac6f-5478a0f1292xc1c01a8f-f724-42bf-ac6f-5478a0f1292xc1c01a8f-f724-42bf-ac6f-5478a0f1292x"
				return m
			}(),
			shouldErr: false,
		},
		{
			name:      "should pass with legitimate UUIDs",
			msg:       genLogOnlyValidationMessage(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate s12ID",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "audit_401d39e12a3c4d8b8021631e63e82492"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate uppercase s12ID",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "audit_401D39E12A3C4D8B8021631E63E82492"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate uppercased UUID",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "401D39E1-2A3C-4D8B-8021-631E63E82492"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate UUID without hyphens",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "401d39e12a3c4d8b8021631e63e82492"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate uppercased UUID without hyphens",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "401D39E12A3C4D8B8021631E63E82492"
				return m
			}(),
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.shouldErr == (err == nil) {
				t.Errorf("%s, supposed to return an error", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("%s, supposed not to return an error, but we received %v", tt.name, err)
			}
		})
	}
}

func TestSoftValidation_ValidateLatitude(t *testing.T) {
	tests := []struct {
		name      string
		msg       *LogOnlyValidationMessage
		shouldErr bool
	}{
		{
			name: "should pass with legitimate positive integer in range latitude",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Latitude = 5
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate negative integer in lower-range latitude",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Latitude = -90
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate negative double in range latitude",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Latitude = -88.8888888888
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate positive double in lower-range latitude",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Latitude = -90.00
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate zero double in latitude",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Latitude = -0.00
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate positive double in upper-range latitude",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Latitude = 90.00
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should fail with positive double higher than upper-range latitude",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Latitude = -91.05
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "should not fail with invalid double higher than upper-range latitude when log-only",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Longitude = -190.12233444
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate zero in range accuracy",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Accuracy = 0
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should pass with legitimate value in upper-range accuracy",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Accuracy = 10_000
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "should not fail with value out-range accuracy",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Accuracy = 10_001
				return m
			}(),
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.shouldErr == (err == nil) {
				t.Errorf("%s, supposed to return an error", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("%s, supposed not to return an error, but we received %v", tt.name, err)
			}
		})
	}
}

func genNumberMessage() *NumberMessage {
	return &NumberMessage{
		NanAllowed:         math.NaN(),
		NanDisallowed:      0,
		OptionalOnly:       0,
		OptionalNoNanValue: 0,
		DefaultNan:         0,
		RangeBasic:         0,
		RangeLow:           0,
		RangeHigh:          0,
		RangeNovalues:      0,
		RangeNotOptional:   5,
		IntTest:            19,
	}
}
func TestNumberValidation_Validate(t *testing.T) {
	tests := []struct {
		name      string
		msg       *NumberMessage
		shouldErr bool
	}{
		{
			"number allows nan",
			func() *NumberMessage {
				m := genNumberMessage()
				m.NanAllowed = math.NaN()
				return m
			}(),
			valid,
		},
		{
			"number allows value when optional",
			func() *NumberMessage {
				m := genNumberMessage()
				m.OptionalNoNanValue = 0.1
				return m
			}(),
			valid,
		},
		{
			"number allows valid value when allow_nan is true",
			func() *NumberMessage {
				m := genNumberMessage()
				m.NanAllowed = 0.1
				return m
			}(),
			valid,
		},
		{
			"number allows valid value when allow_nan is false",
			func() *NumberMessage {
				m := genNumberMessage()
				m.NanDisallowed = 0.1
				return m
			}(),
			valid,
		},
		{
			"number disallows nan when allow_nan is false",
			func() *NumberMessage {
				m := genNumberMessage()
				m.NanDisallowed = math.NaN()
				m.OptionalNoNanValue = math.NaN()
				m.OptionalOnly = math.NaN()
				return m
			}(),
			invalid,
		},
		{
			"number disallows nan when optional",
			func() *NumberMessage {
				m := genNumberMessage()
				m.OptionalNoNanValue = math.NaN()
				m.OptionalOnly = math.NaN()
				return m
			}(),
			invalid,
		},
		{
			"float/double range allows value in range",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeBasic = 2
				return m
			}(),
			valid,
		},
		{
			"float/double range disallows value below range",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeBasic = -1
				return m
			}(),
			invalid,
		},
		{
			"float/double range disallows value above range",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeBasic = 11
				return m
			}(),
			invalid,
		},
		{
			"float/double range validates zero when not optional",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeNotOptional = 0
				return m
			}(),
			invalid,
		},
		{
			"float/double range validates lower bounds",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeLow = -1
				return m
			}(),
			invalid,
		},
		{
			"float/double range validates lower bounds inclusive",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeLow = 1
				return m
			}(),
			valid,
		},
		{
			"float/double range validates upper bounds",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeHigh = 11
				return m
			}(),
			invalid,
		},
		{
			"float/double range validates upper bounds inclusive",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeHigh = 10
				return m
			}(),
			valid,
		},
		{
			"float/double range validates lower bounds with -Infinity",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeLow = math.Inf(-1)
				return m
			}(),
			invalid,
		},
		{
			"float/double range validates upper bounds with -Infinity",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RangeHigh = math.Inf(1)
				return m
			}(),
			invalid,
		},
		{
			"int range allows value in range",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = 50
				return m
			}(),
			valid,
		},
		{
			"int range disallows value below range",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = -50
				return m
			}(),
			invalid,
		},
		{
			"int range disallows value above range",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = 10000
				return m
			}(),
			invalid,
		},
		{
			"int range validates zero when not optional",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = 0
				return m
			}(),
			invalid,
		},
		{
			"int range validates lower bounds",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = -1
				return m
			}(),
			invalid,
		},
		{
			"int range validates lower bounds inclusive",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = 1
				return m
			}(),
			valid,
		},
		{
			"int range validates upper bounds",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = 100
				return m
			}(),
			invalid,
		},
		{
			"int range validates upper bounds inclusive",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = 99
				return m
			}(),
			valid,
		},
		{
			"int range validates lower bounds with MaxInt32",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = math.MaxInt32
				return m
			}(),
			invalid,
		},
		{
			"int range validates upper bounds with -MaxInt32",
			func() *NumberMessage {
				m := genNumberMessage()
				m.IntTest = math.MinInt32
				return m
			}(),
			invalid,
		},
		{
			"int64 range validates lower bounds inclusive valid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.Int64Test = -100
				return m
			}(),
			valid,
		},
		{
			"int64 range validates upper bounds inclusive valid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.Int64Test = 100
				return m
			}(),
			valid,
		},
		{
			"int64 range validates lower bounds inclusive invalid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.Int64Test = -101
				return m
			}(),
			invalid,
		},
		{
			"int64 range validates upper bounds inclusive invalid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.Int64Test = 101
				return m
			}(),
			invalid,
		},
		{
			"uint range validates lower bounds inclusive valid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.UintTest = 5
				return m
			}(),
			valid,
		},
		{
			"uint range validates upper bounds inclusive valid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.UintTest = 99
				return m
			}(),
			valid,
		},
		{
			"uint range validates lower bounds inclusive invalid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.UintTest = 1
				return m
			}(),
			invalid,
		},
		{
			"uint range validates upper bounds inclusive invalid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.UintTest = 100
				return m
			}(),
			invalid,
		},
		{
			"repeated int range validates bounds invalid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RepeatedInt = []int32{1, 2, 3}
				return m
			}(),
			invalid,
		},
		{
			"repeated int range validates bounds valid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RepeatedInt = []int32{10, 15, 20}
				return m
			}(),
			valid,
		},
		{
			"repeated int range validates bounds single invalid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.RepeatedInt = []int32{10, 115, 20}
				return m
			}(),
			invalid,
		},
		{
			"nested int range validates bounds valid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.NestedNumber = &NumberMessage_NestedNumber{Value: 1}
				return m
			}(),
			valid,
		},
		{
			"nested int range validates bounds invalid",
			func() *NumberMessage {
				m := genNumberMessage()
				m.NestedNumber = &NumberMessage_NestedNumber{Value: 9999999}
				return m
			}(),
			invalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.shouldErr && err == nil {
				t.Errorf("%s, supposed to return an error", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("%s, supposed not to return an error, but we received %v", tt.name, err)
			}
		})
	}
}

func TestSoftValidation_ValidateS12Strict(t *testing.T) {
	tests := []struct {
		name      string
		msg       *LogOnlyValidationMessage
		shouldErr bool
	}{
		{
			name: "Should pass with AUDIT",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "audit_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with TEMPLATE",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "template_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with USER",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "user_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with ACTION",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "action_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with NTFMSG",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "ntfmsg_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with EVIDENCE",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "evidence_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with ROLE",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "role_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with LOCATION",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "location_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with RESPONSESET",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "responseset_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with RESPONSE",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "response_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with PREFERENCE",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "preference_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with HEADS_UP",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "heads_up_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with SUBSCRIPTION",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "subscription_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with FOLDER",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "folder_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with SCHEDULEITEM",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "scheduleitem_d848b3a0a1e14dbb8c01fcae8cc2ce12"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should fail with CHICKEN",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.InspectionId = "chicken_a109b645119b4eacb98ace52cce79fa7"
				return m
			}(),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.shouldErr == (err == nil) {
				t.Errorf("%s, supposed to return an error", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("%s, supposed not to return an error, but we received %v", tt.name, err)
			}
		})
	}
}

func TestSoftValidation_ValidateString(t *testing.T) {
	tests := []struct {
		name      string
		msg       *LogOnlyValidationMessage
		shouldErr bool
	}{
		{
			name: "Should pass when size is correct",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Title = "Hey!!"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass when size is incorrect and log-only",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Title = "Hey There!!!"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass when size is correct",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Name = "Hey !"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should fail when size is incorrect",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Name = "Hey !!!!"
				return m
			}(),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.shouldErr == (err == nil) {
				t.Errorf("%s, supposed to return an error", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("%s, supposed not to return an error, but we received %v", tt.name, err)
			}
		})
	}
}

func TestValidation_ValidateLowercaseIDs(t *testing.T) {
	tests := []struct {
		name      string
		msg       *LowercaseValidationMessage
		shouldErr bool
	}{
		{
			name: "Should pass when UUID case is lowercase",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.Uuidv4 = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should not pass UUID when case is uppercase",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.Uuidv4 = "92B6C2F9-ABD8-48BC-A2C9-BF70E969751A"
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "Should pass UUID when case is uppercase and log-only enabled",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.Uuidv4LogOnly = "92B6C2F9-ABD8-48BC-A2C9-BF70E969751A"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass when S12ID case is lowercase",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.S12Id = "audit_92b6c2f9abd848bca2c9bf70e969751a"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should not pass S12ID when case is uppercase",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.S12Id = "audit_92B6C2F9ABD848BCA2C9BF70E969751A"
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "Should pass S12ID when case is uppercase and log-only enabled",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.S12IdLogOnly = "audit_92B6C2F9ABD848BCA2C9BF70E969751A"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass when Legacy case is lowercase",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.Legacy = "92b6c2f9abd848bca2c9bf70e969751a"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should not pass Legacy when case is uppercase",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.Legacy = "92B6C2F9ABD848BCA2C9BF70E969751A"
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "Should pass Legacy when case is uppercase and log-only enabled",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.S12IdLogOnly = "92B6C2F9ABD848BCA2C9BF70E969751A"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with empty value",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.RevId = ""
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with proper value",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.RevId = "4-f6fd91720a594e66a6deb9874aca5fa1"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should fail with invalid value",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.RevId = "f6fd91720a594e66a6deb9874aca5fa1"
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "Should pass with real legacy value #1",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.QuestionId = "e0590e4b-b63b-42f1-b225-4cfcf0afeb9c-101-00000006f"
				return m
			}(),
			shouldErr: false,
		}, {
			name: "Should pass with a real legacy value #2",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.QuestionId = "f3545d44-ea77-11e1-aff1-0800200c9a66"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass with a real legacy value #3",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.QuestionId = "E7497185-5BED-4CDA-B293-81CCF13FEE31-1425-00000050"
				return m
			}(),
			shouldErr: false,
		}, {
			name: "Should pass with a real legacy value #4",
			msg: func() *LowercaseValidationMessage {
				m := genLowercaseValidationMessage()
				m.QuestionId = "ca1cb1f1-d2df-414a-bf4f-8d9a65a52db3-18665-000068a"
				return m
			}(),
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.shouldErr == (err == nil) {
				t.Errorf("%s, supposed to return an error", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("%s, supposed not to return an error, but we received %v", tt.name, err)
			}
		})
	}
}

func TestSoftValidation_ValidateUnsafeString(t *testing.T) {
	tests := []struct {
		name      string
		msg       *LogOnlyValidationMessage
		shouldErr bool
	}{
		{
			name: "Should pass when size is correct",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Title = "Hey!!"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass when size is correct with accents",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Title = "èéśș"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass when size is incorrect and log-only",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Title = "Hey There!!!"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass case when string is optional log-only and falls within the length constraints",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Answer = "Dean O’Callaghan "
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass case when string is optional and empty",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Answer = ""
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass case when string is optional log-only but is greater than max allowed",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Answer = "● Evidence of rat activities were recorded in fe..///...Re Branded 8x V5, Relocated 4x V5 and cladded 4x V"
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should pass case when string is optional and falls within the length constraints",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Note = "Dean O’Callaghan "
				return m
			}(),
			shouldErr: false,
		},
		{
			name: "Should fail case when string is optional and less than 5 min declared",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Note = "123"
				return m
			}(),
			shouldErr: true,
		},
		{
			name: "Should fail case when string is optional and greater than 40 min declared",
			msg: func() *LogOnlyValidationMessage {
				m := genLogOnlyValidationMessage()
				m.Note = "● Evidence of rat activities were recorded in fe..///...Re Branded 8x V5, Relocated 4x V5 and cladded 4x V"
				return m
			}(),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.shouldErr == (err == nil) {
				t.Errorf("%s, supposed to return an error", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("%s, supposed not to return an error, but we received %v", tt.name, err)
			}
		})
	}
}
