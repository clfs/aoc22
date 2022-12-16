package day16

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValve_UnmarshalText(t *testing.T) {
	cases := []struct {
		in   string
		want Valve
	}{
		{
			"Valve AA has flow rate=0; tunnels lead to valves DD, II, BB",
			Valve{"AA", 0, false, []string{"DD", "II", "BB"}},
		},
		{
			"Valve OV has flow rate=10; tunnels lead to valves YW, JT, NN, TK",
			Valve{"OV", 10, false, []string{"YW", "JT", "NN", "TK"}},
		},
		{
			"Valve HH has flow rate=22; tunnel leads to valve GG",
			Valve{"HH", 22, false, []string{"GG"}},
		},
	}

	for _, tc := range cases {
		var got Valve
		err := got.UnmarshalText([]byte(tc.in))
		if err != nil {
			t.Errorf("UnmarshalText(%q) failed: %v", tc.in, err)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Value.UnmarshalText(%q) mismatch (-want,+got):%s\n", tc.in, diff)
		}
	}
}

func readValves(t *testing.T, name string) []Valve {
	t.Helper()
	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("failed to open %q: %v", name, err)
	}
	defer f.Close()

	valves, err := Parse(f)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", name, err)
	}
	return valves
}

func TestVolcano_TimePathTo(t *testing.T) {
	valves := readValves(t, "testdata/small.txt")
	volcano := NewVolcano(valves, 30)

	cases := []struct {
		in   string
		want []string
	}{
		{"AA", []string{"AA"}},
		{"DD", []string{"AA", "DD"}},
		{"JJ", []string{"AA", "II", "JJ"}},
	}

	for _, tc := range cases {
		got, ok := volcano.PathTo(tc.in)
		if !ok {
			t.Errorf("Volcano.PathTo(%q): not ok", tc.in)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Volcano.PathTo(%q) mismatch (-want,+got):%s\n", tc.in, diff)
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 1651},
		//{"testdata/input.txt", 1},
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
				t.Errorf("error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}
