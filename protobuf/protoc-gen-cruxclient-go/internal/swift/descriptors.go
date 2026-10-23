// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type (
	fieldType = descriptorpb.FieldDescriptorProto_Type
)

const (
	typeDouble   = descriptorpb.FieldDescriptorProto_TYPE_DOUBLE
	typeFloat    = descriptorpb.FieldDescriptorProto_TYPE_FLOAT
	typeInt64    = descriptorpb.FieldDescriptorProto_TYPE_INT64
	typeUint64   = descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt32    = descriptorpb.FieldDescriptorProto_TYPE_INT32
	typeFixed64  = descriptorpb.FieldDescriptorProto_TYPE_FIXED64
	typeFixed32  = descriptorpb.FieldDescriptorProto_TYPE_FIXED32
	typeBool     = descriptorpb.FieldDescriptorProto_TYPE_BOOL
	typeString   = descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeGroup    = descriptorpb.FieldDescriptorProto_TYPE_GROUP
	typeMessage  = descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	typeBytes    = descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint32   = descriptorpb.FieldDescriptorProto_TYPE_UINT32
	typeEnum     = descriptorpb.FieldDescriptorProto_TYPE_ENUM
	typeSfixed32 = descriptorpb.FieldDescriptorProto_TYPE_SFIXED32
	typeSfixed64 = descriptorpb.FieldDescriptorProto_TYPE_SFIXED64
	typeSint32   = descriptorpb.FieldDescriptorProto_TYPE_SINT32
	typeSint64   = descriptorpb.FieldDescriptorProto_TYPE_SINT64
)

// file is a proto file with the references between its descriptors resolved.
type file struct {
	proto      *descriptorpb.FileDescriptorProto
	name       string
	pkg        string
	edition    descriptorpb.Edition
	features   *descriptorpb.FeatureSet
	deps       []*file
	publicDeps []*file
	messages   []*message
	enums      []*enum
	extensions []*field
	locations  map[string]*descriptorpb.SourceCodeInfo_Location
}

type message struct {
	proto          *descriptorpb.DescriptorProto
	name           string
	fullName       string
	index          int
	file           *file
	containing     *message
	fields         []*field
	oneofs         []*oneof
	nested         []*message
	enums          []*enum
	extensions     []*field
	extensionRange []*descriptorpb.DescriptorProto_ExtensionRange
	reservedRanges [][2]int32
	reservedNames  []string
	features       *descriptorpb.FeatureSet
}

type enum struct {
	proto      *descriptorpb.EnumDescriptorProto
	name       string
	fullName   string
	index      int
	file       *file
	containing *message
	values     []*enumValue
	features   *descriptorpb.FeatureSet
}

type enumValue struct {
	proto    *descriptorpb.EnumValueDescriptorProto
	name     string
	fullName string
	number   int32
	index    int
	enum     *enum
}

type oneof struct {
	proto      *descriptorpb.OneofDescriptorProto
	name       string
	index      int
	containing *message
	fields     []*field
}

type field struct {
	proto          *descriptorpb.FieldDescriptorProto
	name           string
	fullName       string
	number         int32
	index          int
	typ            fieldType
	label          descriptorpb.FieldDescriptorProto_Label
	features       *descriptorpb.FeatureSet
	isExtension    bool
	defaultValue   *string
	jsonName       string
	oneofIndex     *int32
	proto3Optional bool
	file           *file
	containing     *message
	extScope       *message
	messageType    *message
	enumType       *enum
}

// swiftPrefix is the swift_prefix option, or nil when the file does not set it.
func (f *file) swiftPrefix() *string {
	if options := f.proto.GetOptions(); options != nil {
		return options.SwiftPrefix
	}
	return nil
}

func (f *file) isBundledProto() bool {
	return f.pkg == "google.protobuf" && bundledProtoFiles[f.name]
}

var bundledProtoFiles = setOf([]string{
	"google/protobuf/any.proto",
	"google/protobuf/api.proto",
	"google/protobuf/descriptor.proto",
	"google/protobuf/duration.proto",
	"google/protobuf/empty.proto",
	"google/protobuf/field_mask.proto",
	"google/protobuf/source_context.proto",
	"google/protobuf/struct.proto",
	"google/protobuf/timestamp.proto",
	"google/protobuf/type.proto",
	"google/protobuf/wrappers.proto",
})

