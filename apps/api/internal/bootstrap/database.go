package bootstrap

import (
	"fmt"

	"authentication-project-exam/internal/config"

	pgAdapter "authentication-project-exam/internal/adapter/outbound/persistence/postgres"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: pgAdapter.NewZapGormLogger(logger).LogMode(gormlogger.Info),
	})
}
