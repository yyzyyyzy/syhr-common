package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
	ServiceMode   string
}

func NewGormLogger(c logger.Config, serviceMode string) *GormLogger {
	return &GormLogger{
		SlowThreshold: c.SlowThreshold,
		LogLevel:      c.LogLevel,
		ServiceMode:   serviceMode,
	}
}

var _ logger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = lev
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logx.WithContext(ctx).Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logx.WithContext(ctx).Errorf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logx.WithContext(ctx).Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	logFields := []logx.LogField{
		logx.Field("sql", sql),
		logx.Field("time", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
		logx.Field("rows", rows),
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.WithContext(ctx).Infow("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, logx.Field("error", err.Error()))
			logx.WithContext(ctx).Infow("Database Err", logFields...)
		}
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logFields = append(logFields, logx.Field("Database Slow Log", l.SlowThreshold))
		logx.WithContext(ctx).Sloww("Slow Log", logFields...)
	}
	if l.ServiceMode != service.ProMode {
		logx.WithContext(ctx).Infow("Database Query", logFields...)
	}
}
