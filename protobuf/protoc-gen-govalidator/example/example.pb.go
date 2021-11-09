// Copyright (c) 2018 SafetyCulture Pty Ltd. All Rights Reserved.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: example.proto

package example

import (
	_ "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MyMessageWithEnum_MyEnum int32

const (
	// Default is unspecified.
	MyMessageWithEnum_MY_ENUM_UNSPECIFIED MyMessageWithEnum_MyEnum = 0
	// First valid value.
	MyMessageWithEnum_MY_ENUM_FIRST MyMessageWithEnum_MyEnum = 1
)

// Enum value maps for MyMessageWithEnum_MyEnum.
var (
	MyMessageWithEnum_MyEnum_name = map[int32]string{
		0: "MY_ENUM_UNSPECIFIED",
		1: "MY_ENUM_FIRST",
	}
	MyMessageWithEnum_MyEnum_value = map[string]int32{
		"MY_ENUM_UNSPECIFIED": 0,
		"MY_ENUM_FIRST":       1,
	}
)

func (x MyMessageWithEnum_MyEnum) Enum() *MyMessageWithEnum_MyEnum {
	p := new(MyMessageWithEnum_MyEnum)
	*p = x
	return p
}

func (x MyMessageWithEnum_MyEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MyMessageWithEnum_MyEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_example_proto_enumTypes[0].Descriptor()
}

func (MyMessageWithEnum_MyEnum) Type() protoreflect.EnumType {
	return &file_example_proto_enumTypes[0]
}

func (x MyMessageWithEnum_MyEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MyMessageWithEnum_MyEnum.Descriptor instead.
func (MyMessageWithEnum_MyEnum) EnumDescriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{4, 0}
}

type ExampleMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// returns an error if the string cannot be parsed as a UUID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// bytes can also be parsed as UUID with support for gogo
	UserId []byte `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// should be a valid email address
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	// integers can be greater than a value
	Age int32 `protobuf:"varint,4,opt,name=age,proto3" json:"age,omitempty"`
	// integers can be less than a value
	Speed int64 `protobuf:"varint,5,opt,name=speed,proto3" json:"speed,omitempty"`
	// integers greater/less than or equal, the can also be combined
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
	NoValidation string `protobuf:"bytes,12,opt,name=no_validation,json=noValidation,proto3" json:"no_validation,omitempty"`
	// Types that are assignable to ContactOneof:
	//	*ExampleMessage_Fax
	//	*ExampleMessage_Phone
	ContactOneof isExampleMessage_ContactOneof `protobuf_oneof:"contact_oneof"`
	// Required inner message
	MsgRequired *InnerMessage `protobuf:"bytes,15,opt,name=msg_required,json=msgRequired,proto3" json:"msg_required,omitempty"`
	// returns an error if the string cannot be parsed as a legacy id
	LegacyId string `protobuf:"bytes,16,opt,name=legacy_id,json=legacyId,proto3" json:"legacy_id,omitempty"`
	// InnerMessage can contain a legacy id too
	InnerLegacyId *InnerMessageWithLegacyId `protobuf:"bytes,17,opt,name=inner_legacy_id,json=innerLegacyId,proto3" json:"inner_legacy_id,omitempty"`
	// Trim leading and trailing whitespaces (as defined by Unicode) before doing length check
	Name          string                        `protobuf:"bytes,18,opt,name=name,proto3" json:"name,omitempty"`
	NestedMessage *ExampleMessage_NestedMessage `protobuf:"bytes,19,opt,name=nested_message,json=nestedMessage,proto3" json:"nested_message,omitempty"`
	NotSupported  map[string]string             `protobuf:"bytes,100,rep,name=not_supported,json=notSupported,proto3" json:"not_supported,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// strings can validate against a regular expresion
	Url                            string                    `protobuf:"bytes,21,opt,name=url,proto3" json:"url,omitempty"`
	ContactsWithLengthConstraint   []*ExampleMessage_Contact `protobuf:"bytes,22,rep,name=contacts_with_length_constraint,json=contactsWithLengthConstraint,proto3" json:"contacts_with_length_constraint,omitempty"`
	ContactsWithNoLengthConstraint []*ExampleMessage_Contact `protobuf:"bytes,23,rep,name=contacts_with_no_length_constraint,json=contactsWithNoLengthConstraint,proto3" json:"contacts_with_no_length_constraint,omitempty"`
	ScheduledFor                   *timestamppb.Timestamp    `protobuf:"bytes,24,opt,name=scheduled_for,json=scheduledFor,proto3" json:"scheduled_for,omitempty"`
}

