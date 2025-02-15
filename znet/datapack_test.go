package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责测试datapack拆包和封包功能的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	// 创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:1580")
	if err != nil {
		fmt.Println("net listen err: ", err)
		return
	}

	// 创建服务器goroutine，负责从客户端goroutine读取粘包的数据，然后进行解析
	go func() {
		// 接收客户端连接，从客户端读取数据，进行拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("listener accept err: ", err)
				continue
			}

			// 对于建立的连接，新开一个goroutine读取客户端消息
			go func(conn net.Conn) {
				// 定义封包拆包的对象
				dp := NewDataPack()
				
				for {
					// 首先读取头部的消息，并进行拆包
					headData := make([]byte, dp.GetHaedLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("read head error: ", err)
						return
					}
					msgHead, err := dp.UnPack(headData)   // 拆包
					if err != nil {
						fmt.Println("unpack error: ", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						// 说明后面存在数据内容，则进行第二次读取剩下的消息内容
						msg := msgHead.(*Message)  // 类型断言，从抽象接口类变为具体的类型（msgHead是IMessage类型，msg是Messgae类型）
						msg.Data = make([]byte, msg.GetDataLen())    // 转换为具体的类型，才可以msg.Data
						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("read msg error: ", err)
							return
						}

						// 读取完毕
						fmt.Println("Recv MsgID ", msg.Id, ", datalen = ", msg.DataLen, ", data = ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:1580")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	// 创建一个封包对象
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 1. 封装第一个msg
	msg1 := &Message {
		Id: 1,
		DataLen: 4,
		Data: []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg1 err: ", err)
		return
	}

	// 2. 封装第二个msg
	msg2 := &Message {
		Id: 2,
		DataLen: 7,
		Data: []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 err: ", err)
		return
	}

	// 将两个包粘在一起，一起发送给服务器
	// ... 是一个语法糖，用于将切片（slice）拆解成单独的元素。在你提到的 append(sendData1, sendData2...) 中，s
	// endData2... 的作用是将 sendData2 切片中的所有元素逐个添加到 sendData1 切片中。
	sendData1 = append(sendData1, sendData2...)
	if _, err := conn.Write(sendData1); err != nil {
		fmt.Println("conn write err: ", err)
		return
	}

	select {}  // 阻塞
}