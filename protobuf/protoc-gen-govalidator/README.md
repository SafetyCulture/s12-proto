# Protocol buffer validators compiler

A `protoc` plugin that generates `Validate() error` functions on Go proto `struct`s based on field options inside `.proto` files. The validation functions are code-generated and thus don't suffer on performance from tag-based reflection on deeply-nested messages.

&nbsp;

## Usage

```
$ go get github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator
$ protoc -I. --gogo_out=:. --govalidator_out=. example.proto
```
&nbsp;

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
&nbsp;

## Protobuf Validators
There are a number of common fields for which validators have been developed. For detailed documentation and usage, visit the following link to a private SafetyCulture confluence page - https://safetyculture.atlassian.net/l/cp/wC1Xz6np. 

If new validators are added, please be sure to document them including their usage and options.

Link to the last public version of the readme which had the complete documentation to date - https://github.com/SafetyCulture/s12-proto/blob/f8d868b8135e3f2438bfaca234e27b2f305c59d9/protobuf/protoc-gen-govalidator/README.md 

&nbsp;

## Legacy Validator Fields

__Deprecated: do not define the following validators for new fields. Use one of the newer validation options above.__

- validator.uuid
- validator.int_gt
- validator.int_lt
- validator.int_gte
- validator.int_lte
- validator.length_lte
- validator.length_gte
- validator.legacy_id
- validator.trim_len_check
