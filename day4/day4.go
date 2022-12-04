package day4

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type Pair struct {
	LeftLow   int
	LeftHigh  int
	RightLow  int
	RightHigh int
}

var pairRegexp = regexp.MustCompile(`^(\d+)-(\d+),(\d+)-(\d+)$`)

func (p *Pair) UnmarshalText(text []byte) error {
	matches := pairRegexp.FindStringSubmatch(string(text))
	if matches == nil {
		return fmt.Errorf("invalid pair: %s", text)
	}
	p.LeftLow, _ = strconv.Atoi(matches[1])
	p.LeftHigh, _ = strconv.Atoi(matches[2])
	p.RightLow, _ = strconv.Atoi(matches[3])
	p.RightHigh, _ = strconv.Atoi(matches[4])
	return nil
}

// Redundant returns true if one of the ranges is contained in the other.
func (p *Pair) Redundant() bool {
	return (p.LeftLow >= p.RightLow && p.LeftHigh <= p.RightHigh) ||
		(p.RightLow >= p.LeftLow && p.RightHigh <= p.LeftHigh)
}

// AnyOverlap returns true if the left and right ranges overlap at all.
func (p *Pair) AnyOverlap() bool {
	return (p.LeftLow <= p.RightLow && p.LeftHigh >= p.RightLow) ||
		(p.RightLow <= p.LeftLow && p.RightHigh >= p.LeftLow)
}

func Part1(r io.Reader) int {
	var count int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var p Pair
		if err := p.UnmarshalText(scanner.Bytes()); err != nil {
			panic(err)
		}
		if p.Redundant() {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return count
}

func Part2(r io.Reader) int {
	var count int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var p Pair
		if err := p.UnmarshalText(scanner.Bytes()); err != nil {
			panic(err)
		}
		if p.AnyOverlap() {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return count
}
