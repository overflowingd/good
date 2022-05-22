package good

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	DriverMysql = "mysql"
)

const (
	MigrationsTable     = "migrations"
	MigrationsDirScheme = "file://"
)

func NewMysqlMigrations(sqldb *sql.DB, dir string) (*migrate.Migrate, error) {
	config := new(mysql.Config)
	config.MigrationsTable = MigrationsTable

	driver, err := mysql.WithInstance(sqldb, config)
	if nil != err {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(
		MigrationsDirScheme+dir,
		DriverMysql,
		driver,
	)
}
