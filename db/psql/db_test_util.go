package psql

import (
	"cat-service/db/psql/models"
	"context"
	"math/rand"

	"github.com/google/uuid"
	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v4"
	"github.com/thanhpk/randstr"
)

type TestConn struct {
	conn *pgx.Conn
}

func NewTestConn(conn *pgx.Conn) *TestConn {
	return &TestConn{
		conn,
	}
}


func Truncate(db *TestConn) error {
	_, err := db.conn.Query(context.Background(), "truncate table cats")
	if err != nil {
		return err
	}

	return nil
}

func SeedCats(db *TestConn) (rows [][]interface{}, err error) {
	rows = [][]interface{}{}
	for i := 0; i < 7; i++ {
		id := uuid.New()
		name := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		breed := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		color := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		age := rand.Float32() * 15
		price := rand.Float32() * 15
		rows[i] = append(rows[i], models.Cat{
			ID:    id,
			Name:  name,
			Breed: breed,
			Color: color,
			Age:   age,
			Price: price,
		})
	}
	_, err = db.conn.CopyFrom(context.Background(),
		pgx.Identifier{"cats"},
		[]string{"id", "name", "breed", "color", "age", "price"},
		pgx.CopyFromRows(rows))
	if err != nil {
		return rows, err
	}

	return rows, nil
}

func RandCat() models.Cat {
	id := uuid.New()
	name := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	breed := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	color := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	age := rand.Float32() * 15
	price := rand.Float32() * 15
	return models.Cat{
		ID:    id,
		Name:  name,
		Breed: breed,
		Color: color,
		Age:   age,
		Price: price,
	}
}
