package main

import (
	"fmt"
)

func prepareDB() {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	SetUpDB(db)
}
func main() {
	prepareDB()
	addListeners()
	beginServing()
	fmt.Println("post serve")
}
