# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- `protoc-gen-cruxclient` (C++) is frozen: it does not get the new options. `protoc-gen-cruxclient-go` is the only plugin that changes.

### Added
- `protoc-gen-cruxclient-go`: new option `kotlin_messages`. It writes pure Kotlin message classes with a wire codec. A message is a `data class` with default values, and nullability follows the proto. Each message has `toByteArray()`, `getSerializedSize()` and `parseFrom(ByteArray)`. The codec has no dependency except the Kotlin standard library, and its bytes are the same as the bytes of Java lite. An enum is a `sealed class` with an `UNRECOGNIZED(number)` class, and the reader stops at a nesting depth of 100. `parseFrom` throws for a missing proto2 required field after it has merged the parts of a message, and a string with `utf8_validation = NONE` accepts invalid UTF-8, as Java lite does. A bare `kotlin_messages` uses the root package of `kotlin_package`, or `com.safetyculture.s12kt` if `kotlin_package` is absent, and `kotlin_messages=<Root>` sets it. Code must give a value for each proto2 required field, and `toByteArray()` throws `UninitializedMessageException` for a decoded message that lacks one. Both root options must name the same root package, and each segment must be a Kotlin identifier. The plugin stops with an error when two proto names map to the same Kotlin name in a package or a class, and when the root package is not a Kotlin package name. `parseFrom` throws for a cycle of required message fields that has no value, as Java lite does. With the option off, the output does not change. The new bare option `messages_only` turns off all cruxclient output in a run, so that `kotlin_messages,messages_only` writes only the `.pb.kt` files, the runtime file `S12Wire.kt` and the shrinker keep rules `s12kt.pro`; it needs `kotlin_messages`. The plugin now stops with an error for an option it does not know (it accepts `paths`) and for a value on `messages_only`; master ignored both. `test/wire_compat.sh` compares the codec with Java lite. The plugin also declares support for proto3 `optional` and for editions up to 2023.
- `protoc-gen-cruxclient-go`: new option `swift_messages`. The plugin writes the Swift messages, one `.pb.swift` file for each proto file, and the files are the files of `protoc-gen-swift` 1.38.1 for the same descriptors. The option takes no value. The options `Visibility`, `FileNaming`, `EnumGeneration` and the other `protoc-gen-swift` options work as they do there, and the file option `swift_prefix` sets the type names. The output of the 530 whitelist files and of the 1490 files of the schema without prefixes has no difference from `protoc-gen-swift`, and a set of protos with groups, extensions, message sets and editions has none either. A run with `messages_only` and without a messages option stops with an error. A file of edition 2024 or later is not accepted, while `protoc-gen-swift` 1.38.1 accepts edition 2024. A proto file with a message set works only together with `messages_only`, because the cruxclient files need the runtime descriptors of Go, which refuse a message set. With the option off, the output does not change. A Linux agent can now write the Swift messages, because the plugin is a Go program.
- `protoc-gen-cruxclient-go`: new option `kotlin_package`. A bare `kotlin_package` uses the root package `com.safetyculture.s12kt`. With the option on, each `.djinni.yaml` record has `kotlinTranslator`, `kotlinHeader` and `kotlinTypeSignature` keys in `jni:` and a `jniTypename` key in `kotlin:`. The plugin also writes `.djinni.kotlin.jni.h` with one translator for each message, which binds the pure Kotlin class. With the option off, the output does not change. With the option on, each record has `deriving(eq, parcelable)` without `ord`, and the translator throws for a null or wrong-type object and for a message with a missing required field.
- `protoc-gen-cruxclient-go`: new option `swift_module`. A bare `swift_module` uses the module name `S12Protos`. `swift_module=<Name>` sets the module name. With the option on, each `.djinni.yaml` record has a `swift:` block (`typename`, `module`, `generic`, `translator`, `translator.module`) and a `swiftxx:` block (`translator`, `header`). The plugin also writes `.djinni.swiftxx.h` with one translator for each message. The translator moves a message as serialized bytes between C++ and SwiftProtobuf. The Swift type names follow `protoc-gen-swift`: the `swift_prefix` file option sets the prefix, and a reserved nested name such as `Data` gets the `Message` suffix. The types of the 11 bundled `google/protobuf/*.proto` files use the module `SwiftProtobuf`, because the SwiftProtobuf runtime has them. With the option off, the output is byte-identical to the output of the previous release. The translator throws `ErrorValue` for a message with a missing required field.
- `protoc-gen-cruxclient-go`: with `swift_module` or `kotlin_package`, the translators name a C++ keyword in a package or message, and a message name that would hide a method of the message class, with a trailing `_`, as protoc does. The plugin stops with an error for a nested message named `Translator`, a package with any segment `std` or a message named `std`, with `swift_module` a package whose first segment is `swift`, and with `kotlin_package` a package whose first segment is `kotlin` and a nested message named `JNIInfo`.
- `protoc-gen-govalidator`: U+2062 INVISIBLE TIMES and U+2064 INVISIBLE PLUS are allowed by default in `validator.string` and `validator.unsafe_string`. These two carry the AI content mark Mitti applies to generated text so it is detectable as artificially generated, which is what EU AI Act Article 50(2) asks a provider to make possible. The Regulation names no technique, so the codepoints are Mitti's choice. The allow list has no token for Unicode category Cf and no symbol category maps to it, so before this no field option could admit them and a marked write failed with "value must only have valid characters". Only these two codepoints are added, not the category (AI-1355).
- `protoc-gen-govalidator`: `len` counts the AI content mark like any other character, so a limit means the length of the value that gets stored. A value landing within 24 runes of its maximum therefore cannot carry a mark and is refused with a clear error naming the limit. v1.43.0 measured the mark out of the length, which let a value through that a column with the same limit then rejected, and required every such column to carry headroom nothing in the schema explained (AI-1355).
- `AIMarkStripper` in `s12/protobuf/proto`, a replacer that removes the two mark carriers. `reject_url` now matches against a stripped copy, because an invisible character inside a host name defeats a literal URL match (AI-1355).
- `AIMarkBeforeDotMatcher` in `s12/protobuf/proto`. `break_partial_url` removes only a run of mark carriers sitting in front of a dot, which is the one placement that hides a partial URL from `BreakURLMatcher`. Every other copy in the value survives, so a field that breaks partial URLs stays marked. The Mitti codec inserts only at a word start and never produces that placement, so this removes nothing the codec wrote (AI-1355).
- GitHub Actions workflows: CI, proto lint, and release
- `buf.yaml` and `buf.gen.yaml` for buf toolchain integration
- `.golangci.yml` linter configuration
- `docs/plugins.md` plugin reference documentation
- Extended pre-commit hooks: gofmt, goimports, buf lint, golangci-lint

