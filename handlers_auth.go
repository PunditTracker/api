package main

import (
	"encoding/json"
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
	first_val := r.FormValue("firstname")
	last_val := r.FormValue("lastname")
	if username_val == "" ||
		password_val == "" ||
		first_val == "" ||
		last_val == "" {
		errMessage := "missing values"
		http.Error(w, errMessage, http.StatusBadRequest)
		return
	}
	db, err := getDB()
	if err != nil {
		return
	}
	user := PtUser{
		Username:  username_val,
		Password:  password_val,
		FirstName: first_val,
		LastName:  last_val,
		Created:   time.Now(),
	}
	AddUser(db, &user)

	fmt.Println("user added", user)
}

func RegisterFacebookHandler(w http.ResponseWriter, r *http.Request) {
	fb_token := r.FormValue("fb_token")
	fb_id := r.FormValue("fb_id")
	username_val := r.FormValue("username")
	first_val := r.FormValue("firstname")
	last_val := r.FormValue("lastname")

	if username_val == "" ||
		first_val == "" ||
		last_val == "" ||
		fb_token == "" {
		errMessage := "missing values"
		http.Error(w, errMessage, http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		return
	}
	user := PtUser{
		Username:          username_val,
		FacebookId:        fb_id,
		FacebookAuthToken: fb_token,
		FirstName:         first_val,
		LastName:          last_val,
		Created:           time.Now(),
	}
	AddUser(db, &user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	r.ParseForm()
	username_val := r.FormValue("username")
	password_val := r.FormValue("password")
	num := CheckUser(db, username_val, password_val)

	if num == 0 {
		NotAuthedRedirect(w)
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

func LoginFacebookHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	r.ParseForm()
	fb_id_val := r.FormValue("fb_id")
	if fb_id_val == "" {
		NotAuthedRedirect(w)
		return
	}
	fmt.Println(fb_id_val)
	num := CheckUserFB(db, fb_id_val)
	if num == 0 {
		NotAuthedRedirect(w)
		return
	}
	kv := map[string]string{
		"uid": strconv.Itoa(int(num)),
	}
	setSession(kv, w)
	fmt.Fprintln(w, "logged in as", num)
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	if getSession(r)["uid"] == "" {
		NotAuthedRedirect(w)
	}
}

func NotAuthedRedirect(w http.ResponseWriter) {
	response := map[string]interface{}{"Status": http.StatusUnauthorized, "Message": "Login Failed"}
	j, _ := json.Marshal(response)
	http.Error(w, string(j), http.StatusUnauthorized)
}
