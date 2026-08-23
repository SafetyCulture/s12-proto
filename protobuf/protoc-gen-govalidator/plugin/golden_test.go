package plugin_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/SafetyCulture/s12-proto/protobuf/internal/goldentest"
	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/plugin"

	// Registers the descriptors the golden files are generated from.
	_ "github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/valtest"
)

var update = flag.Bool("update", false, "rewrite the golden files from the current output")

// updateCommand is named in every failure message. The -update flag is defined
// by this package, so a command that spans the whole module does not reach it.
const updateCommand = "go test ./protobuf/protoc-gen-govalidator/plugin/ -update"

// valtestDir holds the proto files the golden output covers. Between them they
// exercise every validator option the plugin supports, and a file added there is
// picked up without editing this test.
const valtestDir = "../valtest"

// TestGolden pins the generated text for the whole validator option surface.
//
// The plan and the emitter are separate, so a change to either can move the
// output without failing a behavioural test. Regenerate with -update
// and read the diff: anything in it is a change to code that ships to every
// service.
func TestGolden(t *testing.T) {
	plugin.Reset()

	sources := goldentest.Sources(t, valtestDir)

	gen, err := protogen.Options{}.New(goldentest.Request(t, sources, "paths=source_relative"))
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
			t.Fatalf("reading %s: %v (run `%s` to create it)", golden, err, updateCommand)
		}
		if got := f.GetContent(); got != string(want) {
			t.Errorf("%s does not match %s; run `%s` and review the diff", name, golden, updateCommand)
		}
	}

	if want := 2 * len(sources); len(produced) != want {
		t.Errorf("generated %d files, want %d (one validator and one regex file per source)", len(produced), want)
	}

	pruneStale(t, produced)
}

// TestGoldenMatchesCheckedIn compares each golden against the file protoc wrote
// into ../valtest for the same source.
//
// TestGolden reads its descriptors from the registry, which the compiled valtest
// package populates, so a .proto edited without `make govalidator-valtest` leaves
// the goldens describing the previous version of it and passing. Holding the two
// copies equal turns that into a failure as soon as either is regenerated, and
// stops the copies drifting apart in the meantime.
func TestGoldenMatchesCheckedIn(t *testing.T) {
	goldens, err := filepath.Glob(filepath.Join("testdata", "*.golden"))
	if err != nil {
		t.Fatalf("listing testdata: %v", err)
	}
	if len(goldens) == 0 {
		t.Fatalf("no golden files in testdata")
	}

	for _, golden := range goldens {
		name := strings.TrimSuffix(filepath.Base(golden), ".golden")
		t.Run(name, func(t *testing.T) {
			if body(t, golden) != body(t, filepath.Join(valtestDir, name)) {
				t.Errorf("%s and %s disagree; run `make govalidator-valtest` then `%s`",
					golden, filepath.Join(valtestDir, name), updateCommand)
			}
		})
	}
}

// body returns the file from its package clause onward. The generated header
// above it carries the protoc and plugin versions of whichever run wrote the
// file, which the two copies have no reason to share.
func body(t *testing.T, path string) string {
	t.Helper()

	text, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	i := strings.Index(string(text), "\npackage ")
	if i < 0 {
		t.Fatalf("%s has no package clause", path)
	}
	return string(text)[i+1:]
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
			t.Errorf("%s is left over from a source that no longer generates it; run `%s`", path, updateCommand)
			continue
		}
		if err := os.Remove(path); err != nil {
			t.Fatalf("removing %s: %v", path, err)
		}
		t.Logf("removed stale %s", path)
	}
}

// TestCommittedOutputMatchesGolden checks the generated Go committed under
// valtestDir against the goldens.
//
// The goldens are produced from descriptors rather than by running protoc, so
// they carry no compiler version and cannot be the committed files themselves.
// That leaves room for the committed files to be left behind by a change to the
// emitter, which this closes: everything but the version block must agree.
func TestCommittedOutputMatchesGolden(t *testing.T) {
	goldens, err := filepath.Glob(filepath.Join("testdata", "*.golden"))
	if err != nil {
		t.Fatalf("listing testdata: %v", err)
	}
	if len(goldens) == 0 {
		t.Fatal("no golden files to compare against")
	}

	for _, golden := range goldens {
		name := strings.TrimSuffix(filepath.Base(golden), ".golden")
		committed := filepath.Join(valtestDir, name)

		want, err := os.ReadFile(golden)
		if err != nil {
			t.Fatalf("reading %s: %v", golden, err)
		}
		got, err := os.ReadFile(committed)
		if err != nil {
			t.Fatalf("reading %s: %v", committed, err)
		}

		if withoutVersions(string(got)) != withoutVersions(string(want)) {
			t.Errorf("%s is behind the emitter; run `make govalidator-valtest`", committed)
		}
	}
}

// withoutVersions drops the generated header's version block, which records the
// protoc that produced the file and so differs between machines.
func withoutVersions(source string) string {
	lines := strings.Split(source, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "// \tprotoc") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}
