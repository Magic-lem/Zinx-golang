package ziface

type IRequest interface {
	// 获取当前连接
	GetConnection() IConnection
	// 获取当前请求消息数据
	GetData() []byte
	// Zinx-V0.5新增：获取当前请求消息的ID
	GetMsgId() uint32
}