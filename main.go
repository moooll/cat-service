package main

import (
	"cat-service/db/psql"
	"context"
	"os"

	"log"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("could not start zap logger\n")
	}

	defer logger.Sync()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		zap.L().Error("could not connect to the psql")
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	catalog := psql.NewCatalog(conn)
	service := &Service{
		catalog,
	}
	e := echo.New()
	e.POST("/cats/add", service.addCat)
	e.GET("/cats", service.getAllCats)
	e.GET("/cats/:id", service.getCat)
	// e.PUT("/cats/:id", service.updateCat)
	e.DELETE("/cats/:id", service.deleteCat)

	if err := e.Start(":8081"); err != nil {
		zap.L().Error("could not start server\n")
	}
}
