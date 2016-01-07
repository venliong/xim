/*
* 好友逻辑服务
 */

package face

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/liuhengloveyou/xim/service"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func init() {
	http.HandleFunc("/friends/list", friendsList)
}

func friendsList(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "GET" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "只支持GET请求")
		return
	}

	if e := r.ParseForm(); e != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, e.Error())
		return
	}

	ver := strings.TrimSpace(r.FormValue("v"))
	if ver == "" {
		ver = "0"
	}

	iver, e := strconv.ParseUint(ver, 10, 64)
	if e != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, e.Error())
		return
	}
	result, e := service.List(uint(iver))
	if e != nil {
		log.Errorln("friendsList ERR:", e.Error())
		gocommon.HttpErr(w, http.StatusInternalServerError, "数据库服务错误.")
		return
	}

	if _, e = w.Write(result); e != nil {
		log.Exitln(e)
	}
}