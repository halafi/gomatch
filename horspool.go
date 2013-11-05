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
	fmt.Printf("\nRunning: Boyer-Moore-Horspool alghoritm.\n\n")
	fmt.Printf("Search word (%d chars long): %q.\n",len(patFile), patFile)
	fmt.Printf("Text        (%d chars long): %q.\n\n",len(textFile), textFile)
	horspool(string(textFile), string(patFile))
}

/*  Function knp performing the Knuth-Morris-Pratt alghoritm.
    Prints whether the word/pattern was found and on what position in the text or not.
	
	@param text string/text to be searched in
	@param word word/pattern to be serached for
*/  
func horspool(text, word string) {
	//Perprocessing
	d := horspool_table(word)  //line 3.
	//Map output
	for key, value := range d {
		fmt.Printf("%c %d\n", key, value)
	}
	//Searching
	pos := 0
	for pos <= len(text) - len(word) {
		/*j := len(word)
		for j>0 && text[pos+j]==word[j] {
			 j--
		}
		if j==0 {
			//fmt.Printf("\n\nWord %q was found at position %d in %q. \n%d comparisons were done.",word, m, text,c)
			fmt.Printf("\n\nWord %q was found at position %d in %q. \n",word, pos+1, text)
		}*/
		pos = pos +  d[text[pos + len(word)]]
		fmt.Printf("%c %d", d[text[pos + len(word)]], pos)
	}
	return
}

/*
	Preprocessing.
	
	@param word word to be analyzed
	@param d map to be filled
*/
func horspool_table(word string)(d map[uint8]int) {
	d = make(map[uint8]int) //map[KeyType]ValueType http://blog.golang.org/go-maps-in-action
	for i := 0; i<len(word); i++ {
		d[word[i]] = len(word)
		//fmt.Printf("%c %d\n",word[i], len(word))
	}
	for i := 1; i < (len(word)-1); i++ {
		d[word[i]] = len(word)-i
		//fmt.Printf("%c %d\n", word[i], len(word)-i)
	}
    return d
}