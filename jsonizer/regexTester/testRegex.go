package main
import ("regexp"; "fmt";"io/ioutil";)

func main() {
	/*The syntax of the regular expressions accepted is the same general syntax used by Perl,
	PYTHON!, and other languages. More precisely, it is the syntax accepted by RE2 and described
	at http://code.google.com/p/re2/wiki/Syntax, except for \C.*/
	patFile,_ := ioutil.ReadFile("testedRegex.txt")
	matched,err := regexp.MatchString((string(patFile)), "50")
	fmt.Println(matched, err)
}
	