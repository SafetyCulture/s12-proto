// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Package swift writes the Swift messages of proto files. The output is the
// output of protoc-gen-swift: declarations only, with SwiftProtobuf as the
// runtime that encodes and decodes the messages.
package swift

import (
	"errors"
	"fmt"
	"runtime/debug"

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
// The function isPluginKey names the options of the other outputs of the plugin, which this output ignores.
func Generate(rw protoplugin.ResponseWriter, request protoplugin.Request, isPluginKey func(key string) bool) error {
	opts, err := parseOptions(request.Parameter(), isPluginKey)
	if err != nil {
		return err
	}
	files, err := render(request.AllFileDescriptorProtos(), toGenerate(request.FileDescriptorProtosToGenerate()), opts)
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

// outputFile is one generated file.
type outputFile struct {
	Name    string
	Content string
}

// render generates the Swift files of the named proto files. The descriptors
// hold every file after the files it imports.
func render(protos []*descriptorpb.FileDescriptorProto, names []string, opts *options) ([]outputFile, error) {
	reg, err := buildFiles(protos)
	if err != nil {
		return nil, err
	}
	g := &generator{opts: opts, storage: newStorageAnalyzer()}
	var out []outputFile
	for _, name := range names {
		f, ok := reg.files[name]
		if !ok {
			return nil, fmt.Errorf("a file to generate has no descriptor: %s", name)
		}
		content, err := g.generateFileRecovering(f)
		if err != nil {
			var internal *internalError
			if len(names) > 1 && !errors.As(err, &internal) {
				err = fmt.Errorf("%s: %w", name, err)
			}
			return nil, err
		}
		out = append(out, outputFile{Name: opts.outputFilename(f.name), Content: content})
	}
	return out, nil
}

// internalError is a panic of the generator. Its text names the file already.
type internalError struct{ text string }

func (e *internalError) Error() string { return e.text }

// generateFileRecovering turns a panic of the generator into an error with the file name and the stack.
func (g *generator) generateFileRecovering(f *file) (content string, err error) {
	defer func() {
		if r := recover(); r != nil {
			content, err = "", &internalError{text: fmt.Sprintf("internal error of the Swift generator on %s: %v\n%s", f.name, r, debug.Stack())}
		}
	}()
	return g.generateFile(f)
}
