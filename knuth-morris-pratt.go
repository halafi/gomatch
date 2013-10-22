package main //package main required for standalone executable
import "fmt" //implements fomratted I/O.
import "os" //accessing command-line arguments
import "log"

/* 	implementation of Knuth-Morris-Pratt alghoritm.
	Requires two command line arguments.
	
	@param string to be searched in (text)
	@param string to be searched for (pattern)
*/

func main() {
	args := os.Args;
	err := len(args);
	if (err > 3) {
		log.Fatal(err, " Too many arguments. Two separated by spaces are required. ");
	}

	fmt.Printf("Text: \n %s\n", args[1]);
	fmt.Printf("Pattern: \n %s\n", args[2]);
}


