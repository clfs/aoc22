package day16

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Valve struct {
	Name    string
	Rate    int
	Tunnels []string
}

var ValveRegexp = regexp.MustCompile(`Valve ([A-Z]{2}) has flow rate=(\d+); (?:tunnels lead|tunnel leads) to (?:valves|valve) ([A-Z]{2}(?:, [A-Z]{2})*)`)

func (v *Valve) UnmarshalText(text []byte) error {
	matches := ValveRegexp.FindStringSubmatch(string(text))
	if matches == nil {
		return fmt.Errorf("no match for %q", text)
	}

	v.Name = matches[1]
	v.Tunnels = strings.Split(matches[3], ", ")

	var err error
	v.Rate, err = strconv.Atoi(matches[2])
	return err
}

type Volcano struct {
	Nodes  []string            // The valves.
	Edges  map[string][]string // A map from valves to their connected valves.
	Rates  map[string]int      // A map from valves to their flow rate.
	Status map[string]bool     // A map from valves to their open status. Closed is false.

	Location string // The valve you're standing at.

	TimeElapsed int // How many minutes have passed.
	TimeLimit   int // The time limit before the volcano explodes.

	// Cached results for Volcano.LenPath.
	LenPathCache map[string]map[string]int
}

func NewVolcano(vs []Valve, limit int) *Volcano {
	volcano := &Volcano{
		Edges:        make(map[string][]string),
		Rates:        make(map[string]int),
		Status:       make(map[string]bool),
		Location:     "AA",
		TimeLimit:    limit,
		LenPathCache: make(map[string]map[string]int),
	}

	for _, v := range vs {
		volcano.Nodes = append(volcano.Nodes, v.Name)

		// For easier testing, sort the tunnels alphabetically.
		slices.Sort(v.Tunnels)
		volcano.Edges[v.Name] = append(volcano.Edges[v.Name], v.Tunnels...)

		volcano.Rates[v.Name] = v.Rate
		volcano.Status[v.Name] = false
	}

	return volcano
}

func (v *Volcano) Rate(name string) int {
	rate, ok := v.Rates[name]
	if !ok {
		panic("bad name " + name)
	}
	return rate
}

func (v *Volcano) Open(name string) {
	status, ok := v.Status[name]
	if !ok {
		panic("bad name " + name)
	}
	if status {
		panic(name + " already open")
	}
	v.Status[name] = true
}

func (v *Volcano) IsOpen(name string) bool {
	status, ok := v.Status[name]
	if !ok {
		panic("bad name " + name)
	}
	return status
}

func (v *Volcano) IsClosed(name string) bool {
	return !v.IsOpen(name)
}

// NextFrom returns the valves that are next from the given valve.
func (v *Volcano) NextFrom(name string) []string {
	edges, ok := v.Edges[name]
	if !ok {
		panic("bad name " + name)
	}
	return edges
}

// Path returns the shortest path between from and to.
func (v *Volcano) Path(from, to string) ([]string, bool) {
	visited := map[string]bool{from: true}
	queue := [][]string{{from}}

	for len(queue) > 0 {
		var nextQueue [][]string
		for _, path := range queue {
			tail := path[len(path)-1]

			visited[tail] = true

			if tail == to {
				return path, true
			}

			for _, next := range v.NextFrom(tail) {
				if visited[next] {
					continue
				}
				tmp := make([]string, len(path))
				copy(tmp, path)
				nextQueue = append(nextQueue, append(tmp, next))
			}
		}
		queue = nextQueue
	}

	return nil, false
}

