package main

import (
	"fmt"
)

func main() {
	db, err := getDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	SetUpDB(db)
	addListeners()
	beginServing()
	fmt.Println("post serve")
}
