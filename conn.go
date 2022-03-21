package main

import (
	"fmt"
	"net"
)

func getTarget(t []Target, client net.Conn) string {
	if c.LoadBalancer == "rr" {
		return roundRobin(t)
	} else if c.LoadBalancer == "lc" {
		return t[0].Ip + ":" + t[0].Port
	} else if c.LoadBalancer == "iphash" {
		ip, _, err := net.SplitHostPort(client.RemoteAddr().String())
		if err != nil {
			fmt.Println("Client ip format error", client.RemoteAddr().String())
		}
		return ipHash(t, ip)
	} else if c.LoadBalancer == "random" {
		return t[0].Ip + ":" + t[0].Port
	}

	// Default algorithm
	return roundRobin(t)
}

func clientRead(client net.Conn, server net.Conn) {
	buf := make([]byte, c.MaxRequest)

	for {
		n, err := client.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println("C: " + string(buf[0:]))

		_, err2 := server.Write(buf[0:n])

		if err2 != nil {
			return
		}
	}
}

func serverRead(server net.Conn, client net.Conn) {
	buf := make([]byte, c.MaxResponse)

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
