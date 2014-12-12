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
	router.HandleFunc("/v1/auth/register", RegisterHandler)           //.Methods("GET")
	router.HandleFunc("/v1/auth/registerfb", RegisterFacebookHandler) //.Methods("GET")
	router.HandleFunc("/v1/auth/login", LoginHandler)                 //.Methods("GET")
	router.HandleFunc("/v1/auth/loginfb", LoginFacebookHandler)       //.Methods("GET")
	router.HandleFunc("/v1/auth/logout", LogoutHandler)               //.Methods("POST")

	//User stuff
	router.HandleFunc("/v1/user", GetAllUsersHandler)               //.Methods("GET")
	router.HandleFunc("/v1/user/featured", GetFeaturedUsersHandler) //.Methods("GET")
	router.HandleFunc("/v1/user/{id}", GetSingleUserHandler)        //.Methods("GET")
	router.HandleFunc("/v1/user/{id}/predictions", GetUserPredictionsHandlers)

	//Prediction stuff
	router.HandleFunc("/v1/prediction", GetAllPredictionsHandler)                             //.Methods("GET")
	router.HandleFunc("/v1/prediction/featured", GetFeaturedPredictionsHandler)               //.Methods("GET")
	router.HandleFunc("/v1/prediction/{id:[0-9]+}", GetSinglePredictionHandler)               //.Methods("GET")
	router.HandleFunc("/v1/prediction/add", AddPredictionHandler)                             //.Methods("GET")
	router.HandleFunc("/v1/prediction/latest/{subcatid:[0-9]+}", GetLatestPredictionsHandler) //.Methods("GET")
	router.HandleFunc("/v1/prediction/subcat/{subcatid:[0-9]+}", GetPredictionsForSubcatHandler)
	router.HandleFunc("/v1/prediction/search", SearchPredictionsHandler)

	//Voting stuff
	router.HandleFunc("/v1/prediction/vote/{id}/{ud}", VoteForPredictionHandler) //.Methods("PUT")

	//Category Stuff
	router.HandleFunc("/v1/category/", GetCategoriesHandler)
	router.HandleFunc("/v1/category/{id:[0-9]+}/sub", GetSubcategoriesHandler)
	router.HandleFunc("/v1/category/{name:[a-z]+}/sub", GetSubcategoriesWithNameHandler)

}

type PTServer struct {
	r *mux.Router
}

func beginServing() {
	fmt.Println("Listening and serving on port", port)
	http.Handle("/", &PTServer{router})
	http.ListenAndServe(port, nil)
}

func (s *PTServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, DELETE, PATCH")
		rw.Header().Set("Content-Type", "application/json")
	}
	if req.Method == "OPTIONS" {
		return
	}
	s.r.ServeHTTP(rw, req)
}
