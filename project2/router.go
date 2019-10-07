package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

var routes []routingTable

//IP4toInt converts an ip address into a binary sequence
func IP4toInt(IPv4Addr string) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(IPv4Addr).To4())
	return IPv4Int.Int64()
}

func handleConnection(conn net.Conn) {
	for {
		var m message
		var tempRoute routingTable

		err := json.NewDecoder(conn).Decode(&m)
		if err != nil {
			log.Fatal("error decoding message ", err)
		}

		temp, err := json.Marshal(m.Msg)
		if err != nil {
			panic(err)
		}
		tempIP := gjson.Get(string(temp), "network")
		tempSubnet := gjson.Get(string(temp), "netmask")
		tempRoute.ip = tempIP.String()
		tempRoute.subnet = tempSubnet.String()
		routes = append(routes, tempRoute)
	}

	//TODO send update messages to the other neighbors
	//In update message, change the src you are given to the destination
	//you are given. Use the neighbors source to send info over tcp
	//You are supposed to keep a running connection to your neighbors.
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
