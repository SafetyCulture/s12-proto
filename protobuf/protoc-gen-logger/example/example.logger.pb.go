// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: example.proto

package example

import (
	"reflect"

	github_com_SafetyCulture_s12_proto_protobuf_s12proto "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	proto "github.com/gogo/protobuf/proto"

	fmt "fmt"

	math "math"

	_ "github.com/SafetyCulture/s12-proto/protobuf/s12proto"

	_ "github.com/gogo/protobuf/gogoproto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *ExampleMessage) Parse(isLevelEnabled func(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level) bool) proto.Message {
	res := &ExampleMessage{}
	if isLevelEnabled(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level_DEBUG) {
		res.Id = this.Id
	}
	res.UserName = this.UserName
	if isLevelEnabled(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level_ERROR) {
		res.Password = this.Password
	}
	if isLevelEnabled(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level_INFO) {
		if this.SomeKindOfInnerValue != nil {
			res.SomeKindOfInnerValue = this.SomeKindOfInnerValue.Parse(isLevelEnabled).(*InnerMessage)
		}
	}
	if reflect.TypeOf(this.TestOneof) == reflect.TypeOf(&ExampleMessage_OneOf1{}) {
		res.TestOneof = this.TestOneof
	}
	if reflect.TypeOf(this.TestOneof) == reflect.TypeOf(&ExampleMessage_OneOf2{}) {
		if isLevelEnabled(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level_ERROR) {
			oneof_ExampleMessage_OneOf2 := this.TestOneof.(*ExampleMessage_OneOf2)
			if oneof_ExampleMessage_OneOf2 != nil {
				oneof_ExampleMessage_OneOf2.OneOf2 = oneof_ExampleMessage_OneOf2.OneOf2.Parse(isLevelEnabled).(*OneOfMessage)
				res.TestOneof = oneof_ExampleMessage_OneOf2
			}
		}
	}
	return res
}
func (this *OneOfMessage) Parse(isLevelEnabled func(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level) bool) proto.Message {
	res := &OneOfMessage{}
	res.Value = this.Value
	return res
}
func (this *InnerMessage) Parse(isLevelEnabled func(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level) bool) proto.Message {
	res := &InnerMessage{}
	if isLevelEnabled(github_com_SafetyCulture_s12_proto_protobuf_s12proto.Level_DEBUG) {
		res.Body = this.Body
	}
	return res
}