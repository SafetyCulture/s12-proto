// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package generator

import (
	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/codegen"
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/common"
)

// GenerateAPISource generates the .crux.api.cc file for a proto file descriptor.
// The output is functionally equivalent to the C++ plugin's api_generator.cc
// GenerateSource function.
func GenerateAPISource(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor) {
	w := codegen.NewWriter()

	path := string(fileDesc.Path())
	pkg := string(fileDesc.Package())

	// Prologue: 3 comment lines + blank line
	WritePrologue(w, path)
	w.WL()

	// Include: basename only (not full path)
	baseParts := common.Tokenize(common.StripProto(path), "/")
	basename := baseParts[len(baseParts)-1]
	w.WL("#include \"%s.crux.api.h\"", basename)

	// Package namespace open
	w.WL()
	w.WL("namespace %s {", common.DotsToColons(pkg))
	w.WL()

	// Per-service blocks
	services := fileDesc.Services()
	for i := 0; i < services.Len(); i++ {
		svc := services.Get(i)
		writeSourceServiceBlock(w, svc)
	}

	// Source epilogue: the leading \n combines with the service NS block's trailing \n\n
	// to produce the double blank line seen in the reference files.
	w.WL()
	w.WL("}  // namespace %s", common.DotsToColons(pkg))
	w.WL()

	baseName := common.StripProto(path)
	rw.AddFile(baseName+".crux.api.cc", w.String())
}

// writeSourceServiceBlock writes the full namespace block for a single service in the .cc file.
func writeSourceServiceBlock(w *codegen.CodeWriter, svc protoreflect.ServiceDescriptor) {
	svcName := string(svc.Name())

	// Namespace open
	w.WL("namespace %sNS {", svcName)

	// Per-method function implementations
	methods := nonClientStreamingMethods(svc)
	for _, method := range methods {
		writeSourceMethodFunctions(w, svc, method)
	}

	// Invoke() routing function
	writeSourceInvokeFunction(w, svc, methods)

	// Namespace close
	w.WL()
	w.WL("}  // namespace %sNS", svcName)
	w.WL()
}

// writeSourceMethodFunctions writes all 7 function implementations for a single method.
func writeSourceMethodFunctions(w *codegen.CodeWriter, svc protoreflect.ServiceDescriptor, method protoreflect.MethodDescriptor) {
	methodName := string(method.Name())
	reqType := common.ClassName(method.Input(), true)
	respType := common.ClassName(method.Output(), true)
	path := "/" + string(svc.FullName()) + "/" + methodName
	isServerStreaming := method.IsStreamingServer()

	// a. Constructor
	w.WL("%sAPI::%sAPI(const std::shared_ptr<crux::engine::RpcTransportProvider>& provider) : mProvider(provider) {}", methodName, methodName)
	w.WL()

	// b. Name()
	w.WL("std::string %sAPI::Name() {", methodName)
	w.WL("  return \"%s_%s\";", common.DotsToUnderscores(string(svc.FullName())), methodName)
	w.WL("}")
	w.WL()

	// c. ServiceName()
	w.WL("std::string %sAPI::ServiceName() {", methodName)
	w.WL("  return kServiceName;")
	w.WL("}")
	w.WL()

	// d. MethodName()
	w.WL("std::string %sAPI::MethodName() {", methodName)
	w.WL("  return \"%s\";", methodName)
	w.WL("}")
	w.WL()

	// e. Path(): the gRPC method path, which is also the gRPC-Web URL path
	w.WL("std::string %sAPI::Path() {", methodName)
	w.WL("  return \"%s\";", path)
	w.WL("}")
	w.WL()

	// f. Execute(unary): if server-streaming -> UNIMPLEMENTED; if unary -> call through the transport
	w.WL("crux::engine::RpcStatus %sAPI::Execute(", methodName)
	w.WL("  crux::engine::RpcContext* context,")
	w.WL("  const %s& request,", reqType)
	w.WL("  %s* response) const {", respType)
	if isServerStreaming {
		w.WL("  return crux::engine::RpcStatus(crux::engine::RpcStatusCode::UNIMPLEMENTED, \"Please call streaming method instead\");")
	} else {
		w.WL("  return crux::engine::UnaryCall(*mProvider, context, Path(), request, response);")
	}
	w.WL("}")
	w.WL()

	// g. Execute(streaming): if server-streaming -> call through the transport; if unary -> UNIMPLEMENTED reader
	w.WL("std::unique_ptr<crux::engine::RpcReader<%s>> %sAPI::Execute(", respType, methodName)
	w.WL("  crux::engine::RpcContext* context,")
	w.WL("  const %s& request) const {", reqType)
	if isServerStreaming {
		w.WL("  return crux::engine::ServerStreamingCall<%s>(*mProvider, context, Path(), request);", respType)
	} else {
		w.WL("  return crux::engine::FailedRpcReader<%s>(crux::engine::RpcStatus(crux::engine::RpcStatusCode::UNIMPLEMENTED, \"Please call non-streaming method instead\"));", respType)
	}
	w.WL("}")
	w.WL()
}

// writeSourceInvokeFunction writes the Invoke() routing function for a service.
func writeSourceInvokeFunction(w *codegen.CodeWriter, svc protoreflect.ServiceDescriptor, methods []protoreflect.MethodDescriptor) {
	w.WL("crux::engine::RpcStatus Invoke(const std::shared_ptr<crux::engine::RpcTransportProvider>& provider, crux::engine::RpcContext* context, const google::protobuf::Any& request_data, const std::string& method_name) {")

	for _, method := range methods {
		methodName := string(method.Name())
		reqType := common.ClassName(method.Input(), true)
		respType := common.ClassName(method.Output(), true)

		w.WL("  if (method_name == \"%s\") {", methodName)
		w.WL("    %s request;", reqType)
		w.WL("    if (!request_data.UnpackTo(&request)) {")
		w.WL("      return crux::engine::RpcStatus(crux::engine::RpcStatusCode::DATA_LOSS, \"Unable to unpack the request data\");")
		w.WL("    }")
		w.WL("    %sAPI api = %sAPI(provider);", methodName, methodName)
		w.WL("    %s response;", respType)
		w.WL("    return api.Execute(context, request, &response);")
		w.WL("  }")
		w.WL()
	}

	w.WL("  return crux::engine::RpcStatus(crux::engine::RpcStatusCode::DATA_LOSS, \"Invalid method name\");")
	w.WL("}")
}
