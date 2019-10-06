package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

// routingTable represents the list of paths for our network.
type routingTable struct {
	ip     string //ip to send data through
	subnet string //subet mask for the ip
	//TODO -> add more fields as needed later.
}

// message represents a message from a neighbor to the router.
type message struct {
	Msg  interface{} `json:"msg"`
	Src  string      `json:"src"`
	Dst  string      `json:"dst"`
	Type string      `json:"type"`
}

func handleConnection(conn net.Conn) {

	//var routes []routingTable
	for {
		var m message

		err := json.NewDecoder(conn).Decode(&m)
		if err != nil {
			log.Fatal("error decoding message ", err)
		}

		network := gjson.Get(m.Src, "network")
		print("The gjson msg value is: ")
		println(network)
	}

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
