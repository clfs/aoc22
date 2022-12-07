package day7

import (
	"bufio"
	"io"
	"math"
	"strconv"
	"strings"
)

func Parse(r io.Reader) (map[string]int64, error) {
	var dirStack []string

	result := map[string]int64{"/": 0}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "$ cd /":
			dirStack = dirStack[:]
		case strings.HasPrefix(line, "$ ls"):
			// skip
		case strings.HasPrefix(line, "$ cd"):
			dir := line[5:]
			if dir == ".." {
				dirStack = dirStack[:len(dirStack)-1]
			} else {
				dirStack = append(dirStack, dir)
			}
		default:
			fields := strings.Fields(line)
			size, _ := strconv.ParseInt(fields[0], 10, 64)
			name := "/" + strings.Join(append(dirStack, fields[1]), "/")
			result[name] = size
		}
	}
	return result, scanner.Err()
}

func Part1(r io.Reader) (int64, error) {
	filesystem, err := Parse(r)
	if err != nil {
		return 0, err
	}

	dirSizes := make(map[string]int64)

	for name, size := range filesystem {
		if size != 0 {
			continue // skip files
		}
		dirSizes[name] = DirSize(name, filesystem)
	}

	var result int64

	for _, size := range dirSizes {
		if size <= 100000 {
			result += size
		}
	}

	return result, nil
}

func DirSize(name string, filesystem map[string]int64) int64 {
	if !strings.HasSuffix(name, "/") {
		name += "/"
	}

	var total int64
	for nn, ss := range filesystem {
		if strings.HasPrefix(nn, name) {
			total += ss
		}
	}
	return total
}

const (
	DiskSpace   = 70000000
	UpdateSpace = 30000000
)

func Part2(r io.Reader) (int64, error) {
	filesystem, err := Parse(r)
	if err != nil {
		return 0, err
	}

	dirSizes := make(map[string]int64)

	for name, size := range filesystem {
		if size != 0 {
			continue // skip files
		}
		dirSizes[name] = DirSize(name, filesystem)
	}

	var (
		bestSize int64 = math.MaxInt64
	)

	currentFreeSpace := DiskSpace - dirSizes["/"]

	// Find the smallest directory that can be deleted to free up space.
	for _, size := range dirSizes {
		if currentFreeSpace+size >= UpdateSpace {
			if size < bestSize {
				bestSize = size
			}
		}
	}

	return bestSize, nil
}
