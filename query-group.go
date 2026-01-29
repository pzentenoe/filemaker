package filemaker

type groupQuery struct {
	queries []*queryFieldOperator
}

func NewGroupQuery(queries ...*queryFieldOperator) *groupQuery {
	return &groupQuery{queries: queries}
}

// AddQuery adds a query field operator to the group.
func (g *groupQuery) AddQuery(name string, operator FieldOperator, value string) *groupQuery {
	g.queries = append(g.queries, NewQueryFieldOperator(name, value, operator))
	return g
}
