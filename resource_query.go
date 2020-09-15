package ctypes

import (
	"strconv"
)

const (
	RQAll = iota
	RQAny
	RQNone
)

const (
	RQEquals = iota
	RQExists
	RQContains
	RQStartsWith
	RQEndsWith
	RQGreaterThan
	RQGreaterThanOrEqual
	RQLessThan
	RQLessThanOrEqual
	RQRegex
)

// ResourceQuery is the format used to perform custom resource queries against memory and execution logs
type ResourceQuery struct {
	Mode    int    `json:"mode"`
	Limit   uint64 `json:"limit"`
	Offset  uint64 `json:"offset"`
	Queries []RQQ  `json:"queries"`
}

func NewResourceQuery(mode int) *ResourceQuery {
	return &ResourceQuery{
		Limit: 10,
		Mode:  mode,
	}
}

func (q *ResourceQuery) ResourceLimit(limit uint64) *ResourceQuery {
	if limit <= 0 {
		limit = 10
	} else if limit > 1000 {
		limit = 1000
	}

	q.Limit = limit

	return q
}

func (q *ResourceQuery) ResourceOffset(offset uint64) *ResourceQuery {
	if offset < 0 {
		offset = 0
	}

	q.Offset = offset

	return q
}

func (q *ResourceQuery) Equals(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQEquals,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) DoesNotEqual(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQEquals,
		Value:    &value,
		Negate:   true,
	})

	return q
}

func (q *ResourceQuery) Exists(field string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQExists,
	})

	return q
}

func (q *ResourceQuery) DoesNotExist(field string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQExists,
		Negate:   true,
	})

	return q
}

func (q *ResourceQuery) Contains(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQContains,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) DoesNotContain(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQContains,
		Value:    &value,
		Negate:   true,
	})

	return q
}

func (q *ResourceQuery) StartsWith(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQStartsWith,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) DoesNotStartWith(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQStartsWith,
		Value:    &value,
		Negate:   true,
	})

	return q
}

func (q *ResourceQuery) EndsWith(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQEndsWith,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) DoesNotEndWith(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQEndsWith,
		Value:    &value,
		Negate:   true,
	})

	return q
}

func (q *ResourceQuery) GreaterThan(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQGreaterThan,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) GreaterThanOrEqualTo(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQGreaterThanOrEqual,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) LessThan(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQLessThan,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) LessThanOrEqualTo(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQLessThanOrEqual,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) MatchesRegEx(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQRegex,
		Value:    &value,
	})

	return q
}

func (q *ResourceQuery) DoesNotMatchRegEx(field, value string) *ResourceQuery {
	q.Queries = append(q.Queries, RQQ{
		Field:    field,
		Operator: RQRegex,
		Value:    &value,
		Negate:   true,
	})

	return q
}

// RQQ is one single query operation
type RQQ struct {
	Field    string  `json:"field"`
	Operator int     `json:"operator"`
	Value    *string `json:"value"`
	Negate   bool    `json:"negate"`
}

// ValueAsTyped will try to convert the value into a number of formats, and return the best matching one
func (r *RQQ) ValueAsTyped() interface{} {
	if r.Value == nil {
		return nil
	}

	// First, try parsing as an int
	if i, err := strconv.ParseInt(*r.Value, 10, 64); err == nil {
		return i
	}

	// Then, try parsing as a float
	if f, err := strconv.ParseFloat(*r.Value, 64); err == nil {
		return f
	}

	// TODO add more parse checks

	return *r.Value
}
