// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: example.proto

package example

import fmt "fmt"
import regexp "regexp"
import github_com_SafetyCulture_s12_proto_protobuf_s12proto "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_ExampleMessage_Email = regexp.MustCompile(`.+\@.+\..+`)

func (m *ExampleMessage) Validate() error {
	if !github_com_SafetyCulture_s12_proto_protobuf_s12proto.IsUUID(m.Id) {
		return fmt.Errorf(`Id: value %q must be parsable as a UUID`, m.Id)
	}
	if len(m.UserID) != github_com_SafetyCulture_s12_proto_protobuf_s12proto.UUIDSize {
		return fmt.Errorf(`UserID: value %q must be exactly 16 bytes long to be a valid UUID`, m.UserID)
	}
	if !_regex_ExampleMessage_Email.MatchString(m.Email) {
		return fmt.Errorf(`Email: value %q must be a string conforming to regex ".+\\@.+\\..+"`, m.Email)
	}
	if !(m.Age > 0) {
		return fmt.Errorf(`Age: value %q must be greater than '0'`, m.Age)
	}
	if !(m.Speed < 110) {
		return fmt.Errorf(`Speed: value %q must be less than '110'`, m.Speed)
	}
	if !(m.Score >= 0) {
		return fmt.Errorf(`Score: value %q must be greater than or equal to '0'`, m.Score)
	}
	if !(m.Score <= 100) {
		return fmt.Errorf(`Score: value %q must be less than or equal to '100'`, m.Score)
	}
	if m.Inner != nil {
		if v, ok := interface{}(m.Inner).(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return github_com_SafetyCulture_s12_proto_protobuf_s12proto.FieldError("Inner", err)
			}
		}
	}
	for _, item := range m.Ids {
		if len(item) != github_com_SafetyCulture_s12_proto_protobuf_s12proto.UUIDSize {
			return fmt.Errorf(`Ids: value %q must be exactly 16 bytes long to be a valid UUID`, item)
		}
	}
	if m.MediaId != "" {
		if !github_com_SafetyCulture_s12_proto_protobuf_s12proto.IsUUID(m.MediaId) {
			return fmt.Errorf(`MediaId: value %q must be parsable as a UUID`, m.MediaId)
		}
	}
	if !(len(m.Description) <= 2000) {
		return fmt.Errorf(`Description: value %q must have length less than or equal to '2000'`, m.Description)
	}
	if !(len(m.Password) >= 8) {
		return fmt.Errorf(`Password: value %q must have length greater than or equal to '8'`, m.Password)
	}
	return nil
}

func (m *InnerMessage) Validate() error {
	if !github_com_SafetyCulture_s12_proto_protobuf_s12proto.IsUUID(m.Id) {
		return fmt.Errorf(`Id: value %q must be parsable as a UUID`, m.Id)
	}
	return nil
}
