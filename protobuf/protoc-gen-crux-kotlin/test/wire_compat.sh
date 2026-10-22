#!/usr/bin/env bash
# Wire compatibility of the Kotlin messages with Java lite.
#
# Generates both sets of classes from the same proto files, writes random encoded samples for
# every message of the root files, and runs test/WireCompat.kt over them.
#
# Needs: go, protoc, kotlinc, javac and java on PATH.
# Environment:
#   APISCHEMA        An APISchema checkout. When set, s12/assets/v1/asset.proto is a root file.
#   PROTOC_INCLUDE   The directory that holds google/protobuf/*.proto. Default: next to protoc.
#   JAVALITE_JAR     The protobuf-javalite jar. Default: downloaded for the version of protoc.
#   WORK             Scratch directory. Default: a new temporary directory, removed on exit.
set -euo pipefail

HERE=$(cd "$(dirname "$0")" && pwd)
PLUGIN_DIR=$(dirname "$HERE")
REPO=$(cd "$PLUGIN_DIR/../.." && pwd)
CRUXCLIENT=$REPO/protobuf/protoc-gen-cruxclient
CRUXCLIENT_GO=$REPO/protobuf/protoc-gen-cruxclient-go

WORK=${WORK:-$(mktemp -d)}
mkdir -p "$WORK"
[ -n "${KEEP_WORK:-}" ] || trap 'rm -rf "$WORK"' EXIT

PROTOC_INCLUDE=${PROTOC_INCLUDE:-$(dirname "$(dirname "$(command -v protoc)")")/include}
[ -f "$PROTOC_INCLUDE/google/protobuf/descriptor.proto" ] || { echo "PROTOC_INCLUDE has no google/protobuf/descriptor.proto: $PROTOC_INCLUDE" >&2; exit 2; }

if [ -z "${JAVALITE_JAR:-}" ]; then
  # protoc 36.2 pairs with protobuf-javalite 4.36.2.
  version=$(protoc --version | sed 's/^libprotoc //')
  JAVALITE_JAR=$WORK/protobuf-javalite.jar
  curl -sSfL -o "$JAVALITE_JAR" "https://repo1.maven.org/maven2/com/google/protobuf/protobuf-javalite/4.$version/protobuf-javalite-4.$version.jar"
fi

includes=(-I"$PLUGIN_DIR/testdata/kitchen" -I"$CRUXCLIENT/proto" -I"$CRUXCLIENT" -I"$CRUXCLIENT_GO/testdata/swiftnames" -I"$PROTOC_INCLUDE")
roots=(s12/kitchen/v1/kitchen.proto legacy/legacy.proto routeguide/v1/message.proto names.proto loose.proto prefixed.proto)
if [ -n "${APISCHEMA:-}" ]; then
  includes+=(-I"$APISCHEMA/protos" -I"$APISCHEMA/third_party" -I"$REPO")
  roots+=(s12/assets/v1/asset.proto)
else
  echo "APISCHEMA is not set: s12.assets.v1.Asset is not in the corpus" >&2
fi

(cd "$PLUGIN_DIR" && go build -o "$WORK/protoc-gen-crux-kotlin" . && go build -o "$WORK/gensamples" ./test/gensamples)

protoc "${includes[@]}" --experimental_allow_proto3_optional --include_imports --descriptor_set_out="$WORK/corpus.desc" "${roots[@]}"
"$WORK/gensamples" -desc "$WORK/corpus.desc" -roots "$(IFS=,; echo "${roots[*]}")" -out "$WORK/samples"

closure=()
java_files=()
while IFS= read -r f; do
  closure+=("$f")
  case $f in google/protobuf/*) ;; *) java_files+=("$f") ;; esac
done <"$WORK/samples/closure.txt"

mkdir -p "$WORK/kt" "$WORK/java" "$WORK/java-classes"
protoc "${includes[@]}" --experimental_allow_proto3_optional --plugin=protoc-gen-crux-kotlin="$WORK/protoc-gen-crux-kotlin" --crux-kotlin_out="$WORK/kt" "${closure[@]}"
protoc "${includes[@]}" --experimental_allow_proto3_optional --java_out=lite:"$WORK/java" "${java_files[@]}"

javac -nowarn -cp "$JAVALITE_JAR" -d "$WORK/java-classes" $(find "$WORK/java" -name '*.java') 2>&1 | grep -v '^Note:' || true
kotlinc -nowarn -include-runtime -cp "$WORK/java-classes:$JAVALITE_JAR" -d "$WORK/compat.jar" $(find "$WORK/kt" -name '*.kt') "$HERE/WireCompat.kt"
java -cp "$WORK/compat.jar:$WORK/java-classes:$JAVALITE_JAR" WireCompatKt "$WORK/samples"
