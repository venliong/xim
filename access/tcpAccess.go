package main

import (
	"fmt"
	"net"
	"time"

	log "github.com/golang/glog"
)

func TcpAccess() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ConfJson["addr"].(string)), int(ConfJson["port"].(float64)), ""})
	if err != nil {
		panic(err)
	}
	fmt.Printf("TCP IM GO... %v:%v", ConfJson["addr"].(string), ConfJson["port"].(float64))

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			// log.Errorln("accept ERR:", err.Error())
			continue
		}
		defer conn.Close()
		log.Infoln("accept OK:", conn.RemoteAddr().String())

		go func() {
			data := make([]byte, 4096)

			for {
				conn.SetReadDeadline(time.Now().Add(time.Duration(HEARTBEAT)))
				i, err := conn.Read(data)
				if err != nil {
					log.Errorln("read from client ERR:", err.Error())
					fmt.Println("are you die?")
				}

				if i == 2 && data[0] == 'B' && data[1] == 'B' {
					fmt.Println("heartbeat...")
					continue
				}

				log.Infoln("read from client:", string(data[0:i]))

				conn.Write([]byte("OK"))
			}

		}()
	}
}
