package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
)

// openSocket opens a target Unix domain socket.
func openSocket(filePath string) net.Conn {
	conn, err := net.Dial("unix", filePath)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// readFully reads every line so far unread from the target socket.
// Returns true if EOF was reached, false otherwise.
func readFully(conn net.Conn) ([]string, bool) {
	toReturn := bytes.NewBuffer(nil)
	var buf [512]byte

	n, err := conn.Read(buf[0:])
	toReturn.Write(buf[0:n])
	if err != nil {
		if err == io.EOF {
			return lineSplit(string(toReturn.Bytes())), true
		}
		log.Fatal(err)
	}

	return lineSplit(string(toReturn.Bytes())), false
}

// lineSplit splits a multi-line string into single lines.
func lineSplit(input string) []string {
	split := make([]string, 1)

	// default return is single line
	split[0] = input
	if strings.Contains(input, "\r\n") { // CR+LF
		split = strings.Split(input, "\r\n")
	} else if strings.Contains(input, "\n") { // LF
		split = strings.Split(input, "\n")
	} else if strings.Contains(input, "\r") { // CR
		split = strings.Split(input, "\r")
	}
	return split
}
