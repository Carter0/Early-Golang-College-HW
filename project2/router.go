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

var queue []message
var routingtable = map[networkTuple][]*rtData{}
var networkMap = make(map[string]networkInfo)
var mutex sync.Mutex
var wg sync.WaitGroup

//addToQueue adds an element to a queue
func addToQueue(msg message) {
	queue = append(queue, msg)
}

//removeFromQueue returns the first element of the list
func removeFromQueue() message {
	temp := queue[0]
	queue[0] = message{}
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

// func dataLogic(m message) {
// 	packetDst := IP4toInt(m.Dst)
// 	for key, value := range routingtable {
// 		netmask := IP4toInt(key.netMask)
// 		possibleToRoute := IP4toInt(key.ip) & netmask
// 		maskedPackDst := packetDst & netmask
// 		if possibleToRoute == maskedPackDst {

// 		}
// 	}
// }

func updateLogic(jsonMsg []byte, m message) {
	tempIP := gjson.GetBytes(jsonMsg, "network").String()
	println("The temp ip is " + tempIP)
	tempSubnet := gjson.GetBytes(jsonMsg, "netmask").String()
	tempTuple := networkTuple{tempIP, tempSubnet}
	tempRoute := createRTData(m, m.Type)
	if val, ok := routingtable[tempTuple]; ok {
		val = append(val, &tempRoute)
	} else {
		rtArray := []*rtData{&tempRoute}
		routingtable[tempTuple] = rtArray
	}
	println(len(routingtable))
	println("end of update")
}

//Updated neighbors forwards update messages to neighbors.
//The most disgusting fucntion ever that hopefully works.
func updateNeighbors() {
	for ip, netInfo := range networkMap {
		for _, msgToSendUnformatted := range netInfo.Msg {
			for ip2, netInfo := range networkMap {
				if ip2 != ip {
					sendMessage := message{msgToSendUnformatted.Msg, msgToSendUnformatted.Dst, ip2, msgToSendUnformatted.Type}
					toSend, err := json.Marshal(sendMessage)
					if err != nil {
						panic(err)
					}
					netInfo.Conn.Write(toSend)
					println("A message was sent.")
				}
			}
		}
	}
}

func handleConnection(conn net.Conn, networkName string) {
	mutex.Lock()

	var msg message
	err := json.NewDecoder(conn).Decode(&msg)
	if err != nil {
		panic(err)
	}

	println(msg.Type)
	addToQueue(msg)

	println("Adding entry to network map")
	if val, ok := networkMap[networkName]; ok {
		val.Msg = append(val.Msg, msg)
	} else {
		temp := []message{}
		temp = append(temp, msg)
		tempNet := networkInfo{temp, conn}
		networkMap[networkName] = tempNet
	}

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

		//Open sockets and start listening.
		conn, err := net.Dial("unixpacket", "./"+ip[i])
		if err != nil {
			panic(err)
		}

		fmt.Println("Starting goroutines")
		println(ip[i])
		wg.Add(1)
		go handleConnection(conn, ip[i])
	}

	wg.Wait()

	println("Start looping through Queue")
	for _, message := range queue {

		println(len(queue))

		jsonMsg, err := json.Marshal(message.Msg)
		if err != nil {
			panic(err)
		}

		println("The message type is " + message.Type)

		println("Start looping through Message type.")
		switch message.Type {
		case "update":
			println("Adding info to routing table.")
			updateLogic(jsonMsg, message)
		case "dump":
			fmt.Println("Dump logic here.")
		case "data":
			println("Forwarding update message to neighbors.")
			updateNeighbors()
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
