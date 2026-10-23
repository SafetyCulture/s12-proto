// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// goldenParameter is the option string of the APISchema crux Swift lane.
const goldenParameter = "Visibility=Public,FileNaming=PathToUnderscores"

func readImage(t *testing.T, path string) []*descriptorpb.FileDescriptorProto {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	if strings.HasSuffix(path, ".gz") {
		zr, err := gzip.NewReader(bytes.NewReader(raw))
		require.NoError(t, err)
		raw, err = io.ReadAll(zr)
		require.NoError(t, err)
	}
	set := &descriptorpb.FileDescriptorSet{}
	require.NoError(t, proto.Unmarshal(raw, set))
	return set.GetFile()
}

func readLines(t *testing.T, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, f.Close()) })
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	require.NoError(t, sc.Err())
	return lines
}

// compareCorpus renders the named files and compares every output with the file of the same name in goldenDir.
// It lists each difference, so one run shows the whole gap.
func compareCorpus(t *testing.T, protos []*descriptorpb.FileDescriptorProto, names []string, goldenDir string) {
	t.Helper()
	opts, err := parseOptions(goldenParameter)
	require.NoError(t, err)
	out, err := render(protos, names, opts)
	require.NoError(t, err)

	rendered := map[string]string{}
	for _, f := range out {
		rendered[f.Name] = f.Content
	}
	entries, err := os.ReadDir(goldenDir)
	require.NoError(t, err)
	want := map[string]bool{}
	for _, e := range entries {
		want[e.Name()] = true
	}

	var differing, missing, extra []string
	for name, content := range rendered {
		if !want[name] {
			extra = append(extra, name)
			continue
		}
		golden, err := os.ReadFile(filepath.Join(goldenDir, name))
		require.NoError(t, err)
		if string(golden) != content {
			differing = append(differing, name)
		}
	}
	for name := range want {
		if _, ok := rendered[name]; !ok {
			missing = append(missing, name)
		}
	}
	sort.Strings(differing)
	sort.Strings(missing)
	sort.Strings(extra)
	t.Logf("%d files rendered, %d golden files: %d differ, %d are missing from the output, %d have no golden file",
		len(rendered), len(want), len(differing), len(missing), len(extra))
	if len(differing)+len(missing)+len(extra) == 0 {
		return
	}
	for _, name := range differing {
		t.Errorf("differs from the golden file: %s", name)
	}
	for _, name := range missing {
		t.Errorf("the output has no file for the golden file: %s", name)
	}
	for _, name := range extra {
		t.Errorf("the output has a file without a golden file: %s", name)
	}
}

// TestWhitelistGolden compares the output with the output of protoc-gen-swift
// for the Crux whitelist image. regen.sh in testdata documents how the image
// and the golden files come from the schema.
func TestWhitelistGolden(t *testing.T) {
	protos := readImage(t, filepath.Join("testdata", "whitelist.binpb.gz"))
	names := readLines(t, filepath.Join("testdata", "whitelist.txt"))
	compareCorpus(t, protos, names, filepath.Join("testdata", "golden"))
}

// TestControlGolden runs the same comparison on the whole schema image without
// prefixes, against the output of protoc-gen-swift for that image. The image
// and the output are too big to keep in the repository, so
// testdata/regen-control.sh makes them. The test skips unless the variables
// CRUX_SWIFT_CONTROL_IMAGE and CRUX_SWIFT_CONTROL_DIR name them.
func TestControlGolden(t *testing.T) {
	image, dir := os.Getenv("CRUX_SWIFT_CONTROL_IMAGE"), os.Getenv("CRUX_SWIFT_CONTROL_DIR")
	if image == "" || dir == "" {
		t.Skip("CRUX_SWIFT_CONTROL_IMAGE and CRUX_SWIFT_CONTROL_DIR are not set")
	}
	protos := readImage(t, image)
	var names []string
	for _, p := range protos {
		if strings.HasPrefix(p.GetName(), "google/protobuf/") {
			continue
		}
		names = append(names, p.GetName())
	}
	compareCorpus(t, protos, names, dir)
}
