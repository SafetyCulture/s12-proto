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
  string id = 1 [(validator.field) = {uuid: true}];
  // Bytes can also be parsed as UUID
  bytes user_id = 2 [(validator.field) = {uuid: true}];
  // Strings can validate against a regular expresion
  string description = 3 [(validator.field) = {regex: "^[a-z]{2,5}$"}];
}
```