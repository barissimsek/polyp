package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var c Config
var tix int = 0
var ipMap = make(map[[16]byte]string, c.HashTableSize)

func main() {
	var configFile string
	var listenPort int

	flag.StringVar(&configFile, "config", "config.json", "configureation file")
	flag.IntVar(&listenPort, "port", 80, "port number")

	flag.Parse()

	ln, err := net.Listen("tcp", ":"+fmt.Sprint(listenPort))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("[%d] Poly is started.\n", os.Getpid())
	fmt.Printf("Listening port %d for incoming requests...\n", listenPort)

	c = Parse(configFile)

	fmt.Println(c)

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("New request from %s...\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}
