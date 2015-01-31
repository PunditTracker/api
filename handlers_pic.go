package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"net/http"
)

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	/*uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}*/

	uid := 1
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

	link := putImageOnS3(b, uniquestring)
	j, _ := json.Marshal(&map[string]interface{}{
		"Message": "Successful Upload",
		"link":    link,
	})
	fmt.Fprintln(w, string(j))
}

func fileformHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><form enctype="multipart/form-data" action='/v1/putprofpic' method='post'><input type='file' name='file'><input type='submit'></form></body>`)
}

func putImageOnS3(data []byte, uniqueIdentifier string) string {
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
	return "https://s3-us-west-1.amazonaws.com/profpics.assets.foretellr.com/" + uniqueIdentifier
}
