package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var db *sql.DB

func ConnectDB(host, port, username, password, dbname string) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	db.SetConnMaxLifetime(time.Nanosecond)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(1)

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}
