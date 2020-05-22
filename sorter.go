package filemaker

type SortOrder string

const (
	Ascending  SortOrder = "ascend"
	Descending SortOrder = "descend"
)

type Sorter struct {
	FieldName string    `json:"fieldName"`
	SortOrder SortOrder `json:"sortOrder"`
}

func NewSorter(fieldName string, sortOrder SortOrder) *Sorter {
	return &Sorter{
		FieldName: fieldName,
		SortOrder: sortOrder,
	}
}

func (so *SortOrder) String() string {
	return string(*so)
}
