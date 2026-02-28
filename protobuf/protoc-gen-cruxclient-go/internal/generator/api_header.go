// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

// Package generator implements the C++ code generators for the
// protoc-gen-cruxclient-go plugin.
package generator

import (
	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/codegen"
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go/internal/common"
)

// GenerateAPIHeader generates the .crux.api.h file for a proto file descriptor.
// The output is functionally equivalent to the C++ plugin's api_generator.cc
// GenerateHeader function.
func GenerateAPIHeader(rw protoplugin.ResponseWriter, fileDesc protoreflect.FileDescriptor) {
	w := codegen.NewWriter()

	path := string(fileDesc.Path())
	pkg := string(fileDesc.Package())

	// Prologue: 3 comment lines + #pragma once (no blank line between them)
	WritePrologue(w, path)
	w.WL("#pragma once")

	// Includes: blank line before STL includes, blank line between STL group and proto group
	w.WL()
	w.WL("#include <vector>")
	w.WL("#include <string>")
	w.WL("#include <memory>")
	w.WL()
	w.WL("#include <google/protobuf/any.pb.h>")
	w.WL("#include \"%s.grpc.pb.h\"", common.StripProto(path))
	w.WL("#include \"crux_engine_client_support.h\"")

	// Package namespace open
	w.WL()
	w.WL("namespace %s {", common.DotsToColons(pkg))
	w.WL()

	// Per-service blocks
	services := fileDesc.Services()
	for i := 0; i < services.Len(); i++ {
		svc := services.Get(i)
		writeHeaderServiceBlock(w, svc)
	}

	// Package namespace close
	w.WL("}  // namespace %s", common.DotsToColons(pkg))
	w.WL()

	baseName := common.StripProto(path)
	rw.AddFile(baseName+".crux.api.h", w.String())
}

// writeHeaderServiceBlock writes the namespace block for a single service.
func writeHeaderServiceBlock(w *codegen.CodeWriter, svc protoreflect.ServiceDescriptor) {
	svcName := string(svc.Name())

	// Namespace open
	w.WL("namespace %sNS {", svcName)

	// kServiceName constant
	w.WL("const char kServiceName[] = \"%s\";", common.DotsToUnderscores(string(svc.FullName())))

	// Per-method class declarations
	methods := nonClientStreamingMethods(svc)
	for _, method := range methods {
		writeHeaderMethodClass(w, svc, method)
	}

	// Invoke() declaration
	w.WL("grpc::Status Invoke(const std::shared_ptr<crux::engine::ChannelProvider>& provider, grpc::ClientContext* context, const google::protobuf::Any& request_data, const std::string& method_name);")

	// Namespace close
	w.WL("}  // namespace %sNS", svcName)
	w.WL()
}

// writeHeaderMethodClass writes a single method's API class declaration.
func writeHeaderMethodClass(w *codegen.CodeWriter, svc protoreflect.ServiceDescriptor, method protoreflect.MethodDescriptor) {
	methodName := string(method.Name())
	reqType := common.ClassName(method.Input(), true)
	respType := common.ClassName(method.Output(), true)
	svcName := string(svc.Name())

	w.WL("class %sAPI {", methodName)
	w.WL(" public:")
	w.WL("  explicit %sAPI(const std::shared_ptr<crux::engine::ChannelProvider>& provider);", methodName)
	w.WL("  static std::string Name();")
	w.WL("  static std::string ServiceName();")
	w.WL("  static std::string MethodName();")

	// Unary Execute overload (multi-line signature)
	w.WL("  grpc::Status Execute(")
	w.WL("    grpc::ClientContext* context,")
	w.WL("    const %s& request,", reqType)
	w.WL("    %s* response) const;", respType)

	// Streaming Execute overload
	w.WL("  // server streaming")
	w.WL("  std::unique_ptr<grpc::ClientReaderInterface<%s>> Execute(", respType)
	w.WL("    grpc::ClientContext* context,")
	w.WL("    const %s& request) const;", reqType)

	w.WL(" private:")
	w.WL("  std::unique_ptr<%s::StubInterface> mStub;", svcName)
	w.WL("};")
	w.WL()
}

// nonClientStreamingMethods returns methods that are not client-streaming.
// Both client-only and bidirectional streaming methods are excluded.
func nonClientStreamingMethods(svc protoreflect.ServiceDescriptor) []protoreflect.MethodDescriptor {
	var methods []protoreflect.MethodDescriptor
	allMethods := svc.Methods()
	for i := 0; i < allMethods.Len(); i++ {
		method := allMethods.Get(i)
		if method.IsStreamingClient() {
			continue
		}
		methods = append(methods, method)
	}
	return methods
}
