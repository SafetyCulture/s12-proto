# protoc-gen-csharpvalidator

Emits C# field validators from the `validator.*` proto annotations, the C# counterpart of
`protoc-gen-govalidator`. Both read the same rule and character-class definitions in
`../protoc-gen-govalidator/plugin/validation_definitions.go`, so the two languages cannot drift.

## Layout

| Path | What |
| --- | --- |
| `runtime/Validation/*.cs` | The C# the generated validators call, the counterpart of `s12/protobuf/proto/validator_helpers.go` |
| `runtime/embed.go` | Embeds those files so the generator emits them with its output |
| `runtime/Runtime.csproj` | Compiles the runtime here, so a break shows up in this repository |
| `plugin/` | The emitter |

The runtime is emitted rather than published separately. The generator writes calls against it by
name and nothing checks those names at build time in this repository, so shipping both from one
version of this generator is what keeps them in step.

## Runtime notes

Two properties are load-bearing and easy to lose:

- **Every pattern uses `RegexOptions.NonBacktracking`.** The Go patterns come from RE2, which is
  linear-time by construction and has no backtracking to exploit. The default .NET engine
  backtracks, and the RFC 5322 email pattern does not terminate on it in reasonable time.
  `NonBacktracking` rejects exactly the constructs RE2 also forbids, so any RE2 pattern compiles
  under it. It cannot be combined with `RegexOptions.Compiled`.
- **`RegexOptions.Compiled` must stay off** for the package to remain AOT- and trim-compatible.

### Translating a pattern from RE2

`runtime/Validation/ValidatorPatterns.cs` holds the fixed patterns. Four substitutions are needed, and each one is a
correctness fix rather than a style choice:

| RE2 | .NET | Why |
| --- | --- | --- |
| `\x{NNNN}` | `\uNNNN` | .NET has no braced form. |
| `\d` | `[0-9]` | .NET `\d` also matches non-ASCII digits. |
| `^` … `$` | `\A` … `\z` | .NET `$` also matches immediately before a trailing newline, so an anchored pattern would accept a value ending in `\n` that Go rejects. |
| `\p{Latin}` | explicit block ranges | .NET has no Unicode script classes. |

The Latin block union covers ASCII, Latin-1 Supplement, Latin Extended-A and -B, and Latin
Extended Additional. It does not cover IPA Extensions, Latin Extended-C through -G, or the
fullwidth forms, all of which are Latin script to Go. The effect is confined to usernames.

`\b` needs care and has no single substitution: RE2 defines it over ASCII word characters and
.NET over Unicode ones. It is dropped from the URL-rejection pattern, where removing a constraint
only widens the set of values recognised as containing a URL, and hand-written in
`StringMutators.BreakPartialUrls`, where the Unicode reading would rewrite values the reference
implementation leaves alone.

### Characters outside the Basic Multilingual Plane

.NET matches over UTF-16 code units, so `\p{L}` never matches a code point above U+FFFF — each
half of the surrogate pair carries category `Surrogate`. Go matches over code points and accepts
them. Call `CharacterClass.IsMatch` rather than `Regex.IsMatch` for a generated character class:
it takes the categories the pattern admits, removes the code points above U+FFFF that fall in
them, and matches what is left.

### Unicode general categories

The Go generator emits two-letter category shorthands such as `\pSc` into its character classes.
RE2 reads `\pXy` as the one-letter class `\pX` followed by a literal `y`, so a field declaring
one symbol category is checked against every symbol category. The C# side targets the declared
category. It is therefore stricter than Go for those fields, never more permissive.

## Where C# is deliberately stricter

Four narrowings are known and intended. None of them accepts a value the Go validators reject, so
no field becomes more permissive by being validated in C#.

| Case | Go | C# |
| --- | --- | --- |
| A field declaring one symbol category | admits every symbol category | admits the declared category |
| Username outside the Latin blocks listed above, such as `ａｂ` or `ɐb` | valid | invalid |
| Timezone `Local` | resolves to the host's zone | invalid |
| Empty timezone | resolves to UTC | invalid |

The two timezone cases are unreachable through a request that carries a timezone worth honouring:
`Local` names the server's zone rather than the caller's, and an empty value is caught by the
required rule.

## What the generator emits against the runtime

- `CharacterClass.IsMatch(pattern, value, categories)` for a character class, never `Regex.IsMatch`,
  with `categories` expanded from the pattern's `\p{…}` tokens — `\p{L}` becomes all five letter
  categories, `\p{Sc}` becomes `UnicodeCategories.Sc`.
- `Regex` instances built with `RegexOptions.NonBacktracking`, patterns translated by the five rules
  above.
- `ValidationError.Create(field, description)` for a failure and `error.Nest(fieldName)` when the
  failure came from a nested message, giving the interceptor a dotted field path.
- `ValidationLog.Report(...)` for a rule marked `log_only`, which records the violation and lets the
  message through.
