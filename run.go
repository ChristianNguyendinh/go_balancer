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
	- add property to the worker struct that points to the ipaddr of the host
	- add a "addWorker method to the host". that will call MakeWorker, add it to the hosts list of recievers,
	and set the worker's host property to the host's ip
	*/

	worker1 := MakeWorker("test1", addr, ":8001")
	worker2 := MakeWorker("test2", addr, ":8002")
	worker3 := MakeWorker("test3", addr, ":8003")
	host := Host{name: "host", recievers: []Worker{worker1, worker2, worker3}, channel: make(chan string)}
	go Listen(host, ":8000")

	time.Sleep(3 * time.Second)

	// current format - :<TYPE>:arg1|arg2|arg3|...
	sendToHost(addr, ":INSTRUCTION:ls -l -h ..|ls -l -h ../../.")

	// wait for message to be recieved
	msg := <-host.channel
	log.Printf(msg + " <<< RECIEVED BY HOST ")
}