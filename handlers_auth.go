package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	fmt.Fprintln(w, "logout")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username_val := r.FormValue("username")
	password_val := r.FormValue("password")
	if username_val == "" || password_val == "" {
		fmt.Fprintln(w, "username or password is blank")
		return
	}
	db, err := getDB()
	if err != nil {
		return
	}
	user := PtUser{
		Username: username_val,
		Password: password_val,
		Created:  time.Now(),
	}
	AddUser(db, user)
	fmt.Println("user added", username_val, password_val)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	r.ParseForm()
	username_val := r.FormValue("username")
	password_val := r.FormValue("password")
	num := CheckUser(db, username_val, password_val)

	if num == 0 {
		fmt.Fprintln(w, "failed log in")
		return
	}
	//Set up session or cookie
	kv := map[string]string{
		"uid": strconv.Itoa(int(num)),
	}
	setSession(kv, w)

	//num now set
	fmt.Fprintln(w, "logged in as", num)
}
