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
	if checkEmailForExistence(w, db, user.Email) {
		return
	}
	user.FirstName = userMap["firstName"]
	user.LastName = userMap["lastName"]
	user.Password = userMap["password"]
	user.ResetValidUntil = time.Now()
	user.Created = time.Now()
	err = SetPasswordSalted(db, &user)
	if err != nil {
		DBError(w, err)
		//Password hash error
		return
	}

	err = SaveUser(db, &user)
	if err != nil {
		DBError(w, err)
		return
	}
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
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	if checkEmailForExistence(w, db, user.Email) {
		return
	}

	user.Created = time.Now()
	user.ResetValidUntil = time.Now()
	user.Password = "FB"
	err = SaveUser(db, &user)
	if err != nil {
		DBError(w, err)
		return
	}
	setSessionForUser(w, &user)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
	SendWelcomeEmail(&user)
}

/*
func testEmailHandler(w http.ResponseWriter, r *http.Request) {
	SendWelcomeEmail(&PtUser{Email: "bjlgds@gmail.com"})
}*/

var (
	WelcomeSubject = "It's time to play a part in tomorrow!"
	WelcomeEmail   = `Welcome to PunditTracker.com,

Thanks for joining us! At PunditTracker.com, it is our goal to bring accountability to the prediction industry. Take a look around, get familiar, and begin submitting your own predictions as soon as you're ready. Our system will begin tracking and scoring you as soon as you do.

Remember, #TomorrowMattersToday

Best,
Team PT`
)

func SendWelcomeEmail(user *PtUser) {
	fromEmail := "noreply@pundittracker.com"
	toEmail := user.Email

	_, err := ses.EnvConfig.SendEmail(
		fromEmail,
		toEmail,
		WelcomeSubject,
		WelcomeEmail,
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

	if checkEmailForNonePassword(w, db, userMap["email"]) {
		return
	}

	authedUser, err := CheckUser(db, userMap["email"], userMap["password"])
	if err != nil {
		if err.Error() == "no user" {
			NoUserWithEmailError(w)
			return
		}
		if err.Error() == "wrong password" {
			IncorrectPasswordError(w)
			return
		}
	}

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
	if checkEmailForNonePassword(w, db, userMap["email"]) {
		return
	}
	var authedUser PtUser
	log.Println("attempt login with usermap:", userMap)
	if userMap["facebookId"] == "" {
		authedUser, err = CheckUser(db, userMap["email"], userMap["password"])
		if err != nil {
			if err.Error() == "no user" {
				NoUserWithEmailError(w)
				return
			}
			if err.Error() == "wrong password" {
				IncorrectPasswordError(w)
				return
			}
		}
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
	toEmail := argMap["email"]

	//Send an email to the user
	ForgotPassword(w, toEmail)
}

func ForgotPassword(w http.ResponseWriter, toEmail string) {
	log.Println("Forgot password begin")
	fromEmail := "noreply@pundittracker.com"
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	var user PtUser
	db.Where("email=?", toEmail).First(&user)
	log.Println(user)
	if user.Id == 0 {
		NoUserWithEmailError(w)
		return
	}

	user.ResetKey.String = uuid.New()
	user.ResetKey.Valid = true
	user.ResetValidUntil = time.Now().Add(time.Hour)
	db.Save(&user)
	link := fmt.Sprintf("pundittracker.com/reset/%d/%s", user.Id, user.ResetKey.String)
	message := "Please goto " + link

	_, err := ses.EnvConfig.SendEmail(
		fromEmail,
		toEmail,
		"Password Recovery- Click link to reset",
		message,
	)
	log.Println("mess: ", message)
	if err != nil {
		j, _ := json.Marshal(map[string]string{"Message": "email failed: " + err.Error()})
		fmt.Fprintln(w, string(j))
	}
}

func ResetPasswordEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		NoIdIncludedError(w)
		return
	}
	if uid == 0 {
		NoIdIncludedError(w)
		return
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
	if !user.ResetKey.Valid {
		JsonError(w, http.StatusUnauthorized, "no reset key set")
		return
	}

	if user.ResetKey.String != resetKey {
		JsonError(w, http.StatusUnauthorized, "reset key incorrect")
		return
	}
	user.Password = newPassword
	SetPasswordSalted(db, &user)
	UpdatePassword(db, &user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	var user PtUser
	db.First(&user, uid)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		JsonDecodeError(w, err)
	}

	defer db.Close()
	db.Model(PtUser{Id: uid}).Update(user)

	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}
	var userMap map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}

	//Check if old password is correct
	if userMap["newPassword"] == "" || userMap["oldPassword"] == "" {
		return
	}

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	user, err := CheckUserWithIdAndPass(db, uid, userMap["oldPassword"])
	if err != nil {
		WrongOldPasswordError(w)
		return
	}
	db.First(&user, user.Id)
	user.Password = userMap["newPassword"]
	//salt new password and update
	err = SetPasswordSalted(db, &user)
	if err != nil {
		log.Println("salt pass err", err.Error())
		return
	}
	err = UpdatePassword(db, &user)
	if err != nil {
		log.Println("change pass error", err.Error())
		return
	}
	j, _ := json.Marshal(map[string]string{
		"Message": "Password Changed",
	})
	fmt.Fprintln(w, string(j))
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

func GetUIDOrZero(r *http.Request) int64 {
	voterIdStr := getSession(r)["uid"]
	if voterIdStr == "" {
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
