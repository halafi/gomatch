package main
import ("fmt"; "log"; "os"; "io/ioutil") 

/** 
	User defined.
	
	@true to take two command line arguments
	@false to take two files "pattern.txt" AND "text.txt"
*/
const commandLineInput bool = false

/**
	Implementation of Knuth-Morris-Pratt algorithm (Prefix based aproach).

	If(commandLineInput == true) requires two command line arguments separated by one space.
	@argument string to be searched "for" (pattern, search word), no spaces allowed
	@argument string to be searched "in" (text), single spaces allowed
	
	If(commandLineInput == false) requires two files in the same folder as this file.
	@file pattern.txt containing the pattern to be searched for
	@file text.txt containing the text to be searched in
*/
func main() {
	if (commandLineInput == true) { //in case of command line input
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
		fmt.Printf("\nRunning: Knuth-Morris-Pratt algorithm.\n\n")
		fmt.Printf("Search word (%d chars long): %q.\n",len(args[1]), pattern)
		fmt.Printf("Text        (%d chars long): %q.\n\n",len(s), s)
		knp(s, pattern)
	} else if (commandLineInput == false) { //in case of file input
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
		fmt.Printf("\nRunning: Knuth-Morris-Pratt algorithm.\n\n")
		fmt.Printf("Search word (%d chars long): %q.\n",len(patFile), patFile)
		fmt.Printf("Text        (%d chars long): %q.\n\n",len(textFile), textFile)
		knp(string(textFile), string(patFile))
	}
}

/**
	Function knp performing the Knuth-Morris-Pratt algorithm.
	Prints whether the word/pattern was found and on what position in the text or not.
	
	@param text string/text to be searched in
	@param word word/pattern to be serached for
*/  
func knp(text, word string) {
	m, i, c := 0, 0, 0 //m - current match in text, i - current character in w, c - ammount of comparations
	t := kmp_table(word)
	for  m + i < len(text) {
		fmt.Printf("\n   comparing characters %c %c at positions %d %d",text[m+i],word[i], m+i, i)
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

/**
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