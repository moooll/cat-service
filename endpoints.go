package main

import (
	"cat-service/db/psql"
	"cat-service/db/psql/models"
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/ffjson/ffjson"
)

type Service struct {
	catalog *psql.Catalog
}

func (s *Service) addCat(c echo.Context) error {
	cat := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &cat); err != nil {
		return err
	}

	err := s.catalog.Save(context.Background(), *cat)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) deleteCat(c echo.Context) error {
	var id uuid.UUID
	if err := (&echo.DefaultBinder{}).BindBody(c, &id); err != nil {
		return err
	}

	s.catalog.Delete(id)
	return nil
}

// updateCat receives requests' bodies following the pattern:
// {id:id, []fields{field_name:new_value}}
// func (s *Service) updateCat(c echo.Context) error {
// 	// type field struct {
// 	// 	name  string
// 	// 	value string
// 	// }
// 	// type request struct {
// 	// 	id     uuid.UUID
// 	// 	fields []field
// 	// }
// 	// req := &request{}
// 	// if err := (&echo.DefaultBinder{}).BindBody(c, req); err != nil {
// 	// 	return err
// 	// }

// 	// psql.Update
// 	// return nil
// }

func (s *Service) getCat(c echo.Context) error {
	var id uuid.UUID
	if err := (&echo.DefaultBinder{}).BindBody(c, &id); err != nil {
		return err
	}

	cat, err := s.catalog.Get(id)
	if err != nil {
		return err
	}

	catB, err :=  ffjson.Marshal(cat) 
	if err != nil {
		return err
	}

	err = c.JSON(200, catB)
	if err != nil {
		return err
	}	

	return nil
}
func (s *Service) getAllCats(c echo.Context) error {
	var cats []models.Cat
	if err := (&echo.DefaultBinder{}).BindBody(c, &cats); err != nil {
		return err
	}

	cats, err := s.catalog.GetAll()
	if err != nil {
		return err
	}

	return nil
}

