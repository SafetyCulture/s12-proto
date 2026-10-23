// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

#pragma once
#include <cctype>
#include <functional>
#include <string>
#include <string_view>
#include <set>
#include <sstream>
#include <google/protobuf/descriptor.h>

using google::protobuf::Descriptor;

namespace cruxclient_generator {
// The following 5 functions were copied from
// google/protobuf/src/google/protobuf/stubs/strutil.h

inline bool HasPrefixString(std::string_view str, std::string_view prefix) {
  return str.size() >= prefix.size() &&
         str.compare(0, prefix.size(), prefix) == 0;
}

inline std::string StripPrefixString(std::string_view str, std::string_view prefix) {
  if (HasPrefixString(str, prefix)) {
    return std::string(str.substr(prefix.size()));
  } else {
    return std::string(str);
  }
}

inline bool HasSuffixString(std::string_view str, std::string_view suffix) {
  return str.size() >= suffix.size() &&
         str.compare(str.size() - suffix.size(), suffix.size(), suffix) == 0;
}

inline std::string StripSuffixString(std::string_view str, std::string_view suffix) {
  if (HasSuffixString(str, suffix)) {
    return std::string(str.substr(0, str.size() - suffix.size()));
  } else {
    return std::string(str);
  }
}

inline void ReplaceCharacters(std::string *s, const char *remove, char replacewith) {
  const char *str_start = s->c_str();
  const char *str = str_start;
  for (str = strpbrk(str, remove); str != nullptr;
       str = strpbrk(str + 1, remove)) {
    (*s)[str - str_start] = replacewith;
  }
}

// The following function was copied from
// google/protobuf/src/google/protobuf/compiler/cpp/cpp_helpers.cc

inline std::string StripProto(std::string_view filename) {
  if (HasSuffixString(filename, ".protodevel")) {
    return StripSuffixString(filename, ".protodevel");
  } else {
    return StripSuffixString(filename, ".proto");
  }
}

// The following 3 functions were copied from
// grpc/src/compiler/cpp_generator.cc

inline std::string FilenameIdentifier(std::string_view filename) {
  std::string result;
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

inline std::string StringReplace(std::string_view str_view, std::string_view from, std::string_view to,
                            bool replace_all) {
  std::string str(str_view);
  size_t pos = 0;

  do {
    pos = str.find(from, pos);
    if (pos == std::string::npos) {
      break;
    }
    str.replace(pos, from.length(), to);
    pos += to.length();
  } while (replace_all);

  return str;
}

inline std::string StringReplace(std::string_view str, std::string_view from, std::string_view to) {
  return StringReplace(str, from, to, true);
}

// The following 4 functions were copied from
// grpc/src/compiler/generator_helpers.h
inline std::vector<std::string> tokenize(std::string_view input,
                                    std::string_view delimiters) {
  std::vector<std::string> tokens;
  size_t pos, last_pos = 0;

  for (;;) {
    bool done = false;
    pos = input.find_first_of(delimiters, last_pos);
    if (pos == std::string_view::npos) {
      done = true;
      pos = input.length();
    }

    tokens.push_back(std::string(input.substr(last_pos, pos - last_pos)));
    if (done) return tokens;

    last_pos = pos + 1;
  }
}

// The following 3 functions were copied from
// grpc/src/compiler/cpp_generator_helpers.h

inline std::string DotsToColons(std::string_view name) {
  return StringReplace(name, ".", "::");
}

inline std::string DotsToUnderscores(std::string_view name) {
  return StringReplace(name, ".", "_");
}

inline std::string DotsToSlashs(std::string_view name) {
  return StringReplace(name, ".", "/");
}

inline std::string UnderscoresToDollar(const std::string &name) {
  return StringReplace(name, "_", "$");
}

inline std::string UnderscoresToDots(const std::string &name) {
  return StringReplace(name, "_", ".");
}

inline std::string ToLower(std::string_view input) {
  std::string output(input);
  std::transform(output.begin(), output.end(), output.begin(),
  [](unsigned char c){ return std::tolower(c); });
  return output;
}

inline std::string join(
  const std::vector<std::string> &elements,
  const std::string &separator) {
  if (!elements.empty()) {
    std::stringstream ss;
    auto it = elements.cbegin();
    while (true) {
      ss << *it++;
      if (it != elements.cend())
        ss << separator;
      else
        return ss.str();
    }
  }
  return "";
}

inline std::string ToCamelCase(std::string_view text) {
  std::stringstream result;
  bool is_new_word = true;

  for (const auto ch : text) {
    is_new_word = is_new_word || ch == '_';
    if (ch == '_') continue;
    if (std::isalpha(ch)) {
      if (is_new_word) {
        result << static_cast<char>(std::toupper(ch));
        is_new_word = false;
      } else {
        result << static_cast<char>(std::tolower(ch));
      }
    } else {
      result << ch;
    }
  }

  return result.str();
}

// Swift module of the SwiftProtobuf types that a bare swift_module option selects.
inline constexpr const char *kDefaultSwiftModule = "S12Protos";

// Root package of the pure Kotlin messages that a bare kotlin_package option selects.
inline constexpr const char *kDefaultKotlinPackage = "com.safetyculture.s12kt";

inline std::string_view TrimString(std::string_view text) {
  while (!text.empty() && std::isspace(static_cast<unsigned char>(text.front()))) text.remove_prefix(1);
  while (!text.empty() && std::isspace(static_cast<unsigned char>(text.back()))) text.remove_suffix(1);
  return text;
}

// Reads one option from a comma-separated protoc parameter string. Returns an empty string
// if the key is absent, which means the output of the option is off. A bare key, or a key
// with an empty value, selects the fallback. The last occurrence of a key wins, and other
// keys are ignored.
inline std::string ParseOption(std::string_view parameter, std::string_view key, std::string_view fallback) {
  std::string value;
  for (const auto &item : tokenize(parameter, ",")) {
    const std::string_view entry = TrimString(item);
    const size_t equals = entry.find('=');
    if (TrimString(entry.substr(0, equals)) != key) continue;
    value = equals == std::string_view::npos ? "" : std::string(TrimString(entry.substr(equals + 1)));
    if (value.empty()) value = fallback;
  }
  return value;
}

// Type prefix of the SwiftProtobuf types of a file. The swift_prefix file option wins,
// even when it is empty. Without it, each part of the proto package is in upper camel
// case, and "_" joins and ends the parts. For example "s12.legacy_ratelimits.v1" gives
// "S12_LegacyRatelimits_V1_". An empty package gives "".
inline std::string SwiftProtobufPrefix(const google::protobuf::FileDescriptor *file) {
  if (file->options().has_swift_prefix()) {
    return std::string(file->options().swift_prefix());
  }
  const std::string package(file->package());
  if (package.empty()) return "";
  std::string prefix;
  bool make_upper = true;
  for (char c : package) {
    if (c == '_') {
      make_upper = true;
    } else if (c == '.') {
      make_upper = true;
      prefix += '_';
    } else {
      if (prefix.empty() && c >= '0' && c <= '9') prefix += '_';
      if (make_upper && c >= 'a' && c <= 'z') c = static_cast<char>(c - 'a' + 'A');
      make_upper = false;
      prefix += c;
    }
  }
  return prefix + "_";
}

// Swift module of the SwiftProtobuf types of a file. The types of the bundled google.protobuf files
// (isBundledProto in protoc-gen-swift) are in the SwiftProtobuf runtime, not in the module of the other protos.
inline std::string SwiftProtobufModule(const google::protobuf::FileDescriptor *file, const std::string &module) {
  static const std::set<std::string, std::less<>> bundled = {
    "google/protobuf/any.proto", "google/protobuf/api.proto", "google/protobuf/descriptor.proto",
    "google/protobuf/duration.proto", "google/protobuf/empty.proto", "google/protobuf/field_mask.proto",
    "google/protobuf/source_context.proto", "google/protobuf/struct.proto",
    "google/protobuf/timestamp.proto", "google/protobuf/type.proto", "google/protobuf/wrappers.proto",
  };
  const bool is_bundled = file->package() == "google.protobuf" && bundled.count(std::string_view(file->name())) > 0;
  return is_bundled ? "SwiftProtobuf" : module;
}

// Names that SwiftProtobuf does not use for a type (reservedTypeNames in the NamingUtils of protoc-gen-swift).
inline bool IsSwiftReservedTypeName(std::string_view name) {
  static const std::set<std::string, std::less<>> reserved = {
    "SwiftProtobuf", "Extensions", "protoMessageName", "decodeMessage", "traverse",
    "isInitialized", "unknownFields", "debugDescription", "description", "dynamicType",
    "hashValue", "Type", "Protocol", "Swift", "Equatable", "Hashable", "Sendable",
    // Keywords used in declarations
    "associatedtype", "class", "deinit", "enum", "extension", "fileprivate", "func",
    "import", "init", "inout", "internal", "let", "open", "operator", "private",
    "protocol", "public", "static", "struct", "subscript", "typealias", "var",
    // Keywords used in statements
    "break", "case", "continue", "default", "defer", "do", "else", "fallthrough", "for",
    "guard", "if", "in", "repeat", "return", "switch", "where", "while",
    // Keywords used in expressions and types
    "as", "Any", "catch", "false", "is", "nil", "rethrows", "super", "self", "Self",
    "throw", "throws", "true", "try",
    // Common types and special variables
    "Bool", "Data", "Double", "Float", "Int", "Int32", "Int64", "String", "UInt",
    "UInt32", "UInt64", "__COLUMN__", "__FILE__", "__FUNCTION__", "__LINE__",
  };
  return reserved.count(name) > 0;
}

// Adds the "Message" suffix to a reserved name, as SwiftProtobuf does. A name that
// already ends in "Message" gets the suffix again if its root is reserved.
inline std::string SanitizeSwiftMessageName(std::string_view name) {
  static constexpr std::string_view kDisambiguator = "Message";
  const bool all_underscores = !name.empty() && name.find_first_not_of('_') == std::string_view::npos;
  if (IsSwiftReservedTypeName(name) || all_underscores) {
    return std::string(name) + std::string(kDisambiguator);
  }
  if (HasSuffixString(name, kDisambiguator)) {
    return SanitizeSwiftMessageName(StripSuffixString(name, kDisambiguator)) + std::string(kDisambiguator);
  }
  return std::string(name);
}

// SwiftProtobuf type name of a message. A nested message is a nested Swift type.
inline std::string SwiftProtobufTypeName(const Descriptor *descriptor) {
  if (descriptor->containing_type() != nullptr) {
    return SwiftProtobufTypeName(descriptor->containing_type()) + "." +
           SanitizeSwiftMessageName(std::string(descriptor->name()));
  }
  return SanitizeSwiftMessageName(SwiftProtobufPrefix(descriptor->file()) + std::string(descriptor->name()));
}

// Translator type of a message in the swiftxx header. It has its own namespace
// because the ObjC translator of the same message is in djinni::<type>.
inline std::string SwiftxxTranslatorName(const std::string &cpp_type_name) {
  return "djinni::swift::" + cpp_type_name + "::Translator";
}

// Kotlin package of the messages of a file: the root package, then the proto package. The
// leading "s12" segment of the proto package is not repeated, because the root package
// already names it. For example "s12.assets.v1" gives "com.safetyculture.s12kt.assets.v1".
inline std::string KotlinPackage(const google::protobuf::FileDescriptor *file, const std::string &root) {
  std::string package(file->package());
  if (package == "s12") return root;
  package = StripPrefixString(package, "s12.");
  if (package.empty()) return root;
  return root + "." + package;
}

// The name Companion would clash with the companion object of the enclosing Kotlin class.
inline std::string KotlinClassName(const std::string &name) {
  return name == "Companion" ? name + "Message" : name;
}

// Kotlin class name of a message inside its package. A nested message is a nested Kotlin
// class, so the separator is "." for the source name and "$" for the JVM binary name.
inline std::string KotlinTypeName(const Descriptor *descriptor, const std::string &separator) {
  if (descriptor->containing_type() != nullptr) {
    return KotlinTypeName(descriptor->containing_type(), separator) + separator +
           KotlinClassName(std::string(descriptor->name()));
  }
  return KotlinClassName(std::string(descriptor->name()));
}

// Fully qualified Kotlin class name of a message, for example "com.safetyculture.s12kt.assets.v1.Outer.Inner".
inline std::string KotlinQualifiedName(const Descriptor *descriptor, const std::string &root) {
  return KotlinPackage(descriptor->file(), root) + "." + KotlinTypeName(descriptor, ".");
}

// JVM binary name of the Kotlin class of a message, for example "com/safetyculture/s12kt/assets/v1/Outer$Inner".
inline std::string KotlinBinaryName(const Descriptor *descriptor, const std::string &root) {
  return DotsToSlashs(KotlinPackage(descriptor->file(), root)) + "/" + KotlinTypeName(descriptor, "$");
}

// JNI translator type of the pure Kotlin class of a message. It has its own namespace
// because the translator of the Java class of the same message is in djinni::<type>.
inline std::string KotlinJNITranslatorName(const std::string &cpp_type_name) {
  return "djinni::kotlin::" + cpp_type_name + "::Translator";
}

inline std::string ClassName(const Descriptor *descriptor, bool qualified) {
  // Find "outer", the descriptor of the top-level message in which
  // "descriptor" is embedded.
  const Descriptor *outer = descriptor;
  while (outer->containing_type() != NULL) outer = outer->containing_type();

  const std::string outer_name(outer->full_name());
  std::string inner_name(descriptor->full_name().substr(outer_name.size()));

  if (qualified) {
    return DotsToColons(outer_name) + DotsToUnderscores(inner_name);
  } else {
    return std::string(outer->name()) + DotsToUnderscores(inner_name);
  }
}

}  // namespace cruxclient_generator