// location returns the source information of a descriptor path.
func (f *file) location(path ...int32) *descriptorpb.SourceCodeInfo_Location {
	return f.locations[pathKey(path)]
}

func pathKey(path []int32) string {
	var sb strings.Builder
	for i, p := range path {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(int(p)))
	}
	return sb.String()
}

// registry resolves type names across the files of a request.
type registry struct {
	files    map[string]*file
	messages map[string]*message
	enums    map[string]*enum
}

// buildFiles converts the descriptors of a request into files. The input holds
// every file after the files it imports.
func buildFiles(protos []*descriptorpb.FileDescriptorProto) (*registry, error) {
	reg := &registry{files: map[string]*file{}, messages: map[string]*message{}, enums: map[string]*enum{}}
	for _, fp := range protos {
		f, err := reg.addFile(fp)
		if err != nil {
			return nil, fmt.Errorf("add file %s: %w", fp.GetName(), err)
		}
		reg.files[f.name] = f
	}
	return reg, nil
}

func (r *registry) addFile(fp *descriptorpb.FileDescriptorProto) (*file, error) {
	f := &file{proto: fp, name: fp.GetName(), pkg: fp.GetPackage()}
	switch {
	case fp.Edition != nil:
		f.edition = fp.GetEdition()
	case fp.GetSyntax() == "" || fp.GetSyntax() == "proto2":
		f.edition = descriptorpb.Edition_EDITION_PROTO2
	case fp.GetSyntax() == "proto3":
		f.edition = descriptorpb.Edition_EDITION_PROTO3
	default:
		return nil, fmt.Errorf("%s: unknown syntax %q", f.name, fp.GetSyntax())
	}
	if f.edition < descriptorpb.Edition_EDITION_PROTO2 || f.edition > descriptorpb.Edition_EDITION_2023 {
		return nil, fmt.Errorf("%s: edition %s is not supported", f.name, f.edition)
	}
	f.features = resolveFeatures(defaultFeatures(f.edition), fp.GetOptions().GetFeatures())

	for _, name := range fp.GetDependency() {
		dep, ok := r.files[name]
		if !ok {
			return nil, fmt.Errorf("%s: the import %s has no descriptor", f.name, name)
		}
		f.deps = append(f.deps, dep)
	}
	for _, i := range fp.GetPublicDependency() {
		if i < 0 || int(i) >= len(f.deps) {
			return nil, fmt.Errorf("%s: public dependency index %d is out of range", f.name, i)
		}
		f.publicDeps = append(f.publicDeps, f.deps[i])
	}
	f.locations = map[string]*descriptorpb.SourceCodeInfo_Location{}
	for _, loc := range fp.GetSourceCodeInfo().GetLocation() {
		f.locations[pathKey(loc.GetPath())] = loc
	}

	for i, ep := range fp.GetEnumType() {
		f.enums = append(f.enums, r.newEnum(ep, i, f, nil, fp.GetPackage(), f.features))
	}
	for i, mp := range fp.GetMessageType() {
		m, err := r.newMessage(mp, i, f, nil, fp.GetPackage(), f.features)
		if err != nil {
			return nil, fmt.Errorf("add message %s: %w", mp.GetName(), err)
		}
		f.messages = append(f.messages, m)
	}
	for i, xp := range fp.GetExtension() {
		f.extensions = append(f.extensions, newField(xp, i, f, nil, f.features, true))
	}
	for _, m := range f.messages {
		if err := r.bindMessage(m); err != nil {
			return nil, fmt.Errorf("bind message %s: %w", m.fullName, err)
		}
	}
	for _, x := range f.extensions {
		if err := r.bindField(x); err != nil {
			return nil, err
		}
	}
	return f, nil
}

func scoped(scope, name string) string {
	if scope == "" {
		return name
	}
	return scope + "." + name
}

