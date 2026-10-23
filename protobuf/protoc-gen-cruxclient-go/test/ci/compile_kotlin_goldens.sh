#!/usr/bin/env bash
# Compiles the generated Kotlin corpus with the emitted wire runtime.
set -euo pipefail

HERE=$(cd "$(dirname "$0")" && pwd)
PLUGIN_DIR=$(cd "$HERE/../.." && pwd)
KOTLINC=${KOTLINC:-kotlinc}
PROTOC=${PROTOC:-protoc}
PROTOC_INCLUDE=${PROTOC_INCLUDE:-$(dirname "$(dirname "$(command -v "$PROTOC")")")/include}
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT

[ -f "$PROTOC_INCLUDE/google/protobuf/any.proto" ] || { echo "PROTOC_INCLUDE has no well-known types: $PROTOC_INCLUDE" >&2; exit 1; }
mkdir -p "$WORK/well-known"
(cd "$PLUGIN_DIR" && go build -o "$WORK/protoc-gen-cruxclient-go" .)
"$PROTOC" -I"$PROTOC_INCLUDE" --plugin=protoc-gen-cruxclient-go="$WORK/protoc-gen-cruxclient-go" \
  --cruxclient-go_out=kotlin_messages:"$WORK/well-known" \
  google/protobuf/any.proto google/protobuf/duration.proto google/protobuf/field_mask.proto \
  google/protobuf/struct.proto google/protobuf/timestamp.proto google/protobuf/wrappers.proto

sources=("$PLUGIN_DIR/internal/kotlin/runtime/S12Wire.kt")
while IFS= read -r -d '' source; do
  sources+=("$source")
done < <(find "$PLUGIN_DIR/testdata/kotlin/golden" "$WORK/well-known" -name '*.kt' ! -name S12Wire.kt -print0)

"$KOTLINC" -nowarn -d "$WORK/classes" "${sources[@]}"
echo "compiled ${#sources[@]} Kotlin source files"
