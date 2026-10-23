#!/usr/bin/env bash
# Makes the inputs of TestControlGolden: the whole schema without Crux prefixes, and the output
# of protoc-gen-swift 1.38.1 for it. The output is what the APISchema Swift lane writes
# (make crux-swift), so the only thing that differs between the two sides is the plugin.
# The files are too big for the repository, so the script writes them to a directory.
#
# Usage:
#   APISCHEMA=<APISchema checkout> testdata/regen-control.sh <out-dir>
#
# Then run the test:
#   CRUX_SWIFT_CONTROL_IMAGE=<out-dir>/image.binpb CRUX_SWIFT_CONTROL_DIR=<out-dir>/swift go test ./internal/swift -run TestControlGolden -v
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
: "${APISCHEMA:?set APISCHEMA to a checkout of APISchema}"
OUT="${1:?usage: regen-control.sh <out-dir>}"
mkdir -p "$OUT"
OUT="$(cd "$OUT" && pwd -P)"
WORK="$(mktemp -d "${TMPDIR:-/tmp}/crux-swift-control.XXXXXX")"
trap 'rm -rf "$WORK"' EXIT
. "$HERE/tools.sh"

BUF="$(find_buf)"
PLUGIN="$(find_protoc_gen_swift)"
export BUF_CACHE_DIR="$WORK/bufcache"
echo "protoc-gen-swift: $("$PLUGIN" --version); buf: $("$BUF" --version); APISchema: $(git -C "$APISCHEMA" rev-parse HEAD)"

cat > "$WORK/swift.yaml" <<YAML
version: v1
plugins:
  - name: swift
    path: $PLUGIN
    out: ./swift
    opt: Visibility=Public,FileNaming=PathToUnderscores
YAML
cd "$APISCHEMA"
"$BUF" build "${BUF_EXCLUDES[@]}" -o "$OUT/image.binpb"
rm -rf "$OUT/swift"
"$BUF" generate "${BUF_EXCLUDES[@]}" --exclude-path third_party/google/protobuf --template "$WORK/swift.yaml" -o "$OUT"
echo "$(ls "$OUT/swift" | wc -l) files in $OUT/swift"
