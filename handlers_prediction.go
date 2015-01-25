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
	defer db.Close()
	predictions := GetFeaturedPredictions(db)
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetAllPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	preds := GetAllPredictions(db)
	j, _ := json.Marshal(preds)
	fmt.Fprintln(w, string(j))
}

func GetSinglePredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	prediction := GetPredictionByID(db, uid)
	j, _ := json.Marshal(prediction)
	fmt.Fprintln(w, string(j))
}

func AddPredictionHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var prediction PtPrediction
	err := dec.Decode(&prediction)
	if err != nil {
		fmt.Println("Json Decode Error", err)
		return
	}
	prediction.CreatorId = GetUIDOrRedirect(w, r)
	prediction.Created = time.Now()
	db, _ := getDB()
	defer db.Close()
	AddPrediction(db, &prediction)
	j, _ := json.Marshal(prediction)
	fmt.Fprintln(w, string(j))
}

func GetLatestPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	//Get the 10 latest predictions
	preds := GetLatestPredictions(db, 10)
	j, _ := json.Marshal(preds)
	fmt.Fprintln(w, string(j))
}

func GetPredictionsForSubcatHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	subCatId, _ := strconv.ParseInt(vars["subcatid"], 10, 64)
	preds := GetPredictionsForSubcatId(db, subCatId)
	j, _ := json.Marshal(preds)
	fmt.Fprintln(w, string(j))
}

func StringToTsQuery(input string, connector string) string {
	toReturn := strings.Join(strings.Split(input, " "), connector)
	return toReturn
}

func SearchPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	searchString := vars["searchstr"]
	searchString = StringToTsQuery(searchString, " & ")
	SearchPredictions(db, searchString)
}

func GetUserPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	predictions := GetUserPrediction(db, int64(uid))
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetTaggedPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	tag := vars["tag"]
	predictions := GetPredictionsForTag(db, tag)
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetHeroPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	heros := GetLivePtHeros(db)
	j, _ := json.Marshal(heros)
	fmt.Fprintln(w, string(j))
}

func GetPredictionSetHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	predictionSets := GetLivePredictionSets(db)
	j, _ := json.Marshal(predictionSets)
	fmt.Fprintln(w, string(j))
}
