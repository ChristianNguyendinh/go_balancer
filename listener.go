package main

import (
	"log"
	"net"
)

/*

package? provides an Listen() method that will start listening on a specified port of the
host machine. Listen() takes in a provided ListenerInterface interface, that implements
a reqHandler method. This method is what Listen() will pass any incoming connections to

*/

// ListenerInterface s
type ListenerInterface interface {
	getName() string
	reqHandler(net.Conn)
}

// Listen method
func Listen(ls ListenerInterface, port string) {
	ln, err := net.Listen("tcp4", port)
	defer func() {
		ln.Close()
		log.Printf("Closing Listener for %s\n", ls.getName())
	}()

	if err != nil {
		panic(err)
	}
	log.Printf("%s listening on port %s\n", ls.getName(), port)

	// infinitely accept connections, start goroutine to handle each one
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go ls.reqHandler(conn)
	}
}
