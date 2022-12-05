package day4

import (
	"bytes"
	"testing"

	"github.com/clfs/aoc22"
)

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/input.txt", 582},
		{"testdata/small.txt", 2},
	}

	for _, tc := range cases {
		got := Part1(bytes.NewReader(aoc22.ReadTestFile(t, tc.name)))

		if got != tc.want {
			t.Errorf("Part1(%q) = %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/input.txt", 893},
		{"testdata/small.txt", 4},
	}

	for _, tc := range cases {
		got := Part2(bytes.NewReader(aoc22.ReadTestFile(t, tc.name)))

		if got != tc.want {
			t.Errorf("Part1(%q) = %d, want %d", tc.name, got, tc.want)
		}
	}
}
