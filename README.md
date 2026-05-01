# Protocol Buffer tools used by SafetyCulture

**WARNING:** Public Repo

[![CI](https://github.com/SafetyCulture/s12-proto/actions/workflows/ci.yml/badge.svg)](https://github.com/SafetyCulture/s12-proto/actions/workflows/ci.yml)
[![Proto Lint](https://github.com/SafetyCulture/s12-proto/actions/workflows/proto-lint.yml/badge.svg)](https://github.com/SafetyCulture/s12-proto/actions/workflows/proto-lint.yml)

## Pre-requisites

* [Go](https://golang.org/doc/install) 1.24+
* [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/)
* [buf](https://buf.build/docs/installation) (for proto linting and formatting)

Install `protoc-gen-go`:
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Generating Go code from proto definitions

```sh
make generate
```

To generate and run validator tests:
```sh
make generate && make govalidator && make govalidator-valtest
```

## Plugins

See [docs/plugins.md](docs/plugins.md) for the full plugin reference, including install instructions and usage examples for each plugin.

| Plugin | Description |
|--------|-------------|
| `protoc-gen-govalidator` | Field-level validation for Go |
| `protoc-gen-s12perm` | gRPC permission interceptors |
| `protoc-gen-logger` | Structured logging for gRPC services |
| `protoc-gen-gogrpcmock` | gRPC mock server implementations |
| `protoc-gen-cruxclient` | C++ client + Djinni bindings (C++ toolchain) |
| `protoc-gen-cruxclient-go` | C++ client + Djinni bindings (Go, no C++ toolchain needed) |

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
cd protobuf/protoc-gen-cruxclient-go && go test ./...
```
