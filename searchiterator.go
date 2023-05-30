package grstack

import "context"

// SearchIterator is used to incrementally iterate over a collection of elements.
type SearchIterator struct {
	cmd   *QueryCmd
	pos   int64
	limit *Limit
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

	// Advance cursor, check if we are still within range.
	if it.pos < it.cmd.options.Limit.Offset {
		it.pos++
		return true
	}

	for {

		// Return if there is no more data to fetch.
		if len(it.cmd.Val()) == 0 {
			return false
		}

		// Add the limit to the offset and run the command again.
		it.limit.Offset += it.limit.Num
		it.cmd.options.Limit = it.limit
		it.cmd.args = it.cmd.options.serialize

		// Fetch next page.
		switch it.cmd.args[0] {
		case "scan", "qscan":
			it.cmd.args[1] = it.cmd.cursor
		default:
			it.cmd.args[2] = it.cmd.cursor
		}

		err := it.cmd.process(ctx, it.cmd)
		if err != nil {
			return false
		}

		it.pos = 1

		// Redis can occasionally return empty page.
		if len(it.cmd.page) > 0 {
			return true
		}
	}
}

// Val returns the key/field at the current cursor position.
func (it *SearchIterator) Val() string {
	var v string
	if it.cmd.Err() == nil && it.pos > 0 && it.pos <= len(it.cmd.page) {
		v = it.cmd.page[it.pos-1]
	}
	return v
}
