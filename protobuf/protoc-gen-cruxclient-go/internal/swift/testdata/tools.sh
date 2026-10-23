# Shared by regen.sh and regen-control.sh, which source it. It has no action of its own.
#
# Finds or builds the tools. Versions are pinned here:
#   protoc-gen-swift  swift-protobuf commit 55d7a1cc5666b85c13464aea1c4b4a90feccb4c8
#   buf               1.66.0, the version of the Crux codegen toolchain
#   protoc            any release that reads the editions test file (checked with 36.2)

SWIFT_PROTOBUF_TAG=1.38.1
SWIFT_PROTOBUF_COMMIT=55d7a1cc5666b85c13464aea1c4b4a90feccb4c8
BUF_VERSION=1.66.0

die() { echo "regen: $*" >&2; exit 1; }

have() { command -v "$1" >/dev/null 2>&1; }

# find_buf prints the path of buf, downloading the pinned release into $WORK when buf is not on PATH.
find_buf() {
  if [ -n "${BUF:-}" ]; then echo "$BUF"; return; fi
  if have buf && [ "$(buf --version)" = "$BUF_VERSION" ]; then command -v buf; return; fi
  case "$(uname -s)-$(uname -m)" in
    Darwin-arm64)
      os=Darwin-arm64
      buf_sha256=fada70b88da222684e9afaf0f9e407a673a7d61808fa93c7ffabc8cc11b9ec67
      ;;
    Linux-x86_64)
      os=Linux-x86_64
      buf_sha256=32e8e7b236e1b9da4eb20ad3c7701a404f7909b18d60f6d00837c63b31d4e6cf
      ;;
    *) die "set BUF to a buf $BUF_VERSION binary" ;;
  esac
  curl -sSfL -o "$WORK/buf" "https://github.com/bufbuild/buf/releases/download/v$BUF_VERSION/buf-$os" || die "download of buf failed"
  actual_sha256="$(openssl dgst -sha256 -r "$WORK/buf" | cut -d ' ' -f1)"
  [ "$actual_sha256" = "$buf_sha256" ] || die "download of buf has an unexpected SHA-256"
  chmod +x "$WORK/buf"
  echo "$WORK/buf"
}

# find_protoc_gen_swift prints the path of protoc-gen-swift 1.38.1. It builds the tag when PROTOC_GEN_SWIFT is not set.
find_protoc_gen_swift() {
  if [ -z "${PROTOC_GEN_SWIFT:-}" ]; then
    have swift || die "set PROTOC_GEN_SWIFT, or install the Swift toolchain so that the script can build protoc-gen-swift $SWIFT_PROTOBUF_TAG"
    git clone -q https://github.com/apple/swift-protobuf "$WORK/swift-protobuf" 2>/dev/null || die "clone of swift-protobuf failed"
    git -C "$WORK/swift-protobuf" checkout -q "$SWIFT_PROTOBUF_COMMIT" || die "checkout of swift-protobuf failed"
    [ "$(git -C "$WORK/swift-protobuf" rev-parse HEAD)" = "$SWIFT_PROTOBUF_COMMIT" ] || die "swift-protobuf is not $SWIFT_PROTOBUF_COMMIT"
    swift build -c release --package-path "$WORK/swift-protobuf" --product protoc-gen-swift >&2 || die "build of protoc-gen-swift failed"
    PROTOC_GEN_SWIFT="$WORK/swift-protobuf/.build/release/protoc-gen-swift"
  fi
  [ "$("$PROTOC_GEN_SWIFT" --version)" = "protoc-gen-swift $SWIFT_PROTOBUF_TAG" ] \
    || die "$PROTOC_GEN_SWIFT is not protoc-gen-swift $SWIFT_PROTOBUF_TAG"
  echo "$PROTOC_GEN_SWIFT"
}

# The paths that APISchema make/generate/_shared.mk leaves out of every lane.
BUF_EXCLUDES=(
  --exclude-path protos/s12/thepub
  --exclude-path protos/s12/workato
  --exclude-path protos/s12/common/health.proto
  --exclude-path third_party/grpc/health/v1/health.proto
)
