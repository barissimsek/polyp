package main

import (
	"fmt"
	"os"

	"encoding/json"
)

type Config struct {
	Target       []Target `json:"targets"`
	LoadBalancer string   `json: "loadBalancer"`
	MaxRequest   int      `json:"maxRequest"`
	MaxResponse  int      `json:"maxResponse"`
}

type Target struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
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
