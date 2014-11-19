package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	port   = ":8080"
	router = mux.NewRouter()
)

func addListeners() {

	//Authentication Stuff
	router.HandleFunc("/v1/auth/register", RegisterHandler) //.Methods("GET")
	router.HandleFunc("/v1/auth/login", LoginHandler)       //.Methods("GET")
	router.HandleFunc("/v1/auth/logout", LogoutHandler)     //.Methods("POST")

	//User stuff
	router.HandleFunc("/v1/users", GetAllUsersHandler)               //.Methods("GET")
	router.HandleFunc("/v1/users/featured", GetFeaturedUsersHandler) //.Methods("GET")
	router.HandleFunc("/v1/users/{id}", GetSingleUserHandler)        //.Methods("GET")

	//Prediction stuff
	router.HandleFunc("/v1/predictions", GetAllPredictionsHandler)                    //.Methods("GET")
	router.HandleFunc("/v1/predictions/featured", GetFeaturedPredictionsHandler)      //.Methods("GET")
	router.HandleFunc("/v1/predictions/{id:[0-9]+}", GetSinglePredictionHandler)      //.Methods("GET")
	router.HandleFunc("/v1/predictions/add", AddPredictionHandler)                    //.Methods("GET")
	router.HandleFunc("/v1/predictions/latest/{subcat}", GetLatestPredictionsHandler) //.Methods("GET")

	//Voting stuff
	router.HandleFunc("/v1/predictions/vote/{id}/{ud}", VoteForPredictionHandler) //.Methods("PUT")

	//Category Stuff
	router.HandleFunc("/v1/category/", GetCategoriesHandler)
	router.HandleFunc("/v1/category/{id}/sub", GetSubcategoriesHandler)
}

func beginServing() {
	fmt.Println("Listening and serving on port", port)
	http.Handle("/", router)
	http.ListenAndServe(port, nil)
}
