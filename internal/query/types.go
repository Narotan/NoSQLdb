package query

// Query представляет распаршенный JSON-запрос
type Query struct {
	Conditions map[string]any
}

type Operator string

const (
	OpEq   Operator = "$eq"
	OpGt   Operator = "$gt"
	OpLt   Operator = "$lt"
	OpLike Operator = "$like"
	OpIn   Operator = "$in"
	OpAnd  Operator = "$and"
	OpOr   Operator = "$or"
)
