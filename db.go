package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

var (
	DBID       = "ptdev"
	DBUSERNAME = "pundittracker"
	DBPASSWORD = "ptrack20!!"
	db_logger  *log.Logger
)

func SetPassword(db *gorm.DB, user *PtUser) error {
	passByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passByte)
	return SaveUser(db, user)
}

func SaveUser(db *gorm.DB, user *PtUser) error {
	err := db.Save(user).Error
	if err != nil {
		log.Println("save user error:", err.Error())
		return err
	}
	return nil
}

func CheckUserWithIdAndPass(db *gorm.DB, id int64, password string) (PtUser, error) {
	var user PtUser
	db.First(&user, id)
	if user.Id == 0 {
		return PtUser{}, errors.New("no user")
	}
	hashedPass := []byte(user.Password)
	e := bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	//Password accepted
	if e == nil {
		return user, nil
	} else {
		return PtUser{}, errors.New("wrong password")
	}
}

func CheckUser(db *gorm.DB, email, password string) (PtUser, error) {
	var user PtUser
	db.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		return PtUser{}, errors.New("no user")
	}
	hashedPass := []byte(user.Password)
	e := bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	//Password accepted
	if e == nil {
		return user, nil
	} else {
		return PtUser{}, errors.New("wrong password")
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
	var pred PtPrediction
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
	db.Limit(10).Where(&PtUser{IsFeatured: true}).Find(&users)
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
	db.Create(p)
}

func AddVote(db *gorm.DB, v *PtVote) {
	db.Save(v)
}

//Category Stuff
func GetCategories(db *gorm.DB) []PtCategory {
	var categories []PtCategory
	db.Find(&categories)
	return categories
}

func GetIdForCategoryName(db *gorm.DB, name string) int64 {
	var cat PtCategory
	db.Where("name = ?", name).First(&cat)
	return cat.Id
}

func GetPredictionsForCategoryId(db *gorm.DB, catId int64) []PtPrediction {
	var preds []PtPrediction
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
		log.Println("err in search", err)
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

func GetTagsForPrediction(db *gorm.DB, pid int64) []string {
	ptTags := []PtTag{}
	db.Raw(`select t.* from prediction_tag_map pmap, pt_prediction p, pt_tag t where (pmap.pt_prediction_id = ?) and pmap.pt_tag_id=t.id group by t.id`, pid).Find(&ptTags)
	toReturn := []string{}
	for _, t := range ptTags {
		toReturn = append(toReturn, t.Tag)
	}
	return toReturn
}

func GetPredictionsForTag(db *gorm.DB, tag string) []PtPrediction {
	var predictions []PtPrediction
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
	var bracket PtBracket
	db.Where("CreatorId = ?", User_Id).First(&bracket)
	return bracket
}

func GetAverageVoteForPredictionId(db *gorm.DB, predId int64) float64 {
	var avg float64
	ro := db.Debug().Raw("SELECT avg(vote_value) from (select vote_value FROM pt_vote where voted_on_id=?) as tab", predId).Row()
	ro.Scan(&avg)
	return avg
}

func SetState(db *gorm.DB, predictionId int64, state PtPredictionState) *PtPrediction {
	prediction := PtPrediction{
		Id: predictionId,
	}
	db.First(&prediction).Update("state", state)
	return &prediction
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

func SetPredictorScore(db *gorm.DB, predictorId int64, state PtPredictionState) {
	if state == DidHappen {
		db.Debug().Exec(`update pt_user set PredictionCorrect = PredictionCorrect+1 where id = ?`, predictorId)
	}
	db.Debug().Exec(`update pt_user set PredictionGraded = PredictionGraded+1 where id = ?`, predictorId)
}

func GetLivePredictionSets(db *gorm.DB, catId int64) []PtPredictionSet {
	var predictionSets []PtPredictionSet
	db.Where("is_live=TRUE and category_id = ?", catId).Find(&predictionSets)
	return predictionSets
}

func GetLivePtHeros(db *gorm.DB, catId int64) []PtHero {
	var heros []PtHero
	db.Model(&PtHero{}).Where("is_live=TRUE and category_id = ?", catId).Find(&heros)
	return heros
}

func SetHero(db *gorm.DB, hero *PtHero) {
	db.Save(hero)
}

func SetPredictionSet(db *gorm.DB, predictionSet *PtPredictionSet) {
	db.Save(predictionSet)
}

func checkEmailForExistence(w http.ResponseWriter, db *gorm.DB, email string) bool {
	var testUser PtUser
	db.Where("email = ?", email).First(&testUser)
	if testUser.Id != 0 {
		if testUser.Password == "NONE" {
			ForgotPassword(w, email)
			MustResetPasswordError(w)
			return true
		} else {
			UserAlreadyExistsError(w)
			return true
		}
	}
	return false
}

func checkEmailForNonePassword(w http.ResponseWriter, db *gorm.DB, email string) bool {
	var testUser PtUser
	db.Where("email = ?", email).First(&testUser)
	if testUser.Id != 0 && testUser.Password == "NONE" {
		ForgotPassword(w, email)
		MustResetPasswordError(w)
		return true
	}
	return false
}

func VoteExists(db *gorm.DB, uid, pid int64) bool {
	var v PtVote
	db.Where("voter_id = ? and voted_on_id = ?", uid, pid).First(&v)
	return v.Id != 0
}

func PredictionDeadlinePassedOrGraded(db *gorm.DB, predId int64, w http.ResponseWriter) bool {
	var pred PtPrediction
	db.Where("id = ?", predId).First(&pred)

	if pred.State != InFuture {
		PredictionGradedError(w)
		return true
	}

	if pred.Deadline.Year() == 1 {
		return false
	}

	if pred.Deadline.Before(time.Now()) {
		DeadlinePassedError(w)
		return true
	}
	return false
}

//Returns -1 if vote doesnt exist, otherwise returns the value of the vote
func UpdateVoteValue(db *gorm.DB, uid int64, pred *PtPrediction) {
	var v PtVote
	db.Where("voter_id = ? and voted_on_id = ?", uid, pred.Id).First(&v)
	if v.Id == 0 {
		pred.CurUserVote = -1
		return
	}
	pred.CurUserVote = v.VoteValue
	return
}

func GetVotesForUser(db *gorm.DB, uid int64) []PtVote {
	db = db.Debug()
	var votes []PtVote
	db.Where("voter_id = ?", uid).Find(&votes)
	for i, _ := range votes {
		db.Where("id = ?", votes[i].VotedOnId).First(&votes[i].VotedOn)
	}
	return votes
}
