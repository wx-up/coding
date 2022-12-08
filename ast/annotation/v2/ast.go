package v2

import (
	"go/ast"

	astHelp "github.com/wx-up/coding/ast/annotation"
)

// EnterVisitor 入口 visitor
type EnterVisitor struct {
	file *FileVisitor
}

func (ev *EnterVisitor) Get() File {
	if ev.file == nil {
		return File{}
	}
	return ev.file.Get()
}

func (ev *EnterVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.File:
		ev.file = &FileVisitor{
			ans: astHelp.NewNodeAnnotation(n, n.Doc),
		}
		return ev.file
	}
	return ev
}

type FileVisitor struct {
	ans astHelp.NodeAnnotation

	types []*TypeVisitor
}

func (fv *FileVisitor) Get() File {
	res := File{
		NodeAnnotation: fv.ans,
	}
	for _, typ := range fv.types {
		res.Types = append(res.Types, typ.Get())
	}
	return res
}

func (fv *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.TypeSpec:
		typ := &TypeVisitor{
			ans: astHelp.NewNodeAnnotation(n, n.Doc),
		}
		fv.types = append(fv.types, typ)
		return typ
	}
	return fv
}

type TypeVisitor struct {
	ans    astHelp.NodeAnnotation
	fields []*FieldVisitor
}

func (tv *TypeVisitor) Get() Type {
	res := Type{
		NodeAnnotation: tv.ans,
	}
	for _, field := range tv.fields {
		res.Fields = append(res.Fields, field.Get())
	}
	return res
}

func (tv *TypeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.Field:
		field := &FieldVisitor{
			ans: astHelp.NewNodeAnnotation(node, n.Doc),
		}
		tv.fields = append(tv.fields, field)
		return field
	}
	return tv
}

type FieldVisitor struct {
	ans astHelp.NodeAnnotation
}

func (f *FieldVisitor) Get() Field {
	return Field{
		NodeAnnotation: f.ans,
	}
}

func (f *FieldVisitor) Visit(node ast.Node) (w ast.Visitor) {
	return nil
}

type File struct {
	astHelp.NodeAnnotation
	Types []Type
}

type Type struct {
	astHelp.NodeAnnotation
	Fields []Field
}

type Field struct {
	astHelp.NodeAnnotation
}
