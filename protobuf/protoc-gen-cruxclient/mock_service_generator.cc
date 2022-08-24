// Copyright (c) 2022 SafetyCulture Pty Ltd. All Rights Reserved.
#include "mock_service_generator.h"
#include <google/protobuf/compiler/code_generator.h>
#include <google/protobuf/compiler/plugin.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/descriptor.pb.h>
#include <google/protobuf/io/printer.h>
#include <google/protobuf/io/zero_copy_stream.h>
#include <iostream>
#include <vector>
#include <map>
#include <memory>
#include <set>
#include "common.h"

using google::protobuf::Descriptor;
using google::protobuf::FileDescriptor;
using google::protobuf::MethodDescriptor;
using google::protobuf::ServiceDescriptor;
using google::protobuf::FileOptions;
using google::protobuf::compiler::CodeGenerator;
using google::protobuf::compiler::GeneratorContext;
using google::protobuf::compiler::ParseGeneratorParameter;
using google::protobuf::compiler::PluginMain;
using google::protobuf::io::Printer;
using google::protobuf::io::ZeroCopyOutputStream;

namespace cruxclient_generator {
void MockServiceGenerator::Generate(
  const google::protobuf::FileDescriptor *file, 
  const std::string &parameter,
  google::protobuf::compiler::GeneratorContext *context,
  std::string *error) const {
  GenerateHeader(file, parameter, context);
}

void MockServiceGenerator::GenerateHeader(
  const google::protobuf::FileDescriptor *file,
  const std::string &parameter,
  google::protobuf::compiler::GeneratorContext *context) const {
  std::string file_name = StripProto(file->name());
  std::unique_ptr<ZeroCopyOutputStream> output(
    context->Open(file_name + ".mock.h"));
  Printer printer(output.get(), '$');

  std::map<std::string, std::string> vars;
  vars["filename"] = file->name();
  vars["filename_base"] = StripProto(file->name());
  if (!file->package().empty()) {
    vars["namespace"] = DotsToColons(file->package());
  }
  printer.Print("// Generated by the CRUX Engine C++ plugin.\n");
  printer.Print("// If you make any local change, they will be lost.\n");
  printer.Print(vars, "// source: $filename$\n");
  
  printer.Print("#pragma once\n\n");
  printer.Print("#include <tuple>\n");
  printer.Print(vars, "#include \"$filename_base$.grpc.pb.h\"\n\n");
  printer.Print("#include \"crux_mock_support.h\"\n\n");

  printer.Print(vars, "namespace $namespace$ {\n");

  for (int service_index = 0; service_index < file->service_count(); ++service_index) {
    const ServiceDescriptor *service = file->service(service_index);
    vars["class_name"] = service->name() + "MockImpl";
    vars["service_fullname"] = DotsToColons(service->full_name());

    printer.Print(vars, "class $class_name$ final : public ::$service_fullname$::Service {\n");
    printer.Print(" public:\n");
    printer.Indent();
    printer.Print(vars, "$class_name$(const GRPCServiceCallback& callback): mCallback(callback) {};\n");
    for (int method_index = 0; method_index < service->method_count(); ++method_index) {
      const MethodDescriptor *method = service->method(method_index);
      vars["method_name"] = method->name();
      vars["method_fullname"] = method->full_name();
      vars["request"] = ClassName(method->input_type(), true);
      vars["response"] = ClassName(method->output_type(), true);

      if (method->client_streaming()) {
        // Client Steaming not supported yet
        continue;
      }

      if (method->server_streaming()) {
        printer.Print(vars, "::grpc::Status $method_name$(::grpc::ServerContext* context, const ::$request$* request, ::grpc::ServerWriter<::$response$>* writer) override {\n");
        printer.Indent();
        printer.Print(vars, "auto[status, bytes] = mCallback(\"$method_fullname$\", request->SerializeAsString(), context->client_metadata());\n");
        printer.Print("if (bytes.has_value()) {\n");
        printer.Indent();
        printer.Print(vars, "::$response$ response;\n");
        printer.Print("response->ParseFromString(*bytes);\n");
        printer.Print("writer->Write(response);\n");
        printer.Outdent();
        printer.Print("}\n");
        printer.Print("return status;\n");
        printer.Outdent();
        printer.Print("}\n");
      } else {
        printer.Print(vars, "::grpc::Status $method_name$(::grpc::ServerContext* context, const ::$request$* request, ::$response$* response) override {\n");
        printer.Indent();
        printer.Print(vars, "auto[status, bytes] = mCallback(\"$method_fullname$\", request->SerializeAsString(), context->client_metadata());\n");
        printer.Print("if (bytes.has_value()) {\n");
        printer.Indent();
        printer.Print("response->ParseFromString(*bytes);\n");
        printer.Outdent();
        printer.Print("}\n");
        printer.Print("return status;\n");
        printer.Outdent();
        printer.Print("}\n");
      }
      
    }
    printer.Outdent();
    printer.Print(" private:\n");
    printer.Indent();
    printer.Print(vars, "const GRPCServiceCallback mCallback;\n");
    printer.Outdent();
    printer.Print("};\n\n");
  }

  printer.Print(vars, "}  // namespace $namespace$\n");
}
}  // namespace cruxclient_generator