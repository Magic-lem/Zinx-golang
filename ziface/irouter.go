package ziface

/*
	路由抽象接口
	路由中的数据都是IRequest
*/
type IRouter interface {
	// 在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	// 处理coon业务的主方法Hook
	Handle(request IRequest)
	// 在处理conn业务之后的钩子方法Hook
	PostHandle(request IRequest)
}