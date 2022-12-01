package day1

import (
	"bytes"
	"testing"

	"github.com/clfs/aoc22"
)

func TestPart1(t *testing.T) {
	data := aoc22.ReadTestFile(t, "testdata/input.txt")

	got := Part1(bytes.NewReader(data))
	want := 71924

	if got != want {
		t.Errorf("Part1() = %d, want %d", got, want)
	}
}

func TestPart2(t *testing.T) {
	data := aoc22.ReadTestFile(t, "testdata/input.txt")

	got := Part2(bytes.NewReader(data))
	want := 210406

	if got != want {
		t.Errorf("Part2() = %d, want %d", got, want)
	}
}
