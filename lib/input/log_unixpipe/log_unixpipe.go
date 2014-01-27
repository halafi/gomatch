// Package unixpipe provides input funcionality for STDIN.
package unixpipe

import "io"
import "log"
import "bufio"
import "os"

// Init does nitialization of buffered io.Reader.
func Init() *bufio.Reader {
	reader := bufio.NewReader(os.Stdin)
	return reader
}

// ReadLine reads a single line using the given reader, return the line
// in string with a boolean value 'true' when EOF is reached.
func ReadLine(reader *bufio.Reader) (logLine string, eof bool) {
    for {
        line, _, err:= reader.ReadLine()
        if err != nil {
			if err == io.EOF {
				return string(line), true
			} else {
				log.Fatal(err)
			}
        }
        return string(line), false
    }
}
