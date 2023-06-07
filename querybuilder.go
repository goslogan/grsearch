package grstack

import (
	"time"
)

type QueryBuilder struct {
	opts QueryOptions
}

// NewAggregateBuilder creats a new fluid builder for aggregates
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		opts: *NewQueryOptions(),
	}
}

// Options returns the options struct built with the builder
func (a *QueryBuilder) Options() *QueryOptions {
	return &a.opts
}

// Limit adds a limit to a query, returning the Query with
// the limit added (to allow chaining)
func (q *QueryBuilder) Limit(first int, num int) *QueryBuilder {
	q.opts.Limit = NewLimit(first, num)
	return q
}

// Dialect sets the dialect option for the query. It is NOT checked.
func (q *QueryBuilder) Dialect(version uint8) *QueryBuilder {
	q.opts.Dialect = version
	return q
}

// WithTimeout sets the timeout for the query, overriding the default
func (q *QueryBuilder) Timeout(timeout time.Duration) *QueryBuilder {
	q.opts.Timeout = timeout
	return q
}

// Return appends a field to the return fields list
func (q *QueryBuilder) Return(identifier string, alias string) *QueryBuilder {
	q.opts.Return = append(q.opts.Return, QueryReturn{Name: identifier, As: alias})
	return q
}

// Filter adds a filter to the current set
func (q *QueryBuilder) Filter(attribute string, min, max interface{}) *QueryBuilder {
	q.opts.Filters = append(q.opts.Filters, QueryFilter{
		Attribute: attribute,
		Min:       min,
		Max:       max,
	})
	return q
}

// InKeys sets the keys to be searched, limiting the search
// to only these keys.
func (q *QueryBuilder) InKeys(keys []string) *QueryBuilder {
	q.opts.InKeys = keys
	return q
}

// InFields adds a field to the INFIELDS list limiting the search
// to only the given fields.
func (q *QueryBuilder) InField(field string) *QueryBuilder {
	q.opts.InFields = append(q.opts.InFields, field)
	return q
}

// Summarize sets the Summarize member of the query,
func (q *QueryBuilder) Summarize(fields []string, separator string, length, fragments int32) *QueryBuilder {
	q.opts.Summarize = &QuerySummarize{
		Fields:    fields,
		Separator: separator,
		Len:       length,
		Frags:     fragments,
	}
	return q
}

// Highlight sets the Highlight member of the query
func (q *QueryBuilder) Highlight(fields []string) *QueryBuilder {
	q.opts.HighLight = &QueryHighlight{
		Fields: fields,
	}
	return q
}

// SortBy sets the value of the sortby option to the query.
func (q *QueryBuilder) SortBy(field string) *QueryBuilder {
	q.opts.SortBy = field
	return q
}

// Ascending sets the sort order of the query results to ascending if sortby is set
func (q *QueryBuilder) Ascending() *QueryBuilder {
	q.opts.SortOrder = SortAsc
	return q
}

// Descending sets the sort order of the query results to ascending if sortby is set
func (q *QueryBuilder) Descending() *QueryBuilder {
	q.opts.SortOrder = SortDesc
	return q
}

// NoContent sets the NoContent flag to true.
func (q *QueryBuilder) NoContent() *QueryBuilder {
	q.opts.NoContent = true
	return q
}

// WithScores sets the WITHSCORES option for searches
func (q *QueryBuilder) WithScores() *QueryBuilder {
	q.opts.WithScores = true
	return q
}

// ExplainScore sets the EXPLAINSCORE option for searches.
func (q *QueryBuilder) ExplainScore() *QueryBuilder {
	q.opts.ExplainScore = true
	return q
}

// WithPayloads sets the WITHPAYLOADS option for searches
func (q *QueryBuilder) WithPayloads() *QueryBuilder {
	q.opts.WithPayloads = true
	return q
}

// Verbatim disables stemming.
func (q *QueryBuilder) Verbatim() *QueryBuilder {
	q.opts.Verbatim = true
	return q
}

// Slop sets the slop length.
func (q *QueryBuilder) Slop(slop int8) *QueryBuilder {
	q.opts.Slop = slop
	return q
}

// NoStopWords disables stop word checking
func (q *QueryBuilder) NoStopWords() *QueryBuilder {
	q.opts.NoStopWords = true
	return q
}

// GeoFilter adds a geographic filter to the query
func (q *QueryBuilder) GeoFilter(attribute string, long, lat, radius float64, units string) *QueryBuilder {
	q.opts.GeoFilters = append(q.opts.GeoFilters, GeoFilter{
		Attribute: attribute,
		Long:      long,
		Lat:       lat,
		Radius:    radius,
		Units:     units,
	})
	return q
}

// Param sets or adds the value of a query parameter.
func (q *QueryBuilder) Param(name string, value interface{}) *QueryBuilder {
	q.opts.Params[name] = value
	return q
}

// Params sets the current set parameters
func (q *QueryBuilder) Params(params map[string]interface{}) *QueryBuilder {
	q.opts.Params = params
	return q
}
