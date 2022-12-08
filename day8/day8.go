package day8

import (
	"bufio"
	"io"
)

func Part1(f *Forest) (int, error) {
	var result int

	for r := 0; r < f.SideLen(); r++ {
		for c := 0; c < f.SideLen(); c++ {
			if f.IsVisible(r, c) {
				result++
			}
		}
	}

	return result, nil
}

type Forest struct {
	Grid [][]int
}

func NewForest(r io.Reader) (*Forest, error) {
	var grid [][]int

	s := bufio.NewScanner(r)
	for s.Scan() {
		var row []int
		for _, rn := range s.Text() {
			row = append(row, int(rn-'0'))
		}
		grid = append(grid, row)
	}

	return &Forest{Grid: grid}, nil
}

func (f Forest) SideLen() int {
	return len(f.Grid)
}

func (f Forest) IsEdge(r, c int) bool {
	return r == 0 || c == 0 || r == f.SideLen()-1 || c == f.SideLen()-1
}

// IsVisible returns true if the tree at (r, c) is visible from outside the forest.
func (f Forest) IsVisible(r, c int) bool {
	if f.IsEdge(r, c) {
		return true // Trees on the edge of the grid are always visible.
	}

	deltas := []struct{ r, c int }{
		{0, -1}, // left
		{0, 1},  // right
		{-1, 0}, // up
		{1, 0},  // down
	}

	for _, d := range deltas {
		if f.IsVisibleInDirection(r, c, d.r, d.c) {
			return true
		}
	}

	return false
}

// IsVisibleInDirection returns true if the tree at (r, c) has line-of-sight
// to the edge of the forest in the direction (dr, dc).
func (f Forest) IsVisibleInDirection(r, c, dr, dc int) bool {
	height := f.Grid[r][c]

	for {
		if f.IsEdge(r, c) {
			return true
		}

		r += dr
		c += dc

		target := f.Grid[r][c]
		if target >= height {
			return false
		}
	}
}

// ScenicScore is the product of the viewing distances in each cardinal direction.
func (f Forest) ScenicScore(r, c int) int {
	product := 1
	for _, d := range []struct{ r, c int }{
		{0, -1}, // left
		{0, 1},  // right
		{-1, 0}, // up
		{1, 0},  // down
	} {
		product *= f.ViewingDistance(r, c, d.r, d.c)
	}
	return product
}

// ViewingDistance returns the number of trees that can be seen from (r, c) in
// the direction (dr, dc). Disregard trees higher than the starting tree, since
// we can't look to the sky.
func (f Forest) ViewingDistance(r, c, dr, dc int) int {
	height := f.Grid[r][c]
	distance := 0

	for {
		if f.IsEdge(r, c) {
			return distance
		}

		r += dr
		c += dc
		distance++

		target := f.Grid[r][c]
		if target >= height {
			return distance
		}
	}
}

func Part2(f *Forest) int {
	var best int
	for r := 0; r < f.SideLen(); r++ {
		for c := 0; c < f.SideLen(); c++ {
			if score := f.ScenicScore(r, c); score > best {
				best = score
			}
		}
	}
	return best
}
