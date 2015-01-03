package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func AddBracketHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var b PtBracket
	err := dec.Decode(&b)
	if err != nil {
		fmt.Println("Json Decode Error", err)
		return
	}
	b.Created = time.Now()

	db, err := getDB()
	if err != nil {
		fmt.Println("db error", err)
	}
	AddBracket(db, &b)
}

func GetBracketHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	User_Id, _ := strconv.ParseInt(vars["userId"], 10, 64)
	bracket := GetMembersBracket(db, User_Id)
	j, _ := json.Marshal(bracket)
	fmt.Fprintln(w, string(j))
}