func (x *ExampleMessage) Reset() {
	*x = ExampleMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExampleMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExampleMessage) ProtoMessage() {}

func (x *ExampleMessage) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExampleMessage.ProtoReflect.Descriptor instead.
func (*ExampleMessage) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0}
}

func (x *ExampleMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ExampleMessage) GetUserId() []byte {
	if x != nil {
		return x.UserId
	}
	return nil
}

func (x *ExampleMessage) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ExampleMessage) GetAge() int32 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *ExampleMessage) GetSpeed() int64 {
	if x != nil {
		return x.Speed
	}
	return 0
}

func (x *ExampleMessage) GetScore() int32 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *ExampleMessage) GetInner() *InnerMessage {
	if x != nil {
		return x.Inner
	}
	return nil
}

func (x *ExampleMessage) GetIds() [][]byte {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ExampleMessage) GetMediaId() string {
	if x != nil {
		return x.MediaId
	}
	return ""
}

func (x *ExampleMessage) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ExampleMessage) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *ExampleMessage) GetNoValidation() string {
	if x != nil {
		return x.NoValidation
	}
	return ""
}

func (m *ExampleMessage) GetContactOneof() isExampleMessage_ContactOneof {
	if m != nil {
		return m.ContactOneof
	}
	return nil
}

func (x *ExampleMessage) GetFax() string {
	if x, ok := x.GetContactOneof().(*ExampleMessage_Fax); ok {
		return x.Fax
	}
	return ""
}

func (x *ExampleMessage) GetPhone() string {
	if x, ok := x.GetContactOneof().(*ExampleMessage_Phone); ok {
		return x.Phone
	}
	return ""
}

func (x *ExampleMessage) GetMsgRequired() *InnerMessage {
	if x != nil {
		return x.MsgRequired
	}
	return nil
}

func (x *ExampleMessage) GetLegacyId() string {
	if x != nil {
		return x.LegacyId
	}
	return ""
}

func (x *ExampleMessage) GetInnerLegacyId() *InnerMessageWithLegacyId {
	if x != nil {
		return x.InnerLegacyId
	}
	return nil
}

func (x *ExampleMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ExampleMessage) GetNestedMessage() *ExampleMessage_NestedMessage {
	if x != nil {
		return x.NestedMessage
	}
	return nil
}

func (x *ExampleMessage) GetNotSupported() map[string]string {
	if x != nil {
		return x.NotSupported
	}
	return nil
}

func (x *ExampleMessage) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ExampleMessage) GetContactsWithLengthConstraint() []*ExampleMessage_Contact {
	if x != nil {
		return x.ContactsWithLengthConstraint
	}
	return nil
}

func (x *ExampleMessage) GetContactsWithNoLengthConstraint() []*ExampleMessage_Contact {
	if x != nil {
		return x.ContactsWithNoLengthConstraint
	}
	return nil
}

func (x *ExampleMessage) GetScheduledFor() *timestamppb.Timestamp {
	if x != nil {
		return x.ScheduledFor
	}
	return nil
}

type isExampleMessage_ContactOneof interface {
	isExampleMessage_ContactOneof()
}

type ExampleMessage_Fax struct {
	Fax string `protobuf:"bytes,13,opt,name=fax,proto3,oneof"`
}

type ExampleMessage_Phone struct {
	Phone string `protobuf:"bytes,14,opt,name=phone,proto3,oneof"`
}

func (*ExampleMessage_Fax) isExampleMessage_ContactOneof() {}

func (*ExampleMessage_Phone) isExampleMessage_ContactOneof() {}

type OuterMessageUsingNestedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SomeMessage *ExampleMessage_NestedMessage `protobuf:"bytes,1,opt,name=some_message,json=someMessage,proto3" json:"some_message,omitempty"`
}

func (x *OuterMessageUsingNestedMessage) Reset() {
	*x = OuterMessageUsingNestedMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OuterMessageUsingNestedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OuterMessageUsingNestedMessage) ProtoMessage() {}

func (x *OuterMessageUsingNestedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OuterMessageUsingNestedMessage.ProtoReflect.Descriptor instead.
func (*OuterMessageUsingNestedMessage) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{1}
}

