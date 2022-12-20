package v2

import (
	"html/template"
	"testing"
)

func TestGoTemplateEngin_Render(t *testing.T) {
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
