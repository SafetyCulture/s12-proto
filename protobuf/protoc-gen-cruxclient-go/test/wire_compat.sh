#!/usr/bin/env bash
# Wire compatibility of the Kotlin messages with Java lite.
#
# Generates both sets of classes from the same proto files, writes random encoded samples for
# every message of the root files, and runs test/WireCompat.kt over them.
#
# Needs: go, protoc, kotlinc, javac, java, curl and shasum on PATH.
# Environment:
#   APISCHEMA        An APISchema checkout. When set, s12/assets/v1/asset.proto is a root file.
#   ROOTS_FILE       A file with one more root proto file in each line, relative to APISCHEMA/protos.
#   SAMPLES          Samples for each message. Default: 12.
#   PROTOC_INCLUDE   The directory that holds google/protobuf/*.proto. Default: next to protoc.
#   JAVALITE_JAR     The protobuf-javalite jar. Default: downloaded, at JAVALITE_VERSION.
#   JAVALITE_VERSION The protobuf-javalite version to download. Default: 4.36.2.
#   JAVALITE_SHA256  The SHA-256 of the downloaded jar. Default: the digest of the default version.
#   WORK             Scratch directory. A directory that you give is kept. Default: a new temporary directory, removed on exit.
#   KEEP_WORK        Set it to keep the default temporary directory too.
set -euo pipefail

HERE=$(cd "$(dirname "$0")" && pwd)
PLUGIN_DIR=$(dirname "$HERE")
REPO=$(cd "$PLUGIN_DIR/../.." && pwd)
CRUXCLIENT=$REPO/protobuf/protoc-gen-cruxclient

if [ -z "${WORK:-}" ]; then
  WORK=$(mktemp -d)
  [ -n "${KEEP_WORK:-}" ] || trap 'rm -rf "$WORK"' EXIT
fi
mkdir -p "$WORK"

PROTOC_INCLUDE=${PROTOC_INCLUDE:-$(dirname "$(dirname "$(command -v protoc)")")/include}
[ -f "$PROTOC_INCLUDE/google/protobuf/descriptor.proto" ] || { echo "PROTOC_INCLUDE has no google/protobuf/descriptor.proto: $PROTOC_INCLUDE" >&2; exit 2; }

if [ -z "${JAVALITE_JAR:-}" ]; then
  JAVALITE_VERSION=${JAVALITE_VERSION:-4.36.2}
  JAVALITE_SHA256=${JAVALITE_SHA256:-8f4f93d08cb37130ecb0bc69d376d7d27c36ff04ddb1bb27430b06caa9fa062e}
  JAVALITE_JAR=$WORK/protobuf-javalite.jar
  curl -sSfL -o "$JAVALITE_JAR" "https://repo1.maven.org/maven2/com/google/protobuf/protobuf-javalite/$JAVALITE_VERSION/protobuf-javalite-$JAVALITE_VERSION.jar"
  actual=$(shasum -a 256 "$JAVALITE_JAR" | cut -d' ' -f1)
  [ "$actual" = "$JAVALITE_SHA256" ] || { echo "protobuf-javalite $JAVALITE_VERSION has SHA-256 $actual, not $JAVALITE_SHA256" >&2; exit 2; }
fi

includes=(-I"$PLUGIN_DIR/testdata/kotlin/kitchen" -I"$CRUXCLIENT/proto" -I"$CRUXCLIENT" -I"$PLUGIN_DIR/testdata/swiftnames" -I"$PROTOC_INCLUDE")
roots=(s12/kitchen/v1/kitchen.proto legacy/legacy.proto editions/closed.proto routeguide/v1/message.proto names.proto loose.proto prefixed.proto)
if [ -n "${APISCHEMA:-}" ]; then
  includes+=(-I"$APISCHEMA/protos" -I"$APISCHEMA/third_party" -I"$REPO")
  roots+=(s12/assets/v1/asset.proto)
  if [ -n "${ROOTS_FILE:-}" ]; then
    while IFS= read -r f || [ -n "$f" ]; do
      if [ -n "$f" ] && [ "$f" != s12/assets/v1/asset.proto ]; then roots+=("$f"); fi
    done <"$ROOTS_FILE"
  fi
else
  echo "APISCHEMA is not set: s12.assets.v1.Asset is not in the corpus" >&2
fi

(cd "$PLUGIN_DIR" && go build -o "$WORK/protoc-gen-cruxclient-go" . && go build -o "$WORK/gensamples" ./test/gensamples)

protoc "${includes[@]}" --experimental_allow_proto3_optional --include_imports --descriptor_set_out="$WORK/corpus.desc" "${roots[@]}"
"$WORK/gensamples" -desc "$WORK/corpus.desc" -roots "$(IFS=,; echo "${roots[*]}")" -n "${SAMPLES:-12}" -out "$WORK/samples"

closure=()
java_files=()
while IFS= read -r f; do
  closure+=("$f")
  case $f in google/protobuf/*) ;; *) java_files+=("$f") ;; esac
done <"$WORK/samples/closure.txt"

mkdir -p "$WORK/kt" "$WORK/java" "$WORK/java-classes"
protoc "${includes[@]}" --experimental_allow_proto3_optional --plugin=protoc-gen-cruxclient-go="$WORK/protoc-gen-cruxclient-go" --cruxclient-go_out=kotlin_messages:"$WORK/kt" "${closure[@]}"
protoc "${includes[@]}" --experimental_allow_proto3_optional --java_out=lite:"$WORK/java" "${java_files[@]}"

find "$WORK/java" -name '*.java' >"$WORK/java-files.txt"
javac -nowarn -cp "$JAVALITE_JAR" -d "$WORK/java-classes" @"$WORK/java-files.txt" >"$WORK/javac.log" 2>&1 || { cat "$WORK/javac.log" >&2; exit 1; }
JAVA_OPTS=${JAVA_OPTS:--Xmx6g} kotlinc -nowarn -include-runtime -cp "$WORK/java-classes:$JAVALITE_JAR" -d "$WORK/compat.jar" "$WORK/kt" "$HERE/WireCompat.kt" "$HERE/FieldDump.kt"
java -Xmx4g -cp "$WORK/compat.jar:$WORK/java-classes:$JAVALITE_JAR" WireCompatKt "$WORK/samples"
