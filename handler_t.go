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
	SetPassword(db, &PtUser{
		Password:   "pass",
		FirstName:  "Ben",
		LastName:   "Levy",
		Email:      "bjlgds@gmail.com",
		Avatar_URL: "https://scontent-a-lga.xx.fbcdn.net/hphotos-xfp1/v/t1.0-9/10310651_10152412999714913_8367590480124742920_n.jpg?oh=ba1ebdc8ab067a410fb3f030bfbddd71&oe=555A0190",
		IsFeatured: true,
		Created:    time.Now(),
	})
	SetPassword(db, &PtUser{
		Password:   "pass",
		FirstName:  "Howard",
		LastName:   "Akumiah",
		Email:      "hakumiah@gmail.com",
		Avatar_URL: "https://fbcdn-sphotos-f-a.akamaihd.net/hphotos-ak-xap1/v/t1.0-9/10891425_10152970931527673_2036531108072244462_n.jpg?oh=f489e2f5f924e50278fa16ac754c3148&oe=5566E2EE&__gda__=1432879035_27f356f3d62060990e8ff694f8e64284",
		IsFeatured: true,
		Created:    time.Now(),
	})
	SetPassword(db, &PtUser{
		FirstName:  "Jake",
		LastName:   "Marsh",
		Email:      "jakemmarsh@gmail.com",
		Avatar_URL: "https://scontent-b-lga.xx.fbcdn.net/hphotos-xfp1/v/t1.0-9/10433831_10152797092583173_4991169485836557026_n.jpg?oh=a8527fd4974df690a81779e3c3384030&oe=55574156",
		FacebookId: "621883172",
		IsFeatured: true,
		Created:    time.Now(),
	})

	AddPrediction(db, &PtPrediction{
		CreatorId: 1,
		Title:     "X Will Win the Super Bowl",
		/*Tags: []PtTag{
			{
				Tag: "Sports",
			},
			{
				Tag: "NFL",
			},
		},*/
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
