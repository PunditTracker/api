package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	voterIdStr := getSession(r)["uid"]
	if voterIdStr == "" {
		NotAuthedRedirect(w)
		return
	}
	voterId, _ := strconv.ParseInt(voterIdStr, 10, 64)
	vars := mux.Vars(r)
	predID, _ := strconv.ParseInt(vars["pred_id"], 10, 64)
	vVal, _ := strconv.Atoi(vars["value"])

	db, _ := getDB()
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
	db, _ := getDB()
	var avg float64
	ro := db.Debug().Raw("SELECT avg(vote_value) from (select vote_value FROM pt_vote where voted_on_id=?) as tab", predId).Row()
	ro.Scan(&avg)
	fmt.Fprintln(w, predId, avg)
}
