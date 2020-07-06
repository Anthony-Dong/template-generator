package orm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/anthony-dong/template-generator/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"go/format"
	"strings"
	"text/template"
)

const (
	CreateTable = "Create Table"
)

const (
	specifiedTableNamesSql = `select table_name from information_schema.tables where table_schema = ? and table_name in ('%s') and table_type = 'base table';`
	tableColumnsSql        = `select column_name,
is_nullable, if(column_type = 'tinyint(1)', 'boolean', data_type) as file_type
from information_schema.columns
where table_schema = ? and  table_name = ?
order by ordinal_position;
`
)

type FieldMeta struct {
	Name       string
	FieldType  string
	IsNullable string //OK  YES
	Tags       []string
}

func (this *FieldMeta) GetTag() string {
	var (
		tags = make([]string, 0)
	)
	for _, elem := range this.Tags {
		if strings.Compare(elem, "xorm") == 0 {
			if strings.Compare(this.Name, "id") == 0 {
				tags = append(tags, fmt.Sprintf("xorm:\"%s\"", "pk autoincr id"))
				continue
			}
			if strings.Compare(this.Name, "create_time") == 0 {
				tags = append(tags, fmt.Sprintf("xorm:\"%s %s\"", "created", this.Name))
				continue
			}
			if strings.Compare(this.Name, "update_time") == 0 {
				tags = append(tags, fmt.Sprintf("xorm:\"%s %s\"", "updated", this.Name))
				continue
			}
		}
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", elem, this.Name))
	}
	if len(tags) == 0 {
		return ""
	}
	var tagStr = ""
	for _, elem := range tags {
		tagStr += fmt.Sprintf("%s ", elem)
	}
	if len(tagStr) > 0 {
		tagStr = strings.TrimRight(tagStr, " ")
	}
	return fmt.Sprintf("`%s`", tagStr)
}

func (this *FieldMeta) GetGoField() string {
	return utils.Marshal(this.Name)
}

func (this *FieldMeta) GetGoType() string {
	switch this.FieldType {
	case "bit", "tinyint", "boolean":
		return "uint8"
	case "smallint", "year":
		return "uint16"
	case "integer", "mediumint", "int":
		return "int"
	case "bigint":
		return "uint64"
	case "date", "timestamp without time zone", "timestamp with time zone", "time with time zone", "time without time zone",
		"timestamp", "datetime", "time":
		return "time.Time"
	case "byte",
		"binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return "[]byte"
	case "text", "character", "character varying", "tsvector", "bit varying", "money", "json", "jsonb", "xml", "point", "interval", "line", "ARRAY",
		"char", "varchar", "tinytext", "mediumtext", "longtext":
		return "string"
	case "real":
		return "float32"
	case "numeric", "decimal", "double precision", "float", "double":
		return "float64"
	default:
		return "string"
	}
}

type ModelMeta struct {
	DbName      string
	TableName   string
	PackageName string
	Fields      []FieldMeta
	Db          *xorm.Engine
	Tags        []string
}

func NewModelMeta(dbName string, tableName string, packageName string, db *xorm.Engine, tags []string) *ModelMeta {
	return &ModelMeta{
		DbName:      dbName,
		TableName:   tableName,
		PackageName: packageName,
		Db:          db,
		Tags:        tags,
	}
}

func (this *ModelMeta) ModelName() string {
	return utils.Marshal(this.TableName)
}

func (this *ModelMeta) Run(template *template.Template) ([]byte, error) {
	if this.Db == nil {
		return nil, errors.New("the db is nil")
	}
	err := this.FindField()
	if err != nil {
		return nil, err
	}
	var body = &bytes.Buffer{}
	err = template.Execute(body, this)
	if err != nil {
		return nil, err
	}
	return format.Source(body.Bytes())
}

func (this *ModelMeta) FindField() error {
	queryString, err := this.Db.SQL(tableColumnsSql, this.DbName, this.TableName).QueryString()
	if err != nil {
		return err
	}
	metas := make([]FieldMeta, 0, len(queryString))
	for _, elem := range queryString {
		metas = append(metas, FieldMeta{
			Name:       elem["column_name"],
			FieldType:  elem["file_type"],
			IsNullable: elem["is_nullable"],
			Tags:       this.Tags,
		})
	}
	this.Fields = metas
	return nil
}

func (this *ModelMeta) GetCreateSql() (string, error) {
	result, err := this.Db.SQL(fmt.Sprintf("SHOW CREATE TABLE %s", this.TableName)).QueryString()
	if err != nil {
		return "", err
	}
	var sql = ""
	for _, e := range result {
		str, isExist := e[CreateTable]
		if isExist {
			sql = str
		}
	}
	if sql == "" {
		return "", errors.New(fmt.Sprintf("can not fond %s create table sql", this.TableName))
	}
	return sql, nil
}

func (this *ModelMeta) HasTimeFiled() bool {
	for _, e := range this.Fields {
		if strings.Compare(e.GetGoType(), "time.Time") == 0 {
			return true
		}
	}
	return false
}
