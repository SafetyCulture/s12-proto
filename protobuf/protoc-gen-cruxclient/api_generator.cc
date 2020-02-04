// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.
#include "api_generator.h"
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
  printer->Print(vars, "#include \"$filename_base$.grpc.pb.h\"\n");
  printer->Print(vars, "#include \"crux_engine_client_support.h\"\n\n");

  PrintNamespace(printer, file, false);
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

    printer->Print(vars, "namespace $service_name$NS {\n");

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
      "explicit $api_name$(const std::shared_ptr<crux::engine::ChannelProvider>& provider);\n");
      printer->Print("static std::string Name();\n");
      printer->Print("static std::string ServiceName();\n");
      printer->Print("static std::string MethodName();\n");
      printer->Print(vars, "grpc::Status Execute(\n");
      printer->Indent();
      printer->Print(vars, "grpc::ClientContext* context,\n"
        "const $request$& request,\n"
        "$response$* response) const;\n");
      printer->Outdent();
      printer->Print("// server streaming\n");
      printer->Print(vars, "std::unique_ptr<grpc::ClientReaderInterface<$response$>> Execute(\n");
      printer->Indent();
      printer->Print(vars, "grpc::ClientContext* context,\n"
        "const $request$& request) const;\n");
      printer->Outdent();
      printer->Outdent();
      printer->Print(" private:\n");
      printer->Indent();
      printer->Print(vars, "std::unique_ptr<$service_name$::StubInterface> mStub;\n");
      printer->Outdent();
      printer->Print("};\n\n");
    }

    printer->Print("template<typename RESPONSE>\n");
    printer->Print("grpc::Status Invoke("
    "const std::shared_ptr<crux::engine::ChannelProvider>& provider, "
    "grpc::ClientContext* context, "
    "const google::protobuf::Any& request_data, "
    "const std::string& method_name) {\n");
    printer->Indent();
    for (int method_index = 0; method_index < service->method_count();
        ++method_index) {
      const MethodDescriptor *method = service->method(method_index);
      vars["method_name"] = method->name();
      vars["request"] = ClassName(method->input_type(), true);
      vars["response"] = ClassName(method->output_type(), true);
      vars["api_name"] = method->name() + "API";

      printer->Print(vars, "if (method_name == \"$method_name$\") {\n");
      printer->Indent();
      printer->Print(vars, "$request$ request;\n");
      printer->Print("if (!request_data.UnpackTo(&request)) {\n");
      printer->Indent();
      printer->Print("return grpc::Status(grpc::StatusCode::DATA_LOSS, \"Unable to unpack the request data\");\n");
      printer->Outdent();
      printer->Print("}\n");
      printer->Print(vars, "$api_name$ api = $api_name$(provider);\n");
      printer->Print(vars, "$response$ response;\n");
      printer->Print("return api.Execute(context, request, &response);\n");
      printer->Outdent();
      printer->Print("}\n\n");
    }
    printer->Print("return grpc::Status(grpc::StatusCode::DATA_LOSS, \"Invalid method name\");\n");
    printer->Outdent();
    printer->Print("}\n\n");

    printer->Print(vars, "}  // namespace $service_name$NS\n\n");
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
    vars["service_fullname_underscore"] = DotsToUnderscores(service->full_name());

    printer->Print(vars, "namespace $service_name$NS {\n");
    printer->Print("template <typename R>\n");
    printer->Print("class UnimplementedClientReader final: public grpc::ClientReaderInterface<R> {\n");
    printer->Print(" public:\n");
    printer->Indent();
    printer->Print("void WaitForInitialMetadata() override {}\n");
    printer->Print("bool NextMessageSize(uint32_t* sz) override { return true; }\n");
    printer->Print("bool Read(R* msg) override { return false; }\n");
    printer->Print("grpc::Status Finish() override { return grpc::Status(grpc::StatusCode::UNIMPLEMENTED, \"Please call non-streaming method instead\"); }\n");
    printer->Outdent();
    printer->Print("};\n\n");
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
      "$api_name$::$api_name$(const std::shared_ptr<crux::engine::ChannelProvider>& provider) {\n");
      printer->Indent();
      printer->Print(vars, "mStub = $service_fullname$::NewStub(provider->ConnectionChannel());\n");
      printer->Outdent();
      printer->Print("}\n\n");

      printer->Print(vars, "std::string $api_name$::Name() {\n");
      printer->Indent();
      printer->Print(vars, "return \"$service_fullname_underscore$_$method_name$\";\n");
      printer->Outdent();
      printer->Print("}\n\n");

      printer->Print(vars, "std::string $api_name$::ServiceName() {\n");
      printer->Indent();
      printer->Print(vars, "return \"$service_fullname_underscore$\";\n");
      printer->Outdent();
      printer->Print("}\n\n");

      printer->Print(vars, "std::string $api_name$::MethodName() {\n");
      printer->Indent();
      printer->Print(vars, "return \"$method_name$\";\n");
      printer->Outdent();
      printer->Print("}\n\n");

      printer->Print(vars, "grpc::Status $api_name$::Execute(\n");
      printer->Indent();
      printer->Print(vars, "grpc::ClientContext* context,\n"
        "const $request$& request,\n"
        "$response$* response) const {\n");
      if (method->server_streaming()) {
        printer->Print(vars, "return grpc::Status(grpc::StatusCode::UNIMPLEMENTED, \"Please call streaming method instead\");\n");
      } else {
        printer->Print(vars, "return mStub->$method_name$(context, request, response);\n");
      }
      printer->Outdent();
      printer->Print("}\n\n");

      printer->Print(vars, "std::unique_ptr<grpc::ClientReaderInterface<$response$>> $api_name$::Execute(\n");
      printer->Indent();
      printer->Print(vars, "grpc::ClientContext* context,\n"
        "const $request$& request) const {\n");
      if (method->server_streaming()) {
        printer->Print(vars, "return mStub->$method_name$(context, request);\n");
      } else {
        printer->Print(vars, "return std::make_unique<UnimplementedClientReader<$response$>>();\n");
      }
      printer->Outdent();
      printer->Print("}\n\n");
    }
    printer->Print(vars, "}  // namespace $service_name$NS\n\n");
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

