package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// openFile initializes buffered reader for file.
func openFile(filePath string) *bufio.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return bufio.NewReader(file)
}

// readLine reads a single text line into string using the given reader.
// Returns the line and true or false (whether EOF was reached or not).
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

// createFile creates file at filePath.
func createFile(filePath string) *os.File {
	newFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return newFile
}

// writeFile writes a string to file.
func writeFile(file *os.File, data string) {
	_, err := file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
}

// closeFile closes file.
func closeFile(file *os.File) {
	file.Close()
}
