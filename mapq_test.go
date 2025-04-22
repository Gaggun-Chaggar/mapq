package mapq_test

import (
	"mapq"
	. "mapq/testutils"
	"testing"

	. "github.com/smarty/assertions"
)

func TestMapQFromSlice(t *testing.T) {
	builder := mapq.FromSlice([]map[string]any{
		{"level": "error", "nested": map[string]any{"object": "hi"}, "all": true},
		{"level": "info", "butt": "big", "all": true},
		{"level": "warn", "array": []any{1, 2}, "all": true},
		{"level": "warn", "butt": "big", "all": true},
		{"level": "error", "all": true, "nested": []any{"object", "hi"}},
	})

	testQueryBuilder(t, builder)
}

func testQueryBuilder(t *testing.T, builder *mapq.Query) {
	t.Run("supports key-value pairs", func(t *testing.T) {
		query := builder.Where(
			mapq.Assert("level", ShouldEqual, "error"),
		)

		filteredSlogs := mapq.Filter(query)
		hasErr := mapq.Exists(query)
		onlyHasErr := mapq.All(query)
		hasTwoOnly := mapq.Has(2, query)
		ExpectThat(t, filteredSlogs, ShouldHaveLength, 2)
		ExpectThat(t, filteredSlogs[0]["level"], ShouldEqual, "error")
		ExpectThat(t, hasErr, ShouldBeTrue)
		ExpectThat(t, onlyHasErr, ShouldBeFalse)
		ExpectThat(t, hasTwoOnly, ShouldBeTrue)
	})

	t.Run("supports nested maps", func(t *testing.T) {
		query := builder.Where(
			mapq.Assert("nested.object", ShouldEqual, "hi"),
		)
		filteredSlogs := mapq.Filter(query)
		ExpectThat(t, filteredSlogs, ShouldHaveLength, 1)
		ExpectThat(t, filteredSlogs[0]["level"], ShouldEqual, "error")
		ExpectThat(t, filteredSlogs[0]["nested"], ShouldEqual, map[string]any{"object": "hi"})
	})

	t.Run("supports OR clauses", func(t *testing.T) {
		query := builder.Where(
			mapq.Or(
				mapq.Assert("level", ShouldEqual, "warn"),
				mapq.Assert("butt", ShouldContainSubstring, "big"),
			),
		)
		ExpectThat(t, mapq.Has(3, query), ShouldBeTrue)
	})

	t.Run("supports XOR clauses", func(t *testing.T) {
		query := builder.Where(
			mapq.XOr(
				mapq.Assert("level", ShouldEqual, "warn"),
				mapq.Assert("butt", ShouldContainSubstring, "big"),
			),
		)
		ExpectThat(t, mapq.Has(2, query), ShouldBeTrue)
	})

	t.Run("supports nested clauses", func(t *testing.T) {
		query := builder.Where(
			mapq.Or(
				mapq.And(
					mapq.Assert("level", ShouldEqual, "warn"),
					mapq.Assert("butt", ShouldContainSubstring, "big"),
				),
				mapq.Assert("level", ShouldEqual, "error"),
			),
		)
		ExpectThat(t, mapq.Has(3, query), ShouldBeTrue)
	})

	t.Run("supports arrays", func(t *testing.T) {
		query := builder.Where(mapq.Assert("array.0", ShouldEqual, 1))
		ExpectThat(t, mapq.Has(1, query), ShouldBeTrue)
	})

	t.Run("all collector returns true when all logs fit query", func(t *testing.T) {
		query := builder.Where(mapq.Assert("all", ShouldBeTrue))
		ExpectThat(t, mapq.All(query), ShouldBeTrue)
	})

	t.Run("exists results false when no logs fit query", func(t *testing.T) {
		query := builder.Where(mapq.Assert("all", ShouldBeFalse))
		ExpectThat(t, mapq.Exists(query), ShouldBeFalse)
	})

	t.Run("handles query for non existent nested fields gracefully", func(t *testing.T) {
		query := builder.Where(mapq.Assert("level.nested", ShouldEqual, "bad"))
		ExpectThat(t, mapq.Exists(query), ShouldBeFalse)
	})
}
