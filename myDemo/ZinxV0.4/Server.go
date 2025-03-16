package main

import (
	"fmt"
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// 重写各个函数
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping .... \n"))
	if err != nil {
		fmt.Println("call back ping ping ping error: ", err)
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error: ", err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping .... \n"))
	if err != nil {
		fmt.Println("call back ping ping ping error: ", err)
	}
}

func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.4]")

	//配置路由
	s.AddRouter(1, &PingRouter{})

	//开启服务
	s.Serve()
}
