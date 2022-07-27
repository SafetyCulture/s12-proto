// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: s12/flags/permissions/permissions.proto

package permissions

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_s12_flags_permissions_permissions_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         20100,
		Name:          "s12.flags.permissions.required_flags",
		Tag:           "bytes,20100,rep,name=required_flags",
		Filename:      "s12/flags/permissions/permissions.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// A list for required permissions for the given method.
	// For example:
	//     rpc GetUser(Request) returns (Response) {
	//       option (s12.flags.permissions.required_flags) = "read:users"
	//       option (google.api.http).get = "/users/v1/profiles:GetUser"
	//     }
	//
	// repeated string required_flags = 20100;
	E_RequiredFlags = &file_s12_flags_permissions_permissions_proto_extTypes[0]
)

var File_s12_flags_permissions_permissions_proto protoreflect.FileDescriptor

var file_s12_flags_permissions_permissions_proto_rawDesc = []byte{
	0x0a, 0x27, 0x73, 0x31, 0x32, 0x2f, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x2f, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x73, 0x31, 0x32, 0x2e, 0x66,
	0x6c, 0x61, 0x67, 0x73, 0x2e, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x3a, 0x47, 0x0a, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x5f, 0x66,
	0x6c, 0x61, 0x67, 0x73, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x84, 0x9d, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x42, 0x3a, 0x5a, 0x38, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x61, 0x66, 0x65, 0x74, 0x79,
	0x43, 0x75, 0x6c, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x73, 0x31, 0x32, 0x2d, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x31, 0x32, 0x2f, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x2f, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
}

var file_s12_flags_permissions_permissions_proto_goTypes = []interface{}{
	(*descriptorpb.MethodOptions)(nil), // 0: google.protobuf.MethodOptions
}
var file_s12_flags_permissions_permissions_proto_depIdxs = []int32{
	0, // 0: s12.flags.permissions.required_flags:extendee -> google.protobuf.MethodOptions
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_s12_flags_permissions_permissions_proto_init() }
func file_s12_flags_permissions_permissions_proto_init() {
	if File_s12_flags_permissions_permissions_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_s12_flags_permissions_permissions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_s12_flags_permissions_permissions_proto_goTypes,
		DependencyIndexes: file_s12_flags_permissions_permissions_proto_depIdxs,
		ExtensionInfos:    file_s12_flags_permissions_permissions_proto_extTypes,
	}.Build()
	File_s12_flags_permissions_permissions_proto = out.File
	file_s12_flags_permissions_permissions_proto_rawDesc = nil
	file_s12_flags_permissions_permissions_proto_goTypes = nil
	file_s12_flags_permissions_permissions_proto_depIdxs = nil
}