func (x *OuterMessageUsingNestedMessage) GetSomeMessage() *ExampleMessage_NestedMessage {
	if x != nil {
		return x.SomeMessage
	}
	return nil
}

type InnerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *InnerMessage) Reset() {
	*x = InnerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InnerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InnerMessage) ProtoMessage() {}

func (x *InnerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InnerMessage.ProtoReflect.Descriptor instead.
func (*InnerMessage) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{2}
}

func (x *InnerMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type InnerMessageWithLegacyId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *InnerMessageWithLegacyId) Reset() {
	*x = InnerMessageWithLegacyId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InnerMessageWithLegacyId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InnerMessageWithLegacyId) ProtoMessage() {}

func (x *InnerMessageWithLegacyId) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InnerMessageWithLegacyId.ProtoReflect.Descriptor instead.
func (*InnerMessageWithLegacyId) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{3}
}

func (x *InnerMessageWithLegacyId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type MyMessageWithEnum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enum MyMessageWithEnum_MyEnum `protobuf:"varint,1,opt,name=enum,proto3,enum=example.MyMessageWithEnum_MyEnum" json:"enum,omitempty"`
}

func (x *MyMessageWithEnum) Reset() {
	*x = MyMessageWithEnum{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MyMessageWithEnum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MyMessageWithEnum) ProtoMessage() {}

func (x *MyMessageWithEnum) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MyMessageWithEnum.ProtoReflect.Descriptor instead.
func (*MyMessageWithEnum) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{4}
}

func (x *MyMessageWithEnum) GetEnum() MyMessageWithEnum_MyEnum {
	if x != nil {
		return x.Enum
	}
	return MyMessageWithEnum_MY_ENUM_UNSPECIFIED
}

