package main

import (
	"net/http"

	log "github.com/golang/glog"
)

type PushMessageHandler struct{}

func (this *PushMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *PushMessageHandler) doGet(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /sayhi?name=xxx&msg=xxx"

	r.ParseForm()
	name, msg := r.FormValue("name"), r.FormValue("msg")
	if name == "" || msg == "" {
		log.Errorln("pushmessage ERR:", name, msg)
		WriteErr(w, http.StatusBadRequest, []byte(USAGE))
		return
	}

	users[name] <- []byte(msg)

	return
}
