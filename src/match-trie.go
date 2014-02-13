// match-trie.go contains funcionality for prefix tree (trie).
package main

import (
	"log"
	"strings"
)

// createNewTrie initializes a new prefix tree. State is the number of
// first state to be created, i is the number of first pattern to be
// added.
func createNewTrie() (trie map[int]map[string]int, finalFor []int, state int, i int) {
	return make(map[int]map[string]int), make([]int, 1), 1, 1
}

// appendPattern creates all the necessary transitions for a single
// pattern to the given trie.
func appendPattern(tokens map[string]string, pattern string, trie map[int]map[string]int, finalFor []int, state int, i int) (map[int]map[string]int, []int, int, int) {
	patternsNameSplit := separatePatternFromName(pattern)
	words := strings.Split(patternsNameSplit[1], " ")
	current := 0
	j := 0
	for j < len(words) && getTransition(current, words[j], trie) != -1 {
		current = getTransition(current, words[j], trie)
		j++
	}
	for j < len(words) {
		finalFor = intArraySizeUp(finalFor, 1)
		finalFor[state] = 0
		if len(getTransitionWords(current, trie)) > 0 && words[j][0] == '<' && words[j][len(words[j])-1] == '>' {
			// conflict check when adding regex transition
			transitionWords := getTransitionWords(current, trie)
			for w := range transitionWords {
				tokenWithoutBrackets := cutWord(1, len(words[j])-2, words[j])
				tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
				switch len(tokenWithoutBracketsSplit) {
				case 2:
					{
						if matchToken(tokens, tokenWithoutBracketsSplit[0], transitionWords[w]) {
							log.Fatal("pattern conflict: token \"" + words[j] + "\" matches word \"" + transitionWords[w] + "\"")
						}
					}
				case 1:
					{
						if matchToken(tokens, tokenWithoutBrackets, transitionWords[w]) {
							log.Fatal("pattern conflict: token \"" + words[j] + "\" matches word \"" + transitionWords[w] + "\"")
						}
					}
				default:
					log.Fatal("invalid token definition: \"<" + tokenWithoutBrackets + ">\"")
				}
			}
		} else if len(getTransitionTokens(current, trie)) > 0 && words[j][0] != '<' && words[j][len(words[j])-1] != '>' {
			// conflict check when adding word transition
			transitionTokens := getTransitionTokens(current, trie)
			for t := range transitionTokens {
				tokenWithoutBrackets := cutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
				tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
				switch len(tokenWithoutBracketsSplit) {
				case 2:
					{
						if matchToken(tokens, tokenWithoutBracketsSplit[0], words[j]) {
							log.Fatal("pattern conflict: token \"" + transitionTokens[t] + "\" matches word \"" + words[j] + "\"")
						}
					}
				case 1:
					{
						if matchToken(tokens, tokenWithoutBrackets, words[j]) {
							log.Fatal("pattern conflict: token \"" + transitionTokens[t] + "\" matches word \"" + words[j] + "\"")
						}
					}
				default:
					log.Fatal("invalid token definition: \"<" + tokenWithoutBrackets + ">\"")
				}
			}
		}
		createTransition(current, words[j], state, trie)
		current = state
		j++
		state++
	}
	if finalFor[current] != 0 {
		log.Fatal("duplicate pattern detected: \"", pattern, "\"")
	} else {
		// mark current state as terminal for pattern number i
		finalFor[current] = i
	}
	i++ // increment pattern number
	return trie, finalFor, state, i
}
