package psql

import (
	"context"
	"log"
	"os"

	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v4"
)

func Connect() interface{} {
	conn, err := pgx.Connect(context.Background(), DatabaseURL)
	if err != nil {
		log.Print("could not connect to the psql ", err.Error())
		os.Exit(1)
	}

	log.Println("after connect ", err)
	err = conn.Ping(context.Background())
	if err != nil {
		log.Print("could not connect to the psql ", err.Error())
	}

	log.Println("after ping ", err)
	return conn
}

func Close(conn *pgx.Conn) error {
	return conn.Close(context.Background())
}
