package main

import (
	"net/http"
	"time"
)

func LoadTestDataHandler(w http.ResponseWriter, r *http.Request) {

	db, _ := getDB()
	AddUser(db, &PtUser{
		Username:  "ben",
		Password:  "pass",
		FirstName: "ben",
		LastName:  "levy",
		Email:     "emailthree",
		Created:   time.Now(),
	})
	AddUser(db, &PtUser{
		Username:  "howie",
		Password:  "password",
		FirstName: "howard",
		LastName:  "akumiah",
		Email:     "emailtwo",
		Created:   time.Now(),
	})
	AddUser(db, &PtUser{
		Username:          "jake",
		FirstName:         "jake",
		LastName:          "marsh",
		Email:             "emailone",
		FacebookId:        "fbid",
		FacebookAuthToken: "auth",
		Created:           time.Now(),
	})

	AddPrediction(db, &PtPrediction{
		CreatorId: 1,
		SubcatId:  0,
		Title:     "Prediction One",
		Deadline:  time.Now(),
	})

	AddPrediction(db, &PtPrediction{
		CreatorId: 2,
		SubcatId:  0,
		Title:     "Prediction Two",
		Deadline:  time.Now(),
	})
	AddPrediction(db, &PtPrediction{
		CreatorId: 3,
		SubcatId:  0,
		Title:     "Prediction Three",
		Deadline:  time.Now(),
	})

	AddVote(db, &PtVote{
		VoterId:   1,
		VotedOnId: 1,
		Created:   time.Now(),
	})
	AddVote(db, &PtVote{
		VoterId:   1,
		VotedOnId: 2,
		Created:   time.Now(),
	})
	AddVote(db, &PtVote{
		VoterId:   1,
		VotedOnId: 3,
		Created:   time.Now(),
	})

}
