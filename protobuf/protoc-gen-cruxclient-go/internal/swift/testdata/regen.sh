#!/usr/bin/env bash
# Regenerates the golden files of the Swift output in this directory:
#
#   whitelist.binpb.gz   the s12 schema image of the Crux whitelist, with the message prefixes of Crux
#   whitelist.txt        the proto files of the whitelist, which the test generates
#   golden/              the output of protoc-gen-swift 1.38.1 for those files
#   kitchen/golden/      the output of protoc-gen-swift 1.38.1 for the protos in kitchen/, one directory for each line of kitchen/variants.txt
#   large/               the proto with an enum and a oneof over the limit of one switch statement, and its golden file
#   upstream/golden/     the output of protoc-gen-swift 1.38.1 for the protos of swift-protobuf in upstream/
#   GOLDEN_VERSION       the version line of the protoc-gen-swift that wrote all golden files
#   ../../../testdata/swiftmessages/  the output of protoc-gen-swift 1.38.1 for the route guide message proto,
#                        which the test of the plugin handler compares
#
# Usage:
#   APISCHEMA=<APISchema checkout> MITTI_MOBILE=<mitti-mobile checkout> testdata/regen.sh
#
# APISCHEMA is a checkout of SafetyCulture/APISchema. MITTI_MOBILE is a checkout of
# mitti-mobile, which has the prefix table (crux/deps/s12-prefixes.yaml), the whitelist
# (crux/tools/s12_modules_*.txt) and the tool that sets the prefixes (tools/lib/s12_prefixes.py).
# The whitelist image is what the Crux codegen gives to every plugin, and the golden
# files use the options of the APISchema Swift lane: Visibility=Public,FileNaming=PathToUnderscores.
# Set PROTOC_GEN_SWIFT to a protoc-gen-swift 1.38.1 binary. Without it the script builds
# the tag, which needs the Swift toolchain. protoc is read from PATH.
#
# After the run, record the inputs in README.md: the commit of APISchema and of mitti-mobile.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
: "${APISCHEMA:?set APISCHEMA to a checkout of APISchema}"
: "${MITTI_MOBILE:?set MITTI_MOBILE to a checkout of mitti-mobile}"
WORK="$(mktemp -d "${TMPDIR:-/tmp}/crux-swift-regen.XXXXXX")"
trap 'rm -rf "$WORK"' EXIT
. "$HERE/tools.sh"

BUF="$(find_buf)"
PLUGIN="$(find_protoc_gen_swift)"
export BUF_CACHE_DIR="$WORK/bufcache"
echo "protoc-gen-swift: $("$PLUGIN" --version); buf: $("$BUF" --version); APISchema: $(git -C "$APISCHEMA" rev-parse HEAD 2>/dev/null || echo unknown); mitti-mobile: $(git -C "$MITTI_MOBILE" rev-parse HEAD 2>/dev/null || echo unknown)"

cd "$APISCHEMA"
"$BUF" build "${BUF_EXCLUDES[@]}" -o "$WORK/image.binpb"
python3 "$MITTI_MOBILE/tools/lib/s12_prefixes.py" "$MITTI_MOBILE/crux/deps/s12-prefixes.yaml" "$MITTI_MOBILE/crux/tools" \
  "$WORK/image.binpb" "$WORK/prefixed.binpb" > "$WORK/files.txt"
LC_ALL=C sort "$WORK/files.txt" > "$HERE/whitelist.txt"

paths=()
while IFS= read -r file; do paths+=(--path "$file"); done < "$HERE/whitelist.txt"
"$BUF" build "$WORK/prefixed.binpb" "${paths[@]}" -o "$WORK/whitelist.binpb"
gzip -9 -n -c "$WORK/whitelist.binpb" > "$HERE/whitelist.binpb.gz"

cat > "$WORK/swift.yaml" <<YAML
version: v1
plugins:
  - name: swift
    path: $PLUGIN
    out: ./swift
    opt: Visibility=Public,FileNaming=PathToUnderscores
YAML
rm -rf "$HERE/golden" "$WORK/out"
"$BUF" generate "$WORK/prefixed.binpb" "${paths[@]}" --template "$WORK/swift.yaml" -o "$WORK/out"
mkdir -p "$HERE/golden"
cp "$WORK/out/swift/"*.pb.swift "$HERE/golden/"
echo "$(ls "$HERE/golden" | wc -l) golden files for $(wc -l < "$HERE/whitelist.txt" | tr -d ' ') proto files"

# The version line that the tests check, so that a golden file cannot come from another program.
"$PLUGIN" --version > "$HERE/GOLDEN_VERSION"

cd "$HERE/kitchen"
rm -rf golden
while IFS='|' read -r name options; do
  case "$name" in ''|'#'*) continue ;; esac
  mkdir -p "golden/$name"
  protoc -I . --plugin=protoc-gen-swift="$PLUGIN" --swift_out="$options":"golden/$name" kitchen/*.proto
done < variants.txt

cd "$HERE/large"
{
  echo '// Written by regen.sh. An enum and a oneof with more cases than the limit of one switch statement.'
  echo 'syntax = "proto3";'
  echo
  echo 'package kitchen.large;'
  echo
  echo 'enum Big {'
  echo '  BIG_UNSPECIFIED = 0;'
  for i in $(seq 1 519); do echo "  BIG_V$i = $i;"; done
  echo '}'
  echo
  echo 'message Wide {'
  echo '  oneof pick {'
  for i in $(seq 1 510); do echo "    string f$i = $i;"; done
  echo '  }'
  echo '}'
} > large.proto
rm -rf golden
mkdir golden
protoc -I . --plugin=protoc-gen-swift="$PLUGIN" --swift_out=Visibility=Public,FileNaming=PathToUnderscores:golden large.proto

cd "$HERE/upstream"
rm -rf golden
mkdir golden
protoc -I . --plugin=protoc-gen-swift="$PLUGIN" --swift_out=golden *.proto

# The route guide proto of the test of the plugin handler.
MODULE_TESTDATA="$HERE/../../../testdata"
mkdir -p "$WORK/routeguide/routeguide/v1" "$MODULE_TESTDATA/swiftmessages"
cp "$MODULE_TESTDATA/routeguide/message.proto" "$WORK/routeguide/routeguide/v1/message.proto"
cp "$MODULE_TESTDATA/routeguide/wire_options.proto" "$WORK/routeguide/wire_options.proto"
rm -f "$MODULE_TESTDATA/swiftmessages/"*.pb.swift
protoc -I "$WORK/routeguide" --plugin=protoc-gen-swift="$PLUGIN" \
  --swift_out=Visibility=Public,FileNaming=PathToUnderscores:"$MODULE_TESTDATA/swiftmessages" routeguide/v1/message.proto
