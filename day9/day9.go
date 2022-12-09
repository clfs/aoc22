package day9

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
)

type Vec2 struct {
	X, Y int
}

func (v Vec2) Add(w Vec2) Vec2 {
	return Vec2{v.X + w.X, v.Y + w.Y}
}

func (v Vec2) Sub(w Vec2) Vec2 {
	return Vec2{v.X - w.X, v.Y - w.Y}
}

func (v Vec2) Dist(w Vec2) float64 {
	return math.Sqrt(float64((v.X-w.X)*(v.X-w.X) + (v.Y-w.Y)*(v.Y-w.Y)))
}

func (v Vec2) Mag() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

type Rope struct {
	knots []Vec2
	seen  map[Vec2]bool
}

func NewRope(n int) (*Rope, error) {
	if n < 2 {
		return nil, fmt.Errorf("invalid rope length: %d", n)
	}
	return &Rope{
		knots: make([]Vec2, n),
		seen:  make(map[Vec2]bool),
	}, nil
}

func (r *Rope) Tail() Vec2 {
	return r.knots[len(r.knots)-1]
}

func (r *Rope) TailsSeen() []Vec2 {
	seen := make([]Vec2, 0, len(r.seen))
	for v := range r.seen {
		seen = append(seen, v)
	}
	return seen
}

func (r *Rope) update(dir Vec2) {
	r.knots[0] = r.knots[0].Add(dir)

	for i := 1; i < len(r.knots); i++ {
		gap := r.knots[i-1].Sub(r.knots[i])

		// easy case, we're close enough to the next knot
		if gap.Mag() < 2 {
			continue
		}

		// otherwise, we need to move the knot
		switch gap {
		case Vec2{2, 0}:
			dir = Vec2{1, 0}
		case Vec2{-2, 0}:
			dir = Vec2{-1, 0}
		case Vec2{0, 2}:
			dir = Vec2{0, 1}
		case Vec2{0, -2}:
			dir = Vec2{0, -1}
		case Vec2{2, 1}, Vec2{1, 2}, Vec2{2, 2}:
			dir = Vec2{1, 1}
		case Vec2{2, -1}, Vec2{1, -2}, Vec2{2, -2}:
			dir = Vec2{1, -1}
		case Vec2{-2, 1}, Vec2{-1, 2}, Vec2{-2, 2}:
			dir = Vec2{-1, 1}
		case Vec2{-2, -1}, Vec2{-1, -2}, Vec2{-2, -2}:
			dir = Vec2{-1, -1}
		default:
			panic(fmt.Sprintf("%v: bad gap %v", r.knots, gap))
		}

		r.knots[i] = r.knots[i].Add(dir)
	}

	r.seen[r.Tail()] = true

	log.Printf("%v\n", r.Debug(0, 0, 5, 5))
}

func (r *Rope) Debug(x0, y0, x1, y1 int) string {
	// make an (x1-x0+1) x (y1-y0+1) grid
	grid := make([][]rune, x1-x0+1)
	for i := range grid {
		grid[i] = make([]rune, y1-y0+1)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	// mark the starting point
	grid[0][0] = 's'

	// mark the knots with their index, iterating backwards
	for i := len(r.knots) - 1; i >= 0; i-- {
		x, y := r.knots[i].X, r.knots[i].Y
		if x < x0 || x > x1 || y < y0 || y > y1 {
			continue
		}
		switch i {
		case 0:
			grid[y-y0][x-x0] = 'H'
		default:
			grid[y-y0][x-x0] = rune('0' + i)
		}
	}

	// turn the grid into a string
	var b strings.Builder
	for i := len(grid) - 1; i >= 0; i-- {
		for j := range grid[i] {
			b.WriteRune(grid[i][j])
		}
		b.WriteRune('\n')
	}
	return b.String()
}

type Instruction struct {
	Direction string
	Count     int
}

func (i *Instruction) UnmarshalText(text []byte) error {
	fields := strings.Fields(string(text))
	if len(fields) != 2 {
		return fmt.Errorf("invalid instruction: %s", text)
	}

	var err error
	i.Direction = fields[0]
	i.Count, err = strconv.Atoi(fields[1])
	return err
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s %d", i.Direction, i.Count)
}

func (r *Rope) Follow(ins Instruction) error {
	var v Vec2
	switch ins.Direction {
	case "U":
		v = Vec2{Y: 1}
	case "L":
		v = Vec2{X: -1}
	case "D":
		v = Vec2{Y: -1}
	case "R":
		v = Vec2{X: 1}
	default:
		return fmt.Errorf("bad direction: %s", ins.Direction)
	}

	log.Printf("== %s ==\n\n", ins)
	log.Printf("%v\n", r.Debug(0, 0, 5, 5))

	for i := 0; i < ins.Count; i++ {
		r.update(v)
	}

	return nil
}

func Part1(r io.Reader) (int, error) {
	return solve(r, 2)
}

func Part2(r io.Reader) (int, error) {
	return solve(r, 10)
}

func solve(r io.Reader, n int) (int, error) {
	log.SetOutput(io.Discard) // disable logging

	rope, err := NewRope(n)
	if err != nil {
		return 0, err
	}

	var ins Instruction

	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := ins.UnmarshalText(s.Bytes()); err != nil {
			return 0, err
		}
		if err := rope.Follow(ins); err != nil {
			return 0, err
		}
	}

	return len(rope.TailsSeen()), s.Err()
}
