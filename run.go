package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

func sendToHost(addr string, msg string) {
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Write([]byte(msg))
}

func main() {
	ip := "127.0.0.1"
	addr := ip + ":" + strconv.Itoa(HOST_PORT)

	/*
	- add property to worker and host that holds its own ipaddr
	- add a "addWorker method to the host". that will call MakeWorker, add it to the hosts list of recievers,
	and set the worker's host property to the host's ip
	*/

	host := Host{name: "host", addr: addr, recievers: []Worker{}, channel: make(chan string)}

	w1 := MakeWorker("test1", addr, ":8001")
	w2 := MakeWorker("test2", addr, ":8002")
	w3 := MakeWorker("test3", addr, ":8003")
	
	host.addWorker(w1)
	host.addWorker(w2)
	host.addWorker(w3)

	go Listen(host, ":8000")

	time.Sleep(3 * time.Second)

	// send some work to all host's workers
	// will send a message on the host's channel when done
	// current format - :<TYPE>:arg1|arg2|arg3|...
	// go host.sendToWorkers(":INSTRUCTION:ls -lh|ls ..")

	// wait for message to be recieved
	// msg := <- host.channel
	// log.Printf(msg + " <<< RECIEVED BY HOST ")

	go host.splitToWorkers("ls", []string{".", "..", "../..", "../../.."})

	msg := <- host.channel
	log.Printf(msg + " <<< RECIEVED BY HOST ")
}