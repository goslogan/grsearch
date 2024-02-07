package grsearch

import "time"

type AggregateBuilder struct {
	opts AggregateOptions
}

type GroupByBuilder struct {
	group AggregateGroupBy
}

// NewAggregateBuilder creats a new fluid builder for aggregates
func NewAggregateBuilder() *AggregateBuilder {
	return &AggregateBuilder{
		opts: *NewAggregateOptions(),
	}
}

// Options returns the options struct built with the builder
func (a *AggregateBuilder) Options() *AggregateOptions {
	return &a.opts
}

// Dialect sets the dialect option for the aggregate. It is NOT checked.
func (a *AggregateBuilder) Dialect(version uint8) *AggregateBuilder {
	a.opts.Dialect = version
	return a
}

// Timeout sets the timeout for the aggregate, overriding the default
func (a *AggregateBuilder) Timeout(timeout time.Duration) *AggregateBuilder {
	a.opts.Timeout = timeout
	return a
}

// Param sets the value of a aggregate parameter.
func (a *AggregateBuilder) Param(name string, value interface{}) *AggregateBuilder {
	a.opts.Params[name] = value
	return a
}

// Params sets all current parameters
func (a *AggregateBuilder) Params(params map[string]interface{}) *AggregateBuilder {
	a.opts.Params = params
	return a
}

// Verbatim sets the verbatim flag, disabling stemming
func (a *AggregateBuilder) Verbatim() *AggregateBuilder {
	a.opts.Verbatim = true
	return a
}

// Limit sets the result limit
func (a *AggregateBuilder) Limit(offset, num int64) *AggregateBuilder {
	a.opts.Steps = append(a.opts.Steps, &Limit{Offset: offset, Num: num})
	return a
}

// Filter adds a result filter
func (a *AggregateBuilder) Filter(filter string) *AggregateBuilder {
	a.opts.Steps = append(a.opts.Steps, AggregateFilter(filter))
	return a
}

// Apply appends a transform to the apply list
func (a *AggregateBuilder) Apply(expression, name string) *AggregateBuilder {
	a.opts.Steps = append(a.opts.Steps, &AggregateApply{
		Expression: expression,
		As:         name,
	})
	return a
}

// WithCursor creates a cursor for the aggregate to scan parts of the result
func (a *AggregateBuilder) Cursor(count uint64, timeout time.Duration) *AggregateBuilder {
	a.opts.Cursor = &AggregateCursor{
		Count:   count,
		MaxIdle: timeout,
	}
	return a
}

// Load adds a field to the load list for the aggregate. The alias can be the
// empty string.
func (a *AggregateBuilder) Load(name string, as string) *AggregateBuilder {
	l := AggregateLoad{Name: name, As: as}
	a.opts.Load = append(a.opts.Load, l)
	return a
}

// LoadAll sets the load list for this aggregate to "LOAD *".
func (a *AggregateBuilder) LoadAll() *AggregateBuilder {
	a.opts.Load = []AggregateLoad{LoadAll}
	return a
}

// SortBy adds a sorting step to this aggregate.
func (a *AggregateBuilder) SortBy(keys []AggregateSortKey) *AggregateBuilder {

	a.opts.Steps = append(a.opts.Steps, &AggregateSort{
		Keys: keys,
	})

	return a
}

// SortByMax sets the MAX limit on an aggregate sort key. This will be
// ignored if not sort keys have been supplied.
func (a *AggregateBuilder) SortByMax(keys []AggregateSortKey, max int64) *AggregateBuilder {
	a.opts.Steps = append(a.opts.Steps, &AggregateSort{
		Keys: keys,
		Max:  max,
	})
	return a
}

// GroupBy adds a new group by statement (constructed with a GroupByBuilder)
func (a *AggregateBuilder) GroupBy(g AggregateGroupBy) *AggregateBuilder {
	a.opts.Steps = append(a.opts.Steps, &g)
	return a
}

/*******************************************************************************
*   GROUP BY builder												          *
*******************************************************************************/

// NewGroupByBuilder creates a builder for group by statements in aggregates.
func NewGroupByBuilder() *GroupByBuilder {
	return &GroupByBuilder{
		group: AggregateGroupBy{},
	}
}

