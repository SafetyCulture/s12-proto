package plugin_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/SafetyCulture/s12-proto/protobuf/internal/goldentest"
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-csharpvalidator/plugin"

	// Registers the descriptors the golden files are generated from.
	_ "github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/valtest"
)

var update = flag.Bool("update", false, "rewrite the generated files from the current output")

// updateCommand is named in every failure message. The -update flag is defined
// by this package, so a command that spans the whole module does not reach it.
const updateCommand = "go test ./protobuf/protoc-gen-csharpvalidator/plugin/ -update"

// valtestDir holds the proto files the golden output covers. Between them they
// exercise every validator option the plugin supports, and a file added there is
// picked up without editing this test.
const valtestDir = "../../protoc-gen-govalidator/valtest"

// generatedDir holds the output the conformance project compiles and runs. The
// comparison is against those files rather than a copy of them, so output that
// has moved on from what is committed fails here instead of reaching the
// conformance run as a stale result.
const generatedDir = "../conformance/generated"

// TestGolden pins the generated C# for the whole validator option surface.
//
// Regenerate with `go test -update` and read the diff: anything in it is a change
// to code that ships to every service.
func TestGolden(t *testing.T) {
	sources := goldentest.Sources(t, valtestDir)

	gen, err := protogen.Options{}.New(goldentest.Request(t, sources, "paths=source_relative"))
	if err != nil {
		t.Fatalf("building the generator: %v", err)
	}

	// The conformance project compiles the runtime from its own source, so only the
	// validators are written here.
	if err := plugin.Generate(gen, plugin.Options{Include: []string{"*.proto"}, SkipRuntime: true}); err != nil {
		t.Fatalf("generating: %v", err)
	}

	resp := gen.Response()
	if resp.Error != nil {
		t.Fatalf("generating: %s", resp.GetError())
	}

	produced := map[string]bool{}
	for _, f := range resp.File {
		name := filepath.Base(f.GetName())
		path := filepath.Join(generatedDir, name)
		produced[name] = true

		if *update {
			if err := os.WriteFile(path, []byte(f.GetContent()), 0o644); err != nil {
				t.Fatalf("writing %s: %v", path, err)
			}
			continue
		}

		want, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v (run `%s` to create it)", path, err, updateCommand)
		}
		if got := f.GetContent(); got != string(want) {
			t.Errorf("%s does not match %s; run `%s` and review the diff", name, path, updateCommand)
		}
	}

	if len(produced) != len(sources) {
		t.Errorf("generated %d files, want %d (one validator per source)", len(produced), len(sources))
	}

	pruneStale(t, produced)
}

// pruneStale reports generated files that nothing produces any more, and removes
// them when the output is being rewritten.
func pruneStale(t *testing.T, produced map[string]bool) {
	t.Helper()

	existing, err := filepath.Glob(filepath.Join(generatedDir, "*.g.cs"))
	if err != nil {
		t.Fatalf("listing %s: %v", generatedDir, err)
	}

	for _, path := range existing {
		if produced[filepath.Base(path)] {
			continue
		}
		if !*update {
			t.Errorf("%s is left over from a source that no longer generates it; run `%s`", path, updateCommand)
			continue
		}
		if err := os.Remove(path); err != nil {
			t.Fatalf("removing %s: %v", path, err)
		}
		t.Logf("removed stale %s", path)
	}
}
