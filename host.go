package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

// Host will recieve and pass on other work to children hosts/workers
type Host struct {
	name      string
	recievers []Worker
}

// implement reqHandler for ListenerInterface
func (Host) reqHandler(conn net.Conn) {
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

func (hs Host) getName() string {
	return hs.name
}

func main() {
	ip := "127.0.0.1"

	worker1 := MakeWorker("test", ip+":8000", ":8001")
	host := Host{name: "host", recievers: []Worker{worker1}}
	go Listen(host, ":8000")

	time.Sleep(3 * time.Second)

	for _, w := range host.recievers {
		log.Printf("%s - %s\n", w.getName(), w.status)
	}

	time.Sleep(3 * time.Second)

	for _, w := range host.recievers {
		// use channels here to avoid race below?
		w.sendResult("<SOME RESULT>")
	}

	// sometimes this program exits before worker finishes sending message - fix with chan
	time.Sleep(3 * time.Second)
}
