package dao

import (
	"time"

	"github.com/liuhengloveyou/xim/common"
)

type User struct {
	Userid  string    `xorm:"VARCHAR(45)"`
	AddTime time.Time `xorm:"not null TIMESTAMP default 'CURRENT_TIMESTAMP'"`
	Version int       `xorm:"INT(11) version"`
}

func (p *User) Insert() (e error) {
	_, e = common.DBs["xim"].Insert("INSERT INFO user values(?,?,?)", p.Userid, time.Now(), time.Now().Unix())

	return
}

func (p *User) Update() (e error) {
	// _, e = common.DBs["xim"].Id(p.Userid).Update(p)

	return
}
