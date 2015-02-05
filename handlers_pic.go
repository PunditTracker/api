package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"net/http"
)

func AdminUploadImageHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	file, h, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	uniquestring := fmt.Sprintf("images/%s", h.Filename)
	bucketName := "assets.foretellr.com"
	contType := h.Header.Get("Content-Type")
	link := putImageOnS3(bucketName, data, contType, uniquestring)
	fmt.Fprintln(w, link)
}

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}
	r.ParseForm()
	file, h, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	uniquestring := fmt.Sprintf("prof_pic/%d", uid)

	bucketName := "assets.foretellr.com"
	contType := h.Header.Get("Content-Type")
	link := putImageOnS3(bucketName, b, contType, uniquestring)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	var user PtUser
	db.First(&user, uid)
	user.Avatar_URL = link
	db.Save(&user)

	j, _ := json.Marshal(&map[string]interface{}{
		"Message": "Successful Upload",
		"link":    link,
	})
	fmt.Fprintln(w, string(j))
}

func fileformHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><form enctype="multipart/form-data" action='/v1/putprofpic' method='post'><input type='file' name='file'><input type='submit'></form></body>`)
}

func putImageOnS3(bucketName string, data []byte, imageType string, uniqueIdentifier string) string {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	s := s3.New(auth, aws.USWest)
	b := s.Bucket(bucketName)
	err = b.Put(uniqueIdentifier, data, imageType, s3.PublicReadWrite)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("https://s3-us-west-1.amazonaws.com/%s/%s", bucketName, uniqueIdentifier)
}
