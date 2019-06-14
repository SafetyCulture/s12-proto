// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: grpcmock.proto

package s12proto

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	descriptor "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
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

type FieldMock struct {
	// Set the mocked field to a set string value. If repeated will choose one at
	// random.
	String_              []string `protobuf:"bytes,1,rep,name=string" json:"string,omitempty"`
	Word                 *bool    `protobuf:"varint,2,opt,name=word" json:"word,omitempty"`
	Words                *bool    `protobuf:"varint,3,opt,name=words" json:"words,omitempty"`
	Wordsn               *int32   `protobuf:"varint,4,opt,name=wordsn" json:"wordsn,omitempty"`
	Intn                 *int32   `protobuf:"varint,5,opt,name=intn" json:"intn,omitempty"`
	Fullname             *bool    `protobuf:"varint,6,opt,name=fullname" json:"fullname,omitempty"`
	Firstname            *bool    `protobuf:"varint,7,opt,name=firstname" json:"firstname,omitempty"`
	Lastname             *bool    `protobuf:"varint,8,opt,name=lastname" json:"lastname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FieldMock) Reset()         { *m = FieldMock{} }
func (m *FieldMock) String() string { return proto.CompactTextString(m) }
func (*FieldMock) ProtoMessage()    {}
func (*FieldMock) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8c8c9b9c4c355de, []int{0}
}
func (m *FieldMock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldMock.Unmarshal(m, b)
}
func (m *FieldMock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldMock.Marshal(b, m, deterministic)
}
func (m *FieldMock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldMock.Merge(m, src)
}
func (m *FieldMock) XXX_Size() int {
	return xxx_messageInfo_FieldMock.Size(m)
}
func (m *FieldMock) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldMock.DiscardUnknown(m)
}

var xxx_messageInfo_FieldMock proto.InternalMessageInfo

func (m *FieldMock) GetString_() []string {
	if m != nil {
		return m.String_
	}
	return nil
}

func (m *FieldMock) GetWord() bool {
	if m != nil && m.Word != nil {
		return *m.Word
	}
	return false
}

func (m *FieldMock) GetWords() bool {
	if m != nil && m.Words != nil {
		return *m.Words
	}
	return false
}

func (m *FieldMock) GetWordsn() int32 {
	if m != nil && m.Wordsn != nil {
		return *m.Wordsn
	}
	return 0
}

func (m *FieldMock) GetIntn() int32 {
	if m != nil && m.Intn != nil {
		return *m.Intn
	}
	return 0
}

func (m *FieldMock) GetFullname() bool {
	if m != nil && m.Fullname != nil {
		return *m.Fullname
	}
	return false
}

func (m *FieldMock) GetFirstname() bool {
	if m != nil && m.Firstname != nil {
		return *m.Firstname
	}
	return false
}

func (m *FieldMock) GetLastname() bool {
	if m != nil && m.Lastname != nil {
		return *m.Lastname
	}
	return false
}

var E_Field = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*FieldMock)(nil),
	Field:         66300,
	Name:          "grpcmock.field",
	Tag:           "bytes,66300,opt,name=field",
	Filename:      "grpcmock.proto",
}

func init() {
	proto.RegisterType((*FieldMock)(nil), "grpcmock.FieldMock")
	proto.RegisterExtension(E_Field)
}

func init() { proto.RegisterFile("grpcmock.proto", fileDescriptor_f8c8c9b9c4c355de) }

var fileDescriptor_f8c8c9b9c4c355de = []byte{
	// 280 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xcf, 0x4a, 0xf4, 0x30,
	0x14, 0xc5, 0xe9, 0x37, 0xd3, 0xf9, 0xda, 0x08, 0x2e, 0xa2, 0x48, 0x18, 0x14, 0x8a, 0xab, 0x6e,
	0x4c, 0x99, 0x22, 0x2e, 0x5c, 0x2a, 0xb8, 0x10, 0x44, 0xa8, 0x3b, 0x77, 0x9d, 0x36, 0xad, 0x61,
	0xd2, 0xdc, 0x92, 0xa4, 0x88, 0x0f, 0xe0, 0xf3, 0xf9, 0x42, 0x2e, 0x24, 0xb7, 0x7f, 0x66, 0xd5,
	0xf3, 0x3b, 0xb7, 0xf7, 0x90, 0x73, 0xc9, 0x69, 0x6b, 0xfa, 0xaa, 0x83, 0xea, 0xc0, 0x7b, 0x03,
	0x0e, 0x68, 0x34, 0xf3, 0x36, 0x69, 0x01, 0x5a, 0x25, 0x32, 0xf4, 0xf7, 0x43, 0x93, 0xd5, 0xc2,
	0x56, 0x46, 0xf6, 0x0e, 0xcc, 0xf8, 0xef, 0xf5, 0x4f, 0x40, 0xe2, 0x27, 0x29, 0x54, 0xfd, 0x02,
	0xd5, 0x81, 0x5e, 0x90, 0x8d, 0x75, 0x46, 0xea, 0x96, 0x05, 0xc9, 0x2a, 0x8d, 0x8b, 0x89, 0x28,
	0x25, 0xeb, 0x4f, 0x30, 0x35, 0xfb, 0x97, 0x04, 0x69, 0x54, 0xa0, 0xa6, 0xe7, 0x24, 0xf4, 0x5f,
	0xcb, 0x56, 0x68, 0x8e, 0xe0, 0x13, 0x50, 0x68, 0xb6, 0x4e, 0x82, 0x34, 0x2c, 0x26, 0xf2, 0x09,
	0x52, 0x3b, 0xcd, 0x42, 0x74, 0x51, 0xd3, 0x2d, 0x89, 0x9a, 0x41, 0x29, 0x5d, 0x76, 0x82, 0x6d,
	0x30, 0x64, 0x61, 0x7a, 0x49, 0xe2, 0x46, 0x1a, 0xeb, 0x70, 0xf8, 0x1f, 0x87, 0x47, 0xc3, 0x6f,
	0xaa, 0x72, 0x1a, 0x46, 0xe3, 0xe6, 0xcc, 0xf7, 0xcf, 0x24, 0x6c, 0x7c, 0x21, 0x7a, 0xc5, 0xc7,
	0xf6, 0x7c, 0x6e, 0xcf, 0xb1, 0xe8, 0x6b, 0xef, 0x24, 0x68, 0xcb, 0x7e, 0xbf, 0xfd, 0x0b, 0x4f,
	0xf2, 0x33, 0xbe, 0x9c, 0x6f, 0x39, 0x44, 0x31, 0x46, 0x3c, 0xdc, 0xbd, 0xdf, 0xb6, 0xd2, 0x7d,
	0x0c, 0x7b, 0x5e, 0x41, 0x97, 0xbd, 0x95, 0x8d, 0x70, 0x5f, 0x8f, 0x83, 0x72, 0x83, 0x11, 0x99,
	0xdd, 0xe5, 0x37, 0x98, 0x7c, 0xbc, 0xae, 0xdd, 0xe5, 0xa8, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x0c, 0x16, 0xbe, 0xaa, 0x92, 0x01, 0x00, 0x00,
}
