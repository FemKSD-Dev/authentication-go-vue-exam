package postgres

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type ZapGormLogger struct {
	logger        *zap.Logger
	slowThreshold time.Duration
	logLevel      gormlogger.LogLevel
}

func NewZapGormLogger(logger *zap.Logger) gormlogger.Interface {
	return &ZapGormLogger{
		logger:        logger,
		slowThreshold: 500 * time.Millisecond,
		logLevel:      gormlogger.Warn,
	}
}

func (l *ZapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {

}

func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel < gormlogger.Warn {
		return
	}
	l.logger.Warn("gorm_warn",
		zap.String("message", msg),
	)
}

func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel < gormlogger.Error {
		return
	}
	l.logger.Error("gorm_error",
		zap.String("message", msg),
	)
}

func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.logLevel >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		l.logger.Error("gorm_query_failed",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)

	case elapsed > l.slowThreshold && l.logLevel >= gormlogger.Warn:
		l.logger.Warn("gorm_slow_query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