// BestMove returns the best next move. It returns either the current
// valve (indicating it should be opened), or the next valve to step to.
//
// If it's no longer useful to move, ok is false.
func (v *Volcano) BestMove() (string, bool) {
	// Get the shortest path to every useful valve.
	paths := make(map[string][]string)
	for _, name := range v.Nodes {
		if v.IsOpen(name) || v.Rate(name) == 0 {
			continue // useless
		}
		path, ok := v.Path(v.Location, name)
		if !ok {
			continue // unreachable
		}
		paths[name] = path
	}

	// Find the path with the best score.
	var (
		bestScore int
		bestPath  []string
	)
	for _, path := range paths {
		score := v.Score(path)
		if score > bestScore {
			bestScore = score
			bestPath = path
		}
	}

	log.Printf("- best path with score %d: %v", bestScore, bestPath)

	switch l := len(bestPath); l {
	case 0:
		return "", false // no more closed valves
	case 1:
		return bestPath[0], true // open the current valve
	default:
		return bestPath[1], true // move to the next valve
	}
}

// Score returns the amount of pressure that would be relieved before
// the volcano explodes, if you followed the path and opened the valve
// at the end.
//
// If the valve at the end is already open, or if you can't reach the
// valve before the volcano explodes, the score is 0.
//
// If the path is empty, Score panics.
func (v *Volcano) Score(path []string) int {
	timeNeeded := len(path)
	timeLeft := v.TimeLimit - v.TimeElapsed

	tail := path[len(path)-1]
	tailRate := v.Rate(tail)

	if v.IsOpen(tail) {
		return 0
	}

	if timeNeeded >= timeLeft {
		return 0
	}

	return (timeLeft - timeNeeded) * tailRate
}

// Tick advances the simulation by one minute.
// It returns the pressure released in that minute.
func (v *Volcano) Tick() int {
	defer func() { v.TimeElapsed++ }()

	log.Printf("==== minute %d ====", v.TimeElapsed)

	var (
		pressure   int
		openValves []string
	)
	for _, name := range v.Nodes {
		if v.IsOpen(name) {
			openValves = append(openValves, name)
			pressure += v.Rate(name)
		}
	}
	if len(openValves) > 0 {
		log.Printf(
			"open valve(s) %s released %d pressure",
			strings.Join(openValves, ", "), pressure,
		)
	}

	move, ok := v.BestMove()
	if ok {
		if move == v.Location {
			log.Printf("opening valve %s", move)
			v.Open(move)
		} else {
			log.Printf("moving to valve %s", move)
			v.Location = move
		}
	}

	return pressure
}

// Run runs the simulation until the time limit is reached.
// It returns the total amount of pressure released.
func (v *Volcano) Run() int {
	var sum int
	for v.TimeElapsed < v.TimeLimit {
		sum += v.Tick()
	}
	return sum
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
	return NewVolcano(valves, 30).Solve2(), nil
}

// RandSample returns a sample of size n from pop. It alters the order
// of elements in pop.
func RandSample(pop []string, n int) ([]string, bool) {
	if n > len(pop) {
		return nil, false
	}
	rand.Shuffle(len(pop), func(i, j int) { pop[i], pop[j] = pop[j], pop[i] })
	return pop[:n], true
}

