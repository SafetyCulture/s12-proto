// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package proto

import (
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

const (
	// UUIDSize is the size of a UUID in bytes.
	UUIDSize int = 16
)

const (
	uuid     string = "^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$"
	legacyId string = "(?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}(-?[0-9a-f]{2,5}-?[0-9a-f]{16})?$"
	rePua    string = `[\x{E000}-\x{F8FF}]` // Private Use Codepoints in the Basic Multilingual Plane (not including planes 15, 16)
)

var (
	rxUUID                      = regexp.MustCompile(uuid)
	rxLegacyId                  = regexp.MustCompile(legacyId)
	UnsafeCharReplacer          = strings.NewReplacer("\u0021", "\uFF01", "\u0022", "\u201D", "\u0023", "\u0020", "\u0025", "\u2052", "\u0026", "\u0020", "\u0027", "\u2019", "\u002A", "\u2217", "\u002B", "\uFF0B", "\u002D", "\u2212", "\u002F", "\u2215", "\u003B", "\u037E", "\u003C", "\u02C2", "\u003D", "\u2E40", "\u003E", "\u02C3", "\u005C", "\uFF3C", "\u0060", "\u2019", "\u007C", "\uFFE8", "\u3164", "\u0020")
	SymbolCharReplacer          = strings.NewReplacer("\u00A0", "\u0020", "\u1680", "\u0020", "\u2000", "\u0020", "\u2001", "\u0020", "\u2002", "\u0020", "\u2003", "\u0020", "\u2004", "\u0020", "\u2005", "\u0020", "\u2006", "\u0020", "\u2007", "\u0020", "\u2008", "\u0020", "\u2009", "\u0020", "\u200A", "\u0020", "\u202F", "\u0020", "\u205F", "\u0020", "\u3000", "\u0020", "\u200C", "", "\u200D", "", "\uFEFF", "", "\u2014", "\u2013", "\u2018", "\u2019", "\u3002", "\u002E", "\uFF0C", "\u002C", "\uFF1A", "\u003A", "\u0009", "\u0020", "\u000A", "\u0020", "\u000D", "\u0020")
	SymbolCharReplacerMultiline = strings.NewReplacer("\u00A0", "\u0020", "\u1680", "\u0020", "\u2000", "\u0020", "\u2001", "\u0020", "\u2002", "\u0020", "\u2003", "\u0020", "\u2004", "\u0020", "\u2005", "\u0020", "\u2006", "\u0020", "\u2007", "\u0020", "\u2008", "\u0020", "\u2009", "\u0020", "\u200A", "\u0020", "\u202F", "\u0020", "\u205F", "\u0020", "\u3000", "\u0020", "\u200C", "", "\u200D", "", "\uFEFF", "", "\u2014", "\u2013", "\u2018", "\u2019", "\u3002", "\u002E", "\uFF0C", "\u002C", "\uFF1A", "\u003A", "\u0009", "\u0020", "\u000D", "\u0020")
	RegexPua                    = regexp.MustCompile(rePua)
)

// IsUUID checks if the string is a UUID (version 3, 4 or 5).
func IsUUID(str string) bool {
	return rxUUID.MatchString(str)
}

// IsUUID checks if the string is a UUIDv4
func IsUUIDv4(str string) bool {
	return govalidator.IsUUIDv4(str)
}

// A legacyId does not contain the document prefix; the prefix is accounted for in the service implementation code
func IsLegacyID(str string) bool {
	return rxLegacyId.MatchString(str)
}

// IsValidEmail checks if an email address is a valid RFC 5322 address
func IsValidEmail(str string, checkDomain bool) bool {
	valid := govalidator.IsEmail(str)
	if valid && checkDomain {
		return govalidator.IsExistingEmail(str)
	}
	return valid
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
