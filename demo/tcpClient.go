package main

import (
	"fmt"
	"net"
	"time"
)

const (
	addr = "127.0.0.1:8080"
)

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial ERR:", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("dial ok.")

	go func() {
		for {
			time.Sleep(30 * time.Second)
			if _, err = conn.Write([]byte{'B', 'B'}); err != nil {
				panic(err)
			}
		}
	}()

	sms := make([]byte, 4096)
	for {
		fmt.Print("请输入要发送的消息:")
		_, err := fmt.Scan(&sms)
		if err != nil {
			fmt.Println("数据输入异常:", err.Error())
			continue
		}

		_, err = conn.Write(sms)
		if err != nil {
			panic(err)
		}

		c, err := conn.Read(sms)
		if err != nil {
			fmt.Println("读取服务器数据异常:", err.Error())
		}

		fmt.Println(string(sms[0:c]))
	}
}
