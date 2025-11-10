package w2iface

// IRouter 路由抽象接口 使用框架者给该链接自定义的处理业务方法
type IRouter interface {
	// PreHandle 在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	// Handle 处理conn业务的主方法
	Handle(request IRequest)
	// PostHandle 在处理conn业务之后的钩子方法Hook
	PostHandle(request IRequest)
}
