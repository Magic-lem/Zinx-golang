package main

import (
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/znet"
	"workspace/src/zinx/mmo_game_zinx/core"
	"workspace/src/zinx/mmo_game_zinx/apis"
	"fmt"
)

// 接收客户端请求并建立连接之后的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)

	// 向客户端发送MsgID=1的消息，同步玩家ID
	player.SyncPid()

	// 向客户端发送MsgID=200的消息，广播玩家登录初始地位置
	player.BrodCastStartPosition()

	// 将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	// 将此连接绑定上一个player id的属性，为了后面能根据连接确定玩家ID
	conn.SetProperty("pid", player.Pid)

	fmt.Println("======> Player pid = ", player.Pid, " is arrived <=======")
}

func main() {
	// 创建一个服务器句柄
	s := znet.NewServer("[zinx for MOM Game]")

	// 绑定连接创建和销毁时的钩子函数
	s.SetOnConnStart(OnConnectionAdd)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})   // 为MsgID=2的消息注册路由

	// 启动服务
	s.Serve()
}