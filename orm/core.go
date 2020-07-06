package orm

import (
	"fmt"
	"github.com/anthony-dong/template-generator/logger"
	"github.com/go-xorm/xorm"
	"os"
	"strings"
	"sync"
	"text/template"
)

type DbType string

// type
const (
	Mysql = DbType("mysql")
)

// template name
const (
	daoTemplateName   = "daoTemplateName"
	modelTemplateName = "modelTemplateName"
)

// other
const (
	sep        = string(os.PathSeparator)
	fileSuffix = ".go"
)

type Template interface {
	Run(template *template.Template) ([]byte, error)
}

type DbMeta interface {
	GetTables() ([]string, error)
}

type Config struct {
	wg            sync.WaitGroup
	daoTemplate   *template.Template
	modelTemplate *template.Template

	// db
	DbType     DbType
	DbName     string
	DbPort     int
	DbHost     string
	DbUserName string
	DbPassword string
	DbCharset  string

	// cnn
	engine *xorm.Engine

	// save
	SaveFile string

	// table
	TableNames []string

	// template
	GeneratorModel   bool
	GeneratorDao     bool
	DaoPackageName   string
	ModelPackageName string
	Tags             []string
}

func FillDNS(userName string, password string, host string, port int, dbName string, charset string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		userName,
		password,
		host,
		port,
		dbName,
		charset,
	)
}

func GetDbType(str string) DbType {
	lower := strings.ToLower(str)
	switch DbType(lower) {
	case Mysql:
		return Mysql
	default:
		logger.FatalF("not support type")
		return Mysql
	}
}
