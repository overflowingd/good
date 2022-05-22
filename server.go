package good

import (
	"database/sql"
	"os"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/overflowingd/good/db"
	"github.com/overflowingd/good/migrations"
	"github.com/overflowingd/good/router"
	"gorm.io/gorm"
)

type Server interface {
	ListenAndServe() (err error)
}

type AfterServerApplication func(conn *gorm.DB, rawsqlconn *sql.DB) error

func ServerApplication(
	buildRoutes router.RouteBuilder,
	ensureMigrations migrations.Hanler,
	migrationsDir string,
	provideMysqlOptions db.MysqlOptionsProvider,
	connPoolOptions db.ConnPoolOptions,
	after AfterServerApplication,
) (*gorm.DB, *sql.DB, Server, error) {
	if migrations.MigrateOnStartup() {
		if err := ensureMigrations(provideMysqlOptions(true), migrationsDir); err != nil {
			return nil, nil, nil, err
		}
	}

	mysql, err := db.NewMysql(provideMysqlOptions(false))
	if err != nil {
		return nil, nil, nil, err
	}

	r, err := NewRouter()
	if err != nil {
		return nil, nil, nil, err
	}

	buildRoutes(r)

	connection, rawconn, err := db.NewConnection(mysql, connPoolOptions, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := after(connection, rawconn); err != nil {
		defer rawconn.Close()
		return nil, nil, nil, err
	}

	return connection, rawconn, NewServer(r), nil
}

func NewServer(router *gin.Engine) Server {
	return endless.NewServer(os.ExpandEnv("$HOST:$PORT"), router)
}
