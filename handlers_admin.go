package main

import (
	"encoding/json"
	_ "encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SetStateHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin set state handler")
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	predictionId, _ := strconv.ParseInt(vars["predId"], 10, 64)
	stateVal, _ := strconv.Atoi(vars["state"])
	newState := PtPredictionState(stateVal)

	prediction := SetState(db, predictionId, newState)

	if newState == DidHappen || newState == DidNotHappen {
		//SetScoreForVotes(db, predictionId, newState)
		SetPredictorScore(db, prediction.CreatorId, newState)
	}

	fmt.Fprintln(w, "state set", newState)
}

func SetHeroHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin set hero handler")
	//Parse the Json
	dec := json.NewDecoder(r.Body)
	var hero PtHero
	err := dec.Decode(&hero)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	SetHero(db, &hero)
	j, err := json.Marshal(hero)
	if err != nil {
		JsonDecodeError(w, err)
	}
	fmt.Fprintln(w, string(j))
}

func SetPredictionSetHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin set prediction set handler")
	//Parse the Json
	dec := json.NewDecoder(r.Body)
	var predictionSet PtPredictionSet
	err := dec.Decode(&predictionSet)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	SetPredictionSet(db, &predictionSet)
	j, _ := json.Marshal(predictionSet)
	fmt.Fprintln(w, string(j))
}

func SetPredictionLocationHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin set prediction location handler")
	dec := json.NewDecoder(r.Body)
	var predictionLoc PtPredictionLocation
	err := dec.Decode(&predictionLoc)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	SetPredictionLocation(db, &predictionLoc)
	j, _ := json.Marshal(predictionLoc)
	fmt.Fprintln(w, string(j))
}

func GetPredictionLocationHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin get prediction location handler")
	vars := mux.Vars(r)
	cat_id, err := strconv.ParseInt(vars["cat_id"], 10, 64)
	if err != nil {
		log.Println("cat_id err", err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	db = db.Debug()
	var locs []PtPredictionLocation
	db.Where("category_id = ?", cat_id).Order("location_num").Find(&locs)
	for i, _ := range locs {
		db.First(&locs[i].Prediction, locs[i].PredictionId)
	}
	j, _ := json.Marshal(locs)
	fmt.Fprintln(w, string(j))
}

func GetHeroHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin get hero handler")
	vars := mux.Vars(r)
	cat_id, err := strconv.ParseInt(vars["cat_id"], 10, 64)
	if err != nil {
		log.Println("cat_id err", err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	db = db.Debug()
	var heros []PtHero
	db.Where("category_id = ?", cat_id).Order("location_num").Find(&heros)
	j, _ := json.Marshal(heros)
	fmt.Fprintln(w, string(j))
}

func GetPredictionSetHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin get prediction set handler")
	vars := mux.Vars(r)
	cat_id, err := strconv.ParseInt(vars["cat_id"], 10, 64)
	if err != nil {
		log.Println("cat_id err", err)
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	db = db.Debug()
	var sets []PtPredictionSet
	db.Where("category_id = ?", cat_id).Order("location_num").Find(&sets)
	j, _ := json.Marshal(sets)
	fmt.Fprintln(w, string(j))
}

func AdminPunditCreateHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin pundit create handler")
	var userMap map[string]string
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&userMap)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	newUser := PtUser{
		Email:     "None@None.com",
		Password:  "None",
		Created:   time.Now(),
		Location:  userMap["location"],
		FirstName: userMap["first_name"],
		LastName:  userMap["last_name"],
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	db = db.Debug()
	SaveUser(db, &newUser)
	j, _ := json.Marshal(newUser)
	fmt.Fprintln(w, string(j))
}

type PtSpecialEventReq struct {
	Year      int
	Category  string
	Selection string
}

func AdminSetResultForCategory(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		return
	}
	log.Println("begin set result handler")
	var req PtSpecialEventReq
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	if req.Selection == "" {
		JsonDecodeError(w, errors.New("Must specify a selection"))
		return
	}

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	var predictionsCorrect []PtPrediction
	var predictionsIncorrect []PtPrediction
	correct := db.Where("special_event_year = ? and special_event_category = ? and special_event_selection = ?", req.Year, req.Category, req.Selection).Find(&predictionsCorrect).RowsAffected
	for _, v := range predictionsCorrect {
		SetPredictorScore(db, v.CreatorId, DidHappen)
		SetState(db, v.Id, DidHappen)
	}
	incorrect := db.Where("special_event_year = ? and special_event_category = ? and special_event_selection != ?", req.Year, req.Category, req.Selection).Find(&predictionsIncorrect).RowsAffected
	for _, v := range predictionsIncorrect {
		SetPredictorScore(db, v.CreatorId, DidNotHappen)
		SetState(db, v.Id, DidNotHappen)
	}
	log.Println("correct:", correct, "incorrect:", incorrect)
	j, _ := json.Marshal(map[string]string{
		"Message": "Set Result",
	})
	fmt.Fprintln(w, string(j))
	return
}
