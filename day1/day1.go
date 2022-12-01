package day1

import (
	"bufio"
	"io"
	"strconv"

	"golang.org/x/exp/slices"
)

func parse(r io.Reader) [][]int {
	var result [][]int
	var group []int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			dst := make([]int, len(group))
			copy(dst, group)
			group = group[:0]

			result = append(result, dst)
			continue
		}

		n, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		group = append(group, n)
	}

	return result
}

func Part1(r io.Reader) int {
	var sum, best int

	for _, group := range parse(r) {
		for _, n := range group {
			sum += n
		}
		if sum > best {
			best = sum
		}
		sum = 0
	}

	return best
}

func Part2(r io.Reader) int {
	var sums []int

	for _, group := range parse(r) {
		var sum int
		for _, n := range group {
			sum += n
		}
		sums = append(sums, sum)
	}

	slices.Sort(sums)

	return sums[len(sums)-1] + sums[len(sums)-2] + sums[len(sums)-3]
}
