package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

var (
	port      = ":8080"
	templates = template.Must(template.ParseFiles("templates/footer.html"))
	router    = mux.NewRouter()
)

func addListeners() {
	//Un-authenticated general stuff
	router.HandleFunc("/v1/predictions/featured", GetFeaturedPredictionsHandler).Methods("GET")
	router.HandleFunc("/v1/users/featured", GetFeaturedUsersHandler).Methods("GET")
	//Blog post goes here

	//Authentication Stuff
	router.HandleFunc("/v1/auth/register", RegisterHandler).Methods("PUT")
	router.HandleFunc("/v1/auth/login", LoginHandler).Methods("POST")
	router.HandleFunc("/v1/auth/logout", LogoutHandler).Methods("POST")

	//User stuff
	router.HandleFunc("/v1/users", GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/v1/users/{id}", GetSingleUserHandler).Methods("GET")

	//Prediction stuff
	router.HandleFunc("/v1/predictions", GetAllPredictionsHandler).Methods("GET")
	router.HandleFunc("/v1/predictions/{id}", GetSinglePredictionHandler).Methods("GET")
	router.HandleFunc("/v1/predictions/add", AddPredictionHandler).Methods("GET")
	router.HandleFunc("/v1/predictions/latest/{subcat}", GetLatestPredictionsHandler).Methods("GET")
}

func beginServing() {
	fmt.Println("Listening and serving on port", port)
	http.Handle("/", router)
	http.ListenAndServe(port, nil)
}
