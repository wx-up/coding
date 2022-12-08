package v1

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	astHelp "github.com/wx-up/coding/ast/annotation"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	type File struct {
		ans  []astHelp.Annotation
		typs []struct {
			ans    []astHelp.Annotation
			fields []struct {
				ans []astHelp.Annotation
			}
		}
	}
	tests := []struct {
		src  string
		want File
	}{
		{
			src: `
// annotation go through the source code and extra the annotation
// @author Deng Ming
// @date 2022/04/02
package annotation

type (
	Interface interface {
		// MyFunc is a test func
		// @parameter arg1 int
		// @parameter arg2 int32
		// @return string
		MyFunc(arg1 int, arg2 int32) string

		// second is a test func
		// @return string
		second() string
	}
)
`,
			want: File{
				ans: []astHelp.Annotation{
					{
						Key:   "author",
						Value: "Deng Ming",
					},
					{
						Key:   "date",
						Value: "2022/04/02",
					},
				},
				typs: []struct {
					ans    []astHelp.Annotation
					fields []struct {
						ans []astHelp.Annotation
					}
				}{
					{
						ans: nil,
						fields: []struct {
							ans []astHelp.Annotation
						}{
							{
								ans: []astHelp.Annotation{
									{
										Key:   "parameter",
										Value: "arg1 int",
									},
									{
										Key:   "parameter",
										Value: "arg2 int32",
									},
									{
										Key:   "return",
										Value: "string",
									},
								},
							},
							{
								ans: []astHelp.Annotation{
									{
										Key:   "return",
										Value: "string",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		fSet := token.NewFileSet()
		f, err := parser.ParseFile(fSet, "src.go", tt.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		v := &FileVisitor{}
		ast.Walk(v, f)
		assert.Equal(t, tt.want.ans, v.ans.Ans)
		for i, typ := range v.types {
			assert.Equal(t, tt.want.typs[i].ans, typ.ans.Ans)

			for ii, filed := range typ.fields {
				assert.Equal(t, tt.want.typs[i].fields[ii].ans, filed.ans.Ans)
			}
		}
	}
}