### Changed
- Bumped Go minimum version from 1.18 to 1.24 in root module and `protoc-gen-s12perm/example`

### Fixed
- `protoc-gen-govalidator`: a message with two non-optional `simple_string` fields carrying `min_len`/`max_len` generated a `Validate` func declaring `length` twice, so the package did not compile. The length variable is now named per field (FG-6837). Generated code only, validation behaviour is unchanged.

## [v1.38.0] - 2026-03-03

### Added
- `protoc-gen-cruxclient-go`: Go port of the C++ cruxclient plugin, generating equivalent C++/ObjC/JNI/YAML output using the `bufbuild/protoplugin` framework with golden-file test coverage (#153)

### Changed
- Upgraded `protoc-gen-cruxclient` for protobuf v33 compatibility (#154)

## [v1.37.0]

### Added
- `protoc-gen-govalidator`: support for `prefix` in string validation rule (PA-1735, #151)

### Changed
- `protoc-gen-cruxclient`: added kotlin type metadata to generated Djinni `.yaml` files (SFT-117, #150)
- Bumped validator URL max length to 2048 per documentation (#149)

## [v1.36.0]

### Removed
- Deleted legacy protobuf generator (MOB-3701, #148)

### Changed
- `protoc-gen-govalidator`: when both `reject` and `break_url` are defined, `reject` takes precedence (UT-5331, #145)
- `protoc-gen-govalidator`: added URL sanitization for user fields (UT-5331, #143)
