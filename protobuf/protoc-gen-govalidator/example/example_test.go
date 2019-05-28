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
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
			},
			false,
		}, {
			"EmptyUUID",
			&ExampleMessage{
				Id:       "",
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:       "notauuid",
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:       id + "a",
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:       "92b6c2f9-abd8-48bc-a2c9-bf70e969751h",
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"EmptyByteUUID",
			&ExampleMessage{
				Id:       id,
				UserID:   []byte{},
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"InvalidByteUUID",
			&ExampleMessage{
				Id:       id,
				UserID:   []byte{53},
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"InvalidByteUUID",
			&ExampleMessage{
				Id:       id,
				UserID:   append(byteID, 67),
				Email:    email,
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"EmptyRegex",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    "",
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"InvalidRegex",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    "something@else",
				Age:      18,
				Password: password,
			},
			true,
		}, {
			"InvalidGreaterThan",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      0,
				Password: password,
			},
			true,
		}, {
			"ValidLessThan",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Speed:    10,
			},
			false,
		}, {
			"ValidLessThan",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Speed:    -10,
			},
			false,
		}, {
			"InvalidLessThan",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Speed:    120,
			},
			true,
		}, {
			"ValidGreaterThanEquals",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Score:    1,
			},
			false,
		}, {
			"InvalidGreaterThanEquals",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Score:    -1,
			},
			true,
		}, {
			"ValidLessThanEquals",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Score:    100,
			},
			false,
		}, {
			"ValidLessThanEquals",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Score:    99,
			},
			false,
		}, {
			"InvalidLessThanEquals",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Score:    111,
			},
			true,
		}, {
			"ValidRepeated",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Ids:      [][]byte{byteID},
			},
			false,
		}, {
			"ValidRepeated",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Ids:      [][]byte{byteID, byteID},
			},
			false,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Ids:      [][]byte{[]byte{}},
			},
			true,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Ids:      [][]byte{[]byte{121}},
			},
			true,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Ids:      [][]byte{byteID, []byte{121}},
			},
			true,
		}, {
			"ValidOptional",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				MediaId:  "",
			},
			false,
		}, {
			"InalidOptional",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				MediaId:  "notauuid",
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
				Description: sb.String(),
			},
			true,
		}, {
			"InvalidMinLength",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: "",
			},
			true,
		}, {
			"InvalidMinLength",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: "1234567",
			},
			true,
		}, {
			"ValidInnerMsg",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Inner: &InnerMessage{
					Id: id,
				},
			},
			false,
		}, {
			"InvalidInnerMsg",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Inner:    &InnerMessage{},
			},
			true,
		}, {
			"InvalidInnerMsg",
			&ExampleMessage{
				Id:       id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				Inner: &InnerMessage{
					Id: "notauuid",
				},
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
