package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	//TODO, figure out how to read in json in golang.

	fmt.Println(conn)
	dec := json.NewDecoder(conn)
	for {
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			log.Println(err)
			return
		}

		for k := range v {
			println(k)
		}
	}

	// message, err := bufio.NewReader(conn).ReadString('\n')
	// if err != nil {
	// 	println("reader did not find anything")
	// }
	// print(message)

	//TODO, I think you might need to return something here. Perhaps a channel.

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

	fmt.Println("Argument Parsing")

	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]

		fmt.Println("Socket connection")

		//Open sockets and start listening.
		conn, err := net.Dial(networkType, "./"+ip[i])
		if err != nil {
			panic(err)
		}

		fmt.Println("Starting goroutines")
		go handleConnection(conn)

	}
}
