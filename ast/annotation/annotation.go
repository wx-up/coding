package annotation

import (
	"go/ast"
	"strings"
)

type NodeAnnotation struct {
	Node ast.Node
	Ans  []Annotation
}

type Annotation struct {
	Key   string
	Value string
}

/*
注解格式：
	// @return string
*/

// NewNodeAnnotation 解析注解
func NewNodeAnnotation(node ast.Node, cg *ast.CommentGroup) NodeAnnotation {
	if cg == nil || len(cg.List) == 0 {
		return NodeAnnotation{Node: node}
	}

	ans := make([]Annotation, 0, len(cg.List))

	for _, c := range cg.List {
		text, ok := parseContent(c)
		if !ok {
			continue
		}
		if strings.Contains(text, "@") {
			seg := strings.SplitN(text, " ", 2)
			if len(seg) != 2 {
				continue
			}
			key := seg[0][1:]
			ans = append(ans, Annotation{
				Key:   key,
				Value: seg[1],
			})
		}
	}
	return NodeAnnotation{
		Node: node,
		Ans:  ans,
	}
}

func parseContent(c *ast.Comment) (string, bool) {
	text := c.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	} else if strings.HasPrefix(text, "/* ") {
		length := len(text)
		return text[3 : length-2], true
	}
	return "", false
}
