package database

import (
	"fmt"
	"log"
	"os"
	"shadowify/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGromDatabase(config config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%v",
		config.Host,
		config.Port,
		config.User,
		config.DBName,
		config.Password,
		config.SSLMode,
	)
	logLevelMap := map[string]logger.LogLevel{
		"SILENT": logger.Silent,
		"ERROR":  logger.Error,
		"WARN":   logger.Warn,
		"INFO":   logger.Info,
	}
	ormLogLevel := "INFO"
	gormLogLevel := logger.Silent
	if logLv, found := logLevelMap[ormLogLevel]; found {
		gormLogLevel = logLv
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  gormLogLevel, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)
	gormConfig := &gorm.Config{Logger: newLogger, TranslateError: true}
	// https://github.com/go-gorm/postgres
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), gormConfig)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
