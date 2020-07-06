package temp

var Dao = `
{{$model_mame:=.GetModelPathName}}
{{$dao_mame:=.GetStructName}}

package {{.DaoPackage}}

import (
	"context"
)

type {{$dao_mame}} struct {
}

func {{.GetNewStructFunc}}() *{{$dao_mame}} {
	return &{{$dao_mame}}{}
}

var (
	_{{.GetModelName}} = new({{$model_mame}})
)

func (*{{$dao_mame}}) TableName() string {
	return _{{.GetModelName}}.TableName()
}

func (this *{{$dao_mame}}) GetById(ctx context.Context, id uint64) (*{{$model_mame}}, cerror.Cerror) {
	result := &{{$model_mame}}{}
	isExist, err := db.SlaveDb().Where("id=?", id).Cols("*").Get(result)
	if err != nil {
		logger.Errorc(ctx, "[GetById] err,id=%d", id)
		return nil, exception.DbExecError(err)
	}
	if !isExist {
		return nil, nil
	}
	return result, nil
}

func (this *{{$dao_mame}}) GetByIds(ctx context.Context, ids []uint64) ([]{{$model_mame}}, cerror.Cerror) {
	list := make([]{{$model_mame}}, 0)
	err := db.SlaveDb().In("id", ids).Cols("*").Find(&list)
	if err != nil {
		logger.Errorc(ctx, "[GetByIds] err,ids=%v", ids)
		return nil, exception.DbExecError(err)
	}
	return list, nil
}

func (this *{{$dao_mame}}) DeleteById(ctx context.Context, session *xorm.Session, id uint64) (int64, cerror.Cerror) {
	effectRow, err := session.Where("id=?", id).Delete(this)
	if err != nil {
		logger.Errorc(ctx, "[DeleteById] err,id=%d", id)
		return 0, exception.DbDeleteError(err)
	}
	return effectRow, nil
}

func (this *{{$dao_mame}}) UpdateById(ctx context.Context, session *xorm.Session, id uint64, params map[string]interface{}) (int64, cerror.Cerror) {
	effectRow, err := session.Table(this).Where("id=?", id).Update(params)
	if err != nil {
		logger.Errorc(ctx, "[UpdateById] err,id=%d,params=%+v", id, params)
		return 0, exception.DbUpdateError(err)
	}
	return effectRow, nil
}

func (this *{{$dao_mame}}) UpdateByIds(ctx context.Context, session *xorm.Session, ids []uint64, params map[string]interface{}) (int64, cerror.Cerror) {
	effectRow, err := session.Table(this).In("id", ids).Update(params)
	if err != nil {
		logger.Errorc(ctx, "[UpdateByIds] err,ids=%d,params=%+v", ids, params)
		return 0, exception.DbUpdateError(err)
	}
	return effectRow, nil
}

func (this *{{$dao_mame}}) SaveOne(ctx context.Context, session *xorm.Session, param *{{$model_mame}}) (int64, cerror.Cerror) {
	effectRow, err := session.InsertOne(param)
	if err != nil {
		logger.Errorc(ctx, "[SaveOne] err,param=%+v", param)
		return 0, exception.DbInsertError(err)
	}
	return effectRow, nil
}

func (this *{{$dao_mame}}) SaveMany(ctx context.Context, session *xorm.Session, params *[]{{$model_mame}}) (int64, cerror.Cerror) {
	effectRow, err := session.Insert(params)
	if err != nil {
		logger.Errorc(ctx, "[SaveMany] err,params=%+v", params)
		return 0, exception.DbInsertError(err)
	}
	return effectRow, nil
}

func (this *{{$dao_mame}}) GetAll(ctx context.Context) ([]{{$model_mame}}, cerror.Cerror) {
	list := make([]{{$model_mame}}, 0)
	err := db.SlaveDb().Cols("*").Find(&list)
	if err != nil {
		logger.Errorc(ctx, "[GetAll] err")
		return nil, exception.DbExecError(err)
	}
	return list, nil
}
`

var Model = `
package {{.PackageName}}

{{if .HasTimeFiled}}
import (
	"time"
)
{{end}}

/**
{{.GetCreateSql}}
*/

type {{.ModelName}} struct {
    {{range .Fields}}
    {{.GetGoField}} {{.GetGoType}} {{if .GetTag}}{{.GetTag}}{{end}}{{end}}
}

func (*{{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}
`
