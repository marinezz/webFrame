package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
		// 模板对象，所有得模板加载进入内存
		htmlTemplates *template.Template
		// 模板渲染函数（可以自定义）
		funcMap template.FuncMap
	}
)

// New 构造函数
func New() *Engine {
	engine := &Engine{Router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default 默认函数，会调用日志和错误恢复 两个中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// SetFuncMap 自定义渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 加载模板方法
func (engine *Engine) LoadHTMLGlob(pattern string) {
	// Must() 检查模板是否正确 Funcs()添加模板函数 	ParseGlob()解析指定路径下模板
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// 创建静态文件方法
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	// 拼接绝对路径
	absolutePath := path.Join(group.prefix, relativePath)
	// http.FileServer(),返回以一个handler，这个handler向request提供文件系统的内容，直接会定位到index.html
	// http.StripPrefix()（第一个参数前缀，第二个参数处理函数）去掉前缀，然后交给handler处理
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		// 获取传递的文件名称
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.setCode(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 将静态文件加入路由
func (group *RouterGroup) Static(relativePath string, root string) {
	// http.Dir() 将字符串路径转换成文件系统
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	// 拼接url
	urlPattern := path.Join(relativePath, "/*filepath")
	// 加入路由表
	group.GET(urlPattern, handler)
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

// Use 将中间件应用到group中
func (group *RouterGroup) Use(middlerwares ...HandleFunc) {
	// 添加进入分组的中间件中
	group.middlewares = append(group.middlewares, middlerwares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandleFunc
	// 根据具体请求，找到所有使用的中间件
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	// 把所有中间件放入context中
	c.handlers = middlewares
	c.engine = engine
	engine.Router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
