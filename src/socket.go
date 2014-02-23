package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
)

// openSocket connects to target socket file.
func openSocket(filePath string) net.Conn {
	connection, err := net.Dial("unix", filePath)
	if err != nil {
		log.Fatal(err)
	}
	return connection
}

// readFully reads everything so far unread from target socket into
// array of string lines.
// Returns true if EOF was reached, false otherwise.
func readFully(connection net.Conn) ([]string, bool) {
	toReturn := bytes.NewBuffer(nil)
	var buf [512]byte
	n, err := connection.Read(buf[0:])
	toReturn.Write(buf[0:n])
	if err != nil {
		if err == io.EOF {
			return lineSplit(string(toReturn.Bytes())), true
		}
		log.Fatal(err)
	}
	return lineSplit(string(toReturn.Bytes())), false
}

// closeSocket closes the target socket connection.
func closeSocket(connection net.Conn) {
	connection.Close()
}

// lineSplit parses a mutli-line string into single lines (array of
// strings).
func lineSplit(input string) []string {
	inputSplit := make([]string, 1)
	inputSplit[0] = input                // default single line
	if strings.Contains(input, "\r\n") { //CR+LF
		inputSplit = strings.Split(input, "\r\n")
	} else if strings.Contains(input, "\n") { //LF
		inputSplit = strings.Split(input, "\n")
	} else if strings.Contains(input, "\r") { //CR
		inputSplit = strings.Split(input, "\r")
	}
	return inputSplit
}
