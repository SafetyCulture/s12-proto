// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

#pragma once
#include <string>

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
class LegacyGenerator {
 public:
  void Generate(
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
  std::string GetMethodSignature(
    const std::string& service_name,
    const std::string& method_name) const;
  void PrintMethodNames(
    google::protobuf::io::Printer *printer,
    const google::protobuf::ServiceDescriptor *service) const;
  void PrintHeaderIncludes(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintHeaderMethods(
    google::protobuf::io::Printer *printer,
    const google::protobuf::ServiceDescriptor *service,
    bool is_virtual,
    bool is_override = false) const;
  void PrintMockHeaderMethods(
    google::protobuf::io::Printer *printer,
    const google::protobuf::ServiceDescriptor *service) const;
  void PrintHeaderInterfaces(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintHeaderClients(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintHeaderMockClients(
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
  void PrintInvokeMethod(
    google::protobuf::io::Printer *printer,
    const google::protobuf::ServiceDescriptor *service) const;
  void PrintSourceClients(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
  void PrintSourceEpilogue(
    google::protobuf::io::Printer *printer,
    const google::protobuf::FileDescriptor *file) const;
};

}  // namespace cruxclient_generator
