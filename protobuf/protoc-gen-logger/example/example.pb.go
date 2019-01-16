// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: example.proto

package example

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
import _ "github.com/gogo/protobuf/gogoproto"

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
	Id                   string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserName             string        `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Password             string        `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	SomeKindOfInnerValue *InnerMessage `protobuf:"bytes,4,opt,name=some_kind_of_inner_value,json=someKindOfInnerValue,proto3" json:"some_kind_of_inner_value,omitempty"`
	// Types that are valid to be assigned to TestOneof:
	//	*ExampleMessage_OneOf1
	//	*ExampleMessage_OneOf2
	TestOneof            isExampleMessage_TestOneof `protobuf_oneof:"test_oneof"`
	MapField             map[string]string          `protobuf:"bytes,7,rep,name=map_field,json=mapField,proto3" json:"map_field,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ExampleMessage) Reset()         { *m = ExampleMessage{} }
func (m *ExampleMessage) String() string { return proto.CompactTextString(m) }
func (*ExampleMessage) ProtoMessage()    {}
func (*ExampleMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_example_293142d358651c05, []int{0}
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

type isExampleMessage_TestOneof interface {
	isExampleMessage_TestOneof()
}

type ExampleMessage_OneOf1 struct {
	OneOf1 string `protobuf:"bytes,5,opt,name=one_of1,json=oneOf1,proto3,oneof"`
}
type ExampleMessage_OneOf2 struct {
	OneOf2 *OneOfMessage `protobuf:"bytes,6,opt,name=one_of2,json=oneOf2,proto3,oneof"`
}

func (*ExampleMessage_OneOf1) isExampleMessage_TestOneof() {}
func (*ExampleMessage_OneOf2) isExampleMessage_TestOneof() {}

func (m *ExampleMessage) GetTestOneof() isExampleMessage_TestOneof {
	if m != nil {
		return m.TestOneof
	}
	return nil
}

func (m *ExampleMessage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ExampleMessage) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *ExampleMessage) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *ExampleMessage) GetSomeKindOfInnerValue() *InnerMessage {
	if m != nil {
		return m.SomeKindOfInnerValue
	}
	return nil
}

func (m *ExampleMessage) GetOneOf1() string {
	if x, ok := m.GetTestOneof().(*ExampleMessage_OneOf1); ok {
		return x.OneOf1
	}
	return ""
}

func (m *ExampleMessage) GetOneOf2() *OneOfMessage {
	if x, ok := m.GetTestOneof().(*ExampleMessage_OneOf2); ok {
		return x.OneOf2
	}
	return nil
}

func (m *ExampleMessage) GetMapField() map[string]string {
	if m != nil {
		return m.MapField
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ExampleMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ExampleMessage_OneofMarshaler, _ExampleMessage_OneofUnmarshaler, _ExampleMessage_OneofSizer, []interface{}{
		(*ExampleMessage_OneOf1)(nil),
		(*ExampleMessage_OneOf2)(nil),
	}
}

func _ExampleMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ExampleMessage)
	// test_oneof
	switch x := m.TestOneof.(type) {
	case *ExampleMessage_OneOf1:
		_ = b.EncodeVarint(5<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.OneOf1)
	case *ExampleMessage_OneOf2:
		_ = b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.OneOf2); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ExampleMessage.TestOneof has unexpected type %T", x)
	}
	return nil
}

func _ExampleMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ExampleMessage)
	switch tag {
	case 5: // test_oneof.one_of1
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.TestOneof = &ExampleMessage_OneOf1{x}
		return true, err
	case 6: // test_oneof.one_of2
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(OneOfMessage)
		err := b.DecodeMessage(msg)
		m.TestOneof = &ExampleMessage_OneOf2{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ExampleMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ExampleMessage)
	// test_oneof
	switch x := m.TestOneof.(type) {
	case *ExampleMessage_OneOf1:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.OneOf1)))
		n += len(x.OneOf1)
	case *ExampleMessage_OneOf2:
		s := proto.Size(x.OneOf2)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type OneOfMessage struct {
	Value                int32    `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OneOfMessage) Reset()         { *m = OneOfMessage{} }
func (m *OneOfMessage) String() string { return proto.CompactTextString(m) }
func (*OneOfMessage) ProtoMessage()    {}
func (*OneOfMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_example_293142d358651c05, []int{1}
}
func (m *OneOfMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OneOfMessage.Unmarshal(m, b)
}
func (m *OneOfMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OneOfMessage.Marshal(b, m, deterministic)
}
func (dst *OneOfMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OneOfMessage.Merge(dst, src)
}
func (m *OneOfMessage) XXX_Size() int {
	return xxx_messageInfo_OneOfMessage.Size(m)
}
func (m *OneOfMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_OneOfMessage.DiscardUnknown(m)
}

var xxx_messageInfo_OneOfMessage proto.InternalMessageInfo

