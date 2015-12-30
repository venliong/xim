/*
* 用户信息逻辑服务
 */

package face

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func init() {
	http.HandleFunc("/user/", userHandler)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
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
