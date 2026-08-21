package plugin_test

import (
	"flag"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/plugin"

	// Registers the descriptors the golden files are generated from.
	_ "github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/valtest"
)

var update = flag.Bool("update", false, "rewrite the golden files from the current output")

// valtestDir holds the proto files the golden output covers. Between them they
// exercise every validator option the plugin supports, and a file added there is
// picked up without editing this test.
const valtestDir = "../valtest"

// TestGolden pins the generated text for the whole validator option surface.
//
// The plan and the emitter are separate, so a change to either can move the
// output without failing a behavioural test. Regenerate with `go test -update`
// and read the diff: anything in it is a change to code that ships to every
// service.
func TestGolden(t *testing.T) {
	sources := goldenSources(t)

	gen, err := protogen.Options{}.New(request(t, sources))
	if err != nil {
		t.Fatalf("building the generator: %v", err)
	}

	for _, f := range gen.Files {
		if f.Generate {
			plugin.GenerateFile(gen, f)
		}
	}

	resp := gen.Response()
	if resp.Error != nil {
		t.Fatalf("generating: %s", resp.GetError())
	}

	produced := map[string]bool{}
	for _, f := range resp.File {
		name := filepath.Base(f.GetName())
		golden := filepath.Join("testdata", name+".golden")
		produced[filepath.Base(golden)] = true

		if *update {
			if err := os.WriteFile(golden, []byte(f.GetContent()), 0o644); err != nil {
				t.Fatalf("writing %s: %v", golden, err)
			}
			continue
		}

		want, err := os.ReadFile(golden)
		if err != nil {
			t.Fatalf("reading %s: %v (run `go test ./... -update` to create it)", golden, err)
		}
		if got := f.GetContent(); got != string(want) {
			t.Errorf("%s does not match %s; run `go test ./... -update` and review the diff", name, golden)
		}
	}

	if want := 2 * len(sources); len(produced) != want {
		t.Errorf("generated %d files, want %d (one validator and one regex file per source)", len(produced), want)
	}

	pruneStale(t, produced)
}

// goldenSources lists the proto files under valtestDir, so adding one there
// extends the golden coverage without touching this test.
func goldenSources(t *testing.T) []string {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join(valtestDir, "*.proto"))
	if err != nil {
		t.Fatalf("listing %s: %v", valtestDir, err)
	}
	if len(paths) == 0 {
		t.Fatalf("no proto files under %s", valtestDir)
	}

	sources := make([]string, 0, len(paths))
	for _, path := range paths {
		sources = append(sources, filepath.Base(path))
	}
	sort.Strings(sources)
	return sources
}

// pruneStale reports golden files that nothing generates any more, and removes
// them when the goldens are being rewritten.
func pruneStale(t *testing.T, produced map[string]bool) {
	t.Helper()

	existing, err := filepath.Glob(filepath.Join("testdata", "*.golden"))
	if err != nil {
		t.Fatalf("listing testdata: %v", err)
	}

	for _, path := range existing {
		if produced[filepath.Base(path)] {
			continue
		}
		if !*update {
			t.Errorf("%s is left over from a source that no longer generates it; run `go test ./... -update`", path)
			continue
		}
		if err := os.Remove(path); err != nil {
			t.Fatalf("removing %s: %v", path, err)
		}
		t.Logf("removed stale %s", path)
	}
}

// request assembles a generator request for sources out of the descriptors their
// generated Go packages registered, so the golden files track the .proto files
// without a checked-in descriptor set to keep in step.
func request(t *testing.T, sources []string) *pluginpb.CodeGeneratorRequest {
	t.Helper()

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: sources,
		Parameter:      proto.String("paths=source_relative"),
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
