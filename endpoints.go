package main

import (
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
		return c.JSON(http.StatusInternalServerError, err)
	}

	created, err := s.catalog.Save(*cat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, created)
}

func (s *Service) deleteCat(c echo.Context) error {
	var id uuid.UUID
	if err := (&echo.DefaultBinder{}).BindBody(c, &id); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.catalog.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "deleted")

}

func (s *Service) updateCat(c echo.Context) error {
	req := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, req); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	newCat, err := s.catalog.Update(req.ID, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, newCat)
}

func (s *Service) getCat(c echo.Context) error {
	var id uuid.UUID
	if err := (&echo.DefaultBinder{}).BindBody(c, &id); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cat, err := s.catalog.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cat)
}
func (s *Service) getAllCats(c echo.Context) error {
	var cats []models.Cat
	if err := (&echo.DefaultBinder{}).BindBody(c, &cats); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cats, err := s.catalog.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cats)
}