type ExampleMessage_NestedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val          string   `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
	NestedEmail  string   `protobuf:"bytes,2,opt,name=nested_email,json=nestedEmail,proto3" json:"nested_email,omitempty"`
	MemberEmails []string `protobuf:"bytes,3,rep,name=member_emails,json=memberEmails,proto3" json:"member_emails,omitempty"`
}

func (x *ExampleMessage_NestedMessage) Reset() {
	*x = ExampleMessage_NestedMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExampleMessage_NestedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExampleMessage_NestedMessage) ProtoMessage() {}

func (x *ExampleMessage_NestedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExampleMessage_NestedMessage.ProtoReflect.Descriptor instead.
func (*ExampleMessage_NestedMessage) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ExampleMessage_NestedMessage) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *ExampleMessage_NestedMessage) GetNestedEmail() string {
	if x != nil {
		return x.NestedEmail
	}
	return ""
}

func (x *ExampleMessage_NestedMessage) GetMemberEmails() []string {
	if x != nil {
		return x.MemberEmails
	}
	return nil
}

type ExampleMessage_Contact struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Phone string `protobuf:"bytes,1,opt,name=phone,proto3" json:"phone,omitempty"`
	Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *ExampleMessage_Contact) Reset() {
	*x = ExampleMessage_Contact{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExampleMessage_Contact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExampleMessage_Contact) ProtoMessage() {}

func (x *ExampleMessage_Contact) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExampleMessage_Contact.ProtoReflect.Descriptor instead.
func (*ExampleMessage_Contact) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0, 2}
}

func (x *ExampleMessage_Contact) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *ExampleMessage_Contact) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type ExampleMessage_NestedMessage_InnerNestedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InnerVal string `protobuf:"bytes,1,opt,name=inner_val,json=innerVal,proto3" json:"inner_val,omitempty"`
}

func (x *ExampleMessage_NestedMessage_InnerNestedMessage) Reset() {
	*x = ExampleMessage_NestedMessage_InnerNestedMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExampleMessage_NestedMessage_InnerNestedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExampleMessage_NestedMessage_InnerNestedMessage) ProtoMessage() {}

func (x *ExampleMessage_NestedMessage_InnerNestedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExampleMessage_NestedMessage_InnerNestedMessage.ProtoReflect.Descriptor instead.
func (*ExampleMessage_NestedMessage_InnerNestedMessage) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0, 0, 0}
}

func (x *ExampleMessage_NestedMessage_InnerNestedMessage) GetInnerVal() string {
	if x != nil {
		return x.InnerVal
	}
	return ""
}

var File_example_proto protoreflect.FileDescriptor

var file_example_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x22, 0x73, 0x31, 0x32, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x0d,
	0x0a, 0x0e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x14, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0x80, 0xeb,
	0x1f, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x04, 0x80, 0xeb, 0x1f, 0x01, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xf2, 0xeb, 0x1f, 0x02, 0x08, 0x01, 0x52, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x42, 0x04, 0x90, 0xeb, 0x1f, 0x00, 0x52, 0x03, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x05, 0x73,
	0x70, 0x65, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x42, 0x04, 0x98, 0xeb, 0x1f, 0x6e,
	0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x42, 0x08, 0xa0, 0xeb, 0x1f, 0x00, 0xa8, 0xeb, 0x1f, 0x64,
	0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x2b, 0x0a, 0x05, 0x69, 0x6e, 0x6e, 0x65, 0x72,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x2e, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x05, 0x69,
	0x6e, 0x6e, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x0c, 0x42, 0x04, 0x80, 0xeb, 0x1f, 0x01, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x23, 0x0a, 0x08,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0x80, 0xeb, 0x1f, 0x01, 0xc0, 0xeb, 0x1f, 0x01, 0x52, 0x07, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x49,
	0x64, 0x12, 0x27, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x42, 0x05, 0xb8, 0xeb, 0x1f, 0xd0, 0x0f, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xb0, 0xeb,
	0x1f, 0x08, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x6e, 0x6f, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x6e, 0x6f, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x03, 0x66, 0x61, 0x78, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x03, 0x66, 0x61, 0x78, 0x12, 0x1c, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x0e,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xb0, 0xeb, 0x1f, 0x0b, 0x48, 0x00, 0x52, 0x05, 0x70, 0x68,
	0x6f, 0x6e, 0x65, 0x12, 0x3e, 0x0a, 0x0c, 0x6d, 0x73, 0x67, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x42, 0x04, 0xc8, 0xeb, 0x1f, 0x01, 0x52, 0x0b, 0x6d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x09, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x5f, 0x69, 0x64,
	0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xd0, 0xeb, 0x1f, 0x01, 0x52, 0x08, 0x6c, 0x65,
	0x67, 0x61, 0x63, 0x79, 0x49, 0x64, 0x12, 0x49, 0x0a, 0x0f, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f,
	0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x21, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79,
	0x49, 0x64, 0x52, 0x0d, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x49,
	0x64, 0x12, 0x20, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x0c, 0xb0, 0xeb, 0x1f, 0x06, 0xb8, 0xeb, 0x1f, 0x0a, 0xd8, 0xeb, 0x1f, 0x01, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x4c, 0x0a, 0x0e, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x0d, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x4e, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x5f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x64, 0x18, 0x64, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x4e, 0x6f, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x12, 0x3c, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2a,
	0x8a, 0xeb, 0x1f, 0x26, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x5c, 0x2f, 0x5c, 0x2f, 0x77, 0x77,
	0x77, 0x5c, 0x2e, 0x73, 0x61, 0x66, 0x65, 0x74, 0x79, 0x63, 0x75, 0x6c, 0x74, 0x75, 0x72, 0x65,
	0x5c, 0x2e, 0x28, 0x69, 0x6f, 0x7c, 0x63, 0x6f, 0x6d, 0x29, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x70, 0x0a, 0x1f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x5f, 0x77, 0x69, 0x74, 0x68,
	0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x5f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69,
	0x6e, 0x74, 0x18, 0x16, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x42, 0x08, 0xe0, 0xeb, 0x1f, 0x00, 0xe8,
	0xeb, 0x1f, 0x0a, 0x52, 0x1c, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x57, 0x69, 0x74,
	0x68, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e,
	0x74, 0x12, 0x6b, 0x0a, 0x22, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x5f, 0x77, 0x69,
	0x74, 0x68, 0x5f, 0x6e, 0x6f, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x5f, 0x63, 0x6f, 0x6e,
	0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x18, 0x17, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x1e,
	0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x57, 0x69, 0x74, 0x68, 0x4e, 0x6f, 0x4c, 0x65,
	0x6e, 0x67, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x12, 0x45,
	0x0a, 0x0d, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x5f, 0x66, 0x6f, 0x72, 0x18,
	0x18, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x42, 0x04, 0xc0, 0xeb, 0x1f, 0x01, 0x52, 0x0c, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x64, 0x46, 0x6f, 0x72, 0x1a, 0xda, 0x02, 0x0a, 0x0d, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x0c, 0xb0, 0xeb, 0x1f, 0x01, 0xb8, 0xeb, 0x1f, 0x28, 0xd8, 0xeb,
	0x1f, 0x01, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x12, 0x31, 0x0a, 0x0c, 0x6e, 0x65, 0x73, 0x74, 0x65,
	0x64, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0x8a,
	0xeb, 0x1f, 0x0a, 0x2e, 0x2b, 0x5c, 0x40, 0x2e, 0x2b, 0x5c, 0x2e, 0x2e, 0x2b, 0x52, 0x0b, 0x6e,
	0x65, 0x73, 0x74, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0xb4, 0x01, 0x0a, 0x0d, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x09, 0x42, 0x8e, 0x01, 0x8a, 0xeb, 0x1f, 0x81, 0x01, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d,
	0x39, 0x21, 0x23, 0x24, 0x26, 0x27, 0x2a, 0x2b, 0x2f, 0x3d, 0x3f, 0x5e, 0x5f, 0x7b, 0x7c, 0x7d,
	0x7e, 0x2d, 0x5d, 0x2b, 0x28, 0x3f, 0x3a, 0x5c, 0x2e, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39,
	0x21, 0x23, 0x24, 0x26, 0x27, 0x2a, 0x2b, 0x2f, 0x3d, 0x3f, 0x5e, 0x5f, 0x7b, 0x7c, 0x7d, 0x7e,
	0x2d, 0x5d, 0x2b, 0x29, 0x2a, 0x40, 0x28, 0x3f, 0x3a, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39,
	0x5d, 0x28, 0x3f, 0x3a, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39, 0x2d, 0x5d, 0x2a, 0x5b, 0x61,
	0x2d, 0x7a, 0x30, 0x2d, 0x39, 0x5d, 0x29, 0x3f, 0x5c, 0x2e, 0x29, 0x2b, 0x5b, 0x61, 0x2d, 0x7a,
	0x30, 0x2d, 0x39, 0x5d, 0x28, 0x3f, 0x3a, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39, 0x2d, 0x5d,
	0x2a, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39, 0x5d, 0x29, 0x3f, 0xe0, 0xeb, 0x1f, 0x02, 0xe8,
	0xeb, 0x1f, 0x05, 0x52, 0x0c, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x73, 0x1a, 0x3f, 0x0a, 0x12, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x29, 0x0a, 0x09, 0x69, 0x6e, 0x6e, 0x65, 0x72,
	0x5f, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0c, 0xb0, 0xeb, 0x1f, 0x01,
	0xb8, 0xeb, 0x1f, 0x28, 0xd8, 0xeb, 0x1f, 0x01, 0x52, 0x08, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x56,
	0x61, 0x6c, 0x1a, 0x3f, 0x0a, 0x11, 0x4e, 0x6f, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x1a, 0x43, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x12, 0x1a,
	0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xc0,
	0xeb, 0x1f, 0x01, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x1c, 0x0a, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xf2, 0xeb, 0x1f, 0x02, 0x08,
	0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x0f, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x5f, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x22, 0x6a, 0x0a, 0x1e, 0x4f, 0x75, 0x74,
	0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x55, 0x73, 0x69, 0x6e, 0x67, 0x4e, 0x65,
	0x73, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x48, 0x0a, 0x0c, 0x73,
	0x6f, 0x6d, 0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x25, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x45, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65,
	0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x0b, 0x73, 0x6f, 0x6d, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x24, 0x0a, 0x0c, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x04, 0x80, 0xeb, 0x1f, 0x01, 0x52, 0x02, 0x69, 0x64, 0x22, 0x30, 0x0a, 0x18, 0x49,
	0x6e, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x4c,
	0x65, 0x67, 0x61, 0x63, 0x79, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x04, 0xd0, 0xeb, 0x1f, 0x01, 0x52, 0x02, 0x69, 0x64, 0x22, 0x86, 0x01,
	0x0a, 0x11, 0x4d, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x45,
	0x6e, 0x75, 0x6d, 0x12, 0x3b, 0x0a, 0x04, 0x65, 0x6e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x21, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x4d, 0x79, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x45, 0x6e, 0x75, 0x6d, 0x2e, 0x4d, 0x79,
	0x45, 0x6e, 0x75, 0x6d, 0x42, 0x04, 0x90, 0xec, 0x1f, 0x01, 0x52, 0x04, 0x65, 0x6e, 0x75, 0x6d,
	0x22, 0x34, 0x0a, 0x06, 0x4d, 0x79, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x59,
	0x5f, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4d, 0x59, 0x5f, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x46,
	0x49, 0x52, 0x53, 0x54, 0x10, 0x01, 0x42, 0x4c, 0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x61, 0x66, 0x65, 0x74, 0x79, 0x43, 0x75, 0x6c, 0x74, 0x75,
	0x72, 0x65, 0x2f, 0x73, 0x31, 0x32, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x67, 0x6f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_example_proto_rawDescOnce sync.Once
	file_example_proto_rawDescData = file_example_proto_rawDesc
)

func file_example_proto_rawDescGZIP() []byte {
	file_example_proto_rawDescOnce.Do(func() {
		file_example_proto_rawDescData = protoimpl.X.CompressGZIP(file_example_proto_rawDescData)
	})
	return file_example_proto_rawDescData
}

var file_example_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_example_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_example_proto_goTypes = []interface{}{
	(MyMessageWithEnum_MyEnum)(0),                           // 0: example.MyMessageWithEnum.MyEnum
	(*ExampleMessage)(nil),                                  // 1: example.ExampleMessage
	(*OuterMessageUsingNestedMessage)(nil),                  // 2: example.OuterMessageUsingNestedMessage
	(*InnerMessage)(nil),                                    // 3: example.InnerMessage
	(*InnerMessageWithLegacyId)(nil),                        // 4: example.InnerMessageWithLegacyId
	(*MyMessageWithEnum)(nil),                               // 5: example.MyMessageWithEnum
	(*ExampleMessage_NestedMessage)(nil),                    // 6: example.ExampleMessage.NestedMessage
	nil,                                                     // 7: example.ExampleMessage.NotSupportedEntry
	(*ExampleMessage_Contact)(nil),                          // 8: example.ExampleMessage.Contact
	(*ExampleMessage_NestedMessage_InnerNestedMessage)(nil), // 9: example.ExampleMessage.NestedMessage.InnerNestedMessage
	(*timestamppb.Timestamp)(nil),                           // 10: google.protobuf.Timestamp
}
var file_example_proto_depIdxs = []int32{
	3,  // 0: example.ExampleMessage.inner:type_name -> example.InnerMessage
	3,  // 1: example.ExampleMessage.msg_required:type_name -> example.InnerMessage
	4,  // 2: example.ExampleMessage.inner_legacy_id:type_name -> example.InnerMessageWithLegacyId
	6,  // 3: example.ExampleMessage.nested_message:type_name -> example.ExampleMessage.NestedMessage
	7,  // 4: example.ExampleMessage.not_supported:type_name -> example.ExampleMessage.NotSupportedEntry
	8,  // 5: example.ExampleMessage.contacts_with_length_constraint:type_name -> example.ExampleMessage.Contact
	8,  // 6: example.ExampleMessage.contacts_with_no_length_constraint:type_name -> example.ExampleMessage.Contact
	10, // 7: example.ExampleMessage.scheduled_for:type_name -> google.protobuf.Timestamp
	6,  // 8: example.OuterMessageUsingNestedMessage.some_message:type_name -> example.ExampleMessage.NestedMessage
	0,  // 9: example.MyMessageWithEnum.enum:type_name -> example.MyMessageWithEnum.MyEnum
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_example_proto_init() }
func file_example_proto_init() {
	if File_example_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_example_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExampleMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_example_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OuterMessageUsingNestedMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_example_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InnerMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_example_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InnerMessageWithLegacyId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_example_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MyMessageWithEnum); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_example_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExampleMessage_NestedMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_example_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExampleMessage_Contact); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_example_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExampleMessage_NestedMessage_InnerNestedMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_example_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ExampleMessage_Fax)(nil),
		(*ExampleMessage_Phone)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_example_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_example_proto_goTypes,
		DependencyIndexes: file_example_proto_depIdxs,
		EnumInfos:         file_example_proto_enumTypes,
		MessageInfos:      file_example_proto_msgTypes,
	}.Build()
	File_example_proto = out.File
	file_example_proto_rawDesc = nil
	file_example_proto_goTypes = nil
	file_example_proto_depIdxs = nil
}
