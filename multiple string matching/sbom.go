package main
import ("fmt"; "log"; /*"os";*/"strings"; "io/ioutil")

/**
 	Implementation of Set Backward Oracle Matching algorithm (Factor based aproach).
	Searches for a set of strings in file "patterns.txt" in text file text.txt.
	
	Requires two files in the same folder as the algorithm
	@file patterns.txt containing the patterns to be searched for separated by ", " 
		  !!!cannot end with ", " followed by nothing
	@file text.txt containing the text to be searched in
*/
func main() {
	patFile, err := ioutil.ReadFile("patterns.txt")
	if err != nil {
		log.Fatal(err)
	}
	textFile, err := ioutil.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	patterns := strings.Split(string(patFile), ", ")
	fmt.Printf("\nRunning: Set Backward Oracle Matching algorithm.\n\n")
	fmt.Printf("Searching for %d patterns/words:\n",len(patterns))
	for i := 0; i < len(patterns); i++ {
		if (len(patterns[i]) > len(textFile)) {
			log.Fatal("There is a pattern that is longer than text! Pattern number:", i+1)
		}
		fmt.Printf("%q ", patterns[i])
	}
	fmt.Printf("\n\nIn text (%d chars long): \n%q\n\n",len(textFile), textFile)
	sbom(string(textFile), patterns)
}

/**
	Function sbom performing the Set Backward Oracle Matching alghoritm.
	
	@param t string/text to be searched in
	@param p list of patterns to be serached for
*/  
func sbom(t string, p []string) {

	return
}

/**	
	Function that takes string and reverses it.
	
	@author Walter http://stackoverflow.com/a/10043083
	@param s string to be reversed
	@return reversed string
*/
func reverse(s string) string {
    l := len(s)
    m := make([]rune, l)
    for _, c := range s {
        l--
        m[l] = c
    }
    return string(m)
}