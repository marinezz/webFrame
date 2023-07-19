package gee

import "strings"

type node struct {
	// 完整路径，只有叶子节点才会存放完整路径
	pattern string
	// 路径的一部分,也就是当前的部分，比如 /hello/index  中的hello
	part string
	// 子节点
	children []*node
	// 是否模糊匹配，当以 : 开头 或者 * 开头为true
	// : 表示必选  * 表示可选
	isWild bool
}

// 匹配子节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 匹配成功，或者子节点中有模糊匹配，也算匹配成功
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 匹配子节点，用于查询,返回精准匹配和模糊匹配的所有子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归构建子树
func (n *node) insert(pattern string, parts []string, height int) {
	// 遍历到根节点，返回
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	//拿到当前遍历的节点
	part := parts[height]
	child := n.matchChild(part)
	// 如果没有匹配成功，则为节点新增子节点
	if child == nil {
		child = &node{
			part: part,
			// 是否模糊匹配
			isWild: part[0] == ':' || part[0] == '*',
		}
		// 加入子节点
		n.children = append(n.children, child)
	}

	// 继续递归
	child.insert(pattern, parts, height+1)
}

// 递归查找路由树，返回一个叶子节点
func (n *node) search(parts []string, height int) *node {
	// 递归到最底层,或者节点以*开头
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// 返回所有匹配成功的节点，包括精准匹配和模糊匹配
	children := n.matchChildren(part)

	for _, child := range children {
		// 继续迭代
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}
