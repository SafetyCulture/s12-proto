# Plugin Reference

This document describes each `protoc` plugin provided by this repository.

---

## protoc-gen-govalidator

**Purpose**: Generates field-level validation code for Go from proto annotations defined in `s12/protobuf/proto/`.

**Generates**: `*.validator.pb.go` — a `Validate()` method on each message that enforces rules such as `required`, `min_length`, `max_length`, `regex`, `prefix`, `url`, and more.

**Install**:
```sh
go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator@latest
```

**Usage**:
```sh
protoc \
  -I. \
  --go_out=paths=source_relative:. \
  --govalidator_out=paths=source_relative:. \
  path/to/*.proto
```

---

## protoc-gen-s12perm

**Purpose**: Generates gRPC permission interceptors from proto service annotations. Enforces RBAC rules at the transport layer for gRPC services using SafetyCulture's permission model.

**Generates**: `*.perm.pb.go` — server interceptors that check permissions defined via `s12/flags/permissions/` annotations on each RPC method.

**Install**:
```sh
go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-s12perm@latest
```

**Usage**:
```sh
protoc \
  -I. \
  --s12perm_out=paths=source_relative:. \
  path/to/*.proto
```

---

## protoc-gen-logger

**Purpose**: Generates structured logging support for gRPC services.

**Generates**: `*.logger.pb.go` — logging hooks for RPC methods.

**Install**:
```sh
go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-logger@latest
```

**Usage**:
```sh
protoc \
  -I. \
  --logger_out=. \
  path/to/*.proto
```

---

## protoc-gen-gogrpcmock

**Purpose**: Generates mock gRPC server implementations for use in tests.

**Generates**: `*.mock.pb.go` — mock implementations of gRPC service interfaces.

**Install**:
```sh
go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-gogrpcmock@latest
```

**Usage**:
```sh
protoc \
  -I. \
  --gogo_out=plugin=grpc=:. \
  --gogrpcmock_out=. \
  path/to/*.proto
```

---

## protoc-gen-cruxclient (C++)

**Purpose**: Generates C++ client code, mock servers, and Djinni cross-language bindings (ObjC, JNI) from proto service definitions. Targets the CRUX Engine mobile client SDK.

**Generates**: `.crux.api.h`, `.crux.api.cc`, `.mock.h`, `.djinni.yaml`, `.djinni.objc.h`, `.djinni.jni.h`

**Install** (requires C++17 toolchain, protobuf, and gRPC libraries):
```sh
make install-cruxclient
```

---

## protoc-gen-cruxclient-go

**Purpose**: Go port of `protoc-gen-cruxclient`. Generates identical output using the `bufbuild/protoplugin` framework — no C++ toolchain required.

**Generates**: `.crux.api.h`, `.crux.api.cc`, `.mock.h`, `.djinni.yaml`, `.djinni.objc.h`, `.djinni.jni.h`

**Install**:
```sh
go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-cruxclient-go@latest
```

**Usage**:
```sh
protoc \
  -I. \
  --cruxclient-go_out=. \
  path/to/*.proto
```
