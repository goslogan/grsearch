package grsearch

import (
	"strings"

	"github.com/goslogan/grsearch/internal"
)

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
	Name           string
	Alias          string
	Sortable       bool
	UnNormalized   bool
	Separator      string
	CaseSensitive  bool
	WithSuffixTrie bool
	NoIndex        bool
}

type TextAttribute struct {
	Name           string
	Alias          string
	Sortable       bool
	UnNormalized   bool
	Phonetic       string
	Weight         float64
	NoStem         bool
	WithSuffixTrie bool
	NoIndex        bool
}

type NumericAttribute struct {
	Name     string
	Alias    string
	Sortable bool
	NoIndex  bool
}

type GeoAttribute struct {
	Name     string
	Alias    string
	Sortable bool
	NoIndex  bool
}

type VectorAttribute struct {
	Name           string
	Alias          string
	Algorithm      string
	Type           string
	Dim            uint64
	DistanceMetric string
	InitialCap     uint64
	BlockSize      uint64
	M              uint64
	EFConstruction uint64
	EFRuntime      uint64
	Epsilon        float64
}

type GeometryAttribute struct {
	Name  string
	Alias string
}

type SchemaAttribute interface {
	serialize() []interface{}
	parseFromInfo(map[interface{}]interface{})
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

	args := []interface{}{"ON", strings.ToUpper(i.On)}
	args = append(args, internal.SerializeCountedArgs("PREFIX", false, i.Prefix)...)

	if i.Filter != "" {
		args = append(args, "FILTER", i.Filter)
	}

	if i.Language != "" {
		args = append(args, "LANGUAGE", i.Language)
	}

	if i.LanguageField != "" {
		args = append(args, "LANGUAGE_FIELD", i.LanguageField)
	}

	args = append(args, "SCORE", i.Score)

	if i.ScoreField != "" {
		args = append(args, "SCORE_FIELD", i.ScoreField)
	}

	if i.MaxTextFields {
		args = append(args, "MAXTEXTFIELDS")
	}

	if i.NoOffsets {
		args = append(args, "NOOFFSETS")
	}

	if i.Temporary > 0 {
		args = append(args, "TEMPORARY", i.Temporary)
	}

	if i.NoHighlight && !i.NoOffsets {
		args = append(args, "NOHL")
	}

	if i.NoFields {
		args = append(args, "NOFIELDS")
	}

	if i.NoFreqs {
		args = append(args, "NOFREQS")
	}

	if i.UseStopWords {
		args = append(args, internal.SerializeCountedArgs("STOPWORDS", true, i.StopWords)...)
	}

	if i.SkipInitialscan {
		args = append(args, "SKIPINITIALSCAN")
	}

	schema := []interface{}{"SCHEMA"}

	for _, attrib := range i.Schema {
		schema = append(schema, attrib.serialize()...)
	}

	return append(args, schema...)
}

func (a *NumericAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "AS", a.Alias)
	}
	attribs = append(attribs, "NUMERIC")

	if a.Sortable {
		attribs = append(attribs, "SORTABLE")
	}

	if a.NoIndex {
		attribs = append(attribs, "NOINDEX")
	}

	return attribs
}

func (a *TagAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "AS", a.Alias)
	}
	attribs = append(attribs, "TAG")

	if a.Separator != "" {
		attribs = append(attribs, "SEPARATOR", a.Separator)
	}

	if a.Sortable {
		attribs = append(attribs, "SORTABLE")
		if a.UnNormalized {
			attribs = append(attribs, "UNF")
		}
	}

	if a.CaseSensitive {
		attribs = append(attribs, "CASESENSITIVE")
	}
	if a.NoIndex {
		attribs = append(attribs, "NOINDEX")
	}

	return attribs
}

func (a *TextAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "AS", a.Alias)
	}

	attribs = append(attribs, "TEXT")

	if a.Weight != 0 {
		attribs = append(attribs, "WEIGHT", a.Weight)
	}

	if a.Sortable {
		attribs = append(attribs, "SORTABLE")
		if a.UnNormalized {
			attribs = append(attribs, "UNF")
		}
	}
	if a.Phonetic != "" {
		attribs = append(attribs, "PHONETIC", a.Phonetic)
	}
	if a.NoStem {
		attribs = append(attribs, "NOSTEM")
	}

	if a.NoIndex {
		attribs = append(attribs, "NOINDEX")
	}

	return attribs
}

func (a *GeometryAttribute) serialize() []interface{} {
	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "AS", a.Alias)
	}
	attribs = append(attribs, "GEOMETRY")

	return attribs
}

func (a *GeoAttribute) serialize() []interface{} {
	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "AS", a.Alias)
	}

	attribs = append(attribs, "GEO")

	if a.Sortable {
		attribs = append(attribs, "SORTABLE")
	}

	if a.NoIndex {
		attribs = append(attribs, "NOINDEX")
	}
	return attribs
}

func (a *VectorAttribute) serialize() []interface{} {
	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "AS", a.Alias)
	}

	attribs = append(attribs, "VECTOR")
	attribs = append(attribs, a.Algorithm)

	params := []interface{}{"TYPE", a.Type, "DIM", a.Dim, "DISTANCE_METRIC", a.DistanceMetric}
	if a.InitialCap != 0 {
		params = append(params, "INITIAL_CAP", a.InitialCap)
	}
	if strings.ToLower(a.Algorithm) == "FLAT" && a.BlockSize != 0 {
		params = append(params, "BLOCK_SIZE", a.BlockSize)
	}
	if strings.ToLower(a.Algorithm) == "HNSW" {
		if a.M != 0 {
			params = append(params, "M", a.M)
		}
		if a.EFConstruction != 0 {
			params = append(params, "EF_CONSTRUCTION", a.EFConstruction)
		}
		if a.EFRuntime != 0 {
			params = append(params, "EF_RUNTIME", a.EFRuntime)
		}
		if a.Epsilon != 0 {
			params = append(params, "EPSILON", a.Epsilon)
		}
	}
	attribs = append(attribs, len(params))
	attribs = append(attribs, params...)

	return attribs
}
