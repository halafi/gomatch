package main //package main required for standalone executable
import "fmt" //implements fomratted I/O.
import "os" //accessing command-line arguments
import "log" //simple logging package

/* 	implementation of Knuth-Morris-Pratt alghoritm.
	Requires two command line arguments.
	
	@argument string to be searched "for" (pattern, search word), no spaces allowed
	@argument one space
	@argument string to be searched "in" (text), spaces allowed
*/
func main() {
	// Error handling & declaration
	args := os.Args;
	if (len(args) <= 2) {
		log.Fatal("Not enough arguments. Two string arguments separated by spaces are required.");
	}
	pattern := args[1];
	s := args[2];
	for i := 3; i<len(args); i++ {
		s = s +" "+ args[i];
	}
	a := len(args[1]);
	b := len(s);
	if ( a > b ) {
		log.Fatal("Pattern (%d) is longer than text (%d).",a,b);
	}

	// Alghoritm execution
	fmt.Printf("\nRunning: Knuth-Morris-Pratt alghoritm.\n\n");
	fmt.Printf("Search word (%d chars long): '%s'.\n",a, pattern);
	fmt.Printf("Text        (%d chars long): '%s'.\n\n",b, s);

	knp(s, pattern);
}

/*  Function knp performing the Knuth-Morris-Pratt alghoritm.
    Prints whether the word/pattern was found in the text or not.
	
	@param s string/text to be searched in
	@param w word/pattern to be serached for
*/  
func knp(s, w string) {
	/*
	m := 0; //denotes the position within s which is the begining of a propective match for w
	i := 0; //index in w denoting the character currently under consideration
	*/
	fmt.Println("Unsupported Operation");
	
}