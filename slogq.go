package mapq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

// Create a new query builder from a structure log file
func FromSlogFile(filePath string) (*Query, error) {
	errorf := packageErrorf("FromSlogFile")
	file, err := os.Open(filePath)

	if err != nil {
		return nil, errorf("os.Open: could not read file '%s': %w", filePath, err)
	}

	return FromSlogReader(file)
}

// Create a new query builder from structured log data in a string variable
func FromSlogString(logStr string) (*Query, error) {
	return FromSlogReader(strings.NewReader(logStr))
}

// Create a new query builder from structured log data in a []byte variable
func FromSlogBytes(logBytes []byte) (*Query, error) {
	return FromSlogReader(bytes.NewReader(logBytes))
}

// Create a new query builder from structured log data in a reader
func FromSlogReader(logReader io.Reader) (*Query, error) {
	scanner := bufio.NewScanner(logReader)

	maps := []map[string]any{}

	var errs error
	for scanner.Scan() {
		line := scanner.Bytes()

		slog := map[string]any{}
		err := json.Unmarshal(line, &slog)

		errs = errors.Join(errs, err)
		if err == nil {
			maps = append(maps, slog)
		}
	}

	return FromSlice(maps), errs
}
