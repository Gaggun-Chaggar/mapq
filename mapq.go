package mapq

import (
	"iter"
	"slices"
)

type Query struct {
	it    iter.Seq2[int, map[string]any]
	stack []any
}

func FromSlice(maps []map[string]any) *Query {
	return &Query{slices.All(maps), []any{}}
}

func (q *Query) Where(asserts ...any) *Query {
	newQ := &Query{
		it:    q.it,
		stack: q.stack,
	}
	newQ.stack = append(q.stack, asserts...)
	return newQ
}

func Assert(property string, op Operation, value ...any) assertion {
	return assertion{property, op, value}
}

func And(asserts ...any) joiner {
	return joiner{andJoin, asserts}
}

func Or(asserts ...any) joiner {
	return joiner{orJoin, asserts}
}

func XOr(asserts ...any) joiner {
	return joiner{xOrJoin, asserts}
}

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

func All(q *Query) bool {
	for _, s := range q.it {
		joiner := &joiner{andJoin, q.stack}
		if !joiner.collect(s) {
			return false
		}
	}
	return true
}

func Exists(q *Query) bool {
	for _, s := range q.it {
		joiner := &joiner{andJoin, q.stack}
		if joiner.collect(s) {
			return true
		}
	}
	return false
}

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