// GroupBy returns the grouping defined by the builder
func (g *GroupByBuilder) GroupBy() AggregateGroupBy {
	return g.group
}

// Property appends a property to the properties list, not adding it if
// it already exists
func (g *GroupByBuilder) Property(name string) *GroupByBuilder {
	g.group.Properties = append(g.group.Properties, name)
	return g
}

// Properties sets all the property for a group by at one time.
func (g *GroupByBuilder) Properties(properties []string) *GroupByBuilder {
	g.group.Properties = properties
	return g
}

// Reduce adds a reducer function to the group by.
func (g *GroupByBuilder) Reduce(r AggregateReducer) *GroupByBuilder {
	g.group.Reducers = append(g.group.Reducers, r)
	return g
}

/*******************************************************************************
*   Reducer shortcuts
*******************************************************************************/

// ReduceCount returns a Reducer configured to count records
func ReduceCount(as string) AggregateReducer {
	return AggregateReducer{Name: "count", As: as}
}

// ReduceCountDistinct returns a Reducer configured to count distinct values of a property
func ReduceCountDistinct(property, as string) AggregateReducer {
	return AggregateReducer{Name: "count_distinct", Args: []interface{}{property}, As: as}
}

// ReduceCountDistinctIsh returns a Reducer configured to count distinct values of a property approximately
func ReduceCountDistinctIsh(property, as string) AggregateReducer {
	return AggregateReducer{Name: "count_distinctish", Args: []interface{}{property}, As: as}
}

// ReduceSum returns a Reducer configured to return the sum of the values of the given property.
func ReduceSum(property, as string) AggregateReducer {
	return AggregateReducer{Name: "sum", Args: []interface{}{property}, As: as}
}

// ReduceMin returns a Reducer configured to return the minimum value of the given property.
func ReduceMin(property, as string) AggregateReducer {
	return AggregateReducer{Name: "min", Args: []interface{}{property}, As: as}
}

// ReduceMax returns a Reducer configured to return the maximum value of the given property.
func ReduceMax(property, as string) AggregateReducer {
	return AggregateReducer{Name: "max", Args: []interface{}{property}, As: as}
}

// ReduceAvg returns a Reducer configured to return the mean value of the given property.
func ReduceAvg(property, as string) AggregateReducer {
	return AggregateReducer{Name: "avg", Args: []interface{}{property}, As: as}
}

// ReduceAvg returns a Reducer configured to return the mean value of the given property.
func ReduceStdDev(property, as string) AggregateReducer {
	return AggregateReducer{Name: "stddev", Args: []interface{}{property}, As: as}
}

// ReduceAvg returns a Reducer configured to return the mean value of the given property.
func ReduceQuantile(property string, quantile float64, as string) AggregateReducer {
	return AggregateReducer{Name: "stddev", Args: []interface{}{property, quantile}, As: as}
}

// ReduceToList returns a reducer configured to merge all distinct values of the property into an array
func ReduceToList(property, as string) AggregateReducer {
	return AggregateReducer{Name: "tolist", Args: []interface{}{property}, As: as}
}

// ReduceFirstValue returns a Reducer configured to get the first value of a given property with optional
// sorting.
func ReduceFirstValue(property, order, as string) AggregateReducer {
	reduceFn := AggregateReducer{Name: "first_value", Args: []interface{}{property}, As: as}
	if order != SortNone {
		reduceFn.Args = append(reduceFn.Args, order)
	}
	return reduceFn
}

// ReduceFirstValue returns a Reducer configured to get the first value of a given property with optional
// sorting using another property as the comparator
func ReduceFirstValueBy(property, comparator, order, as string) AggregateReducer {
	reduceFn := AggregateReducer{Name: "first_value", Args: []interface{}{property, "BY", comparator}, As: as}
	if order != SortNone {
		reduceFn.Args = append(reduceFn.Args, order)
	}
	return reduceFn
}

// ReduceRandomSample returns a Reducder configured to perform a random sampling of values of the
// property with a given sample size
func ReduceRandomSample(property string, sampleSize int64, as string) AggregateReducer {
	return AggregateReducer{Name: "random_sample", Args: []interface{}{property, sampleSize}, As: as}
}