func (r *registry) newEnum(ep *descriptorpb.EnumDescriptorProto, index int, f *file, containing *message, scope string, parent *descriptorpb.FeatureSet) *enum {
	e := &enum{
		proto: ep, name: ep.GetName(), fullName: scoped(scope, ep.GetName()), index: index,
		file: f, containing: containing,
		features: resolveFeatures(parent, ep.GetOptions().GetFeatures()),
	}
	for i, vp := range ep.GetValue() {
		e.values = append(e.values, &enumValue{
			proto: vp, name: vp.GetName(), fullName: e.fullName + "." + vp.GetName(),
			number: vp.GetNumber(), index: i, enum: e,
		})
	}
	r.enums[e.fullName] = e
	return e
}

func (r *registry) newMessage(mp *descriptorpb.DescriptorProto, index int, f *file, containing *message, scope string, parent *descriptorpb.FeatureSet) (*message, error) {
	m := &message{
		proto: mp, name: mp.GetName(), fullName: scoped(scope, mp.GetName()), index: index,
		file: f, containing: containing,
		features:       resolveFeatures(parent, mp.GetOptions().GetFeatures()),
		extensionRange: mp.GetExtensionRange(),
		reservedNames:  mp.GetReservedName(),
	}
	for _, rr := range mp.GetReservedRange() {
		m.reservedRanges = append(m.reservedRanges, [2]int32{rr.GetStart(), rr.GetEnd()})
	}
	oneofFeatures := make([]*descriptorpb.FeatureSet, len(mp.GetOneofDecl()))
	for i, op := range mp.GetOneofDecl() {
		oneofFeatures[i] = resolveFeatures(m.features, op.GetOptions().GetFeatures())
		m.oneofs = append(m.oneofs, &oneof{proto: op, name: op.GetName(), index: i, containing: m})
	}
	for i, ep := range mp.GetEnumType() {
		m.enums = append(m.enums, r.newEnum(ep, i, f, m, m.fullName, m.features))
	}
	for i, np := range mp.GetNestedType() {
		nested, err := r.newMessage(np, i, f, m, m.fullName, m.features)
		if err != nil {
			return nil, err
		}
		m.nested = append(m.nested, nested)
	}
	for i, fp := range mp.GetField() {
		oneofIndex := int(fp.GetOneofIndex())
		if fp.OneofIndex != nil && (oneofIndex < 0 || oneofIndex >= len(m.oneofs)) {
			return nil, fmt.Errorf("%s: field %s has oneof index %d out of range", f.name, fp.GetName(), oneofIndex)
		}
		parentFeatures := m.features
		if fp.OneofIndex != nil && !fp.GetProto3Optional() {
			parentFeatures = oneofFeatures[oneofIndex]
		}
		fld := newField(fp, i, f, m, parentFeatures, false)
		m.fields = append(m.fields, fld)
		if fp.OneofIndex != nil {
			m.oneofs[oneofIndex].fields = append(m.oneofs[oneofIndex].fields, fld)
		}
	}
	for i, xp := range mp.GetExtension() {
		m.extensions = append(m.extensions, newField(xp, i, f, m, m.features, true))
	}
	r.messages[m.fullName] = m
	return m, nil
}

func newField(fp *descriptorpb.FieldDescriptorProto, index int, f *file, containing *message, parent *descriptorpb.FeatureSet, isExtension bool) *field {
	fld := &field{
		proto: fp, name: fp.GetName(), number: fp.GetNumber(), index: index,
		typ: fp.GetType(), label: fp.GetLabel(), isExtension: isExtension, file: f,
		proto3Optional: fp.GetProto3Optional(), oneofIndex: fp.OneofIndex,
		jsonName: fp.GetJsonName(),
	}
	if fp.JsonName == nil {
		fld.jsonName = toJSONFieldName(fp.GetName())
	}
	if fp.DefaultValue != nil {
		fld.defaultValue = fp.DefaultValue
	}
	fld.features = resolveFieldFeatures(f.edition, fp, parent)
	if fp.GetType() == typeMessage && fld.features.GetMessageEncoding() == descriptorpb.FeatureSet_DELIMITED {
		fld.typ = typeGroup
	}
	if fld.features.GetFieldPresence() == descriptorpb.FeatureSet_LEGACY_REQUIRED && fp.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL {
		fld.label = descriptorpb.FieldDescriptorProto_LABEL_REQUIRED
	}
	if isExtension {
		fld.fullName = scoped(extensionScopeName(f, containing), fp.GetName())
		fld.extScope = containing
	} else {
		fld.containing = containing
		fld.fullName = containing.fullName + "." + fp.GetName()
	}
	return fld
}

