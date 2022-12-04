package aoc22

import (
	"bufio"
	"encoding"
	"io"
	"os"
	"testing"
)

func ReadTestFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}
	return data
}

func ParseLines[T encoding.TextUnmarshaler](r io.Reader, x T) ([]T, error) {
	var result []T

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var x T
		if err := x.UnmarshalText(scanner.Bytes()); err != nil {
			return nil, err
		}
		result = append(result, x)
	}

	return result, scanner.Err()
}
