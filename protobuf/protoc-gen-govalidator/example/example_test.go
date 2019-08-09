package example

import (
	"strings"
	"testing"

	"github.com/SafetyCulture/s12-proto/protobuf/s12proto"
)

func TestValidationRules(t *testing.T) {

	const (
		id       string = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
		email    string = "email@address.co"
		password string = "12345678"
	)
	var (
		byteID []byte = []byte{53, 30, 208, 165, 196, 219, 75, 61, 142, 60, 101, 84, 229, 43, 61, 108}
	)

	sb := strings.Builder{}
	for i := 0; i < 2001; i++ {
		sb.WriteRune('a')
	}

	tests := []struct {
		name        string
		input       s12proto.Validator
		shouldError bool
	}{
		{
			"ValidMessage",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			false,
		}, {
			"UUIDWithoutDashes",
			&ExampleMessage{
				Id:          strings.Replace(id, "-", "", -1),
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			false,
		}, {
			"EmptyUUID",
			&ExampleMessage{
				Id:          "",
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:          "notauuid",
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:          id + "a",
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:          "92b6c2f9-abd8-48bc-a2c9-bf70e969751h",
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"EmptyByteUUID",
			&ExampleMessage{
				Id:          id,
				UserID:      []byte{},
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidByteUUID",
			&ExampleMessage{
				Id:          id,
				UserID:      []byte{53},
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidByteUUID",
			&ExampleMessage{
				Id:          id,
				UserID:      append(byteID, 67),
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"EmptyRegex",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       "",
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidRegex",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       "something@else",
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidGreaterThan",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         0,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"ValidLessThan",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Speed:       10,
			},
			false,
		}, {
			"ValidLessThan",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Speed:       -10,
			},
			false,
		}, {
			"InvalidLessThan",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Speed:       120,
			},
			true,
		}, {
			"ValidGreaterThanEquals",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       1,
			},
			false,
		}, {
			"InvalidGreaterThanEquals",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       -1,
			},
			true,
		}, {
			"ValidLessThanEquals",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       100,
			},
			false,
		}, {
			"ValidLessThanEquals",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       99,
			},
			false,
		}, {
			"InvalidLessThanEquals",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       111,
			},
			true,
		}, {
			"ValidRepeated",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{byteID},
			},
			false,
		}, {
			"ValidRepeated",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{byteID, byteID},
			},
			false,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{[]byte{}},
			},
			true,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{[]byte{121}},
			},
			true,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{byteID, []byte{121}},
			},
			true,
		}, {
			"ValidOptional",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				MediaId:     "",
			},
			false,
		}, {
			"InalidOptional",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				MediaId:     "notauuid",
			},
			true,
		}, {
			"ValidMaxLength",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Description: "",
			},
			false,
		}, {
			"ValidMaxLength",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Description: "Some text here",
			},
			false,
		}, {
			"InvalidMaxLength",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Description: sb.String(),
			},
			true,
		}, {
			"InvalidMinLength",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    "",
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"InvalidMinLength",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    "1234567",
				MsgRequired: &InnerMessage{Id: id},
			},
			true,
		}, {
			"ValidInnerMsg",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Inner: &InnerMessage{
					Id: id,
				},
			},
			false,
		}, {
			"InvalidInnerMsg",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Inner:       &InnerMessage{},
			},
			true,
		}, {
			"InvalidInnerMsg",
			&ExampleMessage{
				Id:          id,
				UserID:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Inner: &InnerMessage{
					Id: "notauuid",
				},
			},
			true,
		}, {
			"RequiredMessage",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.input.Validate()
			if test.shouldError && err == nil {
				t.Error("expected error, but nil")
			}
			if !test.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

}
