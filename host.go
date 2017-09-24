package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type hostStruct struct {
	name string
}

// implement reqHandler for ListenerInterface
func (hostStruct) reqHandler(conn net.Conn) {
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

func (hs hostStruct) getName() string {
	return hs.name
}

func main() {
	host := hostStruct{name: "host"}
	Listen(host, ":8000")
}
