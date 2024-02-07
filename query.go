// query provides an interface to RedisSearch's query functionality.
package grsearch

import (
	"fmt"
	"math"
	"time"

	"github.com/goslogan/grsearch/internal"
)

type QueryOptions struct {
	NoContent    bool
	Verbatim     bool
	NoStopWords  bool
	WithScores   bool
	WithPayloads bool
	WithSortKeys bool
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
	GeoFilters   []GeoFilter
	Params       map[string]interface{}
	json         bool
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
	SortNone                 = "" // SortNone is used to indicate that no sorting is required if you want be explicit
	defaultDialect           = 2
)

/******************************************************************************
* Functions operating on the query struct itself							  *
******************************************************************************/

// NewQuery creates a new query with defaults set
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Limit:     NewLimit(DefaultOffset, DefaultLimit),
		Slop:      noSlop,
		SortOrder: SortAsc,
		Dialect:   defaultDialect,
		Params:    map[string]interface{}{},
	}
}

/******************************************************************************
* Functions operating on the return arguments								  *
******************************************************************************/

type QueryReturn struct {
	Name string
	As   string
}

// serialize converts a query struct to a slice of  interface{}
// ready for execution against Redis
func (q *QueryOptions) serialize() []interface{} {
	var args = []interface{}{}

	args = q.appendFlagArg(args, q.NoContent, "NOCONTENT")
	args = q.appendFlagArg(args, q.Verbatim, "VERBATIM")
	args = q.appendFlagArg(args, q.NoStopWords, "NOSTOPWORDS")
	args = q.appendFlagArg(args, q.WithScores, "WITHSCORES")
	args = q.appendFlagArg(args, q.WithPayloads, "WITHPAYLOADS")
	args = q.appendFlagArg(args, q.WithSortKeys, "WITHSORTKEYS")
	args = append(args, q.serializeFilters()...)
	for _, gf := range q.GeoFilters {
		args = append(args, gf.serialize()...)
	}
	args = append(args, q.serializeReturn()...)
	if q.Summarize != nil {
		args = append(args, q.Summarize.serialize()...)
	}
	if q.HighLight != nil {
		args = append(args, q.HighLight.serialize()...)
	}

	if q.Slop != noSlop {
		args = internal.AppendStringArg(args, "SLOP", fmt.Sprintf("%d", q.Slop))
	}

	if q.Timeout != 0 {
		args = internal.AppendStringArg(args, "TIMEOUT", fmt.Sprintf("%d", q.Timeout.Milliseconds()))
	}
	args = q.appendFlagArg(args, q.InOrder, "INORDER")
	args = internal.AppendStringArg(args, "LANGUAGE", q.Language)

	args = append(args, internal.SerializeCountedArgs("INKEYS", false, q.InKeys)...)
	args = append(args, internal.SerializeCountedArgs("INFIELDS", false, q.InFields)...)

	args = q.appendFlagArg(args, q.ExplainScore && q.WithScores, "EXPLAINSCORE")

	if q.SortBy != "" {
		args = append(args, "SORTBY", q.SortBy)
		if q.SortOrder != "" {
			args = append(args, q.SortOrder)
		}
	}

	if q.Limit != nil {
		args = append(args, q.Limit.serialize()...)
	}

	if len(q.Params) != 0 {
		args = append(args, "params", len(q.Params))
		for n, v := range q.Params {
			args = append(args, n, v)
		}
	}

	if q.Dialect != defaultDialect {
		args = append(args, "DIALECT", q.Dialect)
	}

	return args
}

func (q *QueryOptions) serializeReturn() []interface{} {
	if len(q.Return) > 0 {
		fields := []interface{}{}
		for _, ret := range q.Return {
			if ret.As == "" {
				fields = append(fields, ret.Name)
			} else {
				fields = append(fields, ret.Name, "AS", ret.As)
			}
		}
		return append([]interface{}{"RETURN", len(fields)}, fields...)
	} else {
		return nil
	}
}

// setResultSize uses the query to work out how many entries
// in the query raw results slice are used per result.
func (q *QueryOptions) resultSize() int {
	count := 2 // default to 2 - key and value

	if q.WithScores { // one more if returning scores
		count += 1
	}

	if q.NoContent { // one less if not content
		count -= 1
	}

	return count
}

// appendFlagArg appends the values to args if flag is true. args is returned
func (q *QueryOptions) appendFlagArg(args []interface{}, flag bool, value string) []interface{} {
	if flag {
		return append(args, value)
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

/*****
	QueryFilters
*****/

type QueryFilter struct {
	Attribute string
	Min       interface{} // either a numeric value or +inf, -inf or "(" followed by numeric
	Max       interface{} // as above
}

// NewQueryFilter returns a filter with the min and max properties to set + and - infinity.
func NewQueryFilter(attribute string) QueryFilter {
	qf := QueryFilter{Attribute: attribute}
	qf.Min = FilterValue(math.Inf(-1), false)
	qf.Max = FilterValue(math.Inf(1), false)
	return qf
}

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

// serialize converts a filter list to an array of interface{} objects for execution
func (q *QueryFilter) serialize() []interface{} {
	return []interface{}{"FILTER", q.Attribute, q.Min, q.Max}
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

// Serialize the limit for output in an FT.SEARCH
func (ql *Limit) serialize() []interface{} {
	if ql.Offset == DefaultOffset && ql.Num == DefaultLimit {
		return nil
	} else {
		return []interface{}{"LIMIT", ql.Offset, ql.Num}
	}
}

// Limit also needs to implement AggregateStep
func (ql *Limit) serializeStep() []interface{} {
	return ql.serialize()
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

// serialize prepares the summarisation to be passed to Redis.
func (s *QuerySummarize) serialize() []interface{} {
	args := []interface{}{"SUMMARIZE"}
	args = append(args, internal.SerializeCountedArgs("FIELDS", false, s.Fields)...)
	args = append(args, "FRAGS", s.Frags)
	args = append(args, "LEN", s.Len)
	args = append(args, "SEPARATOR", s.Separator)
	return args
}

// QueryHighlight allows the user to define optional query highlighting
type QueryHighlight struct {
	Fields   []string
	OpenTag  string
	CloseTag string
}

// serialize prepares the highlighting to be passed to Redis.
func (h *QueryHighlight) serialize() []interface{} {
	args := []interface{}{"HIGHLIGHT"}
	args = append(args, internal.SerializeCountedArgs("FIELDS", false, h.Fields)...)
	if h.OpenTag != "" || h.CloseTag != "" {
		args = append(args, "TAGS", h.OpenTag, h.CloseTag)
	}
	return args
}

/******************************************************************************
* Geofilters
******************************************************************************/

// GeoFilter represents a location and radius to be used in a search query
type GeoFilter struct {
	Attribute         string
	Long, Lat, Radius float64
	Units             string
}

func (gf *GeoFilter) serialize() []interface{} {
	return []interface{}{"GEOFILTER", gf.Attribute, gf.Long, gf.Lat, gf.Radius, gf.Units}
}
