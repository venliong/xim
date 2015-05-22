package main

import (
	"fmt"
	"net"

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
			log.Errorln("accept ERR:", err.Error())
			continue
		}
		defer conn.Close()
		log.Infoln("accept OK:", conn.RemoteAddr().String())

		go func() {
			data := make([]byte, 4096)

			for {
				i, err := conn.Read(data)
				if err != nil {
					log.Errorln("read from client ERR:", err.Error())
					break
				}
				log.Infoln("read from client:", string(data[0:i]))

				conn.Write([]byte("OK"))
			}

		}()
	}
}
