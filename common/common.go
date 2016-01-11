package common

import (
	"fmt"

	gocommon "github.com/liuhengloveyou/go-common"
	passport "github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/passport/session"

	_ "github.com/go-sql-driver/mysql"
)

type AccessConfig struct {
	Addr     string      `json:"addr"`
	Port     int         `json:"port"`
	NodeName string      `json:"nodeName"`
	NodeConf string      `json:"nodeConf"`
	Passport string      `json:"passport"`
	Session  interface{} `json:"session"`
	DBs      interface{} `json:"dbs"`
}

type LogicConfig struct {
	NodeConf string `json:"nodeConf"`
	Nodes    []struct {
		Name  string            `json:"name"`
		Works map[string]string `json:"works"`
	} `json:"nodes"`
	Session interface{} `json:"session"`
	DBs     interface{} `json:"dbs"`
}

var (
	AccessConf AccessConfig // 接入层配置信息
	LogicConf  LogicConfig  // 逻辑层系统配置信息

	Passport *passport.Passport

	DBs = make(map[string]*gocommon.DBmysql)
)

func InitAccessServ(confile string) error {
	if e := gocommon.LoadJsonConfig(confile, &AccessConf); e != nil {
		return e
	}

	if e := gocommon.InitDBPool(AccessConf.DBs, DBs); e != nil {
		return e
	}

	if nil == session.InitDefaultSessionManager(AccessConf.Session) {
		return fmt.Errorf("InitDefaultSessionManager err.")
	}

	Passport = &passport.Passport{ServAddr: AccessConf.Passport}

	return nil
}
