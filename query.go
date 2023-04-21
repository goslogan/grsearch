// query provides an interface to RedisSearch's query functionality.
package grstack

import (
	"fmt"
	"math"
	"time"

	"github.com/goslogan/redis-stack/internal"
)

type QueryOptions struct {
	NoContent    bool
	Verbatim     bool
	NoStopWords  bool
	Scores       bool // WITHSCORES but we need that for the API name
	Payloads     bool // WITHPAYLOADS but we need that for the API name
	SortKeys     bool // WITHSORTKEYS but we need that for the API name
	InOrder      bool
	ExplainScore bool
	Limit        *Limit
	Return       []QueryReturn
	Filters      []QueryFilter
	InKeys       []string
	InFields     []string
	Language     string
	Slop         int8
	Expander     string
	Scorer       string
	SortBy       string
	SortOrder    string
	Dialect      uint8
	Timeout      time.Duration
	Summarize    *QuerySummarize
	HighLight    *QueryHighlight
	GeoFilter    *GeoFilter
	Params       map[string]interface{}
	resultSize   int
}

const (
	noSlop                   = -100 // impossible value for slop to indicate none set
	DefaultOffset            = 0    // default first value for return offset
	DefaultLimit             = 10   // default number of results to return
	noLimit                  = 0
	defaultSumarizeSeparator = "..."
	defaultSummarizeLen      = 20
	defaultSummarizeFrags    = 3
	GeoMiles                 = "mi"
	GeoFeet                  = "f"
	GeoKilimetres            = "km"
	GeoMetres                = "m"
	SortAsc                  = "ASC"
	SortDesc                 = "DESC"
	defaultDialect           = 2
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
		Limit:   DefaultQueryLimit(),
		Slop:    noSlop,
		SortBy:  SortAsc,
		Dialect: defaultDialect,
	}
}

// String returns the serialized query as a single string. Any quoting
// required to use it in redis-cli is not done.
func (q *QueryOptions) String() string {
	return fmt.Sprintf("%v", q.serialize())
}

// WithLimit adds a limit to a query, returning the Query with
// the limit added (to allow chaining)
func (q *QueryOptions) WithLimit(first int64, num int64) *QueryOptions {
	q.Limit = NewLimit(first, num)
	return q
}

// WithDialect sets the dialect option for the query. It is NOT checked.
func (q *QueryOptions) WithDialect(version uint8) *QueryOptions {
	q.Dialect = version
	return q
}

// WithTimeout sets the timeout for the query, overriding the dedault
func (q *QueryOptions) WithTimeout(timeout time.Duration) *QueryOptions {
	q.Timeout = timeout
	return q
}

/******************************************************************************
* Functions operating on the return arguments								  *
******************************************************************************/

type QueryReturn struct {
	Name string
	As   string
}

// WithReturn sets the return fields, replacing any which
// might currently be set, returning the updated qry. The fields array
// should consist of pairs of strings (identifier & alias)
func (q *QueryOptions) WithReturn(fields []QueryReturn) *QueryOptions {
	q.Return = fields
	return q
}

// AddReturn appends a single field to the return fields,
// returning the updated query
func (q *QueryOptions) AddReturn(identifier string, alias string) *QueryOptions {
	q.Return = append(q.Return, QueryReturn{Name: identifier, As: alias})
	return q
}

