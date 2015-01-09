package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func SetStateHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
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
