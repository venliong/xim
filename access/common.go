package main

import (
	"net/http"
)

func WriteErr(w http.ResponseWriter, statCode int, body []byte) {
	w.WriteHeader(statCode)
	w.Write(body)
}
