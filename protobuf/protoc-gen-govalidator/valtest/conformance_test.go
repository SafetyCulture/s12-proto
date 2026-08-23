package valtest

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	vproto "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
)

// The conformance vectors are what a port of these validators is checked
// against, so they live outside this package. Each language asserts against the
// same committed files rather than running the other language's toolchain.
const conformanceDir = "../../validation/conformance/testdata"

var updateVectors = flag.Bool("update", false, "rewrite the conformance vectors from the current behaviour")

// TestConformanceFields records what the generated validators do to each input,
// for every string field of ValTestMessage.
//
// A verdict alone is not enough: the sanitising rules rewrite the field, so a
// port could agree on accept or reject while leaving a different value behind.
// Each record therefore carries the verdict, the error text, and the value as it
// stands afterwards.
func TestConformanceFields(t *testing.T) {
	inputs := conformanceInputs(t)
	base := conformanceBase(t)

	type record struct {
		Field string `json:"field"`
		In    string `json:"in"`
		OK    bool   `json:"ok"`
		Err   string `json:"err,omitempty"`
		Out   string `json:"out"`
	}

	var records []record
	fields := (&ValTestMessage{}).ProtoReflect().Descriptor().Fields()
	for i := range fields.Len() {
		fd := fields.Get(i)
		if fd.Kind() != protoreflect.StringKind {
			continue
		}

		for _, in := range inputs {
			// Start from a message that passes, so the only field that can fail is the
			// one under test.
			msg := proto.Clone(base).(*ValTestMessage)
			r := msg.ProtoReflect()
			if fd.IsList() {
				list := r.Mutable(fd).List()
				list.Truncate(0)
				list.Append(protoreflect.ValueOfString(in))
			} else {
				r.Set(fd, protoreflect.ValueOfString(in))
			}

			rec := record{Field: string(fd.Name()), In: encode(in)}
			if err := msg.Validate(); err != nil {
				rec.Err = err.Error()
			} else {
				rec.OK = true
			}

			after := ""
			if fd.IsList() {
				if list := r.Get(fd).List(); list.Len() > 0 {
					after = list.Get(0).String()
				}
			} else {
				after = r.Get(fd).String()
			}
			rec.Out = encode(after)

			records = append(records, rec)
		}
	}

	compareVectors(t, "fields.jsonl", records)
}

// TestConformanceNestedFields records the same for a string field one level down,
// where the failure is re-rooted under the field holding the message.
//
// The re-rooted text is its own shape rather than the leaf text with a prefix, so
// a port can agree on every leaf and still render the nesting differently.
func TestConformanceNestedFields(t *testing.T) {
	inputs := conformanceInputs(t)
	base := conformanceBase(t)

	type record struct {
		Field string `json:"field"`
		In    string `json:"in"`
		OK    bool   `json:"ok"`
		Err   string `json:"err,omitempty"`
	}

	var records []record
	fields := base.ProtoReflect().Descriptor().Fields()
	for i := range fields.Len() {
		fd := fields.Get(i)
		if fd.Kind() != protoreflect.MessageKind || fd.IsList() || fd.IsMap() {
			continue
		}
		// Only messages the base already carries. Creating one here would leave its
		// own required fields empty, and those would fail ahead of the field under test.
		if !base.ProtoReflect().Has(fd) {
			continue
		}

		nested := fd.Message().Fields()
		for j := range nested.Len() {
			sfd := nested.Get(j)
			if sfd.Kind() != protoreflect.StringKind {
				continue
			}

			for _, in := range inputs {
				msg := proto.Clone(base).(*ValTestMessage)
				inner := msg.ProtoReflect().Mutable(fd).Message()
				if sfd.IsList() {
					list := inner.Mutable(sfd).List()
					list.Truncate(0)
					list.Append(protoreflect.ValueOfString(in))
				} else {
					inner.Set(sfd, protoreflect.ValueOfString(in))
				}

				rec := record{Field: fmt.Sprintf("%s.%s", fd.Name(), sfd.Name()), In: encode(in)}
				if err := msg.Validate(); err != nil {
					rec.Err = err.Error()
				} else {
					rec.OK = true
				}
				records = append(records, rec)
			}
		}
	}

	if len(records) == 0 {
		t.Fatal("no nested string field to record")
	}
	compareVectors(t, "nested.jsonl", records)
}

