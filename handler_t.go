package main

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func LoadTestDataHandler(w http.ResponseWriter, r *http.Request) {

	db, _ := getDB()
	defer db.Close()
	AddBaseCategories(db)

	AddPrediction(db, &PtPrediction{
		CreatorId: 1,
		Title:     "X Will Win the Super Bowl",
		TagVal: []PtTag{
			{
				Tag: "Sports",
			},
			{
				Tag: "NFL",
			},
		},
		IsFeatured: true,
		Deadline:   time.Now(),
	})

	AddPrediction(db, &PtPrediction{
		CreatorId:  2,
		Title:      "Prediction Two",
		IsFeatured: true,
		Deadline:   time.Now(),
	})
	AddPrediction(db, &PtPrediction{
		CreatorId:  3,
		Title:      "Prediction Three",
		IsFeatured: true,
		Deadline:   time.Now(),
	})

	AddVote(db, &PtVote{
		VoterId:   1,
		VotedOnId: 1,
		VoteValue: 1,
		Created:   time.Now(),
	})
	AddVote(db, &PtVote{
		VoterId:   1,
		VotedOnId: 1,
		VoteValue: 2,
		Created:   time.Now(),
	})
	AddVote(db, &PtVote{
		VoterId:   1,
		VotedOnId: 3,
		VoteValue: 1,
		Created:   time.Now(),
	})

}

func AddBaseCategories(db *gorm.DB) {
	AddCategory(db, PtCategory{
		Name:   "Sports",
		IsLive: false,
	})
	AddCategory(db, PtCategory{
		Name:   "Politics",
		IsLive: false,
	})
	AddCategory(db, PtCategory{
		Name:   "Economics",
		IsLive: false,
	})
	db.Save(&PtHero{
		IsLive:       true,
		LocationNum:  2,
		Title:        "HERO CALLOUT",
		PredictionId: 1,
	})
	db.Save(&PtHero{
		IsLive:       true,
		LocationNum:  3,
		Title:        "HERO CALLOUT2",
		PredictionId: 2,
	})
	db.Save(&PtPredictionSet{
		IsLive:        true,
		Title:         "SET ONE",
		ImageUrl:      "something.jpg",
		Prediction1Id: 1,
		Prediction2Id: 2,
		Prediction3Id: 3,
	})
}

func AddHero(db *gorm.DB, newHero *PtHero) int64 {
	db.Save(newHero)
	return newHero.Id
}
func AddPredictionSet(db *gorm.DB, newSet *PtPredictionSet) int64 {
	db.Save(newSet)
	return newSet.Id
}

func AddCategory(db *gorm.DB, newCategory PtCategory) int64 {
	db.Save(newCategory)
	return newCategory.Id
}
