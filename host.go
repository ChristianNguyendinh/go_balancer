package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"os/exec"
)

// Host will recieve and pass on other work to children hosts/workers
type Host struct {
	name      string
	recievers []Worker
	channel   chan string
}

// implement reqHandler for ListenerInterface
func (h Host) reqHandler(conn net.Conn) {
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
		if strings.HasPrefix(data, ":INSTRUCTION:") {
			log.Println("Recieved instruction: ", data)
			buff = buff[13:]
		} else if strings.HasPrefix(data, ":RESULT:") {
			log.Println("Recieved result: ", data)
			buff = buff[8:]
			/*} else if strings.HasSuffix(data, "\r\n\r\n") {
			log.Println("Recieved end data chunk: ", data[0:len(data)-4])
			break */
		} else {
			log.Println("Dunno what to do with data chunk: ", data)
		}
	}

	// send msg back
	w.Write([]byte("this is from the server"))
	w.Flush()
	// fill waiting channel with recieved message
	cmd := exec.Command("ls", "-l", "-h")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalln("CMD Line Arg Messed up")
	}
	log.Printf(string(out))	
	// handle this stuff in the loop? incase we overflow
	h.channel <- string(buff) //strings.Split(string(buff), "|")
}

func (hs Host) getName() string {
	return hs.name
}

// func main() {
// 	ip := "127.0.0.1"
// 	addr := ip + ":" + strconv.Itoa(HOST_PORT)

// 	worker1 := MakeWorker("test1", addr, ":8001")
// 	worker2 := MakeWorker("test2", addr, ":8002")
// 	worker3 := MakeWorker("test3", addr, ":8003")
// 	host := Host{name: "host", recievers: []Worker{worker1, worker2, worker3}, channel: make(chan string)}
// 	go Listen(host, ":8000")

// 	time.Sleep(3 * time.Second)

// 	for _, w := range host.recievers {
// 		log.Printf("%s - %s\n", w.getName(), w.status)
// 	}

// 	time.Sleep(3 * time.Second)

// 	for _, w := range host.recievers {
// 		go w.sendResult("<SOME RESULT>")
// 	}

// 	// wait for messages to be recieved
// 	for _, w := range host.recievers {
// 		msg := <-host.channel
// 		log.Printf(msg + " from " + w.name)
// 	}
// }
