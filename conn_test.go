package main

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

// startEchoServer starts a TCP server that echoes all received bytes back.
// Returns the listener address and a stop function.
func startEchoServer(t *testing.T) (addr string, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()
	return ln.Addr().String(), func() { ln.Close(); wg.Wait() }
}

// startProxy starts the proxy listener pointing at targetAddr.
// Returns the proxy address and a stop function.
func startProxy(t *testing.T, targetAddr string) (addr string, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	// Wire global config so handleConnection picks the right target.
	ip, port, _ := net.SplitHostPort(targetAddr)
	c = Config{
		Target:       []Target{{Ip: ip, Port: port}},
		LoadBalancer: "rr",
	}
	tix = 0

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn)
		}
	}()
	return ln.Addr().String(), func() { ln.Close(); wg.Wait() }
}

func TestProxyForwardsClientToServer(t *testing.T) {
	echoAddr, stopEcho := startEchoServer(t)
	defer stopEcho()

	proxyAddr, stopProxy := startProxy(t, echoAddr)
	defer stopProxy()

	conn, err := net.DialTimeout("tcp", proxyAddr, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	payload := []byte("hello proxy")
	if _, err := conn.Write(payload); err != nil {
		t.Fatal(err)
	}

	conn.(*net.TCPConn).CloseWrite()

	got, err := io.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Errorf("got %q, want %q", got, payload)
	}
}

func TestProxyForwardsLargePayload(t *testing.T) {
	echoAddr, stopEcho := startEchoServer(t)
	defer stopEcho()

	proxyAddr, stopProxy := startProxy(t, echoAddr)
	defer stopProxy()

	conn, err := net.DialTimeout("tcp", proxyAddr, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	payload := make([]byte, 1<<16) // 64 KiB
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	if _, err := conn.Write(payload); err != nil {
		t.Fatal(err)
	}
	conn.(*net.TCPConn).CloseWrite()

	got, err := io.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(payload) {
		t.Errorf("got %d bytes, want %d", len(got), len(payload))
	}
}

func TestProxyUnreachableTarget(t *testing.T) {
	// Point proxy at a port nothing is listening on.
	c = Config{
		Target:       []Target{{Ip: "127.0.0.1", Port: "1"}},
		LoadBalancer: "rr",
	}
	tix = 0

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		handleConnection(conn)
	}()

	conn, err := net.DialTimeout("tcp", ln.Addr().String(), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// handleConnection closes the client conn on dial failure; read must return EOF/error.
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err == nil {
		t.Error("expected connection to be closed by proxy, but read succeeded")
	}
}

func TestRoundRobin(t *testing.T) {
	targets := []Target{
		{Ip: "10.0.0.1", Port: "80"},
		{Ip: "10.0.0.2", Port: "80"},
		{Ip: "10.0.0.3", Port: "80"},
	}
	tix = 0

	want := []string{"10.0.0.1:80", "10.0.0.2:80", "10.0.0.3:80", "10.0.0.1:80"}
	for i, w := range want {
		got := roundRobin(targets)
		if got != w {
			t.Errorf("call %d: got %q, want %q", i, got, w)
		}
	}
}

func TestRoundRobinWrapsAt255(t *testing.T) {
	targets := []Target{{Ip: "10.0.0.1", Port: "80"}}
	tix = 255
	got := roundRobin(targets)
	if got != "10.0.0.1:80" {
		t.Errorf("got %q", got)
	}
	if tix != 1 {
		t.Errorf("tix should be 1 after wrap, got %d", tix)
	}
}

func TestIpHashReturnsSameTargetForSameIP(t *testing.T) {
	targets := []Target{
		{Ip: "10.0.0.1", Port: "80"},
		{Ip: "10.0.0.2", Port: "80"},
	}
	c = Config{HashTableSize: 16}
	ipMap = make(map[[16]byte]string, c.HashTableSize)
	tix = 0

	first := ipHash(targets, "192.168.1.1")
	second := ipHash(targets, "192.168.1.1")
	if first != second {
		t.Errorf("ipHash not stable: %q vs %q", first, second)
	}
}
