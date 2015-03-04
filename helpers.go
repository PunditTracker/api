package main

import (
	"net/http"
	"strconv"
	"strings"
)

func GetQueryValueInt(r *http.Request, name string, defaultVal int) int {
	urlValues := r.URL.Query()
	valStr, exists := urlValues[name]
	if exists {
		val, err := strconv.Atoi(valStr[0])
		if err != nil {
			return defaultVal
		}
		return val
	} else {
		return defaultVal
	}
}

func GetQueryValueInt64(r *http.Request, name string, defaultVal int64) int64 {
	urlValues := r.URL.Query()
	valStr, exists := urlValues[name]
	if exists {
		val, err := strconv.ParseInt(valStr[0], 10, 64)
		if err != nil {
			return defaultVal
		}
		return val
	} else {
		return defaultVal
	}
}

func StringToTsQuery(input string, connector string) string {
	toReturn := strings.Join(strings.Split(input, " "), connector)
	return toReturn
}
