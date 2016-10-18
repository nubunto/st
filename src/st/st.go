package st

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var mutableTemplate = `
package {{.PkgName}}

type Option func(*{{.ObjectName}}) error

type {{.ObjectName}} struct {
	{{range .Fields}}
	{{.Name}} {{.Type -}}
	{{end}}
}

{{$objName := .ObjectName}}

{{range $field := .Fields}}
func {{$field.Name}}(b {{$field.Type}}) Option {
	return func(c *{{$objName}}) error {
		c.{{$field.Name}} = b
		return nil
	}
}
{{end}}

func New(opts ...Option) (*{{.ObjectName}}, error) {
	a := new({{.ObjectName}})
	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, err
		}
	}
	return a, nil
}`

type Info struct {
	Filename   string
	PkgName    string
	ObjectName string
	Fields     Fields
}

type Fields []Field
type Field struct {
	Name, Type string
}

func Generate(in Info) error {
	b := new(bytes.Buffer)
	t, err := template.New("st").Parse(mutableTemplate)
	if err != nil {
		return err
	}
	if err := t.Execute(b, in); err != nil {
		return err
	}
	if err := createFileAndFmt(in, b); err != nil {
		return err
	}
	return nil
}

func createFileAndFmt(in Info, contents io.Reader) error {
	if err := os.Mkdir(in.PkgName, os.ModeDir|os.ModePerm); err != nil {
		return err
	}
	cmd := exec.Command("goimports")
	cmd.Stdin = contents
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(in.PkgName, in.Filename), out, 0666); err != nil {
		return err
	}
	return nil
}
