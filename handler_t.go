package main

import (
	"net/http"
)

func LoadTestDataHandler(w http.ResponseWriter, r *http.Request) {

	db, _ := getDB()
	AddUser(db, &PtUser{
		Username:  "ben",
		Password:  "pass",
		FirstName: "ben",
		LastName:  "levy"})
	AddUser(db, &PtUser{
		Username:  "howie",
		Password:  "password",
		FirstName: "howard",
		LastName:  "akumiah"})

}
