package apis

import (
	"fmt"
	"workspace/src/zinx/mmo_game_zinx/core"
	"workspace/src/zinx/mmo_game_zinx/pb"
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/znet"

	"google.golang.org/protobuf/proto"
)

// 玩家移动的路由
type MoveApi struct {
	znet.BaseRouter
}

// 重写Handle方法
func (m *MoveApi) Handle(request ziface.IRequest) {
	// 解析客户端传入的proto序列化后的消息
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("unmarshal position msg err: ", err)
		return
	}

	// 得到当前发送位置移动的是哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid err: ", err)
	}
	fmt.Printf("Player pid = %d, move (%f, %f, %f, %f)",
		pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 当前玩家向其他玩家广播当前玩家的移动信息，更新位置
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
