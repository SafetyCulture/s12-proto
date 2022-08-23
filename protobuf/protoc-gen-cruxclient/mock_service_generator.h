// Copyright (c) 2022 SafetyCulture Pty Ltd. All Rights Reserved.

#pragma once
#include <string>
#include <map>
#include <vector>

namespace google::protobuf {
class FileDescriptor;
class ServiceDescriptor;
class Descriptor;

namespace compiler {
class GeneratorContext;
}

namespace io {
class Printer;
}
}  // namespace google::protobuf

namespace cruxclient_generator {
class MockServiceGenerator final {
 public:
  void Generate(
    const google::protobuf::FileDescriptor *file,
    const std::string &parameter,
    google::protobuf::compiler::GeneratorContext *context,
    std::string *error) const;

 private:
  void GenerateHeader(
    const google::protobuf::FileDescriptor *file,
    const std::string &parameter,
    google::protobuf::compiler::GeneratorContext *context) const;
};

}  // namespace cruxclient_generator
