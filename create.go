// implements the functions and data structures required to implement FT.CREATE
package grstack

import "github.com/goslogan/grstack/internal"

// SearchIndex defines an index to be created with FT.CREATE
type IndexOptions struct {
	On              string
	Prefixes        []string
	Filter          string
	Language        string
	LanguageField   string
	Score           float64
	ScoreField      string
	MaxTextFields   bool
	NoOffsets       bool
	Temporary       uint64 // If this is a temporary index, number of seconds until expiry
	NoHighlight     bool
	NoFields        bool
	NoFreqs         bool
	StopWords       []string
	UseStopWords    bool
	SkipInitialscan bool
	Schema          []SchemaAttribute
}

type TagAttribute struct {
	Name         string
	Alias        string
	Sortable     bool
	UnNormalized bool
	Separator    string
	CaseSenstive bool
}

type TextAttribute struct {
	Name         string
	Alias        string
	Sortable     bool
	UnNormalized bool
	Phonetic     string
	Weight       float32
	NoStem       bool
}

type NumericAttribute struct {
	Name         string
	Alias        string
	Sortable     bool
	UnNormalized bool
}

type SchemaAttribute interface {
	serialize() []interface{}
}

/* -- FLUENT INTERFACE to index options; schema options too simple to benefit -- */

// NewIndexOptions returns an initialised IndexOptions struct with defaults set
func NewIndexOptions() *IndexOptions {
	return &IndexOptions{
		On:    "hash", // Default
		Score: 1,      // Default
	}
}

// AddSchemaAttribute appends a schema attribute to the IndexOptions' Schema array
func (i *IndexOptions) AddSchemaAttribute(t SchemaAttribute) *IndexOptions {
	i.Schema = append(i.Schema, t)
	return i
}

// WithSchema sets the IndexOptions' Schema array to the provided values overwriting
// any existing schema.
func (i *IndexOptions) WithSchema(s []SchemaAttribute) *IndexOptions {
	i.Schema = s
	return i
}

// AddPrefix appends a prefix to the IndexOptions' Prefixes array
func (i *IndexOptions) AddPrefix(prefix string) *IndexOptions {
	i.Prefixes = append(i.Prefixes, prefix)
	return i
}

// WithPrefixes sets the IndexOptions' Prefix array to the provided values overwriting
// any existing prefixes.
func (i *IndexOptions) WithPrefixes(prefixes []string) *IndexOptions {
	i.Prefixes = prefixes
	return i
}

// WithFilter sets the IndexOptions' Filter field to the provided value
func (i *IndexOptions) WithFilter(filter string) *IndexOptions {
	i.Filter = filter
	return i
}

// WithLanguage sets the IndexOptions' Language field to the provided value, setting
// the default language for the index
func (i *IndexOptions) WithLanguage(language string) *IndexOptions {
	i.Language = language
	return i
}

// WithLanguageField sets the IndexOptions' LanguageField field to the provided value, setting
// the field definining language in the index
func (i *IndexOptions) WithLanguageField(field string) *IndexOptions {
	i.LanguageField = field
	return i
}

// WithScore sets the IndexOptions' Score field to the provided value, setting
// the default score for documents (this should be zero to 1.0 and is not
// checked)
func (i *IndexOptions) WithScore(score float64) *IndexOptions {
	i.Score = score
	return i
}

// WithScoreField sets the IndexOptions' ScoreField field to the provided value, setting
// the field defining document score in the index
func (i *IndexOptions) WithScoreField(field string) *IndexOptions {
	i.ScoreField = field
	return i
}

/*

	StopWords       []string
	UseStopWords    bool
	SkipInitialscan bool

*/

// WithMaxTextFields sets the IndexOptions' MaxTextFields field to true
func (i *IndexOptions) WithMaxTextFields() *IndexOptions {
	i.MaxTextFields = true
	return i
}

// WithNoOffsets sets the IndexOptions' NoOffsets field to true
func (i *IndexOptions) WithNoOffsets() *IndexOptions {
	i.NoOffsets = true
	return i
}

// AsTempoary sets the Temporary  field to the given number of seconds.
func (i *IndexOptions) AsTemporary(secs uint64) *IndexOptions {
	i.Temporary = secs
	return i
}

