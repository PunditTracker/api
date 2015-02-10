package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func GetVoteHandler(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	vars := mux.Vars(r)
	predId, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	v := GetVote(db, uid, predId)
	j, _ := json.Marshal(v)
	fmt.Fprintln(w, string(j))
}

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	voterId := GetUIDOrRedirect(w, r)
	vars := mux.Vars(r)
	predId, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	vVal, _ := strconv.Atoi(vars["value"])

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	avg := GetAverageVoteForPredictionId(db, predId)
	vote := PtVote{
		VoterId:       voterId,
		VotedOnId:     predId,
		VoteValue:     vVal,
		AverageAtTime: avg,
		Created:       time.Now(),
	}
	AddVote(db, &vote)
}

func AverageForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	predId, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	avg := GetAverageVoteForPredictionId(db, predId)
	response := map[string]interface{}{
		"predictionId": predId,
		"average":      avg,
	}
	j, _ := json.Marshal(response)
	fmt.Fprintln(w, string(j))
}
