package main

import (
	"fmt"
	_ "github.com/gorilla/mux"
	"net/http"
)

func GetFeaturedPredictionsHandler(w http.ResponseWriter, r *http.Request) {

}

func GetFeaturedUsersHandler(w http.ResponseWriter, r *http.Request) {

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username_val := r.PostFormValue("username")
	password_val := r.PostFormValue("password")
	/*if username_val == "" || password_val == "" {
		http.Redirect(w, r, "AddUserFail.html", http.StatusFound)
	}*/
	db, err := getDB()
	if err != nil {
		return
	}
	username_val = "USER"
	password_val = "password"
	addUser(db, username_val, password_val)
	fmt.Println("user added", username_val, password_val)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	GetAllUsers(db)
}

func GetSingleUserHandler(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	_ := vars["id"]*/
}

func GetAllPredictionsHandler(w http.ResponseWriter, r *http.Request) {

}

func GetSinglePredictionHandler(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	_ := vars["id"]*/
}

func AddPredictionHandler(w http.ResponseWriter, r *http.Request) {

}

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}
