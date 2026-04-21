package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func proc_smtp(payload []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(payload))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(parts[0])
		arg := ""
		if len(parts) == 2 {
			arg = parts[1]
		}
		fmt.Printf("[smtp] %s %s\n", cmd, arg)
	}
	return 0
}
