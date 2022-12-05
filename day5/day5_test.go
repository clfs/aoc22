package day5

import (
	"bytes"
	"testing"

	"github.com/clfs/aoc22"
)

func TestPart1(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"testdata/small.txt", "CMZ"},
		{"testdata/input.txt", "SHQWSRBDL"},
	}

	for _, c := range cases {
		data := aoc22.ReadTestFile(t, c.path)
		got := Part1(bytes.NewReader(data))
		if got != c.want {
			t.Errorf("Part1(%q) == %q, want %q", c.path, got, c.want)
		}
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"testdata/small.txt", "MCD"},
		{"testdata/input.txt", "SHQWSRBDL"},
	}

	for _, c := range cases {
		data := aoc22.ReadTestFile(t, c.path)
		got := Part2(bytes.NewReader(data))
		if got != c.want {
			t.Errorf("Part1(%q) == %q, want %q", c.path, got, c.want)
		}
	}
}