func (m *OneOfMessage) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type InnerMessage struct {
	Body                 string   `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InnerMessage) Reset()         { *m = InnerMessage{} }
func (m *InnerMessage) String() string { return proto.CompactTextString(m) }
func (*InnerMessage) ProtoMessage()    {}
func (*InnerMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_example_293142d358651c05, []int{2}
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

func (m *InnerMessage) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*ExampleMessage)(nil), "example.ExampleMessage")
	proto.RegisterMapType((map[string]string)(nil), "example.ExampleMessage.MapFieldEntry")
	proto.RegisterType((*OneOfMessage)(nil), "example.OneOfMessage")
	proto.RegisterType((*InnerMessage)(nil), "example.InnerMessage")
}

func init() { proto.RegisterFile("example.proto", fileDescriptor_example_293142d358651c05) }

var fileDescriptor_example_293142d358651c05 = []byte{
	// 399 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0xdd, 0x8a, 0xd4, 0x30,
	0x14, 0x76, 0x3a, 0x7f, 0x3b, 0x67, 0x67, 0x45, 0xc2, 0x08, 0x71, 0xbc, 0x29, 0x83, 0xc2, 0xdc,
	0xec, 0x0c, 0xad, 0x20, 0xa2, 0x57, 0xae, 0xac, 0xac, 0xc8, 0x3a, 0xd0, 0x05, 0x6f, 0x43, 0x6a,
	0x4f, 0x6a, 0xd9, 0x36, 0x29, 0x4d, 0xab, 0xf6, 0x75, 0xbc, 0xf3, 0x2d, 0x7c, 0x34, 0x49, 0xd2,
	0xa9, 0x9d, 0xbd, 0xcb, 0xf9, 0xfe, 0x92, 0x7c, 0x07, 0x2e, 0xf0, 0x17, 0x2f, 0xca, 0x1c, 0x77,
	0x65, 0xa5, 0x6a, 0x45, 0xe6, 0xdd, 0xb8, 0xbe, 0x4c, 0xb3, 0xfa, 0x7b, 0x13, 0xef, 0xbe, 0xa9,
	0x62, 0x9f, 0xaa, 0x54, 0xed, 0x2d, 0x1f, 0x37, 0xc2, 0x4e, 0x76, 0xb0, 0x27, 0xe7, 0x5b, 0xbf,
	0x1f, 0xc8, 0xef, 0xb8, 0xc0, 0xba, 0xfd, 0xd0, 0xe4, 0x75, 0x53, 0xe1, 0x5e, 0x07, 0xe1, 0xa5,
	0x93, 0xf7, 0x09, 0x3a, 0x08, 0x1d, 0x92, 0xab, 0x34, 0xc5, 0xca, 0x45, 0x6c, 0x7e, 0x8f, 0xe1,
	0xf1, 0xb5, 0xbb, 0xfd, 0x16, 0xb5, 0xe6, 0x29, 0x92, 0x15, 0x78, 0x59, 0x42, 0x47, 0xfe, 0x68,
	0xbb, 0xb8, 0x9a, 0xfc, 0xfd, 0xe3, 0x4f, 0x23, 0x2f, 0x4b, 0xc8, 0x73, 0x58, 0x34, 0x1a, 0x2b,
	0x26, 0x79, 0x81, 0xd4, 0x33, 0x64, 0x74, 0x66, 0x80, 0x2f, 0xbc, 0x40, 0xe2, 0xc3, 0x59, 0xc9,
	0xb5, 0xfe, 0xa9, 0xaa, 0x84, 0x8e, 0x7b, 0xa3, 0x17, 0xf5, 0x28, 0xb9, 0x03, 0xaa, 0x55, 0x81,
	0xec, 0x3e, 0x93, 0x09, 0x53, 0x82, 0x65, 0x52, 0x62, 0xc5, 0x7e, 0xf0, 0xbc, 0x41, 0x3a, 0xf1,
	0x47, 0xdb, 0xf3, 0xf0, 0xe9, 0xee, 0x58, 0xca, 0x27, 0xc3, 0x75, 0xaf, 0xb1, 0x41, 0x93, 0x68,
	0x65, 0xcc, 0x9f, 0x33, 0x99, 0x1c, 0x84, 0x65, 0xbf, 0x1a, 0x23, 0x79, 0x06, 0x73, 0x25, 0x91,
	0x29, 0x11, 0xd0, 0xa9, 0xb9, 0xf5, 0xe6, 0x51, 0x34, 0x53, 0x12, 0x0f, 0x22, 0x20, 0xaf, 0x8f,
	0x54, 0x48, 0x67, 0x0f, 0xe2, 0x0f, 0x46, 0x31, 0x8c, 0xf7, 0x7a, 0x5f, 0x48, 0x6e, 0x60, 0x51,
	0xf0, 0x92, 0x89, 0x0c, 0xf3, 0x84, 0xce, 0xfd, 0xf1, 0xf6, 0x3c, 0x7c, 0xd9, 0x3b, 0x4f, 0x8b,
	0xda, 0xdd, 0xf2, 0xf2, 0xa3, 0xd1, 0x5d, 0xcb, 0xba, 0x6a, 0x8f, 0x3f, 0x2e, 0x3a, 0x70, 0xfd,
	0x0e, 0x2e, 0x4e, 0x04, 0xe4, 0x09, 0x8c, 0xef, 0xb1, 0x75, 0xc5, 0x46, 0xe6, 0x48, 0x56, 0x30,
	0x75, 0x0d, 0xb8, 0x3e, 0xdd, 0xf0, 0xd6, 0x7b, 0x33, 0xba, 0x5a, 0x02, 0xd4, 0xa8, 0x6b, 0xa6,
	0x24, 0x2a, 0xb1, 0x79, 0x01, 0xcb, 0xe1, 0xa3, 0xff, 0xfb, 0x4c, 0xd6, 0xb4, 0xf3, 0x6d, 0xb6,
	0xb0, 0x1c, 0x36, 0x47, 0x28, 0x4c, 0x62, 0x95, 0xb4, 0x27, 0x9b, 0xb4, 0x48, 0x3c, 0xb3, 0xbb,
	0x7f, 0xf5, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xad, 0xc8, 0xf8, 0xf9, 0x87, 0x02, 0x00, 0x00,
}
