// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: example.proto

package example

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/SafetyCulture/s12-proto/protobuf/s12proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ExampleMessage struct {
	// Returns an error if the string cannot be parsed as a UUID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Bytes can also be parsed as UUID
	UserId []byte `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// Strings can validate against a regular expresion
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// integers can be greater than a value
	Age int32 `protobuf:"varint,4,opt,name=age,proto3" json:"age,omitempty"`
	// intergers can be less than a value
	Speed int64 `protobuf:"varint,5,opt,name=speed,proto3" json:"speed,omitempty"`
	// intergers greater/less than or equal, the can also be combined
	Score                int32         `protobuf:"varint,6,opt,name=score,proto3" json:"score,omitempty"`
	Inner                *InnerMessage `protobuf:"bytes,7,opt,name=inner" json:"inner,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ExampleMessage) Reset()         { *m = ExampleMessage{} }
func (m *ExampleMessage) String() string { return proto.CompactTextString(m) }
func (*ExampleMessage) ProtoMessage()    {}
func (*ExampleMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_example_fcdf8b3cfcb7b6c2, []int{0}
}
func (m *ExampleMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExampleMessage.Unmarshal(m, b)
}
func (m *ExampleMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExampleMessage.Marshal(b, m, deterministic)
}
func (dst *ExampleMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExampleMessage.Merge(dst, src)
}
func (m *ExampleMessage) XXX_Size() int {
	return xxx_messageInfo_ExampleMessage.Size(m)
}
func (m *ExampleMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ExampleMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ExampleMessage proto.InternalMessageInfo

func (m *ExampleMessage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ExampleMessage) GetUserId() []byte {
	if m != nil {
		return m.UserId
	}
	return nil
}

func (m *ExampleMessage) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ExampleMessage) GetAge() int32 {
	if m != nil {
		return m.Age
	}
	return 0
}

func (m *ExampleMessage) GetSpeed() int64 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *ExampleMessage) GetScore() int32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *ExampleMessage) GetInner() *InnerMessage {
	if m != nil {
		return m.Inner
	}
	return nil
}

type InnerMessage struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InnerMessage) Reset()         { *m = InnerMessage{} }
func (m *InnerMessage) String() string { return proto.CompactTextString(m) }
func (*InnerMessage) ProtoMessage()    {}
func (*InnerMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_example_fcdf8b3cfcb7b6c2, []int{1}
}
func (m *InnerMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InnerMessage.Unmarshal(m, b)
}
func (m *InnerMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InnerMessage.Marshal(b, m, deterministic)
}
func (dst *InnerMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InnerMessage.Merge(dst, src)
}
func (m *InnerMessage) XXX_Size() int {
	return xxx_messageInfo_InnerMessage.Size(m)
}
func (m *InnerMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_InnerMessage.DiscardUnknown(m)
}

var xxx_messageInfo_InnerMessage proto.InternalMessageInfo

func (m *InnerMessage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*ExampleMessage)(nil), "example.ExampleMessage")
	proto.RegisterType((*InnerMessage)(nil), "example.InnerMessage")
}

func init() { proto.RegisterFile("example.proto", fileDescriptor_example_fcdf8b3cfcb7b6c2) }

var fileDescriptor_example_fcdf8b3cfcb7b6c2 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x4d, 0x4b, 0xc3, 0x30,
	0x1c, 0xc6, 0x97, 0x6e, 0xed, 0x66, 0x36, 0x77, 0x08, 0x28, 0x41, 0xc4, 0x96, 0x21, 0x12, 0xd0,
	0x6d, 0x3a, 0xf5, 0x0b, 0xf8, 0x72, 0xd8, 0xc1, 0x4b, 0x3c, 0x8a, 0x4a, 0xb6, 0xfc, 0x37, 0x03,
	0x5b, 0x53, 0x92, 0x54, 0x7c, 0xc1, 0x9b, 0xdf, 0x73, 0xd0, 0x4f, 0x22, 0x6d, 0xa6, 0xec, 0xe2,
	0x25, 0x90, 0xdf, 0xf3, 0xe4, 0x07, 0x79, 0xf0, 0x36, 0xbc, 0x89, 0x65, 0xb6, 0x80, 0x41, 0x66,
	0xb4, 0xd3, 0xa4, 0xb9, 0xbe, 0xee, 0xdd, 0xcc, 0x95, 0x7b, 0xc9, 0x27, 0x83, 0xa9, 0x5e, 0x0e,
	0xef, 0xc5, 0x0c, 0xdc, 0xfb, 0x75, 0xbe, 0x70, 0xb9, 0x81, 0xa1, 0x3d, 0x1b, 0xf5, 0xab, 0xf2,
	0xb0, 0x3a, 0x27, 0xf9, 0xac, 0x44, 0x9e, 0xbc, 0x8a, 0x85, 0x92, 0xc2, 0x69, 0xe3, 0x75, 0xbd,
	0xef, 0x00, 0x77, 0x6f, 0xbd, 0xf1, 0x0e, 0xac, 0x15, 0x73, 0x20, 0xbb, 0x38, 0x50, 0x92, 0xa2,
	0x04, 0xb1, 0xad, 0xab, 0xa8, 0x58, 0xc5, 0x41, 0x0b, 0xf1, 0x40, 0x49, 0x12, 0xe3, 0x66, 0x6e,
	0xc1, 0x3c, 0x2b, 0x49, 0x83, 0x04, 0xb1, 0xce, 0x5f, 0x18, 0x95, 0x78, 0x2c, 0xc9, 0x05, 0x6e,
	0x4b, 0xb0, 0x53, 0xa3, 0x32, 0xa7, 0x74, 0x4a, 0xeb, 0x95, 0x81, 0x14, 0xab, 0xb8, 0x4b, 0x3a,
	0x4f, 0x0f, 0xa2, 0xff, 0xf1, 0xf8, 0x39, 0x3a, 0xb9, 0xfc, 0x3a, 0xe4, 0x9b, 0x35, 0x42, 0x71,
	0x5d, 0xcc, 0x81, 0x36, 0x12, 0xc4, 0x42, 0xaf, 0xa4, 0x35, 0x5e, 0x22, 0xb2, 0x8f, 0x43, 0x9b,
	0x01, 0x48, 0x1a, 0x26, 0x88, 0xd5, 0x7d, 0x96, 0xa4, 0xdc, 0x43, 0x72, 0x80, 0x43, 0x3b, 0xd5,
	0x06, 0x68, 0x54, 0xbd, 0x6c, 0x15, 0xab, 0xb8, 0xc1, 0x6a, 0xa7, 0x92, 0x7b, 0x4c, 0x8e, 0x71,
	0xa8, 0xd2, 0x14, 0x0c, 0x6d, 0x26, 0x88, 0xb5, 0x47, 0x3b, 0x83, 0xdf, 0x1d, 0xc7, 0x25, 0x5d,
	0x7f, 0x96, 0xfb, 0x4e, 0xef, 0x08, 0x77, 0x36, 0xf1, 0x7f, 0x1b, 0x4c, 0xa2, 0x6a, 0xb5, 0xf3,
	0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe5, 0xd6, 0x3b, 0x19, 0x95, 0x01, 0x00, 0x00,
}
