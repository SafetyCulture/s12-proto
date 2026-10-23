#!/usr/bin/env bash
# Checks generated glue headers against their pinned native dependencies.
set -euo pipefail

HERE=$(cd "$(dirname "$0")" && pwd)
PLUGIN_DIR=$(cd "$HERE/../.." && pwd)
CACHE_DIR=${CRUXCLIENT_CI_CACHE:-"${TMPDIR:-/tmp}/cruxclient-go-ci"}
PROTOBUF_VERSION=36.1
DJINNI_COMMIT=62d8899174a027f85c8f1e4cb0964fc3ff8b7de2
mkdir -p "$CACHE_DIR"

download() {
  local destination=$1
  local url=$2
  local expected=$3
  if [ ! -f "$destination" ]; then
    curl --fail --location --silent --show-error --output "$destination" "$url"
  fi
  local actual
  actual=$(shasum -a 256 "$destination" | awk '{print $1}')
  if [ "$actual" != "$expected" ]; then
    echo "SHA-256 mismatch for $destination: got $actual, want $expected" >&2
    exit 1
  fi
}

case "$(uname -s)-$(uname -m)" in
  Darwin-arm64)
    protoc_platform=osx-aarch_64
    protoc_sha256=de56d57afe30c5d191b11d24ff93dd4025728d7fb43b773886b2d3613e0bdbb2
    jni_platform=darwin
    ;;
  Linux-x86_64)
    protoc_platform=linux-x86_64
    protoc_sha256=c4bc672d9d49214dc8cafdceadf4df92182d6ca8e3ec65a56b2d7de5602669b4
    jni_platform=linux
    ;;
  *)
    echo "unsupported platform: $(uname -s)-$(uname -m)" >&2
    exit 1
    ;;
esac

protoc_zip="$CACHE_DIR/protoc-$PROTOBUF_VERSION-$protoc_platform.zip"
download "$protoc_zip" "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protoc-$PROTOBUF_VERSION-$protoc_platform.zip" "$protoc_sha256"
protoc_root="$CACHE_DIR/protoc-$PROTOBUF_VERSION-$protoc_platform"
if [ ! -x "$protoc_root/bin/protoc" ]; then
  rm -rf "$protoc_root"
  mkdir -p "$protoc_root"
  unzip -q "$protoc_zip" -d "$protoc_root"
fi
protoc="$protoc_root/bin/protoc"

protobuf_archive="$CACHE_DIR/protobuf-$PROTOBUF_VERSION.tar.gz"
download "$protobuf_archive" "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protobuf-$PROTOBUF_VERSION.tar.gz" dc74fa582f559cbd31614ddfefb4868f43c919d7184bde514bb47f90c6025eb8
protobuf_root="$CACHE_DIR/protobuf-$PROTOBUF_VERSION"
if [ ! -d "$protobuf_root/src/google/protobuf" ]; then
  rm -rf "$protobuf_root"
  mkdir -p "$protobuf_root"
  tar -xzf "$protobuf_archive" --strip-components=1 -C "$protobuf_root"
fi

abseil_archive="$CACHE_DIR/abseil-cpp-20250512.1.tar.gz"
download "$abseil_archive" https://github.com/abseil/abseil-cpp/archive/refs/tags/20250512.1.tar.gz 9b7a064305e9fd94d124ffa6cc358592eb42b5da588fb4e07d09254aa40086db
abseil_root="$CACHE_DIR/abseil-cpp-20250512.1"
if [ ! -d "$abseil_root/absl" ]; then
  rm -rf "$abseil_root"
  mkdir -p "$abseil_root"
  tar -xzf "$abseil_archive" --strip-components=1 -C "$abseil_root"
fi

djinni_root="$CACHE_DIR/snapchat-djinni-$DJINNI_COMMIT"
if [ ! -d "$djinni_root/.git" ]; then
  rm -rf "$djinni_root"
  git init -q "$djinni_root"
  git -C "$djinni_root" remote add origin https://github.com/SafetyCulture/snapchat-djinni.git
  git -C "$djinni_root" fetch -q --depth 1 origin "$DJINNI_COMMIT"
fi
git -C "$djinni_root" checkout -q --detach FETCH_HEAD 2>/dev/null || git -C "$djinni_root" checkout -q --detach "$DJINNI_COMMIT"
[ "$(git -C "$djinni_root" rev-parse HEAD)" = "$DJINNI_COMMIT" ] || { echo "Djinni checkout is not $DJINNI_COMMIT" >&2; exit 1; }

if [ -z "${JAVA_HOME:-}" ]; then
  javac_home=$(cd "$(dirname "$(command -v javac)")/.." && pwd)
  if [ -f "$javac_home/include/jni.h" ]; then
    JAVA_HOME=$javac_home
  elif [ "$(uname -s)" = Darwin ]; then
    JAVA_HOME=$(/usr/libexec/java_home)
  else
    JAVA_HOME=$javac_home
  fi
fi
[ -f "$JAVA_HOME/include/jni.h" ] || { echo "JDK headers not found under JAVA_HOME=$JAVA_HOME" >&2; exit 1; }

output="$CACHE_DIR/generated-$PROTOBUF_VERSION"
rm -rf "$output"
mkdir -p "$output/protos/routeguide/v1" "$output/protos/swiftnames" "$output/protos/shadow" "$output/protos/required"
cp "$PLUGIN_DIR/testdata/routeguide/message.proto" "$output/protos/routeguide/v1/message.proto"
cp "$PLUGIN_DIR/testdata/routeguide/wire_options.proto" "$output/protos/wire_options.proto"
cp "$PLUGIN_DIR/testdata/swiftnames/names.proto" "$output/protos/swiftnames/names.proto"
cp "$PLUGIN_DIR/testdata/swiftnames/loose.proto" "$output/protos/swiftnames/loose.proto"
cp "$PLUGIN_DIR/testdata/kotlin/kitchen/shadow/shadow.proto" "$output/protos/shadow/shadow.proto"
cp "$PLUGIN_DIR/testdata/required/required.proto" "$output/protos/required/required.proto"
protoc_log="$output/protoc.log"
if ! "$protoc" -I"$output/protos" -I"$protoc_root/include" --cpp_out="$output" \
  routeguide/v1/message.proto wire_options.proto swiftnames/names.proto swiftnames/loose.proto \
  shadow/shadow.proto required/required.proto >"$protoc_log" 2>&1; then
  cat "$protoc_log" >&2
  exit 1
fi

headers=()
while IFS= read -r -d '' header; do
  headers+=("$header")
done < <(find "$PLUGIN_DIR/testdata" \( -name '*.djinni.swiftxx.h' -o -name '*.djinni.kotlin.jni.h' \) -print0)
[ "${#headers[@]}" -gt 0 ] || { echo "no generated glue headers found" >&2; exit 1; }

for header in "${headers[@]}"; do
  support_dir="$djinni_root/support-lib/jni"
  if [[ "$header" == *.djinni.swiftxx.h ]]; then
    support_dir="$djinni_root/support-lib/swiftxx"
  fi
  driver="$output/header-check.cc"
  printf '#include "%s"\n' "$header" >"$driver"
  clang++ -std=c++20 -fsyntax-only -Wall -Wextra -Wno-deprecated-declarations "$driver" \
    -I"$output" -I"$protobuf_root/src" -I"$protobuf_root/third_party/utf8_range" -I"$abseil_root" \
    -I"$JAVA_HOME/include" -I"$JAVA_HOME/include/$jni_platform" \
    -I"$support_dir" -I"$djinni_root/support-lib"
  echo "checked ${header#$PLUGIN_DIR/}"
done
