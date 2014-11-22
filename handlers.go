package main

import (
	"fmt"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	if getSession(r)["uid"] == "" {
		fmt.Fprintln(w, "not auth")
	} else {
		fmt.Fprintln(w, "is auth")
	}
}
