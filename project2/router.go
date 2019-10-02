package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

//Keep track of how many goroutines are running
var wg sync.WaitGroup

func handleConnection(conn net.Conn) {
	defer wg.Done()
	//TODO, figure out how to read in json in golang.

	// fmt.Println(conn)
	// fmt.Println("Starting jsonDecoder")

	// for {
	// 	dec := json.NewDecoder(conn)
	// 	var v map[string]interface{}
	// 	if err := dec.Decode(&v); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	for k := range v {
	// 		println(k)
	// 	}
	// }

	//TODO, I think you might need to return something here. Perhaps a channel.

}

func main() {

	args := os.Args

	fmt.Println("Argument Parsing")

	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args[1:] {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]

		fmt.Println(ip[i])
		fmt.Println(port[i])

		fmt.Println("Socket connection")

		//Open sockets and start listening.
		conn, err := net.Dial("unixpacket", "./"+ip[i])
		if err != nil {
			panic(err)
		}

		//Create a lamba function that starts immediatly.
		//Add on to the waitgroup
		//When goroutine is done subtract one from the waitgroup
		go func() {
			fmt.Println("Starting goroutines")
			wg.Add(1)
			go handleConnection(conn)
		}()

	}

	//Wait till all goroutines are done before leaving program.(waitgroup # = 0)
	wg.Wait()
}
