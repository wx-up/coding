package v2

import (
	"strings"
)

type route struct {
	// key 表示 http method
	// value 表示路由树根节点
	trees map[string]*node
}

func newRoute() route {
	return route{
		trees: make(map[string]*node),
	}
}

func (r *route) find(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return root, true
	}

	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")
	cur := root
	for _, seg := range segments {
		if cur.children == nil {
			if cur.starChild != nil {
				return cur.starChild, true
			}
			return nil, false
		}
		child, ok := cur.children[seg]
		if !ok {
			if cur.starChild != nil {
				return cur.starChild, true
			}
			return nil, false
		}
		cur = child
	}

	// 另一种方案选择：找到了 node 就直接返回，调用方取判断 node.handler 是否为 nil
	// 这里选择的方案：即使找到 node，只要 handler 为 nil 仍旧认为没有找到
	if cur.handler == nil {
		return nil, false
	}
	return cur, true
}

func (r *route) add(method string, path string, handler HandleFunc) {
	// 根节点处理
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	if path == "/" {
		root.handler = handler
		return
	}

	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	// 循环的写法，其实可以递归
	cur := root

	// 当注册 /user/*/info 时， * 后面的段都忽略
	for _, seg := range segments {
		cur = cur.findOrCreate(seg)
		if seg == "*" {
			break
		}
	}

	// 最后一个 segment 才需要添加 handler
	cur.handler = handler
}

func (n *node) findOrCreate(seg string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	if seg == "*" {
		if n.starChild == nil {
			n.starChild = &node{
				path: "*",
			}
		}
		return n.starChild
	}

	// 不存在则创建一个，存在就返回
	child, ok := n.children[seg]
	if !ok {
		child = &node{
			path: seg,
		}
		n.children[seg] = child
	}

	return child
}

type node struct {
	handler HandleFunc
	path    string // 路由切分后的段

	// 通配符
	starChild *node

	// map 结构可以根据 path 迅速定位到子node节点
	children map[string]*node
}
