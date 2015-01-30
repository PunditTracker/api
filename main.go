package main

import (
	"fmt"
)

func main() {

	addListeners()
	beginServing()
	fmt.Println("post serve")
}
