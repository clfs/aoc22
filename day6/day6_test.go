package day6

import (
	"testing"

	"github.com/clfs/aoc22"
)

func TestPart1_Examples(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 7},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 5},
		{"nppdvjthqldpwncqszvftbrmjlhg", 6},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 10},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 11},
	}

	for _, c := range cases {
		got := Part1(c.in)
		if got != c.want {
			t.Errorf("Part1(%q) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestPart1(t *testing.T) {
	data := aoc22.ReadTestFile(t, "testdata/input.txt")

	got := Part1(string(data))
	want := 1544

	if got != want {
		t.Errorf("Part1() == %d, want %d", got, want)
	}
}

func TestPart2_Examples(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 19},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 23},
		{"nppdvjthqldpwncqszvftbrmjlhg", 23},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 29},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 26},
	}

	for _, c := range cases {
		got := Part2(c.in)
		if got != c.want {
			t.Errorf("Part2(%q) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestPart2(t *testing.T) {
	data := aoc22.ReadTestFile(t, "testdata/input.txt")

	got := Part2(string(data))
	want := 2145

	if got != want {
		t.Errorf("Part2() == %d, want %d", got, want)
	}
}
