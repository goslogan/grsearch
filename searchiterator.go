package grstack

// iterators - how to implement. search results do not support cursors like aggregates
// so we can implement an iterator insipired by the Scan iterator. In this version we need
// to work with search offsets and limits instead of iterator values though.
// How...
// 	1 do an initial search and call Iterator.
//  2 QueryCmd will then initialise an iterator with itself
//  3 Each time we call next we need to check if pos >= the number of results (from Limit in the options)
//  4 If equal we need to take the search options to create a new QueryCmd where we update the offset
//    and run the search again.

import "context"

// SearchIterator is used to incrementally iterate over a collection of elements.
type SearchIterator struct {
	options *QueryOptions
	index   string
	query   string
	pos     int64
	curPos  int64
	maxPos  int64
	process cmdable
	cmd     *QueryCmd
}

// NewSearchIterator returns a configured iterator for QueryCmd
func NewSearchIterator(ctx context.Context, cmd *QueryCmd, process cmdable) *SearchIterator {
	return &SearchIterator{
		cmd:     cmd,
		options: cmd.options,
		index:   cmd.Args()[1].(string),
		query:   cmd.Args()[2].(string),
		process: process,
		pos:     0,
		curPos:  1,
		maxPos:  cmd.Count(),
	}
}

// Err returns the last iterator error, if any.
func (it *SearchIterator) Err() error {
	return it.cmd.Err()
}

// Next advances the cursor and returns true if more values can be read.
func (it *SearchIterator) Next(ctx context.Context) bool {
	// Instantly return on errors.
	if it.cmd.Err() != nil {
		return false
	}

	for {

		if it.cmd.Len() == 0 {
			return false
		}

		if it.pos < it.cmd.Len() {
			it.pos++
			return true
		}

		if it.options.Limit == nil {
			it.options.Limit = &Limit{
				Num:    DefaultLimit,
				Offset: DefaultOffset + DefaultLimit,
			}
		} else {
			it.options.Limit.Offset += it.options.Limit.Num
		}

		it.cmd = it.process.FTSearch(ctx, it.index, it.query, it.options)
		if it.Err() != nil {
			return false
		}
		it.pos = 0

	}
}

// Val returns the key/field at the current cursor position.
func (it *SearchIterator) Val() *QueryResult {
	var v *QueryResult
	if it.cmd.Err() == nil && it.pos > 0 && it.pos <= it.options.Limit.Num {
		v = it.cmd.Val().Results[it.pos-1]
	}
	return v
}
