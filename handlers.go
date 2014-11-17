package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func GetFeaturedPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	fmt.Fprintln(w, GetFeaturedPredictions(db))
}

func GetFeaturedUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	users := GetFeaturedUsers(db)
	fmt.Fprintln(w, users)
}

//Get correct form values
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username_val := r.PostFormValue("username")
	password_val := r.PostFormValue("password")
	db, err := getDB()
	if err != nil {
		return
	}
	username_val = "USER"
	password_val = "password"
	user := PT_User{
		Username: username_val,
		Password: password_val,
		Created:  time.Now(),
	}
	addUser(db, user)
	fmt.Println("user added", username_val, password_val)
}

//Marshal the data to json
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	users := GetAllUsers(db)
	fmt.Fprintln(w, users)
}

//Add error handling
func GetSingleUserHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	uid, _ := strconv.Atoi(vars["id"])
	user := GetUserByID(db, uid)
	fmt.Fprintln(w, user)
}

func GetAllPredictionsHandler(w http.ResponseWriter, r *http.Request) {

}

func GetSinglePredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	uid, _ := strconv.Atoi(vars["id"])
	prediction := GetPredictionByID(db, uid)
	fmt.Fprintln(w, prediction)
}

func GetLatestPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	//db, _ := getDB()

}

func AddPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	title_to_add := "title"
	cId := int64(1)
	pred := PT_Prediction{
		Title:     title_to_add,
		Created:   time.Now(),
		CreatorId: cId,
	}
	AddPrediction(db, &pred)
}

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	//Fill in real values here
	vId := int64(1)
	voId := int64(1)
	vote := PT_Vote{
		VoterId:   vId,
		VotedOnId: voId,
		Created:   time.Now(),
	}
	AddVote(db, &vote)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}
