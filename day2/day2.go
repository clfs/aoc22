package day2

import (
	"bufio"
	"fmt"
	"io"
)

const (
	OppRock     = 'A'
	OppPaper    = 'B'
	OppScissors = 'C'
	MeRock      = 'X'
	MePaper     = 'Y'
	MeScissors  = 'Z'
)

type Round struct {
	Opponent, Me rune
}

func (r *Round) UnmarshalText(text []byte) error {
	if len(text) < 3 {
		return fmt.Errorf("invalid round: %s", string(text))
	}

	r.Opponent = rune(text[0])
	r.Me = rune(text[2])

	return nil
}

func (r Round) String() string {
	return fmt.Sprintf("%c %c", r.Me, r.Opponent)
}

func (r *Round) Score() int {
	var score int

	switch r.Me {
	case MeRock:
		score += 1
	case MePaper:
		score += 2
	default: // scissors
		score += 3
	}

	switch {
	case r.Me == MeRock && r.Opponent == OppScissors || r.Me == MePaper && r.Opponent == OppRock || r.Me == MeScissors && r.Opponent == OppPaper:
		score += 6 // win
	case r.Me == MeRock && r.Opponent == OppPaper || r.Me == MePaper && r.Opponent == OppScissors || r.Me == MeScissors && r.Opponent == OppRock:
		score += 0 // lose
	default:
		score += 3 // tie
	}

	return score
}

type RoundAlt struct {
	Opponent, Requirement rune
}

const (
	ReqLose = 'X'
	ReqDraw = 'Y'
	ReqWin  = 'Z'
)

func (r *RoundAlt) UnmarshalText(text []byte) error {
	if len(text) < 3 {
		return fmt.Errorf("invalid round: %s", string(text))
	}

	r.Opponent = rune(text[0])
	r.Requirement = rune(text[2])

	return nil
}

func (r RoundAlt) String() string {
	return fmt.Sprintf("%c %c", r.Opponent, r.Requirement)
}

func (r RoundAlt) Score() int {
	var score int

	switch r.Requirement {
	case ReqLose:
		score += 0
	case ReqDraw:
		score += 3
	default: // win
		score += 6
	}

	switch {
	case r.Requirement == ReqWin && r.Opponent == OppScissors || r.Requirement == ReqDraw && r.Opponent == OppRock || r.Requirement == ReqLose && r.Opponent == OppPaper:
		score += 1 // rock
	case r.Requirement == ReqWin && r.Opponent == OppPaper || r.Requirement == ReqDraw && r.Opponent == OppScissors || r.Requirement == ReqLose && r.Opponent == OppRock:
		score += 3 // scissors
	default:
		score += 2 // paper
	}

	return score
}

func parse(r io.Reader) []Round {
	var result []Round

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var rd Round
		if err := rd.UnmarshalText(scanner.Bytes()); err != nil {
			panic(err)
		}
		result = append(result, rd)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

func parseAlt(r io.Reader) []RoundAlt {
	var result []RoundAlt

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var rd RoundAlt
		if err := rd.UnmarshalText(scanner.Bytes()); err != nil {
			panic(err)
		}
		result = append(result, rd)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

func Part1(r io.Reader) int {
	var score int

	for _, rd := range parse(r) {
		score += rd.Score()
	}

	return score
}

func Part2(r io.Reader) int {
	var score int

	for _, rd := range parseAlt(r) {
		score += rd.Score()
	}

	return score
}
