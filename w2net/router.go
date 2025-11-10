package w2net

import "github.com/lnnlmario/w2-inx/w2iface"

// 实现router时，先嵌入BaseRouter基类
// BaseRouter 是路由的基类，提供了默认的 PreHandle、Handle 和 PostHandle 方法
// 这样做的目的是为了让用户在实现自己的路由时，可以选择只重写需要的方法，而不必每次都实现所有方法
// 用户只需要继承 BaseRouter 并重写需要的方法即可
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request w2iface.IRequest) {
	// 默认什么都不做
}

func (br *BaseRouter) Handle(request w2iface.IRequest) {
	// 默认什么都不做
}

func (br *BaseRouter) PostHandle(request w2iface.IRequest) {
	// 默认什么都不做
}
