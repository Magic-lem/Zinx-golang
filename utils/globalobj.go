package utils

import (
	"workspace/src/zinx/ziface"
	"os"
	"encoding/json"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server相关参数
	*/
	TcpServer ziface.IServer         // 当前Zinx的全局Server对象
	Host 	  string				 // 当前服务器的IP地址
	TcpPort	  int					 // 当前服务器监听的端口
	Name	  string				 // 当前服务器名称

	/*
		Zinx相关参数
	*/
	Version   string				 // 当前Zinx版本号  
	MaxConn	  int					 // 当前服务器主机允许的最大连接数量
	MaxPacketSize uint32			 // 当前Zinx框架数据包的最大值
	WorkerPoolSize uint32			 // 当前Zinx框架业务Worker工作池数量
	MaxWorkerTaskLen uint32			 // 业务工作Worker对应负责的消息队列管道的最大长度
	MaxMsgChanLen uint32			 // 带有缓冲的消息通道SendBuffMsg中的最大长度
}


/*
	定义一个GlobalObj的全局对象
*/
var GlobalObject *GlobalObj

/*
	提供一个Reload方法，从JSON文件中加载用户提供的配置文件
*/
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)   
	}
	// 解析json到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供init方法，默认加载配置，初始化GlobalObject
*/
func init() {
	GlobalObject = & GlobalObj {
		Name:    "ZinxServerApp",
		Version: "V0.9",
		TcpPort: 1580,
		Host:    "0.0.0.0",
		MaxConn: 2,
		MaxPacketSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen: 1024,
	}

	// 从配置文件中加载用户自定义的配置
	GlobalObject.Reload()
}