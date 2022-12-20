package v2

import (
	"html/template"
	"testing"
)

func TestGoTemplateEngin_Render(t *testing.T) {
	// 测试数据可以放在当前目录 testdata（ 不存在就创建 ）
	// gohtml 后缀的话 goland 可以解析 html，如果是 tpl 则不行
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	s := NewServer(WithTemplateEngine(&GoTemplateEngine{T: tpl}))

	s.Get("/hello", func(ctx *Context) {
		err = ctx.Render("hello.gohtml", nil)
		if err != nil {
			t.Fatal(err)
		}
	})
	err = s.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}
