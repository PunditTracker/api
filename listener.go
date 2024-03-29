package main

import (
	_ "expvar"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var (
	port   = ":3000"
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

	//Forgot Password Workflow
	router.HandleFunc("/v1/auth/forgot", ForgotPasswordEndpoint)
	router.HandleFunc("/v1/auth/reset/{id:[0-9]+}/{resetKey}", ResetPasswordEndpoint)

	//Email
	router.HandleFunc("/v1/marchmadness/add", SubscribeToMarchMadnessHandler)

	//User
	router.HandleFunc("/v1/user", GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/v1/user", UpdateUserHandler).Methods("PATCH")
	router.HandleFunc("/v1/user/featured", GetFeaturedUsersHandler) //.Methods("GET")
	router.HandleFunc("/v1/user/{id:[0-9]+}", GetSingleUserHandler) //.Methods("GET")
	router.HandleFunc("/v1/user/{id:[0-9]+}/vote", GetVotesForUserHandler)
	router.HandleFunc("/v1/user/password", UpdatePasswordHandler).Methods("POST")

	//Image Functions
	router.HandleFunc("/v1/putprofpic", UploadImageHandler)

	//Prediction
	router.HandleFunc("/v1/prediction", GetAllPredictionsHandler)               //.Methods("GET")
	router.HandleFunc("/v1/prediction/featured", GetFeaturedPredictionsHandler) //.Methods("GET")
	router.HandleFunc("/v1/prediction/{id:[0-9]+}", GetSinglePredictionHandler) //.Methods("GET")
	router.HandleFunc("/v1/prediction/add", AddPredictionHandler)               //.Methods("GET")
	router.HandleFunc("/v1/prediction/category/{cat_id:[0-9]+}", GetPredictionsForCategoryHandler)
	router.HandleFunc("/v1/prediction/category/{cat_name:[a-zA-z]+}", GetPredictionsForCategoryNameHandler)
	router.HandleFunc("/v1/prediction/page/{cat_id:[0-9]+}", GetHomePagePredictionsHandler)

	//Search
	router.HandleFunc("/v1/prediction/search/{searchstr}", SearchPredictionsHandler)
	router.HandleFunc("/v1/user/search/{searchstr}", SearchUsersHandler)
	router.HandleFunc("/v1/user/prediction/search/{searchstr}", SearchUsersPredictionsHandler)

	//User Predictions
	router.HandleFunc("/v1/prediction/user/{id:[0-9]+}", GetUserPredictionsHandler)

	//Tags
	router.HandleFunc("/v1/predictions/cat/{id:[0-9]+}/tag/{tag}", GetTaggedPredictionHandler)
	router.HandleFunc("/v1/tag/getSuggestedTags", GetSuggestedTagHandler)
	router.HandleFunc("/v1/prediction/tag/{tag}", GetTaggedPredictionHandler)

	//Voting
	router.HandleFunc("/v1/prediction/vote/{pred_id:[0-9]+}/{value:[0-9]}", VoteForPredictionHandler) //.Methods("PUT")
	router.HandleFunc("/v1/prediction/vote/{pred_id:[0-9]+}/avg", AverageForPredictionHandler)
	router.HandleFunc("/v1/prediction/vote/getown/{pred_id:[0-9]+}", GetVoteHandler)

	//Category
	router.HandleFunc("/v1/category", GetCategoriesHandler)

	//March Madness
	router.HandleFunc("/v1/user/bracket/{userId:[0-9]+}", GetBracketHandler)
	router.HandleFunc("/v1/user/bracket/add", AddBracketHandler)

	//Homepage Functions
	router.HandleFunc("/v1/hero", GetLiveHeroPredictionHandler)
	router.HandleFunc("/v1/predictionSet", GetLivePredictionSetHandler)

	router.HandleFunc("/v1/hero/{cat_id:[0-9]+}", GetLiveHeroPredictionHandler)
	router.HandleFunc("/v1/predictionSet/{cat_id:[0-9]+}", GetLivePredictionSetHandler)

	//Special
	router.HandleFunc("/v1/event/{event_name:[a-zA-z]+}/{event_year:[0-9]+}", GetSpecialEventPredictionHandler)

	//Admin Functons
	router.HandleFunc("/v1/admin/hero/set", SetHeroHandler)
	router.HandleFunc("/v1/admin/predictionSet/set", SetPredictionSetHandler)
	router.HandleFunc("/v1/admin/predictionLoc/set", SetPredictionLocationHandler)

	router.HandleFunc("/v1/admin/prediction/add", AddPredictionAdminHandler)
	router.HandleFunc("/v1/admin/pundit/add", AdminPunditCreateHandler)
	router.HandleFunc("/v1/admin/image/add/{type}", AdminUploadImageHandler)

	router.HandleFunc("/v1/admin/hero/{cat_id:[0-9]+}", GetHeroHandler)
	router.HandleFunc("/v1/admin/predictionSet/{cat_id:[0-9]+}", GetPredictionSetHandler)
	router.HandleFunc("/v1/admin/predictionLoc/{cat_id:[0-9]+}", GetPredictionLocationHandler)
	router.HandleFunc("/v1/admin/prediction/setstate/{predId:[0-9]+}/{state:[0-9]}", SetStateHandler)

	router.HandleFunc("/v1/admin/special_event/result/set", AdminSetResultForCategory)

	router.Handle("/debug/vars", http.DefaultServeMux)
	router.Handle("/debug/pprof/", http.DefaultServeMux)
	router.Handle("/debug/pprof/{x}", http.DefaultServeMux)

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
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
