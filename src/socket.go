// socket.go provides funcionality for unix domain sockets.
package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

// openSocket connects to target socket file, returns connection.
func openSocket(filePath string) net.Conn {
	con, err := net.Dial("unix", filePath)
	if err != nil {
		log.Fatal(err)
	}
	return con
}

// readFully reads everything so far unread from target socket file
// using estabilished connection into array of string lines.
// Returns true if EOF was reached, false otherwise.
func readFully(conn net.Conn) ([]string, bool) {
	result := bytes.NewBuffer(nil)
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	result.Write(buf[0:n])
	if err != nil {
		if err == io.EOF {
			return lineSplit(string(result.Bytes())), true
		}
		log.Fatal(err)
	}
	return lineSplit(string(result.Bytes())), false
}

// closeSocket closes the target connection.
func closeSocket(con net.Conn) {
	con.Close()
}