// TestConformanceDeepNesting records a failure at each depth of a three-level
// message chain.
//
// Re-rooting accumulates: the path grows a segment per level while the innermost
// text stays as it was. A port can get one level right and still repeat the prefix
// or lose a segment further down, so every depth is recorded.
func TestConformanceDeepNesting(t *testing.T) {
	inputs := conformanceInputs(t)
	base := conformanceDeepBase(t)

	type record struct {
		Field string `json:"field"`
		In    string `json:"in"`
		OK    bool   `json:"ok"`
		Err   string `json:"err,omitempty"`
	}

	var records []record
	for _, path := range stringPaths(base.ProtoReflect(), nil) {
		for _, in := range inputs {
			msg := proto.Clone(base).(*MyReqMessage)
			holder := msg.ProtoReflect()
			for _, step := range path[:len(path)-1] {
				holder = holder.Mutable(fieldByName(holder, step)).Message()
			}
			holder.Set(fieldByName(holder, path[len(path)-1]), protoreflect.ValueOfString(in))

			rec := record{Field: strings.Join(path, "."), In: encode(in)}
			if err := msg.Validate(); err != nil {
				rec.Err = err.Error()
			} else {
				rec.OK = true
			}
			records = append(records, rec)
		}
	}

	if len(records) == 0 {
		t.Fatal("no string field to record")
	}
	compareVectors(t, "deep.jsonl", records)
}

// stringPaths lists every string field reachable from m, descending into the
// message fields it already carries.
func stringPaths(m protoreflect.Message, prefix []string) [][]string {
	var paths [][]string

	fields := m.Descriptor().Fields()
	for i := range fields.Len() {
		fd := fields.Get(i)
		here := append(append([]string{}, prefix...), string(fd.Name()))

		switch {
		case fd.IsList() || fd.IsMap():
		case fd.Kind() == protoreflect.StringKind:
			paths = append(paths, here)
		case fd.Kind() == protoreflect.MessageKind && m.Has(fd):
			paths = append(paths, stringPaths(m.Get(fd).Message(), here)...)
		}
	}

	return paths
}

func fieldByName(m protoreflect.Message, name string) protoreflect.FieldDescriptor {
	return m.Descriptor().Fields().ByName(protoreflect.Name(name))
}

// conformanceDeepBase returns the chain every record starts from, and writes it to
// the vector directory so a port begins from the same bytes.
func conformanceDeepBase(t *testing.T) *MyReqMessage {
	t.Helper()

	base := &MyReqMessage{
		UserId: "ab",
		OrgNested: &NestedLevel1Message{
			OrgId3: "abc",
			OrgNested: &NestedLevel2Message{
				OrgId4:    "abcd",
				OrgNested: &NestedLevel3Message{OrgId5: "abcde"},
			},
		},
	}
	if err := base.Validate(); err != nil {
		t.Fatalf("the base chain does not pass validation: %v", err)
	}

	encoded, err := proto.MarshalOptions{Deterministic: true}.Marshal(base)
	if err != nil {
		t.Fatalf("marshalling the base chain: %v", err)
	}

	writeBase(t, "deep_base.b64", encoded)
	return base
}

// TestConformanceHelpers records what each validation helper returns for each
// input. The generated validators call these, and a port reimplements them, so
// they are checked directly rather than only through the fields that use them.
func TestConformanceHelpers(t *testing.T) {
	inputs := conformanceInputs(t)

	type record struct {
		In              string `json:"in"`
		Uuid            bool   `json:"uuid"`
		UuidV4          bool   `json:"uuidV4"`
		LegacyID        bool   `json:"legacyId"`
		LegacyIDLower   bool   `json:"legacyIdLower"`
		S12ID           bool   `json:"s12Id"`
		S12IDLower      bool   `json:"s12IdLower"`
		LongPrefixed    bool   `json:"longPrefixedLegacyId"`
		Email           bool   `json:"email"`
		NonEmail        bool   `json:"nonEmail"`
		URLHTTPS        bool   `json:"urlHttps"`
		URLHTTPSErr     string `json:"urlHttpsErr,omitempty"`
		URLMulti        bool   `json:"urlMulti"`
		URLMultiErr     string `json:"urlMultiErr,omitempty"`
		ContainsURL     bool   `json:"containsUrl"`
		BreakPartialURL string `json:"breakPartialUrl"`
		StripPUA        string `json:"stripPua"`
		ReplaceUnsafe   string `json:"replaceUnsafe"`
		ReplaceOther    string `json:"replaceOther"`
		ReplaceOtherML  string `json:"replaceOtherMultiline"`
	}

	https := []string{"https"}
	multi := []string{"https", "http"}

	records := make([]record, 0, len(inputs))
	for _, in := range inputs {
		rec := record{
			In:              encode(in),
			Uuid:            vproto.IsUUID(in),
			UuidV4:          vproto.IsUUIDv4(in),
			LegacyID:        vproto.IsLegacyID(in, false),
			LegacyIDLower:   vproto.IsLegacyID(in, true),
			S12ID:           vproto.IsS12ID(in, false),
			S12IDLower:      vproto.IsS12ID(in, true),
			LongPrefixed:    vproto.IsLongPrefixedLegacyID(in),
			Email:           vproto.IsValidEmail(in, false),
			NonEmail:        vproto.IsValidNonEmail(in),
			ContainsURL:     vproto.RejectURLMatcher.MatchString(in),
			BreakPartialURL: encode(vproto.BreakURLMatcher.ReplaceAllString(in, ". $1")),
			StripPUA:        encode(vproto.RegexPua.ReplaceAllString(in, "")),
			ReplaceUnsafe:   encode(vproto.UnsafeCharReplacer.Replace(in)),
			ReplaceOther:    encode(vproto.SymbolCharReplacer.Replace(in)),
			ReplaceOtherML:  encode(vproto.SymbolCharReplacerMultiline.Replace(in)),
		}

		if _, err := vproto.IsValidURL(in, https, false); err != nil {
			rec.URLHTTPSErr = err.Error()
		} else {
			rec.URLHTTPS = true
		}
		if _, err := vproto.IsValidURL(in, multi, true); err != nil {
			rec.URLMultiErr = err.Error()
		} else {
			rec.URLMulti = true
		}

		records = append(records, rec)
	}

	compareVectors(t, "helpers.jsonl", records)
}

