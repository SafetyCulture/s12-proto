# Validator conformance vectors

What the Go validators do, recorded so a port can be held to it.

Every language asserts against the same committed files. Nothing here runs
another language's toolchain, so the Go half is checked by `go test ./...` and the
C# half by its own build.

| File | One record per | Records |
| --- | --- | --- |
| `testdata/inputs.txt` | input | the shared corpus |
| `testdata/base.b64` | - | the message every field record starts from |
| `testdata/deep_base.b64` | - | the three-level chain the depth records start from |
| `testdata/fields.jsonl` | (string field of `ValTestMessage`, input) | verdict, error text, and the value afterwards |
| `testdata/nested.jsonl` | (string field one level down, input) | verdict and error text, re-rooted under the field holding the message |
| `testdata/deep.jsonl` | (string field at each depth of a three-level chain, input) | the same, at every depth |
| `testdata/helpers.jsonl` | input | what each validation helper returns for it |
| `testdata/tables.jsonl` | substitution | the character each replacer swaps, and what for |

## Why three files

`fields.jsonl` covers the generated code: it is produced by running the real
validators, so it moves when the generator moves.

`helpers.jsonl` covers the hand-written helpers those validators call. A port
reimplements them, and reaching them only through the fields that happen to use
them would leave most of their behaviour unchecked.

`tables.jsonl` is found by running every code point through each replacer, rather
than by reading the tables. A port transcribes those tables by hand, and this is
what says the copy is complete.

## Why nesting is recorded separately

A failure inside a nested message is re-rooted under the field that holds it, and
the result is its own shape rather than the leaf text with a prefix: the path grows
a segment per level while the innermost text stays as it was, naming the leaf. A
port can agree on every leaf field and still render that wrongly, so `nested.jsonl`
covers one level and `deep.jsonl` covers a three-level chain, where losing a segment
or repeating the prefix shows up.

## Why the base message is committed

`fields.jsonl` sets one field at a time on a message that otherwise passes, so the
only field that can fail is the one under test. That message is committed as wire
bytes rather than transcribed into each language, so both start from the same
value and a change to the fixture cannot quietly mean two different things.

## Why each record carries the value afterwards

Several rules rewrite the field before the remaining checks read it: NFC
normalisation, trimming, URL breaking, the two replacers, and private-use-area
stripping. Two implementations can agree on accept or reject and still leave
different values on the message, so the verdict alone would not catch it.

## Inputs with no C# spelling

A Go string holds bytes and a C# string holds UTF-16 code units, so the two do not
describe quite the same set of values. An input whose bytes are not valid UTF-8 has
no C# spelling: decoding it substitutes the replacement character, which is a
different value. The C# runner steps over those inputs and names them in its
output rather than comparing something the two languages never agreed on. Two of
the inputs are in that position today, and they stay in the corpus because the Go
half still checks them.

## Writing an input

One per line, in `testdata/inputs.txt`. A line opening with `b64:` carries a
base64 payload, which is how a value with leading or trailing space, a line break,
or malformed encoding is written. A line opening with `#` is a comment.

## Regenerating

```bash
# review the diff before committing it - a change here is a change in behaviour
go test ./protobuf/protoc-gen-govalidator/valtest/ -update
```

The vectors are rewritten from current behaviour, so a diff on them is the
behaviour change, stated in the form a reviewer can read.

## Checking a language against them

```bash
go test ./...                # Go, through the valtest package
mise run test:csharp         # C#, through protoc-gen-csharpvalidator/conformance
```
