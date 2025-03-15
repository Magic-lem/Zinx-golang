package apis

import (
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/znet"
	"workspace/src/zinx/mmo_game_zinx/pb"
	"google.golang.org/protobuf/proto"
	"workspace/src/zinx/mmo_game_zinx/core"
	"fmt"
)

// 世界聊天的路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

// 重写基类路由中的Handle函数，处理世界聊天业务
func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 1. 解析客户端传递进来的proto序列化后的消息
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("unmarshal talk msg err: ", err)
		return
	}

	// 2. 确定当前的聊天数据的玩家PID
	pid, err := request.GetConnection().GetProperty("pid")

	// 3. 根据pid获得player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))	// 类型断言，将接口类型转换为其具体的底层类型。

	// 4. 服务端从客户端收到了这个消息，那在服务端将这个世界消息广播给其他全部在线玩家
	player.Talk(proto_msg.Content)
}