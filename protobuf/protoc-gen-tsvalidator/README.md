# Protocol buffer validators compiler

A `protoc` plugin that generates `Validate() error` functions on Go proto `struct`s based on field options inside `.proto` files. The validation functions are code-generated and thus don't suffer on performance from tag-based reflection on deeply-nested messages.

&nbsp;

## Usage

```
$ go get github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-tsvalidator
$ protoc -I. --gogo_out=:. --tsvalidator_out=. example.proto
```
&nbsp;

## Development

Use these commands at the root folder of this repository for testing:
```bash
cd s12-proto

# regenerate examples and run tests
make tsvalidator
make tsvalidator-test

# regenerate valtest and run tests
make tsvalidator-valtest
make tsvalidator-valtest-test
```
&nbsp;

## Playground

The validator plugin is invoked by the protobuf compiler tool directly to generate the respective validator code.
As a result, it makes it very hard to attach the plugin to a debugger and inspect the code during that process.

A sample code has been provided that allows you to attach and step through the plugin during the code generation process. 
The basic premise of the playground approach is
- generate the `CodeGeneratorRequest` object of the proto files (via makefile)
- run the validator plugin on the request object (via the .go code)
- the resulting files are printed to `stdout`

The `CodeGeneratorRequest` object is what is consumed by the protobuf compiler as well. 
The code in `playground_main.go` is a standard go file with a main function.

There are a couple of Makefile targets that help us to achieve this 
- `make playground-req-object`: generate the request object from the `playground.proto`
- `make playground-generated-code`: generate the eventual code in the respective `validator_regex.pb.go` and `validator.pb.go`

At minimum, you are required to run the `playground-req-object` target before you run the code in `playground_main.go`

You can modify the `playground.proto` to add in a new type of message or use existing validators to see how the plugin works in action.

## Protobuf Validators
There are a number of common fields for which validators have been developed. For detailed documentation and usage, visit the following link to a private SafetyCulture confluence page - https://safetyculture.atlassian.net/l/cp/wC1Xz6np. 

If new validators are added, please be sure to document them including their usage and options.

Link to the last public version of the readme which had the complete documentation to date - https://github.com/SafetyCulture/s12-proto/blob/f8d868b8135e3f2438bfaca234e27b2f305c59d9/protobuf/protoc-gen-tsvalidator/README.md 

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
