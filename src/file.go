// file.go does everything IO related for files.
package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// openFile does the initialization of buffered reader for a single file
// located at given filePath and returns it.
func openFile(filePath string) *bufio.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	return reader
}

// readLine reads a single text line using the given reader, returns the
// line and true when EOF is reached, line and false otherwise.
func readLine(reader *bufio.Reader) (string, bool) {
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

// createFile creates a single file at given filePath, returns pointer
// to that file.
func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// writeFile writes a string to file.
func writeFile(file *os.File, data string) {
	_, err := file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
}

// closeFile closes the given file.
func closeFile(file *os.File) {
	file.Close()
}
