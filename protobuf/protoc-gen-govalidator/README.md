# protocol buffer validators compiler

A `protoc` plugin that generates `Validate() error` functions on Go proto `struct`s based on field options inside `.proto` files. The validation functions are code-generated and thus don't suffer on performance from tag-based reflection on deeply-nested messages.

## Usage

```
$ go get github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator
$ protoc -I. --gogo_out=:. --govalidator_out=. example.proto
```

## Development

Use these commands at the root folder of this repository for testing:
```bash
cd s12-proto

# regenerate examples and run tests
make govalidator
make govalidator-test

# regenerate valtest and run tests
make govalidator-valtest
make govalidator-valtest-test
```

## Table of Contents
- [Email validation: validator.email](#validator.email)
- [ID validation: validator.id](#validator.id)
- [String validation: validator.string and validator.unsafe_string](#validator.string)
- [Enum validation: validator.enum_required](#validator.enum_required)
- [URL validation: validator.url](#validator.url)
- [Testing](#validator.testing)
- [Legacy Validator Fields](#validator.legacy)

&nbsp;

## Email validation: validator.email (EmailRules) <a name="validator.email"></a>
Validate email address format against RFC 5322.

| Option       | Type | Default | Description |
|--------------|------|---------|-------------|
| optional     | bool | false   | Set this as an optional field. It will allow the string to be empty, validation is skipped in that case. |

Example usage:
```
    string email = 1 [(validator.email) = {}];

string opt_email = 2 [(validator.email) = { optional: true }];
```

&nbsp;

## ID validation: validator.id (IdRules) <a name="validator.id"></a>
Validate ID format against UUIDv4, legacy ID format or S12 SafetyCulture ID `prefix_id` format.

| Option       | Type   | Default | Description |
|--------------|--------|---------|-------------|
| optional     | bool   | false | Set this as an optional field. It will allow the string to be empty, validation is skipped in that case. |
| version      | string | "v4"  | The UUID version, only v4 is supported currently. |
| legacy       | bool   | false  | Also allow legacy_id format, similar to `validator.legacy_id`. UUID validation will be attempted first, if it fails it will fall back to `IsLegacyID` method. This option can be combined with the `s12id` option but only enabling `legacy` will not also accept `s12` ids. |
| s12id       | bool   | false  | Also allow S12 id format with prefixes (e.g. `template_fffaaaccc33340739b1e67ca70f4cf4d`). UUID validation will be attempted first, if it fails it will fall back to `IsS12ID` method. This option can be combined with the `legacy` option but only enabling `s12id` will not also accept `legacy` ids. |

Example usage:
```
          string id = 1 [(validator.id) = {}];

repeated string ids = 2 [(validator.id) = {}];

   string legacy_id = 3 [(validator.id) = { legacy: true }];
```

&nbsp;
## String validation: validator.string and validator.unsafe_string (StringRules) <a name="validator.string"></a>
Validate string fields that do not represent a common or predefined format. This should not be used if there is a better suitable validation method for common formats like email address, UUID, URLs, and the like.

&nbsp;

| Option             | Type   | Default | Description |
|--------------------|--------|---------|-------------|
| optional           | bool   | false   | Set this as an optional field. It will allow the string to be empty, validation is skipped in that case. |
| trim               | bool   | false   | Trim whitespace including leading and trailing newlines before validation. WARNING: this will mutate the data permanently. Any leading and trailing whitespace will be removed. |
| len                | string | 1:130   | Length (range) counted in bytes in format "min:max" (inclusive).<br>Min value is `1`. Value `<1` is invalid (use `optional` instead). <br>Max value `string`: `1000`. Max value `unsafe_string`: `10000`.<br> Can be set to a single value for fixed length strings, e.g. `16`. Note this is a string type, so must be wrapped in quotes, e.g. `"2:50"`. |
| runes              | bool   | false   | Validate the length by counting runes (number of Unicode codepoints) instead of bytes.  |
| [replace_unsafe](#validator.string.replace_unsafe) | bool   | false   | Replace common unsafe characters with a safe equivalent or a space. For `validator.string` this will also allow all characters in the [restricted characters](#validator.string.restricted_characters) to be allowed by default (but they will be replaced). <br>__Note:__ enabling this option will mutate input data. This is NOT a suitable replacement for output encoding and the use of safe APIs and does NOT magically make your application imune to vulnerabilities. [More details](#validator.string.replace_unsafe).  |
| replace_other      | bool   | false   | Replace several uncommon symbols with more common alternatives that look similar. This allows us to accept a wider range of characters. This will permanently mutate the data. <br>Refer to `stringSymbolReplacerMap` in [validation_definitions.go](plugin/validation_definitions.go). |
| allow              | string | _none_    | Use this option to add a few individual characters that are required in the string field. Will default to [acceptable characters list](#validator.string.acceptable_characters).<br>With `replace_unsafe` disabled, [restricted characters](#validator.string.restricted_characters) are not permitted in this option for `validator.string`. <br>With `replace_unsafe` enabled, all characters in [restricted characters](#validator.string.restricted_characters) are allowed by default unless `allow` option is defined. In that case it will override that behaviour and restrict permitted characters to the defined chars only. Restricted characters will still be replaced. <br>Does not accept regex tokens. All provided characters are used as literals, see examples.  |
| [symbols](#validator.string.symbols) | SymbolCategory | _none_ | Allow zero, one or more [predefined symbol categories](#validator.string.symbols) (`repeated`). By default, only the symbols included in the [acceptable characters list](#validator.string.acceptable_characters) are allowed.  |
| multiline          | bool   | false   | Allow newline characters for multiline strings. By default, multiline strings are not allowed. Note that `\r` will be stripped if this option is enabled (`\n` will be retained). With `trim` option enabled, leading and trailing newlines will be removed.  |
| validate_encoding  | bool   | __true__ | Check for invalid encoding and reject the input in that case. If this is causing issues with normal data, the validation implementation should be improved rather than disabling this option. |
| sanitise_pua  | bool   | false   | Sanitise (remove) Private Use Area Codepoints in the Basic Multilingual Plane. This should normally be enabled and might want to enable this by default later. With this option disabled, PAU characters will be rejected instead of removed. |



### Restricted characters <a name="validator.string.restricted_characters"></a>
Strings pose a potential security risk as by default they can contain any character including control characters and potentially dangerous special characters.
<details>
  <summary>More details</summary>
  If an insecure application is using input data in (sub)system commands or unescaped output, it can result in high impact vulnerabilities like HTML injection (e.g. Cross-Site Scripting, XSS), Command injection, SQL injection, and Path traversal amongst others. The underlying problem is that user supplied data is interpreted as code instead of data, often as the result of a specific character that is used to escape the data context. Therefore, reserved characters that have a meaning in a specific context such as `< > ' " | / & ` should be treaten with care.      
     
  Not allowing certain potentially dangerous special characters or replacing them with a safe equivalent can prevent or impede successfull exploitation of a vulnerable application. While this should be considered an additional layer of protection and certainly not the primary method to prevent against injection vulnerabilities (use secure coding practices like output escaping and safe APIs), it is still an effective method and more often than not a field shouldn't even need to accept all those special characters. For instance, a person's name is highly unlikely to contain a `<` character (yes, I know, [names are tricky](https://www.w3.org/International/questions/qa-personal-names)) so such characters should not be accepted in the first place. 
   
  If an application is accepting one of these characters in a string field, it has to properly output-encode the data for the right context and also never use the untrusted data to dynamically create strings for communication with subsystems (e.g. SQL queries, OS commands). When following secure coding best practices, an application should be able to handle these characters correctly, interpreting them as data, not code but reality shows that this is not always the case.
</details>
   
&nbsp;   
The following is a list of restricted characters that are considered potentially dangerous as they are often used in web app attacks. 

| Restricted characters  |    |                        |    |                        |
|------------------------|----|------------------------|----|------------------------|
| `!` Exclamation mark   |    | `"` Quotation mark     |    | `#` Number sign        |
| `%` Percent sign       |    | `&` Ampersand          |    | `'` Apostrophe         |
| `*` Asterisk           |    | `+` Plus sign          |    | `-` Hyphen-minus       |
| `/` Solidus            |    | `;` Semicolon          |    | `<` Less-than sign     |
| `=` Equals sign        |    | `>` Greater-than sign  |    | `\` Reverse solidus    |
| `` ` `` Grave accent   |    | ``\|`` Vertical line   |    |  |

Refer to `stringUnsafeReplacerMap` in [validation_definitions.go](plugin/validation_definitions.go) for the full list.

### Accepting restricted characters with `replace_unsafe` <a name="validator.string.replace_unsafe"></a>
Enable the `replace_unsafe` option in the `string` validator to accept any of these restricted characters. The validator will replace the character with a safe equivalent or space character instead. Refer to `stringUnsafeReplacerMap` in [validation_definitions.go](plugin/validation_definitions.go) for the mapping. Note that enabling this option will mutate the input data (can be undone by reversing the process if required, except for characters that are mapped to a space).   

__Note: not suitable for fields that require one of the restricted characters to be represented in their original form.__ 
<details>
  <summary>Embedded URL example</summary>
  For instance, if an application needs to accept URLs as part of generic text, the `replace_unsafe` option might break the URL as the `/` character in `https://` will be replace by `∕` (it's not visible, but it is a different character) and any `&`-characters in the URL parameters will be replaced by a space if not URL-encoded. This could result in a URL that is not working. For instance, Chrome browser will replace `http:∕∕example.com` with `http://xn--example-bf0da.com/`.
</details>
&nbsp;

### Acceptable characters list <a name="validator.string.acceptable_characters"></a>
For `validator.string`, the following characters/categories are considered valid and accepted by the validator by default:

| Character/category        | Regex    | From Unicode Category              |
|---------------------------|----------|---------------------------------------------|
| Letter                    | `\pL`      | Letter: any kind of letter from any language        |
| Number                    | `\pN`      | Number: any kind of numeric character in any script |
| ` `  Space                  | `\x{0020}` | Space Separator (Zs) |
| `(`  Left Parenthesis       | `\x{0028}` | Open Punctuation (Ps) |
| `)`  Right Parenthesis      | `\x{0029}` | Close Punctuation (Pe) |
| `,`  Comma                  | `\x{002C}` | Other Punctuation (Po) |
| `.`  Full Stop              | `\x{002E}` | Other Punctuation (Po) |
| `:`  Colon                  | `\x{003A}` | Other Punctuation (Po) |
| `?`  Question Mark          | `\x{003F}` | Other Punctuation (Po) |
| `@`  Commercial At          | `\x{0040}` | Other Punctuation (Po) |
| `[`  Left Square Bracket    | `\x{005B}` | Open Punctuation (Ps) |
| `]`  Right Square Bracket   | `\x{005D}` | Close Punctuation (Pe) |
| `_`  Low Line               | `\x{005F}` | Connector Punctuation (Pc) |
| `¿`  Inverted Question Mark | `\x{00BF}` | Other Punctuation (Po) |
| `–`  En Dash                | `\x{2013}` | Dash Punctuation (Pd) |
| `�`  RuneError              | `\x{FFFD)` | Only when `validate_encoding` option is disabled |

For `validator.unsafe_string`, all of the above are accepted, plus:
| Character/category        | Regex    | From Unicode Category |
|---------------------------|----------|-----------------------|
| `!`  Exclamation Mark | `\x{0021}` | Other Punctuation (Po)    |
| `"`  Quotation Mark   | `\x{0022}` | Other Punctuation (Po)    |
| `#`  Number Sign      | `\x{0023}` | Other Punctuation (Po)    |
| `%`  Percent Sign     | `\x{0025}` | Other Punctuation (Po)    |
| `&`  Ampersand        | `\x{0026}` | Other Punctuation (Po)    |
| `'`  Apostrophe       | `\x{0027}` | Other Punctuation (Po)    |
| `*`  Asterisk         | `\x{002A}` | Other Punctuation (Po)    |
| `-`  Hyphen-Minus     | `\x{002D}` | Dash Punctuation (Pd)     |
| `/`  Solidus          | `\x{002F}` | Other Punctuation (Po)    |
| `\`  Reverse Solidus  | `\x{005C}` | Other Punctuation (Po)    |


Refer to `stringReDefaultSafe` and `stringReDefaultUnsafe` in [validation_definitions.go](plugin/validation_definitions.go) for the full list.

### Choosing between string or unsafe_string
Rule of thumb: use `validator.string` with `replace_unsafe` option.

The main difference between the two is that `unsafe_string` will allow you to accept certain reserved characters and accepts longer length strings. Really only use `unsafe_string` if you need to accept restricted characters and cannot use `replace_unsafe` option or if dealing with strings longer than 1000 bytes. There should not be many fields where this applies.   

Reality is that often at least few of the reserved characters are required. While generally less than 2% of data contains such characters (at least in SafetyCulture's iAuditor), it is not always possible to just reject the input that contains such characters for legit use like `Daily vehicle inspection (<5 ton GVM)`.
By defining the `replace_unsafe` option, the validator will accept restricted characters and replaces these will a similar looking Unicode character, e.g. `Daily vehicle inspection (˂5 ton GVM)`.

### Symbols <a name="validator.string.symbols"></a>
To accept a certain range of symbols, Unicode categories and a custom allow list is used to accept a wide range of symbol-like characters. When accepting additional symbol categories, it is advised to properly test the possible characters/symbols as they might not be available in all web/system fonts or cause certain application components to misbehave. Refer to the provided links for a list of possible symbols in a category.

| Symbol name | Description |
|-------------|-------------|
| `PUNCTUATION` | Only permitted with `replace_unsafe` option for `validator.string`<br>[Unicode Sm](https://www.fileformat.info/info/unicode/category/Sm/list.htm) - Math symbols like `+ < > = \| ~ 𝛁`<br>[Unicode Po](https://www.fileformat.info/info/unicode/category/Po/list.htm) - Punctuation, other category symbols like `! # % ? - ` |
| `CURRENCY` | [Unicode Sc](https://www.fileformat.info/info/unicode/category/Sc/list.htm) - Currency symbols like `$ £ ¥ €` |
| `MODIFIER` | Only permitted with `replace_unsafe` option for `validator.string`<br>[Unicode Sk](https://www.fileformat.info/info/unicode/category/Sk/list.htm) - Modifier symbols like `^ ˥ 🏻` |
| `OTHER` | [Unicode So](https://www.fileformat.info/info/unicode/category/So/list.htm) - Symbol, other symbols like `© ♬ 🐳 😗` |
| `MARK` | [Unicode M](https://www.fileformat.info/info/unicode/category/index.htm) - Mark and combining symbols (all sub cats Spacing Combining, Enclosing, Nonspacing) |
| `COMMON` | List of symbols we often see but often allowing the entire symbol class is not necessary. This list is dynamic in nature and symbols can be added if needed (when safe).<br>Refer to `validator.SymbolCategory_COMMON` in [validation_definitions.go](plugin/validation_definitions.go) for the current list. |

### Restrictions
During generation of the validators (protoc compilation), some checks are in place to ensure that `validator.string` has safe properties. You might run into compilation errors as the result of a panic generated by the validator plugin. The error message should be descriptive enough to work out the reason. There are several additional checks in place such as a deny list of control characters (see `charDenyList` in [validation_definitions.go](plugin/validation_definitions.go)) that cannot be accepted as part of the `allow` option.

Example usage:
```
string title = 1 [(validator.string) = { len: "2:30" }];
 
// In a name, restricted ' and - characters are required, so enable replace_unsafe
string name = 2 [(validator.string) = { len: "1:50", replace_unsafe: true, allow: "'-" }];
 
// In case many restricted characters are required, omit the allow option to allow all of them
string permissive = 3 [(validator.string) = { len: ":100", replace_unsafe: true  }];

// Example with all options, for illustration purposes only
string all_options = 4 [(validator.string) = { 
    optional: true, 
    trim: true, 
    len: "6-20", 
    runes: true, 
    replace_unsafe: true, 
    replace_other: true, 
    sanitise_pua: true,
    allow: "+#/-",
    symbols: [COMMON, CURRENCY, MARK, PUNCTUATION, MODIFIER, OTHER], 
    multiline: true,
    validate_encoding: false
  }];

// Example of a very permissive string that will accept almost any legit input.
// Try to make validation more restrictive, but could be used as a starting point.
// It will still reject bad inputs like invalid encoded data and control characters.
string sc_permissive = 5 [(validator.string) = { 
     len: ":1000",
     replace_unsafe: true, 
     replace_other: true,
     sanitise_pua: true, 
     symbols: [COMMON, CURRENCY, MARK, PUNCTUATION, MODIFIER, OTHER], 
  }];

```
Refer to [valtest.proto](valtest/valtest.proto) for more examples and test cases.

### String validation steps (internal workings) <a name="validator.string.internals"></a>
String validation comprises of the following steps:
1. Preparations: based on defined values and constants in `validation_definitions.go`, the `genStringGenerics()` and `genStringValidator()` methods in `plugin.go` prepare some of the reusable patterns and replacers.
1. Normalisation and canonicalisation
   1. Normalise Unicode NFD strings to NFC strings: this basically transforms the string to the most compact representation. For instance, combining runes such as e + ´ will be represented as a single rune é. If normalisation fails, the input is rejected. Checking is done using `norm.NFD.IsNormalString` and conversion is done using `transform.String(transform.Chain(norm.NFD, norm.NFC), <input>)`. Error message: `must be normalisable to NFC`.
   1. Check for rune errors: with `validate_encoding` enabled, the input is rejected if it contains the the `utf8.RuneError` rune as that indicates that the encoding is incorrect resulting in malformed `� (U+FFFD)` characters. Error message: `must have valid encoding`.
   1. Validate UTF-8 encoding: with `validate_encoding` enabled, test if the string is in valid UTF-8 encoding using `utf8.ValidString()`. Error message: `must be a valid UTF-8-encoded string`.
1. Sanitisation
   1. Sanitise whitespace: with `trim` enabled, strip whitespace using `strings.TrimSpace`.
   1. Replace restricted characters: with `replace_unsafe` enabled:
      1. With `allow` option defined: allow and replace restricted characters defined in the allow option. Will not accept other restricted chracters. Uses `string.ReplaceAll` to replace the characters with the alternative that is defined in `stringUnsafeReplacerMap`.
      1. With `allow` option not defined: allow and replace all restricted characters in the [restricted characters](#validator.string.restricted_characters) list. Uses a previously prepared replacer `_unsafe_char_replacer` to replace all characters in `stringUnsafeReplacerMap`.
   1. Strip carriage return characters: with `multiline` enabled, removes `\r` characters and leaves `\n` for consistent new line behaviour (some systems use `\r\n` for new lines).
   1. Replace other, rare symbols (non-restricted) with a more common alternative: replace symbols in the `stringSymbolReplacerMap` list using a previously prepared replacer `replacer_symbol_allowed`.
   1. Sanitise Private Use Area Codepoints in the Basic Multilingual Plane: with `sanitise_pau` enabled, runes matching `_regex_pua` (`[\x{E000}-\x{F8FF}]` currently) will be removed.
1. Validation (in this order)
   1. Validate length using `len` or `utf8.RuneCountInString` when `runes` option is enabled. Error message: `value must have length X` or `value must have length between X and Y`.
   1. Validate input data against the allow list using `regex.MatchString`. The regex pattern is dynamically generated (based on validator options and `allow`) and then written to `<package>.validator_regex.pb.go`. Patterns are reused for other fields if possible. Error message: `value must only have valid characters`.

&nbsp;

## Enum validation <a name="validator.enum_required"></a>
Validate enum values.

| Option       | Type | Default | Description |
|--------------|------|---------|-------------|
| enum_required | bool | false  | Set this as an optional field. If it's false, the validation is skipped. If it's true, the value for the enum field is required or must be a non-zero value. |

Example usage:
```
    MyEnum enum = 1 [(validator.enum_required) = true];
```

## URL validation: validator.url (URLRules) <a name="validator.url"></a>
Validate a string for a valid URL format. This is a loose validation, focusing on safe characters in the URL and basic format. This validator does not validate if domains are valid and also accepts IP addresses and localhost values.   
*Warning:* when validation passes, it does _not_ result in a URL that is necessarily safe to fetch/resolve, only that the characters in the provided string are expected in a URL. Additional validations are required if you need to fetch the URL on the server to prevent SSRF including access to internal URLs.   

Implementation details including acceptable characters and default min/max length in IsValidURL method in [validation_helpers.go](plugin/validation_helpers.go).    

| Option       | Type | Default | Description |
|--------------|------|---------|-------------|
| optional     | bool | false   | Set this as an optional field. It will allow the string to be empty, validation is skipped in that case. |
| schemes     | repeated string | ["https"]   | Define one or more valid schemes. |
| allow_fragment     | bool | false   | Allow fragments in the URL. |
| allow_http     | bool | false   | Allow http scheme in additional to https as default schemes. Shortcut to `schemes: ["http", "https"]`. |

Example usage:
```
    string url = 1 [(validator.url) = {}];

    string url_demo = 2 [(validator.url) = { 
            optional: true, 
             schemes: ["ftp", "ftps"],
      allow_fragment: true,
          allow_http: true
    }];
```

&nbsp;

### Testing <a name="validator.testing"></a>
There is a new testing suite available in the [valtest](valtest/) folder that can be invoked to test almost any of the validator options. Run it as follows:

```make generate && make govalidator-valtest && make govalidator-valtest-test```

Only errors are displayed. Add the verbose (`-v`) flag in the `Makefile` for more details. It is possible to add additional test payloads, read from a file (one test case per line). Check [valtest_test.go](valtest/valtest_test.go) for details and examples.


## Legacy Validator Fields <a name="validator.legacy"></a>

__Deprecated: do not define the following validators for new fields. Use one of the newer validation options above.__

```
message ExampleMessage {
  // returns an error if the string cannot be parsed as a UUID
  string id = 1 [(validator.uuid) = true];
  // bytes can also be parsed as UUID with support for gogo
  bytes user_id = 2 [(gogoproto.customname) = "UserID", (validator.uuid) = true];
  // strings can validate against a regular expresion
  string email = 3 [(validator.regex) = ".+\\@.+\\..+"];
  // integers can be greater than a value
  int32 age = 4 [(validator.int_gt) = 0];
  // intergers can be less than a value
  int64 speed = 5 [(validator.int_lt) = 110];
  // intergers greater/less than or equal, the can also be combined
  int32 score = 6 [(validator.int_gte) = 0, (validator.int_lte) = 100];
  // validation is created for all messages
  InnerMessage inner = 7;
  // can validate each repeated item too
  repeated bytes ids = 8 [(validator.uuid) = true];
  // only validate if non-zero value
  string media_id = 9 [(validator.uuid) = true, (validator.optional) = true];
  // validate the max length of a string
  string description = 10 [(validator.length_lte) = 2000];
  // validate the min length
  string password = 11 [(validator.length_gte) = 8];
  // You don't need to validate everything
  string no_validation = 12;
  // Trim leading and trailing whitespaces (as defined by Unicode) before doing length check
  string name = 13 [(validator.length_gte) = 6, (validator.length_lte) = 10, (validator.trim_len_check) = true];
}

message InnerMessage {
  string id = 1 [(validator.uuid) = true];
}
```

