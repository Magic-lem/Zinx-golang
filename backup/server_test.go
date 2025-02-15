package znet

import (
    "fmt"
    "net"
    "testing"
    "time"
	"workspace/src/zinx/ziface"
)

/*
    模拟客户端
*/
func ClientTest() {
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

// ============= Zinx V0.3 测试 ================  //

// 定义一个路由类，继承自BaseRouter
type PingRouter struct {
	BaseRouter
}

// 重写各个函数
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnetcion().GetTCPConnection().Write([]byte("before ping .... \n"))
	if err != nil {
		fmt.Println("call back ping ping ping error: ", err)
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnetcion().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error: ", err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnetcion().GetTCPConnection().Write([]byte("after ping .... \n"))
	if err != nil {
		fmt.Println("call back ping ping ping error: ", err)
	}
}

// Server模块的测试函数
func TestServer(t *testing.T) {
    /*
        服务端测试
    */
    // 创建一个server对象
    s := NewServer("[zinx V0.2]")

	// 创建router对象并添加到server中
	s.AddRouter(&PingRouter{})

    /*
        客户端测试
    */
    go ClientTest()

    // 开启服务
    s.Serve()
}