func extensionScopeName(f *file, scope *message) string {
	if scope != nil {
		return scope.fullName
	}
	return f.pkg
}

func trimDot(name string) string { return strings.TrimPrefix(name, ".") }

func (r *registry) bindMessage(m *message) error {
	for _, n := range m.nested {
		if err := r.bindMessage(n); err != nil {
			return fmt.Errorf("bind nested message %s: %w", n.fullName, err)
		}
	}
	for _, f := range m.fields {
		if err := r.bindField(f); err != nil {
			return fmt.Errorf("bind field %s: %w", f.fullName, err)
		}
	}
	for _, x := range m.extensions {
		if err := r.bindField(x); err != nil {
			return fmt.Errorf("bind extension %s: %w", x.fullName, err)
		}
	}
	return nil
}

func (r *registry) bindField(f *field) error {
	if f.isExtension {
		extendee, ok := r.messages[trimDot(f.proto.GetExtendee())]
		if !ok {
			return fmt.Errorf("%s: unknown extendee %s", f.fullName, f.proto.GetExtendee())
		}
		f.containing = extendee
	}
	switch f.typ {
	case typeEnum:
		e, ok := r.enums[trimDot(f.proto.GetTypeName())]
		if !ok {
			return fmt.Errorf("%s: unknown enum %s", f.fullName, f.proto.GetTypeName())
		}
		f.enumType = e
	case typeMessage, typeGroup:
		m, ok := r.messages[trimDot(f.proto.GetTypeName())]
		if !ok {
			return fmt.Errorf("%s: unknown message %s", f.fullName, f.proto.GetTypeName())
		}
		f.messageType = m
		if f.typ == typeGroup && (m.proto.GetOptions().GetMapEntry() || (f.containing != nil && f.containing.proto.GetOptions().GetMapEntry())) {
			f.typ = typeMessage
		}
	}
	return nil
}

func (f *field) isRequired() bool {
	return f.features.GetFieldPresence() == descriptorpb.FeatureSet_LEGACY_REQUIRED
}

func (f *field) isRepeated() bool {
	return f.label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}

func (f *field) isPackable() bool {
	if !f.isRepeated() {
		return false
	}
	switch f.typ {
	case typeString, typeGroup, typeMessage, typeBytes:
		return false
	}
	return true
}

func (f *field) isPacked() bool {
	return f.isPackable() && f.features.GetRepeatedFieldEncoding() == descriptorpb.FeatureSet_PACKED
}

func (f *field) isMap() bool {
	return f.typ == typeMessage && f.messageType.proto.GetOptions().GetMapEntry()
}

func (f *field) hasPresence() bool {
	if f.isRepeated() {
		return false
	}
	switch f.typ {
	case typeGroup, typeMessage:
		return true
	}
	return f.isExtension || f.oneofIndex != nil || f.features.GetFieldPresence() != descriptorpb.FeatureSet_IMPLICIT
}

// realContainingOneof is the oneof of the field, unless the oneof is the synthetic one of a proto3 optional field.
func (f *field) realContainingOneof() *oneof {
	if f.oneofIndex == nil || f.isExtension {
		return nil
	}
	o := f.containing.oneofs[*f.oneofIndex]
	if o.isSynthetic() {
		return nil
	}
	return o
}

func (o *oneof) isSynthetic() bool {
	return len(o.fields) == 1 && o.fields[0].proto3Optional
}

func (m *message) realOneofs() []*oneof {
	var out []*oneof
	for _, o := range m.oneofs {
		if !o.isSynthetic() {
			out = append(out, o)
		}
	}
	return out
}

