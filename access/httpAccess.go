package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim"
)

func HttpAccess() {
	http.Handle("/sayhi", &SayhiHandler{})
	http.Handle("/push", &PushMessageHandler{})
	http.HandleFunc("/send", sendMessage)

	s := &http.Server{
		Addr:           fmt.Sprintf("%v:%v", ConfJson["addr"].(string), ConfJson["port"].(float64)),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("HTTP IM GO... %v:%v\n", ConfJson["addr"].(string), ConfJson["port"].(float64))
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// 发消息
func sendMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	to, msg := r.FormValue("to"), r.FormValue("msg")

	iMsg := xim.Message{xim.MSG_SENDMSG, xim.Message_SendMsg{"", string(to), string(msg)}}
	fmt.Println("iMsg: ", iMsg)

	g := nodenet.GetGraphByName("sendmsg")
	cMsg, _ := nodenet.NewMessage(ConfJson["nodeName"].(string), g, iMsg)
	fmt.Println("cMsg: ", cMsg)

	err := nodenet.SendMsgToNext(g[0], cMsg)
	fmt.Println(err)
}
