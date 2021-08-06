package main

import (
	"cat-service/db/psql"

	"log"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	conn := psql.Connect().(*pgx.Conn)
	defer psql.Close(conn)

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
	if err := e.Start(":8081"); err != nil {
		log.Print("could not start server\n", err.Error())
	}
}
