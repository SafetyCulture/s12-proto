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
	legacyId string = "(?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}(-?[0-9a-f]{2,5}-?[0-9a-f]{16})?$"
)

var (
	rxUUID = regexp.MustCompile(uuid)
	rxLegacyId = regexp.MustCompile(legacyId)
)

// IsUUID checks if the string is a UUID (version 3, 4 or 5).
func IsUUID(str string) bool {
	return rxUUID.MatchString(str)
}

// A legacyId does not contain the document prefix; the prefix is accounted for in the service implementation code
func IsLegacyID(str string) bool {
	return rxLegacyId.MatchString(str)
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
