package day12

import (
	"os"
	"testing"
)

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 31},
		{"testdata/input.txt", 497},
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
				t.Errorf("error: %v", err)
			}
			if got != tc.want {
				t.Errorf("%v, want %v", got, tc.want)
			}
		})
	}
}

func readTopo(t *testing.T, name string) *Topo {
	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	topo, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	return &topo
}

func TestCanMove(t *testing.T) {
	topo := readTopo(t, "testdata/small.txt")
	cases := []struct {
		from, to Point
		want     bool
	}{
		{Point{0, 0}, Point{0, 1}, true},
		{Point{0, 0}, Point{1, 0}, true},
		{Point{0, 0}, Point{1, 1}, false},
	}

	for _, tc := range cases {
		if got := topo.CanMove(tc.from, tc.to); got != tc.want {
			t.Errorf("CanMove(%v, %v) = %v, want %v", tc.from, tc.to, got, tc.want)
		}
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 29},
		{"testdata/input.txt", 492},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			got, err := Part2(f)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			if got != tc.want {
				t.Errorf("%v, want %v", got, tc.want)
			}
		})
	}
}
