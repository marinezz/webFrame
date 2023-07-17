package gee

import "net/http"

// HandleFunc 框架入口

type HandleFunc func(c *Context)

// Engine 第一步： 定义结构体
type Engine struct {
	router *Router
}

// New 第二步 构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 加入路由
func (engine *Engine) addRouter(method string, pattern string, handler HandleFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET 第三步：实现GET和POST函数
func (engine *Engine) GET(pattern string, handler HandleFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, hanlder HandleFunc) {
	engine.addRouter("POST", pattern, hanlder)
}

// 第四步：重写ServeHTTP
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 创建context
	c := newContext(w, req)
	// 将c（也就是w 和 req 放入具体方法中）
	engine.router.handle(c)
}

// Run 第五步：实现run方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
