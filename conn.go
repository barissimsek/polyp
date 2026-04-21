package main

import (
	"fmt"
	"io"
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

func inspect(payload []byte) {
	if c.Processor == nil {
		return
	}
	switch c.Processor.Protocol {
	case "http":
		proc_http(payload)
	case "smtp":
		proc_smtp(payload)
	}
}

func handleConnection(client net.Conn) {
	defer client.Close()

	target := getTarget(c.Target, client)
	fmt.Println("Chosen target: ", target)

	server, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("Can't establish relay connection:", err)
		return
	}
	defer server.Close()

	if c.Processor != nil {
		buf := make([]byte, c.Processor.MaxRequest)
		n, err := client.Read(buf)
		if n > 0 {
			inspect(buf[:n])
			if _, werr := server.Write(buf[:n]); werr != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}

	done := make(chan struct{}, 2)
	go func() {
		io.Copy(server, client)
		if tc, ok := server.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
		done <- struct{}{}
	}()
	go func() {
		io.Copy(client, server)
		if tc, ok := client.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
		done <- struct{}{}
	}()
	<-done
	<-done
}
