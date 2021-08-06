package psql

import (
	"cat-service/db/psql/models"
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
)

func TestGetAll(t *testing.T) {
	if err := Truncate(db); err != nil {
		log.Print("cannot truncate table cats ", err.Error())
	}

	expectedCats, err := SeedCats(db)
	if err != nil {
		log.Print("cannot seed cats :/ ", err.Error())
	}

	gotCats, err := catalog.GetAll()
	if err != nil {
		log.Print("error gettinf all cats ", err.Error())
	}

	for i, v := range expectedCats {
		if gotCats[i] != v[0].(models.Cat) {
			log.Print("expected cat: " + v[0].(string) + "got cat: " + fmt.Sprint(gotCats[i]))
			t.Fail()
		}
	}
}

// todo
func TestGet(t *testing.T) {
	if err := Truncate(db); err != nil {
		log.Print("cannot truncate table cats ", err.Error())
	}

	cats, err := SeedCats(db)
	if err != nil {
		log.Print("cannot seed cats :/ ", err.Error())
	}

	n := rand.Intn(6)
	gotCat, err := catalog.Get(cats[n][0].(models.Cat).ID)
	if err != nil {
		log.Print("error getting items from db: ", err.Error())
	}

	expectedCat := cats[n][0].(models.Cat)
	if gotCat != expectedCat {
		log.Print("expected cat: " + fmt.Sprint(expectedCat) + "got cat: " + fmt.Sprint(gotCat))
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	if err := Truncate(db); err != nil {
		zap.L().Error("cannot truncate table cats ", zap.Error(err))
	}

	cats, err := SeedCats(db)
	if err != nil {
		zap.L().Error("cannot seed cats :/ ", zap.Error(err))
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
		zap.L().Error("cannot truncate table cats ", zap.Error(err))
	}
	
	cats, err := SeedCats(db)
	if err != nil {
		zap.L().Error("cannot seed cats :/ ", zap.Error(err))
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
	err = catalog.Update(cats[n][0].(models.Cat).ID, expectedCat)
	if err != nil {
		log.Print("error: ", err.Error())
	}
}

// doubt if the test gonna pass
func TestSave(t *testing.T) {
	if err := Truncate(db); err != nil {
		zap.L().Error("cannot truncate table cats ", zap.Error(err))
	}

	expectedCat := RandCat()
	err := catalog.Save(expectedCat)
	if err != nil {
		log.Print("err ", err)
		t.Fail()
	}
}