// isGroupLike models the group-like check of the protobuf compilers.
func (f *field) isGroupLike() bool {
	if f.typ != typeGroup {
		return false
	}
	if f.name != strings.ToLower(f.messageType.name) {
		return false
	}
	if f.isExtension {
		if f.extScope == nil {
			return f.messageType.file == f.file && f.messageType.containing == nil
		}
		return f.messageType.containing == f.extScope
	}
	return f.messageType.containing == f.containing
}

// namingBase is the name that the generated names of a field come from.
func (f *field) namingBase() string {
	if f.isGroupLike() {
		return f.messageType.name
	}
	return f.name
}

func (m *message) mapKeyAndValue() (key, value *field, ok bool) {
	if !m.proto.GetOptions().GetMapEntry() {
		return nil, nil, false
	}
	return m.fields[0], m.fields[1], true
}

func (m *message) useMessageSetWireFormat() bool {
	return m.proto.GetOptions().GetMessageSetWireFormat()
}

func (e *enum) isClosed() bool {
	return e.features.GetEnumType() == descriptorpb.FeatureSet_CLOSED
}

func (m *message) isWellKnownAny() bool {
	return m.fullName == "google.protobuf.Any"
}

// protoGenericType is the name part of the runtime methods that handle a field.
func (f *field) protoGenericType() string {
	switch f.typ {
	case typeDouble:
		return "Double"
	case typeFloat:
		return "Float"
	case typeInt64:
		return "Int64"
	case typeUint64:
		return "UInt64"
	case typeInt32:
		return "Int32"
	case typeFixed64:
		return "Fixed64"
	case typeFixed32:
		return "Fixed32"
	case typeBool:
		return "Bool"
	case typeString:
		return "String"
	case typeGroup:
		return "Group"
	case typeMessage:
		return "Message"
	case typeBytes:
		return "Bytes"
	case typeUint32:
		return "UInt32"
	case typeEnum:
		return "Enum"
	case typeSfixed32:
		return "SFixed32"
	case typeSfixed64:
		return "SFixed64"
	case typeSint32:
		return "SInt32"
	case typeSint64:
		return "SInt64"
	}
	panic("unknown field type")
}

// containsRequiredFields reports whether the message holds a required field at any depth.
func (m *message) containsRequiredFields() bool {
	seen := map[string]bool{}
	var helper func(*message) bool
	helper = func(d *message) bool {
		if seen[d.fullName] {
			return false
		}
		seen[d.fullName] = true
		if len(d.extensionRange) > 0 {
			return true
		}
		for _, f := range d.fields {
			if f.isRequired() {
				return true
			}
			if f.messageType != nil && helper(f.messageType) {
				return true
			}
		}
		return false
	}
	return helper(m)
}

func (m *message) needsFoundationImport() bool {
	for _, f := range m.fields {
		if f.typ == typeBytes {
			return true
		}
	}
	for _, x := range m.extensions {
		if x.typ == typeBytes {
			return true
		}
	}
	for _, n := range m.nested {
		if n.needsFoundationImport() {
			return true
		}
	}
	return false
}

func (f *file) needsFoundationImport() bool {
	for _, x := range f.extensions {
		if x.typ == typeBytes {
			return true
		}
	}
	for _, m := range f.messages {
		if m.needsFoundationImport() {
			return true
		}
	}
	return false
}

// extensionRanges returns the extension ranges sorted by start and merged where they touch.
func (m *message) normalizedExtensionRanges() [][2]int32 {
	ranges := make([][2]int32, 0, len(m.extensionRange))
	for _, r := range m.extensionRange {
		ranges = append(ranges, [2]int32{r.GetStart(), r.GetEnd()})
	}
	return sortAndMergeContinuous(ranges)
}

func sortAndMergeContinuous(ranges [][2]int32) [][2]int32 {
	ordered := append([][2]int32(nil), ranges...)
	// Insertion sort keeps equal starts in order, like a stable sort.
	for i := 1; i < len(ordered); i++ {
		for j := i; j > 0 && ordered[j][0] < ordered[j-1][0]; j-- {
			ordered[j], ordered[j-1] = ordered[j-1], ordered[j]
		}
	}
	if len(ordered) > 1 {
		for i := len(ordered) - 2; i >= 0; i-- {
			if ordered[i][1] == ordered[i+1][0] {
				ordered[i] = [2]int32{ordered[i][0], ordered[i+1][1]}
				ordered = append(ordered[:i+1], ordered[i+2:]...)
			}
		}
	}
	return ordered
}

