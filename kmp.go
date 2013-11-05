package main
import (
	"fmt" //implements fomratted I/O.
	"log" //simple logging package
	"os" //accessing command-line arguments
)

/* 	Implementation of Knuth-Morris-Pratt algorithm (Prefix based aproach).
	Requires two command line arguments.
	
	@argument string to be searched "for" (pattern, search word), no spaces allowed
	@argument one space
	@argument string to be searched "in" (text), single spaces allowed
*/
func main() {
	// Error handling & declaration
	args := os.Args
	if (len(args) <= 2) {
		log.Fatal("Not enough arguments. Two string arguments separated by spaces are required!")
	}
	pattern := args[1]
	s := args[2]
	for i := 3; i<len(args); i++ {
		s = s +" "+ args[i]
	}
	if ( len(args[1]) > len(s) ) {
		log.Fatal("Pattern  is longer than text!")
	}
	// Alghoritm execution
	fmt.Printf("\nRunning: Knuth-Morris-Pratt alghoritm.\n\n")
	fmt.Printf("Search word (%d chars long): %q.\n",len(args[1]), pattern)
	fmt.Printf("Text        (%d chars long): %q.\n\n",len(s), s)
	knp(s, pattern)
}

/*  Function knp performing the Knuth-Morris-Pratt alghoritm.
    Prints whether the word/pattern was found and on what position in the text or not.
	
	@param text string/text to be searched in
	@param word word/pattern to be serached for
*/  
func knp(text, word string) {
	m, i, c := 0, 0, 0 //m - current match in text, i - current character in w, c - ammount of comparations
	t := kmp_table(word)
	for  m + i < len(text) {
		fmt.Printf("\n   comparing characters %c %c",word[i],text[m+i])
		c++
		if (word[i] == text[m+i]) {
			fmt.Printf(" - match")
			if (i == len(word) - 1) {
				fmt.Printf("\n\nWord %q was found at position %d in %q. \n%d comparisons were done.",word, m, text,c)
				return
			}
			i++
		} else {
			m = m + i - t[i]
			if (t[i] > -1) {
				i = t[i]
			} else {
				i = 0
			} 
		}
	}
	fmt.Printf("\n\nWord was not found.\n%d comparisons were done.",c)
	return
}

/*
	Table building alghoritm.
	
	@param word word to be analyzed
	@param t table to be filled
*/
func kmp_table(word string)(t []int) {
	t = make([]int, len(word))
    pos, cnd := 2, 0
	t[0], t[1] = -1, 0
	for pos < len(word) {
		if (word[pos-1] == word[cnd]) {
			cnd++
			t[pos] = cnd
			pos++
		} else if (cnd > 0) {
			cnd = t[cnd]
		} else {
			t[pos] = 0
			pos++
		}
	}
    return t
}