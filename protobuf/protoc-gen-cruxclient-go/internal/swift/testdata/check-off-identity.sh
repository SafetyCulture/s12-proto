#!/usr/bin/env bash
# Checks that the plugin writes the same files as an earlier version when the
# messages options are off. It builds the plugin at <base-ref> and at the
# working tree, runs both over the whitelist image with several option strings, and
# compares the output trees with diff -r. An empty diff for every option string is a pass.
#
# Usage:
#   testdata/check-off-identity.sh <base-ref>
#
# Needs git, go and buf on PATH, or BUF set to a buf binary.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
MODULE="$(cd "$HERE/../../.." && pwd -P)"
REPO="$(git -C "$MODULE" rev-parse --show-toplevel)"
BASE="${1:?usage: check-off-identity.sh <base-ref>}"
WORK="$(mktemp -d "${TMPDIR:-/tmp}/crux-swift-identity.XXXXXX")"
trap 'rm -rf "$WORK"' EXIT
. "$HERE/tools.sh"

BUF="$(find_buf)"
export BUF_CACHE_DIR="$WORK/bufcache"
gunzip -c "$HERE/whitelist.binpb.gz" > "$WORK/image.binpb"

mkdir -p "$WORK/base"
git -C "$REPO" archive "$BASE" "${MODULE#"$REPO"/}" | tar -x -C "$WORK/base"
(cd "$WORK/base/${MODULE#"$REPO"/}" && go build -o "$WORK/plugin-base" .)
(cd "$MODULE" && go build -o "$WORK/plugin-new" .)

status=0
for opt in "" "swift_module=S12Protos" "swift_module,kotlin_package"; do
  for version in base new; do
    cat > "$WORK/template-$version.yaml" <<YAML
version: v1
plugins:
  - name: cruxclient
    path: $WORK/plugin-$version
    out: ./out
    opt: "$opt"
YAML
    rm -rf "$WORK/out-$version"
    mkdir "$WORK/out-$version"
    "$BUF" generate "$WORK/image.binpb" --template "$WORK/template-$version.yaml" -o "$WORK/out-$version"
  done
  echo "opt='$opt': $(find "$WORK/out-base" -type f | wc -l | tr -d ' ') files; diff -r:"
  if diff -r "$WORK/out-base" "$WORK/out-new"; then echo "(empty)"; else status=1; fi
done
exit "$status"
