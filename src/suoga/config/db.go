package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DBConfig *Database //nolint:gochecknoglobals

// 连接配置
type Database struct {
	DB_Driver string //nolint:goimports,gofmt
	DB_Url    string
}

func loadDBDefaultConfig() {
	ViperConfig.SetDefault("DB_Driver", "postgres")
	ViperConfig.SetDefault("DB_Url", "postgresql://root:123456@localhost:5432/suoga?sslmode=disable")
}

func SetUpDB() *sql.DB {
	var (
		err error
		DB  *sql.DB
	)
	loadDBDefaultConfig()
	ViperConfig.Unmarshal(&DBConfig) //nolint:errcheck

	DB, err = sql.Open(DBConfig.DB_Driver, DBConfig.DB_Url)
	if err != nil {
		fmt.Println("[+] Get db connecting error: ", err.Error())
		os.Exit(0)

	}
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(100)
	DB.SetConnMaxLifetime(5 * time.Minute)
	return DB
}
