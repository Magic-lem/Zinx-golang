package main

import "workspace/src/zinx/znet"

/*
	基于zinx框架开发的服务端应用程序
*/

func main() {
	// 1. 创建一个server句柄
	s := znet.NewServer("[zinx V0.1]")
	// 2. 启动server
	s.Serve()
}