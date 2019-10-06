package gomod

import "rsc.io/quote"

// //Hello -> Simple hello function
// func Hello() string {
// 	return "Hello, world."
// }

func Hello() string {
	return quote.Hello()
}
