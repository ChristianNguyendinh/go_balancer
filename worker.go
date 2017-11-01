package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"os/exec"
)

// Worker will do the work
type Worker struct {
	name   string
	host   string
	addr   string
	status string
	channel chan string
}

// implement reqHandler for ListenerInterface
func (wr Worker) reqHandler(conn net.Conn) {
	defer conn.Close()

	var (
		buff = make([]byte, 1024)
		r    = bufio.NewReader(conn)
		//w    = bufio.NewWriter(conn)
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

		// we are ignoring length for now
		// after we get the basics working, have a function that will handle the passed in instructions and continue
		if strings.HasPrefix(data, ":INSTRUCTION:") {
			log.Println(wr.name, " Recieved instruction: ", data)
			buff = buff[13:n]
		} else if strings.HasPrefix(data, ":RESULT:") {
			log.Println(wr.name, " Recieved result: ", data)
			buff = buff[8:n]
		} else {
			log.Println("Dunno what to do with data chunk: ", data)
			wr.channel <- "invalid command"
			return
		}
	}

	commands := strings.Split(string(buff), "|")
	log.Println(commands)
	for _, c := range(commands) {
		chunk := strings.Split(c, " ")
		cmd := exec.Command(chunk[0], chunk[1:]...)
		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("CMD Line Arg Messed up\n%s", err)
		}
		log.Printf(string(out))
	}

	// fill waiting channel with recieved message
	wr.channel <- string(buff)	
}

func (wr Worker) sendResult(msg string) {
	conn, err := net.Dial("tcp", wr.host)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Write([]byte(wr.name + " reporting back results: " + msg))
	conn.Write([]byte("\r\n\r\n"))
	log.Println("Writing to chan from worker send")
	//wr.channel <- (wr.name + " written")
}

func (wr Worker) getName() string {
	return wr.name
}

// MakeWorker will initialize and return a new worker
func MakeWorker(n string, h string, ip string, port string) Worker {
	address := ip + port
	worker := Worker{name: n, addr: address, host: h, status: "GOOD", channel: make(chan string)}
	go Listen(worker, port)
	return worker
}
