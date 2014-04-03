package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// openFile initializes a buffered reader for a file.
func openFile(filePath string) *bufio.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return bufio.NewReader(file)
}

// readLine reads a single text line using the given reader.
// Returns the line and true if EOF was reached, false otherwise.
func readLine(reader *bufio.Reader) ([]byte, bool) {
	line, _, err := reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return line, true
		} else {
			log.Fatal(err)
		}
	}
	return line, false
}

// createFile creates a file.
func createFile(filePath string) *os.File {
	newFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return newFile
}

// writeFile writes a string into a file.
func writeFile(file *os.File, data string) {
	_, err := file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
}
