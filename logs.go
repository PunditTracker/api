package main

import (
	"log"
	"log/syslog"
	"os"
)

var (
	request_f *os.File
)

func init() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "sys")
	if e == nil {
		log.SetOutput(logwriter)
	}

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
