package day11

import (
	"os"
	"testing"
)

func TestParseOperation(t *testing.T) {
	cases := []struct {
		in   string
		want func(int) int
	}{
		{"+ 12", func(old int) int { return old + 12 }},
		{"* 12", func(old int) int { return old * 12 }},
		{"* old", func(old int) int { return old * old }},
	}

	for _, c := range cases {
		got := ParseOperation(c.in)
		for i := 0; i < 3; i++ {
			if got(i) != c.want(i) {
				t.Errorf("ParseOperation(%q)(%d) = %d, want %d", c.in, i, got(i), c.want(i))
			}
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 10605},
		{"testdata/input.txt", 76728},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			got, err := Part1(f)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("Part1() = %d, want %d", got, tc.want)
			}
		})
	}

}
