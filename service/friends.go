package service

import (
	"fmt"

	"github.com/liuhengloveyou/xim/dao"
)

func List(version uint) (result string, e error) {
	one := &dao.Friends{}
	has, e := one.GetOneByVersion(version)
	fmt.Println(has, e, one.Friends)
	if e != nil {
		return "", e
	}
	if has && one.Friends != nil {
		return *one.Friends, nil
	}

	return "", nil
}
