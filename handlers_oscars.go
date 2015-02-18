package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var OscarCategories = []string{
	"Best Film",
}

type OscarStruct map[string]string

func GetOscarPredictions(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	var oscars OscarStruct
	for _, v := range OscarCategories {
		var pred PtPrediction
		db.Where("special_event_category = ? and creator_id = ?", v, uid).First(&pred)
		if pred.Id == 0 {
			oscars[v] = ""
		} else {
			oscars[v] = pred.Title
		}
	}

	j, _ := json.Marshal(&oscars)
	fmt.Fprintln(w, string(j))
}
