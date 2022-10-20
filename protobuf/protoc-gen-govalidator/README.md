# Protocol buffer validators compiler

A `protoc` plugin that generates `Validate() error` functions on Go proto `struct`s based on field options inside `.proto` files. The validation functions are code-generated and thus don't suffer on performance from tag-based reflection on deeply-nested messages.

## Table of Contents
- [Usage] (#usage)
- [Development] (#development)
- [Validators] (#validator.list)
- [Legacy Validator Fields](#validator.legacy)
- [Testing Validators](#validator.testing)

&nbsp;

## Usage <a name="usage"></a>

```
$ go get github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator
$ protoc -I. --gogo_out=:. --govalidator_out=. example.proto
```
&nbsp;

## Development <a name="development"></a>

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

## Validators <a name="validator.list"></a>
There are a number of common fields for which validators have been developed. For detailed documentation and usage, visit the links to the confluence pages.
- Email Validation - https://safetyculture.atlassian.net/l/cp/KadCa73c 
- ID Validation - https://safetyculture.atlassian.net/l/cp/Esbr10vG 
- String Validation - https://safetyculture.atlassian.net/l/cp/V0zvBPA1 
- Enum Validation - https://safetyculture.atlassian.net/l/cp/Ed1DCibC 
- URL Validation - https://safetyculture.atlassian.net/l/cp/4VmN3deB 
- Number Validation - https://safetyculture.atlassian.net/l/cp/Bf3iGcNt 

If new validators are added, please be sure to document them including their usage and options. 

Link to the last version of the readme which had the complete documentation to date - https://github.com/SafetyCulture/s12-proto/blob/f8d868b8135e3f2438bfaca234e27b2f305c59d9/protobuf/protoc-gen-govalidator/README.md

&nbsp;

## Legacy Validator Fields <a name="validator.legacy"></a>

__Deprecated: do not define the following validators for new fields. Use one of the newer validation options above.__

```
message ExampleMessage {
  // returns an error if the string cannot be parsed as a UUID
  string id = 1 [(validator.uuid) = true];
  // bytes can also be parsed as UUID with support for gogo
  bytes user_id = 2 [(gogoproto.customname) = "UserID", (validator.uuid) = true];
  // strings can validate against a regular expression
  string email = 3 [(validator.regex) = ".+\\@.+\\..+"];
  // integers can be greater than a value
  int32 age = 4 [(validator.int_gt) = 0];
  // integers can be less than a value
  int64 speed = 5 [(validator.int_lt) = 110];
  // integers greater/less than or equal, the can also be combined
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
&nbsp;

## Testing <a name="validator.testing"></a>
There is a new testing suite available in the [valtest](valtest/) folder that can be invoked to test almost any of the validator options. Run it as follows:

```make generate && make govalidator-valtest && make govalidator-valtest-test```

Only errors are displayed. Add the verbose (`-v`) flag in the `Makefile` for more details. It is possible to add additional test payloads, read from a file (one test case per line). Check [valtest_test.go](valtest/valtest_test.go) for details and examples.
