package day16

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type Valve struct {
	Name    string
	Rate    int
	Open    bool
	Tunnels []string
}

var ValveRegexp = regexp.MustCompile(`Valve ([A-Z]{2}) has flow rate=(\d+); (?:tunnels lead|tunnel leads) to (?:valves|valve) ([A-Z]{2}(?:, [A-Z]{2})*)`)

func (v *Valve) UnmarshalText(text []byte) error {
	// Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
	matches := ValveRegexp.FindStringSubmatch(string(text))
	if matches == nil {
		return fmt.Errorf("no match for %q", text)
	}
	v.Name = matches[1]
	fmt.Sscanf(matches[2], "%d", &v.Rate)
	v.Tunnels = strings.Split(matches[3], ", ")
	return nil
}

type Volcano struct {
	Valves      map[string]Valve
	Location    string
	Released    int
	TimeLimit   int
	TimeElapsed int
}

func NewVolcano(valves []Valve, timeLimit int) *Volcano {
	m := make(map[string]Valve)
	for _, valve := range valves {
		m[valve.Name] = valve
	}
	return &Volcano{
		Valves:    m,
		Location:  "AA",
		TimeLimit: timeLimit,
	}
}

// Tick advances the volcano simulation by one minute.
func (v *Volcano) Tick() {
	defer func() { v.TimeElapsed++ }()

	log.Printf("==== minute %d ====", v.TimeElapsed)

	// Release flow from open valves.
	var released int
	for _, valve := range v.Valves {
		if valve.Open {
			released += valve.Rate
		}
	}
	if released > 0 {
		log.Printf("valve(s) %s released %d pressure", strings.Join(v.openValves(), ", "), released)
		v.Released += released
	}

	// Find the shortest path to each valve potentially worth opening.
	var paths [][]string
	for name, valve := range v.Valves {
		if valve.Open {
			continue // already open
		}
		if valve.Rate == 0 {
			continue // useless
		}
		path, ok := v.PathTo(name)
		if !ok {
			continue // unreachable
		}
		paths = append(paths, path)
	}

	bestPath := v.bestOfPaths(paths)

	// If no worthwhile paths exist, wait out the clock.
	if bestPath == nil {
		log.Print("no worthwhile actions found")
		return
	}

	// If the best path is to the current location, open the current valve.
	if len(bestPath) == 1 {
		log.Printf("opening valve %s", v.Location)
		v.OpenLocation()
		return
	}

	// Otherwise, the best path is to a different location.
	// Move to the next location in the path.
	log.Printf("moving to %s", bestPath[1])
	v.Location = bestPath[1]
}

// OpenLocation opens the valve at the current location.
func (v *Volcano) OpenLocation() {
	named, ok := v.Valves[v.Location]
	if !ok {
		panic(fmt.Sprintf("bad location %q", v.Location))
	}
	named.Open = true
	v.Valves[v.Location] = named
}

// IsOpen returns true if the named valve is open.
func (v *Volcano) IsOpen(name string) bool {
	return v.Valves[name].Open
}

// IsClosed returns true if the named valve is closed.
func (v *Volcano) IsClosed(name string) bool {
	return !v.Valves[name].Open
}

// openValves returns the names of all open valves, in alphabetical order.
func (v *Volcano) openValves() []string {
	var open []string
	for name, valve := range v.Valves {
		if valve.Open {
			open = append(open, name)
		}
	}
	slices.Sort(open)
	return open
}

// bestOfPaths returns the best path to take from multiple options.
// If no paths are provided, it returns nil.
func (v *Volcano) bestOfPaths(paths [][]string) []string {
	var best []string
	for _, path := range paths {
		if best == nil || len(path) < len(best) {
			best = path
		}
	}
	return best
}

// TimeToOpen returns the time it would take to open the named valve.
// If the valve is already open, it returns 0.
// If the valve isn't reachable, ok is false.
func (v *Volcano) TimeToOpen(name string) (time int, ok bool) {
	if v.IsOpen(name) {
		return 0, true
	}

	path, ok := v.PathTo(name)
	if !ok {
		return 0, false
	}

	return len(path), true
}

// PathTo returns the shortest path from the current location to the given valve.
// It includes both endpoints. If no path exists, ok is false.
func (v *Volcano) PathTo(name string) (path []string, ok bool) {
	// BFS.
	visited := make(map[string]bool)
	queue := [][]string{{v.Location}}
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		last := path[len(path)-1]
		if last == name {
			return path, true
		}

		if visited[last] {
			continue
		}
		visited[last] = true

		for _, tunnel := range v.Valves[last].Tunnels {
			queue = append(queue, append(path, tunnel))
		}
	}
	return nil, false
}

// Run runs the simulation until the time limit is reached.
// It returns the total amount of pressure released.
func (v *Volcano) Run() int {
	for v.TimeElapsed < v.TimeLimit {
		v.Tick()
	}
	return v.Released
}

func Parse(r io.Reader) ([]Valve, error) {
	var valves []Valve
	s := bufio.NewScanner(r)
	for s.Scan() {
		var v Valve
		if err := v.UnmarshalText(s.Bytes()); err != nil {
			return nil, err
		}
		valves = append(valves, v)
	}
	return valves, s.Err()
}

func Part1(r io.Reader) (int, error) {
	valves, err := Parse(r)
	if err != nil {
		return 0, err
	}
	return NewVolcano(valves, 30).Run(), nil
}
