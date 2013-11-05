package main
import (
	"fmt" //implements fomratted I/O.
	"log" //simple logging package
	"io/ioutil" // some I/O utility functions
)

/* 	Implementation of Boyer-Moore-Horspool algorithm (Sufix based aproach).
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
	fmt.Printf("\nRunning: Horspool alghoritm.\n\n")
	fmt.Printf("Search word (%d chars long): %q.\n",len(patFile), patFile)
	fmt.Printf("Text        (%d chars long): %q.\n\n",len(textFile), textFile)
	horspool(string(textFile), string(patFile))
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