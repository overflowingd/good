package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const MysqlDefaultPort = 3306
const MysqlDefaultCharset = "utf8mb4"
const MysqlDefaultTimezone = "UTC"

var (
	ErrMysqlInitialization = errors.New("db:mysql:initialization_failed")
)

type DSN interface {
	fmt.Stringer
}

type MysqlConnectionOptions struct {
	DSN
	DbName          string
	DbCharset       *string
	DbHost          string
	DbPort          *int
	DbUser          string
	DbPassword      string
	DbTimezone      *string
	MultiStatements bool
}

func (mcopts *MysqlConnectionOptions) DefaultPort() int {
	return MysqlDefaultPort
}

func (mcopts *MysqlConnectionOptions) DefaultCharset() string {
	return MysqlDefaultCharset
}

func (mcopts *MysqlConnectionOptions) DefaultTimezone() string {
	return MysqlDefaultTimezone
}

func (mcopts *MysqlConnectionOptions) DefaultTimeout() string {
	return "1m"
}

func (mcopts *MysqlConnectionOptions) String() string {
	port := mcopts.DefaultPort()
	if nil != mcopts.DbPort {
		port = *mcopts.DbPort
	}

	charset := mcopts.DefaultCharset()
	if nil != mcopts.DbCharset {
		charset = *mcopts.DbCharset
	}

	timezone := mcopts.DefaultTimezone()
	if nil != mcopts.DbTimezone {
		timezone = *mcopts.DbTimezone
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s&timeout=%s&multiStatements=%t",
		mcopts.DbUser,
		mcopts.DbPassword,
		mcopts.DbHost,
		port,
		mcopts.DbName,
		charset,
		timezone,
		mcopts.DefaultTimeout(),
		mcopts.MultiStatements,
	)
}

type MysqlOptionsProvider func(multiStatements bool) *MysqlConnectionOptions

type ConnPoolOptions interface {
	MaxIdleConns() int
	MaxOpenConns() int
	ConnMaxLifetime() time.Duration
	ConnMaxIdleTime() time.Duration
}

func NewMysql(dsn DSN) (gorm.Dialector, error) {
	dialector := mysql.Open(dsn.String())
	if nil == dialector {
		return nil, ErrMysqlInitialization
	}

	return dialector, nil
}

func NewConnection(
	dialector gorm.Dialector,
	connPoolOptions ConnPoolOptions,
	config *gorm.Config,
) (*gorm.DB, *sql.DB, error) {
	if nil == config {
		config = new(gorm.Config)
		config.NamingStrategy = schema.NamingStrategy{
			SingularTable: true,
		}
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(dialector, config)
	if nil != err {
		return nil, nil, err
	}

	sql, err := db.DB()
	if nil != err {
		return nil, nil, err
	}

	if nil != connPoolOptions {
		sql.SetMaxIdleConns(connPoolOptions.MaxIdleConns())
		sql.SetMaxOpenConns(connPoolOptions.MaxOpenConns())
		sql.SetConnMaxLifetime(connPoolOptions.ConnMaxLifetime())
		sql.SetConnMaxIdleTime(connPoolOptions.ConnMaxIdleTime())
	}

	return db, sql, nil
}
