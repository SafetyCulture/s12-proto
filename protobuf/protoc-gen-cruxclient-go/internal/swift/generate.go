// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Package swift writes the Swift messages of proto files. The output is the
// output of protoc-gen-swift: declarations only, with SwiftProtobuf as the
// runtime that encodes and decodes the messages.
package swift

import (
	"fmt"

	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/types/descriptorpb"
)

// generator holds the state that the files of one run share.
type generator struct {
	opts       *options
	namer      *namer
	storage    *storageAnalyzer
	extensions *extensionSet
}

// Generate adds the Swift file of each file to generate to the response.
func Generate(rw protoplugin.ResponseWriter, request protoplugin.Request) error {
	opts, err := parseOptions(request.Parameter())
	if err != nil {
		return err
	}
	files, err := Render(request.AllFileDescriptorProtos(), toGenerate(request.FileDescriptorProtosToGenerate()), opts)
	if err != nil {
		return err
	}
	for _, f := range files {
		rw.AddFile(f.Name, f.Content)
	}
	return nil
}

func toGenerate(protos []*descriptorpb.FileDescriptorProto) []string {
	names := make([]string, len(protos))
	for i, p := range protos {
		names[i] = p.GetName()
	}
	return names
}

// OutputFile is one generated file.
type OutputFile struct {
	Name    string
	Content string
}

// Render generates the Swift files of the named proto files. The descriptors
// hold every file after the files it imports.
func Render(protos []*descriptorpb.FileDescriptorProto, names []string, opts *options) ([]OutputFile, error) {
	reg, err := buildFiles(protos)
	if err != nil {
		return nil, err
	}
	g := &generator{opts: opts, storage: newStorageAnalyzer()}
	var out []OutputFile
	for _, name := range names {
		f, ok := reg.files[name]
		if !ok {
			return nil, fmt.Errorf("a file to generate has no descriptor: %s", name)
		}
		content, err := g.generateFile(f)
		if err != nil {
			if len(names) > 1 {
				err = fmt.Errorf("%s: %w", name, err)
			}
			return nil, err
		}
		out = append(out, OutputFile{Name: opts.outputFilename(f.name), Content: content})
	}
	return out, nil
}

// ParseParameter reads the options of a request parameter.
func ParseParameter(parameter string) (*options, error) { return parseOptions(parameter) }
