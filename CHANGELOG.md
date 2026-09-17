# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `protoc-gen-govalidator`: U+2062 INVISIBLE TIMES and U+2064 INVISIBLE PLUS are allowed by default in `validator.string` and `validator.unsafe_string`. These two carry the AI content mark that every AI generated string must hold under EU AI Act Article 50(2). The allow list has no token for Unicode category Cf and no symbol category maps to it, so before this no field option could admit them and a marked write failed with "value must only have valid characters". Only these two codepoints are added, not the category (AI-1355).
- `protoc-gen-govalidator`: `len` on `validator.string` and `validator.unsafe_string` is measured on the value with the AI content mark removed, so the mark costs a field nothing. Without this a value already at a field maximum failed the moment it was marked, between 3 and 72 units over depending on the cap and on whether the field counts runes or bytes (AI-1355).
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
