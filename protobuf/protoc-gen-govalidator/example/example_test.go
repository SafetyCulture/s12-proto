package example

import (
	"strings"
	"testing"

	"github.com/SafetyCulture/s12-proto/protobuf/s12proto"
)

func TestValidationRules(t *testing.T) {

	const (
		id               string = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a"
		legacyId         string = "56341C6E-35A7-4C97-9C5E-7AC79673EAB2"
		legacyLongIdFail string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-9-00000190327E0675"     // length = 49 (without dashes)
		legacyLongId1    string = "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-90-00000190327E0675"    // length = 50 (without dashes)
		legacyLongId2    string = "005F2E38-8426-48AF-94DE-5FEA3A396EEA-891-00000153F68896DC"   // length = 51 (without dashes)
		legacyLongId3    string = "007B516E-53F1-4AA0-ABAF-8C78342A2C82-2388-00000221F1C2BD1E"  // length = 52 (without dashes)
		legacyLongId4    string = "00709A17-151F-4CFC-B412-F080343ED84D-11977-000010227B4C60A9" // length = 53 (without dashes)
		email            string = "email@address.co"
		password         string = "12345678"
		name             string = "safety"
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
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"UUIDWithoutDashes",
			&ExampleMessage{
				Id:          strings.Replace(id, "-", "", -1),
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"EmptyUUID",
			&ExampleMessage{
				Id:          "",
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:          "notauuid",
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:          id + "a",
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidUUID",
			&ExampleMessage{
				Id:          "92b6c2f9-abd8-48bc-a2c9-bf70e969751h",
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"EmptyByteUUID",
			&ExampleMessage{
				Id:          id,
				UserId:      []byte{},
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidByteUUID",
			&ExampleMessage{
				Id:          id,
				UserId:      []byte{53},
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidByteUUID",
			&ExampleMessage{
				Id:          id,
				UserId:      append(byteID, 67),
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"EmptyRegex",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       "",
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidRegex",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       "something@else",
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidGreaterThan",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         0,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"ValidLessThan",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Speed:       10,
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"ValidLessThan",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Speed:       -10,
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"InvalidLessThan",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Speed:       120,
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"ValidGreaterThanEquals",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       1,
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"InvalidGreaterThanEquals",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       -1,
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"ValidLessThanEquals",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       100,
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"ValidLessThanEquals",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       99,
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"InvalidLessThanEquals",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Score:       111,
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"ValidRepeated",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{byteID},
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"ValidRepeated",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{byteID, byteID},
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{[]byte{}},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{[]byte{121}},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidRepeated",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Ids:         [][]byte{byteID, []byte{121}},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"ValidOptional",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				MediaId:     "",
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"InalidOptional",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				MediaId:     "notauuid",
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"ValidMaxLength",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Description: "",
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"ValidMaxLength",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Description: "Some text here",
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"InvalidMaxLength",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Description: sb.String(),
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidMinLength",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    "",
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidMinLength",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    "1234567",
				MsgRequired: &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"ValidInnerMsg",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Inner:       &InnerMessage{Id: id},
				LegacyId:    legacyId,
				Name:        name,
			},
			false,
		}, {
			"InvalidInnerMsg",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Inner:       &InnerMessage{},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"InvalidInnerMsg",
			&ExampleMessage{
				Id:          id,
				UserId:      byteID,
				Email:       email,
				Age:         18,
				Password:    password,
				MsgRequired: &InnerMessage{Id: id},
				Inner:       &InnerMessage{Id: "notauuid"},
				LegacyId:    legacyId,
				Name:        name,
			},
			true,
		}, {
			"RequiredMessage",
			&ExampleMessage{
				Id:       id,
				UserId:   byteID,
				Email:    email,
				Age:      18,
				Password: password,
				LegacyId: legacyId,
				Name:     name,
			},
			true,
		}, {
			"LegacyId",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyId,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			false,
		}, {
			"UUIDAsLegacyID",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      id,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: id},
			},
			false,
		}, {
			"LegacyIDWithoutDashes",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      strings.Replace(legacyId, "-", "", -1),
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			false,
		}, {
			"EmptyLegacyID",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      "",
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"InvalidLegacyID",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      "totally-invalid",
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"InvalidLegacyID",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyId + "1",
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"LegacyIDWithInvalidInnerLegacyID",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyId,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId + "a"},
			},
			true,
		}, {
			"LegacyLongIDWithLength50",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongIdFail,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			true,
		}, {
			"LegacyLongIDWithLength50",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId1,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength51",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId2,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength52",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId3,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength53",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId4,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		},
		{
			"LegacyLongIDWithLength54",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId4 + "1",
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			true,
		},
		{
			"LegacyLongIDWithoutDashes",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      strings.Replace(legacyLongId1, "-", "", -1),
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		}, {
			"InvalidLegacyLongID",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId1 + "1",
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			true,
		}, {
			"LegacyLongIDWithInvalidInnerLegacyLongID",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId1,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1 + "a"},
			},
			true,
		}, {
			"MixOfLegacyIDAndLegacyLongID-1",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyId,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyLongId1},
			},
			false,
		}, {
			"MixOfLegacyIDAndLegacyLongID-2",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId1,
				Name:          name,
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			false,
		}, {
			"OnlyWhitespaceInputOnName",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId1,
				Name:          "  \n\t\r  ",
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"LessThanMininumEffectiveLengthOnName",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId1,
				Name:          " \t\t12345  ", // min is 6
				InnerLegacyId: &InnerMessageWithLegacyId{Id: legacyId},
			},
			true,
		}, {
			"EffectiveLengthInRangeOnName",
			&ExampleMessage{
				Id:            id,
				UserId:        byteID,
				Email:         email,
				Age:           18,
				Password:      password,
				MsgRequired:   &InnerMessage{Id: id},
				LegacyId:      legacyLongId1,
				Name:          " \t\t1234567890\t\t\t\r  ", // min is 6, max 10
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
