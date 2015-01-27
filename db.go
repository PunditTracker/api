package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	DBID       = "ptdev"
	DBUSERNAME = "pundittracker"
	DBPASSWORD = "ptrack20!!"
	db_logger  *log.Logger
)

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
			fmt.Println(err)
		}
		db.DB()
		db.SingularTable(true)
		db.LogMode(true)
		db.SetLogger(db_logger)
		return &db, err
	}
	return nil, errors.New("No SERV specified")
}

func AddUser(db *gorm.DB, user *PtUser) error {
	passByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passByte)
	db.Save(user)
	return nil
}

func CheckUser(db *gorm.DB, username, password string) PtUser {
	var user PtUser
	db.Where("username = ?", username).First(&user)
	hashedPass := []byte(user.Password)
	e := bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	//Password accepted
	if e == nil {
		return user
	} else {
		var notUser PtUser
		return notUser
	}
}

func CheckUserFB(db *gorm.DB, fb_id string) PtUser {
	var user PtUser
	db.Where("facebook_id = ?", fb_id).First(&user)
	return user
}

func GetUserByID(db *gorm.DB, uid int64) PtUser {
	var user PtUser
	db.First(&user, uid)
	return user
}

func GetUserPrediction(db *gorm.DB, uid int64) []PtPrediction {
	var preds []PtPrediction
	db.Where("creator_id = ?", uid).Find(&preds)
	return preds
}

func GetPredictionByID(db *gorm.DB, uid int64) PtPrediction {
	pred := PtPrediction{}
	db.Where("id = ?", uid).First(&pred)
	return pred
}

func GetAllUsers(db *gorm.DB) []PtUser {
	users := []PtUser{}
	db.Find(&users)
	return users
}

func GetAllPredictions(db *gorm.DB) []PtPrediction {
	var preds []PtPrediction
	db.Find(&preds)
	return preds
}

func GetFeaturedUsers(db *gorm.DB) []PtUser {
	var users []PtUser
	db.Where(&PtUser{IsFeatured: true}).Find(&users)
	return users
}

func GetFeaturedPredictions(db *gorm.DB, l int) []PtPrediction {
	var predictions []PtPrediction
	db.Where(&PtPrediction{IsFeatured: true}).Limit(l).Find(&predictions)
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
	fmt.Println("category id")
	db.Where("parent_cat_id = ?", catId).Find(&subcats)
	return subcats
}

func GetSubcategoriesWithCategoryName(db *gorm.DB, name string) []PtSubcategory {
	var category PtCategory
	var subcats []PtSubcategory
	fmt.Println("category name")
	db.Where("name = ?", name).First(&category)
	if category.Id == 0 {
		return subcats
	}
	db.Where("parent_cat_id = ?", category.Id).Find(&subcats)
	return subcats
}

func GetPredictionsForSubcatId(db *gorm.DB, subcatId int64) []PtPrediction {
	preds := []PtPrediction{}
	db.Where("subcat_id = ?", subcatId).Find(&preds)
	return preds
}

func GetPredictionsForCategoryId(db *gorm.DB, catId int64) []PtPrediction {
	preds := []PtPrediction{}
	db.Where("category_id = ?", catId).Find(&preds)
	return preds
}

func SearchPredictions(db *gorm.DB, searchString string) []int64 {
	toReturn := []int64{}
	rows, err := db.Raw(`SELECT pid
			FROM (SELECT pred.id as pid,
				  to_tsvector(pred.title) as document
			FROM pt_prediction as pred
			GROUP BY pred.id) p_search
			WHERE p_search.document @@ to_tsquery(?);`, searchString).Rows()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var pid int64
		rows.Scan(&pid)
		toReturn = append(toReturn, pid)
	}
	return toReturn
}

func GetPredictionsForTag(db *gorm.DB, tag string) []PtPrediction {
	predictions := []PtPrediction{}
	db.Raw(`select p.* 
	from prediction_tag_map pmap, pt_prediction p, pt_tag t
	where pmap.pt_tag_id=t.id
	and (t.tag = ?)
	and pmap.pt_prediction_id = p.id
	group by p.id;`, tag).Find(&predictions)
	return predictions
}

func AddBracket(db *gorm.DB, b *PtBracket) {
	db.Save(b)
}

func GetMembersBracket(db *gorm.DB, User_Id int64) PtBracket {
	bracket := PtBracket{}
	db.Where("CreatorId = ?", User_Id).First(&bracket)
	return bracket
}

func GetAverageVoteForPredictionId(db *gorm.DB, predId int64) float64 {
	var avg float64
	ro := db.Debug().Raw("SELECT avg(vote_value) from (select vote_value FROM pt_vote where voted_on_id=?) as tab", predId).Row()
	ro.Scan(&avg)
	return avg
}

func SetState(db *gorm.DB, predictionId int64, state PtPredictionState) {
	prediction := PtPrediction{
		Id: predictionId,
	}
	db.First(&prediction).Update("state", state)
}

func SetScoreForPrediction(db *gorm.DB, predictionId int64, state PtPredictionState) {
	var score int
	if state == DidHappen {
		score = 1
	} else if state == DidNotHappen {
		score = -1
	}
	db.Debug().Exec(`update pt_user set score=score+? WHERE id IN (select voter_id FROM pt_vote where voted_on_id=?)`, score, predictionId)
}

func GetLivePredictionSets(db *gorm.DB) []PtPredictionSet {
	var predictionSets []PtPredictionSet
	db.Where("is_live=TRUE").Find(&predictionSets)
	return predictionSets
}

func GetLivePtHeros(db *gorm.DB) []PtHero {
	var heros []PtHero
	db.Model(&PtHero{}).Where("is_live=TRUE").Find(&heros)
	return heros
}

func SetHero(db *gorm.DB, hero *PtHero) {
	db.Save(hero)
}

func SetPredictionSet(db *gorm.DB, predictionSet *PtPredictionSet) {
	db.Save(predictionSet)
}
