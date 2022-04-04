package main

import (
	"crypto/md5"
)

func roundRobin(t []Target) string {

	if tix == 255 || tix == len(t) {
		tix = 0
	}

	targetString := t[tix].Ip + ":" + t[tix].Port

	tix = tix + 1

	return targetString
}

func ipHash(t []Target, ip string) string {
	hash := md5.Sum([]byte(ip))

	// Cache hit
	if val, ok := ipMap[hash]; ok {
		return val
	} else { // cache miss
		if len(ipMap) == c.HashTableSize { // cache full, evict
			for k, _ := range ipMap {
				delete(ipMap, k)
				break
			}
		}
	}

	// Add missed hash key
	ipMap[hash] = roundRobin(t)

	return ipMap[hash]
}
