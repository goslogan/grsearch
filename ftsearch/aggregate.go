package ftsearch

import "time"

// aggregate.go contains functionality used to implement FT.AGGREGATE

type AggregateOptions struct {
	Verbatim bool
	Load     []Load
	Timeout  time.Duration
	LoadAll  bool
	GroupBy  []AggregateGroupBy
	Apply    []AggregateApply
	Limit    *Limit
	Filter   string
	Cursor   *Cursor // nil means no cursor
	Params   map[string]interface{}
	Dialect  uint8
}

type AggregateGroupBy struct {
	parent     *AggregateOptions
	Properties []string
	Reducers   []Reducer
}

type Reducer struct {
	Name string
	As   string
	Args []string
}
type AggregateApply struct {
	Name string
	As   string
}

type Cursor struct {
	Count   int
	MaxIdle int
}

var LoadAll Load = Load{Name: "*"}

/******************************************************************************
* Methods operating on the aggregate struct itself and simple operations on *
* the struct																  *
******************************************************************************/

// NewAggregateOptions creates a new query with defaults set
func NewAggregateOptions() *AggregateOptions {
	return &AggregateOptions{
		Limit:   DefaultAggregateLimit(),
		Dialect: defaultDialect,
	}
}

// NewGroupBy returns a new group by struct with the parent set.
func (a *AggregateOptions) NewGroupBy() AggregateGroupBy {
	return AggregateGroupBy{parent: a}
}

// WithDialect sets the dialect option for the aggregate. It is NOT checked.
func (a *AggregateOptions) WithDialect(version uint8) *AggregateOptions {
	a.Dialect = version
	return a
}

// WithTimeout sets the timeout for the aggregate, overriding the dedault
func (q *AggregateOptions) WithTimeout(timeout time.Duration) *AggregateOptions {
	q.Timeout = timeout
	return q
}

/******************************************************************************
* Methods operating on Load arguments       								  *
******************************************************************************/

// Load represents parameters to the LOAD argument
type Load struct {
	Name string
	As   string
}

/******************************************************************************
* Methods operating on the parameters       								  *
******************************************************************************/

// AddParam sets the value of a aggregate parameter.
func (a *AggregateOptions) AddParam(name string, value interface{}) *AggregateOptions {
	a.Params[name] = value
	return a
}

// RemoveParam removes a parameter from aggregate options
func (a *AggregateOptions) RemoveParam(name string) *AggregateOptions {
	delete(a.Params, name)
	return a
}

// ClearParams clears all the currently set parameters
func (a *AggregateOptions) ClearParams() *AggregateOptions {
	a.Params = make(map[string]interface{}, 0)
	return a
}

// WithParams sets the current parameters
func (a *AggregateOptions) WithParams() *AggregateOptions {
	a.Params = make(map[string]interface{}, 0)
	return a
}
