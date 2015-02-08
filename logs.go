package main

import (
	"log"
	"os"
)

var (
	request_f *os.File
)

func init() {
	db_f, err := os.OpenFile("db_log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic("no db log")
	}
	db_logger = log.New(db_f, "db: ", log.LstdFlags|log.Lshortfile)
	request_f, err = os.OpenFile("url_log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic("no url routing log")
	}
}
