// Package token_util provides supporting functions for tokens.
package token_util

import "regexp"

// MatchToken returns true if 'word' matches given 'token' regex, false
// otherwise.
func MatchToken(tokens map[string]string, token, word string) bool {
	regex := regexp.MustCompile(tokens[token])
	if regex.MatchString(word) {
		return true
	}
	return false
}
