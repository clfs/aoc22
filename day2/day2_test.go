package day2

import (
	"bytes"
	"testing"

	"github.com/clfs/aoc22"
)

func TestPart1(t *testing.T) {
	data := aoc22.ReadTestFile(t, "testdata/input.txt")

	got := Part1(bytes.NewReader(data))
	want := 15523

	if got != want {
		t.Errorf("Part1() = %d, want %d", got, want)
	}
}

func TestPart2(t *testing.T) {
	data := aoc22.ReadTestFile(t, "testdata/input.txt")

	got := Part2(bytes.NewReader(data))
	want := 15702

	if got != want {
		t.Errorf("Part2() = %d, want %d", got, want)
	}
}

func TestRound(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"A X", 4},
		{"A Y", 8},
		{"A Z", 3},
		{"B X", 1},
		{"B Y", 5},
		{"B Z", 9},
		{"C X", 7},
		{"C Y", 2},
		{"C Z", 6},
	}
	for _, c := range cases {
		var r Round
		if err := r.UnmarshalText([]byte(c.in)); err != nil {
			t.Errorf("Round.UnmarshalText(%q) = %v", c.in, err)
		}
		got := r.Score()
		if got != c.want {
			t.Errorf("Round.Score(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestRoundAlt(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"A X", 3},
		{"A Y", 4},
		{"A Z", 8},
		{"B X", 1},
		{"B Y", 5},
		{"B Z", 9},
		{"C X", 2},
		{"C Y", 6},
		{"C Z", 7},
	}
	for _, c := range cases {
		var r RoundAlt
		if err := r.UnmarshalText([]byte(c.in)); err != nil {
			t.Errorf("RoundAlt.UnmarshalText(%q) = %v", c.in, err)
		}
		got := r.Score()
		if got != c.want {
			t.Errorf("RoundAlt.Score(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}
