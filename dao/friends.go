package dao

import (
	"github.com/liuhengloveyou/xim/common"
)

type Friends struct {
	Userid  *string `xorm:"not null pk INT(11)"`
	Friends *string `xorm:"JSON"`
	Version int     `xorm:"INT(11)"`
}

func (p *Friends) Insert() (e error) {
	_, e = common.DBs["xim"].Insert("INSERT INTO friends values(?,?,?)", p.Userid, p.Friends, 1)

	return
}

func (p *Friends) Find() (one []*Friends, e error) {
	//	e = common.DBs["xim"].Query(sqlStr string, args ...interface{})

	return
}

func (p *Friends) GetOne() (has bool, e error) {
	//	has, e = common.DBs["xim"].Get(p)

	return
}

func (p *Friends) GetOneByVersion(ver uint) (has bool, e error) {
	//	has, e = common.DBs["xim"].Where("version>?", ver).Get(p)

	return
}
