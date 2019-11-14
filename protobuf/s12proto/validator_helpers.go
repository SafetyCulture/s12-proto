// Copyright (c) 2019 SafetyCulture Pty Ltd. All Rights Reserved.

package s12proto

import (
	"regexp"
	"strings"
)

const (
	// UUIDSize is the size of a UUID in bytes.
	UUIDSize int = 16
)

const (
	uuid string = "^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$"
	upper_case_uuid string = "^[[:xdigit:]]{8}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{12}$"
	globally_unique_ios_string string = "^[[:xdigit:]]{8}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{12}-?[[:xdigit:]]{2,5}-?[[:xdigit:]]{16}$"
)

var (
	rxUUID = regexp.MustCompile(uuid)
	rxUUIDUpperCase = regexp.MustCompile(upper_case_uuid)
	rxUniqueIOsString = regexp.MustCompile(globally_unique_ios_string)
)

// IsUUID checks if the string is a UUID (version 3, 4 or 5).
func IsUUID(str string) bool {
	return rxUUID.MatchString(str)
}

func IsLegacyID(str string) bool {
	return rxUUIDUpperCase.MatchString(str) || rxUniqueIOsString.MatchString(str)
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