// WithNoHighlight sets the IndexOptions' NoHighlight field to true
func (i *IndexOptions) WithNoHighlight() *IndexOptions {
	i.NoHighlight = true
	return i
}

// WithNoHighlight sets the IndexOptions' NoFields field to true
func (i *IndexOptions) WithNoFields() *IndexOptions {
	i.NoFields = true
	return i
}

// WithNoHighlight sets the IndexOptions' NoFreqs field to true.
func (i *IndexOptions) WithNoFreqs() *IndexOptions {
	i.NoFreqs = true
	return i
}

// SkipInitialscan sets the IndexOptions' SkipInitialscan field to true.
func (i *IndexOptions) WithSkipInitialscan() *IndexOptions {
	i.SkipInitialscan = true
	return i
}

// AddStopWord appends a new stopword to the IndexOptions' stopwords array
// and sets UseStopWords to true
func (i *IndexOptions) AddStopWord(word string) *IndexOptions {
	i.StopWords = append(i.StopWords, word)
	i.UseStopWords = true
	return i
}

// WithStopWords sets the IndexOptions' StopWords array to a new value
// and sets UseStopWords to true if the array has any entries
func (i *IndexOptions) WithStopWords(words []string) *IndexOptions {
	i.StopWords = words
	i.UseStopWords = len(words) > 0
	return i
}

// WithNoStopWords sets IndexOptions' StopWords array to empty and
// sets UseStopWords to true to ensure the index uses no stopwords at all
func (i *IndexOptions) WithNoStopWords() *IndexOptions {
	i.UseStopWords = true
	i.StopWords = []string{}
	return i
}

/* ---- SERIALIZATION METHODS */

func (i *IndexOptions) serialize() []interface{} {

	args := []interface{}{"on", i.On}
	args = append(args, internal.SerializeCountedArgs("prefix", false, i.Prefixes)...)

	if i.Filter != "" {
		args = append(args, "filter", i.Filter)
	}

	if i.Language != "" {
		args = append(args, "language", i.Language)
	}

	if i.LanguageField != "" {
		args = append(args, "language_field", i.LanguageField)
	}

	args = append(args, "score", i.Score)

	if i.ScoreField != "" {
		args = append(args, "score_field", i.ScoreField)
	}

	if i.MaxTextFields {
		args = append(args, "maxtextfields")
	}

	if i.NoOffsets {
		args = append(args, "nooffsets")
	}

	if i.Temporary > 0 {
		args = append(args, "temporary", i.Temporary)
	}

	if i.NoHighlight && !i.NoOffsets {
		args = append(args, "nohl")
	}

	if i.NoFields {
		args = append(args, "nofields")
	}

	if i.NoFreqs {
		args = append(args, "nofreqs")
	}

	if i.UseStopWords {
		args = append(args, internal.SerializeCountedArgs("stopwords", true, i.StopWords)...)
	}

	if i.SkipInitialscan {
		args = append(args, "skipinitialscan")
	}

	schema := []interface{}{"schema"}

	for _, attrib := range i.Schema {
		schema = append(schema, attrib.serialize()...)
	}

	return append(args, schema...)
}

func (a NumericAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}
	attribs = append(attribs, "numeric")

	if a.Sortable {
		attribs = append(attribs, "sortable")
		if a.UnNormalized {
			attribs = append(attribs, "sortable", "unf")
		}
	}

	return attribs
}

func (a TagAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}
	attribs = append(attribs, "tag")
	if a.Sortable {
		attribs = append(attribs, "sortable")
		if a.UnNormalized {
			attribs = append(attribs, "unf")
		}
	}
	if a.Separator != "" {
		attribs = append(attribs, "separator", a.Separator)
	}
	if a.CaseSenstive {
		attribs = append(attribs, "casesensitive")
	}

	return attribs
}

func (a TextAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}

	attribs = append(attribs, "text")

	if a.Sortable {
		attribs = append(attribs, "sortable")
		if a.UnNormalized {
			attribs = append(attribs, "unf")
		}
	}
	if a.Phonetic != "" {
		attribs = append(attribs, "phonetic", a.Phonetic)
	}
	if a.NoStem {
		attribs = append(attribs, "nostem")
	}
	if a.Weight != 0 {
		attribs = append(attribs, "weight", a.Weight)
	}

	return attribs
}
