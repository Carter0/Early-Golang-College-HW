package main

import (
	"encoding/json"
	"fmt"
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

type msgConnection struct {
	Msg  message
	conn net.Conn
}

//var broadcastMessages = map[networkTuple][]*message{}
var queue []net.Conn
var routingtable = map[networkTuple][]*rtData{}
var mutex sync.Mutex
var wg sync.WaitGroup

//addToQueue adds an element to a queue
func addToQueue(conn net.Conn) {
	queue = append(queue, conn)
}

//removeFromQueue returns the first element of the list
func removeFromQueue() net.Conn {
	temp := queue[0]
	queue[0] = nil
	queue = queue[1:]
	return temp
}

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

func updateLogic(jsonMsg []byte, conn net.Conn, m message) {
	tempIP := gjson.GetBytes(jsonMsg, "network").String()
	tempSubnet := gjson.GetBytes(jsonMsg, "netmask").String()
	tempTuple := networkTuple{tempIP, tempSubnet}
	tempRoute := createRTData(conn, m, m.Type)
	if val, ok := routingtable[tempTuple]; ok {
		val = append(val, &tempRoute)
	} else {
		rtArray := []*rtData{&tempRoute}
		routingtable[tempTuple] = rtArray
	}

	// //Broadcast logic here.
	// if val, ok := broadcastMessages[tempTuple]; ok {
	// 	val = append(val, &m)
	// } else {
	// 	toRouteArray := []*message{&m}
	// 	broadcastMessages[tempTuple] = toRouteArray
	// }
}

// func updateNeighbors() {
// 	for messageKey, messageList := range broadcastMessages {
// 		for key, value := range routingtable {
// 			if !networkTupleEquals(messageKey, key) {
// 				for _, routeMessage := range messageList {
// 					temp := *routeMessage
// 					sendMessage := message{temp.Msg, temp.Dst, key.ip, temp.Type}
// 					toSend, err := json.Marshal(sendMessage)
// 					if err != nil {
// 						panic(err)
// 					}
// 					for _, rtdata := range value {
// 						rtdata.Conn.Write(toSend)
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func handleConnection(conn net.Conn) {
	mutex.Lock()
	addToQueue(conn)
	mutex.Unlock()
	wg.Done()
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
		wg.Add(1)
		go handleConnection(conn)
	}

	wg.Wait()

	println("The counter is now 0 and it is printing a message.")

	/*
		TODO: Make a loop through the queue one at a time. Make sure to test this works properly.
	*/

	for _, conn := range queue {

		var message message
		err := json.NewDecoder(conn).Decode(&message)
		if err != nil {
			panic(err)
		}

		jsonMsg, err := json.Marshal(message.Msg)
		if err != nil {
			panic(err)
		}

		switch message.Type {
		case "update":
			//Either add new info to the routing table
			//Or append data to the routing table values.
			println("Adding info to routing table.")
			updateLogic(jsonMsg, conn, message)
		case "dump":
			fmt.Println("Dump logic here.")
		case "data":
			//Since all data functionality happens after update funcitonality,
			//put send data to neighbors stuff here.
			println("Forwarding update message to neighbors.")
			//updateNeighbors()
			fmt.Println("Data logic here.")
		}
	}

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
