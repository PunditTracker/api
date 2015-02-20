package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func NoInfoAtEndpointError(w http.ResponseWriter) {
	JsonError(w, http.StatusNotFound, "No data found")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	NotFoundError(w)
}

func GetDBOrPrintError(w http.ResponseWriter) *gorm.DB {
	db, err := getDB()
	if err != nil {
		DBError(w, err)
		return nil
	}
	return db
}

func DeadlinePassedError(w http.ResponseWriter) {
	JsonError(w, http.StatusConflict, "Can't vote on that prediction, the deadline has passed")
}

func PredictionGradedError(w http.ResponseWriter) {
	JsonError(w, http.StatusConflict, "Can't vote on that prediction, its already been graded")
}

func NotFoundError(w http.ResponseWriter) {
	JsonError(w, http.StatusNotFound, "Not Found")
}

func NoCredentialError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Not logged in as an admin.")
}

func NoUserWithEmailError(w http.ResponseWriter) {
	JsonError(w, http.StatusBadRequest, "No user exists with that email address.")
}

func NoIdIncludedError(w http.ResponseWriter) {
	JsonError(w, http.StatusBadRequest, "No ID provided.")
}

func UsernameDoesNotExistError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "That username does not exist.")
}

func IncorrectPasswordError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "That password is incorrect.")
}

func MultiVoteError(w http.ResponseWriter) {
	JsonError(w, http.StatusConflict, "Cannot vote on a prediction twice.")
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
	JsonError(w, http.StatusUnauthorized, "Not authorized.")
}

func MustResetPasswordError(w http.ResponseWriter) {
	JsonError(w, http.StatusExpectationFailed, "Due to a recent server change, your password has been reset. Please check your email for instructions and a link to choose your new password.")
}

func UserAlreadyExistsError(w http.ResponseWriter) {
	JsonError(w, http.StatusConflict, "That user already exists.")
}

func WrongOldPasswordError(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, "Old Password Incorrect")
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
