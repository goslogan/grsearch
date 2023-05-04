package grstack

import (
	"fmt"
	"time"

	"github.com/goslogan/grstack/internal"
)

// aggregate.go contains functionality used to implement FT.AGGREGATE

type AggregateOptions struct {
	Verbatim bool
	Load     []AggregateLoad
	Timeout  time.Duration
	GroupBy  []AggregateGroupBy
	SortBy   *AggregateSort
	Apply    []AggregateApply
	Limit    *Limit
	Filter   string
	Cursor   *AggregateCursor // nil means no cursor
	Params   map[string]interface{}
	Dialect  uint8
}

type AggregateGroupBy struct {
	Properties []string
	Reducers   []AggregateReducer
}

type AggregateReducer struct {
	Name string
	As   string
	Args []interface{}
}
type AggregateApply struct {
	Expression string
	As         string
}

// Load represents parameters to the LOAD argument
type AggregateLoad struct {
	Name string
	As   string
}

type AggregateCursor struct {
	Count   uint64
	MaxIdle time.Duration
}

type AggregateSort struct {
	Keys []AggregateSortKey
	Max  int64
}
type AggregateSortKey struct {
	Name  string
	Order string
}

// LoadAll can be used to indicate FT.AGGREGATE idx LOAD *
var LoadAll = AggregateLoad{Name: "*"}

/******************************************************************************
* Methods operating on the aggregate struct itself and simple operations on  *
* the struct																  *
******************************************************************************/

// NewAggregateOptions creates a new query with defaults set
func NewAggregateOptions() *AggregateOptions {
	return &AggregateOptions{
		Dialect: defaultDialect,
	}
}

func (a *AggregateOptions) serialize() []interface{} {
	args := []interface{}{}
	if a.Timeout != defaultTimeout {
		args = internal.AppendStringArg(args, "timeout", fmt.Sprintf("%d", a.Timeout.Milliseconds()))
	}

	args = append(args, a.serializeLoad()...)
	for _, g := range a.GroupBy {
		args = append(args, g.serialize())
	}
	args = append(args, a.SortBy.serialize()...)
	for _, a := range a.Apply {
		args = append(args, a.serialize()...)
	}

	if a.Limit != nil {
		args = append(args, a.Limit.serialize()...)
	}

	if a.Filter != "" {
		args = append(args, "filter", a.Filter)
	}

	if a.Cursor != nil {
		args = append(args, a.Cursor.serialize()...)
	}

	if len(a.Params) != 0 {
		args = append(args, "params", len(a.Params))
		for n, v := range a.Params {
			args = append(args, n, v)
		}
	}

	if a.Dialect != defaultDialect {
		args = append(args, "dialect", a.Dialect)
	}

	return args

}

func (a *AggregateOptions) serializeLoad() []interface{} {

	if len(a.Load) == 1 && a.Load[0].Name == "*" {
		return []interface{}{"load", "*"}
	}
	loads := []interface{}{"load", len(a.Load)}
	for _, l := range a.Load {
		loads = append(loads, l.serialize())
	}
	return loads
}

func (a *AggregateApply) serialize() []interface{} {

	args := []interface{}{"apply", a.Expression}
	if a.As != "" {
		args = append(args, "as", a.As)
	}
	return args
}

func (l AggregateLoad) serialize() []interface{} {
	if l.As == "" {
		return []interface{}{l.Name}
	} else {
		return []interface{}{l.Name, "as", l.As}
	}
}

func (s AggregateSort) serialize() []interface{} {
	nArgs := len(s.Keys)
	if nArgs > 0 && s.Max != 0 {
		nArgs += 2
	}
	keys := []interface{}{nArgs}
	for _, k := range s.Keys {
		keys = append(keys, k.Name, k.Order)
	}
	if nArgs != 0 {
		keys = append(keys, "max", s.Max)
	}
	return keys
}

func (c *AggregateCursor) serialize() []interface{} {
	args := []interface{}{"withcursor"}
	if c.Count != 0 {
		args = append(args, "count", c.Count)
	}
	if c.MaxIdle != 0 {
		args = append(args, "maxidle", c.MaxIdle.Milliseconds())
	}

	return args
}

/******************************************************************************
* Methods for Groupby
******************************************************************************/

func (r *AggregateReducer) serialize() []interface{} {
	args := []interface{}{r.Name, len(r.Args)}
	args = append(args, r.Args...)
	if r.As != "" {
		args = append(args, "as", r.As)
	}

	return args
}

func (g *AggregateGroupBy) serialize() []interface{} {
	args := []interface{}{"GROUPBY", len(g.Properties)}
	args = append(args, g.Properties)
	for _, r := range g.Reducers {
		args = append(args, r.serialize()...)
	}
	return args
}
