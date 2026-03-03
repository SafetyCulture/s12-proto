// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import (
	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/codegen"
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/common"
)

// GenerateMock generates the .mock.h file for a proto file descriptor.
// The output is functionally equivalent to the C++ plugin's mock_generator.cc
// GenerateMock function. It produces a MockImpl class for each service that
// uses the GRPCServiceCallback pattern.
func GenerateMock(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor) {
	w := codegen.NewWriter()

	path := string(fileDesc.Path())
	pkg := string(fileDesc.Package())
	baseName := common.StripProto(path)

	// Prologue
	WritePrologue(w, path)
	w.WL("#pragma once")
	w.WL()
	w.WL("#include <tuple>")
	w.WL("#include \"%s.grpc.pb.h\"", baseName)
	w.WL()
	w.WL("#include \"crux_mock_support.h\"")

	// Package namespace open
	w.WL()
	w.WL("namespace %s {", common.DotsToColons(pkg))

	// Per-service MockImpl classes
	services := fileDesc.Services()
	for i := 0; i < services.Len(); i++ {
		svc := services.Get(i)
		writeMockServiceClass(w, svc)
	}

	// Package namespace close
	w.WL("}  // namespace %s", common.DotsToColons(pkg))

	rw.AddFile(baseName+".mock.h", w.String())
}

// writeMockServiceClass writes a single MockImpl class for a service.
func writeMockServiceClass(w *codegen.CodeWriter, svc protoreflect.ServiceDescriptor) {
	className := string(svc.Name()) + "MockImpl"
	base := "::" + common.DotsToColons(string(svc.FullName())) + "::Service"

	// Class declaration
	w.WL("class %s final : public %s {", className, base)
	w.WL(" public:")

	// Constructor
	w.WL("  %s(const GRPCServiceCallback& callback): mCallback(callback) {};", className)

	// Method overrides
	methods := nonClientStreamingMethods(svc)
	for _, method := range methods {
		writeMockMethodOverride(w, method)
	}

	// Private section
	w.WL(" private:")
	w.WL("  const GRPCServiceCallback mCallback;")
	w.WL("};")
	w.WL()
}

// writeMockMethodOverride writes a single method override in the MockImpl class.
func writeMockMethodOverride(w *codegen.CodeWriter, method protoreflect.MethodDescriptor) {
	methodName := string(method.Name())
	reqType := common.ClassName(method.Input(), true)
	respType := common.ClassName(method.Output(), true)
	// Callback identifier uses proto dots (method.FullName()), NOT C++ colons
	callbackID := string(method.FullName())

	if method.IsStreamingServer() {
		// Server-streaming override
		w.WL("  ::grpc::Status %s(::grpc::ServerContext* context, const ::%s* request, ::grpc::ServerWriter<::%s>* writer) override {", methodName, reqType, respType)
		w.WL("    auto[status, bytes_list] = mCallback(\"%s\", request->SerializeAsString(), context->client_metadata());", callbackID)
		w.WL("    for (auto& bytes : bytes_list) {")
		w.WL("      ::%s response;", respType)
		w.WL("      response.ParseFromString(bytes);")
		w.WL("      writer->Write(response);")
		w.WL("    }")
		w.WL("    return status;")
		w.WL("  }")
	} else {
		// Unary override
		w.WL("  ::grpc::Status %s(::grpc::ServerContext* context, const ::%s* request, ::%s* response) override {", methodName, reqType, respType)
		w.WL("    auto[status, bytes_list] = mCallback(\"%s\", request->SerializeAsString(), context->client_metadata());", callbackID)
		w.WL("    if (!bytes_list.empty()) {")
		w.WL("      response->ParseFromString(bytes_list[0]);")
		w.WL("    }")
		w.WL("    return status;")
		w.WL("  }")
	}
}
