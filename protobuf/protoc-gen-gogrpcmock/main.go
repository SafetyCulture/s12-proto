// Copyright (c) 2019 SafetyCulture Pty Ltd. All Rights Reserved.

package main

import (
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-gogrpcmock/plugin"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	resp := command.GeneratePlugin(req, plugin.New(), ".mock.go")
	command.Write(resp)
}
