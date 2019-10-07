package main

import "fmt"

const MAX int = 3

func main() {

	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j

	a := []int{10, 100, 200}
	var i2 int
	var ptr [MAX]*int

	for i = 0; i2 < MAX; i2++ {
		ptr[i2] = &a[i2] /* assign the address of integer. */
	}
	for i2 = 0; i2 < MAX; i2++ {
		fmt.Printf("Value of a[%d] = %d\n", i, *ptr[i])
	}
}
