package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var c Config
var tix int
var ipMap map[[16]byte]string

func main() {
	var configFile string
	var listenPort int

	flag.StringVar(&configFile, "config", "config.json", "configuration file")
	flag.IntVar(&listenPort, "port", 80, "port number")
	flag.Parse()

	c = Parse(configFile)
	ipMap = make(map[[16]byte]string, c.HashTableSize)

	ln, err := net.Listen("tcp", ":"+fmt.Sprint(listenPort))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[%d] Poly is started.\n", os.Getpid())
	fmt.Printf("Listening port %d for incoming requests...\n", listenPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("New request from %s...\n", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}
