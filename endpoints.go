package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"cat-service/db/psql"
	"cat-service/db/psql/models"
)

type Service struct {
	catalog *psql.Catalog
}

func (s *Service) addCat(c echo.Context) error {
	cat := &models.Cat{}

	if err := (&echo.DefaultBinder{}).BindBody(c, &cat); err != nil {
		log.Print("bind body ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.catalog.Save(*cat)
	if err != nil {
		log.Print("save ", err)
		return c.JSON(http.StatusInternalServerError, "error saving cat(")
	}

	return c.JSON(http.StatusCreated, nil)
}

func (s *Service) deleteCat(c echo.Context) error {
	var id uuid.UUID
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &id); err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.catalog.Delete(id)
	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "deleted")

}

func (s *Service) updateCat(c echo.Context) error {
	req := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, req); err != nil {
		log.Print("bind body: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.catalog.Update(req.ID, *req)
	if err != nil {
		log.Print("update: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *Service) getCat(c echo.Context) error {
	var id uuid.UUID
	q := c.Param("id")
	id, err := uuid.Parse(q)
	log.Print(id)
	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	cat, err := s.catalog.Get(id)
	if err != nil {
		log.Print("get catalog", err)
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cat)
}
func (s *Service) getAllCats(c echo.Context) error {
	var cats []models.Cat
	cats, err := s.catalog.GetAll()
	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cats)
}

func getRandCat(c echo.Context) error {
	cat := psql.RandCat()
	return c.JSON(200, cat)
}
