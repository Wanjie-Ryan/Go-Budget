package main

import "fmt"

// why pointers exist
// sometimes when creating variables and you want to pass them over, sometimes they just don't get passed over, only a copy of them is passed over.
// pointers ensure that the original value is always passed over, its a direct refrence to a memory address, hence when passign the value, it passed the memory address, hence ensuring that the actual value is passed.
func main() {

	// var a *int
	// fmt.Println("value of a", a)
	myNumber := 10
	fmt.Println("value of myNumber", myNumber)
	// pointing to a memory address where the value 10 is stored.
	// & is used to reference to memory
	var a = &myNumber
	fmt.Println("memory address of myNumber", a)
	fmt.Println("value inside the memory address of myNumber", *a)
	// *a is used to get the value inside the memory address
	// fmt.Println("the memory address of myNumber2", *&myNumber)

}
