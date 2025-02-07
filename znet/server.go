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
    Router ziface.IRouter   // 当前Server注册的连接所对应的路由
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
    
            // 3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
    
            // 3.3 处理该新连接请求的业务方法，此时应该有handler和conn是绑定的，得到连接模块对象
            dealConn := NewConnection(conn, cid, s.Router)
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

    //TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

// 路由功能，给当前的服务注册一个路由方法，供客户端连接使用
func (s *Server) AddRouter(router ziface.IRouter) {
    s.Router = router
    fmt.Println("Add Router Succ!!")
}

func NewServer(name string) ziface.IServer {
    // 创建并返回Server类对象
    return &Server{
        Name: utils.GlobalObject.Name,    // 使用全局配置类的参数
        IPVersion: "tcp4",
        IP: utils.GlobalObject.Host,
        Port: utils.GlobalObject.TcpPort,
        Router: nil,   // 路由方法
    }
}