package day14

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/clfs/aoc22"
)

type Cave struct {
	tiles [][]int
}

const (
	CaveWidth       = 2000
	CaveWidthOffset = 1000
	CaveHeight      = 600
	LeakX           = 500
)

func NewCave(r io.Reader) (*Cave, error) {
	c := &Cave{
		tiles: make([][]int, CaveHeight),
	}
	for i := range c.tiles {
		c.tiles[i] = make([]int, CaveWidth)
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		nums := aoc22.ReadInts(s.Text())
		for i := 0; i < len(nums)-3; i += 2 {
			x0, y0, x1, y1 := nums[i], nums[i+1], nums[i+2], nums[i+3]

			if x0 == x1 {
				if y1 < y0 {
					y0, y1 = y1, y0
				}
				for y := y0; y <= y1; y++ {
					c.Set(NewPoint(x0, y), Rock) // vertical line
				}
			} else {
				if x1 < x0 {
					x0, x1 = x1, x0
				}
				for x := x0; x <= x1; x++ {
					c.Set(NewPoint(x, y0), Rock) // horizontal line
				}
			}
		}
	}

	c.Set(NewPoint(LeakX, 0), Leak)

	return c, s.Err()
}

func (c *Cave) Set(p Point, t int) {
	c.tiles[p.Y][p.X+CaveWidthOffset] = t
}

func (c *Cave) SetOffset(p Point, t int) {
	c.tiles[p.Y][p.X] = t
}

const (
	Air = iota
	Rock
	Sand
	Leak
)

// AddFloor adds a floor two rows below the lowest Rock.
// It returns the y coordinate of the new floor.
func (c *Cave) AddFloor() int {
	var best int
	for i, row := range c.tiles {
		for _, tile := range row {
			if tile == Rock {
				best = i
			}
		}
	}

	for x := 0; x < CaveWidth; x++ {
		c.SetOffset(NewPoint(x, best+2), Rock)
	}

	return best + 2
}

// NumSand returns the number of Sand tiles in the cave.
func (c *Cave) NumSand() int {
	var n int
	for _, row := range c.tiles {
		for _, tile := range row {
			if tile == Sand {
				n++
			}
		}
	}
	return n
}

type Point struct {
	X, Y int
}

func (p Point) Down() Point {
	return Point{p.X, p.Y + 1}
}

func (p Point) DownLeft() Point {
	return Point{p.X - 1, p.Y + 1}
}

func (p Point) DownRight() Point {
	return Point{p.X + 1, p.Y + 1}
}

func NewPoint(x, y int) Point {
	return Point{x, y}
}

// At returns the tile at the given point. If the point is out of bounds, At
// returns Air.
func (c *Cave) At(p Point) int {
	if !c.InBounds(p) {
		return Air
	}
	return c.tiles[p.Y][p.X+CaveWidthOffset]
}

// AtOffset returns the tile at the given point, but does not offset the x
// coordinate.
func (c *Cave) AtOffset(p Point) int {
	if !c.InBoundsOffset(p) {
		return Air
	}
	return c.tiles[p.Y][p.X]
}

func (c *Cave) InBounds(p Point) bool {
	if p.Y < 0 {
		return false
	}

	if p.Y >= CaveHeight {
		return false
	}

	if p.X+CaveWidthOffset < 0 {
		return false
	}

	if p.X+CaveWidthOffset >= CaveWidth {
		return false
	}

	return true
}

func (c *Cave) InBoundsOffset(p Point) bool {
	if p.Y < 0 {
		return false
	}

	if p.Y >= CaveHeight {
		return false
	}

	if p.X < 0 {
		return false
	}

	if p.X >= CaveWidth {
		return false
	}

	return true
}

// Next returns the next point that the sand will flow to. If the sand cannot
// flow any further, the point is unchanged.
func (c *Cave) Next(p Point) Point {
	switch {
	case c.AtOffset(p.Down()) == Air:
		return p.Down()
	case c.AtOffset(p.DownLeft()) == Air:
		return p.DownLeft()
	case c.AtOffset(p.DownRight()) == Air:
		return p.DownRight()
	default:
		return p
	}
}

func (c *Cave) Tick() (ok bool) {
	// log.Print("called tick")

	curr := NewPoint(LeakX+CaveWidthOffset, 0)
	log.Printf("curr: %v", curr)
	log.Printf("curr tile: %v", c.AtOffset(curr))

	for {
		next := c.Next(curr)
		log.Printf("next: %v, next tile: %v", next, c.AtOffset(next))

		// If the next tile is the leak itself, we're plugged up. Return false.
		if tile := c.AtOffset(next); tile == Sand {
			return false
		}

		if next == curr {
			log.Print("next == curr, so setting curr to sand, then return true")
			c.SetOffset(curr, Sand)
			return true
		}

		log.Print("next != curr")

		if !c.InBoundsOffset(next) {
			log.Print("next out of bounds, returning false")
			return false
		}

		log.Print("setting curr to next")
		curr = next
	}
}

// TickUntilStable returns the number of sand tiles after sand starts flowing
// into the abyss below.
func (c *Cave) TickUntilStable() int {
	var n int

	log.Printf("==== %d ====", n)
	log.Print(c.Debug(494+CaveWidthOffset-5, 0, 503+CaveWidthOffset+5, 11))
	for ; c.Tick(); n++ {
		log.Printf("==== %d ====", n)
		log.Print(c.Debug(494+CaveWidthOffset-5, 0, 503+CaveWidthOffset+5, 11))
	}
	return c.NumSand()
}

func Part1(r io.Reader) (int, error) {
	c, err := NewCave(r)
	if err != nil {
		return 0, err
	}
	return c.TickUntilStable(), nil
}

func Part2(r io.Reader) (int, error) {
	c, err := NewCave(r)
	if err != nil {
		return 0, err
	}
	n := c.AddFloor()
	log.Printf("added floor on y=%d", n)
	return c.TickUntilStable(), nil
}

func (c Cave) Debug(x0, y0, x1, y1 int) string {
	var b strings.Builder
	for r, row := range c.tiles[y0 : y1+1] {
		for c, tile := range row[x0 : x1+1] {
			switch tile {
			case Air:
				b.WriteRune('.')
			case Rock:
				b.WriteRune('#')
			case Sand:
				b.WriteRune('o')
			case Leak:
				b.WriteRune('+')
			default:
				panic(fmt.Sprintf("unknown tile %d at row %d, col %d", tile, r, c))
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func init() {
	log.SetFlags(0)
}
