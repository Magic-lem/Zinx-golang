package core

import (
	"fmt"
	"math/rand"
	"sync"
	"workspace/src/zinx/mmo_game_zinx/pb"
	"workspace/src/zinx/ziface"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 连接对象
	X    float32            // 平面的X坐标
	Y    float32            // 高度
	Z    float32            // 平面的Y坐标
	V    float32            // 旋转的角度 0 - 360
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
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // X坐标随机在160偏移，10以内
		Y:    0,                            // 默认高度为0
		Z:    float32(140 + rand.Intn(20)), // X坐标随机在160偏移，10以内
		V:    0,                            // 默认角度为0
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
	if p.Conn == nil {
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

// 获取当前玩家周围的玩家（九宫格）
func (p *Player) GetSuroundingPlayers() []*Player {
	// 得到当前AOI九宫格内所有玩家的Pid
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)

	// 利用pids获取所有player
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player := WorldMgrObj.GetPlayerByPid(int32(pid))
		players = append(players, player)
	}

	return players
}

// 发送MsgID=200的消息：同步客户端玩家的出生地点
func (p *Player) BrodCastStartPosition() {
	// 构建MsgID=200的proto数据
	proto_msg := &pb.BrodCast{
		Pid: p.Pid,
		Tp:  2, // TP为2代表广播坐标位置消息
		Data: &pb.BrodCast_P{
			P: &pb.Position{
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
		Tp:  1, // TP为1代表世界聊天消息
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

// 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	// 1. 获取当前玩家周围的玩家有哪些（九宫格）
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	// 2. 将当前玩家的位置信息通过MsgID=200的消息发送给周围的玩家（让其他玩家看到自己）
	// -- 2.1 构建Msg=200的广播消息
	proto_msg := &pb.BrodCast{
		Pid: p.Pid,
		Tp:  2, // TP为2代表广播坐标位置消息
		Data: &pb.BrodCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// -- 2.2 全部周围的玩家向各自的客户端发送这个Msg=200消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	// 3. 将周围玩家的全部位置信息发送给当前玩家客户端（让自己看到其他玩家）
	// -- 3.1 构造MsgID=202的位置同步消息
	// ---- 3.1.1 制作player切片
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		player_msg := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, player_msg)
	}
	// ---- 3.1.2 制作proto消息
	syncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:], // 切片赋值
	}

	// -- 3.2 将同步消息发送给客户端
	p.SendMsg(202, syncPlayers_proto_msg)
}

// 向其他玩家广播移动信息，更新位置
func (p *Player) UpdatePos(x, y, z, v float32) {
	// 1. 更新当前玩家player对象的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	// 2. 构建广播proto消息（MsgID=200，Tp=2）
	proto_msg := &pb.BrodCast{
		Pid: p.Pid,
		Tp:  4, // TP为4代表移动之后的坐标信息
		Data: &pb.BrodCast_P{
			P: &pb.Position{
				X: x,
				Y: y,
				Z: z,
				V: v,
			},
		},
	}

	// 3. 获取当前玩家周边感兴趣的玩家（九宫格）并发送消息
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(x, z)
	for _, pid := range pids {
		player := WorldMgrObj.GetPlayerByPid(int32(pid))
		player.SendMsg(200, proto_msg)
	}
}

// 玩家下线，向周围玩家通知下线
func (p *Player) Offline() {
	// 构造MsgID=201的proto消息
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 得到当前玩家周围玩家（九宫格）
	players := p.GetSuroundingPlayers()

	// 给周围玩家广播MsgID=201的消息，通知下线
	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	// 将当前玩家从世界管理模块和AOI管理模块中删除
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
