package grsearch

import (
	"fmt"
	"strings"
	"time"

	"github.com/goslogan/grsearch/internal"
)

// AggregateOptions represents the options that can be passed to [FT.AGGREGATE].
// This can be built by calling the [NewAggregateOptions] function or via [AggregateOptionsBuilder.Options]
// using the Builder API.
type AggregateOptions struct {
	Verbatim bool             // Set to true if stemming should not be used
	Load     []AggregateLoad  // Values for the LOAD subcommand; use the [LoadAll] variable to represent "LOAD *"
	Timeout  time.Duration    // Sets the query timeout. If zero, no TIMEOUT subcommmand is used
	Cursor   *AggregateCursor // nil means no cursor
	Params   map[string]interface{}
	Dialect  uint8
	Steps    []AggregateStep // The steps to be executed in order
}

// AggregateGroupBy represents a single GROUPBY statement in a
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

type AggregateFilter string

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

type AggregateStep interface {
	serializeStep() []interface{}
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
		Params:  map[string]interface{}{},
		Dialect: defaultDialect,
	}
}

func (a *AggregateOptions) serialize() []interface{} {
	args := []interface{}{}

	if a.Verbatim {
		args = append(args, "verbatim")
	}
	if a.Timeout != 0 {
		args = internal.AppendStringArg(args, "timeout", fmt.Sprintf("%d", a.Timeout.Milliseconds()))
	}
	args = append(args, a.serializeLoad()...)

	for _, step := range a.Steps {
		args = append(args, step.serializeStep()...)
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

func (f AggregateFilter) serializeStep() []interface{} {
	return []interface{}{f}
}

func (a *AggregateOptions) serializeLoad() []interface{} {

	if len(a.Load) == 0 {
		return []interface{}{}
	}

	if len(a.Load) == 1 && a.Load[0].Name == "*" {
		return []interface{}{"load", "*"}
	}
	loads := []interface{}{"load", 0}
	for _, l := range a.Load {
		loads = append(loads, l.serialize()...)
	}
	loads[1] = len(loads) - 2
	return loads
}

func (a *AggregateApply) serializeStep() []interface{} {

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

func (s AggregateSort) serializeStep() []interface{} {

	keys := []interface{}{"SORTBY", 0}
	for _, k := range s.Keys {
		namePrefix := ""
		if !strings.HasPrefix(k.Name, "@") {
			namePrefix = "@"
		}
		if k.Order == "" {
			keys = append(keys, namePrefix+k.Name, "ASC")
		} else {
			keys = append(keys, namePrefix+k.Name, k.Order)
		}

	}

	keys[1] = len(keys) - 2

	if len(keys) != 0 && s.Max != 0 {
		keys = append(keys, "MAX", s.Max)
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
	args := []interface{}{"reduce", r.Name, len(r.Args)}
	args = append(args, r.Args...)
	if r.As != "" {
		args = append(args, "as", r.As)
	}

	return args
}

func (g *AggregateGroupBy) serializeStep() []interface{} {
	args := []interface{}{"GROUPBY", len(g.Properties)}
	for _, arg := range g.Properties {
		args = append(args, arg)
	}
	for _, r := range g.Reducers {
		args = append(args, r.serialize()...)
	}
	return args
}

// Agre
