package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	voterId := GetUIDOrRedirect(w, r)
	vars := mux.Vars(r)
	predID, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	vVal, _ := strconv.Atoi(vars["value"])

	db, err := getDB()
	if err != nil {
		DBError(w)
	}
	//Fill in real values here
	vote := PtVote{
		VoterId:   voterId,
		VotedOnId: predID,
		VoteValue: vVal,
		Created:   time.Now(),
	}
	AddVote(db, &vote)
}

func AverageForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	predId, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	db, err := getDB()
	if err != nil {
		DBError(w)
	}
	var avg float64
	ro := db.Debug().Raw("SELECT avg(vote_value) from (select vote_value FROM pt_vote where voted_on_id=?) as tab", predId).Row()
	ro.Scan(&avg)
	response := map[string]interface{}{
		"predictionId": predId,
		"average":      avg,
	}
	j, _ := json.Marshal(response)
	fmt.Fprintln(w, string(j))
}
