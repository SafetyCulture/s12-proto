// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

#pragma once
#include <string>
#include <google/protobuf/descriptor.h>

using google::protobuf::Descriptor;

namespace cruxclient_generator {
// The following 5 functions were copied from
// google/protobuf/src/google/protobuf/stubs/strutil.h

inline bool HasPrefixString(const std::string &str, const std::string &prefix) {
  return str.size() >= prefix.size() &&
         str.compare(0, prefix.size(), prefix) == 0;
}

inline std::string StripPrefixString(const std::string &str, const std::string &prefix) {
  if (HasPrefixString(str, prefix)) {
    return str.substr(prefix.size());
  } else {
    return str;
  }
}

inline bool HasSuffixString(const std::string &str, const std::string &suffix) {
  return str.size() >= suffix.size() &&
         str.compare(str.size() - suffix.size(), suffix.size(), suffix) == 0;
}

inline std::string StripSuffixString(const std::string &str, const std::string &suffix) {
  if (HasSuffixString(str, suffix)) {
    return str.substr(0, str.size() - suffix.size());
  } else {
    return str;
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

inline std::string StripProto(const std::string &filename) {
  if (HasSuffixString(filename, ".protodevel")) {
    return StripSuffixString(filename, ".protodevel");
  } else {
    return StripSuffixString(filename, ".proto");
  }
}

// The following 3 functions were copied from
// grpc/src/compiler/cpp_generator.cc

inline std::string FilenameIdentifier(const std::string &filename) {
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

inline std::string StringReplace(std::string str, const std::string &from, const std::string &to,
                            bool replace_all) {
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

inline std::string StringReplace(std::string str, const std::string &from, const std::string &to) {
  return StringReplace(str, from, to, true);
}

// The following 4 functions were copied from
// grpc/src/compiler/generator_helpers.h
inline std::vector<std::string> tokenize(const std::string &input,
                                    const std::string &delimiters) {
  std::vector<std::string> tokens;
  size_t pos, last_pos = 0;

  for (;;) {
    bool done = false;
    pos = input.find_first_of(delimiters, last_pos);
    if (pos == std::string::npos) {
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

inline std::string DotsToColons(const std::string &name) {
  return StringReplace(name, ".", "::");
}

inline std::string DotsToUnderscores(const std::string &name) {
  return StringReplace(name, ".", "_");
}

inline std::string ClassName(const Descriptor *descriptor, bool qualified) {
  // Find "outer", the descriptor of the top-level message in which
  // "descriptor" is embedded.
  const Descriptor *outer = descriptor;
  while (outer->containing_type() != NULL) outer = outer->containing_type();

  const std::string &outer_name = outer->full_name();
  std::string inner_name = descriptor->full_name().substr(outer_name.size());

  if (qualified) {
    return DotsToColons(outer_name) + DotsToUnderscores(inner_name);
  } else {
    return outer->name() + DotsToUnderscores(inner_name);
  }
}

}  // namespace cruxclient_generator
