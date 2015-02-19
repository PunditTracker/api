package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var OscarCategories = []string{
	"Picture",
	"Director",
	"Actor",
	"Actress",
	"Supporting Actor",
	"Supporting Actress",
	"Original Screenplay",
	"Adapted Screenplay",
	"Animated Feature Film",
	"Foreign Language Film",
	"Documentary - Feature",
	"Documentary - Short",
	"Live Action Short Film",
	"Animated Short Film",
	"Original Score",
	"Original Song",
	"Sound Editting",
	"Sound Mixing",
	"Production Design",
	"Cinematography",
	"Makeup and Hairstyling",
	"Costume Design",
	"Film Editting",
	"Visual Effects",
}

type OscarStruct map[string]string

func GetSpecialEventPredictionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}
	eventName := vars["event_name"]
	eventYear, err := strconv.Atoi(vars["event_year"])
	if err != nil {
		fmt.Fprintln(w, "error in url")
		return
	}

	if eventName == "oscars" {
		GetOscarPredictions(w, uid, eventYear)
		return
	}

}

func GetOscarPredictions(w http.ResponseWriter, uid int64, year int) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	oscars := OscarStruct{}
	for _, v := range OscarCategories {
		var pred PtPrediction
		db.Where("special_event_category = ? and creator_id = ? and special_event_year = ?", v, uid, year).First(&pred)
		if pred.Id != 0 && pred.SpecialEventSelection.Valid {
			oscars[v] = pred.SpecialEventSelection.String
		} else {
			oscars[v] = ""
		}
	}
	j, _ := json.Marshal(&oscars)
	fmt.Fprintln(w, string(j))
}
