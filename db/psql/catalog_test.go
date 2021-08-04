package psql

import (
	"cat-service/db/psql/models"
	"context"
	"fmt"
	"math/rand"
	"testing"

	//"github.com/ory/dockertest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
)

func TestGetAll(t *testing.T) {
	if err := Truncate(db); err != nil {
		zap.L().Error("cannot truncate table cats")
	}

	expectedCats, err := SeedCats(db)
	if err != nil {
		zap.L().Error("cannot seed cats :/")
	}

	var gotCats []models.Cat
	err = db.conn.QueryRow(context.Background(), "select * from cats").Scan(&gotCats)
	if err != nil {
		zap.L().Error("cannot query the database: ", zap.Error(err))
	}

	for i, v := range expectedCats {
		if gotCats[i] != v[0].(models.Cat) {
			zap.L().Error("expected cat: " + v[0].(string) + "got cat: " + fmt.Sprint(gotCats[i]))
			t.Fail()
		}
	}
}

// todo
func TestGet(t *testing.T) {
	if err := Truncate(db); err != nil {
		zap.L().Error("cannot truncate table cats")
	}

	cats, err := SeedCats(db)
	if err != nil {
		zap.L().Error("cannot seed cats :/")
	}

	n := rand.Intn(6)
	gotCat, err := catalog.Get(cats[n][0].(models.Cat).ID)
	if err != nil {
		zap.L().Error("error getting items from db: ", zap.Error(err))
	}

	expectedCat := cats[n][0].(models.Cat)
	if gotCat != expectedCat {
		zap.L().Error("expected cat: " + fmt.Sprint(expectedCat) + "got cat: " + fmt.Sprint(gotCat))
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	if err := Truncate(db); err != nil {
		zap.L().Error("cannot truncate table cats")
	}

	cats, err := SeedCats(db)
	if err != nil {
		zap.L().Error("cannot seed cats :/")
	}

	n := rand.Intn(6)
	err = catalog.Delete(cats[n][0].(models.Cat).ID)
	if err != nil {
		zap.L().Error("cannot delete item: ", zap.Error(err))
	}

	gotCat := models.Cat{}
	err = db.conn.QueryRow(context.Background(), "select * from cats where id = $1", cats[n][0].(models.Cat).ID).Scan(&gotCat)
	if err != nil {
		zap.L().Info("could not get cat ", zap.Error(err))
	}

	if gotCat.Name != "" && gotCat.Breed != "" && gotCat.Color != "" {
		zap.L().Error("can't delete value: name " + gotCat.Name + " breed " + gotCat.Breed + " color " + gotCat.Color)
		t.Fail()
	}

}

func TestUpdate(t *testing.T) {
	if err := Truncate(db); err != nil {
		zap.L().Error("cannot truncate table cats")
	}

	cats, err := SeedCats(db)
	if err != nil {
		zap.L().Error("cannot seed cats :/")
	}

	name := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	breed := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	color := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	age := rand.Float32() * 15
	price := rand.Float32() * 15
	n := rand.Intn(6)
	expectedCat := models.Cat{
		ID:    cats[n][0].(models.Cat).ID,
		Name:  name,
		Breed: breed,
		Color: color,
		Age:   age,
		Price: price,
	}
	// err = db.conn.QueryRow(context.Background(), "update cats set name=$2, breed=$3, color=$4, age=$5, price=$6 where id $1",
	// 	cats[n][0].(models.Cat).ID, name, breed, color, age, price).Scan(&gotCats)

	gotCat, err := catalog.Update(cats[n][0].(models.Cat).ID, expectedCat)
	if err != nil {
		zap.L().Error("expected cats: " + fmt.Sprint(expectedCat) + " got cats: " + fmt.Sprint(gotCat))
	}

	if expectedCat != gotCat {
		zap.L().Error("expected cat: " + fmt.Sprint(expectedCat) + "got cat: " + fmt.Sprint(gotCat))
		t.Fail()
	}
}

func TestSave(t *testing.T) {
	if err := Truncate(db); err != nil {
		zap.L().Error("cannot truncate table cats")
	}

	expectedCat := randCat()
	gotCat := models.Cat{}
	// err := db.conn.QueryRow(context.Background(), "insert into cats(id, name, breed, color, age, price) values ($1, $2, $3, $4, $5, $6)",
		// expectedCat.ID, expectedCat.Name, expectedCat.Breed, expectedCat.Color, expectedCat.Age, expectedCat.Price).Scan(&gotCat)
	gotCat, err := catalog.Save(expectedCat)
	if err != nil {
		zap.L().Error("expected cat: " + fmt.Sprint(expectedCat) + "got cat: " + fmt.Sprint(gotCat))
		t.Fail()
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

func randCat() models.Cat {
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
