// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package proto

import (
	"regexp"
	"strings"
)

const (
	// UUIDSize is the size of a UUID in bytes.
	UUIDSize int = 16
)

const (
	uuid     string = "^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$"
	legacyId string = "(?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}(-?[0-9a-f]{2,5}-?[0-9a-f]{16})?$"
	s12Id    string = "^[a-z]*_([0-9a-fA-F]){32}$"
	rePua    string = `[\x{E000}-\x{F8FF}]` // Private Use Codepoints in the Basic Multilingual Plane (not including planes 15, 16)
	// govalidator package used before seems unmaintained at the moment and we needed changes to the regex so copied it from https://github.com/asaskevich/govalidator/blob/f21760c49a8d602d863493de796926d2a5c1138d/patterns.go#L77
	// ensure all checks are performmed in IsValidEmail so we can revert this later if we want to
	reEmail string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	reUUID4 string = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
)

var (
	rxUUID                      = regexp.MustCompile(uuid)
	rxLegacyId                  = regexp.MustCompile(legacyId)
	rxS12Id                     = regexp.MustCompile(s12Id)
	UnsafeCharReplacer          = strings.NewReplacer("\u0021", "\uFF01", "\u0022", "\u201D", "\u0023", "\u0020", "\u0025", "\u2052", "\u0026", "\u0020", "\u0027", "\u2019", "\u002A", "\u2217", "\u002B", "\uFF0B", "\u002D", "\u2212", "\u002F", "\u2215", "\u003B", "\u037E", "\u003C", "\u02C2", "\u003D", "\u2E40", "\u003E", "\u02C3", "\u005C", "\uFF3C", "\u0060", "\u2019", "\u007C", "\uFFE8", "\u3164", "\u0020")
	SymbolCharReplacer          = strings.NewReplacer("\u00A0", "\u0020", "\u1680", "\u0020", "\u2000", "\u0020", "\u2001", "\u0020", "\u2002", "\u0020", "\u2003", "\u0020", "\u2004", "\u0020", "\u2005", "\u0020", "\u2006", "\u0020", "\u2007", "\u0020", "\u2008", "\u0020", "\u2009", "\u0020", "\u200A", "\u0020", "\u202F", "\u0020", "\u205F", "\u0020", "\u3000", "\u0020", "\u200C", "", "\u200D", "", "\uFEFF", "", "\u2014", "\u2013", "\u2018", "\u2019", "\u3002", "\u002E", "\uFF0C", "\u002C", "\uFF1A", "\u003A", "\u0009", "\u0020", "\u000A", "\u0020", "\u000D", "\u0020")
	SymbolCharReplacerMultiline = strings.NewReplacer("\u00A0", "\u0020", "\u1680", "\u0020", "\u2000", "\u0020", "\u2001", "\u0020", "\u2002", "\u0020", "\u2003", "\u0020", "\u2004", "\u0020", "\u2005", "\u0020", "\u2006", "\u0020", "\u2007", "\u0020", "\u2008", "\u0020", "\u2009", "\u0020", "\u200A", "\u0020", "\u202F", "\u0020", "\u205F", "\u0020", "\u3000", "\u0020", "\u200C", "", "\u200D", "", "\uFEFF", "", "\u2014", "\u2013", "\u2018", "\u2019", "\u3002", "\u002E", "\uFF0C", "\u002C", "\uFF1A", "\u003A", "\u0009", "\u0020", "\u000D", "\u0020")
	RegexPua                    = regexp.MustCompile(rePua)
	RegexEmail                  = regexp.MustCompile(reEmail)
	RegexUUID4                  = regexp.MustCompile(reUUID4)
)

// IsUUID checks if the string is a UUID (version 3, 4 or 5).
func IsUUID(str string) bool {
	return rxUUID.MatchString(str)
}

// IsUUIDv4 checks if the string is a UUIDv4
func IsUUIDv4(str string) bool {
	return RegexUUID4.MatchString(str)

}

// IsLegacyID A legacyId does not contain the document prefix; the prefix is accounted for in the service implementation code
// will validate LONG ID
func IsLegacyID(str string) bool {
	return rxLegacyId.MatchString(str)
}

// IsS12ID checks if the string has the format of S12ID or UUID
func IsS12ID(str string) bool {
	return rxS12Id.MatchString(str)
}

// IsValidEmail checks if an email address is a valid RFC 5322 address
func IsValidEmail(str string, checkDomain bool) bool {
	// ignore checkDomain for now
	return RegexEmail.MatchString(str)
}

type Validator interface {
	Validate() error
}

type fieldError struct {
	fieldStack []string
	nestedErr  error
}

func (f *fieldError) Error() string {
	return "invalid field " + strings.Join(f.fieldStack, ".") + ": " + f.nestedErr.Error()
}

// FieldError wraps a given Validator error providing a message call stack.
func FieldError(fieldName string, err error) error {
	if ferr, ok := err.(*fieldError); ok {
		ferr.fieldStack = append([]string{fieldName}, ferr.fieldStack...)
		return ferr
	}
	return &fieldError{
		fieldStack: []string{fieldName},
		nestedErr:  err,
	}
}
