package day7

import (
	"bytes"
	"os"
	"testing"

	"github.com/clfs/aoc22"
	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	f, err := os.Open("testdata/small.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	want := map[string]int64{
		"/":        0,
		"/a":       0,
		"/a/e":     0,
		"/a/e/i":   584,
		"/a/f":     29116,
		"/a/g":     2557,
		"/a/h.lst": 62596,
		"/b.txt":   14848514,
		"/c.dat":   8504156,
		"/d":       0,
		"/d/d.ext": 5626152,
		"/d/d.log": 8033020,
		"/d/j":     4060174,
		"/d/k":     7214296,
	}

	got, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Parse() mismatch (-want, got):\n%s", diff)
	}

	// Also test DirSize while we're here.

	wantSizes := map[string]int64{
		"/":    48381165,
		"/a":   94853,
		"/a/e": 584,
		"/d":   24933642,
	}

	for name, want := range wantSizes {
		got := DirSize(name, got)
		if got != want {
			t.Errorf("DirSize(%q) = %d, want %d", name, got, want)
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int64
	}{
		{"testdata/input.txt", 1778099},
		{"testdata/small.txt", 95437},
	}

	for _, c := range cases {
		data := aoc22.ReadTestFile(t, c.name)
		got, err := Part1(bytes.NewReader(data))
		if err != nil {
			t.Errorf("%q: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("%q: got %d, want %d", c.name, got, c.want)
		}
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		name string
		want int64
	}{
		{"testdata/input.txt", 1623571},
		{"testdata/small.txt", 24933642},
	}

	for _, c := range cases {
		data := aoc22.ReadTestFile(t, c.name)
		got, err := Part2(bytes.NewReader(data))
		if err != nil {
			t.Errorf("%q: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("%q: got %d, want %d", c.name, got, c.want)
		}
	}
}
