package main

import (
	"fmt"
	"os"

	"encoding/json"
)

type Config struct {
	Target        []Target  `json:"targets"`
	LoadBalancer  string    `json:"loadBalancer"`
	HashTableSize int       `json:"hashTableSize"`
	Processor     Processor `json:"processor"`
}

type Target struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type Processor struct {
	Protocol    string `json:"protocol"`
	MaxRequest  int    `json:"maxRequest"`
	MaxResponse int    `json:"maxResponse"`
}

func Parse(configFile string) Config {
	var c Config

	byteData, _ := os.ReadFile(configFile)

	err := json.Unmarshal(byteData, &c)

	if err != nil {
		fmt.Println(err)
	}

	return c
}
