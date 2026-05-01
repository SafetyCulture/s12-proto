# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GitHub Actions workflows: CI, proto lint, and release
- `buf.yaml` and `buf.gen.yaml` for buf toolchain integration
- `.golangci.yml` linter configuration
- `docs/plugins.md` plugin reference documentation
- Extended pre-commit hooks: gofmt, goimports, buf lint, golangci-lint

### Changed
- Bumped Go minimum version from 1.18 to 1.24 in root module and `protoc-gen-s12perm/example`

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
