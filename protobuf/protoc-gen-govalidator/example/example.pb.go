// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: example.proto

package example

import (
	fmt "fmt"
	_ "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	math "math"
)

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
	// returns an error if the string cannot be parsed as a UUID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// bytes can also be parsed as UUID with support for gogo
	UserID []byte `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// strings can validate against a regular expresion
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	// integers can be greater than a value
	Age int32 `protobuf:"varint,4,opt,name=age,proto3" json:"age,omitempty"`
	// intergers can be less than a value
	Speed int64 `protobuf:"varint,5,opt,name=speed,proto3" json:"speed,omitempty"`
	// intergers greater/less than or equal, the can also be combined
	Score int32 `protobuf:"varint,6,opt,name=score,proto3" json:"score,omitempty"`
	// validation is created for all messages
	Inner *InnerMessage `protobuf:"bytes,7,opt,name=inner,proto3" json:"inner,omitempty"`
	// can validate each repeated item too
	Ids [][]byte `protobuf:"bytes,8,rep,name=ids,proto3" json:"ids,omitempty"`
	// only validate if non-zero value
	MediaId string `protobuf:"bytes,9,opt,name=media_id,json=mediaId,proto3" json:"media_id,omitempty"`
	// validate the max length of a string
	Description string `protobuf:"bytes,10,opt,name=description,proto3" json:"description,omitempty"`
	// validate the min length
	Password string `protobuf:"bytes,11,opt,name=password,proto3" json:"password,omitempty"`
	// You don't need to validate everything
	NoValidation         string   `protobuf:"bytes,12,opt,name=no_validation,json=noValidation,proto3" json:"no_validation,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExampleMessage) Reset()         { *m = ExampleMessage{} }
