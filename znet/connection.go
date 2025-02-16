package znet

import (
	"fmt"
	"net"
	"errors"
	"io"
	"workspace/src/zinx/ziface"
)


type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	// 当前连接的状态	
	isClosed bool

	// 告知该连接已经停止的channel
	ExitBuffChan chan bool

	// 当前连接所对应的路由
	Router ziface.IRouter
}

// 构造函数：创建一个连接
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn: conn,
		ConnID: connID,
		Router: router,
		isClosed: false,
		ExitBuffChan: make(chan bool, 1),
	}
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	// 当此方法结束后执行退出
	defer fmt.Println("connID = ", c.ConnID, "Reader is eixt, remote addr is ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// Zinx-V0.5修改，集成消息封装类型、拆包机制
		// 1. 读取客户端的消息数据，并进行拆包
		dp := NewDataPack()
		
		// - 1.1 读取消息头部并拆包
		headData := make([]byte, dp.GetHaedLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err: ", err)
			break
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("msg head unpack err: ", err)
			break
		}

		// - 1.2 根据消息头部中记录的消息长度，读取消息内容
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err !=nil {
				fmt.Println("read msg err: ", err)
				break
			}
		}
		msg.SetData(data)

		// 2. 基于连接和数据构建Request对象
		req := Request{
			conn: c,
			msg: msg,
		}

		// 3. 从当前连接的路由中找到对应的处理方法，并执行
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		} (&req)

	}
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	// 打印日志
	fmt.Println("Conn Start() ... ConnID = ", c.ConnID)

	// 开启一个从当前连接读数据的业务
	go c.StartReader()

	// TODO 开启一个从当前连接写数据的业务

}

// 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	// 打印日志
	fmt.Println("Conn Stop() ... ConnID = ", c.ConnID)

	// 检查当前状态
	if c.isClosed == true {
		return
	}

	// 关闭当前连接
	c.Conn.Close()

	// 回收连接中的管道
	close(c.ExitBuffChan)
}

// 获取当前连接的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态：IP和Port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// ZinxV0.5 update：提供一个SendMsg方法，实现对要发送给客户端的数据进行封包，然后发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	// 判断下连接是否关闭
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将数据进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack msg err: ", err)
		return errors.New("pack msg error")
	}

	// 发送数据
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("conn write msg id ", msgId, ", error: ", err)
		return errors.New("conn Write error")
	}

	return nil
}