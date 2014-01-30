// Package file provides input funcionality for a single log file.
package logdata

import "log"
import "io"
import "bufio"
import "os"

// Open does the initialization of buffered io.Reader.
func Open(filePath string) *bufio.Reader {
	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(fi)
	return reader
}

// ReadLine reads a single line using the given reader, returns the line
// and 'true' when EOF is reached, 'false' otherwise.
func ReadLine(reader *bufio.Reader) (logLine string, eof bool) {
	for {
		line, _, err := reader.ReadLine()
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
