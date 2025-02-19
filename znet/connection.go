package znet

import (
	"fmt"
	"net"
	"errors"
	"io"
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/utils"
)


type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	// 当前连接的状态	
	isClosed bool

	// 告知该连接已经停止的channel（由Reader告知Writer）
	ExitBuffChan chan bool

	// ZinxV0.6 update：消息管理模块
	MsgHandler ziface.IMsgHandle

	// ZinV0.7 update：读和写Goroutine之间的通信管道
	msgChan chan []byte
}

// 构造函数：创建一个连接
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn: conn,
		ConnID: connID,
		MsgHandler: msgHandler,
		isClosed: false,
		ExitBuffChan: make(chan bool, 1),
		msgChan: make(chan []byte),
	}
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")

	// 当此方法结束后执行退出
	defer fmt.Println("[Reader is eixt], connID = ", c.ConnID, ", remote addr is ", c.GetRemoteAddr().String())
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

		// 3. ZinxV0.6 update: 
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// -- 3.1 若是启动了协程池，则将消息添加到消息队列中
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// -- 3.2 否则，直接从消息管理模块中找到对应的处理方法，并执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
		
	}
}

// 连接的写业务方法，专门向客户端发送消息
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println("[Writer is eixt], connID = ", c.ConnID, ", remote addr is ", c.GetRemoteAddr().String())

	// 不断阻塞等待channel的消息，一旦读取到消息则发送给客户端
	for {
		select {
		case data := <- c.msgChan:
			// 有数据要发送给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("conn write err: ", err, " Conn Writer exit")
				return
			}
		case <- c.ExitBuffChan:
			// Reader告知Writer当前连接已关闭
			return
		}
	}
	
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	// 打印日志
	fmt.Println("Conn Start() ... ConnID = ", c.ConnID)

	// 开启一个从当前连接读数据的业务
	go c.StartReader()

	// 开启一个从当前连接写数据的业务
	go c.StartWriter()

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

	// 告知Writer关闭
	c.ExitBuffChan <- true

	// 回收连接中的管道
	close(c.ExitBuffChan)
	close(c.msgChan)
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

	// Zinx V0.7 update：将要发送给客户端的消息写到msgChan，由写Goroutine异步写回
	c.msgChan <- binaryMsg

	return nil
}