package template

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`hello, {{.Name}}`)
	if err != nil {
		t.Fatal(err)
	}

	var bs bytes.Buffer
	err = tpl.Execute(&bs, struct {
		Name string
	}{
		Name: "wx",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello, wx", bs.String())
}
