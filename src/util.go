// util.go provides some utility functions (strings, slices, regexp).
package main

import (
	"regexp"
	"strings"
	"log"
)

// logLineSplit splits a single log line string into words, words can
// only be separated by any ammount of spaces.
func logLineSplit(line string) []string {
	words := make([]string, 0)
	if line == "" {
		return words
	}
	words = stringArraySizeUp(words, 1)
	wordIndex := 0
	chars := []uint8(line)
	for c := range chars {
		if chars[c] == ' ' && c < len(chars)-1 {
			if words[wordIndex] != "" {
				words = stringArraySizeUp(words, 1)
				wordIndex++
			}
		} else if chars[c] != ' ' {
			words[wordIndex] = words[wordIndex] + string(chars[c])
		}
	}
	return words
}

// matchToken returns true if a word matches token, false otherwise.
func matchToken(tokens map[string]*regexp.Regexp, token, word string) bool {
	if tokens[token] == nil {
		log.Fatal("token: ", token, " undefined")
	}
	return tokens[token].MatchString(word)
}

// cutWord for a given word performs a cut (both prefix and sufix).
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

// contains checks if an array of strings contains given word.
func contains(s []string, word string) bool {
	for i := range s {
		if s[i] == word {
			return true
		}
	}
	return false
}

// lineSplit parses a mutli-line string into single lines (array of
// strings).
func lineSplit(input string) []string {
	inputSplit := make([]string, 1)
	inputSplit[0] = input
	if strings.Contains(input, "\r\n") { //CR+LF
		inputSplit = strings.Split(input, "\r\n")
	} else if strings.Contains(input, "\n") { //LF
		inputSplit = strings.Split(input, "\n")
	} else if strings.Contains(input, "\r") { //CR
		inputSplit = strings.Split(input, "\r")
	}
	return inputSplit
}

// stringArraySizeUp creates a new string array with old values and
// increased maximum size by the ammnout given.
func stringArraySizeUp(array []string, c int) []string {
	newA := make([]string, cap(array)+c)
	copy(newA, array)
	return newA
}

// intArraySizeUp creates a new int array with old values and increased
// maximum size by the ammnout given.
func intArraySizeUp(array []int, c int) []int {
	newA := make([]int, cap(array)+c)
	copy(newA, array)
	return newA
}