void APIGenerator::GenerateDjinniYAML(
  const google::protobuf::FileDescriptor *file,
  google::protobuf::compiler::GeneratorContext *context) const {
  string file_name = StripProto(file->name());
  std::vector<std::string> paths = tokenize(StripProto(file_name), "/");
  paths.pop_back();
  std::string dir = join(paths, "/");

  std::unique_ptr<ZeroCopyOutputStream> output(
      context->Open(file_name + ".djinni.yaml"));
  Printer printer(output.get(), '$');
  printer.Print("# This file is generated by s12-proto api generator. Please DO NOT modify.\n");

  for (int dep_index = 0; dep_index < file->dependency_count(); ++dep_index) {
    auto imported_file = file->dependency(dep_index);
    std::vector<std::string> dep_paths = tokenize(StripProto(imported_file->name()), "/");
    dep_paths.pop_back();
    std::string dep_dir = join(dep_paths, "/");
    if (dep_dir == dir) {
      PrintDjinniYAML(&printer, imported_file, file_name);
    }
  }

  PrintDjinniYAML(&printer, file, file_name);
}

void APIGenerator::GenerateDjinniObjcSupport(
  const google::protobuf::FileDescriptor *file,
  google::protobuf::compiler::GeneratorContext *context) const {
  string file_name = StripProto(file->name());
  std::vector<std::string> paths = tokenize(StripProto(file_name), "/");
  paths.pop_back();
  std::string dir = join(paths, "/");

  std::unique_ptr<ZeroCopyOutputStream> output(
      context->Open(file_name + ".djinni.objc.h"));
  Printer printer(output.get(), '$');

  printer.Print("// This file is generated by s12-proto api generator. Please DO NOT modify.\n");
  printer.Print("#pragma once\n");
  printer.Print("#include <string>\n");

  for (int dep_index = 0; dep_index < file->dependency_count(); ++dep_index) {
    auto imported_file = file->dependency(dep_index);
    std::vector<std::string> dep_paths = tokenize(StripProto(imported_file->name()), "/");
    dep_paths.pop_back();
    std::string dep_dir = join(dep_paths, "/");
    if (dep_dir == dir) {
      PrintDjinniObjcSupport(&printer, imported_file);
    }
  }
  PrintDjinniObjcSupport(&printer, file);
}

