package main

import "net"

type routingTable struct {
	updateMessage string
	conn          net.Conn
}

type ipInfo struct {
	networkIP   string
	networkMask string
}

func main() {

	var routing map[ipInfo]routingTable

	err, conn := net.Dial("tcp", "192.168.0.1")
	if err != nil {
		panic(err)
	}
	testRoute := routingTable{"update msg here", conn}
	testIP := ipInfo{"192.168.10.0", "255.255.255.0"}

	routing[testIP] = testRoute

}

/*
//This is not thread safe.
var routingtable = map[networkTuple][]*rtData{}


// rtData contains a copy of the original message sent to the router
// as well as necessary data for the routing table at the top level.
type rtData struct {
	relationshipType string
	Msg              []byte
	Conn             net.Conn
}
*/
