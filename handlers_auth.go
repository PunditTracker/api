package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sourcegraph/go-ses"
	"log"
	"net/http"
	"strconv"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	fmt.Fprintln(w, "logout")
}

func isStringAllNumbers(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var userMap map[string]string
	err := dec.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	var user PtUser
	user.Email = userMap["email"]
	user.Password = userMap["password"]
	user.ResetValidUntil = time.Now()
	user.Created = time.Now()
	SetPassword(db, &user)

	db.First(&user, user.Id)

	setSessionForUser(w, &user)

	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
	SendWelcomeEmail(&user)
}

func RegisterFacebookHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var user PtUser
	err := dec.Decode(&user)
	if err != nil {
		JsonDecodeError(w, err)
	}
	user.Created = time.Now()
	user.ResetValidUntil = time.Now()
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	SetPassword(db, &user)
	db.First(&user, user.Id)
	setSessionForUser(w, &user)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
	SendWelcomeEmail(&user)
}

func SendWelcomeEmail(user *PtUser) {
	fromEmail := "noreply@pundittracker.com"
	toEmail := user.Email

	_, err := ses.EnvConfig.SendEmail(
		fromEmail,
		toEmail,
		"Welcome to Pundit Tracker",
		"welcome!",
	)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Welcome message sent to:", user.Email)
	}
}

func setSessionForUser(w http.ResponseWriter, authedUser *PtUser) {
	//Set up session or cookie
	kv := map[string]string{
		"uid": strconv.Itoa(int(authedUser.Id)),
	}
	if authedUser.IsAdmin {
		kv["isadmin"] = "true"
	}
	setSession(kv, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	userMap := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	authedUser := CheckUser(db, userMap["email"], userMap["password"])

	if authedUser.Id == 0 {
		NotAuthedRedirect(w)
		return
	}
	setSessionForUser(w, &authedUser)

	//num now set
	j, err := json.Marshal(authedUser)
	fmt.Fprintln(w, string(j))
}

func LoginFacebookHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	decoder := json.NewDecoder(r.Body)
	var userMap map[string]string
	err := decoder.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	fmt.Println(userMap)
	var authedUser PtUser
	if userMap["facebookId"] == "" {
		authedUser = CheckUser(db, userMap["email"], userMap["password"])
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
	if authedUser.IsAdmin {
		kv["isadmin"] = "true"
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

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	var user PtUser
	db.First(&user, uid)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
}

func ForgotPasswordEndpoint(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var argMap map[string]string
	err := decoder.Decode(&argMap)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	fromEmail := "noreply@pundittracker.com"
	toEmail := argMap["email"]

	//Send an email to the user
	var user PtUser
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	db.Where("email=?", toEmail).First(&user)
	if user.Id == 0 {
		NoUserWithEmailError(w)
		return
	}

	user.ResetKey = uuid.New()
	user.ResetValidUntil = time.Now().Add(time.Hour)
	db.Save(&user)
	message := "Please goto " + fmt.Sprintf("foretellr.com/reset/%d/%s", user.Id, user.ResetKey)

	_, err = ses.EnvConfig.SendEmail(
		fromEmail,
		toEmail,
		"Password Recovery- Click link to reset",
		message,
	)
	fmt.Println("mess: ", message)
	if err == nil {
		j, _ := json.Marshal(map[string]string{"Message": "email sent"})
		fmt.Fprintln(w, string(j))
	} else {
		j, _ := json.Marshal(map[string]string{"Message": "email failed: " + err.Error()})
		fmt.Fprintln(w, string(j))
	}

}

func ResetPasswordEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		NoIdIncludedError(w)
	}
	resetKey := vars["resetKey"]
	decoder := json.NewDecoder(r.Body)
	var argMap map[string]string
	err = decoder.Decode(&argMap)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	newPassword := argMap["password"]
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	var user PtUser
	db.First(&user, uid)

	if time.Now().After(user.ResetValidUntil) {
		JsonError(w, http.StatusUnauthorized, "reset key expired")
		return
	}

	if user.ResetKey != resetKey {
		JsonError(w, http.StatusUnauthorized, "reset key incorrect")
		return
	}
	user.Password = newPassword
	SetPassword(db, &user)
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

//IsAdminOrRedirect returns true if it redirects
func IsAdminOrRedirect(w http.ResponseWriter, r *http.Request) bool {
	isAdmin := getSession(r)["isadmin"]
	if isAdmin != "true" {
		NotAuthedRedirect(w)
		return true
	}
	return false
}
