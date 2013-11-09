package main
import (
	"fmt" //implements fomratted I/O.
	"log" //simple logging package
	"os" //accessing command-line arguments
)

/* 	Implementation of Boyer-Moore-Horspool algorithm (Sufix based aproach).
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
	fmt.Printf("\nRunning: Horspool alghoritm.\n\n")
	fmt.Printf("Search word (%d chars long): %q.\n",len(args[1]), pattern)
	fmt.Printf("Text        (%d chars long): %q.\n\n",len(s), s)
	horspool(s, pattern)
}

/*  Function horspool performing the Horspool algorithm
    Prints whether the word/pattern was found and on what position in the text or not.
	
	@param t string/text to be searched in
	@param p word/pattern to be serached for
*/  
func horspool(t, p string) {
	m, n, c, pos := len(p), len(t), 0, 0
	//Perprocessing
	d := preprocess(t,p)
	//Map output
	fmt.Printf("Precomputed shifts per symbol: ")
	for key, value := range d {
		fmt.Printf("%c:%d; ", key, value)
	}
	fmt.Println()
	//Searching
	for pos <= n - m {
		j := m
		if (t[pos+j-1] != p[j-1]) {
			fmt.Printf("\n   comparing characters %c %c at positions %d %d",t[pos+j-1],p[j-1], pos+j-1, j-1)
			c++
		}
		for j > 0 && t[pos+j-1] == p[j-1] {
			fmt.Printf("\n   comparing characters %c %c at positions %d %d",t[pos+j-1],p[j-1], pos+j-1, j-1)
			c++
			fmt.Printf(" - match")
			j--
		}
		if j==0 {
				fmt.Printf("\n\nWord %q was found at position %d in %q. \n%d comparisons were done.",p, pos, t, c)
				return
			}
		pos = pos + d[t[pos + m ]]
	}
	fmt.Printf("\n\nWord was not found.\n%d comparisons were done.",c)
	return
}

/* 	Function that precomputes map with Key: uint8 (char) Value: int. Values determine safe shifting of search window.

	@Return map[uint8]int d filled map
*/ 
func preprocess(t, p string)(d map[uint8]int) {
	d = make(map[uint8]int)
	for i := 0; i < len(t); i++ {
		d[t[i]] = len(p)
	}
	for i := 0; i < len(p); i++ {
		d[p[i]] = len(p)-i
	}
	return d
}