void APIGenerator::PrintDjinniYAML(
  google::protobuf::io::Printer *printer,
  const google::protobuf::FileDescriptor *file,
  const std::string& main_file_name) const {
  string file_name = StripProto(file->name());
  std::vector<std::string> paths = tokenize(StripProto(file_name), "/");
  paths.pop_back();
  std::string dir = join(paths, "/");

  std::map<string, string> vars;
  vars["dir"] = dir;
  for (int message_index = 0; message_index < file->message_type_count();
       ++message_index) {
    const Descriptor *message = file->message_type(message_index);
    vars["message_name"] = message->name();
    vars["cpp_type_name"] = DotsToColons(message->full_name());
    vars["objc_header"] = DotsToSlashs(message->full_name());
    vars["file_name"] = StripProto(file->name());
    vars["objc_file_name"] = ToCamelCase(tokenize(StripProto(file->name()), "/").back());
    vars["main_file_name"] = main_file_name;

    const auto options = file->options();
    vars["java_package"] = options.java_package();
    vars["objc_class_prefix"] = options.objc_class_prefix();

    printer->Print("---\n");
    printer->Print(vars, "name: pb_$message_name$\n");
    printer->Print("typedef: 'record deriving(eq, ord, parcelable)'\n");
    printer->Print("params: []\n");
    printer->Print("prefix: 'pb'\n");

    printer->Print("cpp:\n");
    printer->Indent();
    printer->Print(vars, "typename: '$cpp_type_name$'\n");
    printer->Print(vars, "header: '\"$file_name$.pb.h\"'\n");
    printer->Print("byValue: false\n");
    printer->Outdent();

    printer->Print("objc:\n");
    printer->Indent();
    printer->Print(vars, "typename: '$objc_class_prefix$$message_name$'\n");
    printer->Print(vars, "header: '\"$dir$/$objc_file_name$.pbobjc.h\"'\n");
    printer->Print(vars, "boxed: '$objc_class_prefix$$message_name$'\n");
    printer->Print("pointer: true\n");
    printer->Print("hash: '%s.hash()'\n");
    printer->Outdent();

    printer->Print("objcpp:\n");
    printer->Indent();
    printer->Print(vars, "translator: 'djinni::$cpp_type_name$::Translator'\n");
    printer->Print(vars, "header: '\"$main_file_name$.djinni.objc.h\"'\n");
    printer->Outdent();

    printer->Print("java:\n");
    printer->Indent();
    printer->Print(vars, "typename: '$java_package$.$message_name$'\n");
    printer->Print(vars, "boxed: '$java_package$.$message_name$'\n");
    printer->Print("reference: true\n");
    printer->Print("generic: false\n");
    printer->Print("hash: '%s.hashCode()'\n");
    printer->Outdent();

    printer->Print("jni:\n");
    printer->Indent();
    printer->Print(vars, "translator: 'djinni::$cpp_type_name$::Translator'\n");
    printer->Print(vars, "header: '\"$main_file_name$.djinni.jni.h\"'\n");
    printer->Print(vars, "typename: '$message_name$'\n");
    printer->Print(vars, "typeSignature: 'L$java_package$.$message_name$;'\n");
    printer->Outdent();

    printer->Print("\n");
  }
}

