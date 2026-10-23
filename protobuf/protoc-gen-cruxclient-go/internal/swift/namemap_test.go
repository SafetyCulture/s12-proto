// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestNameInstructionWriterRejectsDescendingNumbers(t *testing.T) {
	w := newNameInstructionWriter()
	w.writeSame(2, "second")
	w.writeSame(1, "first")

	require.EqualError(t, w.err, "field and case numbers must be written in ascending order")
}

func TestNameMapOrderErrorNamesFileAndType(t *testing.T) {
	f := &file{name: "broken.proto"}
	g := &generator{opts: &options{}, namer: newNamer("", f)}
	mg := &messageGen{
		g:    g,
		desc: &message{fullName: "example.Message"},
		fieldsByNumber: []fieldGenerator{
			&messageFieldGen{desc: &field{name: "second", jsonName: "second", number: 2}},
			&messageFieldGen{desc: &field{name: "first", jsonName: "first", number: 1}},
		},
	}
	mg.generateProtoNameProviding(newPrinter())
	require.EqualError(t, g.fileError(f), "internal error of the Swift generator on broken.proto: example.Message: field and case numbers must be written in ascending order")

	g.err = nil
	eg := &enumGen{
		g:         g,
		desc:      &enum{fullName: "example.State"},
		aliasInfo: &aliasInfo{},
		mainSorted: []*enumValue{
			{proto: &descriptorpb.EnumValueDescriptorProto{}, name: "second", number: 2},
			{proto: &descriptorpb.EnumValueDescriptorProto{}, name: "first", number: 1},
		},
	}
	eg.generateProtoNameProviding(newPrinter())
	require.EqualError(t, g.fileError(f), "internal error of the Swift generator on broken.proto: example.State: field and case numbers must be written in ascending order")
}

func TestRenderReportsNameMapOrderError(t *testing.T) {
	defer func(original func() *nameInstructionWriter) { newNameWriter = original }(newNameWriter)
	newNameWriter = func() *nameInstructionWriter {
		w := newNameInstructionWriter()
		w.writeSame(2, "second")
		w.writeSame(1, "first")
		return w
	}
	file := &descriptorpb.FileDescriptorProto{
		Name:        proto.String("order.proto"),
		Package:     proto.String("order.v1"),
		Syntax:      proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("M")}},
	}
	opts, err := parseOptions("", nil)
	require.NoError(t, err)
	_, err = render([]*descriptorpb.FileDescriptorProto{file}, []string{"order.proto"}, opts)
	var internal *internalError
	require.ErrorAs(t, err, &internal)
	require.EqualError(t, err, "internal error of the Swift generator on order.proto: order.v1.M: field and case numbers must be written in ascending order")
}
