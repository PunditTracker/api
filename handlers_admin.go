package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
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
		//SetScoreForVotes(db, predictionId, newState)
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
	j, _ := json.Marshal(predictionSet)
	fmt.Fprintln(w, string(j))
}

func SetPredictionLocationHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	dec := json.NewDecoder(r.Body)
	var predictionLoc PtPredictionLocation
	err := dec.Decode(&predictionLoc)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	SetPredictionLocation(db, &predictionLoc)
	j, _ := json.Marshal(predictionLoc)
	fmt.Fprintln(w, string(j))
}

func GetPredictionLocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cat_id, err := strconv.ParseInt(vars["cat_id"], 10, 64)
	if err != nil {
		cat_id = 0
		log.Println("cat_id err", err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	db = db.Debug()
	var locs []PtPredictionLocation
	db.Where("category_id = ?", cat_id).Order("location_num").Find(&locs)
	j, _ := json.Marshal(locs)
	fmt.Fprintln(w, string(j))
}
