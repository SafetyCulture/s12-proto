// Package plan turns the s12 validator options on a proto file into an ordered,
// language-neutral description of the checks and rewrites each field needs.
// Emitters render that description and do not read the options themselves, so
// every target language runs the same rules in the same order.
//
// Any changes made in this package will require Security Engineering consultation/review
package plan

import (
	"google.golang.org/protobuf/compiler/protogen"
)

// Unset marks an integer bound that carries no value.
const Unset = int64(-1)

// File is the plan for one proto file.
type File struct {
	Proto    *protogen.File
	Messages []*Message
}

// Message is the plan for one message and the messages declared inside it.
type Message struct {
	Proto      *protogen.Message
	IsMapEntry bool
	Fields     []*Field
	Nested     []*Message

	// LegacyRegex holds the fields carrying a (validator.regex) pattern, in
	// declaration order. The pattern is compiled into a package-level variable
	// rather than shared through the deduplicated pattern pool.
	LegacyRegex []*LegacyRegexField
}

// LegacyRegexField is one (validator.regex) pattern and the field it came from.
type LegacyRegexField struct {
	Proto   *protogen.Field
	Pattern string
}

// Field is the plan for one message field: the scaffolding that selects the
// values to check, followed by the steps to run against each of them.
type Field struct {
	Proto *protogen.Field

	// Name is the proto field name, which is what error messages report.
	Name string

	// RepeatedLenGte and RepeatedLenLte bound the number of elements, and are
	// Unset when absent. They apply to the field itself rather than its elements.
	RepeatedLenGte int64
	RepeatedLenLte int64

	// Loop runs Steps against each element instead of the field itself.
	Loop bool

	// Unsupported names a field shape that carries no validation. Steps is empty.
	Unsupported string

	// Oneof is set for a member of a oneof declared in the proto, which reaches
	// its value through the wrapper type. Synthetic oneofs are reported through
	// Optional instead.
	Oneof bool

	// Optional is set for a field with the optional keyword, whose value is a
	// pointer that must be present before Steps run.
	Optional bool

	Steps []Step
}

// Step is one validator option's worth of work against a single value.
type Step interface{ isStep() }

// LegacyString covers the field options that predate (validator.string):
// (validator.regex), (validator.uuid) and (validator.legacy_id).
type LegacyString struct {
	Optional   bool
	MatchRegex bool
	UUID       bool
	LegacyID   bool
	Length     *Length
}

// Length bounds the byte length of a value, optionally after trimming it for the
// comparison only.
type Length struct {
	TrimFirst bool
	Gte       int64
	Lte       int64
}

// SimpleString bounds the rune count of a value and does nothing else.
type SimpleString struct {
	Optional bool
	Bounded  bool
	MinLen   int32
	MaxLen   int32
	LogOnly  bool
}

// String covers (validator.string) and (validator.unsafe_string). The two share
// their machinery and differ only in which characters the plan allows through.
type String struct {
	Optional bool
	Ops      []Op
}

// Op is one rewrite or check inside a String step. Order is significant: a
// rewrite changes the value the later ops see, and the value is mutated in place
// on the message.
type Op interface{ isOp() }

// OpNormalizeNFC rewrites a decomposed value to its composed form.
type OpNormalizeNFC struct{ LogOnly bool }

// OpCheckEncoding rejects a value carrying the replacement character or any
// invalid UTF-8.
type OpCheckEncoding struct{ LogOnly bool }

// OpRejectURL rejects a value containing a URL anywhere in it.
type OpRejectURL struct{}

// OpBreakPartialURL inserts a space after each dot that joins two characters.
type OpBreakPartialURL struct{}

// OpTrim removes leading and trailing whitespace.
type OpTrim struct{}

// OpReplaceUnsafe swaps every character in the unsafe table for its alternative.
type OpReplaceUnsafe struct{}

// OpReplaceLiteral swaps one character for another.
type OpReplaceLiteral struct{ From, To rune }

// OpStripCR removes carriage returns, leaving line feeds in place.
type OpStripCR struct{}

// OpReplaceOther swaps every character in the rare-symbol table for its
// alternative. Multiline keeps line feeds.
type OpReplaceOther struct{ Multiline bool }

// OpStripPUA removes Private Use Area code points in the Basic Multilingual Plane.
type OpStripPUA struct{}

// OpLength bounds the size of a value, counted in runes when Runes is set and in
// bytes otherwise.
type OpLength struct {
	Runes   bool
	Min     uint32
	Max     uint32
	LogOnly bool
}

// OpHasPrefix requires the value to start with Prefix.
type OpHasPrefix struct{ Prefix string }

// OpMatchPattern requires every character of the value to be in the allow list
// the plan built for this field.
type OpMatchPattern struct {
	Pattern string
	LogOnly bool
}

// ID covers (validator.id): a UUID, plus whichever legacy formats are enabled.
type ID struct {
	Optional      bool
	Version       string
	Legacy        bool
	S12           bool
	LongPrefixed  bool
	LowercaseOnly bool
	LogOnly       bool
	Requirement   string
}

// Email requires a parsable email address.
type Email struct{ Optional bool }

// Username requires an email address, or a non-email username when AllowNonEmail
// is set.
type Username struct {
	Optional      bool
	AllowNonEmail bool
	LogOnly       bool
}

// URL requires an absolute URL using one of Schemes.
type URL struct {
	Optional      bool
	Schemes       []string
	AllowFragment bool
}

// Timezone requires an IANA TZ database name. A value that is not Optional is
// also required to be present.
type Timezone struct{ Optional bool }

// Bytes bounds a byte field and optionally requires it to be UUID-sized.
type Bytes struct {
	Optional bool
	UUIDSize bool
	Length   *Length
}

// Int bounds an integer field. Each bound is Unset when absent.
type Int struct {
	Optional bool
	Gt       int64
	Gte      int64
	Lt       int64
	Lte      int64
}

// Number covers (validator.number) for both integer and floating point fields.
// Min and Max carry the range bounds as written, and are empty when unbounded.
type Number struct {
	Optional  bool
	RejectNaN bool
	Min       string
	Max       string
	LogOnly   bool
}

// MessageField recurses into a nested message that implements the validator
// interface.
type MessageField struct {
	Required bool
	Repeated bool
}

// EnumField requires a non-zero enum value.
type EnumField struct{}

func (*LegacyString) isStep() {}
func (*SimpleString) isStep() {}
func (*String) isStep()       {}
func (*ID) isStep()           {}
func (*Email) isStep()        {}
func (*Username) isStep()     {}
func (*URL) isStep()          {}
func (*Timezone) isStep()     {}
func (*Bytes) isStep()        {}
func (*Int) isStep()          {}
func (*Number) isStep()       {}
func (*MessageField) isStep() {}
func (*EnumField) isStep()    {}

func (*OpNormalizeNFC) isOp()    {}
func (*OpCheckEncoding) isOp()   {}
func (*OpRejectURL) isOp()       {}
func (*OpBreakPartialURL) isOp() {}
func (*OpTrim) isOp()            {}
func (*OpReplaceUnsafe) isOp()   {}
func (*OpReplaceLiteral) isOp()  {}
func (*OpStripCR) isOp()         {}
func (*OpReplaceOther) isOp()    {}
func (*OpStripPUA) isOp()        {}
func (*OpLength) isOp()          {}
func (*OpHasPrefix) isOp()       {}
func (*OpMatchPattern) isOp()    {}
