// Code generated by protoc-gen-govalidator. DO NOT EDIT.
// versions:
// 	protoc-gen-govalidator v2.7.1
// 	protoc                 v5.28.2
// source: example.proto

package example

import (
	_ "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
	regexp "regexp"
)

// Pattern for ExampleMessage_StringOptional
const _regex_val_d4db71516b8749dc594e5bf604c6a110 = `^[\pL\pN\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{2013}]+$`

var _regex_d4db71516b8749dc594e5bf604c6a110 = regexp.MustCompile(_regex_val_d4db71516b8749dc594e5bf604c6a110)