// ClearReturn removes any return fields set and returns the updated query
func (q *QueryOptions) ClearReturn() *QueryOptions {
	q.Return = []QueryReturn{}
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
func (q *QueryOptions) WithSummarize(s *QuerySummarize) *QueryOptions {
	q.Summarize = s
	return q
}

// WithHighlight sets the Highlight member of the query, returning the updated query.
func (q *QueryOptions) WithHighlight(h *QueryHighlight) *QueryOptions {
	q.HighLight = h
	return q
}

// WithSortBy sets the value of the sortby option to the query.
func (q *QueryOptions) WithSortBy(field string) *QueryOptions {
	q.SortBy = field
	return q
}

// Ascending sets the sort order of the query results to ascending if sortby is set
func (q *QueryOptions) Ascending() *QueryOptions {
	q.SortOrder = SortAsc
	return q
}

// Descending sets the sort order of the query results to ascending if sortby is set
func (q *QueryOptions) Descending() *QueryOptions {
	q.SortOrder = SortDesc
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

// WithScores sets the WITHSCORES option for searches
func (q *QueryOptions) WithScores() *QueryOptions {
	q.Scores = true
	return q
}

// WithScores clears the WITHSCORES option for searches
func (q *QueryOptions) WithoutScores() *QueryOptions {
	q.Scores = false
	return q
}

// WithExplainScore sets the EXPLAINSCORE option for searches.
func (q *QueryOptions) WithExplainScore() *QueryOptions {
	q.ExplainScore = true
	return q
}

// WithoutExplainScore clears the EXPLAINSCORE option for searches.
func (q *QueryOptions) WithoutExplainScore() *QueryOptions {
	q.ExplainScore = false
	return q
}

// WithPayloads sets the PAYLOADS option for searches
func (q *QueryOptions) WithPayloads() *QueryOptions {
	q.Payloads = true
	return q
}

// WithooutPayloads sets the PAYLOADS option for searches
func (q *QueryOptions) WithoutPayloads() *QueryOptions {
	q.Payloads = false
	return q
}

// WithGeoFilter adds a geographic filter to the query
func (q *QueryOptions) WithGeoFilter(gf *GeoFilter) *QueryOptions {
	q.GeoFilter = gf
	return q
}

// WithoutGeoFilter removes a georgaphic filter if one is set
func (q *QueryOptions) WithoutGeoFilter() *QueryOptions {
	q.GeoFilter = nil
	return q
}

// AddParam sets the value of a query parameter.
func (q *QueryOptions) AddParam(name string, value interface{}) *QueryOptions {
	q.Params[name] = value
	return q
}

// RemoveParam removes a parameter from search options
func (q *QueryOptions) RemoveParam(name string) *QueryOptions {
	delete(q.Params, name)
	return q
}

// ClearParams clears all the current set parameters
func (q *QueryOptions) ClearParams() *QueryOptions {
	q.Params = make(map[string]interface{}, 0)
	return q
}

// WithParams sets the current set parameters
func (q *QueryOptions) WithParams() *QueryOptions {
	q.Params = make(map[string]interface{}, 0)
	return q
}

// serialize converts a query struct to a slice of  interface{}
// ready for execution against Redis
func (q *QueryOptions) serialize() []interface{} {
	var args = []interface{}{}

	args = q.appendFlagArg(args, q.NoContent, "nocontent")
	args = q.appendFlagArg(args, q.Verbatim, "verbatim")
	args = q.appendFlagArg(args, q.NoStopWords, "nostopwords")
	args = q.appendFlagArg(args, q.Scores, "withscores")
	args = q.appendFlagArg(args, q.Payloads, "withpayloads")
	args = q.appendFlagArg(args, q.SortKeys, "withsortkeys")
	args = append(args, q.serializeFilters()...)
	if q.GeoFilter != nil {
		args = append(args, q.GeoFilter.serialize()...)

	}
	args = append(args, q.serializeReturn()...)
	if q.Summarize != nil {
		args = append(args, q.Summarize.serialize()...)
	}
	if q.HighLight != nil {
		args = append(args, q.HighLight.serialize()...)
	}

	if q.Slop != noSlop {
		args = q.appendStringArg(args, "slop", fmt.Sprintf("%d", q.Slop))
	}
	args = q.appendFlagArg(args, q.InOrder, "inorder")
	args = q.appendStringArg(args, "language", q.Language)

	args = append(args, internal.SerializeCountedArgs("inkeys", false, q.InKeys)...)
	args = append(args, internal.SerializeCountedArgs("infields", false, q.InFields)...)

	args = q.appendFlagArg(args, q.ExplainScore && q.Scores, "EXPLAINSCORE")

	if q.Limit != nil {
		args = append(args, q.Limit.serializeForSearch()...)
	}

	return args
}

func (q *QueryOptions) serializeReturn() []interface{} {
	if len(q.Return) > 0 {
		fields := []interface{}{"return", len(q.Return)}
		for _, ret := range q.Return {
			if ret.As == "" {
				fields = append(fields, ret.Name)
			} else {
				fields = append(fields, ret.Name, "as", ret.As)
			}
		}
		return fields
	} else {
		return nil
	}
}

// setResultSize uses the query to work out how many entries
// in the query raw results slice are used per result.
func (q *QueryOptions) setResultSize() {
	count := 2 // default to 2 - key and value

	if q.Scores { // one more if returning scores
		count += 1
	}

	if q.NoContent { // one less if not content
		count -= 1
	}

	q.resultSize = count
}

// appendFlagArg appends the values to args if flag is true. args is returned
func (q *QueryOptions) appendFlagArg(args []interface{}, flag bool, value string) []interface{} {
	if flag {
		return append(args, value)
	} else {
		return args
	}
}

// appendStringArg appends the name and value if value is not empty
func (q *QueryOptions) appendStringArg(args []interface{}, name, value string) []interface{} {
	if value != "" {
		return append(args, name, value)
	} else {
		return args
	}
}

// serialize the filters
func (q *QueryOptions) serializeFilters() []interface{} {
	args := []interface{}{}
	for _, f := range q.Filters {
		args = append(args, f.serialize()...)
	}
	return args
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

func NewQueryFilter(attribute string) QueryFilter {
	qf := &QueryFilter{Attribute: attribute}
	return qf.
		WithMinInclusive(math.Inf(-1)).
		WithMaxInclusive(math.Inf(1))
}

// WithMinInclusive sets an inclusive minimum for the query filter value and
// returns it
func (qf QueryFilter) WithMinInclusive(val float64) QueryFilter {
	qf.Min = FilterValue(val, false)
	return qf
}

// WithMaxInclusive sets an inclusive maximum for the query filter value and
// returns it
func (qf QueryFilter) WithMaxInclusive(val float64) QueryFilter {
	qf.Max = FilterValue(val, false)
	return qf
}

// WithMinExclusive sets an exclusive minimum for the query filter value and
// returns it
func (qf QueryFilter) WithMinExclusive(val float64) QueryFilter {
	qf.Min = FilterValue(val, true)
	return qf
}

// WithMaxExclusive sets an exclusive maximum for the query filter value and
// returns it
func (qf QueryFilter) WithMaxExclusive(val float64) QueryFilter {
	qf.Max = FilterValue(val, true)
	return qf
}

// serialize converts a filter list to an array of interface{} objects for execution
func (q QueryFilter) serialize() []interface{} {
	return []interface{}{"filter", q.Attribute, q.Min, q.Max}
}

/******************************************************************************
* Functions operating on QueryLimit structs                                   *
******************************************************************************/

// queryLimit defines the results by offset and number.
type Limit struct {
	Offset int64
	Num    int64
}

// NewQueryLimit returns an initialized limit struct
func NewLimit(first int64, num int64) *Limit {
	return &Limit{Offset: first, Num: num}
}

// DefaultQueryLimit returns an initialzied limit struct with the
// default limit range for use in FT.SEARCH
func DefaultQueryLimit() *Limit {
	return NewLimit(DefaultOffset, DefaultLimit)
}

// DefaultAggregateLimit returns an blank limit struct for use
// in FT.AGGREGATE
func DefaultAggregateLimit() *Limit {
	return NewLimit(DefaultOffset, noLimit)
}

// Serialize the limit for output in an FT.SEARCH
func (ql *Limit) serializeForSearch() []interface{} {
	if ql.Offset == DefaultOffset && ql.Num == DefaultLimit {
		return nil
	} else {
		return []interface{}{"limit", ql.Offset, ql.Num}
	}
}

// Serialize the limit for output in an FT.AGGREGATE
func (ql *Limit) serializeForAggregate() []interface{} {
	if ql.Offset == DefaultOffset && ql.Num == noLimit {
		return nil
	} else {
		return []interface{}{"limit", ql.Offset, ql.Num}
	}
}

/********************************************************************
Functions and structs used to set up summarization and highlighting.
********************************************************************/

type QuerySummarize struct {
	Fields    []string
	Frags     int32
	Len       int32
	Separator string
}

func DefaultQuerySummarize() *QuerySummarize {
	return &QuerySummarize{
		Separator: defaultSumarizeSeparator,
		Len:       defaultSummarizeLen,
		Frags:     defaultSummarizeFrags,
	}
}

func NewQuerySummarize() *QuerySummarize {
	return &QuerySummarize{}
}

func NewQueryHighlight() *QueryHighlight {
	return &QueryHighlight{}
}

// WithLen sets the length of the query summarization fragment (in words)
// The modified struct is returned to support chaining
func (s *QuerySummarize) WithLen(len int32) *QuerySummarize {
	s.Len = len
	return s
}

// WithFrags sets the number of the fragements to create and return
// The modified struct is returned to support chaining
func (s *QuerySummarize) WithFrags(n int32) *QuerySummarize {
	s.Frags = n
	return s
}

// WithSeparator sets the fragment separator to be used.
// The modified struct is returned to support chaining
func (s *QuerySummarize) WithSeparator(sep string) *QuerySummarize {
	s.Separator = sep
	return s
}

// WithFields sets the fields to be summarized. Leaving it empty
// (the default) will cause all fields to be summarized
// The modified struct is returned to support chaining
func (s *QuerySummarize) WithFields(fields []string) *QuerySummarize {
	s.Fields = fields
	return s
}

// AddField adds a new field to the list of those to be summarised.
// The modified struct is returned to support chaining
func (s *QuerySummarize) AddField(field string) *QuerySummarize {
	s.Fields = append(s.Fields, field)
	return s
}

// serialize prepares the summarisation to be passed to Redis.
func (s *QuerySummarize) serialize() []interface{} {
	args := []interface{}{"summarize"}
	args = append(args, internal.SerializeCountedArgs("fields", false, s.Fields)...)
	args = append(args, "frags", s.Frags)
	args = append(args, "len", s.Len)
	args = append(args, "separator", s.Separator)
	return args
}

// QueryHighlight allows the user to define optional query highlighting
type QueryHighlight struct {
	Fields   []string
	OpenTag  string
	CloseTag string
}

// WithFields sets the fields to be highlighting. Leaving it empty
// (the default) will cause all fields to be highlighted
// The modified struct is returned to support chaining
func (h *QueryHighlight) WithFields(fields []string) *QueryHighlight {
	h.Fields = fields
	return h
}

// AddField adds a new field to the list of those to be highlighted.
// The modified struct is returned to support chaining
func (h *QueryHighlight) AddField(field string) *QueryHighlight {
	h.Fields = append(h.Fields, field)
	return h
}

// SetTags sets the start and end tags. Both must be non empty or
// both empty. This is not enforced in this code to keep the API consistent
// but will lead to a Redis error if not set correctly.
func (h *QueryHighlight) SetTags(open string, close string) *QueryHighlight {
	h.OpenTag = open
	h.CloseTag = close
	return h
}

// serialize prepares the highlighting to be passed to Redis.
func (h *QueryHighlight) serialize() []interface{} {
	args := []interface{}{"HIGHLIGHT"}
	args = append(args, internal.SerializeCountedArgs("fields", false, h.Fields)...)
	if h.OpenTag != "" || h.CloseTag != "" {
		args = append(args, "tags", h.OpenTag, h.CloseTag)
	}
	return args
}

/******************************************************************************
* Parameters
******************************************************************************/

/******************************************************************************
* Geofilters
******************************************************************************/

// GeoFilter represents a location and radius to be used in a search query
type GeoFilter struct {
	Attribute         string
	Long, Lat, Radius float64
	Units             string
}

func NewGeoFilter(Attribute string, Long, Lat, Radius float64, Units string) *GeoFilter {
	return &GeoFilter{
		Attribute: Attribute,
		Long:      Long,
		Lat:       Lat,
		Radius:    Radius,
		Units:     Units,
	}
}

func (gf *GeoFilter) serialize() []interface{} {
	return []interface{}{"geofilter", gf.Attribute, gf.Long, gf.Lat, gf.Radius, gf.Units}
}

/******************************************************************************
* Internal utilities                                                          *
******************************************************************************/
