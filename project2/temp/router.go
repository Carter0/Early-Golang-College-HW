package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {

	var useUnix = flag.Bool("u", false, "Use plain unix socket instead of SOCK_SEQPACKET")
	flag.Parse()

	args := flag.Args()

	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]
	}
	var network string
	if *useUnix {
		network = "unix"
	} else {
		network = "unixpacket"
	}

	//TODO do some argument testing to make sure user puts in correct fields
	//TODO Open up len(os.args) many domain sockets and listen to them at the same time
	//This requires goroutines or polling?

	// Something is wrong here but I am not sure what.
	conn, err := net.Dial(network, "./"+ip[0])

	if err != nil {
		println("Error message for connection")
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", conn)

	//TODO, figure out how to read in json in golang.
	//Not sure if the reader is a good way to do it

	// //This reads the TCP servers response.
	// reader := bufio.NewReader(conn)
	// if reader == nil {
	// 	println("reader did not find anything")
	// }

	// message, err := reader.ReadString('\n')
	// if err != nil {
	// 	println("Error message for message")
	// 	fmt.Println(err)
	// }

	//print(message)
}
