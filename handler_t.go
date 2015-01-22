package main

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func LoadTestDataHandler(w http.ResponseWriter, r *http.Request) {

	db, _ := getDB()
	AddBaseCategories(db)
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
		Title:     "X Will Win the Super Bowl",
		Tags: []PtTag{
			{
				Tag: "Sports",
			},
			{
				Tag: "NFL",
			},
		},
		Deadline: time.Now(),
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
	AddSubcategory(db, PtSubcategory{
		Name:        "NBA",
		ParentCatId: 1,
		IsLive:      false,
	})
	AddSubcategory(db, PtSubcategory{
		Name:        "NCAA Basketball",
		ParentCatId: 1,
		IsLive:      false,
	})
	AddSubcategory(db, PtSubcategory{
		Name:        "NFL",
		ParentCatId: 1,
		IsLive:      false,
	})
	AddSubcategory(db, PtSubcategory{
		Name:        "Supreme Court Decisions",
		ParentCatId: 2,
		IsLive:      false,
	})
	db.Save(&PtHero{
		IsLive:       true,
		Title:        "HERO CALLOUT",
		PredictionId: 1,
	})
	db.Save(&PtHero{
		IsLive:       true,
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
func AddSubcategory(db *gorm.DB, newSub PtSubcategory) int64 {
	db.Save(newSub)
	return newSub.Id
}
