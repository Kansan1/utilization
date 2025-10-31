package config

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"go.uber.org/zap"
)

var DB *sql.DB

func InitDB() error {
	InitLogger()

	// 使用配置文件中的连接信息
	var err error
	DB, err = sql.Open("sqlserver", GetDSN())
	if err != nil {
		Logger.Fatal("连接数据库失败:", zap.Error(err))
		return err
	}

	err = DB.Ping()
	if err != nil {
		Logger.Fatal("连接数据库失败:", zap.Error(err))
		return err
	}

	Logger.Info("数据库连接成功")
	return nil
}
