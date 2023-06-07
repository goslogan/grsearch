package grstack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goslogan/grstack/internal"
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
	Weight         float32
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
	parse(key string, value interface{})
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

func (a *NumericAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}
	attribs = append(attribs, "numeric")

	if a.Sortable {
		attribs = append(attribs, "sortable")
	}

	if a.NoIndex {
		attribs = append(attribs, "noindex")
	}

	return attribs
}

func (a *TagAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}
	attribs = append(attribs, "tag")

	if a.Separator != "" {
		attribs = append(attribs, "separator", a.Separator)
	}

	if a.Sortable {
		attribs = append(attribs, "sortable")
		if a.UnNormalized {
			attribs = append(attribs, "unf")
		}
	}

	if a.CaseSensitive {
		attribs = append(attribs, "casesensitive")
	}
	if a.NoIndex {
		attribs = append(attribs, "noindex")
	}

	return attribs
}

func (a *TextAttribute) serialize() []interface{} {

	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}

	attribs = append(attribs, "text")

	if a.Weight != 0 {
		attribs = append(attribs, "weight", a.Weight)
	}

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

	if a.NoIndex {
		attribs = append(attribs, "noindex")
	}

	return attribs
}

func (a *GeometryAttribute) serialize() []interface{} {
	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}
	attribs = append(attribs, "geometry")

	return attribs
}

func (a *GeoAttribute) serialize() []interface{} {
	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}

	attribs = append(attribs, "geo")

	if a.Sortable {
		attribs = append(attribs, "sortable")
	}

	if a.NoIndex {
		attribs = append(attribs, "noindex")
	}
	return attribs
}

func (a *VectorAttribute) serialize() []interface{} {
	attribs := []interface{}{a.Name}
	if a.Alias != "" {
		attribs = append(attribs, "as", a.Alias)
	}

	attribs = append(attribs, "vector")
	attribs = append(attribs, a.Algorithm)

	params := []interface{}{"type", a.Type, "dim", a.Dim, "distance_metric", a.DistanceMetric}
	if a.InitialCap != 0 {
		params = append(params, "initial_cap", a.InitialCap)
	}
	if strings.ToLower(a.Algorithm) == "flat" && a.BlockSize != 0 {
		params = append(params, "block_size", a.BlockSize)
	}
	if strings.ToLower(a.Algorithm) == "hnsw" {
		if a.M != 0 {
			params = append(params, "m", a.M)
		}
		if a.EFConstruction != 0 {
			params = append(params, "ef_construction", a.EFConstruction)
		}
		if a.EFRuntime != 0 {
			params = append(params, "ef_runtime", a.EFRuntime)
		}
		if a.Epsilon != 0 {
			params = append(params, "epsilon", a.Epsilon)
		}
	}
	attribs = append(attribs, len(params))
	attribs = append(attribs, params...)

	return attribs
}

/******************************************************************************
* Create IndexOptions from ft.info output
******************************************************************************/
func (i *IndexOptions) parseInfo(input map[string]interface{}) error {

	if data, ok := input["index_definition"]; ok {
		mapped := internal.ToMap(data.([]interface{}))
		if mapped["key_type"] == "JSON" {
			i.On = "json"
		}
		i.Score, _ = mapped["default_score"].(float64)
		i.Prefix = make([]string, len(mapped["prefixes"].([]interface{})))
		for n, p := range mapped["prefixes"].([]interface{}) {
			i.Prefix[n] = p.(string)
		}
	}

	if i.Schema == nil {
		i.Schema = make([]SchemaAttribute, 0)
	}
	if data, ok := input["attributes"]; ok {
		for _, a := range data.([]interface{}) {
			attribInfo := a.([]interface{})
			attribType := "unidentifiable"
			for n, t := range attribInfo {
				if strings.ToLower(t.(string)) == "type" {
					attribType = attribInfo[n+1].(string)
					break
				}
			}
			var attribute SchemaAttribute
			switch strings.ToLower(attribType) {
			case "tag":
				attribute = &TagAttribute{}
			case "text":
				attribute = &TextAttribute{}
			case "numeric":
				attribute = &NumericAttribute{}
			case "geo":
				attribute = &GeoAttribute{}
			case "geometry":
				attribute = &GeometryAttribute{}
			case "vector":
				attribute = &VectorAttribute{}
			default:
				return fmt.Errorf("grstack: unhandled attribute type: %s", attribInfo[5].(string))
			}
			parseAttribute(attribInfo, attribute)
			i.Schema = append(i.Schema, attribute)
		}

	}
	return nil
}

