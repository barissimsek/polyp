package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	var configFile string
	var listenPort int

	flag.StringVar(&configFile, "config", "config.json", "configureation file")
	flag.IntVar(&listenPort, "port", 25, "port number")

	flag.Parse()

	ln, err := net.Listen("tcp", ":"+fmt.Sprint(listenPort))

	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}

	c := Parse(configFile)

	fmt.Println(c)
}
