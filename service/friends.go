package service

import (
	"encoding/json"

	"github.com/liuhengloveyou/xim/dao"
)

func List(version uint) (result []byte, e error) {
	var ones []*dao.Friends
	if ones, e = (&dao.Friends{Version: int(version)}).Find(); e == nil {
		result, e = json.Marshal(ones)
	}

	return
}
