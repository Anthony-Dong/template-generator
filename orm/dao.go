package orm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/anthony-dong/template-generator/utils"
	"go/format"
	"text/template"
)

const (
	defaultModelName = "model"
	defaultDaoName   = "dao"
)

type DaoMeta struct {
	DaoPackage   string
	TableName    string
	ModelPackage string
}

func (this *DaoMeta) Run(temp *template.Template) ([]byte, error) {
	var buf = &bytes.Buffer{}
	err := this.Validate()
	if err != nil {
		return nil, err
	}
	err = temp.Execute(buf, this)
	if err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}

/**
获取 struct name
*/
func (this *DaoMeta) GetStructName() string {
	return utils.LowerCaseFiledFirst(fmt.Sprintf("%sDao", this.GetModelName()))
}

/**
获取new初始化的名称
*/
func (this *DaoMeta) GetNewStructFunc() string {
	return fmt.Sprintf("New%sDao", this.GetModelName())
}

func (this *DaoMeta) GetModelName() string {
	return utils.Marshal(this.TableName)
}

func (this *DaoMeta) GetModelPackage() string {
	return this.ModelPackage
}

func (this *DaoMeta) GetModelPathName() string {
	return fmt.Sprintf("%s.%s", this.GetModelPackage(), this.GetModelName())
}

func (this *DaoMeta) Validate() error {
	if this.TableName == "" {
		return errors.New("table_name is nil")
	}
	if this.ModelPackage == "" {
		this.ModelPackage = defaultModelName
	}
	if this.DaoPackage == "" {
		this.DaoPackage = defaultDaoName
	}
	return nil
}
