package main

func roundRobin(t []Target) string {

	if tix == 255 || tix == len(t) {
		tix = 0
	}

	targetString := t[tix].Ip + ":" + t[tix].Port

	tix = tix + 1

	return targetString
}
