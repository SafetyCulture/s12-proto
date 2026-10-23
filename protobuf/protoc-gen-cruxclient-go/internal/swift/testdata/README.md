# Test data of the Swift messages

The test of `protoc-gen-cruxclient-go` compares the Swift output with the output of
`protoc-gen-swift`. This directory has the inputs and the expected output.

| Path | Content |
|------|---------|
| `whitelist.binpb.gz` | The s12 schema image of the Crux whitelist, with the message prefixes of Crux, and the files that it imports |
| `whitelist.txt` | The 530 proto files of the whitelist. The test writes the Swift file of each one |
| `golden/` | The output of `protoc-gen-swift` for those files, with `Visibility=Public,FileNaming=PathToUnderscores` |
| `kitchen/` | Protos with the constructs that the whitelist does not use (groups, extensions, message sets, defaults, oneofs, editions). `kitchen/variants.txt` lists the options, and `kitchen/golden/<name>/` has the output for each line |
| `large/` | A proto with an enum and a oneof that have more cases than the limit of one switch statement, and its golden file. `regen.sh` writes the proto |
| `upstream/` | Protos of swift-protobuf 1.38.1 that test the naming of types, fields and cases, reserved words, deprecation and groups, copied without a change (license in `upstream/LICENSE.txt`), and the golden files for the default options |
| `GOLDEN_VERSION` | The version line of the `protoc-gen-swift` that wrote the golden files. A test checks it |
| `../../../testdata/swiftmessages/` | The output of `protoc-gen-swift` for the route guide message proto, which the test of the plugin handler uses |

The golden files are in the repository, because CI has no Swift toolchain: the tests compare the output of the plugin with files that `protoc-gen-swift` wrote before. `.gitattributes` marks them as generated and keeps their line endings, because the tests compare bytes. Do not move them to Git LFS: `go test` and `check-off-identity.sh` read them from the checkout.

## Versions of the inputs

| Input | Version |
|-------|---------|
| `protoc-gen-swift` | swift-protobuf tag `1.38.1`, commit `55d7a1cc5666b85c13464aea1c4b4a90feccb4c8`. The iOS app links the same version. |
| `buf` | 1.66.0 |
| `protoc` (kitchen files) | 36.2 |
| APISchema | `868fca1a1f6cd78b0c31ac9ed764c043b22ec919` (`main`, 2026-10-23) |
| mitti-mobile | `49159543207f31dd13e4a20ce3e8196ac34a9fc6` (`tools/lib/s12_prefixes.py`, `crux/deps/s12-prefixes.yaml`, `crux/tools/s12_modules_*.txt`) |

## Scripts

- `regen.sh` writes `whitelist.binpb.gz`, `whitelist.txt`, `golden/`, `kitchen/golden/`, `large/`, `upstream/golden/`, `GOLDEN_VERSION` and `../../../testdata/swiftmessages/` again. It stops if the program is not `protoc-gen-swift` 1.38.1, so the golden files cannot come from the plugin under test. Set `APISCHEMA`, `MITTI_MOBILE` and, if you have a binary, `PROTOC_GEN_SWIFT`. Run it again after you change the version of `protoc-gen-swift`, the prefix table or the whitelist, and update the versions above.
- `regen-control.sh <out-dir>` writes the control corpus: the 1490 files of the whole schema without prefixes, which is the output of `make crux-swift` in APISchema. It is too big for the repository. `TestControlGolden` runs on it when `CRUX_SWIFT_CONTROL_IMAGE` and `CRUX_SWIFT_CONTROL_DIR` are set.
- `check-off-identity.sh <ref>` builds the plugin at `<ref>` and at the working tree. It runs both over the whitelist image with the options off (including `kotlin_messages`) and compares the output trees.
