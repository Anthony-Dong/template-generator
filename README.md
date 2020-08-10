## go-orm

用于生成Go的Model文件，数据库操作太过于麻烦，依靠工具可以直接生成model对象，默认使用的是xorm映射。

> ​	首先声明 ： 公司内部的数据库表字段全部是`下划线模式`，表名称全部是`下划线模式`

- 支持生成xorm的model对象
- 支持生成dao对象
- 支持生成dto对象（time.Time 转化成 int64时间搓）
- 支持外部配置文件，防止重复输入配置文件，默认配置文件在 `go-orm-config.json`，这个优先级低于 命令，如果你命令传入-config，显示申明配置文件，那么它的优先级最高。

需要支持Go mod ，所以版本最好1.11以上

### 下载：

```go
go get -u github.com/anthony-dong/template-generator/cmd/orm
```

### 快速开始

> ​	以xxl_job的表 为例子

配置文件如下:

```json
{
  "db_type": "mysql",
  "tags": [
    "xorm"
  ],
  "db_name": "tests",
  "db_host": "localhost",
  "db_port": 3306,
  "db_user_name": "root",
  "db_password": "12345",
  "db_charset": "utf8"
}
```

```go
~/go/code/framework/tempalte-generator (master*) % bin/go-orm -d=xxl_job  -host=localhost  -port=3306  -u=root -p=123456  -dir=/data/tmp
[GEN-INFO] save path: /data/tmp
[xorm] [info]  2020/07/11 14:29:55.130611 PING DATABASE mysql
[GEN-INFO] save xxl_job_group model success, path=/data/tmp/model/xxl_job_group.go
[GEN-INFO] save xxl_job_user model success, path=/data/tmp/model/xxl_job_user.go
[GEN-INFO] save xxl_job_user dao success, path= /data/tmp/dao/xxl_job_user_dao.go
[GEN-INFO] save xxl_job_group dao success, path= /data/tmp/dao/xxl_job_group_dao.go
[GEN-INFO] save xxl_job_registry model success, path=/data/tmp/model/xxl_job_registry.go
[GEN-INFO] save xxl_job_lock model success, path=/data/tmp/model/xxl_job_lock.go
[GEN-INFO] save xxl_job_logglue model success, path=/data/tmp/model/xxl_job_logglue.go
[GEN-INFO] save xxl_job_log model success, path=/data/tmp/model/xxl_job_log.go
[GEN-INFO] save xxl_job_info model success, path=/data/tmp/model/xxl_job_info.go
[GEN-INFO] save xxl_job_registry dao success, path= /data/tmp/dao/xxl_job_registry_dao.go
[GEN-INFO] save xxl_job_lock dao success, path= /data/tmp/dao/xxl_job_lock_dao.go
[GEN-INFO] save xxl_job_log_report model success, path=/data/tmp/model/xxl_job_log_report.go
[GEN-INFO] save xxl_job_logglue dao success, path= /data/tmp/dao/xxl_job_logglue_dao.go
[GEN-INFO] save xxl_job_log dao success, path= /data/tmp/dao/xxl_job_log_dao.go
[GEN-INFO] save xxl_job_info dao success, path= /data/tmp/dao/xxl_job_info_dao.go
[GEN-INFO] save xxl_job_log_report dao success, path= /data/tmp/dao/xxl_job_log_report_dao.go
[GEN-INFO] save xxl_job dto success, path= /data/tmp/dto/xxl_job_dto.go
generate template finished
```

### 各个模块生成内容介绍

- model 对象

  > 生成注释，公司默认id就是主键，所以我们这次简单暴力，只要是id就是主键(后期改进)

```go
package model

/**
CREATE TABLE `xxl_job_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_name` varchar(64) NOT NULL COMMENT '执行器AppName',
  `title` varchar(12) NOT NULL COMMENT '执行器名称',
  `address_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '执行器地址类型：0=自动注册、1=手动录入',
  `address_list` varchar(512) DEFAULT NULL COMMENT '执行器地址列表，多地址逗号分隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4
*/

type XxlJobGroup struct {
	ID          int    `xorm:"pk autoincr id"`
	AppName     string `xorm:"app_name"`
	Title       string `xorm:"title"`
	AddressType uint8  `xorm:"address_type"`
	AddressList string `xorm:"address_list"`
}

func (*XxlJobGroup) TableName() string {
	return "xxl_job_group"
}
```

- dto 对象

> ​	 json 字段默认是表字段名称，因为表的字段就是下划线模式，同时回生成注释，没有注释的不会添加注释

```go
package dto

type XxlJobGroupDto struct {
	ID          int    `json:"id"`
	AppName     string `json:"app_name"`     //执行器AppName
	Title       string `json:"title"`        //执行器名称
	AddressType uint8  `json:"address_type"` //执行器地址类型：0=自动注册、1=手动录入
	AddressList string `json:"address_list"` //执行器地址列表，多地址逗号分隔
}

