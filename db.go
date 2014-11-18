package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type PT_User struct {
	Id                 int64
	Username           string
	Password           string
	Created            time.Time
	Score              int
	First_Name         string
	Last_Name          string
	Prediction_Graded  int
	Prediction_Correct int
	Avatar_URL         string
	Is_Pundit          bool
	Is_Featured        bool
	Predictions        []PT_Prediction
}

type PT_Category struct {
	Id            int64
	Name          string
	SubCategories []PT_SubCategory
	IsLive        bool
}

type PT_SubCategory struct {
	Id          int64
	Name        string
	ParentCat   PT_Category
	ParentCatId int64
	IsLive      bool
}

type PT_Prediction struct {
	Id    int64
	Title string
	//Category    PT_Category
	Is_Featured bool
	Created     time.Time
	Creator     PT_User
	CreatorId   int64
}

type PT_Vote struct {
	Id        int64
	Voter     PT_User
	VoterId   int64
	VotedOn   PT_Prediction
	VotedOnId int64
	Created   time.Time
}

func SetUpDB(db *gorm.DB) {
	db.AutoMigrate(
		&PT_User{},
		&PT_Category{},
		&PT_SubCategory{},
		&PT_Prediction{},
		&PT_Vote{},
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
		} else {
			fmt.Println("worked")
		}
		db.DB()
		db.SingularTable(true)
		return &db, err
	}
	var e error
	return nil, e
}

func addUser(db *gorm.DB, user PT_User) {
	db.Save(&user)

}

func checkUser(db *gorm.DB, username, password string) int64 {

	return 0
}

func GetUserByID(db *gorm.DB, uid int) PT_User {
	var user PT_User
	db.First(&user, uid)
	return user
}

func GetPredictionByID(db *gorm.DB, uid int) PT_Prediction {
	var pred PT_Prediction
	db.First(&pred, uid)
	return pred
}

func GetAllUsers(db *gorm.DB) []PT_User {
	users := []PT_User{}
	db.Find(&users)
	return users
}

func GetAllPredictions(db *gorm.DB) []PT_Prediction {
	preds := []PT_Prediction{}
	db.Find(&preds)
	return preds
}

func GetFeaturedUsers(db *gorm.DB) []PT_User {
	users := []PT_User{}
	db.Where(&PT_User{Is_Featured: true}).Find(&users)
	return users
}

func GetFeaturedPredictions(db *gorm.DB) []PT_Prediction {
	predictions := []PT_Prediction{}
	db.Where(&PT_Prediction{Is_Featured: true}).Find(&predictions)
	return predictions
}

func GetLatestPredictions(db *gorm.DB, x int) []PT_Prediction {
	predictions := []PT_Prediction{}
	db.Order("created").Limit(x).Find(&predictions)
	return predictions
}

func AddPrediction(db *gorm.DB, p *PT_Prediction) {
	db.Save(p)
}

func AddVote(db *gorm.DB, v *PT_Vote) {
	db.Save(v)
}

func LoginUser(db *gorm.DB, u *PT_User) {
	db.Where("username = ? and password = ?", u.Username, u.Password).First(u)
}
