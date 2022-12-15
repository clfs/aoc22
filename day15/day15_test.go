package day15

import (
	"os"
	"testing"

	"golang.org/x/exp/slices"
)

func TestSensor_Radius(t *testing.T) {
	cases := []struct {
		s    Sensor
		want int
	}{
		{Sensor{Point{2, 2}, Point{3, 3}}, 2},
		{Sensor{Point{-2, 0}, Point{0, -1}}, 3},
	}

	for _, tc := range cases {
		if got := tc.s.Radius(); got != tc.want {
			t.Errorf("%v.Radius() = %d, want %d", tc.s, got, tc.want)
		}
	}
}

func TestLenUnion(t *testing.T) {
	cases := []struct {
		in   []Range
		want int
	}{
		{[]Range{{0, 1}, {2, 3}}, 4},
		{[]Range{{0, 1}, {2, 3}, {4, 5}}, 6},
		{[]Range{{0, 1}, {0, 2}}, 3},
		{[]Range{{-1, 0}, {0, 1}}, 3},
		{[]Range{{-10, 10}, {-10, 10}}, 21},
		{[]Range{{-10, 10}, {-10, 10}, {-10, 10}}, 21},
		{[]Range{}, 0},
		{[]Range{{0, 1}, {-1, 2}}, 4},
	}

	for _, tc := range cases {
		if got := LenUnion(Union(tc.in)); got != tc.want {
			t.Errorf("LenUnion(%v) = %d; want %d", tc.in, got, tc.want)
		}
	}
}

func TestRange_Len(t *testing.T) {
	cases := []struct {
		r    Range
		want int
	}{
		{Range{0, 1}, 2},
		{Range{0, 2}, 3},
		{Range{0, 0}, 1},
		{Range{-1, 0}, 2},
		{Range{-1, -1}, 1},
		{Range{-2, -1}, 2},
	}

	for _, tc := range cases {
		if got := tc.r.Len(); got != tc.want {
			t.Errorf("%v.Len() = %d; want %d", tc.r, got, tc.want)
		}
	}
}

func TestUnion(t *testing.T) {
	cases := []struct {
		in   []Range
		want []Range
	}{
		{
			[]Range{{0, 1}, {2, 3}},
			[]Range{{0, 3}},
		},
		{
			[]Range{{0, 1}, {2, 3}, {4, 5}},
			[]Range{{0, 5}},
		},
		{
			[]Range{{0, 1}, {0, 2}},
			[]Range{{0, 2}},
		},
		{
			[]Range{{-1, 0}, {0, 1}},
			[]Range{{-1, 1}},
		},
		{
			[]Range{{-10, 10}, {-10, 10}},
			[]Range{{-10, 10}},
		},
		{
			[]Range{{-10, 10}, {-10, 10}, {-10, 10}},
			[]Range{{-10, 10}},
		},
		{
			[]Range{},
			[]Range{},
		},
		{
			[]Range{{0, 1}, {-1, 2}},
			[]Range{{-1, 2}},
		},
		{
			[]Range{{0, 1}, {4, 5}},
			[]Range{{0, 1}, {4, 5}},
		},
	}

	for _, tc := range cases {
		in := tc.in // Union modifies its argument if unsorted
		if got := Union(in); !slices.Equal(got, tc.want) {
			t.Errorf("Union of %v = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestSensor_Intersect(t *testing.T) {
	cases := []struct {
		s      Sensor
		y      int
		want   Range
		wantOk bool
	}{
		{
			Sensor{Point{0, 0}, Point{0, 2}},
			-1,
			Range{-1, 1},
			true,
		},
		{
			Sensor{Point{0, 0}, Point{0, 2}},
			0,
			Range{-2, 2},
			true,
		},
	}

	for _, tc := range cases {
		got, gotOk := tc.s.Intersect(tc.y)
		if gotOk != tc.wantOk {
			t.Errorf(
				"%v.Intersect(%d) = %v, %t; want %v, %t",
				tc.s, tc.y, got, gotOk, tc.want, tc.wantOk,
			)
		}
		if gotOk && got != tc.want {
			t.Errorf(
				"%v.Intersect(%d) = %v, %t; want %v, %t",
				tc.s, tc.y, got, gotOk, tc.want, tc.wantOk,
			)
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		y    int
		want int
	}{
		{"testdata/small.txt", 10, 26},
		{"testdata/input.txt", 2000000, 4748135},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			got, err := Part1(f, tc.y)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			if got != tc.want {
				t.Errorf(
					"Part1(%q, %d) = %d; want %d",
					tc.name, tc.y, got, tc.want,
				)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		name  string
		bound int
		want  int
	}{
		{"testdata/small.txt", 20, 56000011},
		{"testdata/input.txt", 4000000, 13743542639657},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			got, err := Part2(f, tc.bound)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			if got != tc.want {
				t.Errorf(
					"Part2(%q, %d) = %d; want %d",
					tc.name, tc.bound, got, tc.want,
				)
			}
		})
	}
}
