package orm

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DBName       string
	Type         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

func (c Config) GetMysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&", c.Username, c.Password, c.Host, c.Port, c.DBName)
}

func (c Config) NewMysql() (*gorm.DB, error) {
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 100
	}
	if c.MaxLifetime == 0 {
		c.MaxLifetime = 3600
	}
	dsn := c.GetMysqlDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewGormLogger(logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Error,
		}, service.DevMode),
	})
	if err != nil {
		return nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxIdleConns(c.MaxIdleConns)
	sdb.SetMaxOpenConns(c.MaxOpenConns)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(c.MaxLifetime))

	err = db.Use(NewTracePlugin())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c Config) MustNewMysql() *gorm.DB {
	db, err := c.NewMysql()
	if err != nil {
		panic(err)
	}
	return db
}
