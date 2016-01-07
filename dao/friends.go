package dao

import (
	"github.com/liuhengloveyou/xim/common"
)

type Friends struct {
	Userid  int    `xorm:"not null pk INT(11)"`
	Friends string `xorm:"JSON"`
	Version int    `xorm:"INT(11)"`
}

func (p *Friends) Insert() (e error) {
	_, e = common.Xorms["xim"].InsertOne(p)

	return
}

func (p *Friends) Find() (one []*Friends, e error) {
	e = common.Xorms["xim"].Find(one, p)

	return
}
