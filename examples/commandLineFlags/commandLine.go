package main

import (
	"flag"
	"fmt"
)

func main() {

	//Flags or Optional Arguments below
	//This is the pointer version
	wordPtr := flag.String("word", "foo", "a string")

	//This is the value version
	var stringVar string
	flag.StringVar(&stringVar, "svar", "bar", "a string var")

	//Make sure to remember to do this before looking at any of the arguments passed in.
	flag.Parse()

	//Testing for positional arguments
	if len(flag.Args()) != 2 {
		println("Please input proper number of positional arguments.")
		return
	}

	//This doesn't seem to work for some reason if you don't test to make sure the proper number above is put in.
	if flag.Arg(1) == "world" {
		println("Don't do that.")
		return
	}

	fmt.Println("word:", *wordPtr)
	fmt.Println("wordValue:", stringVar)
	//Trailing Positional Arguments.
	fmt.Println("Tail/Trailing positional arguments: ", flag.Args())
	fmt.Println("First positional argument: ", flag.Arg(0))
	fmt.Println("Second positional argument: ", flag.Arg(1))

}
