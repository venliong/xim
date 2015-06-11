package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"
)

type SayhiHandler struct{}

func (this *SayhiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *SayhiHandler) doGet(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /sayhi?name=xxx"

	r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		log.Errorln("sayhi ERR:", name)
		WriteErr(w, http.StatusBadRequest, []byte(USAGE))
		return
	}

	user := users.Get(name)
	if user == nil {
		user = &User{make(chan string), time.Now().Unix()}
		users.Set(name, user)
		log.Infoln("login:", name)
	}

	w.Write([]byte(<-user.(*User).ch))

	return
}
