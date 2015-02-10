package main

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getSession(request *http.Request) map[string]string {
	cookieValue := make(map[string]string)
	if cookie, err := request.Cookie("session"); err == nil {
		if err := cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			log.Println(err)
		}
	}
	return cookieValue
}

func setSession(kv map[string]string, response http.ResponseWriter) {
	if encoded, err := cookieHandler.Encode("session", kv); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func ResetSession(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
}

func ViewSession(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, getSession(r))
}
