package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	ip := make([]string, len(os.Args))
	port := make([]string, len(os.Args))
	for i, network := range os.Args[1:] {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]
		//TODO do some argument testing to make sure user puts in correct fields
		//TODO Open up len(os.args) many domain sockets and listen to them at the same time
		//This requires goroutines or polling?
	}

	// Something is wrong here but I am not sure what.
	conn, err := net.Dial("unixpacket", "./"+ip[0])

	if err != nil {
		println("Error message for connection")
		println(err)
	}

	//This reads the TCP servers response.
	reader := bufio.NewReader(conn)
	if reader == nil {
		println("reader did not find anything")
	}

	message, err := reader.ReadString('\n')
	if err != nil {
		println("Error message for message")
		fmt.Println(err)
	}

	print(message)
}
