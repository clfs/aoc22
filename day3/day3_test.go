package day3

import (
	"bytes"
	"testing"

	"github.com/clfs/aoc22"
)

func TestPriority(t *testing.T) {
	cases := []struct {
		in   rune
		want int
	}{
		{'a', 1},
		{'z', 26},
		{'A', 27},
		{'Z', 52},
	}
	for _, c := range cases {
		got := Priority(c.in)
		if got != c.want {
			t.Errorf("Priority(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestPart1(t *testing.T) {
	data := aoc22.ReadTestFile(t, "testdata/input.txt")

	got := Part1(bytes.NewReader(data))
	want := 8109

	if got != want {
		t.Errorf("Part1() = %d, want %d", got, want)
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		path string
		want int
	}{
		{"testdata/input.txt", 2738},
		{"testdata/small.txt", 70},
	}

	for _, c := range cases {
		data := aoc22.ReadTestFile(t, c.path)
		got := Part2(bytes.NewReader(data))
		if got != c.want {
			t.Errorf("Part2(%q) = %d, want %d", c.path, got, c.want)
		}
	}
}
