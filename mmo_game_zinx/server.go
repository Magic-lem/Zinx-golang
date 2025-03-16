package main

import (
	"fmt"
	"workspace/src/zinx/mmo_game_zinx/apis"
	"workspace/src/zinx/mmo_game_zinx/core"
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/znet"
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

	// 同步周边玩家，告知他们当前玩家已经上线，广播当前玩家的位置信息
	player.SyncSurrounding()

	fmt.Println("======> Player pid = ", player.Pid, " is arrived <=======")
}

// 客户端断开连接时触发的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	// 获取当前连接所对应的玩家的
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println("conn GetProperty pid err: ", err)
		return
	}
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 得到当前玩家周围玩家（九宫格），给周围玩家广播MsgID=201的消息，通知下线
	player.Offline()

	fmt.Println("=====> Player pid = ", pid, " offline... <=====")
}

func main() {
	// 创建一个服务器句柄
	s := znet.NewServer("[zinx for MOM Game]")

	// 绑定连接创建和销毁时的钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{}) // 为MsgID=2的消息注册路由 - 世界聊天业务
	s.AddRouter(3, &apis.MoveApi{})      // 为MsgID=3的消息注册路由 - 移动业务

	// 启动服务
	s.Serve()
}
