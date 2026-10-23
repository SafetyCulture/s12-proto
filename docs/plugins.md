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

**Status**: Frozen. This plugin does not get new options. The options `swift_module`, `kotlin_package` and `kotlin_messages` are in `protoc-gen-cruxclient-go` only.

**Purpose**: Generates C++ client code, mock servers, and Djinni cross-language bindings (ObjC, JNI) from proto service definitions. Targets the CRUX Engine mobile client SDK.

**Generates**: `.crux.api.h`, `.crux.api.cc`, `.mock.h`, `.djinni.yaml`, `.djinni.objc.h`, `.djinni.jni.h`

**Install** (requires C++17 toolchain, protobuf, and gRPC libraries):
```sh
make install-cruxclient
```

---

## protoc-gen-cruxclient-go

**Purpose**: Go port of `protoc-gen-cruxclient`. Generates identical output using the `bufbuild/protoplugin` framework — no C++ toolchain required. With no option, the output is the output of `protoc-gen-cruxclient` on `master`. The options below are only in this plugin.

**Generates**: `.crux.api.h`, `.crux.api.cc`, `.mock.h`, `.djinni.yaml`, `.djinni.objc.h`, `.djinni.jni.h`, and `.djinni.swiftxx.h` (only with the `swift_module` option), `.djinni.kotlin.jni.h` (only with the `kotlin_package` option), and one `.pb.kt` file for each proto file with a message or an enum, with the runtime file `<root>/wire/S12Wire.kt` (only with the `kotlin_messages` option)

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

**Option `kotlin_package`**: Turns on the Kotlin output for the Djinni bindings of the pure Kotlin messages. Use the option `kotlin_messages` for the Kotlin messages. A bare `kotlin_package` uses the root package `com.safetyculture.s12kt`. `kotlin_package=<Root>` sets the root package. With the option on, each `.djinni.yaml` record gets `kotlinTranslator`, `kotlinHeader` and `kotlinTypeSignature` keys in its `jni:` block and a `jniTypename` key in its `kotlin:` block, and the plugin writes `.djinni.kotlin.jni.h`. The Kotlin translator moves a message as serialized bytes between C++ and the Kotlin class, with `toByteArray()` and `parseFrom(ByteArray)`. A parse or serialize failure throws `std::runtime_error`, which djinni converts to a Java exception. With the option off, the output does not change.
```sh
protoc \
  -I. \
  --cruxclient-go_out=kotlin_package:. \
  path/to/*.proto
```

**Option `kotlin_messages`**: Writes pure Kotlin message classes with a wire codec. The classes replace the Java lite classes in Kotlin code. They have no dependency except the Kotlin standard library. A bare `kotlin_messages` uses the root package from `kotlin_package`, or `com.safetyculture.s12kt` if `kotlin_package` is absent. `kotlin_messages=<Root>` sets the root package. The option is independent of `kotlin_package`. Give both options the same root package, so that the `jniTypename` keys name the generated classes. With the option off, the output does not change. Only this plugin has the option. The C++ plugin cannot write Kotlin.
```sh
protoc \
  -I. \
  --cruxclient-go_out=kotlin_messages:. \
  path/to/*.proto
```
The plugin writes one `.pb.kt` file for each proto file that has a message or an enum, and one runtime file `<root>/wire/S12Wire.kt` for each run. The proto package follows the root, with the leading `s12.` removed, as the Java lite package does. The Kotlin classes and the Java lite classes can be in one app, because their packages are different. The plugin must get all files of a run in one request, because it writes the runtime file once. With `buf`, set `strategy: all`.

**Option `messages_only`**: Turns off all cruxclient output in the run: the API, mock and Djinni files, and the files of `swift_module` and `kotlin_package`. A run with `kotlin_messages,messages_only` writes only the `.pb.kt` files and the runtime file. The option does not take a value. The output of the other options does not change when `messages_only` is absent.

**Rules**:

| Proto | Kotlin |
|-------|--------|
| message | `data class` with `val` parameters and default values. Use the constructor and `copy()`. |
| plain scalar | Non-null, with the proto default |
| `optional` scalar, singular message | Nullable, `null` by default |
| proto2 `required` | Non-null, with the proto default. `parseFrom` throws `InvalidWireException` if a part has no value for it |
| `repeated` | `List`, empty by default |
| `map` | `Map`, empty by default |
| `oneof` | A nullable property of a `sealed class` type, with one subclass for each member |
| enum | `sealed class` with one `object` for each value, `number`, `name`, `entries`, `fromNumber()` and a class `UNRECOGNIZED(number)` |

Each message has `toByteArray()`, `getSerializedSize()` and `companion object { @JvmStatic fun parseFrom(bytes: ByteArray) }`. A field that the code does not know is kept in `unknownFields` and written again. The well-known types (`Timestamp`, `Duration`, `FieldMask`, `Any`, `Struct`, wrappers, `Empty`) are generated like other messages, when you give their files to the plugin. They have data and codec only.

The wire format is the same as the format of Java lite. A proto3 enum is open: a number that the enum does not declare becomes `UNRECOGNIZED(number)` in the property, and the codec writes the number in the same place. A closed enum (proto2, or an edition enum with `enum_type = CLOSED`) does not reject such a number: the codec keeps it in `unknownFields` and writes it again, as Java lite does. `UNRECOGNIZED` has an `internal` constructor. Get it from `fromNumber`, so that `UNRECOGNIZED(-1)` cannot be made for a declared number. The reader stops at a nesting depth of 100 and throws `InvalidWireException`, as Java lite does.

`parseFrom` checks the required fields of a proto2 message once, after it has merged all parts of the message, and throws `InvalidWireException("missing required field <name>")` for the first field that no part has. A message below the top-level message does not throw on its own, because a later part can have the field. A message can have 63 required fields at most. A string field with the feature `utf8_validation = NONE` (proto2 by default) accepts invalid UTF-8 and replaces the bad bytes with U+FFFD, as Java lite does. A proto3 string, an edition string with `utf8_validation = VERIFY`, and every map entry throw `InvalidWireException` for invalid UTF-8.

A parse returns lists and maps that cannot change. A constructor keeps the list or the map that the caller passes, and the serialized size is cached, so a caller must not change that list or map later. If the caller does, `toByteArray()` throws `IllegalStateException`. An extension (`extend`) is not generated. The bytes of an extension field are kept in `unknownFields` and written again. The codec writes an unknown field with the tags, lengths and varints in their shortest form, as Java lite does. A string with an unpaired surrogate is written with `?` in place of the surrogate, where Java lite throws. The codec throws `InvalidWireException`, a `RuntimeException`, for every wire error. A varint of more than ten bytes is a wire error. Java lite accepts an int32 or an int64 varint of more than ten bytes when it reads from a byte array, and refuses it when it reads from a stream. The codec refuses it in both cases. The plugin does not validate the values of its options.

**Test**: `go test ./...` in `protobuf/protoc-gen-cruxclient-go` runs the golden test and the nullability tests. `test/wire_compat.sh` generates Kotlin and Java lite classes from the same proto files, encodes random messages and compares the bytes in both directions. Set `APISCHEMA` to an APISchema checkout to include `s12.assets.v1.Asset`. Set `ROOTS_FILE` to a list of proto files, relative to `APISCHEMA/protos`, to include more files, such as the whole Crux whitelist.
