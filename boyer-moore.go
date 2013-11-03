package main //package main required for standalone executable
import "fmt" //implements fomratted I/O.
import "os" //accessing command-line arguments
import "log" //simple logging package

/* 	Implementation of Knuth-Morris-Pratt alghoritm.
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
	fmt.Printf("Search word (%d chars long): '%s'.\n",len(args[1]), pattern)
	fmt.Printf("Text        (%d chars long): '%s'.\n\n",len(s), s)
	bm(s, pattern)
}

/*  Function bm performing the Boyer-Moore alghoritm.
    Prints whether the word/pattern was found and on what position in the text or not.
	
	@param s string/text to be searched in
	@param w word/pattern to be serached for
*/  
func bm(text, word string) {
	fmt.Printf("Unsuported operation")
}