func parseAttribute(info []interface{}, attrib SchemaAttribute) {
	var (
		paramKey string
		paramVal interface{}
		ok       bool
	)

	for n := 0; n < len(info); n++ {
		if paramKey, ok = info[n].(string); ok {
			if n+1 < len(info) {
				paramVal = info[n+1]
			}
			attrib.parse(strings.ToLower(paramKey), paramVal)
		}
	}
}

func (a *TagAttribute) parse(key string, val interface{}) {

	switch key {
	case "identifier":
		a.Name = val.(string)
	case "attribute":
		a.Alias = val.(string)
	case "separator":
		a.Separator = val.(string)
	case "sortable":
		a.Sortable = true
	case "unf":
		a.UnNormalized = true
	case "casesensitive":
		a.CaseSensitive = true
	case "withsuffixtrie":
		a.WithSuffixTrie = true
	case "noindex":
		a.NoIndex = true
	}

}

func (a *TextAttribute) parse(key string, val interface{}) {
	switch key {
	case "identifier":
		a.Name = val.(string)
	case "attribute":
		a.Alias = val.(string)
	case "sortable":
		a.Sortable = true
	case "unf":
		a.UnNormalized = true
	case "withsuffixtrie":
		a.WithSuffixTrie = true
	case "noindex":
		a.NoIndex = true
	case "weight":
		a.Weight = float32(val.(float64))
	case "phonetic":
		a.Phonetic = val.(string)
	case "nostem":
		a.NoStem = true
	}

}

func (a *NumericAttribute) parse(key string, val interface{}) {
	switch key {
	case "identifier":
		a.Name = val.(string)
	case "attribute":
		a.Alias = val.(string)
	case "sortable":
		a.Sortable = true
	case "noindex":
		a.NoIndex = true
	}
}

func (a *GeometryAttribute) parse(key string, val interface{}) {
	switch key {
	case "identifier":
		a.Name = val.(string)
	case "attribute":
		a.Alias = val.(string)
	}

}

func (a *GeoAttribute) parse(key string, val interface{}) {
	switch key {
	case "identifier":
		a.Name = val.(string)
	case "attribute":
		a.Alias = val.(string)
	case "sortable":
		a.Sortable = true
	case "noindex":
		a.NoIndex = true
	}
}

func (a *VectorAttribute) parse(key string, val interface{}) {
	switch key {
	case "flat", "hnsw":
		a.Algorithm = key
	case "type":
		a.Type = strings.ToLower(val.(string))
	case "dim":
		i, _ := strconv.ParseUint(val.(string), 10, 64)
		a.Dim = i
	case "distance_metric":
		a.DistanceMetric = strings.ToLower(val.(string))
	case "initial_cap":
		cap, _ := strconv.ParseUint(val.(string), 10, 64)
		a.InitialCap = cap
	case "block_size":
		size, _ := strconv.ParseUint(val.(string), 10, 64)
		a.BlockSize = size
	case "m":
		m, _ := strconv.ParseUint(val.(string), 10, 64)
		a.M = m
	case "ef_construction":
		ef, _ := strconv.ParseUint(val.(string), 10, 64)
		a.EFConstruction = ef
	case "ef_runtime":
		ef, _ := strconv.ParseUint(val.(string), 10, 64)
		a.EFRuntime = ef
	case "epsilon":
		ef, _ := strconv.ParseFloat(val.(string), 64)
		a.Epsilon = ef
	}
}
