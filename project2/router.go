package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	//TODO, figure out how to read in json in golang.

	fmt.Println(conn)

	// message, err := bufio.NewReader(conn).ReadString('\n')
	// if err != nil {
	// 	println("reader did not find anything")
	// }
	// print(message)

}

func main() {

	var useUnix = flag.Bool("u", false, "Use plain unix socket instead of SOCK_SEQPACKET")
	flag.Parse()

	args := flag.Args()

	var networkType string
	if *useUnix {
		networkType = "unix"
	} else {
		networkType = "unixpacket"
	}

	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]

		//Open sockets and start listening.
		conn, err := net.Dial(networkType, "./"+ip[i])
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)

	}
}
