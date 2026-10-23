// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Command gensamples writes random encoded messages for the wire compatibility test. It reads a
// descriptor set, writes one sample file for each sample, and writes an index that names the
// Kotlin class, the Java class, the sample file and the full name of each sample. It also writes
// the list of proto files that define the types reachable from the root files.
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-crux-kotlin/internal/kotlin"
)

func main() {
	descPath := flag.String("desc", "", "descriptor set with imports")
	roots := flag.String("roots", "", "comma-separated proto files whose messages get samples")
	out := flag.String("out", "", "output directory")
	root := flag.String("kotlin_package", kotlin.DefaultPackage, "root package of the Kotlin classes")
	perMessage := flag.Int("n", 12, "samples for each message")
	seed := flag.Int64("seed", 1, "random seed")
	flag.Parse()

	raw, err := os.ReadFile(*descPath)
	check(err)
	set := &descriptorpb.FileDescriptorSet{}
	check(proto.Unmarshal(raw, set))
	files, err := protodesc.NewFiles(set)
	check(err)

	var rootFiles []protoreflect.FileDescriptor
	for _, name := range strings.Split(*roots, ",") {
		fd, err := files.FindFileByPath(name)
		check(err)
		rootFiles = append(rootFiles, fd)
	}

	reach := map[string]bool{}
	var index strings.Builder
	check(os.MkdirAll(*out, 0o755))
	r := rand.New(rand.NewSource(*seed))
	count := 0
	for _, fd := range rootFiles {
		collect(fd, reach)
		walkMessages(fd.Messages(), func(md protoreflect.MessageDescriptor) {
			for i := 0; i < *perMessage; i++ {
				m := dynamicpb.NewMessage(md)
				if i > 0 || proto.CheckInitialized(m) != nil {
					fill(r, m, 0)
				}
				if proto.CheckInitialized(m) != nil {
					continue
				}
				if i%4 == 3 {
					m.SetUnknown(unknownBytes(r))
				}
				data, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
				check(err)
				name := fmt.Sprintf("%05d.bin", count)
				count++
				check(os.WriteFile(filepath.Join(*out, name), data, 0o644))
				fmt.Fprintf(&index, "%s\t%s\t%s\t%s\n", kotlin.BinaryClassName(md, *root), javaClassName(md), name, md.FullName())
			}
		})
	}
	check(os.WriteFile(filepath.Join(*out, "index.tsv"), []byte(index.String()), 0o644))

	var closure []string
	for path := range reach {
		closure = append(closure, path)
	}
	sort.Strings(closure)
	check(os.WriteFile(filepath.Join(*out, "closure.txt"), []byte(strings.Join(closure, "\n")+"\n"), 0o644))
	fmt.Printf("%d samples, %d files in the closure\n", count, len(closure))
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func walkMessages(msgs protoreflect.MessageDescriptors, visit func(protoreflect.MessageDescriptor)) {
	for i := 0; i < msgs.Len(); i++ {
		md := msgs.Get(i)
		if md.IsMapEntry() {
			continue
		}
		visit(md)
		walkMessages(md.Messages(), visit)
	}
}

// collect adds the files that define the types the messages of a file refer to, until no file is
// added.
func collect(fd protoreflect.FileDescriptor, reach map[string]bool) {
	if reach[fd.Path()] {
		return
	}
	reach[fd.Path()] = true
	walkAll(fd.Messages(), func(md protoreflect.MessageDescriptor) {
		for i := 0; i < md.Fields().Len(); i++ {
			f := md.Fields().Get(i)
			switch {
			case f.Message() != nil:
				collect(f.Message().ParentFile(), reach)
			case f.Enum() != nil:
				collect(f.Enum().ParentFile(), reach)
			}
		}
	})
}

func walkAll(msgs protoreflect.MessageDescriptors, visit func(protoreflect.MessageDescriptor)) {
	for i := 0; i < msgs.Len(); i++ {
		visit(msgs.Get(i))
		walkAll(msgs.Get(i).Messages(), visit)
	}
}

// javaClassName returns the JVM binary name of the Java class of a message.
func javaClassName(md protoreflect.MessageDescriptor) string {
	opts := md.ParentFile().Options().(*descriptorpb.FileOptions)
	var chain []string
	var d protoreflect.Descriptor = md
	for {
		chain = append([]string{string(d.Name())}, chain...)
		parent, ok := d.Parent().(protoreflect.MessageDescriptor)
		if !ok {
			break
		}
		d = parent
	}
	pkg := opts.GetJavaPackage()
	if pkg == "" {
		pkg = string(md.ParentFile().Package())
	}
	if !opts.GetJavaMultipleFiles() {
		outer := opts.GetJavaOuterClassname()
		if outer == "" {
			base := strings.TrimSuffix(filepath.Base(md.ParentFile().Path()), ".proto")
			outer = upperCamel(base)
			if declares(md.ParentFile(), outer) {
				outer += "OuterClass"
			}
		}
		chain = append([]string{outer}, chain...)
	}
	if pkg == "" {
		return strings.Join(chain, "$")
	}
	return pkg + "." + strings.Join(chain, "$")
}

func upperCamel(s string) string {
	var b strings.Builder
	up := true
	for _, r := range s {
		if r == '_' {
			up = true
			continue
		}
		if up && r >= 'a' && r <= 'z' {
			r -= 'a' - 'A'
		}
		up = false
		b.WriteRune(r)
	}
	return b.String()
}

func fill(r *rand.Rand, m *dynamicpb.Message, depth int) {
	md := m.Descriptor()
	chosen := map[protoreflect.OneofDescriptor]protoreflect.FieldNumber{}
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if od := fd.ContainingOneof(); od != nil && !od.IsSynthetic() {
			if _, ok := chosen[od]; !ok {
				members := od.Fields()
				pick := r.Intn(members.Len() + 1)
				if pick == members.Len() {
					chosen[od] = 0
				} else {
					chosen[od] = members.Get(pick).Number()
				}
			}
			if chosen[od] != fd.Number() {
				continue
			}
		} else if fd.Cardinality() != protoreflect.Required && r.Intn(10) < 3 {
			continue
		}
		setField(r, m, fd, depth)
	}
}

