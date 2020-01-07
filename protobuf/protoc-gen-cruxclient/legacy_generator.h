// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

#pragma once
#include <string>

namespace google::protobuf::compiler {
class GeneratorContext;
}

namespace google::protobuf {
class FileDescriptor;
}

namespace cruxclient_generator {
class LegacyGenerateor {
 public:
  void Generate(
    const google::protobuf::FileDescriptor *file,
    const std::string &parameter,
    google::protobuf::compiler::GeneratorContext *context,
    std::string *error) const;
};

}  // namespace cruxclient_generator
