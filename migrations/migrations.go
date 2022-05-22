package migrations

import (
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/overflowingd/good/db"
)

type Hanler func(opts *db.MysqlConnectionOptions, migrationsDir string) error

func MigrateOnStartup() bool {
	_, exists := os.LookupEnv("DB_MIGRATE_ON_STARTUP")
	return exists
}

func Up(m *migrate.Migrate) error {
	return m.Up()
}

func SufficientError(err error) bool {
	switch err {
	case migrate.ErrNilVersion:
	case migrate.ErrNoChange:
		return false
	}

	return true
}
