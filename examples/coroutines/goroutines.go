package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//Channels provide a way for goroutines to communicate.
	c := make(chan string) //This line of code makes a channel.
	go boring("boring!", c)
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c) //print out a string (%q)
	}
	fmt.Println("You're boring; I'm leaving.")
}

func boring(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) //Throws both the msg and i into a string and sends
		// and sends it through the channel to the main f(x)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

//Below is a more advanced "pattern"
//Basically it returns the channel
//It is a generator.
//It makes an anonymous goroutine function that fills up the channel in the background
//Functionally does the same thing as in the first boring function, but it
//prevents the user from having to make the channel.
func boring2(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}() //I am assuming this is an anonymous function.
	return c
}
