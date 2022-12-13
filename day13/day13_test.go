package day13

import (
	"os"
	"testing"
)

func toPacket(t *testing.T, s string) []any {
	t.Helper()
	packet, err := ParsePacket([]byte(s))
	if err != nil {
		t.Fatalf("invalid packet %q: %v", s, err)
	}
	return packet
}

func TestCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{
			"[1,1,3,1,1]",
			"[1,1,5,1,1]",
			-1,
		},
		{
			"[[1],[2,3,4]]",
			"[[1],4]",
			-1,
		},
		{
			"[[4,4],4,4]",
			"[[4,4],4,4,4]",
			-1,
		},
	}

	for _, c := range cases {
		a, b := toPacket(t, c.a), toPacket(t, c.b)
		got := Compare(a, b)
		if got != c.want {
			t.Errorf(
				"Compare(%q, %q) = %d, want %d",
				c.a, c.b, got, c.want,
			)
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 13},
		{"testdata/input.txt", 6478},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatalf("failed to open %q: %v", tc.name, err)
			}
			defer f.Close()

			got, err := Part1(f)
			if err != nil {
				t.Fatalf("failed to run Part1: %v", err)
			}
			if got != tc.want {
				t.Errorf("Part1() = %d, want %d", got, tc.want)
			}
		})
	}
}

func FuzzCompare(f *testing.F) {
	f.Fuzz(func(t *testing.T, left, right []byte) {
		a, err := ParsePacket(left)
		if err != nil {
			t.Skip()
		}
		b, err := ParsePacket(right)
		if err != nil {
			t.Skip()
		}

		if n := Compare(a, a); n != 0 {
			t.Errorf("a = %s, cmp(a,a) = %d", a, n)
		}
		if n := Compare(b, b); n != 0 {
			t.Errorf("b = %s, cmp(b,b) = %d", b, n)
		}

		x := Compare(a, b) // -1, 0, 1
		y := Compare(b, a) // 1, 0, -1

		if x+y != 0 {
			t.Errorf(
				"a = %s, b = %s, cmp(a,b) = %d, cmp(b,a) = %d",
				a, b, x, y,
			)
		}
	})
}

func TestPart2(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 140},
		{"testdata/input.txt", 21922},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatalf("failed to open %q: %v", tc.name, err)
			}
			defer f.Close()

			got, err := Part2(f)
			if err != nil {
				t.Fatalf("failed to run Part2: %v", err)
			}
			if got != tc.want {
				t.Errorf("Part2() = %d, want %d", got, tc.want)
			}
		})
	}
}
