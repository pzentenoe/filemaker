package filemaker

type queryGroup struct {
	queries []*queryFieldOperator
}

func NewQueryGroup(queries ...*queryFieldOperator) *queryGroup {
	return &queryGroup{queries: queries}
}
