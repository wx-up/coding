package v1

type route struct {
	// 切片比较小的时候，遍历的速度不一定会比 map 检索慢， gin 就是这么做的
	// 因为 http 的 method 只有 10 个，枚举值在：http/method.go
	trees []*tree
}

func (r *route) add(method string, path string, handler HandleFunc) {
}

type tree struct {
	root   *node
	method string
}

type node struct {
	path string // 路由切分后的段

	// for 循环，比较 path 定位到 子node 节点
	children []*node
}
