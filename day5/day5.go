package day5

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Move struct {
	Count    int // How many crates to move
	Src, Dst int // The source and destination stacks; one-indexed
}

// "move 10 from 4 to 3"
// should return
// Move{Count: 10, Src: 4, Dst: 3}
var moveRegexp = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)

func (m *Move) UnmarshalText(text []byte) error {
	// use moveRegexp to parse the text
	matches := moveRegexp.FindStringSubmatch(string(text))
	if matches == nil {
		return fmt.Errorf("invalid move: %q", text)
	}

	count, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("invalid count: %q", matches[1])
	}

	src, err := strconv.Atoi(matches[2])
	if err != nil {
		return fmt.Errorf("invalid source: %q", matches[2])
	}

	dst, err := strconv.Atoi(matches[3])
	if err != nil {
		return fmt.Errorf("invalid destination: %q", matches[3])
	}

	m.Count = count
	m.Src = src
	m.Dst = dst

	return nil
}

func Rearrange(crates [][]rune, moves []Move) [][]rune {
	for _, m := range moves {
		for i := 0; i < m.Count; i++ {
			crates[m.Dst-1] = append(crates[m.Dst-1], crates[m.Src-1][len(crates[m.Src-1])-1])
			crates[m.Src-1] = crates[m.Src-1][:len(crates[m.Src-1])-1]
		}
	}
	return crates
}

func RearrangeMultipleAtOnce(crates [][]rune, moves []Move) [][]rune {
	for _, m := range moves {
		crates[m.Dst-1] = append(crates[m.Dst-1], crates[m.Src-1][len(crates[m.Src-1])-m.Count:]...)
		crates[m.Src-1] = crates[m.Src-1][:len(crates[m.Src-1])-m.Count]
	}
	return crates
}

func parse(r io.Reader) ([][]rune, []Move) {
	// First, read the stacks of crates.
	//
	//     [D]
	// [N] [C]
	// [Z] [M] [P]
	//
	// should return [][]rune{{'Z', 'N'}, {'M', 'C', 'D'}, {'P'}}

	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	fields := strings.Split(string(data), "\n\n")

	crates := parseCrates(fields[0])
	moves := parseMoves(fields[1])

	return crates, moves
}

func parseCrates(s string) [][]rune {
	//     [D]
	// [N] [C]
	// [Z] [M] [P]
	//  1   2   3
	//
	// should become:
	//
	// [][]rune{{'Z', 'N'}, {'M', 'C', 'D'}, {'P'}}

	scanner := bufio.NewScanner(strings.NewReader(s))
	crates := make([][]rune, 0)

	// 1 5 9

	for scanner.Scan() {
		line := scanner.Text()

		for i := 1; i < len(line); i += 4 {
			id := line[i]
			if id >= 'A' && id <= 'Z' {
				for len(crates) <= i/4 {
					crates = append(crates, make([]rune, 0))
				}
				crates[i/4] = append(crates[i/4], rune(id))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// reverse the order of the stacks

	for i := 0; i < len(crates); i++ {
		crates[i] = reverse(crates[i])
	}

	return crates
}

func reverse(sl []rune) []rune {
	for i := 0; i < len(sl)/2; i++ {
		j := len(sl) - i - 1
		sl[i], sl[j] = sl[j], sl[i]
	}
	return sl
}

func parseMoves(s string) []Move {
	scanner := bufio.NewScanner(strings.NewReader(s))
	moves := make([]Move, 0)
	for scanner.Scan() {
		var m Move
		if err := m.UnmarshalText(scanner.Bytes()); err != nil {
			panic(err)
		}
		moves = append(moves, m)
	}
	return moves
}

func Part1(r io.Reader) string {
	crates, moves := parse(r)
	crates = Rearrange(crates, moves)

	tops := make([]rune, 0)
	for _, stack := range crates {
		tops = append(tops, stack[len(stack)-1])
	}

	return string(tops)
}

func Part2(r io.Reader) string {
	crates, moves := parse(r)
	crates = RearrangeMultipleAtOnce(crates, moves)

	tops := make([]rune, 0)
	for _, stack := range crates {
		tops = append(tops, stack[len(stack)-1])
	}

	return string(tops)
}