// ambitiousExtensionRanges also merges ranges that have no field between them.
func (m *message) ambitiousExtensionRanges() [][2]int32 {
	merged := m.normalizedExtensionRanges()
	if len(merged) <= 1 {
		return merged
	}
	numbers := make([]int, 0, len(m.fields))
	for _, f := range m.fields {
		numbers = append(numbers, int(f.number))
	}
	// Descending order.
	for i := 1; i < len(numbers); i++ {
		for j := i; j > 0 && numbers[j] > numbers[j-1]; j-- {
			numbers[j], numbers[j-1] = numbers[j-1], numbers[j]
		}
	}
	pos := 0
	next := func() (int, bool) {
		if pos >= len(numbers) {
			return 0, false
		}
		n := numbers[pos]
		pos++
		return n, true
	}
	nextNumber, has := next()
	for has && int(merged[len(merged)-1][0]) < nextNumber {
		nextNumber, has = next()
	}
	for i := len(merged) - 2; i >= 0; i-- {
		if !has || int(merged[i][0]) > nextNumber {
			merged[i] = [2]int32{merged[i][0], merged[i+1][1]}
			merged = append(merged[:i+1], merged[i+2:]...)
		} else {
			for has && int(merged[i][0]) < nextNumber {
				nextNumber, has = next()
			}
		}
	}
	return merged
}

// jsonOptionNumber is the field number of the pb.enumvalue.json option, an extension of the enum value options.
const jsonOptionNumber = 998

// jsonOption reads the pb.enumvalue.json option. It reports false when the value has no option.
func (v *enumValue) jsonOption() (string, bool) {
	return readJSONOption(v.proto.GetOptions(), jsonOptionNumber, protoregistry.GlobalTypes)
}

// readJSONOption reads the string of the option from the extension with the given number. The
// option is an extension that a program can register, so the registered type comes first. The unknown
// fields of the options are the other source: a program can register the type after it read the options.
func readJSONOption(options *descriptorpb.EnumValueOptions, number protowire.Number, types interface {
	FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
}) (string, bool) {
	if options == nil {
		return "", false
	}
	if xt, err := types.FindExtensionByNumber(options.ProtoReflect().Descriptor().FullName(), number); err == nil {
		xd := xt.TypeDescriptor()
		if options.ProtoReflect().Has(xd) {
			if xd.Message() == nil {
				return "", true
			}
			message := options.ProtoReflect().Get(xd).Message()
			if field := message.Descriptor().Fields().ByNumber(1); field != nil && field.Kind() == protoreflect.StringKind {
				return message.Get(field).String(), true
			}
			return "", true
		}
	}
	return jsonOptionFromUnknown(options.ProtoReflect().GetUnknown(), number)
}

// jsonOptionFromUnknown finds the option in the wire bytes of the unknown fields.
func jsonOptionFromUnknown(raw []byte, number protowire.Number) (string, bool) {
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return "", false
		}
		raw = raw[n:]
		if num != number || typ != protowire.BytesType {
			m := protowire.ConsumeFieldValue(num, typ, raw)
			if m < 0 {
				return "", false
			}
			raw = raw[m:]
			continue
		}
		body, m := protowire.ConsumeBytes(raw)
		if m < 0 {
			return "", false
		}
		for len(body) > 0 {
			inner, innerType, k := protowire.ConsumeTag(body)
			if k < 0 {
				return "", false
			}
			body = body[k:]
			if inner == 1 && innerType == protowire.BytesType {
				s, l := protowire.ConsumeBytes(body)
				if l < 0 {
					return "", false
				}
				return string(s), true
			}
			l := protowire.ConsumeFieldValue(inner, innerType, body)
			if l < 0 {
				return "", false
			}
			body = body[l:]
		}
		return "", true
	}
	return "", false
}
