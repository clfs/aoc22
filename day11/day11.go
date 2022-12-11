package day11

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Monkey struct {
	Items      []int
	Operation  func(int) int
	Divisor    int
	Pass, Fail int
}

func ParseOperation(s string) func(int) int {
	if s == "* old" {
		return func(old int) int { return old * old }
	}

	if s[0] == '*' {
		n, err := strconv.Atoi(s[2:])
		if err != nil {
			panic(fmt.Sprintf("bad operation: %q", s))
		}
		return func(old int) int { return old * n }
	}

	if s[0] == '+' {
		n, err := strconv.Atoi(s[2:])
		if err != nil {
			panic(fmt.Sprintf("bad operation: %q", s))
		}
		return func(old int) int { return old + n }
	}

	panic(fmt.Sprintf("bad operation: %q", s))
}

func (m *Monkey) UnmarshalText(text []byte) error {
	lines := strings.Split(string(text), "\n")

	m.Items = ReadNumbers(lines[1])

	m.Operation = ParseOperation(
		// just the "+ 12" bit
		strings.TrimPrefix(lines[2], "  Operation: new = old "))

	m.Divisor = ReadNumbers(lines[3])[0]
	m.Pass = ReadNumbers(lines[4])[0]
	m.Fail = ReadNumbers(lines[5])[0]

	return nil
}

var numbersRe = regexp.MustCompile(`\d+`)

func ReadNumbers(s string) []int {
	var result []int
	for _, num := range numbersRe.FindAllString(s, -1) {
		n, err := strconv.Atoi(num)
		if err != nil {
			panic(fmt.Sprintf("bad number: %q", num))
		}
		result = append(result, n)
	}
	return result
}

func Parse(r io.Reader) ([]Monkey, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	texts := bytes.Split(data, []byte("\n\n"))

	var result []Monkey
	for _, text := range texts {
		var m Monkey
		if err := m.UnmarshalText(text); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}

func Part1(r io.Reader) (int, error) {
	monkeys, err := Parse(r)
	if err != nil {
		return 0, err
	}

	nInspections := make(map[int]int)

	for i := 0; i < 20; i++ {
		log.Printf("==== round %d", i)
		for j, m := range monkeys {
			log.Printf("Monkey %d: %v", j, m.Items)
		}

		for j, m := range monkeys {
			// Monkey inspects each item in its list.
			for k, item := range m.Items {
				// Inspect and apply operation.
				m.Items[k] = m.Operation(item)
				nInspections[j]++

				// Get bored with item.
				m.Items[k] /= 3

				if m.Items[k]%m.Divisor == 0 {
					monkeys[m.Pass].Items = append(monkeys[m.Pass].Items, m.Items[k])
				} else {
					monkeys[m.Fail].Items = append(monkeys[m.Fail].Items, m.Items[k])
				}
			}

			// Monkey is done with its list. (This is the deletion.)
			monkeys[j].Items = []int{}
		}
	}

	log.Print(nInspections)

	values := SortedValues(nInspections)
	return values[len(values)-1] * values[len(values)-2], nil
}

// input:
// [1: 30, 2: 40, 3: 50]
// output:
// [50, 40, 30]
func SortedValues(m map[int]int) []int {
	var result []int
	for _, v := range m {
		result = append(result, v)
	}
	sort.Ints(result)
	return result
}

func Part2(r io.Reader) (int, error) {
	monkeys, err := Parse(r)
	if err != nil {
		return 0, err
	}

	megaMod := 1
	for _, m := range monkeys {
		megaMod *= m.Divisor
	}

	nInspections := make(map[int]int)

	for i := 0; i < 10000; i++ {
		log.Printf("==== round %d", i)
		for j, m := range monkeys {
			log.Printf("Monkey %d: %v", j, m.Items)
		}

		for j := range monkeys {
			// Monkey inspects each item in its list.
			for k := range monkeys[j].Items {
				// Shrink the worry
				monkeys[j].Items[k] %= megaMod

				// Inspect and apply operation.
				monkeys[j].Items[k] = monkeys[j].Operation(monkeys[j].Items[k])
				nInspections[j]++

				if monkeys[j].Items[k]%monkeys[j].Divisor == 0 {
					monkeys[monkeys[j].Pass].Items = append(monkeys[monkeys[j].Pass].Items, monkeys[j].Items[k])
				} else {
					monkeys[monkeys[j].Fail].Items = append(monkeys[monkeys[j].Fail].Items, monkeys[j].Items[k])
				}
			}

			// Monkey is done with its list. (This is the deletion.)
			monkeys[j].Items = []int{}
		}
	}

	log.Print(megaMod)
	log.Print(nInspections)

	values := SortedValues(nInspections)
	return values[len(values)-1] * values[len(values)-2], nil
}

/*

mod 8

n = 1 (mod 8)
n*n = 1 (mod 8)

n = 2 (mod 8)
n*n = 4 (mod 8)

n = 3 (mod 8)
n*n = 9 = 1 (mod 8)

n = 4 (mod 8)
n*n = 16 = 0 (mod 8)


*/
