// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.
#include "api_generator.h"
#include <google/protobuf/compiler/code_generator.h>
#include <google/protobuf/compiler/plugin.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/io/printer.h>
#include <google/protobuf/io/zero_copy_stream.h>
#include <iostream>
#include <vector>
#include <map>
#include <memory>
#include "common.h"

using google::protobuf::Descriptor;
using google::protobuf::FileDescriptor;
using google::protobuf::MethodDescriptor;
using google::protobuf::ServiceDescriptor;
using google::protobuf::compiler::CodeGenerator;
using google::protobuf::compiler::GeneratorContext;
using google::protobuf::compiler::ParseGeneratorParameter;
using google::protobuf::compiler::PluginMain;
using google::protobuf::io::Printer;
using google::protobuf::io::ZeroCopyOutputStream;
using std::string;

namespace cruxclient_generator {

void APIGenerator::PrintPrologue(
  Printer *printer,
  const FileDescriptor *file) const {
  std::map<string, string> vars;
  vars["filename"] = file->name();
  printer->Print("// Generated by the CRUX Engine C++ plugin.\n");
  printer->Print("// If you make any local change, they will be lost.\n");
  printer->Print(vars, "// source: $filename$\n");
}

void APIGenerator::PrintNamespace(
  Printer *printer,
  const FileDescriptor *file,
  bool isEpilogue) const {
  std::map<string, string> vars;
  if (!file->package().empty()) {
    vars["namespace"] = DotsToColons(file->package());

    if (!isEpilogue) {
      printer->Print(vars, "namespace $namespace$ {\n");
    } else {
      printer->Print(vars, "}  // namespace $namespace$\n");
    }
    printer->Print(vars, "\n");
  }
}

void APIGenerator::PrintHeaderPrologue(
  Printer *printer,
  const FileDescriptor *file) const {
  std::map<string, string> vars;
  vars["filename_identifier"] = FilenameIdentifier(file->name());

  PrintPrologue(printer, file);
}

std::string APIGenerator::GetMethodSignature(
  const std::string& service_name,
  const std::string& method_name) const {
  return "k" + service_name + method_name;
}

void APIGenerator::PrintHeaderIncludes(
  Printer *printer,
  const FileDescriptor *file) const {
  std::map<string, string> vars;
  vars["filename_base"] = StripProto(file->name());
  auto package_name = file->package();
  ReplaceCharacters(&package_name, ".", '/');

  printer->Print("#pragma once\n\n");
  printer->Print("#include <vector>\n");
  printer->Print("#include <string>\n");
  printer->Print("#include <memory>\n\n");
  printer->Print(vars, "#include <google/protobuf/any.pb.h>\n");
  printer->Print(vars, "#include \"$filename_base$.grpc.pb.h\"\n\n");

  PrintNamespace(printer, file, false);
}

void APIGenerator::PrintHeaderInterface(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const {
  printer->Print("class ChannelProvider {\n");
  printer->Print(" public:\n");
  printer->Indent();
  printer->Print("virtual std::shared_ptr<grpc::Channel> ConnectionChannel() const = 0;\n");
  printer->Outdent();
  printer->Print("};\n\n");
}

void APIGenerator::PrintHeaderAPIs(
  Printer *printer,
  const FileDescriptor *file) const {
  std::map<string, string> vars;

  for (int service_index = 0; service_index < file->service_count();
       ++service_index) {
    const ServiceDescriptor *service = file->service(service_index);
    vars["service_name"] = service->name();
    vars["service_fullname"] = DotsToColons(service->full_name());

    printer->Print(vars, "namespace $service_name$ {\n");
    for (int method_index = 0; method_index < service->method_count();
        ++method_index) {
      const MethodDescriptor *method = service->method(method_index);
      vars["method_name"] = method->name();
      vars["request"] = ClassName(method->input_type(), true);
      vars["response"] = ClassName(method->output_type(), true);
      vars["api_name"] = method->name() + "API";

      if (method->client_streaming()) {
        // [RC]: Client Steaming not supported yet
        continue;
      }


      printer->Print(vars, "class $api_name$ {\n");
      printer->Print(" public:\n");
      printer->Indent();
      printer->Print(vars,
      "explicit $api_name$(const std::shared_ptr<ChannelProvider>& provider);\n");
      printer->Print("std::string Name() const;\n");
      printer->Print(vars, "grpc::Status Execute(\n");
      printer->Indent();
      printer->Print(vars, "grpc::ClientContext* context,\n"
        "const $request$& request,\n"
        "$response$* response) const;\n");
      printer->Outdent();
      printer->Outdent();
      printer->Print(" private:\n");
      printer->Indent();
      printer->Print(vars, "std::unique_ptr<$service_name$::StubInterface> mStub;\n");
      printer->Outdent();
      printer->Print("};\n\n");
    }
    printer->Print(vars, "}  // namespace $service_name$\n\n");
  }
}

void APIGenerator::PrintHeaderEpilogue(
  Printer *printer,
  const FileDescriptor *file) const {
  PrintNamespace(printer, file, true);
}

void APIGenerator::PrintSourcePrologue(
  Printer *printer,
  const FileDescriptor *file) const {
  PrintPrologue(printer, file);
  printer->Print("\n");
}

void APIGenerator::PrintSourceIncludes(
  Printer *printer,
  const FileDescriptor *file) const {
  std::map<string, string> vars;
  vars["filename_base"] = tokenize(StripProto(file->name()), "/").back();
  printer->Print(vars, "#include \"$filename_base$.crux.api.h\"\n\n");

  PrintNamespace(printer, file, false);
}

void APIGenerator::PrintSourceAPIs(
  Printer *printer,
  const FileDescriptor *file) const {
  std::map<string, string> vars;

  for (int service_index = 0; service_index < file->service_count();
       ++service_index) {
    const ServiceDescriptor *service = file->service(service_index);
    vars["service_name"] = service->name();
    vars["service_fullname"] = DotsToColons(service->full_name());

    printer->Print(vars, "namespace $service_name${\n");
    for (int method_index = 0; method_index < service->method_count();
         ++method_index) {
      const MethodDescriptor *method = service->method(method_index);
      vars["method_name"] = method->name();
      vars["request"] = ClassName(method->input_type(), true);
      vars["response"] = ClassName(method->output_type(), true);
      vars["api_name"] = method->name() + "API";

      if (method->client_streaming()) {
        // [RC]: Client Streaming not supported... yet
        continue;
      }

      printer->Print(vars,
      "$api_name$(const std::shared_ptr<ChannelProvider>& provider) {\n");
      printer->Indent();
      printer->Print(vars, "mStub = $service_fullname$::NewStub(provider->ConnectionChannel());\n");
      printer->Outdent();
      printer->Print("}\n\n");

      printer->Print(vars, "std::string $api_name$::Name() const {\n");
      printer->Indent();
      printer->Print(vars, "return \"$service_name$_$method_name$\";\n");
      printer->Outdent();
      printer->Print("}\n\n");

      printer->Print(vars, "grpc::Status $api_name$::Execute(\n");
      printer->Indent();
      printer->Print(vars, "grpc::ClientContext* context,\n"
        "const $request$& request,\n"
        "$response$* response) const {\n");
      printer->Print(vars, "return mStub->$method_name$(context, request, response);\n");
      printer->Outdent();
      printer->Print("}\n\n");
    }
    printer->Print(vars, "}  // namespace $service_name$\n\n");
  }
}

void APIGenerator::PrintSourceEpilogue(
  Printer *printer,
  const FileDescriptor *file) const {
  printer->Print("\n");
  PrintNamespace(printer, file, true);
}

void APIGenerator::GenerateAPIHeader(
  const google::protobuf::FileDescriptor *file,
  const std::string &parameter,
  google::protobuf::compiler::GeneratorContext *context) const {
  string file_name = StripProto(file->name());
  std::unique_ptr<ZeroCopyOutputStream> output(
      context->Open(file_name + ".crux.api.h"));
  Printer printer(output.get(), '$');
  PrintHeaderPrologue(&printer, file);
  PrintHeaderIncludes(&printer, file);
  PrintHeaderInterface(&printer, file);
  PrintHeaderAPIs(&printer, file);
  PrintHeaderEpilogue(&printer, file);
}

void APIGenerator::GenerateAPISource(
  const google::protobuf::FileDescriptor *file,
  const std::string &parameter,
  google::protobuf::compiler::GeneratorContext *context) const {
  string file_name = StripProto(file->name());
  std::unique_ptr<ZeroCopyOutputStream> output(
      context->Open(file_name + ".crux.api.cc"));
  Printer printer(output.get(), '$');
  PrintSourcePrologue(&printer, file);
  PrintSourceIncludes(&printer, file);
  PrintSourceAPIs(&printer, file);
  PrintSourceEpilogue(&printer, file);
}

void APIGenerator::Generate(
  const google::protobuf::FileDescriptor *file,
  const std::string &parameter,
  google::protobuf::compiler::GeneratorContext *context,
  std::string *error) const {
  GenerateAPIHeader(file, parameter, context);
  GenerateAPISource(file, parameter, context);
}
}  // namespace cruxclient_generator