func (m *ExampleMessage) String() string { return proto.CompactTextString(m) }
func (*ExampleMessage) ProtoMessage()    {}
func (*ExampleMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{0}
}
func (m *ExampleMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExampleMessage.Unmarshal(m, b)
}
func (m *ExampleMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExampleMessage.Marshal(b, m, deterministic)
}
func (m *ExampleMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExampleMessage.Merge(m, src)
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

func (m *ExampleMessage) GetUserID() []byte {
	if m != nil {
		return m.UserID
	}
	return nil
}

func (m *ExampleMessage) GetEmail() string {
	if m != nil {
		return m.Email
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

func (m *ExampleMessage) GetIds() [][]byte {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *ExampleMessage) GetMediaId() string {
	if m != nil {
		return m.MediaId
	}
	return ""
}

func (m *ExampleMessage) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ExampleMessage) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *ExampleMessage) GetNoValidation() string {
	if m != nil {
		return m.NoValidation
	}
	return ""
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
	return fileDescriptor_15a1dc8d40dadaa6, []int{1}
}
func (m *InnerMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InnerMessage.Unmarshal(m, b)
}
func (m *InnerMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InnerMessage.Marshal(b, m, deterministic)
}
func (m *InnerMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InnerMessage.Merge(m, src)
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

func init() { proto.RegisterFile("example.proto", fileDescriptor_15a1dc8d40dadaa6) }

var fileDescriptor_15a1dc8d40dadaa6 = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xbf, 0x8e, 0xd4, 0x30,
	0x10, 0xc6, 0x2f, 0x9b, 0xbf, 0xe7, 0xcd, 0x2d, 0x92, 0x05, 0xc8, 0xda, 0x02, 0x47, 0x77, 0x27,
	0x11, 0x29, 0x22, 0x27, 0x8e, 0x17, 0x40, 0x11, 0x14, 0x29, 0x68, 0x2c, 0x41, 0x75, 0xd2, 0xc9,
	0xbb, 0x36, 0xc1, 0xd2, 0x26, 0x8e, 0xec, 0x2c, 0x50, 0x52, 0x53, 0x51, 0x52, 0x22, 0xd1, 0x50,
	0x52, 0x52, 0xf2, 0x08, 0x3c, 0x01, 0x05, 0xad, 0x5f, 0x02, 0xd9, 0x0e, 0xab, 0x2d, 0xe8, 0x3c,
	0xbf, 0xef, 0x9b, 0x91, 0xbf, 0x19, 0x70, 0xc6, 0xdf, 0xd3, 0x7e, 0xdc, 0xf1, 0x7a, 0x54, 0x72,
	0x92, 0x30, 0x9d, 0xcb, 0xf5, 0xa3, 0x4e, 0x4c, 0x6f, 0xf6, 0x9b, 0x7a, 0x2b, 0xfb, 0xab, 0x4e,
	0x76, 0xf2, 0xca, 0xe9, 0x9b, 0xfd, 0x6b, 0x57, 0xb9, 0xc2, 0xbd, 0x7c, 0xdf, 0x1a, 0xe9, 0xc7,
	0xd7, 0x1e, 0xbe, 0xa5, 0x3b, 0xc1, 0xe8, 0x24, 0x95, 0x57, 0xce, 0xbf, 0x86, 0x60, 0xf5, 0xdc,
	0x0f, 0x7d, 0xc1, 0xb5, 0xa6, 0x1d, 0x87, 0x77, 0xc1, 0x42, 0x30, 0x14, 0x14, 0x41, 0x79, 0xda,
	0x44, 0x1f, 0x0c, 0x0e, 0xc8, 0x42, 0x30, 0xf8, 0x10, 0xa4, 0x7b, 0xcd, 0xd5, 0xad, 0x60, 0x68,
	0x51, 0x04, 0x65, 0xde, 0xac, 0xfe, 0xfc, 0xc6, 0xc9, 0x4b, 0xcd, 0x55, 0xfb, 0xcc, 0x99, 0x12,
	0x2b, 0xb7, 0x0c, 0x5e, 0x82, 0x98, 0xf7, 0x54, 0xec, 0x50, 0xe8, 0x26, 0xac, 0x3e, 0x1a, 0x0c,
	0xea, 0xea, 0xe6, 0x69, 0x5d, 0xdd, 0xd4, 0x75, 0x45, 0xbc, 0x08, 0xef, 0x83, 0x90, 0x76, 0x1c,
	0x45, 0x45, 0x50, 0xc6, 0x4d, 0xf4, 0xc9, 0xe0, 0x13, 0x62, 0x01, 0x5c, 0x83, 0x58, 0x8f, 0x9c,
	0x33, 0x14, 0x17, 0x41, 0x19, 0x36, 0xd1, 0x67, 0x83, 0x07, 0xe2, 0x11, 0x7c, 0x00, 0x62, 0xbd,
	0x95, 0x8a, 0xa3, 0xc4, 0x75, 0x65, 0x5f, 0x0c, 0x3e, 0xf9, 0x66, 0x30, 0x23, 0x1e, 0xc3, 0x0a,
	0xc4, 0x62, 0x18, 0xb8, 0x42, 0x69, 0x11, 0x94, 0xcb, 0xeb, 0x7b, 0xf5, 0xbf, 0xe5, 0xb5, 0x96,
	0xce, 0xf1, 0x88, 0xf7, 0xd8, 0x0f, 0x08, 0xa6, 0x51, 0x56, 0x84, 0x65, 0x3e, 0xc7, 0xb4, 0x00,
	0x5e, 0x80, 0xac, 0xe7, 0x4c, 0x50, 0x1b, 0xf4, 0xd4, 0x25, 0xc8, 0xac, 0xf8, 0xd3, 0x1a, 0x52,
	0xa7, 0xb4, 0x76, 0x19, 0x4b, 0xc6, 0xf5, 0x56, 0x89, 0x71, 0x12, 0x72, 0x40, 0xc0, 0xf9, 0xe2,
	0x1f, 0x06, 0xff, 0xba, 0x43, 0x8e, 0x15, 0x58, 0x80, 0x6c, 0xa4, 0x5a, 0xbf, 0x93, 0x8a, 0xa1,
	0xa5, 0xdf, 0xe8, 0x77, 0x83, 0x33, 0x72, 0xa0, 0xf0, 0x02, 0x9c, 0x0d, 0xf2, 0x76, 0x3e, 0x8b,
	0x1d, 0x96, 0x5b, 0x1b, 0xc9, 0x07, 0xf9, 0xea, 0xc0, 0xce, 0x2f, 0x41, 0x7e, 0x9c, 0xe1, 0xff,
	0x27, 0xda, 0x24, 0xee, 0xa4, 0x4f, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0x9c, 0xb6, 0x07, 0x82,
	0x35, 0x02, 0x00, 0x00,
}
