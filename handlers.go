package main

import (
	"encoding/json"
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
	username_val = "USER2"
	password_val = "password2"
	user := PtUser{
		Username: username_val,
		Password: password_val,
		Created:  time.Now(),
	}
	AddUser(db, user)
	fmt.Println("user added", username_val, password_val)
}

//Marshal the data to json
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	users := GetAllUsers(db)
	j, _ := json.Marshal(users)
	fmt.Fprintln(w, string(j))
}

//Add error handling
func GetSingleUserHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	uid, _ := strconv.Atoi(vars["id"])
	user := GetUserByID(db, uid)
	j, _ := json.Marshal(user)
	fmt.Fprintln(w, string(j))
}

func GetAllPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	preds := GetAllPredictions(db)
	j, _ := json.Marshal(preds)
	fmt.Fprintln(w, string(j))
}

func GetSinglePredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	uid, _ := strconv.Atoi(vars["id"])
	prediction := GetPredictionByID(db, uid)
	j, _ := json.Marshal(prediction)
	fmt.Fprintln(w, string(j))
}

func GetLatestPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	//Get the 10 latest predictions
	preds := GetLatestPredictions(db, 10)
	j, _ := json.Marshal(preds)
	fmt.Fprintln(w, string(j))
}

func AddPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	title_to_add := "title"
	cId := int64(1)
	pred := PtPrediction{
		Title:     title_to_add,
		Created:   time.Now(),
		CreatorId: cId,
	}
	AddPrediction(db, &pred)
	fmt.Fprintln(w, "add prediction")
}

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	//Fill in real values here
	voterID := int64(1)
	predID := int64(1)
	vote := PtVote{
		VoterId:   voterID,
		VotedOnId: predID,
		Created:   time.Now(),
	}
	AddVote(db, &vote)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	r.ParseForm()
	username_val := r.FormValue("username")
	password_val := r.FormValue("password")
	u := PtUser{
		Username: username_val,
		Password: password_val,
	}

	LoginUser(db, &u)
	if u.Id == 0 {
		fmt.Fprintln(w, "failed log in")
		return
	}
	//Set up session or cookie

	//SetSession(strconv.Itoa(int(u.Id)), w)
	//u.Id is now set
	fmt.Fprintln(w, u.Id)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "logout")
}

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	categories := GetCategories(db)
	fmt.Fprintln(w, categories)
}

func GetSubcategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	categoryId := int64(1)
	subcats := GetSubcategoriesWithCategoryId(db, categoryId)
	fmt.Fprintln(w, subcats)
}
