package main //package main required for standalone executable
import "fmt" //implements fomratted I/O.
import "os" //accessing command-line arguments
import "log" //simple logging package

/* 	implementation of Knuth-Morris-Pratt alghoritm.
	Requires two command line arguments.
	
	@argument string to be searched "in" (text)
	@argument string to be searched "for" (pattern)
*/

func main() {
	// Error handling & declaration
	args := os.Args;
	if (len(args) > 3) {
		log.Fatal("Too many arguments. Two string arguments separated by spaces are required.");
	} 
	if (len(args) <= 2) {
		log.Fatal("Not enough arguments. Two string arguments separated by spaces are required.");
	}
	text := args[1];
	pattern := args[2];
	if (len(args[1])<len(args[2])) {
		log.Fatal("Pattern is longer than text.");
	}
	// Alghoritm execution
	fmt.Printf("Running: Knuth-Morris-Pratt alghoritm.\nText: %s\nPattern: %s\n", text, pattern);
}


