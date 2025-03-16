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

// 重写Handle函数
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	// ZinxV0.5中封装的发送消息的函数
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("send msg error: ", err)
	}
}

func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.5]")

	//配置路由
	s.AddRouter(1, &PingRouter{})

	//开启服务
	s.Serve()
}
