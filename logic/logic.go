/*
业务逻辑
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/passport/session"
)

var (
	Sig string

	mynodes  map[string]*nodenet.Component
	passport *client.Passport
)

var (
	confile = flag.String("c", "./logic.conf.simple", "配置文件路径.")
)

func init() {
	mynodes = make(map[string]*nodenet.Component)
}

func initNodenet(fn string) error {
	if e := nodenet.BuildFromConfig(fn); e != nil {
		return e
	}

	for i := 0; i < len(Conf.Nodes); i++ {
		name := Conf.Nodes[i].Name
		mynodes[name] = nodenet.GetComponentByName(name)
		if mynodes[name] == nil {
			return fmt.Errorf("No node: ", name)
		}

		for k, v := range Conf.Nodes[i].Works {
			t, w := nodenet.GetMessageTypeByName(k), nodenet.GetWorkerByName(v)
			if t == nil {
				return fmt.Errorf("No message registerd: %s", k)
			}
			if w == nil {
				return fmt.Errorf("No worker registerd: %s", v)
			}
			if reflect.TypeOf(t) != w.Message {
				return fmt.Errorf("worker can't recive message: %v %v", w, k)
			}
			mynodes[name].RegisterHandler(t, w.Handler)
		}

		go mynodes[name].Run()
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
