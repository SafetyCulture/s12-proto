# protocol buffer validators compiler

A `protoc` plugin that generates `Validate() error` functions on Go proto `struct`s based on field options inside `.proto` files. The validation functions are code-generated and thus don't suffer on performance from tag-based reflection on deeply-nested messages.

## Usage

```
$ go get github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator
$ protoc -I. --gogo_out=:. --govalidator_out=. example.proto
```

## Validator Fields

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
}

message InnerMessage {
  string id = 1 [(validator.uuid) = true];
}
```
