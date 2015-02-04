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
	defer db.Close()
	users := GetAllUsers(db)
	j, _ := json.Marshal(users)
	fmt.Fprintln(w, string(j))
}

//Add error handling
func GetSingleUserHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	user := GetUserByID(db, uid)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
}

func GetFeaturedUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	users := GetFeaturedUsers(db)
	j, _ := json.Marshal(users)
	fmt.Fprintln(w, string(j))
}

/*

func GetSingleUserForNameHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	var user PtUser
	vars := mux.Vars(r)
	username := vars["name"]
	db.Where("username = ?", username).First(&user)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
}*/
