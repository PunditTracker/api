package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type PtUser struct {
	Id                int64
	Username          string    `sql:"not null"`
	Password          string    `sql:"not null"`
	Created           time.Time `sql:"not null; DEFAULT:current_timestamp"`
	Score             int       `sql:"not null; DEFAULT:0"`
	PredictionGraded  int       `sql:"not null; DEFAULT:0"`
	PredictionCorrect int       `sql:"not null; DEFAULT:0"`
	IsPundit          bool      `sql:"not null; DEFAULT:FALSE"`
	IsFeatured        bool      `sql:"not null; DEFAULT:FALSE"`
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
	IsLive      bool  `sql:"not null"; DEFAULT: FALSE`
	Predictions []PtPrediction
}

type PtPrediction struct {
	Id    int64
	Title string
	//Category    PT_Category
	IsFeatured bool      `sql:"not null; DEFAULT:FALSE"`
	Created    time.Time `sql:"not null; DEFAULT:current_timestamp"`
	CreatorId  int64     `sql:"not null"`
	SubcatId   int64     `sql:"not null"`
	Creator    PtUser
	Subcat     PtSubcategory
	Deadline   time.Time
}

type PtVote struct {
	Id        int64
	VoterId   int64     `sql:"not null"`
	VotedOnId int64     `sql:"not null"`
	Created   time.Time `sql:"not null; DEFAULT:current_timestamp"`
	Voter     PtUser
	VotedOn   PtPrediction
}

func SetUpDB(db *gorm.DB) {
	db.AutoMigrate(
		&PtUser{},
		&PtCategory{},
		&PtSubcategory{},
		&PtPrediction{},
		&PtVote{},
	)
}

var (
	DBID       = "ptdev"
	DBUSERNAME = "pundittracker"
	DBPASSWORD = "ptrack20!!"
)

func getDB() (*gorm.DB, error) {
	serv := os.Getenv("SERV")
	if serv == "local" {
		db, err := gorm.Open("postgres", "sslmode=disable")
		db.DB()
		db.SingularTable(true)
		return &db, err
		//return sql.Open("postgres", "sslmode=disable")
	}
	if serv == "aws" {
		db, err := gorm.Open("postgres", "host=ptdev.ccm2e8gfsxjt.us-west-2.rds.amazonaws.com dbname=ptdev user=pundittracker password=ptrack20!!")
		if err != nil {
			fmt.Println(err)
		}
		db.DB()
		db.SingularTable(true)
		return &db, err
	}
	var e error
	return nil, e
}

func AddUser(db *gorm.DB, user PtUser) {
	db.Save(&user)
}

func CheckUser(db *gorm.DB, username, password string) int64 {

	return 0
}

func GetUserByID(db *gorm.DB, uid int) PtUser {
	var user PtUser
	db.First(&user, uid)
	return user
}

func GetPredictionByID(db *gorm.DB, uid int) PtPrediction {
	var pred PtPrediction
	db.First(&pred, uid)
	return pred
}

func GetAllUsers(db *gorm.DB) []PtUser {
	users := []PtUser{}
	db.Find(&users)
	return users
}

func GetAllPredictions(db *gorm.DB) []PtPrediction {
	preds := []PtPrediction{}
	db.Find(&preds)
	return preds
}

func GetFeaturedUsers(db *gorm.DB) []PtUser {
	users := []PtUser{}
	db.Where(&PtUser{IsFeatured: true}).Find(&users)
	return users
}

func GetFeaturedPredictions(db *gorm.DB) []PtPrediction {
	predictions := []PtPrediction{}
	db.Where(&PtPrediction{IsFeatured: true}).Find(&predictions)
	return predictions
}

func GetLatestPredictions(db *gorm.DB, x int) []PtPrediction {
	predictions := []PtPrediction{}
	db.Order("created").Limit(x).Find(&predictions)
	return predictions
}

func AddPrediction(db *gorm.DB, p *PtPrediction) {
	db.Save(p)
}

func AddVote(db *gorm.DB, v *PtVote) {
	db.Save(v)
}

func LoginUser(db *gorm.DB, u *PtUser) {
	db.Where("username = ? and password = ?", u.Username, u.Password).First(u)
}

//Category Stuff
func GetCategories(db *gorm.DB) []PtCategory {
	categories := []PtCategory{}
	db.Find(&categories)
	return categories
}

func GetSubcategoriesWithCategoryId(db *gorm.DB, catId int64) []PtSubcategory {
	subcats := []PtSubcategory{}
	db.Where("parent_cat_id = ?", catId).Find(&subcats)
	return subcats
}

func GetSubcategoriesWithCategoryName(db *gorm.DB, name string) []PtSubcategory {
	/*subcats := []PtSubcategory{}

	db.Where("parent_cat_id = ?", uid).Find(&subcats)
	return subcats*/
	return []PtSubcategory{}
}

func GetPredictionsForSubcatId(db *gorm.DB, subcatId int64) []PtPrediction {
	preds := []PtPrediction{}
	db.Where("subcat_id = ?", subcatId).Find(&preds)
	return preds
}
