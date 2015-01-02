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
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	prediction := GetPredictionByID(db, uid)
	j, _ := json.Marshal(prediction)
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

func GetLatestPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	//Get the 10 latest predictions
	preds := GetLatestPredictions(db, 10)
	j, _ := json.Marshal(preds)
	fmt.Fprintln(w, string(j))
}

func GetPredictionsForSubcatHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	subCatId, _ := strconv.ParseInt(vars["subcatid"], 10, 64)
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
	vars := mux.Vars(r)
	searchString := vars["searchstr"]
	searchString = StringToTsQuery(searchString)
	SearchPredictions(db, searchString)
}

func GetUserPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	predictions := GetUserPrediction(db, int64(uid))
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetTaggedPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	tag := vars["tag"]
	predictions := GetPredictionsForTag(db, tag)
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}
