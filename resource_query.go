package ctypes

import (
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	RQAll = iota
	RQAny
	RQNor
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
	Mode    int      `json:"mode"`
	Limit   uint64   `json:"limit"`
	Offset  uint64   `json:"offset"`
	Queries []RQQ    `json:"queries"`
	Sort    []RQSort `json:"sort"`
}

func NewResourceQuery(mode int) *ResourceQuery {
	return &ResourceQuery{
		Limit: 10,
		Mode:  mode,
	}
}

func ResourceQueryFromURL(url *url.URL, fieldPrefix string) (*ResourceQuery, error) {
	queryParams := url.Query()

	resourceQuery := NewResourceQuery(RQAll)

	for key, queryStr := range queryParams {
		if len(queryStr) == 0 {
			continue
		}

		switch strings.TrimSpace(strings.ToLower(key)) {

		// Parse the limit from the query parameter. Even if the limit is not valid, the rq builder will take care of it
		case "limit":
			lim, _ := strconv.Atoi(queryStr[0])
			resourceQuery.ResourceLimit(uint64(lim))

		// Parse the offset from the query parameter. Even if the offset is not valid, the builder will take care of it
		case "offset":
			off, _ := strconv.Atoi(queryStr[0])
			resourceQuery.ResourceOffset(uint64(off))

		// Parse the ascending sort parameters
		case "sort", "sortasc":
			for _, field := range queryStr {
				if !validateFieldName(field) {
					return nil, errors.New("invalid field name " + field)
				}

				resourceQuery.SortAsc(fieldPrefix + field)
			}

		// Parse the descending sort parameters
		case "sortdesc":
			for _, field := range queryStr {
				if !validateFieldName(field) {
					return nil, errors.New("invalid field name " + field)
				}

				resourceQuery.SortDesc(fieldPrefix + field)
			}

		case "mode":
			switch queryStr[0] {
			case "all":
				resourceQuery.ModeAll()
			case "any":
				resourceQuery.ModeAny()
			case "nor":
				resourceQuery.ModeNor()
			}

		// All other query parameters might be other operators
		default:
			if parts := strings.Split(key, "."); len(parts) > 1 {
				prefix := parts[0]

				field := fieldPrefix + strings.TrimPrefix(key, prefix + ".")

				if !validateFieldName(field) {
					return nil, errors.New("invalid field name " + field)
				}

				for _, val := range queryStr {
					switch strings.TrimSpace(strings.ToLower(prefix)) {
					case "eq":
						resourceQuery.Equals(field, val)
					case "!eq":
						resourceQuery.DoesNotEqual(field, val)
					case "ex":
						resourceQuery.Exists(field)
					case "!ex":
						resourceQuery.DoesNotExist(field)
					case "cont":
						resourceQuery.Contains(field, val)
					case "!cont":
						resourceQuery.DoesNotContain(field, val)
					case "sw":
						resourceQuery.StartsWith(field, val)
					case "!sw":
						resourceQuery.DoesNotStartWith(field, val)
					case "ew":
						resourceQuery.EndsWith(field, val)
					case "!ew":
						resourceQuery.DoesNotEndWith(field, val)
					case "gt":
						resourceQuery.GreaterThan(field, val)
					case "!gt":
						resourceQuery.LessThanOrEqualTo(field, val)
					case "gte":
						resourceQuery.GreaterThanOrEqualTo(field, val)
					case "!gte":
						resourceQuery.LessThan(field, val)
					case "lt":
						resourceQuery.LessThan(field, val)
					case "!lt":
						resourceQuery.GreaterThanOrEqualTo(field, val)
					case "lte":
						resourceQuery.LessThanOrEqualTo(field, val)
					case "!lte":
						resourceQuery.GreaterThan(field, val)
					case "rx":
						resourceQuery.MatchesRegEx(field, val)
					case "!rx":
						resourceQuery.DoesNotMatchRegEx(field, val)
					}
				}
			}
		}
	}

	return resourceQuery.sortQueriesByField(), nil
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

func (q *ResourceQuery) ModeAll() *ResourceQuery {
	q.Mode = RQAll
	return q
}

func (q *ResourceQuery) ModeAny() *ResourceQuery {
	q.Mode = RQAny
	return q
}

func (q *ResourceQuery) ModeNor() *ResourceQuery {
	q.Mode = RQNor
	return q
}

func (q *ResourceQuery) SortAsc(field string) *ResourceQuery {
	q.Sort = append(q.Sort, RQSort{
		Field:     field,
		Ascending: true,
	})

	return q
}

func (q *ResourceQuery) SortDesc(field string) *ResourceQuery {
	q.Sort = append(q.Sort, RQSort{
		Field:     field,
		Ascending: false,
	})

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

func (q *ResourceQuery) sortQueriesByField() *ResourceQuery {
	sort.Slice(q.Queries, func(i, j int) bool {
		a, b := q.Queries[i], q.Queries[j]
		return a.Field > b.Field
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

// RQSort is a single sort operation
type RQSort struct {
	Field     string `json:"field"`
	Ascending bool   `json:"asc"`
}

func (r *RQSort) FieldNameValid() bool {
	return validateFieldName(r.Field)
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

	// Then, try parsing as a float.
	// and make sure the string has a decimal in it. An overflow size integer will be parsed as a float in some cases
	if f, err := strconv.ParseFloat(*r.Value, 64); err == nil && strings.Contains(*r.Value, ".") {
		return f
	}

	// Try parsing the value as a uuid
	if id, err := uuid.Parse(*r.Value); err != nil {
		return id
	}

	// TODO add more parse checks

	return *r.Value
}

func (r *RQQ) FieldNameValid() bool {
	return validateFieldName(r.Field)
}

func validateFieldName(name string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z_$@\\-]{1}[a-zA-Z_$@\\-0-9.]+$", name)
	return match
}