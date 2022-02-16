package main

import (
	"flag"
	"fmt"
)

func main() {
	var configFile string
	var listenPort int

	flag.StringVar(&configFile, "config", "config.json", "configureation file")
	flag.IntVar(&listenPort, "port", 25, "port number")

	flag.Parse()

	fmt.Println("Port:", listenPort)
	fmt.Println("Config:", configFile)

	c := Parse(configFile)

	fmt.Println(c)
}
