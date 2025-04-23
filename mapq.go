package mapq

import (
	"iter"
	"slices"
)

type Query struct {
	it    iter.Seq2[int, map[string]any]
	stack []any
}

// Create a query builder from a slice of objects
func FromSlice(maps []map[string]any) *Query {
	return &Query{slices.All(maps), []any{}}
}

// Build a query
// returns a new instance of `*Query` with assertions included
func (q *Query) Where(asserts ...any) *Query {
	newQ := &Query{
		it:    q.it,
		stack: q.stack,
	}
	newQ.stack = append(q.stack, asserts...)
	return newQ
}

// Create a new property assertion
// Supports nested properties with dot notation
// e.g. mapq.Assert("httpRequest.status", ShouldEqual, 200)
func Assert(property string, op Operation, value ...any) assertion {
	return assertion{property, op, value}
}

// On execution, joins assertion results with a logical 'and' operation
func And(asserts ...any) joiner {
	return joiner{andJoin, asserts}
}

// On execution, joins assertion results with a logical 'or' operation
func Or(asserts ...any) joiner {
	return joiner{orJoin, asserts}
}

// On execution, joins assertion results with a logical 'exclusive or' operation
func XOr(asserts ...any) joiner {
	return joiner{xOrJoin, asserts}
}

// Executes a query and returns filtered results of the query
func Filter(q *Query) []map[string]any {
	var result []map[string]any

	for _, s := range q.it {
		joiner := &joiner{andJoin, q.stack}
		if joiner.collect(s) {
			result = append(result, s)
		}
	}

	return result
}

// Returns true if all rows of data matches the query
func All(q *Query) bool {
	for _, s := range q.it {
		joiner := &joiner{andJoin, q.stack}
		if !joiner.collect(s) {
			return false
		}
	}
	return true
}

// Returns true if at least one row of data matches the query
func Exists(q *Query) bool {
	for _, s := range q.it {
		joiner := &joiner{andJoin, q.stack}
		if joiner.collect(s) {
			return true
		}
	}
	return false
}

// Returns true if the number of rows that matches the query is exactly `i`
func Has(i int, q *Query) bool {
	trueCount := 0

	for _, s := range q.it {
		joiner := &joiner{andJoin, q.stack}
		if joiner.collect(s) {
			trueCount++
		}
	}

	return trueCount == i
}
