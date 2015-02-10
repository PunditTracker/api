package main

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func init() {
	db, err := getDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	SetUpDB(db)
}

func getDB() (*gorm.DB, error) {
	serv := os.Getenv("SERV")
	if serv == "local" {
		db, err := gorm.Open("postgres", "sslmode=disable")
		db.DB()
		db.SingularTable(true)
		return &db, err
	}
	if serv == "aws" {
		db, err := gorm.Open("postgres", "host=ptdev.ccm2e8gfsxjt.us-west-2.rds.amazonaws.com dbname=ptdev user=pundittracker password=ptrack20!!")
		if err != nil {
			log.Println(err)
		}
		db.DB()
		db.SingularTable(true)
		db.LogMode(false)
		return &db, err
	}
	return nil, errors.New("No SERV specified")
}

type PtUser struct {
	Id                int64
	Password          string         `sql:"not null" json:"-"`
	ResetKey          sql.NullString `json:"-" sql:"DEFAULT:null"`
	ResetValidUntil   time.Time      `json:"-" sql:"DEFAULT:current_timestamp"`
	Email             string         `sql:"not null; unique"`
	Created           time.Time      `sql:"not null; DEFAULT:current_timestamp"`
	Score             int            `sql:"not null; DEFAULT:0"`
	PredictionGraded  int            `sql:"not null; DEFAULT:0"`
	PredictionCorrect int            `sql:"not null; DEFAULT:0"`
	IsAdmin           bool           `sql:"not null; DEFAULT:FALSE"`
	IsPundit          bool           `sql:"not null; DEFAULT:FALSE"`
	IsFeatured        bool           `sql:"not null; DEFAULT:FALSE"`
	FacebookId        string         `sql:"not null; DEFAULT:''"`
	FacebookAuthToken string         `sql:"not null; DEFAULT:''"`
	FirstName         string         `sql:"not null; DEFAULT:''"`
	LastName          string         `sql:"not null; DEFAULT:''"`
	Avatar_URL        string         `sql:"not null; DEFAULT:''"`
	Location          string         `sql:"not null; DEFAULT:''"`
	Predictions       []PtPrediction
}

type PtCategory struct {
	Id     int64
	Name   string `sql:"not null; unique"`
	IsLive bool   `sql:"not null; DEFAULT:FALSE"`
}

type PtPredictionState int

const (
	InFuture     PtPredictionState = 0
	Ungraded                       = 1
	DidHappen                      = 2
	DidNotHappen                   = 3
)

type PtPrediction struct {
	Id          int64
	CreatorId   int64             `sql:"not null"`
	CategoryId  int64             `sql:"not null"`
	Title       string            `sql:"not null"`
	State       PtPredictionState `sql:"not null";DEFAULT:0`
	IsFeatured  bool              `sql:"not null; DEFAULT:FALSE"`
	Created     time.Time         `sql:"not null; DEFAULT:current_timestamp"`
	Deadline    time.Time
	ImageUrl    string
	Creator     PtUser
	Category    PtCategory `json:"-"`
	Tags        []string   `sql:"-"`
	CurUserVote int        `sql:"-"`
	TagVal      []PtTag    `gorm:"many2many:prediction_tag_map;"`
}

type PtVote struct {
	Id            int64
	VoterId       int64     `sql:"not null"`
	VotedOnId     int64     `sql:"not null"`
	AverageAtTime float64   `sql:"not null"`
	VoteValue     int       `sql:"not null"`
	Created       time.Time `sql:"not null; DEFAULT:current_timestamp"`
	Voter         PtUser
	VotedOn       PtPrediction `json:"-"`
}

type PtTag struct {
	Id  int64
	Tag string `sql:"not null; unique"`
}

type PtBracket struct {
	Id        int64
	CreatorId int64     `sql:"not null"`
	Created   time.Time `sql:"not null; DEFAULT:current_timestamp"`

	FirstRound0  string
	FirstRound1  string
	FirstRound2  string
	FirstRound3  string
	FirstRound4  string
	FirstRound5  string
	FirstRound6  string
	FirstRound7  string
	FirstRound8  string
	FirstRound9  string
	FirstRound10 string
	FirstRound11 string
	FirstRound12 string
	FirstRound13 string
	FirstRound14 string
	FirstRound15 string

	SecondRound0 string
	SecondRound1 string
	SecondRound2 string
	SecondRound3 string
	SecondRound4 string
	SecondRound5 string
	SecondRound6 string
	SecondRound7 string

	ThirdRound0 string
	ThirdRound1 string
	ThirdRound2 string
	ThirdRound3 string

	FourthRound0 string
	FourthRound1 string

	FifthRound0 string
}

type PtPredictionSet struct {
	Id            int64
	IsLive        bool   `sql:"not null; DEFAULT:FALSE"`
	Title         string `sql:"not null"`
	ImageUrl      string `sql:"not null"`
	Prediction1Id int64  `sql:"not null"`
	Prediction2Id int64  `sql:"not null"`
	Prediction3Id int64  `sql:"not null"`
	Prediction1   PtPrediction
	Prediction2   PtPrediction
	Prediction3   PtPrediction
}

type PtHero struct {
	Id           int64
	LocationNum  int64  `sql:"not null"`
	IsLive       bool   `sql:"not null; DEFAULT:FALSE"`
	ImageUrl     string `sql:"not null"`
	Title        string `sql:"not null"`
	PredictionId int64
	ButtonText   string
	ButtonUrl    string
}

func SetUpDB(db *gorm.DB) {
	db.Debug().AutoMigrate(
		&PtUser{},
		&PtCategory{},
		&PtPrediction{},
		&PtVote{},
		&PtTag{},
		&PtBracket{},
		&PtPredictionSet{},
		&PtHero{},
	)
}

func (p *PtPrediction) AfterFind(tx *gorm.DB) {
	tx.First(&p.Creator, p.CreatorId)
	p.Tags = GetTagsForPrediction(tx, p.Id)
}
