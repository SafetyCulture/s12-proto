# Validator conformance vectors

What the Go validators do, recorded so a port can be held to it.

Every language asserts against the same committed files. Nothing here runs
another language's toolchain, so the Go half is checked by `go test ./...` and the
C# half by its own build.

| File | One record per | Records |
| --- | --- | --- |
| `testdata/inputs.txt` | input | the shared corpus |
| `testdata/fields.jsonl` | (string field of `ValTestMessage`, input) | verdict, error text, and the value afterwards |
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

## Why each record carries the value afterwards

Several rules rewrite the field before the remaining checks read it: NFC
normalisation, trimming, URL breaking, the two replacers, and private-use-area
stripping. Two implementations can agree on accept or reject and still leave
different values on the message, so the verdict alone would not catch it.

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
