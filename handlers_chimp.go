package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type MCEmail struct {
	Email string `json:"email"`
}

type MCRequest struct {
	EmailInfo MCEmail `json:"email"`
	Id        string  `json:"id"`
	ApiKey    string  `json:"apikey"`
}

var (
	mc_api_key = "be65d6be6a910f396f00d0ea161b85ed-us6"
	mc_list_id = "cb7ad4c6ff"
)

func SubscribeToMarchMadnessHandler(w http.ResponseWriter, r *http.Request) {
	var valMap map[string]string
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&valMap)
	email := valMap["email"]
	if email == "" {
		return
	}
	mailChimpRequest := MCRequest{
		EmailInfo: MCEmail{Email: email},
		Id:        mc_list_id,
		ApiKey:    mc_api_key,
	}
	parts := strings.Split(mc_api_key, "-")
	dc := parts[1]
	path := "lists/subscribe.json"
	api_url := fmt.Sprintf("https://%s.api.mailchimp.com/2.0/%s", dc, path)

	data, err := json.Marshal(mailChimpRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	resp, err := http.Post(api_url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Fprintln(w, "succesfully saved email", resp.Body)

}
