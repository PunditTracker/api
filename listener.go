package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	port   = ":8080"
	router = mux.NewRouter()
)

func addListeners() {

	router.HandleFunc("/loadData", LoadTestDataHandler)

	//Authentication
	router.HandleFunc("/v1/auth/register", RegisterHandler)           //.Methods("GET")
	router.HandleFunc("/v1/auth/registerfb", RegisterFacebookHandler) //.Methods("GET")
	router.HandleFunc("/v1/auth/login", LoginHandler)                 //.Methods("GET")
	router.HandleFunc("/v1/auth/loginfb", LoginFacebookHandler)       //.Methods("GET")
	router.HandleFunc("/v1/auth/logout", LogoutHandler)               //.Methods("POST")
	router.HandleFunc("/v1/auth/check", CheckAuthHandler)

	//User
	router.HandleFunc("/v1/user", GetAllUsersHandler)               //.Methods("GET")
	router.HandleFunc("/v1/user/featured", GetFeaturedUsersHandler) //.Methods("GET")
	router.HandleFunc("/v1/user/{id:[0-9]+}", GetSingleUserHandler) //.Methods("GET")
	//router.HandleFunc("/v1/user/{name:[a-zA-Z0-9]+}", GetSingleUserForNameHandler) //.Methods("GET")

	//Image Functions
	router.HandleFunc("/v1/putprofpic", uploadImageHandler)

	//Prediction
	router.HandleFunc("/v1/prediction", GetAllPredictionsHandler)                             //.Methods("GET")
	router.HandleFunc("/v1/prediction/featured", GetFeaturedPredictionsHandler)               //.Methods("GET")
	router.HandleFunc("/v1/prediction/{id:[0-9]+}", GetSinglePredictionHandler)               //.Methods("GET")
	router.HandleFunc("/v1/prediction/add", AddPredictionHandler)                             //.Methods("GET")
	router.HandleFunc("/v1/prediction/latest/{subcatid:[0-9]+}", GetLatestPredictionsHandler) //.Methods("GET")
	router.HandleFunc("/v1/prediction/category/{catid:[0-9]+}", GetPredictionsForCategoryHandler)
	router.HandleFunc("/v1/prediction/subcat/{subcatid:[0-9]+}", GetPredictionsForSubcatHandler)
	router.HandleFunc("/v1/prediction/search/{searchstr}", SearchPredictionsHandler)
	router.HandleFunc("/v1/prediction/user/{id:[0-9]+}", GetUserPredictionsHandler)
	router.HandleFunc("/v1/prediction/tag/{tag}", GetTaggedPredictionHandler)

	//Tags
	router.HandleFunc("/v1/predictions/cat/{id:[0-9]+}/tag/{tag}", GetTaggedPredictionHandler)
	router.HandleFunc("/v1/tag/getSuggestedTags", GetSuggestedTagHandler)

	//Voting
	router.HandleFunc("/v1/prediction/vote/{pred_id:[0-9]+}/{value:[0-9]}", VoteForPredictionHandler) //.Methods("PUT")
	router.HandleFunc("/v1/prediction/vote/{pred_id:[0-9]+}/avg", AverageForPredictionHandler)

	//Category
	router.HandleFunc("/v1/category/", GetCategoriesHandler)
	router.HandleFunc("/v1/category/{id:[0-9]+}/sub", GetSubcategoriesHandler)
	router.HandleFunc("/v1/category/{name:[a-zA-Z0-9]+}/sub", GetSubcategoriesWithNameHandler)

	//March Madness
	router.HandleFunc("/v1/user/bracket/{userId:[0-9]+}", GetBracketHandler)
	router.HandleFunc("/v1/user/bracket/add", AddBracketHandler)

	//Homepage Functions
	router.HandleFunc("/v1/homepage/hero", GetHeroPredictionHandler)
	router.HandleFunc("/v1/homepage/predictionSet", GetPredictionSetHandler)

	//Admin Functons
	router.HandleFunc("/v1/admin/homepage/set/hero", SetHeroHandler)
	router.HandleFunc("/v1/admin/homepage/set/predictionSet", SetPredictionSetHandler)
	router.HandleFunc("/v1/admin/prediction/setstate/{predId:[0-9]+}/{state:[0-9]}", SetStateHandler)
}

type PTServer struct {
	r *mux.Router
}

func beginServing() {
	log.Println("Listening and serving on port", port)
	serv := &PTServer{router}
	http.Handle("/", serv)
	http.ListenAndServe(port, handlers.LoggingHandler(request_f, serv))
	//http.ListenAndServe(port, nil)
}

func (s *PTServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, DELETE, PATCH")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Content-Type", "application/json")
	}
	if req.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(rw, req)
}
