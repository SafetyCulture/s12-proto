# Test data of the Swift messages

The test of `protoc-gen-cruxclient-go` compares the Swift output with the output of
`protoc-gen-swift`. This directory has the inputs and the expected output.

| Path | Content |
|------|---------|
| `whitelist.binpb.gz` | The s12 schema image of the Crux whitelist, with the message prefixes of Crux, and the files that it imports |
| `whitelist.txt` | The 530 proto files of the whitelist. The test writes the Swift file of each one |
| `golden/` | The output of `protoc-gen-swift` for those files, with `Visibility=Public,FileNaming=PathToUnderscores` |
| `kitchen/` | Protos with the constructs that the whitelist does not use, and their golden files in `kitchen/golden/` |

The golden files are in the repository, because the set has 21 MB, which is under the 50 MB limit.

## Versions of the inputs

| Input | Version |
|-------|---------|
| `protoc-gen-swift` | swift-protobuf tag `1.38.1`, commit `55d7a1cc5666b85c13464aea1c4b4a90feccb4c8`. The iOS app links the same version. |
| `buf` | 1.66.0 |
| `protoc` (kitchen files) | 36.2 |
| APISchema | `868fca1a1f6cd78b0c31ac9ed764c043b22ec919` (`main`, 2026-10-23) |
| mitti-mobile | `49159543207f31dd13e4a20ce3e8196ac34a9fc6` (`tools/lib/s12_prefixes.py`, `crux/deps/s12-prefixes.yaml`, `crux/tools/s12_modules_*.txt`) |

## Scripts

- `regen.sh` writes `whitelist.binpb.gz`, `whitelist.txt`, `golden/` and `kitchen/golden/` again. Set `APISCHEMA`, `MITTI_MOBILE` and, if you have a binary, `PROTOC_GEN_SWIFT`. Run it again after you change the version of `protoc-gen-swift`, the prefix table or the whitelist, and update the versions above.
- `regen-control.sh <out-dir>` writes the control corpus: the 1490 files of the whole schema without prefixes, which is the output of `make crux-swift` in APISchema. It is too big for the repository. `TestControlGolden` runs on it when `CRUX_SWIFT_CONTROL_IMAGE` and `CRUX_SWIFT_CONTROL_DIR` are set.
- `check-off-identity.sh <ref>` builds the plugin at `<ref>` and at the working tree. It runs both over the whitelist image with the options off and compares the output trees.
