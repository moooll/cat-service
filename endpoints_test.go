package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"cat-service/db/psql"
	"cat-service/db/psql/models"

	"github.com/pquerna/ffjson/ffjson"
	"go.uber.org/zap"
)

func TestAddCat(t *testing.T) {
	cat := psql.RandCat()
	var b *bytes.Buffer
	err := json.NewEncoder(b).Encode(cat)
	if err != nil {
		zap.L().Error("cannot encode cat :/")
	}

	req := httptest.NewRequest(http.MethodPost, "/cats/add", b)
	resp := req.Response

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		zap.L().Error("expected status: " +
			strconv.Itoa(http.StatusCreated) + " got status: " + strconv.Itoa(resp.StatusCode))
		t.Fail()
	}
}

func TestDeleteCat(t *testing.T) {
	err := psql.Truncate(db)
	if err != nil {
		zap.L().Error("error truncating table cats")
	}

	cats, err := psql.SeedCats(db)
	if err != nil {
		zap.L().Error("error seeding cats")
	}

	id := cats[rand.Intn(6)][0].(models.Cat).ID
	req := httptest.NewRequest(http.MethodDelete, "/cats/"+id.String(), nil)
	if req.Response.StatusCode != http.StatusOK {
		zap.L().Error("testing delete handler failed :(")
		t.Fail()
	}
}

func TestUpdateCat(t *testing.T) {
	cat := psql.RandCat()
	var b *bytes.Buffer
	err := json.NewEncoder(b).Encode(cat)
	if err != nil {
		zap.L().Error("cannot encode cat :/")
	}

	req, _ := PrepareDBAndRequest(http.MethodPut, b)
	if req.Response.StatusCode != http.StatusOK {
		zap.L().Error("testing delete handler failed :(")
		t.Fail()
	}

}

func TestGetCat(t *testing.T) {
	req, expectedCat := PrepareDBAndRequest(http.MethodGet, nil)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Response.Body)
	if err != nil {
		zap.L().Error("error reading from buffer")
	}

	var gotCat models.Cat
	err = ffjson.Unmarshal(buf.Bytes(), &gotCat)
	if err != nil {
		zap.L().Error("error unmarshalling cat")
	}

	if req.Response.StatusCode != http.StatusOK || expectedCat != gotCat {
		zap.L().Error("testing get handler failed :(")
		t.Fail()
	}
}

func TestGetAllCats(t *testing.T) {
	err := psql.Truncate(db)
	if err != nil {
		zap.L().Error("error truncating table cats")
	}

	expectedCats, err := psql.SeedCats(db)
	if err != nil {
		zap.L().Error("error seeding cats")
	}

	req := httptest.NewRequest(http.MethodGet, "/cats", nil)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(req.Response.Body)
	if err != nil {
		zap.L().Error("error reading from buffer")
	}

	var gotCats []models.Cat
	err = ffjson.Unmarshal(buf.Bytes(), &gotCats)
	if err != nil {
		zap.L().Error("error unmarshalling cats")
	}

	for i, v := range expectedCats {
		if v[0].(models.Cat) != gotCats[i] {
			zap.L().Error("expected cat: " + fmt.Sprint(v[0].(models.Cat)) + " got cat: " + fmt.Sprint(gotCats[i]))
		}
	}

}

func PrepareDBAndRequest(meth string, body io.Reader) (*http.Request, models.Cat) {
	err := psql.Truncate(db)
	if err != nil {
		zap.L().Error("error truncating table cats")
	}

	cats, err := psql.SeedCats(db)
	if err != nil {
		zap.L().Error("error seeding cats")
	}

	cat := cats[rand.Intn(6)][0].(models.Cat)
	id := cat.ID
	req := httptest.NewRequest(meth, "/cats/"+id.String(), body)
	return req, cat
}
