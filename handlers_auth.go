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

//Get correct form values
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username_val := r.PostFormValue("username")
	password_val := r.PostFormValue("password")
	db, err := getDB()
	if err != nil {
		return
	}
	username_val = "USER2"
	password_val = "password2"
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
	username_val = "USER2"
	password_val = "password2"

	u := PtUser{
		Username: username_val,
		Password: password_val,
	}

	LoginUser(db, &u)
	if u.Id == 0 {
		fmt.Fprintln(w, "failed log in")
		return
	}
	//Set up session or cookie
	kv := map[string]string{
		"uid": strconv.Itoa(int(u.Id)),
	}
	setSession(kv, w)

	//u.Id is now set
	fmt.Fprintln(w, u.Id)
}
