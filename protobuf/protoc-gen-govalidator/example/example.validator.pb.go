// Code generated by protoc-gen-govalidator. DO NOT EDIT.
// versions:
// 	protoc-gen-govalidator v2.1.1
// 	protoc                 v3.17.3
// source: example.proto

package example

import (
	fmt "fmt"
	proto "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
	regexp "regexp"
	strings "strings"
)

const _regex_val_ExampleMessage_Url = `https:\/\/www\.safetyculture\.(io|com)`

var _regex_ExampleMessage_Url = regexp.MustCompile(_regex_val_ExampleMessage_Url)

const _regex_val_ExampleMessage_NestedMessage_NestedEmail = `.+\@.+\..+`

var _regex_ExampleMessage_NestedMessage_NestedEmail = regexp.MustCompile(_regex_val_ExampleMessage_NestedMessage_NestedEmail)

const _regex_val_ExampleMessage_NestedMessage_MemberEmails = `[a-z0-9!#$&'*+/=?^_{|}~-]+(?:\.[a-z0-9!#$&'*+/=?^_{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?`

var _regex_ExampleMessage_NestedMessage_MemberEmails = regexp.MustCompile(_regex_val_ExampleMessage_NestedMessage_MemberEmails)

func (m *ExampleMessage) Validate() error {
	if !proto.IsUUID(m.Id) {
		return fmt.Errorf(`id: value must be parsable as a UUID`)
	}
	if len(m.UserId) != proto.UUIDSize {
		return fmt.Errorf(`user_id: value must be exactly 16 bytes long to be a valid UUID`)
	}
	if m.Email != "" {
		if !proto.IsValidEmail(m.Email, false) {
			return fmt.Errorf(`email: value must be parsable as an email address`)
		}
	}
	if !(m.Age > 0) {
		return fmt.Errorf(`age: value must be greater than 0`)
	}
	if !(m.Speed < 110) {
		return fmt.Errorf(`speed: value must be less than 110`)
	}
	if !(m.Score >= 0) {
		return fmt.Errorf(`score: value must be greater than or equal to 0`)
	}
	if !(m.Score <= 100) {
		return fmt.Errorf(`score: value must be less than or equal to 100`)
	}
	if m.Inner != nil {
		if v, ok := interface{}(m.Inner).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("inner", err)
			}
		}
	}
	for _, item := range m.Ids {
		if len(item) != proto.UUIDSize {
			return fmt.Errorf(`ids: value must be exactly 16 bytes long to be a valid UUID`)
		}
	}
	if m.MediaId != "" {
		if !proto.IsUUID(m.MediaId) {
			return fmt.Errorf(`media_id: value must be parsable as a UUID`)
		}
	}
	if !(len(m.Description) <= 2000) {
		return fmt.Errorf(`description: value must have length less than or equal to 2000`)
	}
	if !(len(m.Password) >= 8) {
		return fmt.Errorf(`password: value must have length greater than or equal to 8`)
	}
	// Validation of oneof fields is unsupported.
	if m.MsgRequired == nil {
		return fmt.Errorf("field msg_required is required")
	}
	if m.MsgRequired != nil {
		if v, ok := interface{}(m.MsgRequired).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("msg_required", err)
			}
		}
	}
	if !proto.IsLegacyID(m.LegacyId) {
		return fmt.Errorf(`legacy_id: value must be parsable as a UUID or a legacy ID`)
	}
	if m.InnerLegacyId != nil {
		if v, ok := interface{}(m.InnerLegacyId).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("inner_legacy_id", err)
			}
		}
	}
	_trim_ExampleMessage_Name := strings.TrimSpace(m.Name)
	_ = _trim_ExampleMessage_Name
	if !(len(_trim_ExampleMessage_Name) >= 6) {
		return fmt.Errorf(`name: value must have length greater than or equal to 6`)
	}
	if !(len(_trim_ExampleMessage_Name) <= 10) {
		return fmt.Errorf(`name: value must have length less than or equal to 10`)
	}
	if m.NestedMessage != nil {
		if v, ok := interface{}(m.NestedMessage).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("nested_message", err)
			}
		}
	}
	// Validation of proto3 map<> fields is unsupported.
	if !_regex_ExampleMessage_Url.MatchString(m.Url) {
		return fmt.Errorf(`url: value must be a string conforming to predefined pattern`)
	}
	if !(len(m.ContactsWithLengthConstraint) >= 0) {
		return fmt.Errorf(`contacts_with_length_constraint: length must be greater than or equal to 0`)
	}
	if !(len(m.ContactsWithLengthConstraint) <= 10) {
		return fmt.Errorf(`contacts_with_length_constraint: length must be lesser than or equal to 10`)
	}
	if len(m.ContactsWithLengthConstraint) > 0 {
		for _, item := range m.ContactsWithLengthConstraint {
			if item != nil {
				if v, ok := interface{}(item).(proto.Validator); ok {
					if err := v.Validate(); err != nil {
						return proto.FieldError("contacts_with_length_constraint", err)
					}
				}
			}
		}
	}
	if len(m.ContactsWithNoLengthConstraint) > 0 {
		for _, item := range m.ContactsWithNoLengthConstraint {
			if item != nil {
				if v, ok := interface{}(item).(proto.Validator); ok {
					if err := v.Validate(); err != nil {
						return proto.FieldError("contacts_with_no_length_constraint", err)
					}
				}
			}
		}
	}
	if m.ScheduledFor != nil {
		if v, ok := interface{}(m.ScheduledFor).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("scheduled_for", err)
			}
		}
	}
	return nil
}

