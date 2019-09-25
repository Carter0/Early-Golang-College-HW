package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

func main() {

	//TODO, get proper command line parsing and then pass in the proper ip as the
	// second argument to net.Dial()

	//Look at the example code for how the proffessor parsed the command line args.

	//Basic argument parsing
	var port string
	var ip string
	flag.StringVar(&port, "p", "27993", "The port to connect to")
	flag.StringVar(&ip, "ip", "", "The ip address to connect to")
	flag.Parse()

	// Something is wrong here but I am not sure what.
	conn, err := net.Dial("unixpacket", "/tmp/echo.sock")

	if err != nil {
		println("Error message for connection")
		println(err)
	}

	//This reads the TCP servers response.
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		println("Error message for message")
		fmt.Println(err)
	}

	print(message)
}
