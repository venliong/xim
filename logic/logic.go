/*
业务逻辑
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/passport/session"
)

type Config struct {
	NodeNames []string    `json:"nodeNames"`
	NodeConf  string      `json:"nodeConf"`
	Session   interface{} `json:"session"`
}

var (
	Sig  string
	Conf Config // 系统配置信息

	nodes    map[string]*nodenet.Component //每个进程可以运行多个woker
	passport *client.Passport
)

var (
	confile = flag.String("c", "./logic.conf.simple", "配置文件路径.")
)

func init() {
	nodes = make(map[string]*nodenet.Component)
}

func initNodenet(fn string) error {
	if e := nodenet.BuildFromConfig(fn); e != nil {
		return e
	}

	for _, name := range Conf.NodeNames {
		nodes[name] = nodenet.GetComponentByName(name)
		if nodes[name] != nil {
			nodes[name].SetHandler(nodenet.GetWorkerByName(name))
			go nodes[name].Run()
		}
	}

	return nil
}

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)
	}()
}

func main() {
	flag.Parse()

	if e := gocommon.LoadJsonConfig(*confile, &Conf); e != nil {
		panic(e)
	}

	if e := initNodenet(Conf.NodeConf); e != nil {
		panic(e)
	}

	session.InitDefaultSessionManager(Conf.Session)

	sigHandler()

	fmt.Println("logic GO...")

	for {
		time.Sleep(3 * time.Second)
	}
}
