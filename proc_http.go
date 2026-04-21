package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
)

func proc_http(payload []byte) int {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(payload)))
	if err != nil {
		fmt.Println("[http] failed to parse request:", err)
		return 0
	}
	fmt.Printf("[http] %s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
	for name, vals := range req.Header {
		for _, v := range vals {
			fmt.Printf("[http]   %s: %s\n", name, v)
		}
	}
	return 0
}
