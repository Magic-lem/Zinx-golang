package main

import (
	"fmt"
	"workspace/src/zinx/ziface"
    "workspace/src/zinx/znet"
)

//test 自定义两个路由
type PingRouter struct {
	znet.BaseRouter
}

// 重写Handle函数
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Ping Router Handle")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendBuffMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("send msg error: ", err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (hr *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Hello Router Handle")
	// 先读取客户端的数据，再回写Hello, welcome to Zinx!
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendBuffMsg(1, []byte("Hello, welcome to Zinx!"))
	if err != nil {
		fmt.Println("send msg error: ", err)
	}
}

// 创建连接成功后执行的回调函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=====> DoConnectionBegin is Called...")
	// 向客户端发送成功建立连接的消息
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println("DoConnectionBegin conn sendmsg error: ", err)
	}

	// 给当前连接设置一些属性
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Magic-lem")
	conn.SetProperty("Home", "https://github.com/Magic-lem")
}

// 断开连接前执行的回调函数
func DoConnectionLost(conn ziface.IConnection) {
	// 查询当前连接的一些自定义属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}

	fmt.Println("=====> DoConnectionLost is Called...")
	fmt.Println("conn ID = ", conn.GetConnID(), "is Lost...")
}

func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.10]")

	// 注册连接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

    //配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//开启服务
	s.Serve()
}