package postgres_config

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormClient(dsn string, logger *zap.Logger) *gorm.DB {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("[Postgres] database connection failed", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = sqlDB.Ping()
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("database connection established")
	return db
}
