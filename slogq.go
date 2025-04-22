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

func FromSlogFile(filePath string) (*Query, error) {
	errorf := packageErrorf("FromFile")
	file, err := os.Open(filePath)

	if err != nil {
		return nil, errorf("os.Open: could not read file '%s': %w", filePath, err)
	}

	return FromSlogReader(file)
}

func FromSlogString(logStr string) (*Query, error) {
	return FromSlogReader(strings.NewReader(logStr))
}

func FromSlogBytes(logBytes []byte) (*Query, error) {
	return FromSlogReader(bytes.NewReader(logBytes))
}

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
