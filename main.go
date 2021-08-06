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
	// logger, err := zap.NewDevelopment()
	// if err != nil {
	// 	log.Fatal("could not start zap logger\n")
	// }

	// defer logger.Sync()

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@0.0.0.0:5433/catalog")
	if err != nil {
		log.Println("could not connect to the psql ", err)
		zap.L().Error("could not connect to the psql ", zap.Error(err))
		os.Exit(1)
	}

	log.Println("after connect ", err)
	err = conn.Ping(context.Background())
	if err != nil {
		log.Println("could not connect to the psql ", err)
		zap.L().Error("could not connect to the psql ", zap.Error(err))
	}

	log.Println("after ping ", err)
	defer conn.Close(context.Background())

	catalog := psql.NewCatalog(conn)
	service := &Service{
		catalog,
	}
	e := echo.New()
	e.POST("/cats", service.addCat)
	e.GET("/cats", service.getAllCats)
	e.GET("/cats/:id", service.getCat)
	e.PUT("/cats/:id", service.updateCat)
	e.DELETE("/cats/:id", service.deleteCat)
	e.GET("/cats/get-rand-cat", getRandCat)
	// e.GET()
	if err := e.Start(":8081"); err != nil {
		log.Println("could not connect to the psql ", err)
		zap.L().Error("could not start server\n", zap.Error(err))
	}

	log.Println("after server ", err)

}
