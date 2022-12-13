package day13

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"sort"
)

func PacketToString(p []any) string {
	b, err := json.Marshal(p)
	if err != nil {
		panic(fmt.Sprintf("un-stringable packet %#v: %v", p, err))
	}
	return string(b)
}

func ParsePacket(b []byte) ([]any, error) {
	// Check for valid characters only. This is for easier fuzzing.
	for _, c := range b {
		switch c {
		case '[', ']':
			continue
		case ',', ' ':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		default:
			return nil, fmt.Errorf("invalid character %q", c)
		}
	}

	var p []any
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("unparseable %q: %v", b, err)
	}
	return p, nil
}

// Parse parses a list of packets from r, skipping empty lines.
func Parse(r io.Reader) ([][]any, error) {
	var result [][]any

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 {
			continue
		}
		packet, err := ParsePacket(line)
		if err != nil {
			return nil, err
		}
		result = append(result, packet)
	}

	return result, s.Err()
}

var (
	typeFloat64 = reflect.TypeOf(float64(0))
	typePacket  = reflect.TypeOf([]any{})
)

// CompareFloat returns -1, 0, or 1 if a < b, a == b, or a > b.
func CompareFloat(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// Compare returns -1, 0, or 1 if a < b, a == b, or a > b.
func Compare(a, b []any) int {
	log.Printf("cmp(%v, %v) = ?", a, b)

	if len(a) == 0 && len(b) == 0 {
		log.Printf("cmp(%v, %v) = 0, since both are empty", a, b)
		return 0
	}

	if len(a) != 0 && len(b) == 0 {
		log.Printf("cmp(%v, %v) = 1, since only b is empty", a, b)
		return 1
	}

	if len(a) == 0 && len(b) != 0 {
		log.Printf("cmp(%v, %v) = -1, since only a is empty", a, b)
		return -1
	}

	var n int

	for i := range a {
		if i >= len(b) {
			log.Printf("cmp(%v, %v) = 1, since b is out of items", a, b)
			return 1
		}

		var (
			ai  = a[i]
			bi  = b[i]
			aiT = reflect.TypeOf(ai)
			biT = reflect.TypeOf(bi)
		)

		log.Printf("inspect elements %v and %v", ai, bi)

		switch {
		case aiT == typeFloat64 && biT == typeFloat64:
			n = CompareFloat(ai.(float64), bi.(float64))
		case aiT == typePacket && biT == typePacket:
			n = Compare(ai.([]any), bi.([]any))
		case aiT == typePacket && biT == typeFloat64:
			n = Compare(ai.([]any), []any{bi})
		case aiT == typeFloat64 && biT == typePacket:
			n = Compare([]any{ai}, bi.([]any))
		}

		if n != 0 {
			log.Printf("cmp(%v, %v) = %d, since elements were unequal", a, b, n)
			return n
		}
	}

	if len(a) < len(b) {
		log.Printf("cmp(%v, %v) = -1, since a is out of items", a, b)
		return -1
	}

	log.Printf("cmp(%v, %v) = %d, after inspecting all elements", a, b, n)
	return n
}

func Part1(r io.Reader) (int, error) {
	packets, err := Parse(r)
	if err != nil {
		return 0, err
	}

	var good []int
	for i := 0; i < len(packets); i += 2 {
		if Compare(packets[i], packets[i+1]) == -1 {
			good = append(good, i/2+1)
		}
	}

	log.Print(good)

	var sum int
	for _, i := range good {
		sum += i
	}
	return sum, nil
}

func Part2(r io.Reader) (int, error) {
	packets, err := Parse(r)
	if err != nil {
		return 0, err
	}

	dividers := []string{"[[2]]", "[[6]]"}
	for _, d := range dividers {
		p, err := ParsePacket([]byte(d))
		if err != nil {
			return 0, err
		}
		packets = append(packets, p)
	}

	// Sort packets.
	sort.Slice(packets, func(i, j int) bool {
		return Compare(packets[i], packets[j]) == -1
	})

	product := 1
	for i := 0; i < len(packets); i++ {
		s := PacketToString(packets[i])
		if s == "[[2]]" || s == "[[6]]" {
			product *= (i + 1)
		}
	}

	return product, nil
}
