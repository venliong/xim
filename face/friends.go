/*
* 好友逻辑服务
 */

package face

import (
	"net/http"
)

func HttpFriends() {
	http.HandleFunc("/friends/list", friendsList)
}

func friendsList(w http.ResponseWriter, r *http.Request) {
}
