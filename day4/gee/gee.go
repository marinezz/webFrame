package gee

import (
	"log"
	"net/http"
)

// HandleFunc 框架入口

type HandleFunc func(c *Context)

type (
	RouterGroup struct {
		prefix string
		// 中间件
		middlewares []HandleFunc
		// 所有的groups共享一个engine
		engine *Engine
	}

	Engine struct {
		// 嵌套类型
		*RouterGroup
		// 路由
		Router *Router
		// 存储所有的分组
		groups []*RouterGroup
	}
)

// New 构造函数
func New() *Engine {
	engine := &Engine{Router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 先拼接前缀
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, pattern string, handler HandleFunc) {
	// 拼接真实路径
	patterns := group.prefix + pattern
	log.Printf("route : %s", patterns)
	group.engine.Router.addRoute(method, patterns, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.Router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
