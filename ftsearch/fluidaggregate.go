package ftsearch

/*******************************************************************************
*   GROUP BY related methods												   *
*******************************************************************************/

// AddProperty appends a property to the properties list, not adding it if
// it already exists
func (g *AggregateGroupBy) AddProperty(name string) *AggregateGroupBy {
	g.Properties = append(g.Properties, name)
	return g
}

// this module implements a fluid interface to the aggregate structure's configuration
func (g *AggregateGroupBy) RemoveField(name string) *AggregateGroupBy {
	g.Properties = append(g.Properties, name)
	return g
}

// AddLoad adds a field to the load list for the aggregate. The alias can be the
// empty string.
func (a *AggregateOptions) AddLoad(name string, as string) *AggregateOptions {
	l := Load{Name: name, As: as}
	a.Load = append(a.Load, l)
	return a
}

// RemoveLoad removes a field from the load list for an aggregate.
func (a *AggregateOptions) RemoveLoad(identifier string) *AggregateOptions {
	nl := make([]Load, 0)
	for _, l := range a.Load {
		if l.Name != identifier {
			nl = append(nl, l)
		}
	}
	a.Load = nl
	return a
}

// WithLoad sets all load parameters for an aggregate. If set to
// ftsearch.LoadAll all fields are loaded.
func (a *AggregateOptions) WithLoad(l []Load) *AggregateOptions {
	a.Load = l
	return a
}

// ClearLoad removes any defined Load fields and clears them.
func (a *AggregateOptions) ClearLoad() *AggregateOptions {
	a.Load = []Load{}
	return a
}
