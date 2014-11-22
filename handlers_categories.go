package main

import (
	"fmt"
	"net/http"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	categories := GetCategories(db)
	fmt.Fprintln(w, categories)
}

func GetSubcategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := getDB()
	categoryId := int64(1)
	subcats := GetSubcategoriesWithCategoryId(db, categoryId)
	fmt.Fprintln(w, subcats)
}
