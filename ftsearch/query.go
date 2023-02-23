// query provides an interface to RedisSearch's query functionality.
package ftsearch

import (
	"fmt"
	"math"
)

type QueryOptions struct {
	Index       string
	NoContent   bool
	Verbatim    bool
	NoStopWords bool
	WithScores  bool
	InOrder     bool
	// ExplainScore bool
	Limit        *queryLimit
	ReturnFields []string
	Filters      []QueryFilter
	InKeys       []string
	InFields     []string
	Language     string
	Slop         int32
	Summarize    *querySummarize
	HighLight    *queryHighlight
	resultSize   int
}

const (
	noSlop                   = -100 // impossible value for slop to indicate none set
	defaultOffset            = 0    // default first value for return offset
	defaultLimit             = 10   // default number of results to return
	defaultSumarizeSeparator = "..."
	defaultSummarizeLen      = 20
	defaultSummarizeFrags    = 3
)

type QueryResult struct {
	Score       float64
	Value       map[string]string
	Explanation []interface{}
}

/******************************************************************************
* Functions operating on the query struct itself							  *
******************************************************************************/

// NewQuery creates a new query with defaults set
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Limit: DefaultQueryLimit(),
		Slop:  noSlop,
	}
}

// String returns the serialized query as a single string. Any quoting
// required to use it in redis-cli is not done.
func (q *QueryOptions) String() string {
	return fmt.Sprintf("%v", q.serialize())
}

// WithIndex sets the index to be search on a query, returning the
// udpated query for chaining
func (q *QueryOptions) WithIndex(index string) *QueryOptions {
	q.Index = index
	return q
}

// WithLimit adds a limit to a query, returning the Query with
// the limit added (to allow chaining)
func (q *QueryOptions) WithLimit(first int64, num int64) *QueryOptions {
	q.Limit = NewQueryLimit(first, num)
	return q
}

// WithReturnFields sets the return fields, replacing any which
// might currently be set, returning the updated qry.
func (q *QueryOptions) WithReturnFields(fields []string) *QueryOptions {
	q.ReturnFields = fields
	return q
}

// AddReturnField appends a single field to the return fields,
// returning the updated query
func (q *QueryOptions) AddReturnField(field string) *QueryOptions {
	q.ReturnFields = append(q.ReturnFields, field)
	return q
}

// WithFilters sets the filters, replacing any which might
// be currently set, returning the updated query
func (q *QueryOptions) WithFilters(filters []QueryFilter) *QueryOptions {
	q.Filters = filters
	return q
}

// WithFilters sets the filters, replacing any which might
// be currently set, returning the updated query
func (q *QueryOptions) AddFilter(filter QueryFilter) *QueryOptions {
	q.Filters = append(q.Filters, filter)
	return q
}

// WithInKeys sets the keys to be searched, limiting the search
// to only these keys. The updated query is returned.
func (q *QueryOptions) WithInKeys(keys []string) *QueryOptions {
	q.InKeys = keys
	return q
}

// AddKey adds a single key to the keys to be searched, limiting the search
// to only these keys. The updated query is returned.
func (q *QueryOptions) AddKey(key string) *QueryOptions {
	q.InKeys = append(q.InKeys, key)
	return q
}

// WithInKeys sets the fields to be searched, limiting the search
// to only these fields. The updated query is returned.
func (q *QueryOptions) WithInFields(fields []string) *QueryOptions {
	q.InFields = fields
	return q
}

// AddField adds a single field to the fields to be searched in, limiting the search
// to only these fields. The updated query is returned.
func (q *QueryOptions) AddField(field string) *QueryOptions {
	q.InFields = append(q.InFields, field)
	return q
}

// WithSummarize sets the Summarize member of the query, returning the updated query.
func (q *QueryOptions) WithSummarize(s *querySummarize) *QueryOptions {
	q.Summarize = s
	return q
}

// WithHighlight sets the Highlight member of the query, returning the updated query.
func (q *QueryOptions) WithHighlight(h *queryHighlight) *QueryOptions {
	q.HighLight = h
	return q
}

// WithContent sets the NoContent flag to false.
func (q *QueryOptions) WithContent() *QueryOptions {
	q.NoContent = false
	return q
}

// WithoutContent sets the NoContent flag to true.
func (q *QueryOptions) WithoutContent() *QueryOptions {
	q.NoContent = true
	return q
}

