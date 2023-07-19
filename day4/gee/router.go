package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	// 记录每种请求的树根节点 key为请求方法
	roots map[string]*node
	// 记录每种请求的handlerFunc key为 方法 - 路径
	handlers map[string]HandleFunc
}

// 构造函数
func newRouter() *Router {
	return &Router{
		roots:    map[string]*node{},
		handlers: map[string]HandleFunc{},
	}
}

// 拆分路径为字符串切片
func parsePattern(pattern string) []string {
	// 按照 / 切分
	split := strings.Split(pattern, "/")

	// 判断切分的每个元素进行判断，直到以 * 开头的结束
	parts := make([]string, 0)
	for _, part := range split {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 加入router
func (r Router) addRoute(method string, pattern string, handler HandleFunc) {
	// 将路径拆分
	parts := parsePattern(pattern)
	// 拼接key
	key := method + "_" + pattern
	// 判断是否有该方法的树存在，没有则创建
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 插入节点，构建路由树
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 根据真实的URL，解析出路由，以及参数
func (r *Router) getRouter(method string, path string) (*node, map[string]string) {
	// 分割URL
	searchParts := parsePattern(path)
	// 拿到该方法的路由树
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 查找路由树，返回一个叶子节点
	n := root.search(searchParts, 0)
	// 构造params
	params := make(map[string]string)
	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				// 等于同位置的路径
				params[part[1:]] = searchParts[i]
			}

			if part[0] == '*' && len(part) > 1 {
				// 把后面所有字符串拼接起来
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {
	// 通过真实的路径来解析路由
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		// 将参数放入context
		c.Params = params
		// 拼接key，是方法 + 解析出来的完整路径
		key := c.Method + "_" + n.pattern
		// 拿到方法
		handle := r.handlers[key]
		handle(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
