# Typescript protocol buffer validators compiler

A `protoc` plugin that generates `Zod` schemas for messages and attributes based on field options inside `.proto` files. The validation functions are code-generated and thus don't suffer on performance from tag-based reflection on deeply-nested messages.

The intent of this validation library is to provide input level validations on the client side. We aim to ensure 100% compatibility with the govalidator package and if it is not possible to validate effectively on the frontend for some of the validations/replacements then it should be captured by the backend validation.

> Backend validation has a primary objective of creating a secure and reliable system
> Frontend validation is used to provide the users immediate feedback where possible and to aid by centralizing the validation logic to the proto specification. Thereby also increasing the developer experience.

## Usage

```shell
go get github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-tsvalidator
protoc -I. --gogo_out=:. --tsvalidator_out=. example.proto
```

## Development

Use these commands at current folder:

```shell
cd s12-proto/protobuf/protoc-gen-tsvalidator

# regenerate examples
task gen-example

# Or any custom folder
make gen FOLDER=folder

# regenerate valtest and run tests
task test
```

## In Progress

- types of s12 id
- upper/lower case uuid
- unsafe/safe string
- replace unsafe
- Latitude/Longitude
- rune length vs string length // will result in 7 bytes and 4 runes
- rune length // this NFD string (e + ´) will be normalised to NFC string before len check so still 4 bytes
- run length // len 4 in runes, 5 in bytes
- SanitiseLength
- SanitisePua
- NotSanitisePua
