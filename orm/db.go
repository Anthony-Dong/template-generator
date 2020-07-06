package orm

import (
	"github.com/go-xorm/xorm"
)

const (
	tableNameKey = "table_name"
)

const (
	tableNamesSql = `select table_name from information_schema.tables where table_schema = ? and table_type = 'base table';`
)

type mysqlMeta struct {
	Db     *xorm.Engine
	DbName string
}

func NewMysqlMeta(dbName string, db *xorm.Engine) DbMeta {
	return &mysqlMeta{
		Db:     db,
		DbName: dbName,
	}
}

func (this *mysqlMeta) GetTables() ([]string, error) {
	list, err := this.Db.SQL(tableNamesSql, this.DbName).QueryString()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(list))
	for _, e := range list {
		tableName, isExist := e[tableNameKey]
		if isExist {
			result = append(result, tableName)
		}
	}
	return result, nil
}