func setField(r *rand.Rand, m *dynamicpb.Message, fd protoreflect.FieldDescriptor, depth int) {
	switch {
	case fd.IsMap():
		if depth >= 3 && fd.MapValue().Message() != nil {
			return
		}
		mp := m.Mutable(fd).Map()
		for i := r.Intn(4); i > 0; i-- {
			key := randomValue(r, mp.NewValue, fd.MapKey(), depth).MapKey()
			mp.Set(key, randomValue(r, mp.NewValue, fd.MapValue(), depth+1))
		}
	case fd.IsList():
		if depth >= 3 && fd.Message() != nil {
			return
		}
		list := m.Mutable(fd).List()
		for i := r.Intn(5); i > 0; i-- {
			list.Append(randomValue(r, list.NewElement, fd, depth+1))
		}
	default:
		if depth >= 3 && fd.Message() != nil && fd.Cardinality() != protoreflect.Required {
			return
		}
		m.Set(fd, randomValue(r, func() protoreflect.Value { return m.NewField(fd) }, fd, depth+1))
	}
}

func randomValue(r *rand.Rand, newMessage func() protoreflect.Value, fd protoreflect.FieldDescriptor, depth int) protoreflect.Value {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(r.Intn(2) == 0)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(randInt32(r))
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(uint32(randInt32(r)))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(randInt64(r))
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(uint64(randInt64(r)))
	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(float32(randFloat(r)))
	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(randFloat(r))
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(randString(r))
	case protoreflect.BytesKind:
		b := make([]byte, r.Intn(20))
		r.Read(b)
		return protoreflect.ValueOfBytes(b)
	case protoreflect.EnumKind:
		values := fd.Enum().Values()
		if r.Intn(8) == 0 {
			return openEnumValue(r)
		}
		return protoreflect.ValueOfEnum(values.Get(r.Intn(values.Len())).Number())
	case protoreflect.MessageKind:
		v := newMessage()
		if fd.Message() != nil {
			fill(r, v.Message().Interface().(*dynamicpb.Message), depth)
		}
		return v
	}
	panic("unsupported kind " + fd.Kind().String())
}

func randInt32(r *rand.Rand) int32 {
	switch r.Intn(6) {
	case 0:
		return 0
	case 1:
		return -1
	case 2:
		return math.MaxInt32
	case 3:
		return math.MinInt32
	case 4:
		return int32(r.Intn(300))
	}
	return int32(r.Uint32())
}

func randInt64(r *rand.Rand) int64 {
	switch r.Intn(6) {
	case 0:
		return 0
	case 1:
		return -1
	case 2:
		return math.MaxInt64
	case 3:
		return math.MinInt64
	case 4:
		return int64(r.Intn(300))
	}
	return int64(r.Uint64())
}

func randFloat(r *rand.Rand) float64 {
	switch r.Intn(6) {
	case 0:
		return 0
	case 1:
		return math.Copysign(0, -1)
	case 2:
		return math.Inf(1 - 2*r.Intn(2))
	case 3:
		return 1.5
	}
	return float64(float32(r.NormFloat64() * 1000))
}

var stringPieces = []string{"", "a", "hello", "h\u00e9llo", "\u65e5\u672c\u8a9e", "\U0001F600", "tab\there", "quote\"s", "$dollar", "line\nbreak"}

func randString(r *rand.Rand) string {
	var b strings.Builder
	for i := r.Intn(4); i >= 0; i-- {
		b.WriteString(stringPieces[r.Intn(len(stringPieces))])
	}
	return b.String()
}

// unknownBytes returns encoded fields that no message of the test schemas declares.
func unknownBytes(r *rand.Rand) []byte {
	var b []byte
	b = protowire.AppendTag(b, 100001, protowire.VarintType)
	b = protowire.AppendVarint(b, r.Uint64())
	b = protowire.AppendTag(b, 100002, protowire.Fixed32Type)
	b = protowire.AppendFixed32(b, r.Uint32())
	b = protowire.AppendTag(b, 100003, protowire.Fixed64Type)
	b = protowire.AppendFixed64(b, r.Uint64())
	b = protowire.AppendTag(b, 100004, protowire.BytesType)
	b = protowire.AppendString(b, randString(r))
	b = protowire.AppendTag(b, 100005, protowire.StartGroupType)
	b = protowire.AppendTag(b, 1, protowire.VarintType)
	b = protowire.AppendVarint(b, 7)
	b = protowire.AppendTag(b, 100005, protowire.EndGroupType)
	return b
}

// openEnumValue returns an enum number that no test enum declares.
func openEnumValue(r *rand.Rand) protoreflect.Value {
	return protoreflect.ValueOfEnum(protoreflect.EnumNumber(100 + r.Intn(50)))
}

// declares reports whether a file has a top-level message, enum or service of the given name.
func declares(fd protoreflect.FileDescriptor, name string) bool {
	n := protoreflect.Name(name)
	return fd.Messages().ByName(n) != nil || fd.Enums().ByName(n) != nil || fd.Services().ByName(n) != nil
}
