package core

import (
	"fmt"
	"sync"
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/mmo_game_zinx/pb"
	"google.golang.org/protobuf/proto"
	"math/rand"
)

type Player struct {
	Pid    int32  // 玩家ID
	Conn   ziface.IConnection   // 连接对象
	X      float32  // 平面的X坐标
	Y      float32  // 高度
	Z      float32  // 平面的Y坐标
	V      float32  // 旋转的角度 0 - 360
}

// Player ID 生成器，维护一个全局变量，并使用锁保护
var PidGen int32 = 1
var IdLock sync.Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// Demo：玩家ID的生成，正式的应该需要数据库管理，ID作为主键
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	// 创建玩家对象
	player := &Player{
		Pid: id,
		Conn: conn,
		X: float32(160 + rand.Intn(10)),   // X坐标随机在160偏移，10以内
		Y: 0,                              // 默认高度为0
		Z: float32(140 + rand.Intn(20)),   // X坐标随机在160偏移，10以内
		V: 0,							   // 默认角度为0
	}

	return player
}

/*
  类方法：向客户端发送消息
  在本方法中需要先使用PortoBuf进行序列化，然后使用zinx中的方法（SendMsg）进行发送
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) { // 所有的proto定义的消息都是继承自proto.Message
	// 将proto.Message消息结构体序列化为二进制数据
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}

	// 通过zinx.SendMsg方法将二进制数据发送给客户端
	if (p.Conn == nil) {
		fmt.Println("connetcion in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("player send msg err: ", err)
		return
	}
}  

// 发送MsgID=1的消息：同步客户端玩家的Pid
func (p *Player) SyncPid() {
	// 构建MsgID=1的proto数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 将消息发送给客户端
	p.SendMsg(1, data)
}

// 发送MsgID=200的消息：同步客户端玩家的出生地点
func (p *Player) BrodCastStartPosition() {
	// 构建MsgID=200的proto数据
	proto_msg := &pb.BrodCast{
		Pid: p.Pid,
		Tp: 2,  // TP为2代表广播坐标位置消息
		Data: &pb.BrodCast_P{
			&pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

// 玩家向世界广播聊天消息
func (p *Player) Talk(content string) {
	// 1. 组建MsgID：200的消息
	proto_msg := &pb.BrodCast{
		Pid: p.Pid,
		Tp: 1,   // TP为2代表世界聊天消息
		Data: &pb.BrodCast_Content{
			Content: content,
		},
	}

	// 2. 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayer()

	// 3. 向所有玩家（包括自己）发送MsgID=200的消息
	for _, player := range players {
		// 每个玩家的服务器向对应客户端发送消息
		player.SendMsg(200, proto_msg)   
	}
}