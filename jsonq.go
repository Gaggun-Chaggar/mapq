package mapq

import (
	"encoding/json"
	"io"
	"os"
)

func FromJSONFile(filePath string) (*Query, error) {
	errorf := packageErrorf("FromFile")
	file, err := os.Open(filePath)

	if err != nil {
		return nil, errorf("os.Open: could not read file '%s': %w", filePath, err)
	}

	return FromJSONReader(file)
}

func FromJSONString(jsonStr string) (*Query, error) {
	return FromJSONBytes([]byte(jsonStr))
}

func FromJSONReader(jsonReader io.Reader) (*Query, error) {
	errorf := packageErrorf("FromReader")
	bytes, err := io.ReadAll(jsonReader)

	if err != nil {
		return nil, errorf("io.ReadAll: could not read json: %w", err)
	}

	return FromJSONBytes(bytes)
}

func FromJSONBytes(jsonBytes []byte) (*Query, error) {
	errorf := packageErrorf("FromBytes")
	maps := []map[string]any{}

	err := json.Unmarshal(jsonBytes, &maps)

	if err != nil {
		return nil, errorf("json.Unmarshal: could not read json bytes: %w", err)
	}

	return FromSlice(maps), nil
}
