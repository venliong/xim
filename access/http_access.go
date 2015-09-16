package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func HttpAccess() {
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/friend/", friendHandler)

	http.HandleFunc("/recv", recvMessage)
	http.HandleFunc("/send", sendMessage)

	s := &http.Server{
		Addr:           fmt.Sprintf("%v:%v", Conf.Addr, Conf.Port),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("HTTP IM GO... %v:%v\n", Conf.Addr, Conf.Port)
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func doOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)

	return
}

func friendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)
	if r.Method == "OPTIONS" {
		doOptions(w, r)
		return
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)
	if r.Method == "OPTIONS" {
		doOptions(w, r)
		return
	}
	if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusNotImplemented, "只支持POST请求.")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}
	log.Infoln(r.RequestURI, string(body))

	stat, cookies, response, e := passport.Execute(r.RequestURI, body, r.Cookies())
	if e != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, e.Error())
		log.Errorln("call passport ERR: ", err)
		return
	}
	fmt.Println(stat, string(response), e)

	if cookies != nil {
		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}
	}

	doOptions(w, r)
	gocommon.HttpErr(w, stat, string(response))

	return
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	api := r.Header.Get("X-API")
	if api == "" {
		log.Errorln("X-API nil")
		gocommon.HttpErr(w, http.StatusBadRequest, "X-API为空.")
		return
	}
	log.Infoln("X-API:", api)

	r.ParseForm()
	userid, msg := r.FormValue("userid"), r.FormValue("msg")
	if userid == "" || msg == "" {
		log.Errorf("param ERR:[%s],[%s].", userid, msg)
		gocommon.HttpErr(w, http.StatusBadRequest, "请求参数错误.")
		return
	}
	log.Infoln(api, userid, msg)

	//SendMsgToUser(from, to, msg)
	gocommon.HttpErr(w, http.StatusOK, "OK")

	return
}

func recvMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid, key := r.FormValue("userid"), r.FormValue("key")
	if userid == "" || key == "" {
		log.Errorln("recvMessage ERR:", userid, key)
		gocommon.HttpErr(w, http.StatusBadRequest, "")
		return
	}
	log.Infoln("recv: ", userid)

	sess, _ := users.GetSessionById(userid)
	user := sess.Get("info")
	if user == nil {
		user = &User{ID: userid, ch: make(chan string), act: time.Now().Unix()}
		sess.Set("info", user)
		log.Infoln("login:", userid)
	}

	gocommon.HttpErr(w, http.StatusOK, <-user.(*User).ch)

	return
}
