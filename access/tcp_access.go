package main

import (
	"fmt"
	"net"
	"time"

	log "github.com/golang/glog"
)

func TcpAccess() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(Conf.Addr), Conf.Port, ""})
	if err != nil {
		panic(err)
	}
	fmt.Printf("TCP IM GO... %v:%v", Conf.Addr, Conf.Port)

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			// log.Errorln("accept ERR:", err.Error())
			continue
		}
		log.Infoln("accept OK:", conn.RemoteAddr().String())

		go func(c *net.TCPConn) {
			defer c.Close()

			for {
				data := make([]byte, 4096)

				conn.SetReadDeadline(time.Now().Add(time.Duration(HEARTBEAT * 3)))
				i, err := conn.Read(data)
				if err != nil {
					log.Errorln("read from client ERR:", err.Error())
					break
				}

				if i == 4 && data[0] == 0 && data[1] == 0 && data[2] == 0 && data[3] == 0 {
					continue
				}

				log.Infoln("read from client:", string(data[0:i]))
				fmt.Println(">>> ", string(data[0:i]))
			}

		}(conn)
	}
}
