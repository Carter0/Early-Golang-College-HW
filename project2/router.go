package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
)

// rtData contains a copy of the original message sent to the router
// as well as necessary data for the routing table at the top level.
type rtData struct {
	relationshipType string
	Msg              []byte
	Conn             net.Conn
}

// networkTuple represents the list of paths for our network.
//TODO, this is not unique, consider port numbers as well.
type networkTuple struct {
	ip      string //ip to send data through
	netMask string //subet mask for the ip
}

// message represents a message from a neighbor to the router.
type message struct {
	Msg  interface{} `json:"msg"`
	Src  string      `json:"src"`
	Dst  string      `json:"dst"`
	Type string      `json:"type"`
}

var routingtable = map[networkTuple][]*rtData{}
var mutex sync.Mutex

//networkTupleEquals just determines if two network tuples are equal
func networkTupleEquals(net1 networkTuple, net2 networkTuple) bool {
	if net1.ip == net2.ip && net1.netMask == net2.netMask {
		return true
	}
	return false
}

//IP4toInt converts an ip address into a binary sequence
func IP4toInt(IPv4Addr string) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(IPv4Addr).To4())
	return IPv4Int.Int64()
}

//CreateRTData creates the routing table information
func createRTData(conn net.Conn, m message, tempType string) rtData {
	message, err := json.Marshal(m.Msg)
	if err != nil {
		panic(err)
	}
	return rtData{tempType, message, conn}
}

//Temp Solution
/*
	1. Create a separate list of messages to send out to other networks and map them to the networkTuple
	2. In the updateNeighbors function, loop through the routing table and correctly set the send out message
	3. Send out the message to all other neighbors but the neighbor that sent you the message
	4. After you are done, remove that message from the list.

	I think this works, but it doesn't feel like the best solution. But it does feel good enough.
*/

func updateLogic(jsonMsg []byte, conn net.Conn, m message) {
	tempIP := gjson.GetBytes(jsonMsg, "network").String()
	tempSubnet := gjson.GetBytes(jsonMsg, "netmask").String()
	tempTuple := networkTuple{tempIP, tempSubnet}
	tempRoute := createRTData(conn, m, m.Type)
	mutex.Lock()
	if val, ok := routingtable[tempTuple]; ok {
		val = append(val, &tempRoute)
	} else {
		rtArray := []*rtData{&tempRoute}
		routingtable[tempTuple] = rtArray
	}
	mutex.Unlock()
}

//Update neighbors forwards update messages to the neighbors.
func updateNeighbors() {
	// for forwardKey, forwardValue := range routingtable {
	// 	for key := range routingtable {
	// 		if !networkTupleEquals(forwardKey, key) {
	// 			for _, rtdata := range forwardValue {
	// 				temp := *rtdata
	// 				sendMessage := message{temp.Msg, temp.mes.Dst, key.ip, temp.mes.Type}
	// 				toSend, err := json.Marshal(sendMessage)
	// 				if err != nil {
	// 					panic(err)
	// 				}
	// 				temp.Conn.Write(toSend)
	// 			}
	// 		}
	// 	}
	// }
}

func handleConnection(conn net.Conn) {
	for {
		var mes message

		err := json.NewDecoder(conn).Decode(&mes)
		if err != nil {
			log.Fatal("error decoding message ", err)
		}

		jsonMsg, err := json.Marshal(mes.Msg)
		if err != nil {
			panic(err)
		}

		switch mes.Type {
		case "update":
			//Either add new info to the routing table
			//Or append data to the routing table values.
			println("Adding info to routing table.")
			updateLogic(jsonMsg, conn, mes)
		case "dump":
			fmt.Println("Dump logic here.")
		case "data":
			//Since all data functionality happens after update funcitonality,
			//put send data to neighbors stuff here.
			println("Forwarding update message to neighbors.")
			mutex.Lock()
			updateNeighbors()
			mutex.Unlock()

			fmt.Println("Data logic here.")
		}
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

	//Makes main thread doesn't close down.
	select {}
}

/* Code below can be useful for testing the contents of the routing table
for key, value := range routingtable {
	println("The key is: ")
	println(key.ip)
	println(key.netMask)
	for _, rtdata := range value {
		temp := *rtdata
		println(temp.relationshipType)
	}
}
*/
