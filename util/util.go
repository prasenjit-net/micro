package util

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func LogAndIgnoreError(e error) {
	if e != nil {
		log.Println(e.Error())
	}
}

func LogAndIgnoreErrorWithMessage(message string, e error) {
	if e != nil {
		log.Println(message, e.Error())
	}
}

func LogFatal(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func LogFatalWithMessage(message string, err error) {
	if err != nil {
		log.Fatalf(message, err.Error())
	}
}

func LogReqTime(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		defer log.Printf("processing time for %s is %s", r.RequestURI, time.Now().Sub(startTime))
	}
}

func Close(res io.Closer) {
	err := res.Close()
	if err != nil {
		log.Fatalf("Failed to close resource %s", err.Error())
	}
}

func CloseWithIgnore(res io.Closer) {
	err := res.Close()
	if err != nil {
		log.Printf("Failed to close resource %s", err.Error())
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
