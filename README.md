# Protocol Buffer tools used by SafetyCulture

**WARNING:** Public Repo

[![CI](https://github.com/SafetyCulture/s12-proto/actions/workflows/ci.yml/badge.svg)](https://github.com/SafetyCulture/s12-proto/actions/workflows/ci.yml)
[![Proto Lint](https://github.com/SafetyCulture/s12-proto/actions/workflows/proto-lint.yml/badge.svg)](https://github.com/SafetyCulture/s12-proto/actions/workflows/proto-lint.yml)

## Pre-requisites

* [mise](https://mise.jdx.dev/getting-started.html) (manages Go, buf, golangci-lint, and protoc-gen-go)
* [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/) (`brew install protobuf`)

Bootstrap tooling:
```sh
mise install
```

## Common tasks

```sh
mise run generate           # Regenerate the protoc-gen-go output under s12/
mise run verify:generated   # Fail if the fixtures differ from a fresh generation or one is uncommitted
mise run update:generated   # Rewrite the committed fixtures from a fresh generation
mise run test               # Run tests (root module and protoc-gen-cruxclient-go)
mise run lint               # Run Go linter
mise run lint:proto         # Lint proto files with buf
```

Generated code is committed, so changing a plugin template or a `.proto` means
regenerating its fixtures in the same commit: `mise run update:generated`, then
`git add -A && git diff --cached` — a new fixture is untracked, so `git diff` on
its own would not show it. CI runs `verify:generated`, which generates into a
temp directory and diffs, so it never writes the tree; a fixture that is out of
date, missing or not yet added fails it. Both accept an optional generator
(`mise run verify:generated -- s12perm`), and
`./scripts/generated-fixtures.sh list` shows which generators are checked.

To regenerate the fixtures and run the validator tests:
```sh
mise run update:generated && go test ./protobuf/protoc-gen-govalidator/...
```

## Plugins

See [docs/plugins.md](docs/plugins.md) for the full plugin reference, including install instructions and usage examples for each plugin.

| Plugin | Description |
|--------|-------------|
| `protoc-gen-govalidator` | Field-level validation for Go |
| `protoc-gen-s12perm` | gRPC permission interceptors |
| `protoc-gen-logger` | Structured logging for gRPC services |
| `protoc-gen-gogrpcmock` | gRPC mock server implementations |
| `protoc-gen-cruxclient` | C++ client + Djinni bindings (C++ toolchain). Frozen: new options are only in `protoc-gen-cruxclient-go` |
| `protoc-gen-cruxclient-go` | C++ client + Djinni bindings (Go, no C++ toolchain needed), and pure Kotlin messages with a wire codec (option `kotlin_messages`) |

## Proto linting with buf

```sh
buf lint
buf breaking --against '.git#branch=master'
```

## Running tests

```sh
# Root module (all plugins except cruxclient-go)
go test ./...

# protoc-gen-cruxclient-go
(cd protobuf/protoc-gen-cruxclient-go && go test ./...)

# protoc-gen-s12perm example (separate module, needs private sc-go.io modules)
(cd protobuf/protoc-gen-s12perm/example && go test ./...)

# Kotlin and Java lite wire compatibility (needs go, protoc, kotlinc, javac, java, curl and shasum)
protobuf/protoc-gen-cruxclient-go/test/wire_compat.sh
```
