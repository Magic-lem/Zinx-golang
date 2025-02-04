package znet

import (
	"fmt"
	"net"
	"workspace/src/zinx/ziface"
)


type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	// 当前连接的状态	
	isClosed bool
	
	// 当前连接的业务处理方法api
	handleAPI ziface.HandleFunc

	// 告知该连接已经停止的channel
	ExitBuffChan chan bool
}

// 构造函数：创建一个连接
func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.HandleFunc) *Connection {
	return &Connection{
		Conn: conn,
		ConnID: connID,
		handleAPI: callbackAPI,
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
		// 读取客户端的数据到buf中，目前只考虑最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("conn read error: ", err)
			continue
		}

		// 调用当前连接所绑定的处理业务的API
		err = c.handleAPI(c.Conn, buf, cnt)
		if err != nil {
			fmt.Println("CoonID", c.ConnID, "handler api is error: ", err)
			break
		}
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

// 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}