type XxlJobLogglueDto struct {
	ID         int    `json:"id"`
	JobID      int    `json:"job_id"`      //任务，主键ID
	GlueType   string `json:"glue_type"`   //GLUE类型
	GlueSource string `json:"glue_source"` //GLUE源代码
	GlueRemark string `json:"glue_remark"` //GLUE备注
	AddTime    int64  `json:"add_time"`
	UpdateTime int64  `json:"update_time"`
}
```

- dao对象

  > ​	生成了 crud各种操作，其中引入不了的包，属于公司内部包。。。

```go
package dao

import (
	"context"
)

type xxlJobGroupDao struct {
}

func NewXxlJobGroupDao() *xxlJobGroupDao {
	return &xxlJobGroupDao{}
}

var (
	_XxlJobGroup = new(model.XxlJobGroup)
)

func (*xxlJobGroupDao) TableName() string {
	return _XxlJobGroup.TableName()
}

func (this *xxlJobGroupDao) GetById(ctx context.Context, id uint64) (*model.XxlJobGroup, cerror.Cerror) {
	result := &model.XxlJobGroup{}
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

func (this *xxlJobGroupDao) GetByIds(ctx context.Context, ids []uint64) ([]model.XxlJobGroup, cerror.Cerror) {
	list := make([]model.XxlJobGroup, 0)
	err := db.SlaveDb().In("id", ids).Cols("*").Find(&list)
	if err != nil {
		logger.Errorc(ctx, "[GetByIds] err,ids=%v", ids)
		return nil, exception.DbExecError(err)
	}
	return list, nil
}

func (this *xxlJobGroupDao) DeleteById(ctx context.Context, session *xorm.Session, id uint64) (int64, cerror.Cerror) {
	effectRow, err := session.Where("id=?", id).Delete(this)
	if err != nil {
		logger.Errorc(ctx, "[DeleteById] err,id=%d", id)
		return 0, exception.DbDeleteError(err)
	}
	return effectRow, nil
}

func (this *xxlJobGroupDao) UpdateById(ctx context.Context, session *xorm.Session, id uint64, params map[string]interface{}) (int64, cerror.Cerror) {
	effectRow, err := session.Table(this).Where("id=?", id).Update(params)
	if err != nil {
		logger.Errorc(ctx, "[UpdateById] err,id=%d,params=%+v", id, params)
		return 0, exception.DbUpdateError(err)
	}
	return effectRow, nil
}

func (this *xxlJobGroupDao) UpdateByIds(ctx context.Context, session *xorm.Session, ids []uint64, params map[string]interface{}) (int64, cerror.Cerror) {
	effectRow, err := session.Table(this).In("id", ids).Update(params)
	if err != nil {
		logger.Errorc(ctx, "[UpdateByIds] err,ids=%d,params=%+v", ids, params)
		return 0, exception.DbUpdateError(err)
	}
	return effectRow, nil
}

func (this *xxlJobGroupDao) SaveOne(ctx context.Context, session *xorm.Session, param *model.XxlJobGroup) (int64, cerror.Cerror) {
	effectRow, err := session.InsertOne(param)
	if err != nil {
		logger.Errorc(ctx, "[SaveOne] err,param=%+v", param)
		return 0, exception.DbInsertError(err)
	}
	return effectRow, nil
}

func (this *xxlJobGroupDao) SaveMany(ctx context.Context, session *xorm.Session, params *[]model.XxlJobGroup) (int64, cerror.Cerror) {
	effectRow, err := session.Insert(params)
	if err != nil {
		logger.Errorc(ctx, "[SaveMany] err,params=%+v", params)
		return 0, exception.DbInsertError(err)
	}
	return effectRow, nil
}

func (this *xxlJobGroupDao) GetAll(ctx context.Context) ([]model.XxlJobGroup, cerror.Cerror) {
	list := make([]model.XxlJobGroup, 0)
	err := db.SlaveDb().Cols("*").Find(&list)
	if err != nil {
		logger.Errorc(ctx, "[GetAll] err")
		return nil, exception.DbExecError(err)
	}
	return list, nil
}
```





## go-build

用于快速构建项目:

首先本克隆项目，到本地，然后

- 构建脚本

```go
go get -u github.com/anthony-dong/template-generator/cmd/build
```

- 查看命令

```shell
bin/go-build -h
```

- 快速开始

```go
bin/go-build -dir=/data/tmp -mod=city-demo -git=git@gihub.com:Anthony-Dong/template.git
```

> // 快速构建：
> -dir 项目本地位置
> -mod 你的项目名称：go mod  的名称 ，本地版本不得 低于go 1.11
> -git 是我的模版地址，会告诉你

## 