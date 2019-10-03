package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	//TODO, figure out how to read in json in golang.

	// fmt.Println(conn)
	// fmt.Println("Starting jsonDecoder")

	// for {
	// 	dec := json.NewDecoder(conn)
	// 	var v map[string]interface{}
	// 	if err := dec.Decode(&v); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	for k := range v {
	// 		println(k)
	// 	}
	// }

	//TODO, I think you might need to return something here. Perhaps a channel.

    fmt.Println("testing")

}

func main() {

	args := os.Args
	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args[1:] {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]

		fmt.Println("Socket connection")

		//Open sockets and start listening.
		conn, err := net.Dial("unixpacket", "./"+ip[i])
		if err != nil {
			panic(err)
		}

		fmt.Println("Starting goroutines")
		go handleConnection(conn)

	}
	select {}
}
