package mapq_test

import (
	"mapq"
	"strings"
	"testing"

	. "mapq/testutils"

	. "github.com/smarty/assertions"
)

func TestSlogQ(t *testing.T) {
	slogs := strings.TrimSpace(`
	{ "level": "error", "nested": { "object": "hi" }, "all": true }
	{ "level": "info", "butt": "big", "all": true }
	{ "level": "warn", "array": [1,2], "all": true }
	{ "level": "warn", "butt": "big", "all": true }
	{ "level": "error", "all": true, "nested": ["object", "hi"] }
	`)

	t.Run("from file", func(t *testing.T) {
		testFileName := "./test.slog.log"
		TestFile(t, testFileName, slogs)

		builder, err := mapq.FromSlogFile(testFileName)
		ExpectThat(t, err, ShouldBeNil)
		testQueryBuilder(t, builder)
	})

	t.Run("from file - should return error on file open failure", func(t *testing.T) {
		_, err := mapq.FromSlogFile("./i-do-not-exist")
		ExpectThat(t, err, ShouldNotBeNil)
	})

	t.Run("from string", func(t *testing.T) {
		builder, _ := mapq.FromSlogString(slogs)
		testQueryBuilder(t, builder)
	})

	t.Run("from bytes", func(t *testing.T) {
		builder, _ := mapq.FromSlogBytes([]byte(slogs))
		testQueryBuilder(t, builder)
	})

	t.Run("should err on unparsable slog", func(t *testing.T) {
		_, err := mapq.FromSlogString(`{`)
		ExpectThat(t, err, ShouldNotBeNil)
	})

}
