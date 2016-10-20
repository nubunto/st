package st

import (
	"bytes"
	"go/format"
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

func (in Info) filepath() string {
	return filepath.Join(".", in.PkgName, in.Filename)
}

type Fields []Field
type Field struct {
	Name, Type string
}

type generator struct {
	b   *bytes.Buffer
	err error
}

func Generate(in Info) error {
	s := &generator{b: new(bytes.Buffer)}
	s.template(in)
	s.removeExistingFile(in)
	s.createPkgDirectory(in)
	s.formatSourceContents()
	s.writeFile(in)
	s.findImportsIfPossible(in)
	return s.err
}

func (s *generator) removeExistingFile(in Info) {
	if s.err == nil {
		if err := os.Remove(in.filepath()); !os.IsNotExist(err) {
			s.err = err
		}
	}
}

func (s *generator) createPkgDirectory(in Info) {
	if s.err == nil {
		if err := os.Mkdir(in.PkgName, os.ModeDir|os.ModePerm); !os.IsExist(err) {
			s.err = err
		}
	}
}

func (s *generator) formatSourceContents() {
	if s.err == nil {
		var out []byte
		out, s.err = format.Source(s.b.Bytes())
		s.b = bytes.NewBuffer(out)
	}
}

func (s *generator) writeFile(in Info) {
	if s.err == nil {
		s.err = ioutil.WriteFile(in.filepath(), s.b.Bytes(), 0666)
	}
}

func (s *generator) findImportsIfPossible(in Info) {
	if s.err == nil {
		cmd := exec.Command("goimports", "-w", in.filepath())
		_, err := cmd.CombinedOutput()
		s.err = err
	}
}

func (s *generator) template(in Info) {
	t, err := template.New("st").Parse(mutableTemplate)
	if err != nil {
		s.err = err
	}
	if err := t.Execute(s.b, in); err != nil {
		s.err = err
	}
}
