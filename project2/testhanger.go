package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
)

var wg *sync.WaitGroup
var runLock *sync.RWMutex
var running *bool

func checkRunning() bool {
	runLock.RLock()
	defer runLock.RUnlock()
	return *running
}

func dumpMessages(ln net.Listener) {

	const SIZE = 2096
	buf := make([]byte, SIZE)

	defer ln.Close()
	defer wg.Done()

	for checkRunning() {
		s, err := ln.Accept()
		if err != nil {
			fmt.Println("33")
			panic(err)
		}

		for n, err := s.Read(buf); n > 0 && err == nil; n, err = s.Read(buf) {
			fmt.Print(string(buf))
		}

		if err != nil {
			fmt.Println("40")
			panic(err)
		}
	}
}

func main() {

	fmt.Println("Parse Args")

	var useUnix = flag.Bool("u", false, "Use plain unix socket instead of SOCK_SEQPACKET")
	flag.Parse()

	args := flag.Args()

	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]
	}
	var network string
	if *useUnix {
		network = "unix"
	} else {
		network = "unixpacket"
	}

	fmt.Println("Init Globals")

	wg = &sync.WaitGroup{}
	runLock = &sync.RWMutex{}
	var r bool = true
	running = &r
	var lc net.ListenConfig

	fmt.Println("Start GoRoutines")

	for _, addr := range ip {
		ln, err := lc.Listen(context.Background(), network, "./"+addr)
		if err != nil {
			fmt.Println("62")
			panic(err)
		}
		wg.Add(1)
		go dumpMessages(ln)
	}

	fmt.Println("Enter Wait")

	wg.Wait()
}
