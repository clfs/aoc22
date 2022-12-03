package day3

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type Rucksack struct {
	Left, Right string
}

type Group []Rucksack

func (r *Rucksack) CommonItem() rune {
	var (
		left  = make(map[rune]bool)
		right = make(map[rune]bool)
	)

	for _, rn := range r.Left {
		left[rn] = true
	}
	for _, rn := range r.Right {
		right[rn] = true
	}

	var common rune

	for rn := range left {
		if right[rn] {
			common = rn
			break
		}
	}

	return common
}

func (r *Rucksack) UnmarshalText(text []byte) error {
	if len(text)%2 == 1 {
		return fmt.Errorf("invalid rucksack: %s", text)
	}

	bound := len(text) / 2
	r.Left, r.Right = string(text[:bound]), string(text[bound:])
	return nil
}

func parse(r io.Reader) []Rucksack {
	var result []Rucksack

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var ruck Rucksack
		if err := ruck.UnmarshalText(scanner.Bytes()); err != nil {
			panic(err)
		}
		result = append(result, ruck)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

func Priority(r rune) int {
	// Lowercase item types a through z have priorities 1 through 26.
	// Uppercase item types A through Z have priorities 27 through 52.

	switch {
	case r >= 'a' && r <= 'z':
		return int(r - 'a' + 1)
	case r >= 'A' && r <= 'Z':
		return int(r - 'A' + 27)
	default:
		panic(r)
	}
}

func Part1(r io.Reader) int {
	var sum int
	for _, ruck := range parse(r) {
		sum += Priority(ruck.CommonItem())
	}
	return sum
}

func parseGroups(r io.Reader) [][]Rucksack {
	var (
		result [][]Rucksack
		group  []Rucksack
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var ruck Rucksack
		if err := ruck.UnmarshalText(scanner.Bytes()); err != nil {
			panic(err)
		}

		group = append(group, ruck)

		log.Println(group)
		if len(group) == 3 {
			dst := make([]Rucksack, len(group))
			copy(dst, group)
			result = append(result, dst)
			group = group[:0]
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return result
}

func BadgeFor(a, b, c Rucksack) rune {
	// For each rucksack, disregard the left and right compartments,
	// and find the only item type that appears in all 3 rucksacks.

	var aM = make(map[rune]bool)
	var bM = make(map[rune]bool)
	var cM = make(map[rune]bool)

	for _, rn := range a.Left {
		aM[rn] = true
	}
	for _, rn := range a.Right {
		aM[rn] = true
	}
	for _, rn := range b.Left {
		bM[rn] = true
	}
	for _, rn := range b.Right {
		bM[rn] = true
	}
	for _, rn := range c.Left {
		cM[rn] = true
	}
	for _, rn := range c.Right {
		cM[rn] = true
	}

	var common rune

	for rn := range aM {
		if bM[rn] && cM[rn] {
			common = rn
			break
		}
	}

	return common
}

func Part2(r io.Reader) int {
	var sum int
	for _, group := range parseGroups(r) {
		sum += Priority(BadgeFor(group[0], group[1], group[2]))
	}
	return sum
}