// TestConformanceTables records the character substitutions each replacer makes,
// found by running every code point through it. A port reproduces the tables by
// hand, and this is what says the copy is complete.
func TestConformanceTables(t *testing.T) {
	type record struct {
		Replacer string `json:"replacer"`
		From     string `json:"from"`
		To       string `json:"to"`
	}

	var records []record
	for _, replacer := range []struct {
		name string
		r    *strings.Replacer
	}{
		{"unsafe", vproto.UnsafeCharReplacer},
		{"other", vproto.SymbolCharReplacer},
		{"otherMultiline", vproto.SymbolCharReplacerMultiline},
	} {
		for cp := rune(0); cp <= 0x10FFFF; cp++ {
			if cp >= 0xD800 && cp <= 0xDFFF {
				continue // a surrogate half is not a code point on its own
			}
			in := string(cp)
			out := replacer.r.Replace(in)
			if out == in {
				continue
			}
			records = append(records, record{
				Replacer: replacer.name,
				From:     fmt.Sprintf("%04X", cp),
				To:       encode(out),
			})
		}
	}

	compareVectors(t, "tables.jsonl", records)
}

// conformanceBase returns the message every field record starts from, and writes
// it to the vector directory so a port begins from the same bytes rather than a
// second transcription of the fixture.
//
// invalid_encoding_string is carried as text here. The fixture holds bytes that
// are not valid UTF-8 in that field, which no wire format will accept, and the
// field opts out of encoding validation so any value passes it. Every record
// overwrites the field under test, so the substitution reaches no result.
func conformanceBase(t *testing.T) *ValTestMessage {
	t.Helper()

	base := proto.Clone(&valMsg).(*ValTestMessage)
	base.InvalidEncodingString = "Accept invalid"

	encoded, err := proto.MarshalOptions{Deterministic: true}.Marshal(base)
	if err != nil {
		t.Fatalf("marshalling the base message: %v", err)
	}

	writeBase(t, "base.b64", encoded)
	return base
}

// writeBase records the bytes a port starts its records from, or checks them.
func writeBase(t *testing.T, name string, encoded []byte) {
	t.Helper()

	path := filepath.Join(conformanceDir, name)
	want := encode(string(encoded))
	if *updateVectors {
		if err := os.WriteFile(path, []byte(want+"\n"), 0o644); err != nil {
			t.Fatalf("writing %s: %v", path, err)
		}
		return
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v (run `go test ./... -update` to create it)", path, err)
	}
	if strings.TrimSpace(string(got)) != want {
		t.Errorf("the base message no longer matches %s; run `go test ./... -update` and review the diff", path)
	}
}

// conformanceInputs reads the shared corpus. A b64: line carries a value that
// cannot be written literally, such as one with edge whitespace or a line break.
func conformanceInputs(t *testing.T) []string {
	t.Helper()

	path := filepath.Join(conformanceDir, "inputs.txt")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	defer f.Close()

	var inputs []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if payload, ok := strings.CutPrefix(line, "b64:"); ok {
			raw, err := base64.StdEncoding.DecodeString(payload)
			if err != nil {
				t.Fatalf("%s: unreadable b64 line %q: %v", path, payload, err)
			}
			inputs = append(inputs, string(raw))
			continue
		}
		inputs = append(inputs, line)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if len(inputs) == 0 {
		t.Fatalf("%s holds no inputs", path)
	}
	return inputs
}

// encode carries a value through the vector file without the file's own line and
// escaping rules changing it.
func encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// compareVectors checks the records against the committed file, or rewrites it.
// Each record is a line of its own, so one changed behaviour shows as one line.
func compareVectors[T any](t *testing.T, name string, records []T) {
	t.Helper()

	var buf strings.Builder
	encoder := json.NewEncoder(&buf)
	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			t.Fatalf("encoding %s: %v", name, err)
		}
	}

	path := filepath.Join(conformanceDir, name)
	if *updateVectors {
		if err := os.WriteFile(path, []byte(buf.String()), 0o644); err != nil {
			t.Fatalf("writing %s: %v", path, err)
		}
		t.Logf("rewrote %s with %d records", path, len(records))
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v (run `go test ./... -update` to create it)", path, err)
	}
	if buf.String() != string(want) {
		t.Errorf("behaviour no longer matches %s; run `go test ./... -update` and review the diff", path)
	}
}
