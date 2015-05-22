package main

import (
	"net/http"

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

	users[name] = make(chan []byte)

	for {
		msg := <-users[name]
		w.Write(msg)
		break
	}

	return
}
