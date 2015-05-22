package main

import (
	"fmt"
	"net/http"
	"time"
)

func HttpAccess() {
	http.Handle("/sayhi", &SayhiHandler{})
	http.Handle("/push", &PushMessageHandler{})

	s := &http.Server{
		Addr:           fmt.Sprintf("%v:%v", ConfJson["addr"].(string), ConfJson["port"].(float64)),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("HTTP IM GO... %v:%v", ConfJson["addr"].(string), ConfJson["port"].(float64))
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
