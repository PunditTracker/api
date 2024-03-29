package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()
	categories := GetCategories(db)
	if categories == nil {
		NoInfoAtEndpointError(w)
		return
	}
	j, _ := json.Marshal(categories)
	fmt.Fprintln(w, string(j))
}
