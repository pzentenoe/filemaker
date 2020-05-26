package filemaker

type groupQuery struct {
	queries []*queryFieldOperator
}

func NewGroupQuery(queries ...*queryFieldOperator) *groupQuery {
	return &groupQuery{queries: queries}
}
