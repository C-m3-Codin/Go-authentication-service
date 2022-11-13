package main

import (
	"authentication/cmd/api/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64 = 0

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	conn := connectToDb()
	if conn == nil {
		log.Panic("Coulnt connect to DB")
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	log.Println("Starting Auth service ")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDb() *sql.DB {

	dsn := os.Getenv("DSN")
	fmt.Println(dsn)
	// dsn := "postgres://postgres:postgres@localhost:5431/"
	for {
		connection, err := openDb(dsn)
		if err != nil {
			log.Println("Postgres not up")
			counts++
		} else {
			log.Println("Connected to postgres")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("sleep between try")
		time.Sleep(2 * time.Second)
		continue
	}

}
