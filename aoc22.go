package aoc22

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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

var intsRegexp = regexp.MustCompile(`\d+`)

// ReadInts returns a slice of all integers in the given string. For example,
// "a 12\n3, b::4" becomes []int{12, 3, 4}.
func ReadInts(s string) []int {
	var result []int
	for _, num := range intsRegexp.FindAllString(s, -1) {
		n, err := strconv.Atoi(num)
		if err != nil {
			panic(fmt.Sprintf("bad number: %q", num))
		}
		result = append(result, n)
	}
	return result
}
