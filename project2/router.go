package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

//Keep track of how many goroutines are running
var wg sync.WaitGroup

func handleConnection(conn net.Conn) {
	//TODO, figure out how to read in json in golang.

	fmt.Println(conn)
	fmt.Println("Starting jsonDecoder")

	for {
		dec := json.NewDecoder(conn)
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			log.Println(err)
			return
		}

		for k := range v {
			println(k)
		}
	}

	//TODO, I think you might need to return something here. Perhaps a channel.

}

func main() {

	var useUnix = flag.Bool("u", false, "Use plain unix socket instead of SOCK_SEQPACKET")
	flag.Parse()

	args := flag.Args()

	var networkType string
	if *useUnix {
		networkType = "unix"
	} else {
		networkType = "unixpacket"
	}

	fmt.Println("Argument Parsing")

	ip := make([]string, len(args))
	port := make([]string, len(args))
	for i, network := range args {
		split := strings.Split(network, "-")
		ip[i] = split[0]
		port[i] = split[1]

		fmt.Println("Socket connection")

		//Open sockets and start listening.
		conn, err := net.Dial(networkType, "./"+ip[i])
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
			wg.Done()
		}()

		//Wait till all goroutines are done before leaving program.(waitgroup # = 0)
		wg.Wait()
	}
}