void APIGenerator::PrintDjinniObjcSupport(
  google::protobuf::io::Printer *printer,
  const google::protobuf::FileDescriptor *file) const {
  string file_name = StripProto(file->name());
  std::vector<std::string> paths = tokenize(StripProto(file_name), "/");
  paths.pop_back();
  std::string dir = join(paths, "/");

  std::map<string, string> vars;
  vars["dir"] = dir;
  for (int message_index = 0; message_index < file->message_type_count();
       ++message_index) {
    const Descriptor *message = file->message_type(message_index);
    vars["message_name"] = message->name();
    vars["objc_file_name"] = ToCamelCase(tokenize(StripProto(file->name()), "/").back());
    vars["cpp_type_name"] = DotsToColons(message->full_name());
    vars["file_name"] = StripProto(file->name());
    vars["cpp_header"] = DotsToSlashs(ToLower(message->full_name()));
    vars["objc_header"] = DotsToSlashs(message->full_name());
    const auto options = file->options();
    vars["objc_class_prefix"] = options.objc_class_prefix();

    printer->Print(vars, "#include \"$file_name$.pb.h\"\n");
    printer->Print(vars, "#import \"$dir$/$objc_file_name$.pbobjc.h\"\n\n");

    printer->Print(vars, "namespace djinni::$cpp_type_name$ {\n");
    printer->Print(vars, "struct Translator {\n");
    printer->Indent();
    printer->Print(vars, "using CppType = ::$cpp_type_name$;\n");
    printer->Print(vars, "using ObjcType = $objc_class_prefix$$message_name$*;\n");
    printer->Print(vars, "using Boxed = Translator;\n\n");

    printer->Print("static CppType toCpp(ObjcType message) {\n");
    printer->Indent();
    printer->Print("assert(message);\n");
    printer->Print("NSData * data = [message data];\n");
    printer->Print("const void *bytes = [data bytes];\n");
    printer->Print("int byte_len = (int)[data length];\n");
    printer->Print("CppType cpp_message;\n");
    printer->Print("cpp_message.ParseFromArray(bytes, byte_len);\n");
    printer->Print("return cpp_message;\n");
    printer->Outdent();
    printer->Print("}\n\n");

    printer->Print("static ObjcType fromCpp(const CppType& message) {\n");
    printer->Indent();
    printer->Print("size_t byte_size = message.ByteSizeLong();\n");
    printer->Print("void *bytes = malloc(byte_size);\n");
    printer->Print("message.SerializeToArray(bytes, (int)byte_size);\n");
    printer->Print("NSData *data = [NSData dataWithBytes: bytes length: (int)byte_size];\n");
    printer->Print("NSError *error;\n");
    printer->Print(vars, "return [$objc_class_prefix$$message_name$ parseFromData:data error:&error];\n");
    printer->Outdent();
    printer->Print("}\n");

    printer->Outdent();
    printer->Print(vars, "};\n");

    printer->Print(vars, "}  //namespace djinni::$cpp_type_name$\n\n");
  }
}

void APIGenerator::GenerateDjinniJavaSupport(
  const google::protobuf::FileDescriptor *file,
  google::protobuf::compiler::GeneratorContext *context) const {

}

void APIGenerator::Generate(
  const google::protobuf::FileDescriptor *file,
  const std::string &parameter,
  google::protobuf::compiler::GeneratorContext *context,
  std::string *error) const {
  GenerateAPIHeader(file, parameter, context);
  GenerateAPISource(file, parameter, context);
}

void APIGenerator::GenerateDjinniSupport(
  const google::protobuf::FileDescriptor *file,
  const std::string &parameter,
  google::protobuf::compiler::GeneratorContext *context,
  std::string *error) const {
  GenerateDjinniYAML(file, context);
  GenerateDjinniObjcSupport(file, context);
  GenerateDjinniJavaSupport(file, context);
}
}  // namespace cruxclient_generator
