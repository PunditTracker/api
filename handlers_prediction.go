package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetFeaturedPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	fmt.Fprintln(w, GetFeaturedPredictions(db))
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

func GetPredictionsForSubcatHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	subCatId := int64(1)
	preds := GetPredictionsForSubcatId(db, subCatId)
	j, _ := json.Marshal(preds)
	fmt.Fprintln(w, string(j))
}

func StringToTsQuery(input string) string {
	toReturn := strings.Join(strings.Split(input, " "), " & ")
	return toReturn
}

func SearchPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	searchString := "test title football"
	searchString = StringToTsQuery(searchString)
	SearchPredictions(db, searchString)
}
