package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func server(port string) {
	ln, err := net.Listen("tcp4", port)
	defer ln.Close()
	if err != nil {
		panic(err)
	}
	log.Printf("Server listening on port %s\n", port)

	// infinitely accept connections, start goroutine to handle each one
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go reqHandler(conn)
	}
}

func client(addr string) {
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	// send data to the server w/ the expected stop string
	conn.Write([]byte("this is from the client"))
	conn.Write([]byte("\r\n\r\n"))
	log.Printf("Sent message")

	// make a buffer slice and read data from the server
	buff := make([]byte, 1024)
	// read will block - if n == 0, connection has been closed by server
	n, err := conn.Read(buff)
	if err != nil {
		panic(err)
	}

	log.Printf("Recieve: %s\n", buff[:n])
}

// goroutine to handle conenctions
func reqHandler(conn net.Conn) {
	defer conn.Close()

	var (
		buff = make([]byte, 1024)
		r    = bufio.NewReader(conn)
		w    = bufio.NewWriter(conn)
	)

	// read until you get an EOF error or the client sends the stop string
	for {
		n, err := r.Read(buff)
		data := string(buff[:n])

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		// check if the data ends with this substring.
		// will fail if buffer ends before string finishes?
		if strings.HasSuffix(data, "\r\n\r\n") {
			log.Println("Recieved data chunk: ", data[0:len(data)-4])
			break
		} else {
			log.Println("Recieved data chunk: ", data)
		}
	}

	w.Write([]byte("this is from the server"))
	w.Flush()
	log.Printf("Sent back message")
}

// accepts exactly one argument, either 'client' or 'server'
func main() {
	if len(os.Args) != 2 {
		panic("need one arg of client of server")
	}

	port := ":8000"
	addr := "127.0.0.1" + port

	// run a client or server based on the cmd line arg
	if os.Args[1] == "client" {
		client(addr)
	} else if os.Args[1] == "server" {
		server(port)
	} else {
		panic("invalid arg")
	}
}