func (m *ExampleMessage_NestedMessage) Validate() error {
	_trim_ExampleMessage_NestedMessage_Val := strings.TrimSpace(m.Val)
	_ = _trim_ExampleMessage_NestedMessage_Val
	if !(len(_trim_ExampleMessage_NestedMessage_Val) >= 1) {
		return fmt.Errorf(`val: value must have length greater than or equal to 1`)
	}
	if !(len(_trim_ExampleMessage_NestedMessage_Val) <= 40) {
		return fmt.Errorf(`val: value must have length less than or equal to 40`)
	}
	if !_regex_ExampleMessage_NestedMessage_NestedEmail.MatchString(m.NestedEmail) {
		return fmt.Errorf(`nested_email: value must be a string conforming to predefined pattern`)
	}
	if !(len(m.MemberEmails) >= 2) {
		return fmt.Errorf(`member_emails: length must be greater than or equal to 2`)
	}
	if !(len(m.MemberEmails) <= 5) {
		return fmt.Errorf(`member_emails: length must be lesser than or equal to 5`)
	}
	for _, item := range m.MemberEmails {
		if !_regex_ExampleMessage_NestedMessage_MemberEmails.MatchString(item) {
			return fmt.Errorf(`member_emails: value must be a string conforming to predefined pattern`)
		}
	}
	return nil
}

func (m *ExampleMessage_NestedMessage_InnerNestedMessage) Validate() error {
	_trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal := strings.TrimSpace(m.InnerVal)
	_ = _trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal
	if !(len(_trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal) >= 1) {
		return fmt.Errorf(`inner_val: value must have length greater than or equal to 1`)
	}
	if !(len(_trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal) <= 40) {
		return fmt.Errorf(`inner_val: value must have length less than or equal to 40`)
	}
	return nil
}

func (m *ExampleMessage_Contact) Validate() error {
	if m.Phone != "" {
	}
	if m.Email != "" {
		if !proto.IsValidEmail(m.Email, false) {
			return fmt.Errorf(`email: value must be parsable as an email address`)
		}
	}
	return nil
}

func (m *OuterMessageUsingNestedMessage) Validate() error {
	if m.SomeMessage != nil {
		if v, ok := interface{}(m.SomeMessage).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("some_message", err)
			}
		}
	}
	return nil
}

func (m *InnerMessage) Validate() error {
	if !proto.IsUUID(m.Id) {
		return fmt.Errorf(`id: value must be parsable as a UUID`)
	}
	return nil
}

func (m *InnerMessageWithLegacyId) Validate() error {
	if !proto.IsLegacyID(m.Id) {
		return fmt.Errorf(`id: value must be parsable as a UUID or a legacy ID`)
	}
	return nil
}
