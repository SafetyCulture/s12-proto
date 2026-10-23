#!/usr/bin/env bash
# Runs the required-field Swift translator against libprotobuf.
set -euo pipefail

HERE=$(cd "$(dirname "$0")" && pwd)
PLUGIN_DIR=$(cd "$HERE/../.." && pwd)
CACHE_DIR=${CRUXCLIENT_CI_CACHE:-"${TMPDIR:-/tmp}/cruxclient-go-ci"}
PROTOC=${PROTOC:-protoc}
CXX=${CXX:-clang++}
DJINNI_COMMIT=bd2e2eee2002df116ef3c5e50843050423693e0d
mkdir -p "$CACHE_DIR"

protoc_version=$($PROTOC --version)
protoc_version=${protoc_version##* }
protobuf_cflags=()
protobuf_libs=()
if command -v pkg-config >/dev/null 2>&1 && pkg-config --exists protobuf; then
  protobuf_version=$(pkg-config --modversion protobuf)
  read -r -a protobuf_cflags <<< "$(pkg-config --cflags protobuf)"
  read -r -a protobuf_libs <<< "$(pkg-config --libs protobuf)"
elif command -v brew >/dev/null 2>&1; then
  protobuf_prefix=$(brew --prefix protobuf)
  protobuf_version=$($protobuf_prefix/bin/protoc --version)
  protobuf_version=${protobuf_version##* }
  protobuf_cflags=(-I"$protobuf_prefix/include" -I"$(brew --prefix abseil)/include")
  protobuf_libs=(-L"$protobuf_prefix/lib" -lprotobuf)
else
  echo "protobuf development files require pkg-config or Homebrew" >&2
  exit 1
fi
[ "$protoc_version" = "$protobuf_version" ] || { echo "protoc $protoc_version does not match libprotobuf $protobuf_version" >&2; exit 1; }

djinni_root="$CACHE_DIR/snapchat-djinni-$DJINNI_COMMIT"
if [ ! -d "$djinni_root/.git" ]; then
  rm -rf "$djinni_root"
  git init -q "$djinni_root"
  git -C "$djinni_root" remote add origin https://github.com/SafetyCulture/snapchat-djinni.git
  git -C "$djinni_root" fetch -q --depth 1 origin "$DJINNI_COMMIT"
fi
git -C "$djinni_root" checkout -q --detach FETCH_HEAD 2>/dev/null || git -C "$djinni_root" checkout -q --detach "$DJINNI_COMMIT"
[ "$(git -C "$djinni_root" rev-parse HEAD)" = "$DJINNI_COMMIT" ] || { echo "Djinni checkout is not $DJINNI_COMMIT" >&2; exit 1; }

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
mkdir -p "$WORK/generated" "$WORK/protos/required"
cp "$PLUGIN_DIR/testdata/required/required.proto" "$WORK/protos/required/required.proto"
"$PROTOC" -I"$WORK/protos" --cpp_out="$WORK/generated" required/required.proto

"$CXX" -O2 -DNDEBUG -std=c++20 \
  "$HERE/required_swiftxx_test.cc" "$WORK/generated/required/required.pb.cc" \
  -I"$PLUGIN_DIR/testdata" -I"$WORK/generated" \
  -I"$djinni_root/support-lib/swiftxx" -I"$djinni_root/support-lib" \
  "${protobuf_cflags[@]}" "${protobuf_libs[@]}" \
  -o "$WORK/required_swiftxx_test"
"$WORK/required_swiftxx_test"
