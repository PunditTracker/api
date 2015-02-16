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
	if IsAdminOrRedirect(w, r) {
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	predictionId, _ := strconv.ParseInt(vars["predId"], 10, 64)
	stateVal, _ := strconv.Atoi(vars["state"])
	newState := PtPredictionState(stateVal)

	prediction := SetState(db, predictionId, newState)

	if newState == DidHappen || newState == DidNotHappen {
		SetScoreForPrediction(db, predictionId, newState)
		SetPredictorScore(db, prediction.CreatorId, newState)
	}

	fmt.Fprintln(w, "state set", newState)
}

func SetHeroHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	//Parse the Json
	dec := json.NewDecoder(r.Body)
	var hero PtHero
	err := dec.Decode(&hero)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	SetHero(db, &hero)
	j, err := json.Marshal(hero)
	if err != nil {
		JsonDecodeError(w, err)
	}
	fmt.Fprintln(w, string(j))
}

func SetPredictionSetHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	//Parse the Json
	dec := json.NewDecoder(r.Body)
	var predictionSet PtPredictionSet
	err := dec.Decode(&predictionSet)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	SetPredictionSet(db, &predictionSet)
	j, err := json.Marshal(predictionSet)
	if err != nil {
		JsonDecodeError(w, err)
	}
	fmt.Fprintln(w, string(j))
}
