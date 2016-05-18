package lib

import (
	"bytes"
	"text/template"
)

type Locals struct {
	Image string
}

const TEMPLATE_SCRIPT string = `
#!/bin/bash
echo hello: {{.Image}}
`

func Render(locals *Locals) string {
	tmpl, err := template.New("script").Parse(TEMPLATE_SCRIPT)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, locals); err != nil {
		panic(err)
	}
	return doc.String()
}
