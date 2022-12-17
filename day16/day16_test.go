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
			Valve{"AA", 0, []string{"DD", "II", "BB"}},
		},
		{
			"Valve OV has flow rate=10; tunnels lead to valves YW, JT, NN, TK",
			Valve{"OV", 10, []string{"YW", "JT", "NN", "TK"}},
		},
		{
			"Valve HH has flow rate=22; tunnel leads to valve GG",
			Valve{"HH", 22, []string{"GG"}},
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

func TestNewVolcano(t *testing.T) {
	valves := readValves(t, "testdata/small.txt")
	limit := 30

	want := &Volcano{
		Nodes: []string{"AA", "BB", "CC", "DD", "EE", "FF", "GG", "HH", "II", "JJ"},
		Edges: map[string][]string{
			"AA": {"BB", "DD", "II"},
			"BB": {"AA", "CC"},
			"CC": {"BB", "DD"},
			"DD": {"AA", "CC", "EE"},
			"EE": {"DD", "FF"},
			"FF": {"EE", "GG"},
			"GG": {"FF", "HH"},
			"HH": {"GG"},
			"II": {"AA", "JJ"},
			"JJ": {"II"},
		},
		Rates: map[string]int{
			"AA": 0,
			"BB": 13,
			"CC": 2,
			"DD": 20,
			"EE": 3,
			"FF": 0,
			"GG": 0,
			"HH": 22,
			"II": 0,
			"JJ": 21,
		},
		Status: map[string]bool{
			"AA": false,
			"BB": false,
			"CC": false,
			"DD": false,
			"EE": false,
			"FF": false,
			"GG": false,
			"HH": false,
			"II": false,
			"JJ": false,
		},
		Location:     "AA",
		TimeLimit:    30,
		LenPathCache: map[string]map[string]int{},
	}
	got := NewVolcano(valves, limit)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewVolcano() mismatch (-want,+got):%s\n", diff)
	}
}

func TestVolcano_Path(t *testing.T) {
	volcano := NewVolcano(readValves(t, "testdata/small.txt"), 30)

	cases := []struct {
		from, to string
		want     []string
	}{
		{"AA", "AA", []string{"AA"}},
		{"AA", "DD", []string{"AA", "DD"}},
		{"AA", "II", []string{"AA", "II"}},
		{"AA", "BB", []string{"AA", "BB"}},
		{"AA", "CC", []string{"AA", "BB", "CC"}},
		{"GG", "DD", []string{"GG", "FF", "EE", "DD"}},
	}

	for _, tc := range cases {
		got, ok := volcano.Path(tc.from, tc.to)
		if !ok {
			t.Errorf("Volcano.Path(%q, %q) not ok", tc.from, tc.to)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf(
				"Volcano.Path(%q, %q) mismatch (-want,+got):%s\n",
				tc.from, tc.to, diff,
			)
		}
	}
}

func TestPart1(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 1651},
		{"testdata/input.txt", 2119},
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

func TestPart2(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"testdata/small.txt", 1707},
		//{"testdata/input.txt", 2615},
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
				t.Errorf("error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestVolcano_Evaluate(t *testing.T) {
	volcano := NewVolcano(readValves(t, "testdata/small.txt"), 30)

	cases := []struct {
		path []string
		want int
	}{
		// testdata/small.txt
		{[]string{"DD", "BB", "JJ", "HH", "EE", "CC"}, 1651},
		{[]string{"DD"}, 560},
		{[]string{"JJ"}, 567},
		{[]string{"DD", "BB"}, 33*28 - ((33 - 20) * 3)},

		// testdata/input.txt
		// {[]string{"FJ", "EL", "ST", "PF", "MD", "DK", "LR"}, 1975},
		// {[]string{"OV", "FJ", "EL", "ST", "PF", "MD", "DK", "FK", "LR", "XX"}, 0},
	}

	for _, tc := range cases {
		got := volcano.Evaluate(tc.path)
		if got != tc.want {
			t.Errorf("%v: got %d, want %d", tc.path, got, tc.want)
		}
	}
}
