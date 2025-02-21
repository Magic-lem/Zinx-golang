package znet

import (
    "fmt"
    "net"
    "time"
    "workspace/src/zinx/ziface"
    "workspace/src/zinx/utils"
    // "errors"
)


// 定义一个服务器类
type Server struct {
    Name string
    IPVersion string
    IP string
    Port int
    MsgHandler ziface.IMsgHandle   // ZinxV0.6 update，消息管理模块，可以保存多个msgId-路由
    ConnMgr ziface.IConnManager  // ZinxV0.9 update，连接管理模块 
    OnConnStart func(conn ziface.IConnection)  // ZinxV0.9 update：Server创建连接后自动调用的Hook函数属性 — OnConnStart
    OnConnStop func(conn ziface.IConnection)  // ZinxV0.9 update：Server销毁连接前自动调用的Hook函数属性 — OnConnStop
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

// 启动服务器
func (s *Server) Start() {
    fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

    // 打印一下配置信息
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

     // 开启一个goroutine去做服务端Linster业务，主goroutine接着返回
    go func() {
        // 0. ZinxV0.6 update: 启动协程池
        s.MsgHandler.StartWorkerPool()

        // 1.获取一个TCP的Addr
        addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
        if err != nil {
            fmt.Println("resolve tcp addr err: ", err)
            return
        }
        
        // 2.监听服务器的地址
        listener, err := net.ListenTCP(s.IPVersion, addr)
        if err != nil {
            fmt.Println("net listen", s.IPVersion, "err: ", err)
            return
        }
    
        fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")  // 监听成功
        
        // 用于创建connID
        var cid uint32
        cid = 0
    
        // 3. 阻塞等待客户端连接，处理客户端连接业务（读写）
        for {
            // 3.1 阻塞等待客户端连接请求
            conn, err := listener.AcceptTCP()
            if err != nil {
                fmt.Println("Accept err", err)
                continue
            }
    
            // 3.2 ZinxV0.9 Update： 设置服务器最大连接控制，如果超过最大连接，那么则关闭此新的连接，不继续
            if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
                // TODO：给客户端发送一个超出一个最大连接的错误包
                fmt.Println("Too Many Connections, MaxConn = ", utils.GlobalObject.MaxConn)
                conn.Close()
                continue
            }
    
            // 3.3 处理该新连接请求的业务方法，此时应该有handler和conn是绑定的，得到连接模块对象
            dealConn := NewConnection(s, conn, cid, s.MsgHandler)
            cid++
            
            go dealConn.Start()  // 启动连接处理
        }
    }()
}

// 运行服务器
func (s *Server) Serve() {
    // 1.启动服务器
    s.Start()

    // TODO 中间可以加一些其他逻辑，也可用于后面的扩展
    
    // 2.阻塞（放到这里再阻塞的原因就是为了上面可以加扩展），不阻塞会导致主程退出，服务结束
    for {
        time.Sleep(time.Second * 10)
    }
}

// 停止服务器
func (s *Server) Stop() {
    fmt.Println("[STOP] Zinx server , name " , s.Name)

    // 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
    s.ConnMgr.ClearConn()  // ZinxV0.9 update 清理连接
}

// 路由功能，给当前的服务注册一个路由方法，供客户端连接使用
// ZinxV0.6 update：删除Router属性，而是增加到消息管理模块MsgHandler中
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
    s.MsgHandler.AddRouter(msgId, router)
    fmt.Println("MsgHandler Add Router Succ!!")
}

// ZinxV0.9 update：获取Server的消息管理模块
func (s *Server) GetConnMgr() ziface.IConnManager {
    return s.ConnMgr
}

func NewServer(name string) ziface.IServer {
    // 创建并返回Server类对象
    return &Server{
        Name: utils.GlobalObject.Name,    // 使用全局配置类的参数
        IPVersion: "tcp4",
        IP: utils.GlobalObject.Host,
        Port: utils.GlobalObject.TcpPort,
        MsgHandler: NewMsgHandle(),  
        ConnMgr: NewConnManager(),
        OnConnStart: nil,
        OnConnStop: nil,
    }
}

// ZinxV0.9 update：注册OnConnStart钩子函数的API
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
    s.OnConnStart = hookFunc
}

// ZinxV0.9 update：注册OnConnStop钩子函数的API
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
    s.OnConnStop = hookFunc
}

// ZinxV0.9 update：调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
    if s.OnConnStart != nil {
        fmt.Println("-----> Call OnConnStart()...")
        s.OnConnStart(connection)
    }
}

// ZinxV0.9 update：调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
    if s.OnConnStop != nil {
        fmt.Println("-----> Call OnConnStop()...")
        s.OnConnStop(connection)
    }
}