package distance

// Levenshtein computes the edit distance between two strings.
func Levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)

	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(
				curr[j-1]+1,
				prev[j]+1,
				prev[j-1]+cost,
			)
		}
		prev, curr = curr, prev
	}

	return prev[len(b)]
}

// ClosestMatch finds the candidate with the smallest edit distance to input,
// returning it only if the distance is at most maxDist.
func ClosestMatch(input string, candidates []string, maxDist int) (string, bool) {
	best := ""
	bestDist := maxDist + 1

	for _, c := range candidates {
		d := Levenshtein(input, c)
		if d < bestDist {
			bestDist = d
			best = c
		}
	}

	if bestDist <= maxDist {
		return best, true
	}
	return "", false
}
