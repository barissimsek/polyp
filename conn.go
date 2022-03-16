package main

import (
	"fmt"
	"net"
)

func getTarget(t []Target) string {
	// TODO: Implement LRU loadbalancer
	return t[0].Ip + ":" + t[0].Port
}

func clientRead(c net.Conn, s net.Conn) {
	buf := make([]byte, 2097152)

	for {
		n, err := c.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println("C: " + string(buf[0:]))

		_, err2 := s.Write(buf[0:n])

		if err2 != nil {
			return
		}
	}
}

func serverRead(s net.Conn, c net.Conn) {
	buf := make([]byte, 2097152)

	for {
		n, err := s.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println("S: ", string(buf[0:]))

		_, err2 := c.Write(buf[0:n])

		if err2 != nil {
			return
		}
	}
}

func handleConnection(client net.Conn, c Config) {
	target := getTarget(c.Target)

	server, err := net.Dial("tcp", target)

	if err != nil {
		fmt.Println("Can't establish relay connections")
		fmt.Println(err)
		return
	}

	go clientRead(client, server)
	go serverRead(server, client)

}
