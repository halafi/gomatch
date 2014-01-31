package main

import "bufio"
import "os"
import "log"
import "io"

// Function openFile does the initialization of buffered io.Reader for 
// file at given filePath.
func openFile(filePath string) *bufio.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	return reader
}

// ReadLine reads a single line using the given reader, returns the line
// and 'true' when EOF is reached, 'false' otherwise.
func readLine(reader *bufio.Reader) (string, bool) {
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

// Creates file for writing output at 'filePath'.
func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// Writes given 'line' string to the given file.
func writeFile(file *os.File, data string) {
	_, err := file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
}

// Closes the given 'file'.
func closeFile(file *os.File) {
	file.Close()
}