// serialize converts a query struct to a slice of  interface{}
// ready for execution against Redis
func (q *QueryOptions) serialize() []interface{} {
	var args = []interface{}{}

	if q.NoContent {
		args = append(args, "NOCONTENT")
	}

	if q.Verbatim {
		args = append(args, "VERBATIM")
	}

	if q.NoStopWords {
		args = append(args, "NOSTOPWORDS")
	}

	if q.WithScores {
		args = append(args, "WITHSCORES")
	}

	for _, filter := range q.Filters {
		args = append(args, filter.serialize())
	}

	args = append(args, serializeCountedArgs("RETURN", false, q.ReturnFields)...)

	if q.Summarize != nil {
		args = append(args, q.Summarize.serialize()...)
	}

	if q.HighLight != nil {
		args = append(args, q.HighLight.serialize()...)
	}

	args = append(args, q.serializeSlop()...)

	if q.InOrder {
		args = append(args, "INORDER")
	}

	args = append(args, q.serializeLanguage()...)
	args = append(args, serializeCountedArgs("INKEYS", false, q.InKeys)...)
	args = append(args, serializeCountedArgs("INFIELDS", false, q.InFields)...)

	//if q.ExplainScore {
	//	args = append(args, INSCORE")
	//}

	if q.Limit != nil {
		args = append(args, q.Limit.serialize()...)
	}

	return args
}

func (q *QueryOptions) serializeSlop() []interface{} {
	if q.Slop != noSlop {
		return []interface{}{"SLOP", q.Slop}
	} else {
		return nil
	}
}

func (q *QueryOptions) serializeLanguage() []interface{} {
	if q.Language != "" {
		return []interface{}{"LANGUAGE", q.Language}
	} else {
		return nil
	}
}

// setResultSize uses the query to work out how many entries
// in the query raw results slice are used per result.
func (q *QueryOptions) setResultSize() {
	count := 2 // default to 2 - key and value

	if q.WithScores { // one more if returning scores
		count += 1
	}

	if q.NoContent { // one less if not content
		count -= 1
	}

	//if q.ExplainScore { // one more if explaining
	//		count += 1
	//	}

	q.resultSize = count
}

/******************************************************************************
* Public utilities                                                            *
******************************************************************************/

// FilterValue formats a value for use in a filter and returns it
func FilterValue(val float64, exclusive bool) interface{} {
	prefix := ""
	if exclusive {
		prefix = "("
	}

	if math.IsInf(val, -1) {
		return prefix + "-inf"
	} else if math.IsInf(val, 1) {
		return prefix + "+inf"
	} else {
		return fmt.Sprintf("%s%f", prefix, val)
	}
}

/*****
	QueryFilters
*****/

type QueryFilter struct {
	Attribute string
	Min       interface{} // either a numeric value or +inf, -inf or "(" followed by numeric
	Max       interface{} // as above
}

func NewQueryFilter(attribute string) *QueryFilter {
	qf := &QueryFilter{Attribute: attribute}
	return qf.
		WithMinInclusive(math.Inf(-1)).
		WithMaxInclusive(math.Inf(1))
}

// WithMinInclusive sets an inclusive minimum for the query filter value and
// returns it
func (qf *QueryFilter) WithMinInclusive(val float64) *QueryFilter {
	qf.Min = FilterValue(val, false)
	return qf
}

// WithMaxInclusive sets an inclusive maximum for the query filter value and
// returns it
func (qf *QueryFilter) WithMaxInclusive(val float64) *QueryFilter {
	qf.Max = FilterValue(val, false)
	return qf
}

// WithMinExclusive sets an exclusive minimum for the query filter value and
// returns it
func (qf *QueryFilter) WithMinExclusive(val float64) *QueryFilter {
	qf.Min = FilterValue(val, true)
	return qf
}

// WithMaxExclusive sets an exclusive maximum for the query filter value and
// returns it
func (qf *QueryFilter) WithMaxExclusive(val float64) *QueryFilter {
	qf.Max = FilterValue(val, true)
	return qf
}

// serialize converts a filter list to an array of interface{} objects for execution
func (q QueryFilter) serialize() []interface{} {
	return []interface{}{"FILTER", q.Attribute, q.Min, q.Max}
}

