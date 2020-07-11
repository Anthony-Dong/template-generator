package main

import (
	"flag"
	"fmt"
	"github.com/anthony-dong/template-generator/logger"
	"github.com/anthony-dong/template-generator/orm"
	"os"
	"strings"
)

const (
	generatorVersion = "1.0.0"
	successExit      = -1
)

var (
	help             bool
	showVersion      bool
	destDir          string
	modelPackageName string
	daoPackageName   string
	dtoPackageName   string
	tableNames       strFlags
	dbType           string
	tags             strFlags
	dbName           string
	dbHost           string
	dbPort           int
	dbUserName       string
	dbPassword       string
	dbCharset        string
)

type strFlags []string

func (i *strFlags) String() string {
	return "table names"
}

func (i *strFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&showVersion, "v", false, "generator version")
	flag.StringVar(&dbType, "db_type", "mysql", "database type, eg: -db_type=mysql")
	flag.StringVar(&dbName, "d", "", "database name, eg: -d=xorm")
	flag.StringVar(&dbHost, "host", "localhost", "database host, eg: -port=localhost")
	flag.IntVar(&dbPort, "port", 3306, "database port, eg: -port=3306")
	flag.StringVar(&dbUserName, "u", "root", "database username, eg: -u=root")
	flag.StringVar(&dbPassword, "p", "123456", "database password, eg: -p=123456")
	flag.Var(&tableNames, "t", "database table names default all tables , eg: -t=class -t=user")
	flag.StringVar(&dbCharset, "charset", "utf8", "database table names, eg: -charset=utf8")
	flag.Var(&tags, "tag", "modle tag names support add many tags,default xorm, eg: -tag=xorm -tag=json")
	flag.StringVar(&destDir, "dir", "./tmp", "generated directory default ./tmp, eg: -dir=./tmp")
	flag.StringVar(&modelPackageName, "model_package", "model", "package name default model, eg:-model_package=model")
	flag.StringVar(&daoPackageName, "dao_package", "dao", "package name default dao, eg:-dao_package=dao")
	flag.StringVar(&dtoPackageName, "dto_package", "dto", "package name default dao, eg:-dto_package=dao")
}

func main() {
	// 解析输入
	initFlag()
	// 初始化属性
	config := initProperties()
	err := config.Generator()
	if err != nil {
		logger.FatalF("generate fail,err=%s", err)
	}
	fmt.Println("generate template finished")
}

// 初始化属性
func initProperties() *orm.Config {
	config := new(orm.Config)
	config.DbType = orm.GetDbType(dbType)
	config.DbName = dbName
	config.DbUserName = dbUserName
	config.DbPassword = dbPassword
	config.DbHost = dbHost
	config.DbPort = dbPort
	config.DbCharset = dbCharset
	config.GeneratorModel = true
	config.GeneratorDao = true
	config.DaoPackageName = daoPackageName
	config.ModelPackageName = modelPackageName
	config.DtoPackageName = dtoPackageName
	config.Tags = addTag(tags)
	config.SaveFile = destDir
	return config
}

// 指示
func initFlag() {
	flag.Parse()
	if help {
		fmt.Println(`generator version: anthony/1.0.0
Usage: generator -host=localhost -port=3306 -d=urban_v -u=root -p=123456 -t=class -t=student -tag=xorm -dir=./tmp
Option:`)
		flag.PrintDefaults()
		os.Exit(successExit)
	}

	if showVersion {
		fmt.Printf("generator version: anthony/%s", generatorVersion)
		os.Exit(successExit)
	}
}

func addTag(tags strFlags) []string {
	str := "xorm"
	result := make([]string, 0)
	result = append(result, str)
	for _, e := range tags {
		if strings.Compare(str, e) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}
