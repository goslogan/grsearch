package grstack

type IndexBuilder struct {
	opts IndexOptions
}

// NewIndexBuilder creats a new fluid builder for indexes
func NewIndexBuilder() *IndexBuilder {
	return &IndexBuilder{
		opts: *NewIndexOptions(),
	}
}

// Options returns the options struct built with the builder
func (a *IndexBuilder) Options() *IndexOptions {
	return &a.opts
}

// On indicates if the index is on hashes (default) or json
func (i *IndexBuilder) On(idxType string) *IndexBuilder {
	i.opts.On = idxType
	return i
}

// Schema appends a schema attribute to the IndexOptions' Schema array
func (i *IndexBuilder) Schema(t SchemaAttribute) *IndexBuilder {
	i.opts.Schema = append(i.opts.Schema, t)
	return i
}

// Prefix appends a prefix to the IndexOptions' Prefix array
func (i *IndexBuilder) Prefix(prefix string) *IndexBuilder {
	i.opts.Prefix = append(i.opts.Prefix, prefix)
	return i
}

// Filter sets the IndexOptions' Filter field to the provided value
func (i *IndexBuilder) Filter(filter string) *IndexBuilder {
	i.opts.Filter = filter
	return i
}

// Language sets the IndexOptions' Language field to the provided value, setting
// the default language for the index
func (i *IndexBuilder) Language(language string) *IndexBuilder {
	i.opts.Language = language
	return i
}

// LanguageField sets the IndexOptions' LanguageField field to the provided value, setting
// the field definining language in the index
func (i *IndexBuilder) LanguageField(field string) *IndexBuilder {
	i.opts.LanguageField = field
	return i
}

// Score sets the IndexOptions' Score field to the provided value, setting
// the default score for documents (this should be zero to 1.0 and is not
// checked)
func (i *IndexBuilder) Score(score float64) *IndexBuilder {
	i.opts.Score = score
	return i
}

// ScoreField sets the IndexOptions' ScoreField field to the provided value, setting
// the field defining document score in the index
func (i *IndexBuilder) ScoreField(field string) *IndexBuilder {
	i.opts.ScoreField = field
	return i
}

// MaxTextFields sets the IndexOptions' MaxTextFields field to true
func (i *IndexBuilder) MaxTextFields() *IndexBuilder {
	i.opts.MaxTextFields = true
	return i
}

// NoOffsets sets the IndexOptions' NoOffsets field to true
func (i *IndexBuilder) NoOffsets() *IndexBuilder {
	i.opts.NoOffsets = true
	return i
}

// Temporary sets the Temporary  field to the given number of seconds.
func (i *IndexBuilder) Temporary(secs uint64) *IndexBuilder {
	i.opts.Temporary = secs
	return i
}

// NoHighlight sets the IndexOptions' NoHighlight field to true
func (i *IndexBuilder) NoHighlight() *IndexBuilder {
	i.opts.NoHighlight = true
	return i
}

// NoFields sets the IndexOptions' NoFields field to true
func (i *IndexBuilder) NoFields() *IndexBuilder {
	i.opts.NoFields = true
	return i
}

// NoFreqs sets the IndexOptions' NoFreqs field to true.
func (i *IndexBuilder) NoFreqs() *IndexBuilder {
	i.opts.NoFreqs = true
	return i
}

// SkipInitialscan sets the IndexOptions' SkipInitialscan field to true.
func (i *IndexBuilder) SkipInitialscan() *IndexBuilder {
	i.opts.SkipInitialscan = true
	return i
}

// topWord appends a new stopword to the IndexOptions' stopwords array
// and sets UseStopWords to true
func (i *IndexBuilder) StopWord(word string) *IndexBuilder {
	i.opts.StopWords = append(i.opts.StopWords, word)
	i.opts.UseStopWords = true
	return i
}
