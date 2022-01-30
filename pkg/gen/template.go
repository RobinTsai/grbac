package gen

import (
	"bytes"
)

var tempStr = `package {{.PkgName}}

var {{.OutputVar}} = ` + "`{{.Content}}`"

type FileConfig struct {
	PkgName   string
	OutputVar string
	Content   *bytes.Buffer
}
