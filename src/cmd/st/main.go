package main

import (
	"log"
	"os"
	"strings"

	"st"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "st - functional style settings generator"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fields",
			Usage: "colon separated list of fields and their names - 'Foo:string Bar:int'. The types can be any native Go datatype, or even custom types. Be sure to prefix it with the name of the package if that's the case",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "The name of the generated type.",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "The name of the generated file.",
			Value: "auto_generated_st_type.go",
		},
	}
	app.Action = func(c *cli.Context) {
		pkg := c.Args().Get(0)
		if pkg == "" {
			log.Fatal("no package name specified -- aborting")
		}
		typeName := c.String("type")
		fieldsStr := c.String("fields")
		fieldsSplit := strings.Split(fieldsStr, " ")
		var fields st.Fields
		for _, fieldInfo := range fieldsSplit {
			var field st.Field
			fieldInfoSplit := strings.Split(fieldInfo, ":")
			if len(fieldInfoSplit) == 2 {
				name, t := fieldInfoSplit[0], fieldInfoSplit[1]
				field.Name = name
				field.Type = t
				fields = append(fields, field)
			}
		}
		in := st.Info{
			Filename:   c.String("file"),
			PkgName:    pkg,
			ObjectName: typeName,
			Fields:     fields,
		}
		if err := st.Generate(in); err != nil {
			log.Fatal(err)
		}
	}
	app.Run(os.Args)
}
