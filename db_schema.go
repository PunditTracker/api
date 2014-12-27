package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
)

type PtUser struct {
	Id                int64
	Username          string    `sql:"not null; unique"`
	Password          string    `sql:"not null"`
	Email             string    `sql:"not null; unique"`
	Created           time.Time `sql:"not null; DEFAULT:current_timestamp"`
	Score             int       `sql:"not null; DEFAULT:0"`
	PredictionGraded  int       `sql:"not null; DEFAULT:0"`
	PredictionCorrect int       `sql:"not null; DEFAULT:0"`
	IsPundit          bool      `sql:"not null; DEFAULT:FALSE"`
	IsFeatured        bool      `sql:"not null; DEFAULT:FALSE"`
	FacebookId        string
	FacebookAuthToken string //`sql:"unique"`?
	FirstName         string
	LastName          string
	Avatar_URL        string
	Predictions       []PtPrediction
}

type PtCategory struct {
	Id            int64
	Name          string
	Subcategories []PtSubcategory
	IsLive        bool `sql:"not null; DEFAULT:FALSE"`
}

type PtSubcategory struct {
	Id          int64
	Name        string
	ParentCat   PtCategory
	ParentCatId int64 `sql:"not null"`
	IsLive      bool  `sql:"not null; DEFAULT: FALSE"`
	Predictions []PtPrediction
}

type PtPrediction struct {
	Id         int64
	CreatorId  int64     `sql:"not null"`
	SubcatId   int64     `sql:"not null"`
	Title      string    `sql:"not null"`
	IsFeatured bool      `sql:"not null; DEFAULT:FALSE"`
	Created    time.Time `sql:"not null; DEFAULT:current_timestamp"`
	Deadline   time.Time `sql:"not null"`
	Creator    PtUser
	Subcat     PtSubcategory
	Tags       []PtTag `gorm:"many2many:prediction_tag_map;"`
}

type PtVote struct {
	Id        int64
	VoterId   int64     `sql:"not null"`
	VotedOnId int64     `sql:"not null"`
	VotedFor  bool      `sql:"not null"`
	Created   time.Time `sql:"not null; DEFAULT:current_timestamp"`
	Voter     PtUser
	VotedOn   PtPrediction
}

type PtTag struct {
	Id  int64
	Tag string `sql:"not null"`
}

func SetUpDB(db *gorm.DB) {
	db.AutoMigrate(
		&PtUser{},
		&PtCategory{},
		&PtSubcategory{},
		&PtPrediction{},
		&PtVote{},
		&PtTag{},
	)
}
