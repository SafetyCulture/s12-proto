// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

#pragma once
#include <string>
#include <map>

namespace google::protobuf {
class FileDescriptor;
class ServiceDescriptor;

namespace compiler {
class GeneratorContext;
}

namespace io {
class Printer;
}
}  // namespace google::protobuf

namespace cruxclient_generator {
class APIGenerator {
 public:
  void Generate(
    const google::protobuf::FileDescriptor *file,
    const std::string &parameter,
    google::protobuf::compiler::GeneratorContext *context,
    std::string *error) const;

  void GenerateDjinniSupport(
    const google::protobuf::FileDescriptor *file,
    const std::string &parameter,
    google::protobuf::compiler::GeneratorContext *context,
    std::string *error) const;

 private:
  void PrintPrologue(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintNamespace(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file,
    bool isEpilogue) const;
  void PrintHeaderPrologue(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintHeaderIncludes(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintHeaderAPIs(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintHeaderEpilogue(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintSourcePrologue(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintSourceIncludes(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintSourceAPIs(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintSourceEpilogue(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;

  void GenerateAPIHeader(
    const google::protobuf::FileDescriptor *file,
    const std::string &parameter,
    google::protobuf::compiler::GeneratorContext *context) const;

  void GenerateAPISource(
    const google::protobuf::FileDescriptor *file,
    const std::string &parameter,
    google::protobuf::compiler::GeneratorContext *context) const;

  // Djinni Support
  void GenerateDjinniYAML(
    const google::protobuf::FileDescriptor *file,
    google::protobuf::compiler::GeneratorContext *context) const;

  void PrintDjinniYAML(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file,
    const std::string& main_file_name) const;

  void GenerateDjinniObjcSupport(
    const google::protobuf::FileDescriptor *file,
    google::protobuf::compiler::GeneratorContext *context) const;

  void PrintDjinniObjcSupport(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;

  void GenerateDjinniJNISupport(
    const google::protobuf::FileDescriptor *file,
    google::protobuf::compiler::GeneratorContext *context) const;

  void PrintDjinniJNISupport(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;

};

}  // namespace cruxclient_generator
