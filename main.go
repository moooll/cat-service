package main

import (
	"cat-crud/storage"
	
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("could not start zap logger\n")
	}

	defer logger.Sync()
	// init storage
	catStore := storage.NewCatalog()
	service := &Service{
		*catStore,
	}
	e := echo.New()

	e.POST("/cats/add", service.addCat)
	e.GET("/cats", service.getAllCats)
	e.GET("/cats/:id", service.getCat)
	e.GET("/cats/:name", service.getCatName)
	e.PUT("/cats/:id", service.updateCat)
	e.PUT("/cats/:name", service.updateCatName)
	e.DELETE("/cats/:id", service.deleteCat)

	if err := e.Start(":8081"); err != nil {
		zap.L().Error("could not start server\n")
	}
}
