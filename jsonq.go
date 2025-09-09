package mapq

import (
	"encoding/json"
	"io"
	"os"
)

// Create a new query builder from a json file
func FromJSONFile(filePath string) (*Query, error) {
	errorf := packageErrorf("FromJSONFile")
	file, err := os.Open(filePath)

	if err != nil {
		return nil, errorf("os.Open: could not read file '%s': %w", filePath, err)
	}

	return FromJSONReader(file)
}

// Create a new query builder from a json string
func FromJSONString(jsonStr string) (*Query, error) {
	return FromJSONBytes([]byte(jsonStr))
}

// Create a new query builder from a json reader
func FromJSONReader(jsonReader io.Reader) (*Query, error) {
	errorf := packageErrorf("FromJSONReader")
	bytes, err := io.ReadAll(jsonReader)

	if err != nil {
		return nil, errorf("io.ReadAll: could not read json: %w", err)
	}

	return FromJSONBytes(bytes)
}

// Create a new query builder from json bytes
func FromJSONBytes(jsonBytes []byte) (*Query, error) {
	errorf := packageErrorf("FromJSONBytes")
	maps := []map[string]any{}

	err := json.Unmarshal(jsonBytes, &maps)

	if err != nil {
		return nil, errorf("json.Unmarshal: could not read json bytes: %w", err)
	}

	return FromSlice(maps), nil
}
