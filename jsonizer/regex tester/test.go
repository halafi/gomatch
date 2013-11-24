package main
import ("regexp"; "fmt";"io/ioutil"; "strings")

func main() {
	/*The syntax of the regular expressions accepted is the same general syntax used by Perl,
	PYTHON!, and other languages. More precisely, it is the syntax accepted by RE2 and described
	at http://code.google.com/p/re2/wiki/Syntax, except for \C.*/
	patFile,_ := ioutil.ReadFile("testedRegex.txt")
	patText := strings.Split(string(patFile), " ")
	matched,err := regexp.MatchString(patText[0], patText[1])
	fmt.Printf("\n%t\n\n\nErrors: ", matched)
	fmt.Println(err)
}
	