/******************************************************************************
* Functions operating on QueryLimit structs                                   *
******************************************************************************/

// queryLimit defines the results by offset and number.
type queryLimit struct {
	First int64
	Num   int64
}

// NewQueryLimit returns an initialized QueryLimit struct
func NewQueryLimit(first int64, num int64) *queryLimit {
	return &queryLimit{First: first, Num: num}
}

// DefaultQueryLimit returns an initialzied QueryLimit struct with the
// default limit range
func DefaultQueryLimit() *queryLimit {
	return NewQueryLimit(defaultOffset, defaultLimit)
}

// Serialize the limit for output
func (ql *queryLimit) serialize() []interface{} {
	if ql.First == defaultOffset && ql.Num == defaultLimit {
		return nil
	} else {
		return []interface{}{"LIMIT", ql.First, ql.Num}
	}
}

/********************************************************************
Functions and structs used to set up summarization and highlighting.
********************************************************************/

type querySummarize struct {
	Fields    []string
	Frags     int32
	Len       int32
	Separator string
}

func DefaultQuerySummarize() *querySummarize {
	return &querySummarize{
		Separator: defaultSumarizeSeparator,
		Len:       defaultSummarizeLen,
		Frags:     defaultSummarizeFrags,
	}
}

func NewQuerySummarize() *querySummarize {
	return &querySummarize{}
}

func NewQueryHighlight() *queryHighlight {
	return &queryHighlight{}
}

// WithLen sets the length of the query summarization fragment (in words)
// The modified struct is returned to support chaining
func (s *querySummarize) WithLen(len int32) *querySummarize {
	s.Len = len
	return s
}

// WithFrags sets the number of the fragements to create and return
// The modified struct is returned to support chaining
func (s *querySummarize) WithFrags(n int32) *querySummarize {
	s.Frags = n
	return s
}

// WithSeparator sets the fragment separator to be used.
// The modified struct is returned to support chaining
func (s *querySummarize) WithSeparator(sep string) *querySummarize {
	s.Separator = sep
	return s
}

// WithFields sets the fields to be summarized. Leaving it empty
// (the default) will cause all fields to be summarized
// The modified struct is returned to support chaining
func (s *querySummarize) WithFields(fields []string) *querySummarize {
	s.Fields = fields
	return s
}

// AddField adds a new field to the list of those to be summarised.
// The modified struct is returned to support chaining
func (s *querySummarize) AddField(field string) *querySummarize {
	s.Fields = append(s.Fields, field)
	return s
}

// serialize prepares the summarisation to be passed to Redis.
func (s *querySummarize) serialize() []interface{} {
	args := []interface{}{"SUMMARIZE"}
	args = append(args, serializeCountedArgs("fields", false, s.Fields)...)
	args = append(args, "FRAGS", s.Frags)
	args = append(args, "LEN", s.Len)
	args = append(args, "SEPARATOR", s.Separator)
	return args
}

// queryHighlight allows the user to define optional query highlighting
type queryHighlight struct {
	Fields   []string
	OpenTag  string
	CloseTag string
}

// WithFields sets the fields to be highlighting. Leaving it empty
// (the default) will cause all fields to be highlighted
// The modified struct is returned to support chaining
func (h *queryHighlight) WithFields(fields []string) *queryHighlight {
	h.Fields = fields
	return h
}

// AddField adds a new field to the list of those to be highlighted.
// The modified struct is returned to support chaining
func (h *queryHighlight) AddField(field string) *queryHighlight {
	h.Fields = append(h.Fields, field)
	return h
}

// SetTags sets the start and end tags. Both must be non empty or
// both empty. This is not enforced in this code to keep the API consistent
// but will lead to a Redis error if not set correctly.
func (h *queryHighlight) SetTags(open string, close string) *queryHighlight {
	h.OpenTag = open
	h.CloseTag = close
	return h
}

// serialize prepares the highlighting to be passed to Redis.
func (h *queryHighlight) serialize() []interface{} {
	args := []interface{}{"HIGHLIGHT"}
	args = append(args, serializeCountedArgs("fields", false, h.Fields)...)
	if h.OpenTag != "" || h.CloseTag != "" {
		args = append(args, "TAGS", h.OpenTag, h.CloseTag)
	}
	return args
}

/******************************************************************************
* Internal utilities                                                          *
******************************************************************************/