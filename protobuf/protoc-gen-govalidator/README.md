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
  // Returns an error if the string cannot be parsed as a UUID
  string id = 1 [(validator.uuid) = true];
  // Bytes can also be parsed as UUID
  bytes user_id = 2 [(validator.uuid) = true];
  // Strings can validate against a regular expresion
  string description = 3 [(validator.regex) = "^[a-z]{2,5}$"];
  // integers can be greater than a value
  int32 age = 4 [(validator.int_gt) = 0];
  // intergers can be less than a value
  int64 speed = 5 [(validator.int_lt) = 110];
  // intergers greater/less than or equal, the can also be combined
  int32 score = 6 [(validator.int_gte) = 0, (validator.int_lte) = 100];
  // Validation is created for all messages
  InnerMessage inner = 7;
}

message InnerMessage {
  string id = 1 [(validator.uuid) = true];
}
```