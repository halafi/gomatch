package main

import "regexp"

// MatchToken returns true if 'word' matches given 'token' regex, false
// otherwise.
func matchToken(tokens map[string]string, token, word string) bool {
	regex := regexp.MustCompile(tokens[token])
	if regex.MatchString(word) {
		return true
	}
	return false
}

// Function checks if a word 'word 'exist in an array of strings, if not
// then it is added. Returns an array of strings containing 'word' and
// all of the old values
func addWord(s []string, word string) []string {
	if !contains(s, word) {
		newS := make([]string, cap(s)+1)
		copy(newS, s)
		newS[len(newS)-1] = word
		return newS
	} else {
		return s
	}
}

// Function cutWord for a given 'word' performs a cut, so that the new
// word (returned) starts at 'begin' position of the old word, and ends
// at 'end' position of the old word.
func cutWord(begin, end int, word string) string {
	if end >= len(word) {
		return ""
	}
	d := make([]uint8, end-begin+1)
	for j, i := 0, begin; i <= end; i, j = i+1, j+1 {
		d[j] = word[i]
	}
	return string(d)
}

// Contains checks if an array of strings 's' contains 'word', if yes
// returns true, false otherwise.
func contains(s []string, word string) bool {
	for i := range s {
		if s[i] == word {
			return true
		}
	}
	return false
}

// Increases size of string array by the ammnout given 'c'.
func stringArraySizeUp(array []string, c int) []string {
	newA := make([]string, cap(array)+c)
	copy(newA, array)
	return newA
}

// Increases size of int array by the ammnout given 'c'.
func intArraySizeUp(array []int, c int) []int {
	newA := make([]int, cap(array)+c)
	copy(newA, array)
	return newA
}
