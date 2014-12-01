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
	categories := GetCategories(db)
	j, _ := json.Marshal(categories)
	fmt.Fprintln(w, string(j))
}

func GetSubcategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	categoryId, _ := strconv.Atoi(vars["id"])
	subcats := GetSubcategoriesWithCategoryId(db, int64(categoryId))
	j, _ := json.Marshal(subcats)
	fmt.Fprintln(w, string(j))
}

func GetSubcategoriesWithNameHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	vars := mux.Vars(r)
	categoryName := vars["name"]
	subcats := GetSubcategoriesWithCategoryName(db, categoryName)
	j, _ := json.Marshal(subcats)
	fmt.Fprintln(w, string(j))
}
