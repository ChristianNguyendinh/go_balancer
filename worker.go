package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

// Worker will do the work
type Worker struct {
	name   string
	host   string
	status string
}

// implement reqHandler for ListenerInterface
func (wr Worker) reqHandler(conn net.Conn) {
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
			log.Println(wr.name+" Recieved data chunk: ", data[0:len(data)-4])
			break
		} else {
			log.Println(wr.name+"Recieved data chunk: ", data)
		}
	}

	w.Write([]byte("this is from " + wr.name))
	w.Flush()
}

func (wr Worker) sendResult(msg string) {
	conn, err := net.Dial("tcp", wr.host)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Write([]byte(wr.name + " reporting back results: " + msg))
	conn.Write([]byte("\r\n\r\n"))
}

func (wr Worker) getName() string {
	return wr.name
}

// MakeWorker will initialize and return a new worker
func MakeWorker(n string, h string, ps string) Worker {
	worker := Worker{name: n, host: h, status: "GOOD"}
	go Listen(worker, ps)
	return worker
}
