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
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
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
	err := request.GetConnection().SendMsg(1, []byte("Hello, welcome to Zinx!"))
	if err != nil {
		fmt.Println("send msg error: ", err)
	}
}

func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.8]")

    //配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//开启服务
	s.Serve()
}