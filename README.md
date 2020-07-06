## go-build

用于快速构建项目:

首先本克隆项目，到本地，然后

- 构建脚本

```go
./build.sh  
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

## go-orm

用于生成Go的Model文件，数据库操作太过于麻烦，依靠工具可以直接生成model对象，默认使用的是xorm映射。

需要支持Go mod ，所以版本最好1.11以上

### 下载：

```go
go get -u github.com/anthony-dong/orm-generator/cmd/orm

或者 down下目录执行
./build.sh
```

### 快速开始

```go
:~/go/code/framework/tempalte-generator/bin % ./go-orm  -d=my_db  -host=xx.x.xxx.xx  -port=3306  -u=root -p=123456  -dir=/data/tmp -tag=xorm -tag=json
[GEN-INFO] save path: /data/tmp
[xorm] [info]  2020/07/06 21:24:10.325751 PING DATABASE mysql
[GEN-INFO] save peccancy_answer_user_submit model success, path=/data/tmp/model/peccancy_answer_user_submit.go
[GEN-INFO] save peccancy_answer_user_submit dao success, path= /data/tmp/dao/peccancy_answer_user_submit_dao.go
// ...
generate template finished
```

