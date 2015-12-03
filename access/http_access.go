package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/liuhengloveyou/xim"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func HttpAccess() {
	http.HandleFunc("/recv", recvMessage)
	http.HandleFunc("/send", sendMessage)

	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/friend/", friendHandler)

	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("/Users/liuheng/go/src/github.com/liuhengloveyou/xim-ionic/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome you!"))
		log.Infoln("RequestURI:", r.RequestURI)
	})

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
	w.Header().Set("Access-Control-Allow-Origin", "*.*")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)

	return
}

func friendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)
	doOptions(w, r)
	if r.Method == "OPTIONS" {
		return
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)
	doOptions(w, r)
	if r.Method == "OPTIONS" {
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

	gocommon.HttpErr(w, stat, string(response))

	return
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	doOptions(w, r)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusNotImplemented, "只支持POST请求.")
		return
	}

	api := r.Header.Get("X-API")
	if api == "" {
		log.Errorln("X-API nil")
		gocommon.HttpErr(w, http.StatusBadRequest, "X-API为空.")
		return
	}
	log.Infoln("sendMessage X-API:", api)

	switch api {
	case xim.API_TEMPGROUP:
		if _, e := tgroup(r, "send"); e != nil {
			log.Errorln("sendMessage tgroup ERR:", e.Error())
			gocommon.HttpErr(w, http.StatusInternalServerError, "临时讨论组系统错误.")
			return
		}
	default:
		log.Errorln("X-API ERR:", api)
		gocommon.HttpErr(w, http.StatusBadRequest, "末知的X-API:"+api)
	}

	gocommon.HttpErr(w, http.StatusOK, "OK")

	return
}

func recvMessage(w http.ResponseWriter, r *http.Request) {
	var (
		user *xim.User
		e    error
	)

	doOptions(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	api := r.Header.Get("X-API")
	if api == "" {
		log.Errorln("X-API nil")
		gocommon.HttpErr(w, http.StatusBadRequest, "X-API为空.")
		return
	}
	log.Infoln("recvMessage X-API:", api)

	r.ParseForm()
	switch api {
	case xim.API_TEMPGROUP:
		if user, e = tgroup(r, "recv"); e != nil {
			gocommon.HttpErr(w, http.StatusInternalServerError, "临时讨论组系统错误.")
			return
		}
	default:
		log.Errorln("X-API ERR:", api)
		gocommon.HttpErr(w, http.StatusBadRequest, "末知的X-API:"+api)
	}
	if user == nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, "系统内部错误.")
		return
	}
	if user.ID == "" {
		gocommon.HttpErr(w, http.StatusInternalServerError, "系统内部错误.")
		return
	}

	ctx := user.HistoryMessage()
	if ctx == "" {
		select {
		case ctx = <-user.MsgChan:
		case <-time.After(1 * time.Minute):
			ctx = "TIMEOUT"
		case <-w.(http.CloseNotifier).CloseNotify():
			log.Warningln("client closed:", api, *user)
			return
		}
	}

	log.Infoln("recvover:", ctx)
	gocommon.HttpErr(w, http.StatusOK, ctx)

	return
}

func tgroup(r *http.Request, logic string) (user *xim.User, e error) {
	if "recv" == logic {
		userid, groupid := r.FormValue("uid"), r.FormValue("gid")
		if userid == "" || groupid == "" {
			return nil, fmt.Errorf("末知的用户或组.")
		}
		log.Infoln("tgroup: ", userid, groupid)

		if user, e = TGroutRecv(userid, groupid); e != nil {
			return nil, e
		}

		return user, nil
	} else if "send" == logic {
		r.ParseForm()

		userid, groupid := r.FormValue("uid"), r.FormValue("gid")
		if userid == "" || groupid == "" {
			return nil, fmt.Errorf("末知的用户或组.")
		}
		bm, e := ioutil.ReadAll(r.Body)
		if e != nil {
			return nil, e
		}
		log.Infoln("tgroup: ", userid, groupid, string(bm))

		if e = TGroutSend(userid, groupid, string(bm)); e != nil {
			return nil, e
		}
	}

	return nil, nil
}
