package ziface

// 定义一个服务器接口
type IServer interface {
    // 启动服务器
    Start()
    // 停止服务器
    Stop()
    // 运行服务器
    Serve()
    // 路由功能，给当前的服务注册一个路由方法，供客户端连接使用
    AddRouter(msgId uint32, router IRouter)
    // ZinxV0.9 update：获取Server的消息管理模块
    GetConnMgr() IConnManager
    // ZinxV0.9 update：注册OnConnStart钩子函数的API
    SetOnConnStart(func(connection IConnection))
    // ZinxV0.9 update：注册OnConnStop钩子函数的API
    SetOnConnStop(func(connection IConnection))
    // ZinxV0.9 update：调用OnConnStart钩子函数的方法
    CallOnConnStart(connection IConnection)
    // ZinxV0.9 update：调用OnConnStop钩子函数的方法
    CallOnConnStop(connection IConnection)
}