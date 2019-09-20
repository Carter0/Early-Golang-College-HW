package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {

	//Dealing with the programs arguments.
	var port string
	var ssl bool
	flag.StringVar(&port, "p", "27993", "The port to connect to")
	flag.BoolVar(&ssl, "s", false, "Whether to use SSL connection are not. Defaults to false.")
	flag.Parse()

	//Some basic error checking.
	if len(flag.Args()) != 2 {
		println("Please input proper number of positional arguments.")
		return
	}

	hostname := flag.Arg(0)
	neuID := flag.Arg(1)

	CONNECT := hostname + ":" + port
	c, err := net.Dial("tcp", CONNECT)

	if err != nil {
		fmt.Println(err)
		return
	}

	//This sends the data to the tcp socket. Talks to tcp server.
	fmt.Fprintf(c, "cs3700fall2019 HELLO "+neuID+"\n")

	// Here is where we read input from the server
	for {

		//This reads the TCP servers response.
		message, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		//Split the string into a string array with a space delimeter.
		splitString := strings.Split(message, " ")

		//Find code below
		if splitString[1] == "FIND" {
			counter := 0
			toSearchString := splitString[2]
			searchThroughString := splitString[3]
			//Search through the string to find the right matches.
			for i := 0; i < len(searchThroughString); i++ {
				if searchThroughString[i] == toSearchString[0] {
					counter++
				}
			}
			fmt.Fprintf(c, "cs3700fall2019 COUNT ")
			fmt.Fprintf(c, "%d", counter)
			fmt.Fprintf(c, "\n")
		}

		//Bye code below, exit condition
		if splitString[1] == "BYE" {
			println(message)
			return
		}
	}

}
