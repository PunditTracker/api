package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func SetStateHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	predictionId, _ := strconv.ParseInt(vars["predId"], 10, 64)
	stateVal, _ := strconv.Atoi(vars["state"])
	newState := PtPredictionState(stateVal)
	SetState(db, predictionId, newState)

	if newState == DidHappen || newState == DidNotHappen {
		SetScoreForPrediction(db, predictionId, newState)
	}

	fmt.Fprintln(w, "state set", newState)
}

func SetHeroHandler(w http.ResponseWriter, r *http.Request) {
	//Parse the Json
	dec := json.NewDecoder(r.Body)
	var hero PtHero
	err := dec.Decode(&hero)
	if err != nil {
		fmt.Println("Json Decode Error", err)
		return
	}

	db, _ := getDB()
	defer db.Close()
	SetHero(db, &hero)
	fmt.Fprintln(w, "hero set", hero)
}

func SetPredictionSetHandler(w http.ResponseWriter, r *http.Request) {
	//Parse the Json
	dec := json.NewDecoder(r.Body)
	var predictionSet PtPredictionSet
	err := dec.Decode(&predictionSet)
	if err != nil {
		fmt.Println("Json Decode Error", err)
		return
	}

	db, _ := getDB()
	defer db.Close()
	SetPredictionSet(db, &predictionSet)
	fmt.Fprintln(w, "hero set", predictionSet)
}
