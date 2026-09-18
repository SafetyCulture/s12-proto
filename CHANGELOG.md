# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
