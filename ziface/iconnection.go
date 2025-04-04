package ziface

import "net"

// 定义封装的链接模块的接口
type IConnection interface {
	// 启动连接，让当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态
	Stop()
	// 获取当前连接的socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接的ID
	GetConnID() uint32
	// 获取远程客户端的TCP状态：IP和Port
	GetRemoteAddr() net.Addr
	// 发送数据，将数据发送给远程的客户端
	SendMsg(uint32, []byte) error
	// 带有缓冲的发送数据，非阻塞地将数据发送给远程的客户端
	SendBuffMsg(uint32, []byte) error
	// ZinxV0.10 update：设置连接属性
	SetProperty(string, interface{}) 
	// ZinxV0.10 update：获取连接属性
	GetProperty(string) (interface{}, error)
	// ZinxV0.10 update：移除连接属性
	RemoveProperty(string) 
}

// 定义一个处理链接的业务
type HandleFunc func(*net.TCPConn, []byte, int) error