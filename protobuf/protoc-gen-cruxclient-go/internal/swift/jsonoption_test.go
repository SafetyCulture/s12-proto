// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// testOptionNumber is a number in the extension range of the enum value options. The number of the real option is not.
const testOptionNumber = 1998

func jsonOptionWire(value string) []byte {
	body := protowire.AppendTag(nil, 1, protowire.BytesType)
	body = protowire.AppendString(body, value)
	raw := protowire.AppendTag(nil, testOptionNumber, protowire.BytesType)
	return protowire.AppendBytes(raw, body)
}

func optionsFromWire(t *testing.T, raw []byte, resolver proto.UnmarshalOptions) *descriptorpb.EnumValueOptions {
	t.Helper()
	options := &descriptorpb.EnumValueOptions{}
	require.NoError(t, resolver.Unmarshal(raw, options))
	return options
}

// TestJSONOption checks that the option is found in the unknown fields and, when a program
// registers the extension, through the registered type, with the same result.
func TestJSONOption(t *testing.T) {
	_, ok := (&enumValue{proto: &descriptorpb.EnumValueDescriptorProto{}}).jsonOption()
	require.False(t, ok)

	empty := &protoregistry.Types{}
	raw := jsonOptionWire("custom_name")
	unknown := optionsFromWire(t, raw, proto.UnmarshalOptions{Resolver: empty})
	require.NotEmpty(t, unknown.ProtoReflect().GetUnknown())
	name, ok := readJSONOption(unknown, testOptionNumber, empty)
	require.True(t, ok)
	require.Equal(t, "custom_name", name)

	name, ok = readJSONOption(optionsFromWire(t, jsonOptionWire(""), proto.UnmarshalOptions{Resolver: empty}), testOptionNumber, empty)
	require.True(t, ok)
	require.Empty(t, name)

	repeated := append(jsonOptionWire("first"), jsonOptionWire("last")...)
	name, ok = readJSONOption(optionsFromWire(t, repeated, proto.UnmarshalOptions{Resolver: empty}), testOptionNumber, empty)
	require.True(t, ok)
	require.Equal(t, "last", name)

	_, ok = readJSONOption(&descriptorpb.EnumValueOptions{Deprecated: proto.Bool(true)}, testOptionNumber, empty)
	require.False(t, ok)

	extension := &descriptorpb.FileDescriptorProto{
		Name:       proto.String("pb/enumvalue/json.proto"),
		Package:    proto.String("pb.enumvalue"),
		Dependency: []string{"google/protobuf/descriptor.proto"},
		MessageType: []*descriptorpb.DescriptorProto{{
			Name: proto.String("JsonEnumValueOptions"),
			Field: []*descriptorpb.FieldDescriptorProto{{
				Name: proto.String("string"), Number: proto.Int32(1), JsonName: proto.String("string"),
				Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:  descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
			}},
		}},
		Extension: []*descriptorpb.FieldDescriptorProto{{
			Name: proto.String("json"), Number: proto.Int32(testOptionNumber), JsonName: proto.String("json"),
			Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
			TypeName: proto.String(".pb.enumvalue.JsonEnumValueOptions"),
			Extendee: proto.String(".google.protobuf.EnumValueOptions"),
		}},
	}
	fd, err := protodesc.NewFile(extension, protoregistry.GlobalFiles)
	require.NoError(t, err)
	registered := &protoregistry.Types{}
	require.NoError(t, registered.RegisterExtension(dynamicpb.NewExtensionType(fd.Extensions().Get(0))))

	known := optionsFromWire(t, raw, proto.UnmarshalOptions{Resolver: registered})
	require.Empty(t, known.ProtoReflect().GetUnknown())
	name, ok = readJSONOption(known, testOptionNumber, registered)
	require.True(t, ok)
	require.Equal(t, "custom_name", name)

	_, ok = readJSONOption(&descriptorpb.EnumValueOptions{Deprecated: proto.Bool(true)}, testOptionNumber, registered)
	require.False(t, ok)

	// Options that were read before the type was registered still have the option in the unknown fields.
	name, ok = readJSONOption(unknown, testOptionNumber, registered)
	require.True(t, ok)
	require.Equal(t, "custom_name", name)
}
