package testutils

import (
	"errors"
	"os"
	"testing"

	. "github.com/smarty/assertions"
)

func ExpectThat(t *testing.T, actual any, assertion SoFunc, expected ...any) {
	t.Helper()
	msg := assertion(actual, expected...)

	if msg != "" {
		t.Log(msg)
		t.FailNow()
	}

	t.Logf("✓")
}

func TestFile(t *testing.T, testFileName string, content string) {
	f, _ := os.Create(testFileName)
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() {
		os.Remove(testFileName)
	})
}

type ErrReader struct{}

func (ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
