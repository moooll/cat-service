package main

import (
	"cat-crud/models"
	"cat-crud/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Service struct {
	Catalog storage.Catalog
}

func (s *Service) addCat(c echo.Context) error {
	type Cat struct {
		Name  string  `json:"name"`
		Breed string  `json:"breed"`
		Color string  `json:"color"`
		Age   float32 `json:"age"`
		Price float32 `json:"price"`
	}
	cat := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &cat); err != nil {
		return err
	}

	s.Catalog.Save(*cat)
	return nil
}

func (s *Service) deleteCat(c echo.Context) error {
	var id uuid.UUID
	if err := (&echo.DefaultBinder{}).BindBody(c, &id); err != nil {
		return err
	}

	s.Catalog.Delete(id)
	return nil
}

// updateCat receives requests' bodies following the pattern:
// {id:id, []fields{field_name:new_value}}
func (s *Service) updateCat(c echo.Context) error {
	type field struct {
		name string
		value string
	}
	type request struct {
		id uuid.UUID
		fields []field 
	}
	req := &request{}
	if err := (&echo.DefaultBinder{}).BindBody(c, req); err != nil {
		return err
	}

	
	return nil
}

func (s *Service) getCat(c echo.Context) error {

	return nil
}
func (s *Service) getAllCats(c echo.Context) error {
	return nil
}

func (s *Service) getCatName(c echo.Context) error {
	return nil
}

func (s *Service) updateCatName(c echo.Context) error {
	return nil
}
