package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func GetDBOrPrintError(w http.ResponseWriter) *gorm.DB {
	db, err := getDB()
	if err != nil {
		DBError(w, err)
		return nil
	}
	return db
}

func NoCredentialError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Not logged in as admin")
}

func NoUserWithEmailError(w http.ResponseWriter) {
	JsonError(w, http.StatusBadRequest, "no user with that email")
}

func NoIdIncludedError(w http.ResponseWriter) {
	JsonError(w, http.StatusBadRequest, "No id provided")
}

func UsernameDoesNotExistError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Username does not exist")
}

func IncorrectPasswordError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Incorrect Password")
}

func JsonDecodeError(w http.ResponseWriter, err error) {
	JsonError(w, http.StatusUnauthorized, "Json Decode Error: "+err.Error())
	log.Println("jsonDecerr: " + err.Error())
}

func DBError(w http.ResponseWriter, err error) {
	JsonError(w, http.StatusConflict, "Database Error: "+err.Error())
	log.Println("dberr: " + err.Error())
}

func NotAuthedRedirect(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Not Authorized")
}

func MustResetPasswordError(w http.ResponseWriter) {
	JsonError(w, http.StatusExpectationFailed, "User must reset password")
}

func UserAlreadyExistsError(w http.ResponseWriter) {
	JsonError(w, http.StatusConflict, "User already exists")
}

func JsonError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	response := map[string]interface{}{"Status": status, "Message": message}
	j, err := json.Marshal(response)
	if err != nil {
		j = []byte("Json Failed")
	}
	fmt.Fprintln(w, string(j))
}
