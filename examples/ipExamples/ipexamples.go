package main

import (
	"math/big"
	"net"
	"os"
	"strconv"
)

//IP4toInt -> converts an ip address into a
func IP4toInt(IPv4Addr string) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(IPv4Addr).To4())
	return IPv4Int.Int64()
}

func main() {

	// //Figure out how to convert an ip address into binary.
	// s1 := net.ParseIP("192.0.2.33")
	// if s1 != nil {
	// 	fmt.Printf("%#v\n", s1)
	// }
	// s2 := net.ParseIP("2001:db8:8714:3a90::12")
	// if s2 != nil {
	// 	fmt.Printf("%#v\n", s2)
	// }

	// for _, item := range s2 {
	// 	var i64 int64
	// 	i64 = int64(item)
	// 	println(strconv.FormatInt(i64, 2))
	// }

	input := os.Args[1:]
	network := IP4toInt(input[0])
	subnet := IP4toInt(input[2])

	//The and of the network to route to and the subnet.
	result := network & subnet
	binaryResult := strconv.FormatInt(result, 2)
	println(binaryResult)

	//The and of something located the network.
	result2 := IP4toInt(input[1])
	binaryResult2 := strconv.FormatInt(result2, 2)
	println(binaryResult2)

	for i := range binaryResult {
		if binaryResult[i] != binaryResult2[i] {
			print("The counter is ")
			println(i)
		}
	}

}
