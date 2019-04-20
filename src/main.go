package main

import (
	wine "WCMS/src/wine"
	"database/sql"
	"net/http"
	"os"
	"io"
	"log"
)

var db *sql.DB

func init() {
	logfile, err := os.OpenFile("wcms.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()

	wrt := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(wrt)

	db, err = sql.Open("postgres", "postgres://project:password@localhost/db_project?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	wine.DB = db

	fmt.Println("Successfully connected to database")
}

func main() {
	http.HandleFunc("/wine/", wine.WineHandler)
	http.ListenAndServe(":8080", nil)
}