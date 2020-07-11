package orm

import (
	"fmt"
	"github.com/anthony-dong/template-generator/file"
	"github.com/anthony-dong/template-generator/logger"
	"github.com/anthony-dong/template-generator/orm/temp"
	"github.com/anthony-dong/template-generator/utils"
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
	{
		parse, err := template.New(dtoTemplateName).Funcs(map[string]interface{}{
			"Upper": utils.Marshal,
		}).Parse(temp.Dto)
		if err != nil {
			return err
		}
		config.dtoTemplate = parse
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
	if config.DtoPackageName == "" {
		config.DtoPackageName = defaultDtoPackageName
	}

	err = config.initDb()
	if err != nil {
		return err
	}
	err = config.initTemplate()
	if err != nil {
		return err
	}
	var dtoMetas *DtoMetas
	if config.GeneratorDto {
		dtoMetas = NewDtoMeta(config.DtoPackageName)
	}
	for _, tableName := range config.TableNames {
		config.wg.Add(1)
		go func(config *Config, tableName string, dtoMetas *DtoMetas) {
			defer func() {
				if err := recover(); err != nil {
					logger.FatalF("generate err: %v", err)
				}
				config.wg.Done()
			}()
			switch config.DbType {
			case Mysql:
				func() {
					err := mysql(config, tableName, dtoMetas)
					if err != nil {
						logger.FatalF("generate err: %v", err)
					}
				}()
			default:
				logger.FatalF("not support %s db type", config.DbType)
			}
		}(config, tableName, dtoMetas)
	}
	config.wg.Wait()

	if dtoMetas == nil {
		return nil
	}
	dtoFile := getSaveFileName(config.SaveFile, config.DtoPackageName, getDtoFileName(config.DbName))
	bytes, err := dtoMetas.Run(config.dtoTemplate)
	if err != nil {
		return err
	}
	err = file.WriteFile(dtoFile, bytes)
	if err != nil {
		return err
	}
	return nil
}

func mysql(config *Config, tableName string, dtoMetas *DtoMetas) error {
	if config.GeneratorModel {
		model := NewModelMeta(config.DbName, tableName, config.ModelPackageName, config.engine, config.Tags)
		modelBody, err := model.Run(config.modelTemplate)
		if err != nil {
			return err
		}
		saveFilePath := getSaveFileName(config.SaveFile, config.ModelPackageName, tableName)
		err = file.WriteFile(saveFilePath, modelBody)
		if err != nil {
			return err
		}
		logger.InfoF("save %s model success, path=%s", tableName, saveFilePath)
		if config.GeneratorDto {
			if dtoMetas != nil {
				dtoMetas.Append(model)
			}
		}
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

func getDtoFileName(dbName string) string {
	return fmt.Sprintf("%s_dto", dbName)
}

func getSaveFileName(dir, packageName, fileName string) string {
	return fmt.Sprintf("%s%s%s%s%s%s", dir, sep, packageName, sep, fileName, fileSuffix)
}
