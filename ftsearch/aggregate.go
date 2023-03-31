package ftsearch

import "time"

// aggregate.go contains functionality used to implement FT.AGGREGATE

type AggregateOptions struct {
	Verbatim      bool
	Load          []string
	Timeout       time.Duration
	LoadAll       bool
	GroupBy       []AggregateGroupBy
	Apply         []AggregateApply
	Limit         *limit
	Filter        string
	Cursor        bool
	CursorLimit   int
	CursorMaxIdle int
	Params        map[string]interface{}
}

type AggregateApply struct {
	Name string
	As   string
}

type AggregateGroupBy struct {
	Properties []string
	Reducers   []Reducer
}

type Reducer struct {
	Name string
	As   string
	Args []string
}
