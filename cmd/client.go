package main

import (
	"fmt"
	"net"
)

func clientRead(client net.Conn, server net.Conn) {
	ret := make(chan int, 0)
	buf := make([]byte, c.Processor.MaxRequest)

	for {
		n, err := client.Read(buf[0:])
		if err != nil {
			return
		}

		fmt.Println("C: " + string(buf))

		switch c.Processor.Protocol {
		case "http":
			go proc_http(string(buf), ret)
		case "smtp":
			go proc_smtp(string(buf), ret)
		default:
			go proc_http(string(buf), ret)
		}

		// Get the result from processor
		// 0: Deliver
		// Non 0: Reject
		result := <-ret

		if result == 0 {
			_, err2 := server.Write(buf[0:n])

			if err2 != nil {
				return
			}
		} else {
			client.Close()
			server.Close()
		}
	}
}
