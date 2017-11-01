package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

// Host will recieve and pass on other work to children hosts/workers
type Host struct {
	name      string
	addr      string
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
			log.Println(h.name, " Recieved instruction: ", data)
			buff = buff[13:n]
			//HOW DOES A HOST SEND DATA TO RECIEVERS???			
		} else if strings.HasPrefix(data, ":RESULT:") {
			log.Println(h.name, "Recieved result: ", data)
			buff = buff[8:n]
		} else {
			log.Println("Dunno what to do with data chunk: ", data)
		}
	}

	// handle this stuff in the loop? incase we overflow
	// send msg back
	w.Write([]byte("this is from the server"))
	w.Flush()

}

func (hs Host) getName() string {
	return hs.name
}

func (hs *Host) addWorker(name string, port string) {
	var newWorker = MakeWorker(name, hs.addr, "127.0.0.1", port)
	hs.recievers = append(hs.recievers, newWorker)
}

func (hs Host) sendToWorkers() {
	for _, w := range hs.recievers {
		conn, err := net.Dial("tcp", w.addr)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
	
		go conn.Write([]byte(":INSTRUCTION:ls -lh|ls .."))

		message := <- w.channel
		log.Println(message)
	}

	hs.channel <- "done"
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
