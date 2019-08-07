// Copyright (c) 2019 SafetyCulture Pty Ltd. All Rights Reserved.

#include <google/protobuf/compiler/code_generator.h>
#include <google/protobuf/compiler/plugin.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/io/printer.h>
#include <google/protobuf/io/zero_copy_stream.h>
#include <iostream>
#include <string>

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
namespace {

// The following 5 functions were copied from
// google/protobuf/src/google/protobuf/stubs/strutil.h

inline bool HasPrefixString(const string &str, const string &prefix) {
  return str.size() >= prefix.size() &&
         str.compare(0, prefix.size(), prefix) == 0;
}

inline string StripPrefixString(const string &str, const string &prefix) {
  if (HasPrefixString(str, prefix)) {
    return str.substr(prefix.size());
  } else {
    return str;
  }
}

inline bool HasSuffixString(const string &str, const string &suffix) {
  return str.size() >= suffix.size() &&
         str.compare(str.size() - suffix.size(), suffix.size(), suffix) == 0;
}

inline string StripSuffixString(const string &str, const string &suffix) {
  if (HasSuffixString(str, suffix)) {
    return str.substr(0, str.size() - suffix.size());
  } else {
    return str;
  }
}

void ReplaceCharacters(string *s, const char *remove, char replacewith) {
  const char *str_start = s->c_str();
  const char *str = str_start;
  for (str = strpbrk(str, remove); str != nullptr;
       str = strpbrk(str + 1, remove)) {
    (*s)[str - str_start] = replacewith;
  }
}

// The following function was copied from
// google/protobuf/src/google/protobuf/compiler/cpp/cpp_helpers.cc

string StripProto(const string &filename) {
  if (HasSuffixString(filename, ".protodevel")) {
    return StripSuffixString(filename, ".protodevel");
  } else {
    return StripSuffixString(filename, ".proto");
  }
}

// The following 3 functions were copied from
// grpc/src/compiler/cpp_generator.cc

string FilenameIdentifier(const string &filename) {
  string result;
  for (unsigned i = 0; i < filename.size(); i++) {
    char c = filename[i];
    if (isalnum(c)) {
      result.push_back(c);
    } else {
      static char hex[] = "0123456789abcdef";
      result.push_back('_');
      result.push_back(hex[(c >> 4) & 0xf]);
      result.push_back(hex[c & 0xf]);
    }
  }
  return result;
}

inline string StringReplace(string str, const string &from, const string &to,
                            bool replace_all) {
  size_t pos = 0;

  do {
    pos = str.find(from, pos);
    if (pos == string::npos) {
      break;
    }
    str.replace(pos, from.length(), to);
    pos += to.length();
  } while (replace_all);

  return str;
}

inline string StringReplace(string str, const string &from, const string &to) {
  return StringReplace(str, from, to, true);
}

// The following 4 functions were copied from
// grpc/src/compiler/generator_helpers.h
inline std::vector<string> tokenize(const string &input,
                                    const string &delimiters) {
  std::vector<string> tokens;
  size_t pos, last_pos = 0;

  for (;;) {
    bool done = false;
    pos = input.find_first_of(delimiters, last_pos);
    if (pos == string::npos) {
      done = true;
      pos = input.length();
    }

    tokens.push_back(input.substr(last_pos, pos - last_pos));
    if (done) return tokens;

    last_pos = pos + 1;
  }
}

// The following 3 functions were copied from
// grpc/src/compiler/cpp_generator_helpers.h

inline string DotsToColons(const string &name) {
  return StringReplace(name, ".", "::");
}

inline string DotsToUnderscores(const string &name) {
  return StringReplace(name, ".", "_");
}

inline string ClassName(const Descriptor *descriptor, bool qualified) {
  // Find "outer", the descriptor of the top-level message in which
  // "descriptor" is embedded.
  const Descriptor *outer = descriptor;
  while (outer->containing_type() != NULL) outer = outer->containing_type();

  const string &outer_name = outer->full_name();
  string inner_name = descriptor->full_name().substr(outer_name.size());

  if (qualified) {
    return DotsToColons(outer_name) + DotsToUnderscores(inner_name);
  } else {
    return outer->name() + DotsToUnderscores(inner_name);
  }
}

// ******************************************* //
//         Code Generation Printing            //
// ******************************************* //

void PrintPrologue(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;
  vars["filename"] = file->name();
  printer->Print("// Generated by the CRUX Service Client C++ plugin.\n");
  printer->Print("// If you make any local change, they will be lost.\n");
  printer->Print(vars, "// source: $filename$\n");
}

void PrintNamespace(Printer *printer, const FileDescriptor *file,
                    bool isEpilogue) {
  std::map<string, string> vars;
  if (!file->package().empty()) {
    std::vector<string> parts = tokenize(StripProto(file->package()), ".");
    for (auto part = parts.begin(); part != parts.end(); part++) {
      vars["part"] = *part;
      if (!isEpilogue) {
        printer->Print(vars, "namespace $part$ {\n");
      } else {
        printer->Print(vars, "}  // namespace $part$\n");
      }
    }
    printer->Print(vars, "\n");
  }
}

void PrintHeaderPrologue(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;
  vars["filename_identifier"] = FilenameIdentifier(file->name());

  PrintPrologue(printer, file);
}

void PrintHeaderIncludes(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;
  vars["filename_base"] = StripProto(file->name());

  printer->Print("#pragma once\n\n");
  printer->Print(vars, "#include \"$filename_base$.grpc.pb.h\"\n\n");
  printer->Print("#include <string>\n");
  printer->Print("#include <memory>\n\n");

  PrintNamespace(printer, file, false);
}

void PrintHeaderMethods(Printer *printer, const ServiceDescriptor *service,
                        bool is_virtual, bool is_override = false) {
  std::map<string, string> vars;
  for (int method_index = 0; method_index < service->method_count();
       ++method_index) {
    const MethodDescriptor *method = service->method(method_index);
    vars["method_name"] = method->name();
    vars["request"] = ClassName(method->input_type(), true);
    vars["response"] = ClassName(method->output_type(), true);

    if (method->client_streaming()) {
      // [RC]: Client Steaming not supported yet
      continue;
    }

    if (method->server_streaming()) {
      vars["response"] = "std::vector<" + vars["response"] + ">";
    }
    if (is_virtual) {
      printer->Print("virtual ");
    }
    printer->Print(vars,
                   "$response$ $method_name$(const $request$& "
                   "request) const");
    if (is_virtual) {
      printer->Print(" = 0");
    }
    if (is_override) {
      printer->Print(" override");
    }
    printer->Print(";\n");
  }
}

void PrintHeaderInterfaces(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;
  for (int service_index = 0; service_index < file->service_count();
       ++service_index) {
    const ServiceDescriptor *service = file->service(service_index);
    vars["service_name"] = service->name();

    printer->Print(vars, "class $service_name$ClientInterface {\n");
    printer->Print(" public:\n");
    printer->Indent();
    printer->Print(vars, "virtual ~$service_name$ClientInterface() {}\n");
    printer->Print("virtual void MakeRequest(const std::string& request_data, const std::string& request_type) const = 0;\n");
    PrintHeaderMethods(printer, service, true);
    printer->Outdent();
    printer->Print("};\n\n");
  }
}

void PrintHeaderClients(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;

  for (int service_index = 0; service_index < file->service_count();
       ++service_index) {
    const ServiceDescriptor *service = file->service(service_index);
    vars["service_name"] = service->name();

    printer->Print(
        vars,
        "class $service_name$Client: public $service_name$ClientInterface {\n");
    printer->Print(" public:\n");
    printer->Indent();
    printer->Print(vars,
                   "explicit $service_name$Client(const "
                   "std::shared_ptr<$service_name$::StubInterface>& "
                   "stub);\n");
    printer->Print("void MakeRequest(const std::string& request_data, const std::string& request_type) const override;\n");
    PrintHeaderMethods(printer, service, false, true);
    printer->Outdent();
    printer->Print("\n");

    printer->Print(" private:\n");
    printer->Indent();
    printer->Print(vars,
                   "std::shared_ptr<$service_name$::StubInterface> mStub;\n\n");
    printer->Outdent();
    printer->Print("};\n\n");
  }
}

void PrintHeaderEpilogue(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;
  vars["filename"] = file->name();
  vars["filename_identifier"] = FilenameIdentifier(file->name());

  PrintNamespace(printer, file, true);
}

void PrintSourcePrologue(Printer *printer, const FileDescriptor *file) {
  PrintPrologue(printer, file);
  printer->Print("\n");
}

void PrintSourceIncludes(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;
  vars["filename_base"] = StripProto(file->name());

  printer->Print(vars, "#include \"$filename_base$.crux.client.pb.h\"\n");
  printer->Print(vars, "#include \"s12_client_support.hpp\"\n\n");

  PrintNamespace(printer, file, false);
}

void PrintSourceClients(Printer *printer, const FileDescriptor *file) {
  std::map<string, string> vars;
  if (!file->package().empty()) {
    vars["package"].append("::");
  }

  for (int service_index = 0; service_index < file->service_count();
       ++service_index) {
    const ServiceDescriptor *service = file->service(service_index);
    vars["service_name"] = service->name();

    printer->Print(vars,
                   "$service_name$Client::$service_name$Client(const "
                   "std::shared_ptr<$service_name$::StubInterface>& "
                   "stub) : mStub(stub) {}\n\n");

    for (int method_index = 0; method_index < service->method_count();
         ++method_index) {
      const MethodDescriptor *method = service->method(method_index);
      vars["method_name"] = method->name();
      vars["request"] = ClassName(method->input_type(), true);
      vars["response"] = ClassName(method->output_type(), true);

      if (method->server_streaming()) {
        vars["response_item"] = vars["response"];
        vars["response"] = "std::vector<" + vars["response"] + ">";
      }

      if (method->client_streaming()) {
        // [RC]: Client Streaming not supported... yet
        continue;
      }

      printer->Print(vars,
                     "$response$ $service_name$Client::$method_name$(const "
                     "$request$& request) const {\n");
      printer->Indent();
      printer->Print(vars, "$response$ response;\n");
      printer->Print("grpc::ClientContext context;\n");
      if (method->server_streaming()) {
        printer->Print(vars, "$response_item$ item;\n");
        printer->Print(
            vars,
            "std::unique_ptr<grpc::ClientReaderInterface<$response_item$>> "
            "stream = mStub->$method_name$(&context, request);\n");
        printer->Print("while (stream->Read(&item)) {\n");
        printer->Indent();
        printer->Print("response.emplace_back(item);\n");
        printer->Outdent();
        printer->Print("}\n");
        printer->Print("grpc::Status status = stream->Finish();\n");
      } else {
        printer->Print(vars,
                       "grpc::Status status = mStub->$method_name$(&context, "
                       "request, &response);\n");
      }
      printer->Print("if (!status.ok()) {\n");
      printer->Indent();
      printer->Print(
          "throw crux::ServiceException(status.error_code(), "
          "status.error_message());\n");
      printer->Outdent();
      printer->Print("}\n");
      printer->Print("return response;\n");
      printer->Outdent();
      printer->Print("}\n");
    }
  }
}

void PrintSourceEpilogue(Printer *printer, const FileDescriptor *file) {
  printer->Print("\n");
  PrintNamespace(printer, file, true);
}

class Generator : public CodeGenerator {
 public:
  Generator() {}
  ~Generator() override {}

