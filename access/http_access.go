package access

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim"
)

func HttpAccess() {
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/friend/", friendHandler)

	http.HandleFunc("/recv", recvMessage)
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}
	log.Infoln(string(body))

	stat, cookies, response, e := passport.Execute(r.RequestURI, body, r.Cookies())
	if e != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(e.Error()))
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
	gocommon.HttpErr(w, stat, response)

	return
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	to, msg := r.FormValue("to"), r.FormValue("msg")
	log.Infoln(to, msg)

	user := users.Get(to)
	if user != nil {
		user.(*User).ch <- []byte(msg)
		gocommon.HttpErr(w, http.StatusOK, nil)
		return
	}

	iMsg := xim.Message{xim.MSG_SENDMSG, xim.Message_SendMsg{"", string(to), string(msg)}}
	fmt.Println("iMsg: ", iMsg)

	g := nodenet.GetGraphByName("sendmsg")
	cMsg, _ := nodenet.NewMessage(ConfJson["nodeName"].(string), g, iMsg)
	fmt.Println("cMsg: ", cMsg)

	err := nodenet.SendMsgToNext(g[0], cMsg)
	fmt.Println(err)
}

func recvMessage(w http.ResponseWriter, r *http.Request) {
	// const USAGE = "GET /recv?userid=xxx&key=xxx"

	r.ParseForm()
	userid, key := r.FormValue("userid"), r.FormValue("key")
	if userid == "" || key == "" {
		log.Errorln("recvMessage ERR:", userid, key)
		gocommon.HttpErr(w, http.StatusBadRequest, nil)
		return
	}

	user := users.Get(userid)
	if user == nil {
		user = &User{ID: userid, ch: make(chan []byte), act: time.Now().Unix()}
		users.Set(userid, user)
		log.Infoln("login:", userid)
	}

	gocommon.HttpErr(w, http.StatusOK, <-user.(*User).ch)

	return
}
