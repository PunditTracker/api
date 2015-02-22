package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

func AdminUploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if IsAdminOrRedirect(w, r) {
		log.Println("not admin")
		return
	}

	vars := mux.Vars(r)
	uploadType := vars["type"]
	var folder string
	if uploadType == "prediction" {
		folder = "prediction_pic"
	} else if uploadType == "hero" {
		folder = "hero_pic"
	} else {
		JsonError(w, http.StatusBadRequest, "specify either or prediction or hero")
		return
	}

	log.Println("begin upload admin image handler")

	data, h, err := GetImageDataFromRequest(w, r)
	if err != nil {
		log.Println("imager err", err.Error())
		return
	}
	uniquestring := fmt.Sprintf("%s/%s", folder, h.Filename)
	bucketName := "assets.pundittracker.com"
	contType := h.Header.Get("Content-Type")
	link, err := putImageOnS3(bucketName, data, contType, uniquestring)
	if err != nil {
		log.Println("upload error:", err.Error())
		fmt.Fprintln(w, "upload error:", err.Error())
	}
	j, _ := json.Marshal(&map[string]interface{}{
		"Message": "Successful Upload",
		"link":    link,
	})
	fmt.Fprintln(w, string(j))
}

func GetImageDataFromRequest(w http.ResponseWriter, r *http.Request) ([]byte, *multipart.FileHeader, error) {
	r.ParseForm()
	file, h, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, "formfile error", err.Error())
		log.Println("formfile error", err.Error())
		return nil, nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(w, "readfile error", err.Error())
		log.Println("readfile error", err.Error())
		return nil, nil, err
	}
	return data, h, nil
}

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	uid := GetUIDOrRedirect(w, r)
	if uid == 0 {
		return
	}
	data, h, err := GetImageDataFromRequest(w, r)
	if err != nil {
		return
	}
	uniquestring := fmt.Sprintf("prof_pic/%d", uid)
	bucketName := "assets.pundittracker.com"
	contType := h.Header.Get("Content-Type")
	link, err := putImageOnS3(bucketName, data, contType, uniquestring)
	if err != nil {
		log.Println("upload error:", err.Error())
		fmt.Fprintln(w, "upload error:", err.Error())
		return
	}
	log.Println(link)
	db := GetDBOrPrintError(w)
	if db == nil {
		return
	}
	defer db.Close()

	var user PtUser
	db.First(&user, uid)
	user.AvatarUrl = link
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

func putImageOnS3(bucketName string, data []byte, imageType string, uniqueIdentifier string) (string, error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	s := s3.New(auth, aws.USEast)
	b := s.Bucket(bucketName)
	err = b.Put(uniqueIdentifier, data, imageType, s3.PublicReadWrite)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return fmt.Sprintf("https://s3.amazonaws.com/%s/%s", bucketName, uniqueIdentifier), nil
}
