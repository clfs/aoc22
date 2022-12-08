package day8

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewForest(t *testing.T) {
	f, err := os.Open("testdata/small.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err := NewForest(f)
	if err != nil {
		t.Fatal(err)
	}

	want := &Forest{Grid: [][]int{
		{3, 0, 3, 7, 3},
		{2, 5, 5, 1, 2},
		{6, 5, 3, 3, 2},
		{3, 3, 5, 4, 9},
		{3, 5, 3, 9, 0},
	}}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewForest() mismatch (-want +got):\n%s", diff)
	}
}

func readForest(t *testing.T, name string) *Forest {
	t.Helper()

	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	forest, err := NewForest(f)
	if err != nil {
		t.Fatal(err)
	}

	return forest
}

func TestForest_IsVisibleInDirection(t *testing.T) {
	forest := readForest(t, "testdata/small.txt")

	cases := []struct {
		r, c   int
		dr, dc int
		want   bool
	}{
		{1, 1, 0, -1, true},  // top-left 5, left
		{1, 1, -1, 0, true},  // top-left 5, up
		{1, 1, 1, 0, false},  // top-left 5, down
		{1, 1, 0, 1, false},  // top-left 5, right
		{1, 2, 0, 1, true},   // top-middle 5, right
		{1, 2, -1, 0, true},  // top-middle 5, up
		{1, 2, 0, -1, false}, // top-middle 5, left
		{1, 2, 1, 0, false},  // top-middle 5, down
	}

	for _, tc := range cases {
		got := forest.IsVisibleInDirection(tc.r, tc.c, tc.dr, tc.dc)
		if got != tc.want {
			t.Errorf("Forest.IsVisibleInDirection(%d, %d, %d, %d) = %t, want %t", tc.r, tc.c, tc.dr, tc.dc, got, tc.want)
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 21},
		{"testdata/input.txt", 1796},
	}

	for _, tc := range cases {
		f := readForest(t, tc.name)

		got, err := Part1(f)
		if err != nil {
			t.Errorf("%q: %v", tc.name, err)
		}
		if got != tc.want {
			t.Errorf("%q: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestForest_ScenicScore(t *testing.T) {
	f := readForest(t, "testdata/small.txt")

	cases := []struct {
		r, c int
		want int
	}{
		{1, 2, 4}, // middle 5 in second row
		{3, 2, 8}, // 5 in middle of fourth row
	}

	for _, tc := range cases {
		got := f.ScenicScore(tc.r, tc.c)
		if got != tc.want {
			t.Errorf("Forest.ScenicScore(%d, %d) = %d, want %d", tc.r, tc.c, got, tc.want)
		}
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 8},
		{"testdata/input.txt", 288120},
	}

	for _, tc := range cases {
		f := readForest(t, tc.name)
		if got := Part2(f); got != tc.want {
			t.Errorf("%q: got %d, want %d", tc.name, got, tc.want)
		}
	}
}
