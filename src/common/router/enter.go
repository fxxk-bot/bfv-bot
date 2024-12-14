package router

type RouterGroup struct {
	EventRouter EventRouter
}

var RouterGroupApp = new(RouterGroup)
