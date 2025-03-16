package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// 1.开始尝试连接，先等待三秒，等待服务启动
	fmt.Println("Client Test ... start")
	time.Sleep(time.Second * 3)

	// 2.建立连接，处理连接错误
	conn, err := net.Dial("tcp", "127.0.0.1:1580")
	if err != nil {
		fmt.Println("net dial err: ", err)
		return
	}

	// 3.向对端发送消息，等待返回消息并打印，间隔一秒无限循环
	for {
		_, err := conn.Write([]byte("hahaha"))
		if err != nil {
			fmt.Println("conn write err: ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn read err: ", err)
			return
		}

		fmt.Printf(" server callbacl : %s, cnt = %d \n", buf[:cnt], cnt)

		time.Sleep(time.Second * 1)
	}
}
