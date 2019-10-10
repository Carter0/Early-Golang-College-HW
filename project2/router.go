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

type networkInfo struct {
	Msg  []message
	Conn net.Conn
}

var queue []net.Conn
var routingtable = map[networkTuple][]*rtData{}
var networkMap map[string]networkInfo
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
func createRTData(m message, tempType string) rtData {
	message, err := json.Marshal(m.Msg)
	if err != nil {
		panic(err)
	}
	return rtData{tempType, message}
}

func updateLogic(jsonMsg []byte, m message) {
	tempIP := gjson.GetBytes(jsonMsg, "network").String()
	tempSubnet := gjson.GetBytes(jsonMsg, "netmask").String()
	tempTuple := networkTuple{tempIP, tempSubnet}
	tempRoute := createRTData(m, m.Type)
	if val, ok := routingtable[tempTuple]; ok {
		val = append(val, &tempRoute)
	} else {
		rtArray := []*rtData{&tempRoute}
		routingtable[tempTuple] = rtArray
	}
}

func contains(s []message, e message) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func updateNeighbors() {

	// for _, mes := range messages {
	// 	if !contains(broadcastMessages, mes) {
	// 		toSendWrong := mes
	// 		broadcastMessages = append(broadcastMessages, mes)
	// 		for key, value := range routingtable {
	// 			if ()
	// 		}
	// 	}
	// }

	// if !networkTupleEquals(messageKey, key) {
	// 	for _, routeMessage := range messageList {
	// 		temp := *routeMessage
	// 		sendMessage := message{temp.Msg, temp.Dst, key.ip, temp.Type}
	// 		toSend, err := json.Marshal(sendMessage)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		for _, rtdata := range value {
	// 			rtdata.Conn.Write(toSend)
	// 		}
	// 	}
	// }

}

func handleConnection(conn net.Conn, networkName string) {

	println("Locking the code.")
	mutex.Lock()
	println("Adding conn to queue")
	addToQueue(conn)
	var msg message
	println("Decoding the message")
	err := json.NewDecoder(conn).Decode(&msg)
	if err != nil {
		panic(err)
	}

	println("Adding entry to network map")
	if val, ok := networkMap[networkName]; ok {
		val.Msg = append(val.Msg, msg)
	} else {
		println("Adding entry to empty map.")
		var temp []message
		temp = append(temp, msg)
		networkMap[networkName] = networkInfo{temp, conn}
	}

	mutex.Unlock()
	wg.Done()
	println("Unlocking the code and finish goroutine.")
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
		go handleConnection(conn, ip[i])
	}

	wg.Wait()

	println("Start looping through Queue")
	for _, conn := range queue {
		for key, value := range networkMap {
			fmt.Println(key)
			fmt.Println(value.Msg[0].Type)
		}

		var message message
		err := json.NewDecoder(conn).Decode(&message)
		if err != nil {
			panic(err)
		}

		jsonMsg, err := json.Marshal(message.Msg)
		if err != nil {
			panic(err)
		}

		println("Start looping through Message type.")
		switch message.Type {
		case "update":
			//Either add new info to the routing table
			//Or append data to the routing table values.
			println("Adding info to routing table.")
			updateLogic(jsonMsg, message)
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
