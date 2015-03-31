package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

/*
	Get User Functions
*/

func GetUserByID(db *gorm.DB, uid int64) PtUser {
	var user PtUser
	db.First(&user, uid)
	return user
}

func GetFeaturedUsers(db *gorm.DB) []PtUser {
	var users []PtUser
	db.Limit(10).Where(&PtUser{IsFeatured: true}).Find(&users)
	return users
}

func GetAllUsers(db *gorm.DB) []PtUser {
	users := []PtUser{}
	db.Find(&users)
	return users
}

/*
	Get Prediction Functions
*/

func GetPredictionByID(db *gorm.DB, uid int64) PtPrediction {
	var pred PtPrediction
	db.Where("id = ?", uid).First(&pred)
	return pred
}

func GetUserPrediction(db *gorm.DB, uid int64) []PtPrediction {
	var preds []PtPrediction
	db.Where("creator_id = ?", uid).Order("created desc").Find(&preds)
	return preds
}

func GetFeaturedPredictions(db *gorm.DB, l int) []PtPrediction {
	var predictions []PtPrediction
	db.Order("created desc").Where(&PtPrediction{IsFeatured: true}).Limit(l).Find(&predictions)
	return predictions
}

func GetPertinentPredictions(db *gorm.DB, CategoryId int64, limit int) []PtPrediction {
	var preds []PtPrediction
	db = db.Debug()
	createdCutOff := time.Now().Add(-time.Duration(3 * 31 * 24 * time.Hour))
	if CategoryId == 0 {
		db.Where("created >= ?", createdCutOff).Order("random()").Limit(limit).Find(&preds)
	} else {
		db.Where("category_id = ? and created>= ?", CategoryId, createdCutOff).Order("random()").Limit(limit).Find(&preds)
	}
	return preds
}

func GetAllPredictions(db *gorm.DB, limit, offset int64) []PtPrediction {
	var preds []PtPrediction
	db.Order("created desc").Limit(limit).Offset(offset).Find(&preds)
	return preds
}

func GetLatestPredictions(db *gorm.DB, x int) []PtPrediction {
	predictions := []PtPrediction{}
	db.Order("created desc").Limit(x).Find(&predictions)
	return predictions
}

/*
	Add Functions
*/

func AddPrediction(db *gorm.DB, p *PtPrediction) {
	db.Create(p)
}

func AddVote(db *gorm.DB, v *PtVote) {
	db.Save(v)
}

/*
	Category Functions
*/

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

func GetPredictionsForCategoryId(db *gorm.DB, catId int64, limit int) []PtPrediction {
	var preds []PtPrediction
	db.Order("created desc").Limit(limit).Where("category_id = ?", catId).Find(&preds)
	return preds
}

/*
	Search and Tags
*/

func SearchPredictions(db *gorm.DB, searchString string, before, after time.Time, limit int) []PtPrediction {
	db = db.Debug()
	var preds []PtPrediction
	err := db.Raw(`SELECT *
			FROM (SELECT *,
				  to_tsvector(pt_prediction.title) as document
			FROM pt_prediction
			WHERE pt_prediction.created > ? AND
			pt_prediction.created < ?
			GROUP BY pt_prediction.id) p_search
			WHERE p_search.document @@ to_tsquery(?)
			ORDER BY p_search.created desc
			LIMIT ?;`, after.UTC(), before.UTC(), searchString, limit).Find(&preds).Error
	if err != nil {
		fmt.Println("err: ", err)
	}
	return preds
}

func SearchUsers(db *gorm.DB, searchString string) []PtUser {
	var users []PtUser
	db.Raw(`SELECT * 
			FROM (SELECT *,
					to_tsvector(pt_user.first_name) ||
					to_tsvector(pt_user.last_name) ||
					to_tsvector(pt_user.affiliation) as document
			FROM pt_user
			GROUP BY pt_user.id) u_search
			WHERE u_search.document @@ to_tsquery(?)
			ORDER BY RANDOM()`, searchString).Find(&users)
	return users
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

func GetIdWithTag(db *gorm.DB, tagName string) int64 {
	var tag PtTag
	db.Where("tag = ?", tagName).First(&tag)
	return tag.Id
}

/*
	Bracket Stuff
*/

func AddBracket(db *gorm.DB, b *PtBracket) {
	db.Save(b)
}

func GetMembersBracket(db *gorm.DB, User_Id int64) PtBracket {
	var bracket PtBracket
	db.Where("CreatorId = ?", User_Id).First(&bracket)
	return bracket
}

/*
	Vote Functions
*/

func VoteExists(db *gorm.DB, uid, pid int64) bool {
	var v PtVote
	db.Where("voter_id = ? and voted_on_id = ?", uid, pid).First(&v)
	return v.Id != 0
}

func GetAverageVoteForPredictionId(db *gorm.DB, predId int64) float64 {
	var avg float64
	ro := db.Debug().Raw("SELECT avg(vote_value) from (select vote_value FROM pt_vote where voted_on_id=?) as tab", predId).Row()
	ro.Scan(&avg)
	return avg
}

//Returns -1 if vote doesnt exist, otherwise returns the value of the vote
func UpdateVoteValue(db *gorm.DB, uid int64, pred *PtPrediction) {
	UpdateCurUserVote(db, uid, pred)
	UpdateHistoricalVoteValues(db, pred)
}

func UpdateCurUserVote(db *gorm.DB, uid int64, pred *PtPrediction) {
	if uid == 0 {
		pred.CurUserVote = -1
		return
	}
	var v PtVote
	db.Where("voter_id = ? and voted_on_id = ?", uid, pred.Id).First(&v)
	if v.Id == 0 {
		pred.CurUserVote = -1
		return
	}
	pred.CurUserVote = v.VoteValue
	return
}

func UpdateHistoricalVoteValues(db *gorm.DB, pred *PtPrediction) {
	db.Order("created").Where("voted_on_id = ?", pred.Id).Find(&pred.VoteHistory)
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

/*
	Score Functions
*/

func SetState(db *gorm.DB, predictionId int64, state PtPredictionState) *PtPrediction {
	prediction := PtPrediction{
		Id: predictionId,
	}
	db.First(&prediction).Update("state", state)
	return &prediction
}

func SetScoreForVotes(db *gorm.DB, predictionId int64, state PtPredictionState) {
	var score int
	if state == DidHappen {
		score = 1
	} else if state == DidNotHappen {
		score = -1
	}
	db.Debug().Exec(`update pt_user set score=score+? WHERE id IN (select voter_id FROM pt_vote where voted_on_id=?)`, score, predictionId)
}

func SetPredictorScore(db *gorm.DB, predictorId int64, state PtPredictionState) {
	log.Println("save predictor score start", predictorId)
	db = db.Debug()
	if state == DidHappen {
		db.Exec(`update pt_user set predictions_correct = predictions_correct+1 where id = ?`, predictorId)
	}
	db.Exec(`update pt_user set predictions_graded = predictions_graded+1 where id = ?`, predictorId)
}

/*
	Prediction Sets and Heros
*/

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

func SetPredictionLocation(db *gorm.DB, predictionLoc *PtPredictionLocation) {
	db.Save(predictionLoc)
}

/*
	Auth and Login Functions
*/

func UpdatePassword(db *gorm.DB, user *PtUser) error {
	return db.First(&PtUser{}, user.Id).Update("password", user.Password).Error
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
