package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetVotesForUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("start get vote handler")
	vars := mux.Vars(r)
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	votes := GetVotesForUser(db, uid)
	if votes == nil {
		NoInfoAtEndpointError(w)
		return
	}
	log.Println("votes:", votes)
	j, _ := json.Marshal(votes)
	fmt.Fprintln(w, string(j))
}

func GetVoteHandler(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	vars := mux.Vars(r)
	predId, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	pred := PtPrediction{Id: predId}
	UpdateVoteValue(db, uid, &pred)
	v := pred.CurUserVote
	j, _ := json.Marshal(v)
	fmt.Fprintln(w, string(j))
}

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	voterId := GetUIDOrRedirect(w, r)
	if voterId == 0 {
		return
	}
	vars := mux.Vars(r)
	predId, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	vVal, _ := strconv.Atoi(vars["value"])

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	if PredictionDeadlinePassed(db, predId) {
		DeadlinePassedError(w)
		return
	}

	if VoteExists(db, voterId, predId) {
		MultiVoteError(w)
		return
	}

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
