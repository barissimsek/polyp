package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Target        []Target   `json:"targets"`
	LoadBalancer  string     `json:"loadBalancer,omitempty"`
	HashTableSize int        `json:"hashTableSize,omitempty"`
	Processor     *Processor `json:"processor,omitempty"`
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

	byteData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(byteData, &c); err != nil {
		fmt.Println("Error parsing config:", err)
		os.Exit(1)
	}

	return c
}
