
# 开发日志
- 数据结构设计
  - 网址：https://dbdiagram.io/d/60d5769bdd6a5971481efd29

- migration cli
- 安装:brew install golang-migrate 
  - 适用于不同的数据库Postgres、Mysql、Mongodb..
  - 使用语法：
        create [-ext E] [-dir D] [-seq] [-digits N] [-format] NAME
                    Create a set of timestamped up/down migrations titled NAME, in directory D with extension E.
                    Use -seq option to generate sequential up/down migrations with N digits.
                    Use -format option to specify a Go time format string.
        goto V       Migrate to version V
        up [N]       Apply all or N up migrations
        down [N]     Apply all or N down migrations
        drop         Drop everything inside database
        force V      Set version V but don't run migration (ignores dirty state)
        version      Print current migration version
  主要使用的是create、up和down, 查看版本migrate -version

  - 初始化数据库命令: 
    - 新建db目录： mkdir -p db/migration
    - 创建两个迁移文件：migrate create -ext sql -dir db/migration -seq init_schema， 初始化版本号为1
    - migrate up 主要是基于当前的数据结构图和migration的up文件，向前更改，如果我们想恢复之前的同一个版本结构需要运行同一个版本的down文件: migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable" -verbose up

    - migrate up 会根据migration文件中的版本号依次更新1_up.sql->2_up.sql->...，down原理相同

- 数据ORM操作
  - 操作数据库Database/sql、SQLX、SQLC和GORM，推荐使用SQLC和SQLX
  - 网址: https://sqlc.dev/
  - 我们使用sqlc来自动生成CRUD的代码
  - 安装cli：brew install sqlc
  - 验证sqlc版本；sqlc version; 查看帮助：sqlc help
    Available Commands:
        compile     Statically check SQL for syntax and type errors
        generate    Generate Go code from SQL
        help        Help about any command
        init        Create an empty sqlc.yaml settings file
        version     Print the sqlc version number

  - 使用sqlc init初始化sqlc的yaml配置文件。
  - 编写query下面的的查询语句，最后使用sqlc generate， 生成代码

- 测试程序操作pg时，需要安装数据库驱动
  - go get github.com/lib/pq

- 安装测试环境
  - 网址：go get github.com/stretchr/testify
  - 进入测试文件所在目录测试命令：go test -run TestCreateAccount ./
  https://github.com/techschool/simplebank