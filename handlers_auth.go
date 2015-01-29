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
	var userMap map[string]string
	err := dec.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w)
		return
	}

	db, err := getDB()
	defer db.Close()
	if err != nil {
		return
	}
	var user PtUser
	user.Username = userMap["username"]
	user.Password = userMap["password"]

	user.Created = time.Now()
	AddUser(db, &user)
	fmt.Println("user added", user)
}

func RegisterFacebookHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var user PtUser
	err := dec.Decode(&user)
	user.Created = time.Now()

	db, err := getDB()
	defer db.Close()
	if err != nil {
		return
	}

	AddUser(db, &user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		DBError(w)
		return
	}
	userMap := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w)
		return
	}
	authedUser := CheckUser(db, userMap["username"], userMap["password"])

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
	defer db.Close()
	if err != nil {
		DBError(w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var userMap map[string]string
	err = decoder.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w)
		return
	}
	fmt.Println(userMap)
	var authedUser PtUser
	if userMap["facebookId"] == "" {
		authedUser = CheckUser(db, userMap["username"], userMap["password"])
	} else {
		authedUser = CheckUserFB(db, userMap["facebookId"])
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
	uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}
	db, _ := getDB()
	defer db.Close()
	var user PtUser
	db.First(&user, uid)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))

}

func JsonDecodeError(w http.ResponseWriter) {
	response := map[string]interface{}{"Status": http.StatusUnauthorized, "Message": "Json Decode Error"}
	j, _ := json.Marshal(response)
	http.Error(w, string(j), http.StatusBadRequest)
}

func DBError(w http.ResponseWriter) {
	response := map[string]interface{}{"Status": http.StatusUnauthorized, "Message": "Database Error"}
	j, _ := json.Marshal(response)
	http.Error(w, string(j), http.StatusConflict)
}

func NotAuthedRedirect(w http.ResponseWriter) {
	response := map[string]interface{}{"Status": http.StatusUnauthorized, "Message": "Not Authorized"}
	j, _ := json.Marshal(response)
	http.Error(w, string(j), http.StatusUnauthorized)
}

func GetUIDOrRedirect(w http.ResponseWriter, r *http.Request) int64 {
	voterIdStr := getSession(r)["uid"]
	if voterIdStr == "" {
		NotAuthedRedirect(w)
		return 0
	}
	voterId, _ := strconv.ParseInt(voterIdStr, 10, 64)
	return voterId
}
