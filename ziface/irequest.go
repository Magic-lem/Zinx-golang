package ziface

type IRequest interface {
	// 获取当前连接
	GetConnetcion() IConnection
	// 获取当前请求消息数据
	GetData() []byte
}