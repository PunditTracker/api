package main

import (
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"net/http"
)

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	uid = 1
	if uid == 0 {
		return
	}
	r.ParseForm()
	file, _, err := r.FormFile("file")
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
	putImageOnS3(b, uniquestring)
}

func fileformHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><form enctype="multipart/form-data" action='/v1/putprofpic' method='post'><input type='file' name='file'><input type='submit'></form></body>`)
}

func putImageOnS3(data []byte, uniqueIdentifier string) {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	s := s3.New(auth, aws.USWest)

	b := s.Bucket("profpics.assets.foretellr.com")
	err = b.Put(uniqueIdentifier, data, "png", s3.PublicReadWrite)
	if err != nil {
		panic(err.Error())
	}
}
