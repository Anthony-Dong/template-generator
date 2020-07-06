package orm

import (
	"fmt"
	"github.com/anthony-dong/template-generator/file"
	"github.com/anthony-dong/template-generator/logger"
	"github.com/anthony-dong/template-generator/orm/temp"
	"github.com/go-xorm/xorm"
	"path/filepath"
	"text/template"
)

func (config *Config) initTemplate() error {
	{
		parse, err := template.New(modelTemplateName).Parse(temp.Model)
		if err != nil {
			return err
		}
		config.modelTemplate = parse
	}
	{
		parse, err := template.New(daoTemplateName).Parse(temp.Dao)
		if err != nil {
			return err
		}
		config.daoTemplate = parse
	}
	return nil
}

func (config *Config) initDb() error {
	engine, err := xorm.NewEngine(string(config.DbType), FillDNS(config.DbUserName, config.DbPassword, config.DbHost, config.DbPort, config.DbName, config.DbCharset))
	if err != nil {
		return err
	}
	err = engine.Ping()
	if err != nil {
		return err
	}
	config.engine = engine
	if config.TableNames == nil || len(config.TableNames) == 0 {
		tableNames, err := NewMysqlMeta(config.DbName, engine).GetTables()
		if err != nil {
			return err
		}
		config.TableNames = tableNames
	}
	return nil
}

func (config *Config) Generator() error {
	dir, err := filepath.Abs(config.SaveFile)
	if err != nil {
		return err
	}
	config.SaveFile = dir
	logger.InfoF("save path: %s", dir)
	if config.ModelPackageName == "" {
		config.ModelPackageName = defaultModelName
	}
	if config.DaoPackageName == "" {
		config.DaoPackageName = defaultDaoName
	}

	err = config.initDb()
	if err != nil {
		return err
	}
	err = config.initTemplate()
	if err != nil {
		return err
	}
	for _, tableName := range config.TableNames {
		config.wg.Add(1)
		go func(config *Config, tableName string) {
			defer func() {
				if err := recover(); err != nil {
					logger.FatalF("generate err: %v", err)
				}
				config.wg.Done()
			}()
			switch config.DbType {
			case Mysql:
				func() {
					err := mysql(config, tableName)
					if err != nil {
						logger.FatalF("generate err: %v", err)
					}
				}()
			default:
				logger.FatalF("not support %s db type", config.DbType)
			}
		}(config, tableName)
	}
	config.wg.Wait()
	return nil
}

func mysql(config *Config, tableName string) error {
	if config.GeneratorModel {
		modelBody, err := NewModelMeta(config.DbName, tableName, config.ModelPackageName, config.engine, config.Tags).Run(config.modelTemplate)
		if err != nil {
			return err
		}
		saveFilePath := getSaveFileName(config.SaveFile, config.ModelPackageName, tableName)
		err = file.WriteFile(saveFilePath, modelBody)
		if err != nil {
			return err
		}
		logger.InfoF("save %s model success, path=%s", tableName, saveFilePath)
	}
	if config.GeneratorDao {
		dao := DaoMeta{
			TableName:    tableName,
			DaoPackage:   config.DaoPackageName,
			ModelPackage: config.ModelPackageName,
		}
		daoBody, err := dao.Run(config.daoTemplate)
		if err != nil {
			return err
		}
		saveFilePath := getSaveFileName(config.SaveFile, config.DaoPackageName, getDaoFileName(tableName))
		err = file.WriteFile(saveFilePath, daoBody)
		if err != nil {
			return err
		}
		logger.InfoF("save %s dao success, path= %s", tableName, saveFilePath)
	}
	return nil
}

func getDaoFileName(tableName string) string {
	return fmt.Sprintf("%s_dao", tableName)
}

func getSaveFileName(dir, packageName, fileName string) string {
	return fmt.Sprintf("%s%s%s%s%s%s", dir, sep, packageName, sep, fileName, fileSuffix)
}
