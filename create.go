// implements the functions and data structures required to implement FT.CREATE
package grstack

import "github.com/goslogan/grstack/internal"

// SearchIndex defines an index to be created with FT.CREATE
type IndexOptions struct {
	On              string
	Prefix          []string
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

// NewIndexOptions returns an initialised IndexOptions struct with defaults set
func NewIndexOptions() *IndexOptions {
	return &IndexOptions{
		On:    "hash", // Default
		Score: 1,      // Default
	}
}

/* ---- SERIALIZATION METHODS */

func (i *IndexOptions) serialize() []interface{} {

	args := []interface{}{"on", i.On}
	args = append(args, internal.SerializeCountedArgs("prefix", false, i.Prefix)...)

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
