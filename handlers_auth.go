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
	dec := json.NewDecoder(r.Body)
	var user PtUser
	err := dec.Decode(&user)
	/*
			if username_val == "" ||
				password_val == "" ||
				firstname_val == "" ||
				lastname_val == "" {
				errMessage := "missing values:"
				if username_val == "" {
					errMessage += " username"
				}
				if password_val == "" {
					errMessage += " password"
				}
				if email_val == "" {
					errMessage += " email"
				}
				if firstname_val == "" {
					errMessage += " firstname"
				}
				if lastname_val == "" {
					errMessage += " lastname"
				}
				http.Error(w, errMessage, http.StatusBadRequest)
				return
			}
		user := PtUser{
			Username:  username_val,
			Password:  password_val,
			Email:     email_val,
			FirstName: firstname_val,
			LastName:  lastname_val,
			Created:   time.Now(),
		}
	*/
	db, err := getDB()
	if err != nil {
		return
	}

	user.Created = time.Now()
	AddUser(db, &user)
	fmt.Println("user added", user)
}

func RegisterFacebookHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var user PtUser
	err := dec.Decode(&user)
	user.Created = time.Now()
	/*
		if username_val == "" ||
			first_val == "" ||
			last_val == "" ||
			fb_token == "" ||
			email_val == "" {
			errMessage := "missing values"
			http.Error(w, errMessage, http.StatusBadRequest)
			return
		}
	*/

	db, err := getDB()
	if err != nil {
		return
	}

	AddUser(db, &user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		fmt.Fprintln(w, "db err:", err)
		return
	}
	var user PtUser
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "decode error:", err)
		return
	}
	authedUser := CheckUser(db, user.Username, user.Password)

	if authedUser.Id == 0 {
		NotAuthedRedirect(w)
		return
	}
	//Set up session or cookie
	kv := map[string]string{
		"uid": strconv.Itoa(int(authedUser.Id)),
	}
	setSession(kv, w)

	//num now set
	j, err := json.Marshal(authedUser)
	fmt.Fprintln(w, string(j))
}

func LoginFacebookHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		fmt.Fprintln(w, "db err: ", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var user PtUser
	err = decoder.Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "json decode error: ", err)
		return
	}
	fmt.Println(user)
	var authedUser PtUser
	if user.FacebookId == "" {
		authedUser = CheckUser(db, user.Username, user.Password)
	} else {
		authedUser = CheckUserFB(db, user.FacebookId)
	}

	//If not able to find user
	if authedUser.Id == 0 {
		NotAuthedRedirect(w)
		return
	}
	kv := map[string]string{
		"uid": strconv.Itoa(int(authedUser.Id)),
	}
	setSession(kv, w)
	j, err := json.Marshal(authedUser)
	fmt.Fprintln(w, string(j))
}

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	uid := getSession(r)["uid"]
	if uid == "" {
		NotAuthedRedirect(w)
	} else {
		fmt.Fprintln(w, "logged in as ", uid)
	}
}

func NotAuthedRedirect(w http.ResponseWriter) {
	response := map[string]interface{}{"Status": http.StatusUnauthorized, "Message": "Not Authorized"}
	j, _ := json.Marshal(response)
	http.Error(w, string(j), http.StatusUnauthorized)
}
