package day6

func Part1(s string) int {
	for i := 0; i < len(s)-4; i++ {
		window := s[i : i+4]
		if isSOP(window) {
			return i + 4 // the answer's one-indexed
		}
	}
	return -1
}

func Part2(s string) int {
	for i := 0; i < len(s)-14; i++ {
		window := s[i : i+14]
		if isSOP(window) {
			return i + 14 // the answer's one-indexed
		}
	}
	return -1
}

func isSOP(s string) bool {
	seen := make(map[rune]bool)
	for _, r := range s {
		if seen[r] {
			return false
		}
		seen[r] = true
	}
	return true
}
