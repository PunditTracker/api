package main

import (
	"net/http"
	"time"
)

func VoteForPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	//Fill in real values here
	voterID := int64(1)
	predID := int64(1)
	vote := PtVote{
		VoterId:   voterID,
		VotedOnId: predID,
		Created:   time.Now(),
	}
	AddVote(db, &vote)
}
