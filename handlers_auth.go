package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
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
		JsonDecodeError(w)
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

	user.Created = time.Now()
	SetPassword(db, &user)

	db.First(&user, user.Id)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
	SendWelcomeEmail(&user)
}

func RegisterFacebookHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var user PtUser
	err := dec.Decode(&user)
	if err != nil {
		JsonDecodeError(w)
	}
	user.Created = time.Now()

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	SetPassword(db, &user)
	db.First(&user, user.Id)
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
		JsonDecodeError(w)
		return
	}
	authedUser := CheckUser(db, userMap["email"], userMap["password"])

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
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	decoder := json.NewDecoder(r.Body)
	var userMap map[string]string
	err := decoder.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w)
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
		JsonDecodeError(w)
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

	user.ResetKey = uuid.New()
	user.ResetValidUntil = time.Now().Add(time.Hour)
	db.Save(&user)

	_, err = ses.EnvConfig.SendEmail(
		fromEmail,
		toEmail,
		"Password Recovery- Click link to reset",
		"Please goto "+fmt.Sprintf("foretellr.com/v1/auth/%s", user.ResetKey),
	)
	if err == nil {
		//Sent email
	} else {
		// Error sending email
	}

}

func ResetPasswordEndpoint(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var argMap map[string]string
	err := decoder.Decode(&argMap)
	if err != nil {
		JsonDecodeError(w)
		return
	}
	resetKey := argMap["resetKey"]
	newPassword := argMap["newpass"]
	uid := 1
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
