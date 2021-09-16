// Code generated by protoc-gen-govalidator. DO NOT EDIT.

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
		return fmt.Errorf(`id: value '%v' must be parsable as a UUID`, m.Id)
	}
	if len(m.UserId) != proto.UUIDSize {
		return fmt.Errorf(`user_id: value '%v' must be exactly 16 bytes long to be a valid UUID`, m.UserId)
	}
	if m.Email != "" {
		if !proto.IsValidEmail(m.Email) {
			return fmt.Errorf(`email: value '%v' must be parsable as a valid email address`, m.Email)
		}
	}
	if !(m.Age > 0) {
		return fmt.Errorf(`age: value '%v' must be greater than '0'`, m.Age)
	}
	if !(m.Speed < 110) {
		return fmt.Errorf(`speed: value '%v' must be less than '110'`, m.Speed)
	}
	if !(m.Score >= 0) {
		return fmt.Errorf(`score: value '%v' must be greater than or equal to '0'`, m.Score)
	}
	if !(m.Score <= 100) {
		return fmt.Errorf(`score: value '%v' must be less than or equal to '100'`, m.Score)
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
			return fmt.Errorf(`ids: value '%v' must be exactly 16 bytes long to be a valid UUID`, item)
		}
	}
	if m.MediaId != "" {
		if !proto.IsUUID(m.MediaId) {
			return fmt.Errorf(`media_id: value '%v' must be parsable as a UUID`, m.MediaId)
		}
	}
	if !(len(m.Description) <= 2000) {
		return fmt.Errorf(`description: value '%v' must have length less than or equal to '2000'`, m.Description)
	}
	if !(len(m.Password) >= 8) {
		return fmt.Errorf(`password: value '%v' must have length greater than or equal to '8'`, m.Password)
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
		return fmt.Errorf(`legacy_id: value '%v' must be parsable as a UUID or a legacy ID`, m.LegacyId)
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
		return fmt.Errorf(`name: value '%v' must have length greater than or equal to '6'`, _trim_ExampleMessage_Name)
	}
	if !(len(_trim_ExampleMessage_Name) <= 10) {
		return fmt.Errorf(`name: value '%v' must have length less than or equal to '10'`, _trim_ExampleMessage_Name)
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
		return fmt.Errorf(`url: value '%v' must be a string conforming to regex '%s'`, m.Url, _regex_val_ExampleMessage_Url)
	}
	return nil
}

func (m *ExampleMessage_NestedMessage) Validate() error {
	_trim_ExampleMessage_NestedMessage_Val := strings.TrimSpace(m.Val)
	_ = _trim_ExampleMessage_NestedMessage_Val
	if !(len(_trim_ExampleMessage_NestedMessage_Val) >= 1) {
		return fmt.Errorf(`val: value '%v' must have length greater than or equal to '1'`, _trim_ExampleMessage_NestedMessage_Val)
	}
	if !(len(_trim_ExampleMessage_NestedMessage_Val) <= 40) {
		return fmt.Errorf(`val: value '%v' must have length less than or equal to '40'`, _trim_ExampleMessage_NestedMessage_Val)
	}
	if !_regex_ExampleMessage_NestedMessage_NestedEmail.MatchString(m.NestedEmail) {
		return fmt.Errorf(`nested_email: value '%v' must be a string conforming to regex '%s'`, m.NestedEmail, _regex_val_ExampleMessage_NestedMessage_NestedEmail)
	}
	if !(len(m.MemberEmails) >= 2) {
		return fmt.Errorf(`member_emails: length '%v' must be greater than or equal to '2'`, len(m.MemberEmails))
	}
	if !(len(m.MemberEmails) <= 5) {
		return fmt.Errorf(`member_emails: length '%v' must be lesser than or equal to '5'`, len(m.MemberEmails))
	}
	for _, item := range m.MemberEmails {
		if !_regex_ExampleMessage_NestedMessage_MemberEmails.MatchString(item) {
			return fmt.Errorf(`member_emails: value '%v' must be a string conforming to regex '%s'`, item, _regex_val_ExampleMessage_NestedMessage_MemberEmails)
		}
	}
	return nil
}

func (m *ExampleMessage_NestedMessage_InnerNestedMessage) Validate() error {
	_trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal := strings.TrimSpace(m.InnerVal)
	_ = _trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal
	if !(len(_trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal) >= 1) {
		return fmt.Errorf(`inner_val: value '%v' must have length greater than or equal to '1'`, _trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal)
	}
	if !(len(_trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal) <= 40) {
		return fmt.Errorf(`inner_val: value '%v' must have length less than or equal to '40'`, _trim_ExampleMessage_NestedMessage_InnerNestedMessage_InnerVal)
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
		return fmt.Errorf(`id: value '%v' must be parsable as a UUID`, m.Id)
	}
	return nil
}

func (m *InnerMessageWithLegacyId) Validate() error {
	if !proto.IsLegacyID(m.Id) {
		return fmt.Errorf(`id: value '%v' must be parsable as a UUID or a legacy ID`, m.Id)
	}
	return nil
}
