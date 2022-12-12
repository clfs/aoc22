package day12

import (
	"bufio"
	"fmt"
	"io"
)

func ToHeight(r rune) int {
	switch {
	case 'a' <= r && r <= 'z':
		return int(r - 'a')
	case r == 'S':
		return ToHeight('a')
	case r == 'E':
		return ToHeight('z')
	default:
		panic(fmt.Sprintf("bad rune: %c", r))
	}
}

type Point struct {
	X, Y int
}

type Topo struct {
	Grid        [][]int
	Start, Goal Point
}

func Parse(r io.Reader) (Topo, error) {
	var topo Topo

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		row := make([]int, len(line))
		for i, r := range line {
			row[i] = ToHeight(r)
			if r == 'S' {
				topo.Start = Point{X: i, Y: len(topo.Grid)}
			}
			if r == 'E' {
				topo.Goal = Point{X: i, Y: len(topo.Grid)}
			}
		}
		topo.Grid = append(topo.Grid, row)
	}

	return topo, s.Err()
}

func (t *Topo) At(p Point) int {
	return t.Grid[p.Y][p.X]
}

func (t *Topo) Neighbors(p Point) []Point {
	var neighbors []Point
	for _, d := range []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		np := Point{X: p.X + d.X, Y: p.Y + d.Y}
		if 0 <= np.X && np.X < len(t.Grid[0]) &&
			0 <= np.Y && np.Y < len(t.Grid) {
			neighbors = append(neighbors, np)
		}
	}
	return neighbors
}

func (t *Topo) CanMove(from, to Point) bool {
	var found bool
	neighbors := t.Neighbors(from)
	for _, n := range neighbors {
		if n == to {
			found = true
			break
		}
	}
	return found && t.At(to)-t.At(from) <= 1
}

func (t *Topo) LengthOfShortestPath() (int, error) {
	seen := make(map[Point]int)
	queue := []struct {
		p Point
		n int
	}{{t.Start, 0}}

	for len(queue) > 0 {
		// Get the first item from the queue.
		head := queue[0]
		queue = queue[1:]

		// If the head is the goal, we're done.
		if head.p == t.Goal {
			return head.n, nil
		}

		// If we've seen this point before, and the number of steps is greater
		// than or equal to the number of steps we've seen before, then we
		// don't need to explore this path.
		if n, ok := seen[head.p]; ok && head.n >= n {
			continue
		}

		// Mark this point as seen.
		seen[head.p] = head.n

		// Add the neighbors to the queue.
		for _, np := range t.Neighbors(head.p) {
			if t.CanMove(head.p, np) {
				queue = append(queue, struct {
					p Point
					n int
				}{np, head.n + 1})
			}
		}
	}

	return 0, fmt.Errorf("no path exists between %v and %v", t.Start, t.Goal)
}

func Part1(r io.Reader) (int, error) {
	topo, err := Parse(r)
	if err != nil {
		return 0, err
	}
	return topo.LengthOfShortestPath()
}

func Part2(r io.Reader) (int, error) {
	topo, err := Parse(r)
	if err != nil {
		return 0, err
	}

	var best int

	for x := 0; x < len(topo.Grid[0]); x++ {
		for y := 0; y < len(topo.Grid); y++ {
			p := Point{X: x, Y: y}
			if topo.At(p) == ToHeight('a') {
				topo.Start = p
				n, err := topo.LengthOfShortestPath()
				if err != nil {
					continue // no path exists
				}
				if best == 0 || n < best {
					best = n
				}
			}
		}
	}

	return best, nil
}
