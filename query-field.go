package filemaker

type FieldOperator string

const (
	Equal            FieldOperator = "eq"
	Contains         FieldOperator = "cn"
	BeginsWith       FieldOperator = "bw"
	EndsWith         FieldOperator = "ew"
	GreaterThan      FieldOperator = "gt"
	GreaterThanEqual FieldOperator = "gte"
	LessThan         FieldOperator = "lt"
	LessThanEqual    FieldOperator = "lte"
)

type queryFieldOperator struct {
	Name     string
	Value    string
	Operator FieldOperator
}

func NewQueryFieldOperator(name, value string, operator FieldOperator) *queryFieldOperator {
	return &queryFieldOperator{
		Name:     name,
		Value:    value,
		Operator: operator,
	}
}

func (qf *queryFieldOperator) valueWithOp() string {
	switch qf.Operator {
	case Equal:
		return "==" + qf.Value
	case Contains:
		return "==*" + qf.Value + "*"
	case BeginsWith:
		return "==" + qf.Value + "*"
	case EndsWith:
		return "==*" + qf.Value
	case GreaterThan:
		return ">" + qf.Value
	case GreaterThanEqual:
		return ">=" + qf.Value
	case LessThan:
		return "<" + qf.Value
	case LessThanEqual:
		return "<=" + qf.Value
	default:
		return qf.Value
	}
}
