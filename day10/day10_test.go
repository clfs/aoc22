package day10

import (
	"os"
	"testing"
)

func TestOp_UnmarshalText(t *testing.T) {
	cases := []struct {
		in   string
		want Op
	}{
		{"noop", Op{"noop", 0}},
		{"addx 11", Op{"addx", 11}},
		{"addx -11", Op{"addx", -11}},
	}

	for _, c := range cases {
		var got Op
		err := got.UnmarshalText([]byte(c.in))
		if err != nil {
			t.Errorf("UnmarshalText(%q) returned error %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("UnmarshalText(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func readProgram(t *testing.T, name string) Program {
	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	p, err := ParseProgram(f)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCPU_Tick(t *testing.T) {
	cases := []struct {
		name    string
		outputs []int // may be truncated
	}{
		{"testdata/small.txt", []int{1, 1, 1, 4, 4}},
		{"testdata/large.txt", []int{1, 1, 16, 16, 5}},
	}

	for _, tc := range cases {
		var cpu CPU
		cpu.Load(readProgram(t, tc.name))

		for i, want := range tc.outputs {
			if got := cpu.Tick(); got != want {
				t.Errorf(
					"%q: tick #%d = %v, want %v",
					tc.name, i, got, want,
				)
			}
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/large.txt", 13140},
		{"testdata/input.txt", 15220},
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
				t.Errorf("Part1() error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Part1() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"testdata/large.txt", "testdata/large_crt.txt"},
		{"testdata/input.txt", "testdata/input_crt.txt"},
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
				t.Errorf("Part2() error: %v", err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			if got != string(want) {
				t.Errorf("Part2() mismatch:\ngot:\n%v\nwant:\n%s", got, want)
			}
		})
	}

}
