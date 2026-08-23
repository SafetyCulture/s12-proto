# C# conformance runner

Runs the generated C# validators and the runtime over the vectors in
[`protobuf/validation/conformance`](../../validation/conformance), and reports
every result that differs from the one Go recorded. Exits non-zero when there is
one.

```bash
mise run test:csharp
```

## What it compiles

| Source | Where it comes from |
| --- | --- |
| `generated/*Validator.g.cs` | this repository's plugin, committed |
| message classes | protoc, at build time via `Grpc.Tools` |
| `S12.Protobuf.Validation` | `../runtime`, by project reference |

The validators are committed because they are this repository's output, and the
plugin's golden test compares its current output against these very files - so
output that has moved on from what is committed fails there rather than reaching
a conformance run as a stale result. Rewrite them with:

```bash
go test ./protobuf/protoc-gen-csharpvalidator/plugin/ -update
```

The message classes are protoc's own output and are generated during the build,
so there is nothing to keep in step.

## What it checks

`tables.jsonl` by running every code point through each replacer, `helpers.jsonl`
by calling each helper directly, and `fields.jsonl` by setting one field at a time
on the committed base message and validating it. The verdict, the error text and
the value left on the field are all compared.

`nested.jsonl` and `deep.jsonl` do the same for a field inside a nested message,
one level down and at every depth of a three-level chain, which is where the
re-rooted error text is checked.
