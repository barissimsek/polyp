package main

import (
	"crypto/md5"
	"sync"
)

var lbMu sync.Mutex

func roundRobin(t []Target) string {
	lbMu.Lock()
	defer lbMu.Unlock()
	s := t[tix%len(t)].Ip + ":" + t[tix%len(t)].Port
	tix = (tix + 1) % len(t)
	return s
}

func ipHash(t []Target, ip string) string {
	hash := md5.Sum([]byte(ip))

	lbMu.Lock()
	defer lbMu.Unlock()

	if val, ok := ipMap[hash]; ok {
		return val
	}

	if c.HashTableSize > 0 && len(ipMap) >= c.HashTableSize {
		for k := range ipMap {
			delete(ipMap, k)
			break
		}
	}

	target := t[tix%len(t)].Ip + ":" + t[tix%len(t)].Port
	tix = (tix + 1) % len(t)
	ipMap[hash] = target
	return target
}
