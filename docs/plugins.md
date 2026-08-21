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

**Method options** (from `s12/flags/permissions/permissions.proto`):

| Option | Semantic |
| --- | --- |
| `required_flags` | The caller must hold **every** listed permission. |
| `required_permissions` | One requirement per option line, all of which must be satisfied. Each requirement sets exactly one of `any_of` (hold at least one) or `all_of` (hold every one). |
| `required_scope` | The credentials scope must contain this value. |

```proto
// admin:training AND (write:training OR write:course)
rpc UpdateCourse(Request) returns (Response) {
  option (s12.flags.permissions.required_permissions) = {all_of: ["admin:training"]};
  option (s12.flags.permissions.required_permissions) = {any_of: ["write:training", "write:course"]};
}
```

Repeated message options do not coalesce, so each line is a separate requirement and adding
one narrows the gate. A requirement setting both operators, or neither, fails generation.
`required_flags` and `required_permissions` may both be set on a method, in which case both
apply. Permission checks (but not the scope check) are short-circuited for an admin-scoped
call.

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
