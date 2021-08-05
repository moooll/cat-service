package main

import (
	"cat-service/db/psql"
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/ory/dockertest"
	"go.uber.org/zap"
)

var e *echo.Echo
var db *psql.TestConn
var catalog *psql.Catalog

func TestMain(m *testing.M) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("could not start zap logger\n")
	}

	defer logger.Sync()

	pool, err := dockertest.NewPool("")
	resource, err := pool.Run("postgres", "13.3", []string{"POSTGRES_PASSWORD=''", "POSTGRES_DB=catalog"})
	if err != nil {
		zap.L().Error("error starting container: ", zap.Error(err))
	}

	var conn *pgx.Conn
	if err = pool.Retry(func() error {
		conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			return err
		}

		err = conn.Ping(context.Background())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		zap.L().Error("error on retry ", zap.Error(err))
	}

	if err = pool.Purge(resource); err != nil {
		zap.L().Error("could not purge resource ", zap.Error(err))
	}

	db = psql.NewTestConn(conn)
	// catalog := psql.NewCatalog(conn)
}
