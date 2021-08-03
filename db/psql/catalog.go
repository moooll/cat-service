package psql

import (
	"cat-service/db/psql/models"
	"context"

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

type stringCats struct {
	ID    string
	Name  string
	Breed string
	Color string
	Age   string
	Price string
}

func (s *Catalog) Save( ctx context.Context, cat models.Cat) error {
	_, err := s.conn.Query(context.Background(),
		"insert into cats(id, name, breed, color, age, price) values $id, $name, $breed, $color, $age, $price", 
		cat.ID, cat.Name, cat.Breed, cat.Color, cat.Age, cat.Price)
	if err != nil {
		return err
	}

	return nil
}

func (s *Catalog) Delete(id uuid.UUID) error {
	_, err := s.conn.Query(context.Background(), 
	"delete from cats where id=$1",
	id)
	if err != nil {
		return err
	}

	return nil
}

// todo
// func (s *Catalog) Update(id uuid.UUID, fields []interface{}) error {
// 	type field struct {
// 		name string
// 		value string
// 	}

// 	_, err := s.conn.Query(context.Background(), 
// 		"update cats set $1=$2 where id=$3", 
// 	)
// 	if err != nil {

// 	}
// 	return nil
// }

func (s *Catalog) Get(id uuid.UUID) (cat models.Cat, err error) {
	err = s.conn.QueryRow(context.Background(), "select * from cats where id=$1", id).Scan(&cat)
	if err != nil {
		return cat, err
	}

	return models.Cat{}, nil
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
			&cat,
		)
		if err != nil {
			return cats, err
		}

		cats = append(cats, cat)
	}
	return cats, nil
}