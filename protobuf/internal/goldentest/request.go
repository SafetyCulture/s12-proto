// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Package goldentest assembles generator requests for the plugin golden tests.
package goldentest

import (
	"path/filepath"
	"sort"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/pluginpb"
)

// Sources lists the proto files in dir by base name, so adding one there extends
// the golden coverage without touching the test that calls this.
func Sources(t *testing.T, dir string) []string {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join(dir, "*.proto"))
	if err != nil {
		t.Fatalf("listing %s: %v", dir, err)
	}
	if len(paths) == 0 {
		t.Fatalf("no proto files under %s", dir)
	}

	sources := make([]string, 0, len(paths))
	for _, path := range paths {
		sources = append(sources, filepath.Base(path))
	}
	sort.Strings(sources)
	return sources
}

// Request assembles a generator request for sources out of the descriptors their
// generated Go packages registered, so the golden files track the .proto files
// without a checked-in descriptor set to keep in step.
//
// The calling test must import the package holding those generated files.
func Request(t *testing.T, sources []string, parameter string) *pluginpb.CodeGeneratorRequest {
	t.Helper()

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: sources,
		Parameter:      proto.String(parameter),
	}

	seen := map[string]bool{}
	var collect func(fd protoreflect.FileDescriptor)
	collect = func(fd protoreflect.FileDescriptor) {
		if seen[fd.Path()] {
			return
		}
		seen[fd.Path()] = true

		imports := fd.Imports()
		for i := 0; i < imports.Len(); i++ {
			collect(imports.Get(i).FileDescriptor)
		}
		req.ProtoFile = append(req.ProtoFile, protodesc.ToFileDescriptorProto(fd))
	}

	for _, source := range sources {
		fd, err := protoregistry.GlobalFiles.FindFileByPath(source)
		if err != nil {
			t.Fatalf("%s has no registered descriptor (%v); run `make govalidator-valtest` so its generated Go exists", source, err)
		}
		collect(fd)
	}

	// Guard against a dependency that never made it into the request, which
	// protogen would otherwise report as a confusing missing-import error.
	for _, fd := range req.ProtoFile {
		for _, dep := range fd.GetDependency() {
			if !seen[dep] {
				t.Fatalf("%s depends on %s, which is not in the request", fd.GetName(), dep)
			}
		}
	}

	return req
}
