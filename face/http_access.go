package face

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"
	"github.com/liuhengloveyou/xim/service"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func init() {
	http.HandleFunc("/recv", recvMessage)
	http.HandleFunc("/send", sendMessage)
}

func HttpAccess() {
	//http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("/Users/liuheng/go/src/github.com/liuhengloveyou/xim-ionic/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome you!"))
		log.Infoln("RequestURI:", r.RequestURI)
	})

	s := &http.Server{
		Addr:           fmt.Sprintf("%v:%v", common.AccessConf.Addr, common.AccessConf.Port),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("HTTP IM GO... %v:%v\n", common.AccessConf.Addr, common.AccessConf.Port)
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func optionsFilter(w http.ResponseWriter, r *http.Request) {
	return

	w.Header().Set("Access-Control-Allow-Origin", "http://web.xim.com:9000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)

	return
}

func authFilter(w http.ResponseWriter, r *http.Request) (sess session.SessionStore, auth bool) {
	if e := r.ParseForm(); e != nil {
		return nil, false
	}

	token := strings.TrimSpace(r.FormValue("token"))
	if token == "" {
		sessionConf := common.AccessConf.Session.(map[string]interface{})
		if cookie, e := r.Cookie(sessionConf["cookie_name"].(string)); e == nil {
			if cookie != nil {
				token = cookie.Value
			}
		}
	}
	if token == "" {
		log.Errorln("token nil")
		return nil, false
	}

	sess, err := session.GetSessionById(token)
	if err != nil {
		log.Warningln("session ERR:", err.Error())
		return nil, false
	}

	if sess.Get("user") == nil {
		log.Errorln("session no user:", sess)
		// passport auth.
		info, err := common.Passport.UserAuth(token)
		if err != nil {
			log.Errorln("passport auth ERR:", err.Error())
			return nil, false
		}

		user := &service.User{}
		if err := json.Unmarshal(info, user); err != nil {
			log.Errorln("passport response ERR:", string(info))
			return nil, false

		}

		sess.Set("user", user)
		log.Errorln("session from passport:", sess)
		return sess, true
	}

	return nil, false
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "只支持POST请求.")
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
	case common.API_TEMPGROUP:
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
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "只支持POST请求.")
		return
	}

	var (
		user *common.UserMessage
		e    error
	)

	api := r.Header.Get("X-API")
	if api == "" {
		log.Errorln("X-API nil")
		gocommon.HttpErr(w, http.StatusBadRequest, "X-API为空.")
		return
	}
	log.Infoln("recvMessage X-API:", api)

	r.ParseForm()
	switch api {
	case common.API_TEMPGROUP:
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

// 正常1对1聊天
func chat(r *http.Request) (user *common.UserMessage, e error) {

	return nil, nil
}

// 临时讨论组
func tgroup(r *http.Request, logic string) (user *common.UserMessage, e error) {
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
