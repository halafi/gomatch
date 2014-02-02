package main

import "net"
import "log"
import "bytes"
import "io"

// client side

// Open target unix domain socket file.
func openSocket(filePath string) (net.Conn){
	c, err := net.Dial("unix", filePath)
	if err != nil {
        log.Fatal(err)
    }
    return c
}

// Client sends data to server.
func write(conn net.Conn, text string) {
	_, err := conn.Write([]byte(text))
	if err != nil {
		log.Fatal(err)
	}
}

// Reads everything from socket into array of string lines, returns true
// if EOF was reached, false otherwise.
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

// Close target connection.
func closeSocket(c net.Conn) {
	c.Close()
}

// server side for testing
func startServer(filePath string) net.Listener {
	l, err := net.Listen("unix", filePath)
    if err != nil {
        log.Fatal("listen error", err)
    }
	return l
}

func echoServer(c net.Conn) {
    for {
        buf := make([]byte, 512)
        nr, err := c.Read(buf)
        if err != nil {
            return
        }
        data := buf[0:nr]
        log.Println("Server got:", string(data))
        _, err = c.Write(data)
        if err != nil {
            log.Fatal("Write: ", err)
        }
    }
}

func closeServer(l net.Listener) {
	l.Close()
}
