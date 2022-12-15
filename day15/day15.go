package day15

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"golang.org/x/exp/slices"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Point struct {
	X, Y int
}

type Sensor struct {
	Location, NearestBeacon Point
}

func (s *Sensor) UnmarshalText(text []byte) error {
	_, err := fmt.Sscanf(
		string(text),
		"Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
		&s.Location.X,
		&s.Location.Y,
		&s.NearestBeacon.X,
		&s.NearestBeacon.Y,
	)
	return err
}

func Parse(r io.Reader) ([]Sensor, error) {
	var sensors []Sensor

	s := bufio.NewScanner(r)
	for s.Scan() {
		var sensor Sensor
		if err := sensor.UnmarshalText(s.Bytes()); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}

	return sensors, s.Err()
}

// Radius returns the distance from the sensor to the nearest beacon.
func (s *Sensor) Radius() int {
	return abs(s.Location.X-s.NearestBeacon.X) + abs(s.Location.Y-s.NearestBeacon.Y)
}

// Intersect returns the range of x values in the intersection of the sensor's
// area and a horizontal line at y.
func (s *Sensor) Intersect(y int) (Range, bool) {
	delta := s.Radius() - abs(s.Location.Y-y)
	if delta < 0 {
		return Range{}, false
	}
	return Range{s.Location.X - delta, s.Location.X + delta}, true
}

// IntersectSegment returns the range of x values in the intersection of the
// sensor's area and the segment from (0, y) to (x, y).
func (s *Sensor) IntersectSegment(x, y int) (Range, bool) {
	delta := s.Radius() - abs(s.Location.Y-y)
	if delta < 0 {
		return Range{}, false
	}
	return Range{
		max(0, s.Location.X-delta),
		min(x, s.Location.X+delta),
	}, true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Range is an inclusive range.
type Range struct {
	Low, High int
}

// Len returns the number of integers in the inclusive range.
func (r Range) Len() int {
	return r.High - r.Low + 1
}

func LenUnion(rs []Range) int {
	var res int
	for _, r := range rs {
		res += r.Len()
	}
	return res
}

func In(rs []Range, x int) bool {
	for _, r := range rs {
		if r.Low <= x && x <= r.High {
			return true
		}
	}
	return false
}

// Union returns the non-overlapping, sorted union of the given ranges.
// If rs is empty, Union returns nil.
// If rs is not sorted, Union modifies the order of rs.
func Union(rs []Range) []Range {
	if len(rs) == 0 {
		return nil
	}

	// First, sort the ranges by low.
	slices.SortFunc(rs, func(a, b Range) bool { return a.Low < b.Low })

	// Then, merge them.
	union := []Range{rs[0]}
	for _, r := range rs[1:] {
		last := &union[len(union)-1]
		if r.Low <= last.High+1 {
			last.High = max(r.High, last.High)
		} else {
			union = append(union, r)
		}
	}

	return union
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r Range) String() string {
	return fmt.Sprintf("[%d,%d]", r.Low, r.High)
}

// NImpossible accepts a horizontal line at y, and returns the number of points
// on the line that cannot contain a beacon.
func NImpossible(sensors []Sensor, y int) int {
	var ranges []Range

	for _, s := range sensors {
		if r, ok := s.Intersect(y); ok {
			ranges = append(ranges, r)
		}
	}

	union := Union(ranges)

	takenX := make(map[int]bool)
	for _, s := range sensors {
		if s.NearestBeacon.Y == y {
			takenX[s.NearestBeacon.X] = true
		}
	}

	return LenUnion(union) - len(takenX)
}

// NImpossibleSegment accepts a horizontal line segment between (0,y) and (x,y)
// inclusive, and returns the number of points on the segment that cannot
// contain a beacon.
func NImpossibleSegment(sensors []Sensor, x, y int) int {
	var ranges []Range

	for _, s := range sensors {
		if r, ok := s.IntersectSegment(x, y); ok {
			ranges = append(ranges, r)
		}
	}

	union := Union(ranges)

	takenX := make(map[int]bool)
	for _, s := range sensors {
		if s.NearestBeacon.Y == y && s.NearestBeacon.X <= x {
			takenX[s.NearestBeacon.X] = true
		}
	}

	return LenUnion(union) - len(takenX)
}

func Part1(r io.Reader, y int) (int, error) {
	sensors, err := Parse(r)
	if err != nil {
		return 0, err
	}

	return NImpossible(sensors, y), nil
}

func Part2(r io.Reader, bound int) (int, error) {
	sensors, err := Parse(r)
	if err != nil {
		return 0, err
	}

	p, err := FindDistressBeacon(sensors, bound)
	if err != nil {
		return 0, err
	}

	return p.TuningFrequency(), nil
}

func (p Point) TuningFrequency() int {
	return p.X*4000000 + p.Y
}

// FindDistressBeacon returns the point of the distress beacon.
// The beacon is within (0,0)x(bound,bound) inclusive.
func FindDistressBeacon(sensors []Sensor, bound int) (Point, error) {
	// Backwards, since Eric probably placed it at the bottom
	for y := bound; y >= 0; y-- {

		nImpossible := NImpossibleSegment(sensors, bound, y)

		if nImpossible > bound {
			continue
		}

		for x := 0; x <= bound; x++ {
			if IsDistressBeacon(sensors, x, y) {
				log.Printf("distress beacon found at (%d,%d)", x, y)
				return Point{x, y}, nil
			}
		}
	}

	return Point{}, fmt.Errorf("no distress beacon found")
}

func IsDistressBeacon(sensors []Sensor, x, y int) bool {
	p := Point{x, y}
	for _, s := range sensors {
		if Distance(s.Location, p) <= s.Radius() {
			return false
		}
	}
	return true
}

func Distance(p1, p2 Point) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}
