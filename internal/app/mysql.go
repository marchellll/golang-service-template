package app

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"

	drv "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host	 string
	Port	 string
	DBName	 string

	Username string
	Password string
}


func ConnectDB(logger *zap.Logger, config Config) *gorm.DB {
	dbConfig := drv.NewConfig()

	dbConfig.Addr = config.MySQLConfig.Host + ":" + config.MySQLConfig.Port
	dbConfig.DBName = config.MySQLConfig.DBName
	dbConfig.User = config.MySQLConfig.Username
	dbConfig.Passwd = config.MySQLConfig.Password
	dbConfig.Net = "tcp"

	gormDB, err := gorm.Open(mysql.Open(dbConfig.FormatDSN()), &gorm.Config{
		Logger: zapgorm2.New(logger),
	})

	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	return gormDB
}
