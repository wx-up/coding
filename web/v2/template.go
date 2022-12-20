package v2

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
)

type TemplateEngine interface {
	Render(ctx context.Context, name string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, name string, data any) ([]byte, error) {
	var res bytes.Buffer
	err := g.T.ExecuteTemplate(&res, name, data)
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

// LoadFrom 系列的方法由具体的引擎自己管理，不定义在 Web 框架的 TemplateEngine 接口中
// Web 框架只关注渲染

func (g *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

func (g *GoTemplateEngine) LoadFromFiles(filenames ...string) error {
	var err error
	g.T, err = template.ParseFiles(filenames...)
	return err
}

func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, patterns ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, patterns...)
	return err
}
