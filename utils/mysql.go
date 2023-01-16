package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

var DB *gorm.DB

func ReadMysqlConfig() string {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("No caller information")
		return ""
	}
	fileName = path.Join(path.Dir(fileName), "../config/mysql.json")

	viper.SetConfigFile(fileName)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Read config file error: %s\n", err)
		return ""
	}

	config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.Get("user"),
		viper.Get("passwd"),
		viper.Get("host"),
		viper.Get("port"),
		viper.Get("database"))

	return config
}

func InitMysql() {
	config := ReadMysqlConfig()

	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	DB, _ = gorm.Open(mysql.Open(config), &gorm.Config{
		Logger: newLogger,
	})

	sqlDB, err := DB.DB()
	if err != nil {
		panic(err.Error())
	}

	sqlDB.SetMaxIdleConns(viper.GetInt("poolConfig.maxIdleConns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("poolConfig.maxOpenConns"))
	sqlDB.SetConnMaxIdleTime(time.Minute)
	sqlDB.SetConnMaxLifetime(time.Minute)
}
