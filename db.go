package main

import (
	_ "fmt"
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
	Id          int64
	Title       string
	Category    PT_Category
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

func getDB() (*gorm.DB, error) {
	serv := os.Getenv("SERV")
	if serv == "local" {
		db, err := gorm.Open("postgres", "sslmode=disable")
		db.DB()
		db.SingularTable(true)
		return &db, err
		//return sql.Open("postgres", "sslmode=disable")
	}
	var e error
	return nil, e
}

func addUser(db *gorm.DB, username, password string) {
	user := PT_User{
		Username: username,
		Password: password,
		Created:  time.Now(),
	}
	db.Save(&user)

}

func checkUser(db *gorm.DB, username, password string) int64 {
	/*var uid int64
	err := db.QueryRow("SELECT uid from pt_user where name = $1 and password = $2", username, password).Scan(&uid)
	if err != nil {
		return -1
	}
	return uid*/
	return 0
}

func GetAllUsers(db *gorm.DB) {
	/*_, err := db.Query("SELECT * FROM pt_user")
	if err != nil {
		return
	}*/

}

func GetFeaturedUsers(db *gorm.DB) {
	//val, _ := db.Query("SELECT uid FROM prediction WHERE is_featured=TRUE")

}

func GetFeaturedPredictions(db *gorm.DB) {

}

func AddPrediction(db *gorm.DB, p PT_Prediction) {
	/*var pid int64
	err := db.QueryRow("INSERT INTO pt_prediction () VALUES() RETURNING id").Scan(&pid)
	if err != nil {
		return
	}*/
}
