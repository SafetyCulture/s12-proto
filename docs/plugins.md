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

**Generates**: `.crux.api.h`, `.crux.api.cc`, `.mock.h`, `.djinni.yaml`, `.djinni.objc.h`, `.djinni.jni.h`, and `.djinni.swiftxx.h` (only with the `swift_module` option)

**Install** (requires C++17 toolchain, protobuf, and gRPC libraries):
```sh
make install-cruxclient
```

**Option `swift_module`**: Turns on the Swift output. A bare `swift_module` uses the module name `S12Protos`. `swift_module=<Name>` sets the module name. With the option on, each `.djinni.yaml` record gets a `swift:` and a `swiftxx:` block, and the plugin writes `.djinni.swiftxx.h`. With the option off, the output does not change. This plugin is the one in the `ghcr.io/safetyculture/protoc-cpp` image. The Go port writes the same bytes.
```sh
protoc \
  -I. \
  --cruxclient_out=swift_module=S12Protos:. \
  path/to/*.proto
```

**Option `kotlin_package`**: Turns on the Kotlin output for the pure Kotlin messages of `protoc-gen-crux-kotlin`. A bare `kotlin_package` uses the root package `com.safetyculture.s12kt`. `kotlin_package=<Root>` sets the root package. With the option on, each `.djinni.yaml` record gets `kotlinTranslator`, `kotlinHeader` and `kotlinTypeSignature` keys in its `jni:` block and a `jniTypename` key in its `kotlin:` block, and the plugin writes `.djinni.kotlin.jni.h`. The Kotlin translator moves a message as serialized bytes between C++ and the Kotlin class, with `toByteArray()`, `getSerializedSize()` and `parseFrom(ByteArray)`. With the option off, the output does not change.
```sh
protoc \
  -I. \
  --cruxclient_out=kotlin_package:. \
  path/to/*.proto
```

---

## protoc-gen-cruxclient-go

**Purpose**: Go port of `protoc-gen-cruxclient`. Generates identical output using the `bufbuild/protoplugin` framework — no C++ toolchain required.

**Generates**: `.crux.api.h`, `.crux.api.cc`, `.mock.h`, `.djinni.yaml`, `.djinni.objc.h`, `.djinni.jni.h`, and `.djinni.swiftxx.h` (only with the `swift_module` option)

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

**Option `swift_module`**: Turns on the Swift output. A bare `swift_module` uses the module name `S12Protos`. `swift_module=<Name>` sets the module name. With the option on, each `.djinni.yaml` record gets a `swift:` and a `swiftxx:` block, and the plugin writes `.djinni.swiftxx.h` next to `.djinni.objc.h`. With the option off, the output does not change.
```sh
protoc \
  -I. \
  --cruxclient-go_out=swift_module=S12Protos:. \
  path/to/*.proto
```

**Option `kotlin_package`**: Turns on the Kotlin output for the pure Kotlin messages of `protoc-gen-crux-kotlin`. A bare `kotlin_package` uses the root package `com.safetyculture.s12kt`. `kotlin_package=<Root>` sets the root package. With the option on, each `.djinni.yaml` record gets `kotlinTranslator`, `kotlinHeader` and `kotlinTypeSignature` keys in its `jni:` block and a `jniTypename` key in its `kotlin:` block, and the plugin writes `.djinni.kotlin.jni.h`. The Kotlin translator moves a message as serialized bytes between C++ and the Kotlin class, with `toByteArray()`, `getSerializedSize()` and `parseFrom(ByteArray)`. With the option off, the output does not change.
```sh
protoc \
  -I. \
  --cruxclient-go_out=kotlin_package:. \
  path/to/*.proto
```

---

## protoc-gen-crux-kotlin

**Purpose**: Generates pure Kotlin message classes with a wire codec. The classes replace the Java lite classes in Kotlin code. They have no dependency except the Kotlin standard library.

**Generates**: one `.pb.kt` file for each proto file, and one runtime file `<root>/wire/S12Wire.kt` for each run.

**Install**:
```sh
make install-crux-kotlin
```

**Usage**:
```sh
protoc \
  -I. \
  --crux-kotlin_out=kotlin_package:. \
  path/to/*.proto
```

**Option `kotlin_package`**: Sets the root package. The default is `com.safetyculture.s12kt`. The proto package follows the root, with the leading `s12.` removed, as the Java lite package does. The Kotlin classes and the Java lite classes can be in one app, because their packages are different.

**Rules**:

| Proto | Kotlin |
|-------|--------|
| message | `data class` with `val` parameters and default values. Use the constructor and `copy()`. |
| plain scalar | Non-null, with the proto default |
| `optional` scalar, singular message | Nullable, `null` by default |
| `repeated` | `List`, empty by default |
| `map` | `Map`, empty by default |
| `oneof` | A nullable property of a `sealed class` type, with one subclass for each member |
| enum | `enum class` with `number`, `fromNumber()` and an `UNRECOGNIZED` entry |

Each message has `toByteArray()`, `getSerializedSize()` and `companion object { @JvmStatic fun parseFrom(bytes: ByteArray) }`. A field that the code does not know is kept in `unknownFields` and written again. The well-known types (`Timestamp`, `Duration`, `FieldMask`, `Any`, `Struct`, wrappers, `Empty`) are generated like other messages, when you give their files to the plugin. They have data and codec only.

The wire format is the same as the format of Java lite. One case is different. A proto3 enum number that the enum does not declare cannot be held by `UNRECOGNIZED`. The codec keeps the number in `unknownFields` and writes it after the known fields. For a map entry, the codec keeps the whole entry there.

**Test**: `go test ./...` in the plugin directory runs the golden test and the nullability tests. `test/wire_compat.sh` generates Kotlin and Java lite classes from the same proto files, encodes random messages and compares the bytes in both directions. Set `APISCHEMA` to an APISchema checkout to include `s12.assets.v1.Asset`.
