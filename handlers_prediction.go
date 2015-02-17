package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetFeaturedPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	urlValues := r.URL.Query()
	var limit int
	var err error
	limitStr, exists := urlValues["limit"]
	if exists {
		limit, err = strconv.Atoi(limitStr[0])
		if err != nil {
			limit = 5
		}
	} else {
		limit = 5
	}

	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	predictions := GetFeaturedPredictions(db, limit)
	if predictions == nil {
		predictions = []PtPrediction{}
		return
	}
	uid := GetUIDOrZero(r)
	if uid != 0 {
		for i, _ := range predictions {
			UpdateVoteValue(db, uid, &predictions[i])
		}
	}
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetAllPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	predictions := GetAllPredictions(db)
	if predictions == nil {
		predictions = []PtPrediction{}
		return
	}
	uid := GetUIDOrZero(r)
	if uid != 0 {
		for i, _ := range predictions {
			UpdateVoteValue(db, uid, &predictions[i])
		}
	}
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetSinglePredictionHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	prediction := GetPredictionByID(db, uid)
	if prediction.Id == 0 {
		NoInfoAtEndpointError(w)
		return
	}
	UpdateVoteValue(db, GetUIDOrZero(r), &prediction)
	j, _ := json.Marshal(prediction)
	fmt.Fprintln(w, string(j))
}

func AddPredictionHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var prediction PtPrediction
	err := dec.Decode(&prediction)
	if err != nil {
		JsonDecodeError(w, err)
		return
	}
	prediction.CreatorId = GetUIDOrRedirect(w, r)
	if prediction.CreatorId == 0 {
		return
	}
	prediction.Created = time.Now()
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	for _, t := range prediction.Tags {
		prediction.TagVal = append(prediction.TagVal, PtTag{
			Id:  GetIdWithTag(db, t),
			Tag: t,
		})
	}
	log.Println("add prediction:\n", prediction)
	AddPrediction(db, &prediction)

	//Cur user hasn't voted yet
	prediction.CurUserVote = -1

	j, _ := json.Marshal(prediction)
	fmt.Fprintln(w, string(j))
}

func GetLatestPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	//Get the 10 latest predictions
	predictions := GetLatestPredictions(db, 10)
	if predictions == nil {
		predictions = []PtPrediction{}
		return
	}
	uid := GetUIDOrZero(r)
	if uid != 0 {
		for i, _ := range predictions {
			UpdateVoteValue(db, uid, &predictions[i])
		}
	}

	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetPredictionsForCategoryHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	catId, err := strconv.ParseInt(vars["cat_id"], 10, 64)
	if err != nil {
		NoInfoAtEndpointError(w)
		return
	}
	predictions := GetPredictionsForCategoryId(db, catId)
	if predictions == nil {
		predictions = []PtPrediction{}
		return
	}
	uid := GetUIDOrZero(r)
	if uid != 0 {
		for i, _ := range predictions {
			UpdateVoteValue(db, uid, &predictions[i])
		}
	}
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetPredictionsForCategoryNameHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	catName := strings.ToUpper(vars["cat_name"])
	catId := GetIdForCategoryName(db, catName)
	if catId == 0 {
		NoInfoAtEndpointError(w)
		return
	}
	log.Println(catId)
	predictions := GetPredictionsForCategoryId(db, catId)
	if predictions == nil {
		predictions = []PtPrediction{}
		return
	}
	uid := GetUIDOrZero(r)
	if uid != 0 {
		for i, _ := range predictions {
			UpdateVoteValue(db, uid, &predictions[i])
		}
	}
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func StringToTsQuery(input string, connector string) string {
	toReturn := strings.Join(strings.Split(input, " "), connector)
	return toReturn
}

func SearchPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	searchString := vars["searchstr"]
	searchString = StringToTsQuery(searchString, " & ")
	predictions := SearchPredictions(db, searchString)
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetUserPredictionsHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	uid, _ := strconv.ParseInt(vars["id"], 10, 64)
	predictions := GetUserPrediction(db, int64(uid))
	if predictions == nil {
		predictions = []PtPrediction{}
		return
	}
	cur_uid := GetUIDOrZero(r)
	if uid != 0 {
		for i, _ := range predictions {
			UpdateVoteValue(db, cur_uid, &predictions[i])
		}
	}
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetTaggedPredictionHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	tag := vars["tag"]
	predictions := GetPredictionsForTag(db, tag)
	if predictions == nil {
		predictions = []PtPrediction{}
		return
	}
	uid := GetUIDOrZero(r)
	if uid != 0 {
		for i, _ := range predictions {
			UpdateVoteValue(db, uid, &predictions[i])
		}
	}
	j, _ := json.Marshal(predictions)
	fmt.Fprintln(w, string(j))
}

func GetHeroPredictionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catId, _ := strconv.ParseInt(vars["cat_id"], 10, 64)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	heros := GetLivePtHeros(db, catId)
	if heros == nil {
		NoInfoAtEndpointError(w)
		return
	}
	j, _ := json.Marshal(heros)
	fmt.Fprintln(w, string(j))
}

func GetPredictionSetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catId, _ := strconv.ParseInt(vars["cat_id"], 10, 64)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	uid := GetUIDOrZero(r)
	predictionSets := GetLivePredictionSets(db, catId)
	if predictionSets == nil {
		NoInfoAtEndpointError(w)
		return
	}
	for i, _ := range predictionSets {
		db.First(&predictionSets[i].Prediction1, predictionSets[i].Prediction1Id)
		db.First(&predictionSets[i].Prediction2, predictionSets[i].Prediction2Id)
		db.First(&predictionSets[i].Prediction3, predictionSets[i].Prediction3Id)
		if uid != 0 {
			UpdateVoteValue(db, uid, &predictionSets[i].Prediction1)
			UpdateVoteValue(db, uid, &predictionSets[i].Prediction2)
			UpdateVoteValue(db, uid, &predictionSets[i].Prediction3)
		}
	}

	j, _ := json.Marshal(predictionSets)
	fmt.Fprintln(w, string(j))
}