  bool Generate(const FileDescriptor *file, const string &parameter,
                GeneratorContext *context, string *error) const override {
    if (!file->service_count()) {
      // No services, nothing to do.
      return true;
    }

    string file_name = StripProto(file->name());

    std::unique_ptr<ZeroCopyOutputStream> header_output(
        context->Open(file_name + ".crux.client.pb.h"));
    Printer header_printer(header_output.get(), '$');
    PrintHeaderPrologue(&header_printer, file);
    PrintHeaderIncludes(&header_printer, file);
    PrintHeaderInterfaces(&header_printer, file);
    PrintHeaderClients(&header_printer, file);
    PrintHeaderEpilogue(&header_printer, file);

    std::unique_ptr<ZeroCopyOutputStream> source_output(
        context->Open(file_name + ".crux.client.pb.cc"));
    Printer source_printer(source_output.get(), '$');
    PrintSourcePrologue(&source_printer, file);
    PrintSourceIncludes(&source_printer, file);
    PrintSourceClients(&source_printer, file);
    PrintSourceEpilogue(&source_printer, file);

    return true;
  }
};

}  // namespace
}  // namespace cruxclient_generator

int main(int argc, char *argv[]) {
  cruxclient_generator::Generator generator;
  PluginMain(argc, argv, &generator);
  return 0;
}
