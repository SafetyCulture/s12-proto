package example

import (
	"strings"
	"testing"

	"github.com/SafetyCulture/s12-proto/protobuf/s12proto"
)

func TestValidationRules(t *testing.T) {

	const (
		id            string = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
		legacyId      string = "56341C6E-35A7-4C97-9C5E-7AC79673EAB2"
		legacyLongIdFail string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-9-00000190327E0675" // length = 49 (without dashes)
		legacyLongId1 string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-90-00000190327E0675" // length = 50 (without dashes)
		legacyLongId2 string = "005F2E38-8426-48AF-94DE-5FEA3A396EEA-891-00000153F68896DC" // length = 51 (without dashes)
		legacyLongId3 string = "007B516E-53F1-4AA0-ABAF-8C78342A2C82-2388-00000221F1C2BD1E" // length = 52 (without dashes)
		legacyLongId4 string = "00709A17-151F-4CFC-B412-F080343ED84D-11977-000010227B4C60A9" // length = 53 (without dashes)
		email         string = "email@address.co"
		password      string = "12345678"
	)
	var (
		byteID []byte = []byte{53, 30, 208, 165, 196, 219, 75, 61, 142, 60, 101, 84, 229, 43, 61, 108}
	)

	sb := strings.Builder{}
	for i := 0; i < 2001; i++ {
		sb.WriteRune('a')
	}

	tests := [...]struct {
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				Inner:       &InnerMessage{Id: id},
				LegacyID:    legacyId,
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
				LegacyID:    legacyId,
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
				Inner:       &InnerMessage{Id: "notauuid"},
				LegacyID:    legacyId,
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
				LegacyID: legacyId,
			},
			true,
		}, {
			"LegacyID",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyId,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			false,
		}, {
			"LegacyIDWithoutDashes",
			&ExampleMessage{
				Id:          id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyID: strings.Replace(legacyId, "-", "", -1),
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			false,
		}, {
			"EmptyLegacyID",
			&ExampleMessage{
				Id:          id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyID: "",
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"InvalidLegacyID",
			&ExampleMessage{
				Id:          id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyID: "totally-invalid",
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"InvalidLegacyID",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyId + "1",
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"LegacyIDWithInvalidInnerLegacyID",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyId,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId + "a"},
			},
			true,
		}, {
			"LegacyLongIDWithLength50",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongIdFail,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			true,
		}, {
			"LegacyLongIDWithLength50",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId1,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength51",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId2,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength52",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId3,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength53",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId4,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength54",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId4 + "1",
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			true,
		},
		{
			"LegacyLongIDWithoutDashes",
			&ExampleMessage{
				Id:          id,
				UserID:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyID: strings.Replace(legacyLongId1, "-", "", -1),
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		}, {
			"InvalidLegacyLongID",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId1 + "1",
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			true,
		}, {
			"LegacyLongIDWithInvalidInnerLegacyLongID",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId1,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1 + "a"},
			},
			true,
		}, {
			"MixOfLegacyIDAndLegacyLongID-1",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyId,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		}, {
			"MixOfLegacyIDAndLegacyLongID-2",
			&ExampleMessage{
				Id:            id,
				UserID:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyID:      legacyLongId1,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			false,
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
