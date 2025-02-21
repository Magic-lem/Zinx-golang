package main

import (
	"fmt"
	"net"
	"time"
	"io"
	"workspace/src/zinx/znet"
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
		// ZinxV0.5 需要发送封包后的消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx-V0.9 Test Message")))
		if err != nil {
			fmt.Println("pack msg err: ", err)
			return
		}

		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("conn write err: ", err)
			return
		}
		
		// ZinxV0.5 需要对接收到的消息进行拆包
		headData := make([]byte, dp.GetHaedLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("read msg head err: ", err)
			return
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack msg head err: ", err)
			return
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("read msg err: ", err)
				return
			}
		}
		msg.SetData(data)

		fmt.Printf("==> Recv Msg ID= %d, data: %s, cnt = %d \n", msg.GetMsgId(), string(msg.GetData()), msg.GetDataLen())

		time.Sleep(time.Second * 1)
	}
}