package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func AddBracketHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var b PtBracket
	err := dec.Decode(&b)
	if err != nil {
		fmt.Println("Json Decode Error", err)
		return
	}

	db, _ := getDB()
	AddBracket(db, &b)
}

func GetBracketHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	num, _ := strconv.Atoi(vars["user"])
	User_Id := int64(num)
	bracket := GetMembersBracket(db, User_Id)
	j, _ := json.Marshal(bracket)
	fmt.Fprintln(w, string(j))
}
