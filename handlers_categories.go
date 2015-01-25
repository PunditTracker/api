package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	categories := GetCategories(db)
	j, _ := json.Marshal(categories)
	fmt.Fprintln(w, string(j))
}

func GetSubcategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	categoryId, _ := strconv.ParseInt(vars["id"], 10, 64)
	subcats := GetSubcategoriesWithCategoryId(db, categoryId)
	j, _ := json.Marshal(subcats)
	fmt.Fprintln(w, string(j))
}

func GetSubcategoriesWithNameHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	defer db.Close()
	vars := mux.Vars(r)
	categoryName := vars["name"]
	subcats := GetSubcategoriesWithCategoryName(db, categoryName)
	j, _ := json.Marshal(subcats)
	fmt.Fprintln(w, string(j))
}