func (v *Volcano) Solve() int {
	var targets []string
	for _, name := range v.Nodes {
		if v.IsClosed(name) && v.Rate(name) > 0 {
			targets = append(targets, name)
		}
	}

	var (
		bestScore  int
		bestSample []string
	)

	n := min(len(targets), 9)

	for i := 0; i < 1000000000; i++ {
		rand.Shuffle(len(targets), func(i, j int) { targets[i], targets[j] = targets[j], targets[i] })
		score := v.Evaluate(targets[:n])
		if score > bestScore {
			bestScore = score
			tmp := make([]string, n)
			copy(tmp, targets[:n])
			bestSample = tmp
			log.Printf("score %d with sample %v", bestScore, bestSample)
		}
	}

	return bestScore
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Evaluate determines how much pressure would be released if you
// visited each target in order and opened them.
//
// It returns 0 if you can't reach all targets before the volcano explodes.
func (v *Volcano) Evaluate(targets []string) int {
	var (
		timeElapsed = 0    // how much time has passed
		burnRate    = 0    // pressure burnt per minute
		totalBurnt  = 0    // total pressure burnt
		location    = "AA" // current location
	)

	for _, target := range targets {
		// log.Printf("*** %d min elapsed, location = %s", timeElapsed, location)

		// Get the length of the path to the next target.
		lenPath := v.LenPath(location, target)

		// log.Printf("%d edges to next target", lenPath)

		// If no path exists, or the path would take too long, stand
		// in place and finish.
		if lenPath == -1 || timeElapsed+lenPath >= v.TimeLimit {
			// log.Print("didn't get to all targets...")
			return 0
		}

		// Move to the next target.
		// log.Printf("moving to %s over %d minute(s)", target, lenPath)
		location = target
		timeElapsed += lenPath
		totalBurnt += burnRate * lenPath

		// Open it.
		totalBurnt += burnRate
		// log.Printf("opening %s over 1 minute", location)
		burnRate += v.Rate(location)
		timeElapsed += 1
	}

	// Add anything left...
	// log.Printf("*** %d min elapsed, location = %s", timeElapsed, location)
	totalBurnt += burnRate * (v.TimeLimit - timeElapsed)

	return totalBurnt
}

func (v *Volcano) LenPathCached(from, to string) (int, bool) {
	f, ok := v.LenPathCache[from]
	if !ok {
		return 0, false
	}

	t, ok := f[to]
	if !ok {
		return 0, false
	}

	return t, true
}

// LenPath returns the number of edges between from and to.
// If from == to, it returns 0.
// If no path exists, it returns -1.
func (v *Volcano) LenPath(from, to string) int {
	if n, ok := v.LenPathCached(from, to); ok {
		return n
	}

	visited := make(map[string]bool)
	queue := []string{from}
	count := 0

	for len(queue) > 0 {
		var nextQueue []string
		for _, name := range queue {
			visited[name] = true
			if name == to {
				if _, ok := v.LenPathCache[from]; !ok {
					v.LenPathCache[from] = make(map[string]int)
				}
				v.LenPathCache[from][to] = count
				return count
			}
			for _, next := range v.NextFrom(name) {
				if visited[next] {
					continue
				}
				nextQueue = append(nextQueue, next)
			}
		}
		queue = nextQueue
		count++
	}

	if _, ok := v.LenPathCache[from]; !ok {
		v.LenPathCache[from] = make(map[string]int)
	}
	v.LenPathCache[from][to] = -1
	return -1
}

/*

new strategy:

- t0, t1, etc.
- eliminate anything with score zero
- [t0,t1],[t1,t2],[t2,t3], etc.
- eliminate anything with score zero
- if I have nothing left to check, return the best path

*/

func (v *Volcano) Solve2() int {
	var targets []string
	for _, name := range v.Nodes {
		if v.IsClosed(name) && v.Rate(name) > 0 {
			targets = append(targets, name)
		}
	}

	var (
		bestScore int
		bestPath  []string
	)

	var queue [][]string
	for _, t := range targets {
		queue = append(queue, []string{t})
	}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		score := v.Evaluate(path)
		// log.Printf("score %d with path %v", score, path)
		if score > bestScore {
			bestScore = score
			bestPath = path
			log.Printf("⭐️ best score %d with path %v", bestScore, bestPath)
		} else if score == 0 {
			continue
		}

		for _, t := range targets {
			tmp := make([]string, len(path))
			copy(tmp, path)
			tmp = append(tmp, t)

			if len(tmp) > len(targets) {
				// log.Printf("too long: %v", tmp)
				break // exit from target appending
			}
			if !allUnique(tmp) {
				// log.Printf("not unique: %v", tmp)
				continue
			}

			queue = append(queue, tmp)
		}
	}

	return bestScore
}

func isUseless(m map[string]bool, path []string) bool {
	k := pathToKey(path)
	for prefix := range m {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

func pathToKey(s []string) string {
	return strings.Join(s, ",")
}

func allUnique(s []string) bool {
	seen := make(map[string]bool)
	for _, v := range s {
		if seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}
