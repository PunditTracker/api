package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	users := GetAllUsers(db)
	j, _ := json.Marshal(users)
	fmt.Fprintln(w, string(j))
}

//Add error handling
func GetSingleUserHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	uid, _ := strconv.Atoi(vars["id"])
	user := GetUserByID(db, uid)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
}

func GetFeaturedUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	users := GetFeaturedUsers(db)
	j, _ := json.Marshal(users)
	fmt.Fprintln(w, string(j))
}
