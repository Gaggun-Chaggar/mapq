package mapq_test

import (
	"strings"
	"testing"

	"github.com/Gaggun-Chaggar/mapq"

	. "github.com/Gaggun-Chaggar/mapq/testutils"

	. "github.com/smarty/assertions"
)

func TestJSONQ(t *testing.T) {
	json := strings.TrimSpace(`[
		{ "level": "error", "nested": { "object": "hi" }, "all": true },
		{ "level": "info", "butt": "big", "all": true },
		{ "level": "warn", "array": [1,2], "all": true },
		{ "level": "warn", "butt": "big", "all": true },
		{ "level": "error", "all": true, "nested": ["object", "hi"]}
	]`)

	t.Run("from file", func(t *testing.T) {
		testFileName := "./test.array.json"
		TestFile(t, testFileName, json)

		builder, err := mapq.FromJSONFile(testFileName)
		ExpectThat(t, err, ShouldBeNil)
		testQueryBuilder(t, builder)
	})

	t.Run("from file - should return error on file open failure", func(t *testing.T) {
		_, err := mapq.FromJSONFile("./i-do-not-exist")
		ExpectThat(t, err, ShouldNotBeNil)
	})

	t.Run("from string", func(t *testing.T) {
		builder, err := mapq.FromJSONString(json)
		ExpectThat(t, err, ShouldBeNil)
		testQueryBuilder(t, builder)
	})

	t.Run("from reader", func(t *testing.T) {
		builder, err := mapq.FromJSONReader(strings.NewReader(json))
		ExpectThat(t, err, ShouldBeNil)
		testQueryBuilder(t, builder)
	})

	t.Run("from reader - should return err on failed read", func(t *testing.T) {
		_, err := mapq.FromJSONReader(ErrReader{})
		ExpectThat(t, err, ShouldNotBeNil)
	})

	t.Run("should return err given invalid json", func(t *testing.T) {
		_, err := mapq.FromJSONString(`[`)
		ExpectThat(t, err, ShouldNotBeNil)
	})

}
