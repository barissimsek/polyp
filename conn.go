package main

import (
	"fmt"
	"net"
)

func getTarget(t []Target, client net.Conn) string {
	if c.LoadBalancer == "rr" {
		return roundRobin(t)
	} else if c.LoadBalancer == "lc" {
		return roundRobin(t)
	} else if c.LoadBalancer == "iphash" {
		ip, _, err := net.SplitHostPort(client.RemoteAddr().String())
		if err != nil {
			fmt.Println("Client ip format error", client.RemoteAddr().String())
		}
		return ipHash(t, ip)
	} else if c.LoadBalancer == "random" {
		return roundRobin(t)
	}

	// Default algorithm
	return roundRobin(t)
}

func serverRead(server net.Conn, client net.Conn) {
	buf := make([]byte, c.Processor.MaxResponse)

	for {
		n, err := server.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println("S: ", string(buf[0:]))

		_, err2 := client.Write(buf[0:n])

		if err2 != nil {
			return
		}
	}
}

func handleConnection(client net.Conn) {
	target := getTarget(c.Target, client)

	fmt.Println("Chosen target: ", target)

	server, err := net.Dial("tcp", target)

	if err != nil {
		fmt.Println("Can't establish relay connections")
		fmt.Println(err)
		return
	}

	go clientRead(client, server)
	go serverRead(server, client)

}
