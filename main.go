package main

import (
	"log"
)

func main() {
	log.Println("Start Main")
	addListeners()
	beginServing()
	log.Println("End Main")
}
