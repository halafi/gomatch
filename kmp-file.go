package main //package main required for standalone executable
import (
	"fmt" //implements fomratted I/O.
	"log" //simple logging package
	"io/ioutil" // some I/O utility functions
)

/* 	Implementation of Knuth-Morris-Pratt algorithm (Prefix based aproach).
	Requires two files in the folder with this file:
	
	@File pattern.txt containing the pattern to be searched for
	@File text.txt containing the text to be searched in
*/
func main() {
	// Error handling & file input
	patFile, err := ioutil.ReadFile("pattern.txt")
	if err != nil {
		log.Fatal(err)
	}
	textFile, err := ioutil.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}

	if (len(patFile) > len(textFile)) {
		log.Fatal("Pattern  is longer than text!")
	}
	// Alghoritm execution
	fmt.Printf("\nRunning: Knuth-Morris-Pratt alghoritm.\n\n")
	fmt.Printf("Search word (%d chars long): %q.\n",len(patFile), patFile)
	fmt.Printf("Text        (%d chars long): %q.\n\n",len(textFile), textFile)
	knp(string(textFile), string(patFile))
}

/*  Function knp performing the Knuth-Morris-Pratt alghoritm.
    Prints whether the word/pattern was found and on what position in the text or not.
	
	@param s string/text to be searched in
	@param w word/pattern to be serached for
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