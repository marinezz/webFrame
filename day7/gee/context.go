package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 重名了H，方便构造json数据
type H map[string]interface{}

// Context context 第一步：构造属性
type Context struct {
	// 请求与响应
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求路径、方法、参数
	Path   string
	Method string
	Params map[string]string
	// 响应码
	StatusCode int
	// 记录所有的执行函数，前面是中间件，最后一个是要执行的函数
	handlers []HandleFunc
	index    int     // 记录执行顺序
	engine   *Engine // 拿到 HTML模板
}

// 第二部：构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// 控制权交给下一个中间件
func (c *Context) next() {
	c.index++
	len := len(c.handlers)
	// 不是所有的handler都会调用next(),这样写兼容性更好，而不是用判断语句
	for ; c.index < len; c.index++ {
		// 执行逻辑
		c.handlers[c.index](c)
	}
}

// Query 第三步 ：获取参数的Query，PostForm，Param方法
// /hello?name=marine
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Fail 中断执行 返回错误
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// PostForm post 请求体中的参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Param param 比如 /hello/name/marine  name = marine 就是参数
func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) setHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) setCode(code int) {
	// 将请求放入响应中
	c.Writer.WriteHeader(code)
	// 放入context中
	c.StatusCode = code
}

// 第四步：实现String、JSON、Data、HTML快速构建方法
// String 用于构建错误信息
func (c *Context) String(code int, format string, value ...interface{}) {
	// 构造响应头
	c.setHeader("Content-Type", "text-plain")
	// 构造状态码
	c.setCode(code)
	// 构造响应体
	// fmt.Sprintf 格式化并返回字符串 ... 放在后面用于解序列
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Tpye", "application/json")
	c.setCode(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.setCode(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.setHeader("Content-Type", "text/html")
	c.setCode(code)
	// 对模板渲染
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
