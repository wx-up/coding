package v1

import (
	"go/ast"

	astHelp "github.com/wx-up/coding/ast/annotation"
)

type FileVisitor struct {
	ans   astHelp.NodeAnnotation
	types []*TypeVisitor
}

func (fv *FileVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		fv.ans = astHelp.NewNodeAnnotation(n, n.Doc)
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

func (tv *TypeVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Field:
		field := &FieldVisitor{
			ans: astHelp.NewNodeAnnotation(n, n.Doc),
		}
		tv.fields = append(tv.fields, field)
		return field
	}
	return tv
}

type FieldVisitor struct {
	ans astHelp.NodeAnnotation
}

func (f *FieldVisitor) Visit(node ast.Node) (w ast.Visitor) {
	return nil
}
