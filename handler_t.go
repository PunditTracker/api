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
}

func AddCategory(db *gorm.DB, newCategory PtCategory) int64 {
	db.Save(newCategory)
	return newCategory.Id
}
func AddSubcategory(db *gorm.DB, newSub PtSubcategory) int64 {
	db.Save(newSub)
	return newSub.Id
}
