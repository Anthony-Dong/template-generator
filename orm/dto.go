package orm

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
)

const (
	defaultDtoPackageName = "dto"
)

type DtoMeta struct {
	TableName  string
	TableField []FieldMeta
}

type DtoMetas struct {
	packageName string
	dtos        []DtoMeta
}

func NewDtoMeta(packageName string) *DtoMetas {
	if packageName == "" {
		packageName = defaultDtoPackageName
	}
	return &DtoMetas{
		packageName: packageName,
		dtos:        []DtoMeta{},
	}
}

func (this *DtoMetas) Append(meta *ModelMeta) {
	dtoMeta := DtoMeta{
		TableName:  meta.TableName,
		TableField: meta.Fields,
	}
	this.dtos = append(this.dtos, dtoMeta)
}

func (this *DtoMetas) Run(tmpl *template.Template) ([]byte, error) {
	var buffer = &bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf(`
	package %s

	`, this.packageName))
	for _, elem := range this.dtos {
		err := tmpl.Execute(buffer, &elem)
		if err != nil {
			return nil, err
		}
	}
	return format.Source(buffer.Bytes())
}
