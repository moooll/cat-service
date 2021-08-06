package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"cat-service/db/psql"
	"cat-service/db/psql/models"

	"github.com/pquerna/ffjson/ffjson"
)

func TestAddCat(t *testing.T) {
	cat := psql.RandCat()
	var b *bytes.Buffer
	err := json.NewEncoder(b).Encode(cat)
	if err != nil {
		log.Print("cannot encode cat :/ ", err.Error())
	}

	req := httptest.NewRequest(http.MethodPost, "/cats", b)
	resp := req.Response

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Print("expected status: " +
			strconv.Itoa(http.StatusCreated) + " got status: " + strconv.Itoa(resp.StatusCode))
		t.Fail()
	}
}

func TestDeleteCat(t *testing.T) {
	err := psql.Truncate(db)
	if err != nil {
		log.Print("error truncating table cats ", err.Error())
	}

	cats, err := psql.SeedCats(db)
	if err != nil {
		log.Print("error seeding cats ", err.Error())
	}

	id := cats[rand.Intn(6)][0].(models.Cat).ID
	req := httptest.NewRequest(http.MethodDelete, "/cats/"+id.String(), nil)
	if req.Response.StatusCode != http.StatusOK {
		log.Print("testing delete handler failed :( ", err.Error())
		t.Fail()
	}
}

func TestUpdateCat(t *testing.T) {
	cat := psql.RandCat()
	var b *bytes.Buffer
	err := json.NewEncoder(b).Encode(cat)
	if err != nil {
		log.Print("cannot encode cat :/ ", err.Error())
	}

	req, _ := PrepareDBAndRequest(http.MethodPut, b)
	if req.Response.StatusCode != http.StatusOK {
		log.Print("testing delete handler failed :( ", err.Error())
		t.Fail()
	}

}

func TestGetCat(t *testing.T) {
	req, expectedCat := PrepareDBAndRequest(http.MethodGet, nil)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Response.Body)
	if err != nil {
		log.Print("error reading from buffer ", err.Error())
	}

	var gotCat models.Cat
	err = ffjson.Unmarshal(buf.Bytes(), &gotCat)
	if err != nil {
		log.Print("error unmarshalling cat ", err.Error())
	}

	if req.Response.StatusCode != http.StatusOK || expectedCat != gotCat {
		log.Print("testing get handler failed :( ", err.Error())
		t.Fail()
	}
}

func TestGetAllCats(t *testing.T) {
	err := psql.Truncate(db)
	if err != nil {
		log.Print("error truncating table cats ", err.Error())
	}

	expectedCats, err := psql.SeedCats(db)
	if err != nil {
		log.Print("error seeding cats ", err.Error())
	}

	req := httptest.NewRequest(http.MethodGet, "/cats", nil)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(req.Response.Body)
	if err != nil {
		log.Print("error reading from buffer ", err.Error())
	}

	var gotCats []models.Cat
	err = ffjson.Unmarshal(buf.Bytes(), &gotCats)
	if err != nil {
		log.Print("error unmarshalling cats ", err.Error())
	}

	for i, v := range expectedCats {
		if v[0].(models.Cat) != gotCats[i] {
			log.Print("expected cat: " + fmt.Sprint(v[0].(models.Cat)) + " got cat: " + fmt.Sprint(gotCats[i]))
		}
	}

}

func PrepareDBAndRequest(meth string, body io.Reader) (*http.Request, models.Cat) {
	err := psql.Truncate(db)
	if err != nil {
		log.Print("error truncating table cats ", err.Error())
	}

	cats, err := psql.SeedCats(db)
	if err != nil {
		log.Print("error seeding cats ", err.Error())
	}

	cat := cats[rand.Intn(6)][0].(models.Cat)
	id := cat.ID
	req := httptest.NewRequest(meth, "/cats/"+id.String(), body)
	return req, cat
}
