// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: other.proto

package s12_routeguide

import (
	fmt "fmt"
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

type Other struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Other) Reset()         { *m = Other{} }
func (m *Other) String() string { return proto.CompactTextString(m) }
func (*Other) ProtoMessage()    {}
func (*Other) Descriptor() ([]byte, []int) {
	return fileDescriptor_b402626a50d68a3b, []int{0}
}
func (m *Other) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Other.Unmarshal(m, b)
}
func (m *Other) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Other.Marshal(b, m, deterministic)
}
func (m *Other) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Other.Merge(m, src)
}
func (m *Other) XXX_Size() int {
	return xxx_messageInfo_Other.Size(m)
}
func (m *Other) XXX_DiscardUnknown() {
	xxx_messageInfo_Other.DiscardUnknown(m)
}

var xxx_messageInfo_Other proto.InternalMessageInfo

func (m *Other) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*Other)(nil), "s12.routeguide.Other")
}

func init() { proto.RegisterFile("other.proto", fileDescriptor_b402626a50d68a3b) }

var fileDescriptor_b402626a50d68a3b = []byte{
	// 81 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x2f, 0xc9, 0x48,
	0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2b, 0x36, 0x34, 0xd2, 0x2b, 0xca, 0x2f,
	0x2d, 0x49, 0x4d, 0x2f, 0xcd, 0x4c, 0x49, 0x55, 0x12, 0xe7, 0x62, 0xf5, 0x07, 0x49, 0x0b, 0xf1,
	0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x31, 0x65, 0xa6, 0x24, 0xb1,
	0x81, 0xd5, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xae, 0x22, 0x8f, 0x13, 0x3e, 0x00, 0x00,
	0x00,
}
