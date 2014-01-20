// Package util provides various funcionality for strings.
package util

// Function checks if a word 'word 'exist in an array of strings, if not
// then it is added. Returns an array of strings containing 'word' and
// all of the old values
func AddWord(s []string, word string) []string {
	for i := range s {
		if s[i] == word {
			return s
		}
	}
	newS := make([]string, cap(s)+1)
	copy(newS, s)
	newS[len(newS)-1] = word
	return newS
}

// Function cutWord for a given 'word' performs a cut, so that the new
// word (returned) starts at 'begin' position of the old word, and ends
// at 'end' position of the old word.
func CutWord(begin, end int, word string) string {
	if end >= len(word) {
		return ""
	}
	d := make([]uint8, end-begin+1)
	for j, i := 0, begin; i <= end; i, j = i+1, j+1 {
		d[j] = word[i]
	}
	return string(d)
}
