// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: permissions.proto

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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

var E_RequiredFlags = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: ([]string)(nil),
	Field:         20100,
	Name:          "s12.flags.permissions.required_flags",
	Tag:           "bytes,20100,rep,name=required_flags",
	Filename:      "permissions.proto",
}

func init() {
	proto.RegisterExtension(E_RequiredFlags)
}

func init() { proto.RegisterFile("permissions.proto", fileDescriptor_46cca66312ac1c30) }

var fileDescriptor_46cca66312ac1c30 = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x48, 0x2d, 0xca,
	0xcd, 0x2c, 0x2e, 0xce, 0xcc, 0xcf, 0x2b, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x2d,
	0x36, 0x34, 0xd2, 0x4b, 0xcb, 0x49, 0x4c, 0x2f, 0xd6, 0x43, 0x92, 0x94, 0x52, 0x48, 0xcf, 0xcf,
	0x4f, 0xcf, 0x49, 0xd5, 0x07, 0x2b, 0x4a, 0x2a, 0x4d, 0xd3, 0x4f, 0x49, 0x2d, 0x4e, 0x2e, 0xca,
	0x2c, 0x28, 0xc9, 0x2f, 0x82, 0x68, 0xb4, 0x72, 0xe7, 0xe2, 0x2b, 0x4a, 0x2d, 0x2c, 0xcd, 0x2c,
	0x4a, 0x4d, 0x89, 0x07, 0xeb, 0x17, 0x92, 0xd3, 0x83, 0x68, 0xd2, 0x83, 0x69, 0xd2, 0xf3, 0x4d,
	0x2d, 0xc9, 0xc8, 0x4f, 0xf1, 0x2f, 0x28, 0x01, 0x99, 0x29, 0xd1, 0x32, 0x97, 0x51, 0x81, 0x59,
	0x83, 0x33, 0x88, 0x17, 0xa6, 0xcf, 0x0d, 0xa4, 0xcd, 0xc9, 0x2c, 0xca, 0x24, 0x3d, 0xb3, 0x24,
	0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x3f, 0x38, 0x31, 0x2d, 0xb5, 0xa4, 0xd2, 0xb9, 0x34,
	0xa7, 0xa4, 0xb4, 0x28, 0x55, 0xbf, 0xd8, 0xd0, 0x48, 0x17, 0x6c, 0x1a, 0xc2, 0x21, 0xc5, 0x86,
	0x46, 0x60, 0x36, 0x20, 0x00, 0x00, 0xff, 0xff, 0x26, 0x45, 0x89, 0xe0, 0xcd, 0x00, 0x00, 0x00,
}
