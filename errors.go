package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetDBOrPrintError(w http.ResponseWriter) *gorm.DB {
	db, err := getDB()
	if err != nil {
		DBError(w)
		return nil
	}
	return db
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

func JsonDecodeError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Json Decode Error")
}

func DBError(w http.ResponseWriter) {
	JsonError(w, http.StatusConflict, "Database Error")
}

func NotAuthedRedirect(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Not Authorized")
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
