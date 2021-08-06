package psql

import (
	"cat-service/db/psql/models"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type Catalog struct {
	conn *pgx.Conn
}

func NewCatalog(conn *pgx.Conn) *Catalog {
	return &Catalog{
		conn,
	}
}

func (s *Catalog) Save(cat models.Cat) (err error) {
	_, err = s.conn.Query(context.Background(),
		"insert into cats(id, name, breed, color, age, price) values ($1, $2, $3, $4, $5, $6);",
		cat.ID, cat.Name, cat.Breed, cat.Color, cat.Age, cat.Price)
	if err != nil {
		log.Print("err ", err)
		return err
	}

	return nil
}

func (s *Catalog) Delete(id uuid.UUID) error {
	_, err := s.conn.Query(context.Background(),
		"delete from cats where id=$1",
		id)
	if err != nil {
		log.Print("delete ", err)
		return err
	}

	return nil
}

func (s *Catalog) Update(id uuid.UUID, newCat models.Cat) (err error) {
	_, err = s.conn.Query(context.Background(),
		"update cats set name=$2, breed=$3, color=$4, age=$5, price=$6 where id=$1",
		newCat.ID, newCat.Name, newCat.Breed, newCat.Color, newCat.Age, newCat.Price)
	if err != nil {
		return err
	}

	return nil
}

func (s *Catalog) Get(id uuid.UUID) (cat models.Cat, err error) {
	err = s.conn.QueryRow(context.Background(), "select * from cats where id=$1", id).Scan(&cat.ID, &cat.Name, &cat.Breed, &cat.Color, &cat.Age, &cat.Price)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func (s *Catalog) GetAll() (cats []models.Cat, err error) {
	rows, err := s.conn.Query(context.Background(), "select * from cats")
	if err != nil {
		return []models.Cat{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var cat models.Cat
		err = rows.Scan(
			&cat.ID, &cat.Name, &cat.Breed, &cat.Color, &cat.Age, &cat.Price,
		)
		if err != nil {
			return cats, err
		}

		cats = append(cats, cat)
	}
	return cats, nil
}
