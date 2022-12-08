package v2

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	astHelp "github.com/wx-up/coding/ast/annotation"

	"github.com/stretchr/testify/assert"
)

func TestFileVisitor_Get(t *testing.T) {
	testCases := []struct {
		src  string
		want File
	}{
		{
			src: `
// annotation go through the source code and extra the annotation
// @author Deng Ming
/* @multiple first line
second line
*/
// @date 2022/04/02
package annotation

type (
	// FuncType is a type
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	FuncType func()
)

type (
	// StructType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	StructType struct {
		// Public is a field
		// @type string
		Public string
	}

	// SecondType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/03
	SecondType struct {
	}
)

type (
	// Interface is a test interface
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/04
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
				NodeAnnotation: astHelp.NodeAnnotation{
					Ans: []astHelp.Annotation{
						{
							Key:   "author",
							Value: "Deng Ming",
						},
						{
							Key:   "multiple",
							Value: "first line\nsecond line\n",
						},
						{
							Key:   "date",
							Value: "2022/04/02",
						},
					},
				},
				Types: []Type{
					{
						NodeAnnotation: astHelp.NodeAnnotation{
							Ans: []astHelp.Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/02",
								},
							},
						},
					},
					{
						NodeAnnotation: astHelp.NodeAnnotation{
							Ans: []astHelp.Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/02",
								},
							},
						},
						Fields: []Field{
							{
								NodeAnnotation: astHelp.NodeAnnotation{
									Ans: []astHelp.Annotation{
										{
											Key:   "type",
											Value: "string",
										},
									},
								},
							},
						},
					},
					{
						NodeAnnotation: astHelp.NodeAnnotation{
							Ans: []astHelp.Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/03",
								},
							},
						},
					},
					{
						NodeAnnotation: astHelp.NodeAnnotation{
							Ans: []astHelp.Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/04",
								},
							},
						},
						Fields: []Field{
							{
								NodeAnnotation: astHelp.NodeAnnotation{
									Ans: []astHelp.Annotation{
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
							},
							{
								NodeAnnotation: astHelp.NodeAnnotation{
									Ans: []astHelp.Annotation{
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
		},
	}
	for _, tc := range testCases {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "src.go", tc.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		tv := &EnterVisitor{}
		ast.Walk(tv, f)
		file := tv.Get()
		assertAnnotations(t, tc.want.NodeAnnotation, file.NodeAnnotation)
		if len(tc.want.Types) != len(file.Types) {
			t.Fatal()
		}
		for i, typ := range file.Types {
			wantType := tc.want.Types[i]
			assertAnnotations(t, wantType.NodeAnnotation, typ.NodeAnnotation)
			if len(wantType.Fields) != len(typ.Fields) {
				t.Fatal()
			}
			for j, fd := range typ.Fields {
				wantFd := wantType.Fields[j]
				assertAnnotations(t, wantFd.NodeAnnotation, fd.NodeAnnotation)
			}
		}
	}
}

func assertAnnotations(t *testing.T, wantAns astHelp.NodeAnnotation, dst astHelp.NodeAnnotation) {
	want := wantAns.Ans
	if len(want) != len(dst.Ans) {
		t.Fatal()
	}
	for i, an := range want {
		val := dst.Ans[i]
		assert.Equal(t, an.Value, val.Value)
	}
}
