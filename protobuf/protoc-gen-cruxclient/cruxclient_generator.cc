// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

#include <google/protobuf/compiler/code_generator.h>
#include <google/protobuf/compiler/plugin.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/io/printer.h>
#include <google/protobuf/io/zero_copy_stream.h>
#include <iostream>
#include <string>
#include "common.h"
#include "api_generator.h"
#include "mock_service_generator.h"

using google::protobuf::Descriptor;
using google::protobuf::FileDescriptor;
using google::protobuf::MethodDescriptor;
using google::protobuf::ServiceDescriptor;
using google::protobuf::compiler::CodeGenerator;
using google::protobuf::compiler::GeneratorContext;
using google::protobuf::compiler::ParseGeneratorParameter;
using google::protobuf::compiler::PluginMain;
using std::string;

namespace cruxclient_generator {

class Generator : public CodeGenerator {
 public:
  Generator() {}
  ~Generator() override {}

  bool Generate(
    const FileDescriptor *file,
    const string &parameter,
    GeneratorContext *context,
    string *error) const override {
    APIGenerator api_generator;
    if (file->service_count() == 0) {
      // No services, generate djinni support
      api_generator.GenerateDjinniSupport(file, parameter, context, error);
      return true;
    }

    api_generator.Generate(file, parameter, context, error);
    api_generator.GenerateDjinniSupport(file, parameter, context, error);

    MockServiceGenerator mock_service_generator;
    mock_service_generator.Generate(file, parameter, context, error);
    return true;
  }
};
}  // namespace cruxclient_generator

int main(int argc, char *argv[]) {
  cruxclient_generator::Generator generator;
  PluginMain(argc, argv, &generator);
  return 0;
}
